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
/// The modules mirror the contracts' own packages, because generated
/// code refers to other packages by that structure. The workspace lints
/// apply to our code, not to generated code, so they are relaxed here.
#[allow(missing_docs, clippy::pedantic)] // name-ok
pub mod proto {
    /// Urshanabi's own contracts.
    pub mod urshanabi {
        /// The build-identity contract.
        pub mod build {
            /// Its first version.
            pub mod v1 {
                crate::libs::rpc::include_proto!("urshanabi.build.v1");
            }
        }
    }

    /// The contract language's own well-known types, generated with
    /// ours so that every message can be written as JSON (R-78).
    pub mod google {
        /// Its types.
        pub mod protobuf {
            crate::libs::rpc::include_proto!("google.protobuf"); // name-ok
        }
    }
}

/// The build-identity contract, by its short name.
pub use proto::urshanabi::build::v1 as build;

/// The contract language's well-known types, by their short name.
pub use proto::google::protobuf as wire; // name-ok
