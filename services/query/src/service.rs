//! The service's contract implementation and its server.

use std::future::Future;

use crate::build::build_service_server::{BuildService, BuildServiceServer};
use crate::build::{BuildInfo, GetBuildInfoRequest, GetBuildInfoResponse};
use crate::identity;
use crate::libs::{log, rpc, runtime};

/// The metadata key carrying the request id (roadmap R-08).
pub const REQUEST_ID: &str = "x-request-id";

/// Answers the build-identity contract with this build's identity.
#[derive(Debug, Clone)]
pub struct Service {
    builds: Vec<BuildInfo>,
}

impl Service {
    /// A service reporting the given build as its own.
    #[must_use]
    pub fn new(this: BuildInfo) -> Self {
        Self { builds: vec![this] }
    }
}

#[rpc::async_trait]
impl BuildService for Service {
    async fn get_build_info(
        &self,
        request: rpc::Request<GetBuildInfoRequest>,
    ) -> Result<rpc::Response<GetBuildInfoResponse>, rpc::Status> {
        let id = request
            .metadata()
            .get(REQUEST_ID)
            .map(|v| v.to_str().unwrap_or(""));
        log::info!(request_id = loggable(id), "GetBuildInfo");
        Ok(rpc::Response::new(GetBuildInfoResponse {
            builds: self.builds.clone(),
        }))
    }
}

/// Returns a request id safe to write to a log. The id comes from the
/// caller, so anything but a plain identifier could forge log lines.
fn loggable(id: Option<&str>) -> &str {
    match id {
        None => "none",
        Some(v)
            if (1..=128).contains(&v.len())
                && v.bytes()
                    .all(|b| b.is_ascii_alphanumeric() || matches!(b, b'-' | b'_' | b'.')) =>
        {
            v
        }
        Some(_) => "invalid",
    }
}

/// Serves the build-identity contract on `listener` until `shutdown`
/// completes.
///
/// # Errors
///
/// Returns the transport's error if the server fails.
pub async fn serve(
    listener: runtime::net::TcpListener,
    shutdown: impl Future<Output = ()>,
) -> Result<(), rpc::transport::Error> {
    serve_with(listener, shutdown, Service::new(identity::this_build())).await
}

/// Serves `answers` on `listener` until `shutdown` completes, then
/// finishes the calls already in flight before returning (roadmap
/// R-60). Taking the answering service as an argument is what lets a
/// test hold a call open across a shutdown.
///
/// # Errors
///
/// Returns the transport's error if the server fails.
pub async fn serve_with<S: BuildService>(
    listener: runtime::net::TcpListener,
    shutdown: impl Future<Output = ()>,
    answers: S,
) -> Result<(), rpc::transport::Error> {
    rpc::transport::Server::builder()
        .add_service(BuildServiceServer::new(answers))
        .serve_with_incoming_shutdown(
            rpc::transport::server::TcpIncoming::from(listener),
            shutdown,
        )
        .await
}

#[cfg(test)]
mod tests {
    use super::loggable;

    #[test]
    fn a_plain_request_id_is_logged_as_given() {
        assert_eq!(loggable(Some("req-42_a.b")), "req-42_a.b");
    }

    #[test]
    fn a_missing_request_id_is_logged_as_none() {
        assert_eq!(loggable(None), "none");
    }

    #[test]
    fn a_request_id_that_could_forge_log_lines_is_refused() {
        assert_eq!(loggable(Some("x\nforged line")), "invalid");
        assert_eq!(loggable(Some("")), "invalid");
        assert_eq!(loggable(Some(&"a".repeat(129))), "invalid");
    }
}
