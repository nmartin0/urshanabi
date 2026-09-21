//! This build's identity, recorded by the build script (roadmap R-121).

use crate::libs::wire::Timestamp;
use crate::proto::build::BuildInfo;

/// The component's directory name.
pub const COMPONENT: &str = "query";

/// The environment variable a deployment sets to the published SHA-256
/// of the running artifact. A binary cannot contain its own digest.
pub const ARTIFACT_SHA256: &str = "URSHANABI_ARTIFACT_SHA256";

/// The identity of the running build. Every value but the artifact
/// digest is fixed when the source is built; the digest is read from
/// the environment, and left empty unless it is a well-formed SHA-256.
///
/// # Panics
///
/// Never in practice: the build script refuses to build without a
/// positive integer commit time.
#[must_use]
pub fn this_build() -> BuildInfo {
    let committed: i64 = env!("URSHANABI_SOURCE_COMMITTED_AT")
        .parse()
        .expect("the build script records a positive integer");
    BuildInfo {
        component: COMPONENT.to_owned(),
        version: env!("URSHANABI_VERSION").to_owned(),
        source_revision: env!("URSHANABI_SOURCE_REVISION").to_owned(),
        source_committed_at: Some(Timestamp {
            seconds: committed,
            nanos: 0,
        }),
        artifact_sha256: artifact_sha256(std::env::var(ARTIFACT_SHA256).ok().as_deref()),
    }
}

/// Accepts a digest only if it is 64 lowercase hexadecimal characters.
fn artifact_sha256(value: Option<&str>) -> String {
    match value {
        Some(d) if d.len() == 64 && d.bytes().all(|b| matches!(b, b'0'..=b'9' | b'a'..=b'f')) => {
            d.to_owned()
        }
        _ => String::new(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn a_well_formed_digest_is_kept() {
        let d = "a".repeat(64);
        assert_eq!(artifact_sha256(Some(&d)), d);
    }

    #[test]
    fn a_missing_or_malformed_digest_is_empty() {
        assert_eq!(artifact_sha256(None), "");
        assert_eq!(artifact_sha256(Some("not-a-digest")), "");
        assert_eq!(artifact_sha256(Some(&"A".repeat(64))), "");
    }

    #[test]
    fn the_identity_describes_the_source() {
        let build = this_build();
        assert_eq!(build.component, "query");
        assert_eq!(build.source_revision.len(), 40);
        assert!(
            build
                .source_committed_at
                .is_some_and(|t| t.seconds > 0 && t.nanos == 0)
        );
    }
}
