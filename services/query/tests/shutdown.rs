//! A shutdown finishes the work already in hand (roadmap R-60).
//!
//! Elysium had no shutdown handling at all, so a restart in the middle
//! of a write could split an audit pair or a commit. The rule here is
//! that a service stops accepting work, finishes what it holds, and only
//! then exits. This test holds one call open across a shutdown and
//! requires it to be answered.

use std::sync::Arc;
use std::time::Duration;

use query::build::build_service_client::BuildServiceClient;
use query::build::build_service_server::BuildService;
use query::build::{GetBuildInfoRequest, GetBuildInfoResponse};
use query::libs::rpc::{Request, Response, Status};
use query::libs::runtime;
use query::service::serve_with;

/// A service whose one call takes long enough to still be in hand when
/// the shutdown is asked for.
struct Slow {
    holding: Arc<runtime::sync::Notify>,
}

#[query::libs::rpc::async_trait]
impl BuildService for Slow {
    async fn get_build_info(
        &self,
        _request: Request<GetBuildInfoRequest>,
    ) -> Result<Response<GetBuildInfoResponse>, Status> {
        // Tell the test the call is in hand, then take a while.
        self.holding.notify_one();
        runtime::time::sleep(Duration::from_millis(300)).await;
        Ok(Response::new(GetBuildInfoResponse { builds: Vec::new() }))
    }
}

#[test]
fn a_shutdown_finishes_the_call_already_in_hand() {
    let runtime = runtime::runtime::Builder::new_multi_thread()
        .enable_all()
        .build()
        .expect("a runtime starts");
    runtime.block_on(async {
        let listener = runtime::net::TcpListener::bind("127.0.0.1:0")
            .await
            .expect("a port is free");
        let address = listener.local_addr().expect("the port is known");
        let holding = Arc::new(runtime::sync::Notify::new());
        let (stop, stopped) = runtime::sync::oneshot::channel::<()>();

        let serving = runtime::spawn(serve_with(
            listener,
            async move {
                let _ = stopped.await;
            },
            Slow {
                holding: Arc::clone(&holding),
            },
        ));

        let mut client = BuildServiceClient::connect(format!("http://{address}"))
            .await
            .expect("the service answers");
        let asking =
            runtime::spawn(async move { client.get_build_info(GetBuildInfoRequest {}).await });

        // Ask it to stop while that call is still in hand.
        holding.notified().await;
        let asked_to_stop = std::time::Instant::now();
        drop(stop);

        serving
            .await
            .expect("the server ran")
            .expect("the server stopped cleanly");
        // Stopping waits for the work in hand: were the service to
        // return at once, the process would exit from under that call.
        // The call takes 300ms, so returning sooner means abandoning it.
        let waited = asked_to_stop.elapsed();
        assert!(
            waited >= Duration::from_millis(250),
            "the service stopped after {waited:?}, abandoning a call still in hand"
        );
        let answered = asking.await.expect("the call ran");
        assert!(
            answered.is_ok(),
            "the shutdown cut off a call already in hand: {answered:?}"
        );

        // Once stopped, it takes no new work.
        assert!(
            BuildServiceClient::connect(format!("http://{address}"))
                .await
                .is_err(),
            "the service still accepts calls after stopping"
        );
    });
}
