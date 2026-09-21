//! The query service answers every question from people and agents,
//! enforcing policy on every answer (see `README.md`). In the walking
//! skeleton it answers one question: which build is running (roadmap
//! R-121).

pub mod identity;
pub mod libs;
pub mod logging;
pub mod service;

/// Code generated from `contracts/` at build time (roadmap R-11).
///
/// The workspace lints apply to our code, not to generated code, so
/// they are relaxed for this module alone.
#[allow(missing_docs, clippy::pedantic)] // name-ok
pub mod proto {
    /// The build-identity contract, `urshanabi.build.v1`.
    pub mod build {
        crate::libs::rpc::include_proto!("urshanabi.build.v1");
    }
}
