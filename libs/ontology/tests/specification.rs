//! Every TOML example in the format's specification parses with this
//! library, so the specification and the code cannot drift apart
//! (roadmap R-15).

use std::fs;
use std::path::Path;

use ontology::model::{ActionType, Catalogue, Header, LinkType, Metric, ObjectType};

#[test]
fn every_example_in_the_specification_parses() {
    let spec = Path::new(env!("CARGO_MANIFEST_DIR")).join("../../docs/ontology-format.md");
    let text = fs::read_to_string(&spec).expect("the specification");
    let mut checked = 0;
    for block in text.split("```toml\n").skip(1) {
        let body = block.split("```").next().expect("a closed block");
        let Some(path) = body
            .lines()
            .next()
            .and_then(|l| l.strip_prefix("# ontology/"))
        else {
            continue;
        };
        let result = if path == "ontology.toml" {
            ontology::parse::<Header>(body).map(drop)
        } else if path.starts_with("types/") {
            ontology::parse::<ObjectType>(body).map(drop)
        } else if path.starts_with("links/") {
            ontology::parse::<LinkType>(body).map(drop)
        } else if path.starts_with("actions/") {
            ontology::parse::<ActionType>(body).map(drop)
        } else if path.starts_with("metrics/") {
            ontology::parse::<Metric>(body).map(drop)
        } else if path.starts_with("languages/") {
            ontology::parse::<Catalogue>(body).map(drop)
        } else {
            panic!("the specification has an example for an unknown file: {path}");
        };
        if let Err(e) = result {
            panic!("the example for {path} does not parse: {e}");
        }
        checked += 1;
    }
    assert!(
        checked >= 6,
        "only {checked} examples were found; the extraction is broken"
    );
}
