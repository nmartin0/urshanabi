//! The producer answers every expectation its consumers have recorded
//! (roadmap R-78). Each file in `contracts/expectations/` says what one
//! consumer relies on; this test calls the real service and checks the
//! answer against every one of them, so a change that breaks a consumer
//! fails this service's own build.

use std::fs;
use std::path::Path;

use query::build::GetBuildInfoRequest;
use query::build::build_service_client::BuildServiceClient;
use query::libs::runtime;
use query::service::serve;
use serde_json::Value; // name-ok

/// Reads the value at a path such as `builds.0.component`.
fn at<'a>(answer: &'a Value, path: &str) -> Option<&'a Value> {
    let mut here = answer;
    for step in path.split('.') {
        here = match step.parse::<usize>() {
            Ok(index) => here.get(index)?,
            Err(_) => here.get(step)?,
        };
    }
    Some(here)
}

/// Checks one recorded expectation against the answer, returning what
/// is wrong when it does not hold.
fn unmet(answer: &Value, expectation: &Value) -> Option<String> {
    let path = expectation["path"].as_str()?;
    let Some(found) = at(answer, path) else {
        return Some(format!("{path} is absent from the answer"));
    };
    if let Some(wanted) = expectation.get("equals").and_then(Value::as_str) {
        let got = found.as_str().unwrap_or_default();
        if got != wanted {
            return Some(format!("{path} is {got:?}, and {wanted:?} was expected"));
        }
    }
    if let Some(pattern) = expectation.get("matches").and_then(Value::as_str) {
        let got = found.as_str().unwrap_or_default();
        if !matches(got, pattern) {
            return Some(format!("{path} is {got:?}, which does not match {pattern}"));
        }
    }
    if expectation.get("present").and_then(Value::as_bool) == Some(true) && found.is_null() {
        return Some(format!("{path} is present but empty"));
    }
    None
}

/// The patterns recorded so far are alternatives of literal text,
/// character classes and counts; this understands exactly those.
fn matches(text: &str, pattern: &str) -> bool {
    pattern.split('|').any(|one| matches_one(text, one))
}

fn matches_one(text: &str, pattern: &str) -> bool {
    let pattern = pattern.trim_start_matches('^').trim_end_matches('$');
    if let Some((class, count)) = hex_run(pattern) {
        return text.len() == count && text.chars().all(|c| class.contains(c));
    }
    if pattern == r"development(\+modified)?" {
        return text == "development" || text == "development+modified";
    }
    if pattern == r"[0-9]+\.[0-9]+\.[0-9]+" {
        let parts: Vec<&str> = text.split('.').collect();
        return parts.len() == 3
            && parts
                .iter()
                .all(|p| !p.is_empty() && p.chars().all(|c| c.is_ascii_digit()));
    }
    panic!("the expectation uses a pattern this test does not understand: {pattern}");
}

/// Reads a pattern of the form `[0-9a-f]{40}`.
fn hex_run(pattern: &str) -> Option<(&'static str, usize)> {
    let rest = pattern.strip_prefix("[0-9a-f]{")?;
    let count = rest.strip_suffix('}')?.parse().ok()?;
    Some(("0123456789abcdef", count))
}

#[test]
fn every_consumer_expectation_holds() {
    let root = Path::new(env!("CARGO_MANIFEST_DIR")).join("../..");
    let recorded = root.join("contracts/expectations");
    let mut files: Vec<_> = fs::read_dir(&recorded)
        .expect("the recorded expectations are readable")
        .map(|entry| entry.expect("an entry is readable").path())
        .collect();
    files.sort();
    assert!(
        !files.is_empty(),
        "no consumer has recorded any expectation"
    );

    let runtime = runtime::runtime::Builder::new_current_thread()
        .enable_all()
        .build()
        .expect("a runtime starts");
    let mut checked = 0;
    for file in files {
        let recorded: Value =
            serde_json::from_str(&fs::read_to_string(&file).expect("the file is readable")) // name-ok
                .expect("the file is well-formed");
        if recorded["producer"] != "query" {
            continue;
        }
        let consumer = recorded["consumer"].as_str().unwrap_or("a consumer");
        for interaction in recorded["interactions"]
            .as_array()
            .expect("interactions are listed")
        {
            let answer = runtime.block_on(ask());
            for expectation in interaction["expect"]
                .as_array()
                .expect("expectations are listed")
            {
                checked += 1;
                if let Some(wrong) = unmet(&answer, expectation) {
                    panic!("{consumer} relies on something this service no longer does: {wrong}");
                }
            }
        }
    }
    assert!(checked > 0, "no expectation of this service was checked");
}

/// Starts the real service and asks it for its build identity.
async fn ask() -> Value {
    let listener = runtime::net::TcpListener::bind("127.0.0.1:0")
        .await
        .expect("a port is free");
    let address = listener.local_addr().expect("the port is known");
    let (stop, stopped) = runtime::sync::oneshot::channel::<()>();
    runtime::spawn(async move {
        let _ = serve(listener, async move {
            let _ = stopped.await;
        })
        .await;
    });
    let mut client = BuildServiceClient::connect(format!("http://{address}"))
        .await
        .expect("the service answers");
    let answer = client
        .get_build_info(GetBuildInfoRequest {})
        .await
        .expect("the call succeeds")
        .into_inner();
    let written = serde_json::to_value(answer).expect("the answer can be written as JSON"); // name-ok
    drop(stop);
    written
}
