//! Build script: generates the service's contract code from `contracts/`
//! at build time, so nothing generated is committed (roadmap R-11), and
//! records the source identity the service reports (R-121).
//!
//! The identity describes the source, never the build moment, so two
//! builds of the same commit report the same identity. Each value comes
//! from the environment when set, as a release build does, and from the
//! source repository otherwise.

use std::env;
use std::fs;
use std::path::{Path, PathBuf};
use std::process::Command;

use prost::Message; // name-ok
use prost_types::FileDescriptorSet; // name-ok
use tonic_prost_build as generator; // name-ok

/// The pinned contract toolchain, installed by `script/bootstrap`.
const CONTRACT_TOOL: &str = ".tools/bin/buf"; // name-ok

/// The source repository's command-line tool, and its option that
/// resolves a path inside the repository's own storage.
const REPOSITORY_TOOL: &str = "git"; // name-ok
const STORAGE_PATH: &str = "--git-path"; // name-ok

/// The prefix of an instruction to the build tool.
const BUILD_TOOL: &str = "cargo"; // name-ok

fn main() {
    let manifest = PathBuf::from(env::var("CARGO_MANIFEST_DIR").expect("set by the builder"));
    let root = manifest
        .parent()
        .and_then(Path::parent)
        .expect("the service lives two levels below the repository root")
        .to_path_buf();
    generate_contracts(&root);
    record_identity(&root);
}

/// Compiles the contracts to a descriptor set with the pinned contract
/// toolchain, then generates the service's code from it.
fn generate_contracts(root: &Path) {
    let contracts = root.join("contracts");
    let tool = env::var_os("URSHANABI_CONTRACT_TOOL")
        .map_or_else(|| root.join(CONTRACT_TOOL), PathBuf::from);
    let out = PathBuf::from(env::var("OUT_DIR").expect("set by the builder"));
    let descriptors = out.join("contracts.binpb");

    let status = Command::new(&tool)
        .arg("build")
        .arg(&contracts)
        .arg("--as-file-descriptor-set")
        .arg("-o")
        .arg(&descriptors)
        .status()
        .unwrap_or_else(|e| panic!("cannot run {}: {e}; run script/bootstrap", tool.display()));
    assert!(status.success(), "the contracts failed to compile");

    let bytes = fs::read(&descriptors).expect("the descriptor set was just written");
    let set = FileDescriptorSet::decode(bytes.as_slice())
        .expect("the contract toolchain writes a valid descriptor set");
    generator::configure()
        .compile_fds(set)
        .expect("code generation from the contracts");

    directive("rerun-if-changed", &contracts.display().to_string());
    directive("rerun-if-env-changed", "URSHANABI_CONTRACT_TOOL");
}

/// Records the version, source revision and commit time as compile-time
/// environment variables.
fn record_identity(root: &Path) {
    let version = env::var("URSHANABI_VERSION").unwrap_or_else(|_| "development".to_owned());

    let revision = env::var("URSHANABI_SOURCE_REVISION")
        .unwrap_or_else(|_| repository(root, &["rev-parse", "HEAD"]));
    assert!(
        revision.len() == 40 && revision.bytes().all(|b| b.is_ascii_hexdigit()),
        "the source revision must be a full 40-character hash, not {revision:?}"
    );

    // SOURCE_DATE_EPOCH is the reproducible-builds standard for this value.
    let committed = env::var("SOURCE_DATE_EPOCH")
        .unwrap_or_else(|_| repository(root, &["log", "-1", "--format=%ct"]));
    assert!(
        committed.parse::<i64>().is_ok_and(|s| s > 0),
        "the commit time must be positive seconds since the epoch, not {committed:?}"
    );

    directive("rustc-env", &format!("URSHANABI_VERSION={version}"));
    directive(
        "rustc-env",
        &format!("URSHANABI_SOURCE_REVISION={revision}"),
    );
    directive(
        "rustc-env",
        &format!("URSHANABI_SOURCE_COMMITTED_AT={committed}"),
    );
    for var in [
        "URSHANABI_VERSION",
        "URSHANABI_SOURCE_REVISION",
        "SOURCE_DATE_EPOCH",
    ] {
        directive("rerun-if-env-changed", var);
    }

    // Rebuild when the checked-out commit changes.
    let head = repository(root, &["rev-parse", STORAGE_PATH, "HEAD"]);
    directive("rerun-if-changed", &root.join(head).display().to_string());
    if let Some(reference) = try_repository(root, &["symbolic-ref", "-q", "HEAD"]) {
        let path = repository(root, &["rev-parse", STORAGE_PATH, &reference]);
        directive("rerun-if-changed", &root.join(path).display().to_string());
    }
}

/// Runs a source-repository command, panicking with the reason if it fails.
fn repository(root: &Path, args: &[&str]) -> String {
    try_repository(root, args).unwrap_or_else(|| {
        panic!(
            "cannot read the source identity ({}); outside a repository, \
             set URSHANABI_SOURCE_REVISION and SOURCE_DATE_EPOCH",
            args.join(" ")
        )
    })
}

/// Runs a source-repository command, returning its trimmed output.
fn try_repository(root: &Path, args: &[&str]) -> Option<String> {
    let output = Command::new(REPOSITORY_TOOL)
        .args(args)
        .current_dir(root)
        .output()
        .ok()?;
    output
        .status
        .success()
        .then(|| String::from_utf8_lossy(&output.stdout).trim().to_owned())
}

/// Sends one instruction to the build tool.
fn directive(kind: &str, value: &str) {
    println!("{BUILD_TOOL}:{kind}={value}");
}
