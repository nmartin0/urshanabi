//! The libraries the service is built on, each named once and used by
//! its role everywhere else (`RULES.md` H1).

pub use tokio as runtime; // name-ok
pub use tonic as rpc; // name-ok
pub use tracing as log;
pub use tracing_subscriber as log_output;
