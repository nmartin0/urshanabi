//! The query service's entry point.

use std::process::ExitCode;

use query::libs::{log, runtime};

/// The environment variable naming the address to listen on.
const LISTEN: &str = "URSHANABI_LISTEN";

fn main() -> ExitCode {
    if let Err(e) = query::logging::init() {
        eprintln!("cannot start logging: {e}");
        return ExitCode::FAILURE;
    }
    let address = std::env::var(LISTEN).unwrap_or_else(|_| "127.0.0.1:50051".to_owned());
    let result = runtime::runtime::Builder::new_multi_thread()
        .enable_all()
        .build()
        .map_err(|e| e.to_string())
        .and_then(|rt| rt.block_on(run(&address)));
    match result {
        Ok(()) => ExitCode::SUCCESS,
        Err(e) => {
            log::error!(error = %e, "the query service stopped");
            ExitCode::FAILURE
        }
    }
}

async fn run(address: &str) -> Result<(), String> {
    let listener = runtime::net::TcpListener::bind(address)
        .await
        .map_err(|e| format!("cannot listen on {address}: {e}"))?;
    let local = listener.local_addr().map_err(|e| e.to_string())?;
    log::info!(address = %local, "listening");
    query::service::serve(listener, shutdown())
        .await
        .map_err(|e| e.to_string())
}

/// Completes on an interrupt or a termination request.
async fn shutdown() {
    let interrupt = runtime::signal::ctrl_c();
    let mut terminate =
        match runtime::signal::unix::signal(runtime::signal::unix::SignalKind::terminate()) {
            Ok(signal) => signal,
            Err(e) => {
                log::error!(error = %e, "cannot watch for termination requests");
                let _ = interrupt.await;
                return;
            }
        };
    runtime::select! {
        _ = interrupt => {}
        _ = terminate.recv() => {}
    }
    log::info!("shutting down");
}
