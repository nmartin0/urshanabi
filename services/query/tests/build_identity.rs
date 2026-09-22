//! Integration: the real server, called over the network as the gateway
//! calls it (walking skeleton, roadmap R-121).

use std::future::Future;

use query::libs::{rpc, runtime};
use query::proto::build::GetBuildInfoRequest;
use query::proto::build::build_service_client::BuildServiceClient;
use query::service::{self, REQUEST_ID};

/// Runs `f` on a fresh runtime.
fn run<F: Future>(f: F) -> F::Output {
    runtime::runtime::Builder::new_current_thread()
        .enable_all()
        .build()
        .expect("a runtime")
        .block_on(f)
}

/// Starts the server on a port the operating system chooses, and
/// returns its address and a handle that stops it.
async fn start() -> (String, runtime::sync::oneshot::Sender<()>) {
    let listener = runtime::net::TcpListener::bind("127.0.0.1:0")
        .await
        .expect("a free port");
    let address = format!("http://{}", listener.local_addr().expect("a bound address"));
    let (stop, stopped) = runtime::sync::oneshot::channel::<()>();
    runtime::spawn(service::serve(listener, async move {
        let _ = stopped.await;
    }));
    (address, stop)
}

#[test]
fn the_server_reports_its_own_build_identity() {
    run(async {
        let (address, stop) = start().await;
        let mut client = BuildServiceClient::connect(address)
            .await
            .expect("a connection");
        let mut request = rpc::Request::new(GetBuildInfoRequest {});
        request
            .metadata_mut()
            .insert(REQUEST_ID, "test-1".parse().expect("valid metadata"));

        let builds = client
            .get_build_info(request)
            .await
            .expect("an answer")
            .into_inner()
            .builds;

        assert_eq!(builds.len(), 1, "the query service reports only itself");
        let build = &builds[0];
        assert_eq!(build.component, "query");
        assert!(
            matches!(
                build.version.as_str(),
                "development" | "development+modified"
            ),
            "unexpected version {:?}",
            build.version
        );
        assert!(build.source_revision.bytes().all(|b| b.is_ascii_hexdigit()));
        assert_eq!(build.source_revision.len(), 40);
        let committed = build.source_committed_at.expect("a commit time");
        assert!(
            committed.seconds > 0 && committed.nanos == 0,
            "whole seconds, UTC"
        );
        let _ = stop.send(());
    });
}

#[test]
fn every_answer_reports_the_same_identity() {
    run(async {
        let (address, stop) = start().await;
        let mut client = BuildServiceClient::connect(address)
            .await
            .expect("a connection");
        let first = client
            .get_build_info(GetBuildInfoRequest {})
            .await
            .expect("an answer");
        let second = client
            .get_build_info(GetBuildInfoRequest {})
            .await
            .expect("an answer");
        assert_eq!(first.into_inner(), second.into_inner());
        let _ = stop.send(());
    });
}
