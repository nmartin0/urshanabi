//! The format's validation rules (`docs/ontology-format.md`). Every
//! rule is checked and every problem reported, so one run shows an
//! author everything to fix.

use std::collections::{BTreeMap, BTreeSet};

use crate::model::{ObjectType, Property};
use crate::{Ontology, Problem};

/// The scalar property types.
const SCALARS: &[&str] = &[
    "string",
    "text",
    "integer",
    "decimal",
    "float",
    "boolean",
    "date",
    "timestamp",
    "civil_time",
];
/// The types a sum or an average can be taken over.
const NUMERIC: &[&str] = &["integer", "decimal", "float"];
const FRESHNESS: &[&str] = &["real-time", "near-real-time", "scheduled"];
const CARDINALITY: &[&str] = &["one-to-one", "one-to-many", "many-to-one", "many-to-many"];
const RISK: &[&str] = &["low", "high"];
const AGGREGATION: &[&str] = &["count", "sum", "average", "min", "max"];

/// Collects problems as the rules are checked.
struct Report(Vec<Problem>);

impl Report {
    fn add(&mut self, place: &str, message: impl Into<String>) {
        self.0.push(Problem {
            place: place.to_owned(),
            message: message.into(),
        });
    }
}

/// Checks every rule, returning every problem found.
pub(crate) fn validate(o: &Ontology) -> Vec<Problem> {
    let mut r = Report(Vec::new());
    header(o, &mut r);
    unique_names(
        o.types.iter().map(|t| t.name.as_str()),
        "object type",
        &mut r,
    );
    unique_names(o.links.iter().map(|l| l.name.as_str()), "link type", &mut r);
    unique_names(
        o.actions.iter().map(|a| a.name.as_str()),
        "action type",
        &mut r,
    );
    unique_names(o.metrics.iter().map(|m| m.name.as_str()), "metric", &mut r);
    let by_name: BTreeMap<&str, &ObjectType> =
        o.types.iter().map(|t| (t.name.as_str(), t)).collect();
    for t in &o.types {
        object_type(t, &mut r);
    }
    links(o, &by_name, &mut r);
    actions(o, &by_name, &mut r);
    metrics(o, &by_name, &mut r);
    translations(o, &mut r);
    r.0
}

fn header(o: &Ontology, r: &mut Report) {
    let place = "ontology.toml";
    if o.header.version == 0 {
        r.add(place, "the version counts up from 1");
    }
    if o.header.languages.is_empty() {
        r.add(place, "at least one language is required");
    }
    let mut seen = BTreeSet::new();
    for tag in &o.header.languages {
        if !language_tag(tag) {
            r.add(
                place,
                format!("{tag:?} is not a language tag such as \"es\" or \"pt-BR\""),
            );
        }
        if !seen.insert(tag) {
            r.add(place, format!("the language {tag} is listed twice"));
        }
    }
}

/// Rule 1: names are unique within their kind and in `UpperCamelCase`.
fn unique_names<'a>(names: impl Iterator<Item = &'a str>, kind: &str, r: &mut Report) {
    let mut seen = BTreeSet::new();
    for name in names {
        if !upper_camel(name) {
            r.add(name, format!("a {kind} name must be UpperCamelCase"));
        }
        if !seen.insert(name) {
            r.add(name, format!("two {kind}s share this name"));
        }
    }
}

/// Rules 1 to 4 for one object type.
fn object_type(t: &ObjectType, r: &mut Report) {
    let place = t.name.as_str();
    // Rule 2: the source is in the curated layer. Whether the table
    // exists is checked against the table catalogue when one is given.
    if !t.source.starts_with("curated.") || t.source.split('.').count() < 3 {
        r.add(place, format!("the source {:?} must name a curated-layer table, such as \"curated.sales.customers\"", t.source));
    }
    if !FRESHNESS.contains(&t.freshness.as_str()) {
        r.add(
            place,
            format!(
                "freshness {:?} must be one of {}",
                t.freshness,
                FRESHNESS.join(", ")
            ),
        );
    }
    let mut seen = BTreeSet::new();
    for p in &t.properties {
        let at = format!("{}.{}", t.name, p.name);
        if !lower_snake(&p.name) {
            r.add(&at, "a property name must be lower_snake_case");
        }
        if !seen.insert(p.name.as_str()) {
            r.add(&at, "the property is declared twice");
        }
        property(p, &at, r);
    }
    if !seen.contains(t.primary_key.as_str()) {
        r.add(
            place,
            format!(
                "the primary key {:?} is not one of its properties",
                t.primary_key
            ),
        );
    }
}

/// Rules 3 and 4 for one property.
fn property(p: &Property, at: &str, r: &mut Report) {
    if !known_type(&p.kind) {
        r.add(at, format!("unknown type {:?}", p.kind));
    }
    let edit_only = p.column.as_deref() == Some("");
    if p.editable {
        match p.reconcile.as_deref() {
            None => r.add(at, "an editable property must declare reconcile"),
            Some(rule) if !reconcile_rule(rule) => r.add(
                at,
                format!("reconcile {rule:?} must be edit-persists, most-recent or source-priority:<source>"),
            ),
            Some(rule) if edit_only && rule != "edit-persists" => {
                r.add(at, "a property that exists only through edits always keeps them: reconcile = \"edit-persists\"");
            }
            Some(_) => {}
        }
    } else {
        if p.reconcile.is_some() {
            r.add(at, "reconcile applies only to an editable property");
        }
        if edit_only {
            r.add(
                at,
                "a property that exists only through edits must be editable",
            );
        }
    }
}

/// Rule 5: links name existing types and properties, with matching types.
fn links(o: &Ontology, types: &BTreeMap<&str, &ObjectType>, r: &mut Report) {
    for l in &o.links {
        let at = l.name.as_str();
        if !CARDINALITY.contains(&l.cardinality.as_str()) {
            r.add(
                at,
                format!(
                    "cardinality {:?} must be one of {}",
                    l.cardinality,
                    CARDINALITY.join(", ")
                ),
            );
        }
        match (l.cardinality == "many-to-many", &l.through) {
            (true, None) => r.add(
                at,
                "a many-to-many link names its curated join table as through",
            ),
            (false, Some(_)) => r.add(at, "only a many-to-many link has a through table"),
            _ => {}
        }
        let (from, to) = (types.get(l.from.as_str()), types.get(l.to.as_str()));
        if from.is_none() {
            r.add(
                at,
                format!("links from {:?}, which is not an object type", l.from),
            );
        }
        if to.is_none() {
            r.add(
                at,
                format!("links to {:?}, which is not an object type", l.to),
            );
        }
        if l.keys.is_empty() && l.through.is_none() {
            r.add(at, "a link needs at least one key");
        }
        if let (Some(from), Some(to)) = (from, to) {
            for k in &l.keys {
                let (a, b) = (find(from, &k.from), find(to, &k.to));
                if a.is_none() {
                    r.add(at, format!("{} has no property {:?}", from.name, k.from));
                }
                if b.is_none() {
                    r.add(at, format!("{} has no property {:?}", to.name, k.to));
                }
                if let (Some(a), Some(b)) = (a, b)
                    && a.kind != b.kind
                {
                    r.add(
                        at,
                        format!(
                            "the key {}.{} is {} but {}.{} is {}",
                            from.name, a.name, a.kind, to.name, b.name, b.kind
                        ),
                    );
                }
            }
        }
    }
}

/// Rule 4 for actions: they set only editable properties of the type
/// they edit, from their own parameters.
fn actions(o: &Ontology, types: &BTreeMap<&str, &ObjectType>, r: &mut Report) {
    for a in &o.actions {
        let at = a.name.as_str();
        if !RISK.contains(&a.risk.as_str()) {
            r.add(at, format!("risk {:?} must be low or high", a.risk));
        }
        let mut params = BTreeSet::new();
        for p in &a.parameters {
            let pat = format!("{}.{}", a.name, p.name);
            if !lower_snake(&p.name) {
                r.add(&pat, "a parameter name must be lower_snake_case");
            }
            if !params.insert(p.name.as_str()) {
                r.add(&pat, "the parameter is declared twice");
            }
            if !known_type(&p.kind) {
                r.add(&pat, format!("unknown type {:?}", p.kind));
            }
        }
        let Some(target) = types.get(a.edits.as_str()) else {
            r.add(
                at,
                format!("edits {:?}, which is not an object type", a.edits),
            );
            continue;
        };
        for rule in &a.rules {
            for (property, param) in &rule.set {
                match find(target, property) {
                    None => r.add(
                        at,
                        format!("sets {property:?}, which {} does not have", target.name),
                    ),
                    Some(p) if !p.editable => r.add(
                        at,
                        format!("sets {}.{property}, which is not editable", target.name),
                    ),
                    Some(_) => {}
                }
                if !params.contains(param.as_str()) {
                    r.add(
                        at,
                        format!(
                            "sets {property} from {param:?}, which is not one of its parameters"
                        ),
                    );
                }
            }
        }
    }
}

/// Rule 8: metrics name existing types and properties.
fn metrics(o: &Ontology, types: &BTreeMap<&str, &ObjectType>, r: &mut Report) {
    for m in &o.metrics {
        let at = m.name.as_str();
        if !AGGREGATION.contains(&m.aggregation.as_str()) {
            r.add(
                at,
                format!(
                    "aggregation {:?} must be one of {}",
                    m.aggregation,
                    AGGREGATION.join(", ")
                ),
            );
        }
        let Some(of) = types.get(m.of.as_str()) else {
            r.add(
                at,
                format!("aggregates {:?}, which is not an object type", m.of),
            );
            continue;
        };
        match (m.aggregation.as_str(), m.expression.as_deref()) {
            ("count", None) => {}
            (_, None) => r.add(at, "only a count may omit its expression"),
            (aggregation, Some(e)) => match find(of, e) {
                None => r.add(at, format!("{} has no property {e:?}", of.name)),
                Some(p)
                    if matches!(aggregation, "sum" | "average")
                        && !NUMERIC.contains(&p.kind.as_str()) =>
                {
                    r.add(
                        at,
                        format!(
                            "cannot take the {aggregation} of {}.{e}, which is {}",
                            of.name, p.kind
                        ),
                    );
                }
                Some(_) => {}
            },
        }
    }
}

/// Rule 7: every language has a catalogue, every catalogue has every
/// key, and no catalogue has a key the ontology lacks.
fn translations(o: &Ontology, r: &mut Report) {
    let mut keys = BTreeSet::new();
    for t in &o.types {
        keys.insert(t.name.clone());
        keys.extend(
            t.properties
                .iter()
                .map(|p| format!("{}.{}", t.name, p.name)),
        );
    }
    keys.extend(o.links.iter().map(|l| l.name.clone()));
    for a in &o.actions {
        keys.insert(a.name.clone());
        keys.extend(
            a.parameters
                .iter()
                .map(|p| format!("{}.{}", a.name, p.name)),
        );
    }
    keys.extend(o.metrics.iter().map(|m| m.name.clone()));

    for tag in &o.header.languages {
        let place = format!("languages/{tag}.toml");
        let Some(catalogue) = o.catalogues.get(tag) else {
            r.add(&place, "is missing: every listed language has a catalogue");
            continue;
        };
        for key in &keys {
            if !catalogue.contains_key(key) {
                r.add(&place, format!("has no entry for {key}"));
            }
        }
        for key in catalogue.keys() {
            if !keys.contains(key) {
                r.add(
                    &place,
                    format!("translates {key}, which the ontology does not define"),
                );
            }
        }
    }
    for tag in o.catalogues.keys() {
        if !o.header.languages.contains(tag) {
            r.add(
                &format!("languages/{tag}.toml"),
                "is for a language ontology.toml does not list",
            );
        }
    }
}

fn find<'t>(t: &'t ObjectType, name: &str) -> Option<&'t Property> {
    t.properties.iter().find(|p| p.name == name)
}

/// Rule 3: a scalar type, or an array of one.
fn known_type(kind: &str) -> bool {
    let inner = kind
        .strip_prefix("array<")
        .and_then(|k| k.strip_suffix('>'))
        .unwrap_or(kind);
    SCALARS.contains(&inner)
}

fn reconcile_rule(rule: &str) -> bool {
    matches!(rule, "edit-persists" | "most-recent")
        || rule
            .strip_prefix("source-priority:")
            .is_some_and(|s| !s.is_empty())
}

fn upper_camel(name: &str) -> bool {
    name.starts_with(|c: char| c.is_ascii_uppercase())
        && name.chars().all(|c| c.is_ascii_alphanumeric())
}

fn lower_snake(name: &str) -> bool {
    name.starts_with(|c: char| c.is_ascii_lowercase())
        && name
            .chars()
            .all(|c| c.is_ascii_lowercase() || c.is_ascii_digit() || c == '_')
        && !name.ends_with('_')
        && !name.contains("__")
}

fn language_tag(tag: &str) -> bool {
    let mut parts = tag.split('-');
    let first_ok = parts
        .next()
        .is_some_and(|p| (2..=3).contains(&p.len()) && p.chars().all(|c| c.is_ascii_lowercase()));
    first_ok
        && parts.all(|p| (2..=8).contains(&p.len()) && p.chars().all(|c| c.is_ascii_alphanumeric()))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn names_follow_their_case_rules() {
        assert!(upper_camel("RaiseCreditLimit"));
        assert!(!upper_camel("raiseCreditLimit"));
        assert!(!upper_camel("Raise_Limit"));
        assert!(lower_snake("credit_limit"));
        assert!(!lower_snake("creditLimit"));
        assert!(!lower_snake("credit__limit"));
        assert!(!lower_snake("limit_"));
    }

    #[test]
    fn types_are_scalars_or_arrays_of_scalars() {
        assert!(known_type("civil_time"));
        assert!(known_type("array<date>"));
        assert!(!known_type("array<array<date>>"));
        assert!(!known_type("datetime"));
    }

    #[test]
    fn reconcile_rules_are_the_three_declared_forms() {
        assert!(reconcile_rule("edit-persists"));
        assert!(reconcile_rule("most-recent"));
        assert!(reconcile_rule("source-priority:billing"));
        assert!(!reconcile_rule("source-priority:"));
        assert!(!reconcile_rule("newest"));
    }

    #[test]
    fn language_tags_follow_bcp_47_shape() {
        assert!(language_tag("es"));
        assert!(language_tag("pt-BR"));
        assert!(!language_tag("Spanish"));
        assert!(!language_tag("e"));
    }
}
