//! The ontology's files as written, one type per kind of file
//! (`docs/ontology-format.md`). Unknown fields are refused, so a
//! misspelled field is an error, never silently ignored.

use std::collections::BTreeMap;

use serde::Deserialize; // name-ok

/// `ontology/ontology.toml`: the ontology's name, version and languages.
#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(deny_unknown_fields)] // name-ok
pub struct Header {
    /// The ontology's stable name.
    pub name: String,
    /// Its version, counting up from 1.
    pub version: u32,
    /// Its languages as BCP 47 tags; the first is the source language.
    pub languages: Vec<String>,
}

/// `ontology/types/<name>.toml`: one object type.
#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(deny_unknown_fields)] // name-ok
pub struct ObjectType {
    /// The type's stable name, in `UpperCamelCase`.
    pub name: String,
    /// The curated-layer table the type is read from (roadmap R-101).
    pub source: String,
    /// The property that identifies each object.
    pub primary_key: String,
    /// The property shown as each object's display name.
    pub title: String,
    /// `real-time`, `near-real-time` or `scheduled` (roadmap R-107).
    pub freshness: String,
    /// A description in the source language.
    pub description: Option<String>,
    /// `active`, `experimental` or `deprecated`; `experimental` when
    /// absent.
    pub status: Option<String>,
    /// The type's properties.
    #[serde(default, rename = "property")] // name-ok
    pub properties: Vec<Property>,
}

/// One property of an object type.
#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(deny_unknown_fields)] // name-ok
pub struct Property {
    /// The property's stable name, in `lower_snake_case`.
    pub name: String,
    /// The property's type, such as `string` or `array<date>`.
    #[serde(rename = "type")] // name-ok
    pub kind: String,
    /// Whether people may edit it.
    #[serde(default)] // name-ok
    pub editable: bool,
    /// How an edit meets fresh source data; required when editable
    /// (roadmap R-102).
    pub reconcile: Option<String>,
    /// The source column; the property's name when absent, and empty
    /// for a property that exists only through edits.
    pub column: Option<String>,
    /// A label raising the property above the one it inherits (R-76).
    pub classification: Option<String>,
    /// `active`, `experimental` or `deprecated`; `experimental` when
    /// absent.
    pub status: Option<String>,
}

/// `ontology/links/<name>.toml`: one link type.
#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(deny_unknown_fields)] // name-ok
pub struct LinkType {
    /// The link's stable name, in `UpperCamelCase`.
    pub name: String,
    /// The type the link starts from.
    pub from: String,
    /// The type the link points to.
    pub to: String,
    /// `one-to-one`, `one-to-many`, `many-to-one` or `many-to-many`.
    pub cardinality: String,
    /// The property pairs that join the two types.
    #[serde(default)] // name-ok
    pub keys: Vec<Key>,
    /// The curated join table of a many-to-many link.
    pub through: Option<String>,
    /// `active`, `experimental` or `deprecated`; `experimental` when
    /// absent.
    pub status: Option<String>,
}

/// One property pair joining a link's two types.
#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(deny_unknown_fields)] // name-ok
pub struct Key {
    /// The property on the `from` type.
    pub from: String,
    /// The property on the `to` type.
    pub to: String,
}

/// `ontology/actions/<name>.toml`: one action type.
#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(deny_unknown_fields)] // name-ok
pub struct ActionType {
    /// The action's stable name, in `UpperCamelCase`.
    pub name: String,
    /// `low` or `high`; high-risk actions need a second person (R-42).
    pub risk: String,
    /// The object type the action edits; exactly one parameter,
    /// typed `object<Type>`, names which object.
    pub edits: String,
    /// `active`, `experimental` or `deprecated`; `experimental` when
    /// absent.
    pub status: Option<String>,
    /// The action's parameters.
    #[serde(default, rename = "parameter")] // name-ok
    pub parameters: Vec<Parameter>,
    /// What the action does.
    #[serde(default, rename = "rule")] // name-ok
    pub rules: Vec<Rule>,
}

/// One parameter of an action.
#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(deny_unknown_fields)] // name-ok
pub struct Parameter {
    /// The parameter's stable name, in `lower_snake_case`.
    pub name: String,
    /// The parameter's type: a property's type, or `object<Type>` to
    /// refer to an existing object.
    #[serde(rename = "type")] // name-ok
    pub kind: String,
}

/// One rule of an action: the properties it sets, each from a parameter.
#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(deny_unknown_fields)] // name-ok
pub struct Rule {
    /// Property name to the parameter that supplies its new value.
    #[serde(default)] // name-ok
    pub set: BTreeMap<String, String>,
}

/// `ontology/metrics/<name>.toml`: one metric.
#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(deny_unknown_fields)] // name-ok
pub struct Metric {
    /// The metric's stable name, in `UpperCamelCase`.
    pub name: String,
    /// The object type it aggregates.
    pub of: String,
    /// `count`, `sum`, `average`, `min` or `max`.
    pub aggregation: String,
    /// The property aggregated; absent only for `count`.
    pub expression: Option<String>,
    /// `active`, `experimental` or `deprecated`; `experimental` when
    /// absent.
    pub status: Option<String>,
}

/// One entry of a translation catalogue (roadmap R-130).
#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(deny_unknown_fields)] // name-ok
pub struct Entry {
    /// The display name.
    pub name: String,
    /// The plural display name.
    pub plural: Option<String>,
    /// The description.
    pub description: Option<String>,
}

/// `ontology/languages/<tag>.toml`: stable name to its translation.
pub type Catalogue = BTreeMap<String, Entry>;
