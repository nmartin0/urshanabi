//! The service's log output: one line per event, with timestamps in UTC
//! (roadmap R-128).

use crate::libs::log;
use crate::libs::log_output::fmt::MakeWriter;

/// The service's log subscriber, writing to `writer`.
pub fn subscriber<W>(writer: W) -> impl log::Subscriber + Send + Sync
where
    W: for<'w> MakeWriter<'w> + Send + Sync + 'static,
{
    crate::libs::log_output::fmt()
        .with_target(false)
        .with_writer(writer)
        .finish()
}

/// Sends the service's logs to standard error.
///
/// # Errors
///
/// Fails if a global subscriber is already set.
pub fn init() -> Result<(), log::subscriber::SetGlobalDefaultError> {
    log::subscriber::set_global_default(subscriber(std::io::stderr))
}

#[cfg(test)]
mod tests {
    use std::io::Write;
    use std::sync::{Arc, Mutex};
    use std::time::{SystemTime, UNIX_EPOCH};

    use crate::libs::log;

    /// A log destination the test can read back.
    #[derive(Clone, Default)]
    struct Captured(Arc<Mutex<Vec<u8>>>);

    impl Write for Captured {
        fn write(&mut self, bytes: &[u8]) -> std::io::Result<usize> {
            self.0
                .lock()
                .expect("an unpoisoned lock")
                .extend_from_slice(bytes);
            Ok(bytes.len())
        }
        fn flush(&mut self) -> std::io::Result<()> {
            Ok(())
        }
    }

    /// The current UTC hour, from the system clock alone.
    fn utc_hour() -> u64 {
        let seconds = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .expect("a clock after 1970")
            .as_secs();
        (seconds / 3600) % 24
    }

    #[test]
    fn log_timestamps_are_utc_whatever_the_host_zone() {
        let captured = Captured::default();
        let writer = captured.clone();
        let before = utc_hour();
        log::subscriber::with_default(super::subscriber(move || writer.clone()), || {
            log::info!("probe");
        });
        let after = utc_hour();

        let bytes = captured.0.lock().expect("an unpoisoned lock").clone();
        let line = String::from_utf8(bytes).expect("text");
        let stamp = line.split_whitespace().next().expect("a timestamp");
        assert!(stamp.ends_with('Z'), "not marked UTC: {line:?}");
        let hour: u64 = stamp[11..13].parse().expect("an hour");
        assert!(
            hour == before || hour == after,
            "logged hour {hour} is not the UTC hour ({before} to {after}): {line:?}"
        );
    }
}
