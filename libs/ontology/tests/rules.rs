//! Integration: a complete ontology loads, and each rule of the format
//! (`docs/ontology-format.md`) refuses an ontology that breaks it.

use std::fs;
use std::path::{Path, PathBuf};
use std::sync::atomic::{AtomicUsize, Ordering};

/// The complete, valid fixture ontology.
fn fixture() -> PathBuf {
    Path::new(env!("CARGO_MANIFEST_DIR")).join("tests/fixtures/sales")
}

/// A private copy of the fixture, for one test to break.
fn copy() -> PathBuf {
    static NEXT: AtomicUsize = AtomicUsize::new(0);
    let dir = std::env::temp_dir().join(format!(
        "ontology-test-{}-{}",
        std::process::id(),
        NEXT.fetch_add(1, Ordering::Relaxed)
    ));
    let _ = fs::remove_dir_all(&dir);
    copy_dir(&fixture(), &dir);
    dir
}

fn copy_dir(from: &Path, to: &Path) {
    fs::create_dir_all(to).expect("a writable temporary directory");
    for entry in fs::read_dir(from).expect("the fixture") {
        let path = entry.expect("an entry").path();
        let target = to.join(path.file_name().expect("a name"));
        if path.is_dir() {
            copy_dir(&path, &target);
        } else {
            fs::copy(&path, &target).expect("a copy");
        }
    }
}

/// Replaces `from` with `to` in one file of the copy; `from` must occur.
fn edit(dir: &Path, file: &str, from: &str, to: &str) {
    let path = dir.join(file);
    let text = fs::read_to_string(&path).expect("the file");
    assert!(text.contains(from), "{file} does not contain {from:?}");
    fs::write(&path, text.replacen(from, to, 1)).expect("a write");
}

/// Loads `dir`, requiring refusal; returns every problem as text.
fn refused(dir: &Path) -> String {
    let problems = ontology::load(dir).expect_err("the ontology should be refused");
    let _ = fs::remove_dir_all(dir);
    problems
        .iter()
        .map(ToString::to_string)
        .collect::<Vec<_>>()
        .join("\n")
}

fn assert_refused(dir: &Path, expected: &str) {
    let found = refused(dir);
    assert!(
        found.contains(expected),
        "expected {expected:?} among:\n{found}"
    );
}

#[test]
fn the_complete_ontology_loads() {
    let o = ontology::load(&fixture()).unwrap_or_else(|p| panic!("refused: {p:?}"));
    assert_eq!(o.header.languages, ["en", "es", "pt"]);
    assert_eq!(
        (
            o.types.len(),
            o.links.len(),
            o.actions.len(),
            o.metrics.len()
        ),
        (2, 1, 1, 1)
    );
    assert_eq!(o.catalogues["es"]["Customer"].name, "Cliente");
}

#[test]
fn an_unknown_field_is_refused_not_ignored() {
    let d = copy();
    edit(
        &d,
        "types/order.toml",
        "freshness   = \"real-time\"",
        "freshnes = \"real-time\"",
    );
    assert_refused(&d, "unknown field");
}

#[test]
fn rule_1_names_follow_their_case() {
    let d = copy();
    edit(
        &d,
        "types/order.toml",
        "name        = \"Order\"",
        "name = \"order\"",
    );
    assert_refused(&d, "must be UpperCamelCase");
}

#[test]
fn rule_1_names_are_unique_within_their_kind() {
    let d = copy();
    fs::copy(d.join("types/order.toml"), d.join("types/order2.toml")).expect("a copy");
    assert_refused(&d, "two object types share this name");
}

#[test]
fn rule_2_the_source_is_a_curated_table() {
    let d = copy();
    edit(
        &d,
        "types/order.toml",
        "curated.sales.orders",
        "raw.sales.orders",
    );
    assert_refused(&d, "must name a curated-layer table");
}

#[test]
fn rule_2_the_primary_key_is_a_property() {
    let d = copy();
    edit(
        &d,
        "types/order.toml",
        "primary_key = \"order_id\"",
        "primary_key = \"order_number\"",
    );
    assert_refused(
        &d,
        "the primary key \"order_number\" is not one of its properties",
    );
}

#[test]
fn rule_3_property_types_are_known() {
    let d = copy();
    edit(
        &d,
        "types/order.toml",
        "type = \"decimal\"",
        "type = \"money\"",
    );
    assert_refused(&d, "unknown type \"money\"");
}

#[test]
fn rule_4_an_editable_property_declares_reconcile() {
    let d = copy();
    edit(
        &d,
        "types/customer.toml",
        "editable  = true\nreconcile = \"edit-persists\"\n\n",
        "editable  = true\n\n",
    );
    assert_refused(&d, "an editable property must declare reconcile");
}

#[test]
fn rule_4_an_action_sets_only_editable_properties() {
    let d = copy();
    edit(
        &d,
        "actions/raise_credit_limit.toml",
        "credit_limit = ",
        "customer_id = ",
    );
    assert_refused(&d, "sets Customer.customer_id, which is not editable");
}

#[test]
fn rule_4_an_edit_only_property_keeps_its_edits() {
    let d = copy();
    edit(
        &d,
        "types/customer.toml",
        "reconcile = \"edit-persists\"\ncolumn",
        "reconcile = \"most-recent\"\ncolumn",
    );
    assert_refused(&d, "always keeps them");
}

#[test]
fn rule_5_links_name_existing_types() {
    let d = copy();
    edit(
        &d,
        "links/customer_orders.toml",
        "to          = \"Order\"",
        "to = \"Invoice\"",
    );
    assert_refused(&d, "links to \"Invoice\", which is not an object type");
}

#[test]
fn rule_5_link_keys_have_matching_types() {
    let d = copy();
    edit(
        &d,
        "types/order.toml",
        "name = \"customer_id\"\ntype = \"string\"",
        "name = \"customer_id\"\ntype = \"integer\"",
    );
    assert_refused(&d, "is string but Order.customer_id is integer");
}

#[test]
fn rule_7_every_catalogue_has_every_key() {
    let d = copy();
    edit(
        &d,
        "languages/pt.toml",
        "\"Customer.credit_limit\" = { name = \"Limite de crédito\" }\n",
        "",
    );
    assert_refused(
        &d,
        "languages/pt.toml: has no entry for Customer.credit_limit",
    );
}

#[test]
fn rule_7_every_language_has_a_catalogue() {
    let d = copy();
    fs::remove_file(d.join("languages/pt.toml")).expect("a removal");
    assert_refused(&d, "languages/pt.toml: is missing");
}

#[test]
fn rule_8_a_sum_is_taken_over_a_number() {
    let d = copy();
    edit(
        &d,
        "metrics/total_credit.toml",
        "expression  = \"credit_limit\"",
        "expression = \"customer_id\"",
    );
    assert_refused(
        &d,
        "cannot take the sum of Customer.customer_id, which is string",
    );
}

#[test]
fn freshness_is_one_of_the_declared_tiers() {
    let d = copy();
    edit(&d, "types/order.toml", "\"real-time\"", "\"instant\"");
    assert_refused(&d, "freshness \"instant\" must be one of");
}

#[test]
fn every_problem_is_reported_in_one_run() {
    let d = copy();
    edit(
        &d,
        "types/order.toml",
        "type = \"decimal\"",
        "type = \"money\"",
    );
    edit(
        &d,
        "metrics/total_credit.toml",
        "aggregation = \"sum\"",
        "aggregation = \"median\"",
    );
    let found = refused(&d);
    assert!(
        found.contains("unknown type \"money\"") && found.contains("aggregation \"median\""),
        "{found}"
    );
}
