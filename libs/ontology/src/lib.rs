//! The one place the ontology's meaning is implemented (roadmap R-51):
//! it loads an ontology written in the format `docs/ontology-format.md`
//! specifies, and refuses one that breaks any of that format's rules.
//! The compiler, the indexer, the query service and the write service
//! all use this library, so no two components can disagree about what
//! an ontology means.

pub mod model;
mod validate;

use std::collections::BTreeMap;
use std::fmt;
use std::fs;
use std::path::Path;

use serde::de::DeserializeOwned; // name-ok

use model::{ActionType, Catalogue, Header, LinkType, Metric, ObjectType};

/// A loaded ontology that satisfies every rule.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Ontology {
    /// Its name, version and languages.
    pub header: Header,
    /// Its object types.
    pub types: Vec<ObjectType>,
    /// Its link types.
    pub links: Vec<LinkType>,
    /// Its action types.
    pub actions: Vec<ActionType>,
    /// Its metrics.
    pub metrics: Vec<Metric>,
    /// Language tag to that language's translation catalogue.
    pub catalogues: BTreeMap<String, Catalogue>,
}

/// One reason an ontology was refused.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Problem {
    /// Where: a file, or the definition the problem belongs to.
    pub place: String,
    /// What is wrong.
    pub message: String,
}

impl fmt::Display for Problem {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}: {}", self.place, self.message)
    }
}

/// Parses the text of one ontology file.
///
/// # Errors
///
/// Returns the parser's description of the first syntax error or
/// misplaced field.
pub fn parse<T: DeserializeOwned>(text: &str) -> Result<T, String> {
    toml::from_str(text).map_err(|e| e.to_string())
}

/// Loads and validates the ontology in `dir`, laid out as the format
/// specifies.
///
/// # Errors
///
/// Returns every problem found: files that cannot be read or parsed
/// first, then every broken rule.
pub fn load(dir: &Path) -> Result<Ontology, Vec<Problem>> {
    let mut problems = Vec::new();
    let header = read_one::<Header>(&dir.join("ontology.toml"), &mut problems);
    let types = read_all::<ObjectType>(&dir.join("types"), &mut problems);
    let links = read_all::<LinkType>(&dir.join("links"), &mut problems);
    let actions = read_all::<ActionType>(&dir.join("actions"), &mut problems);
    let metrics = read_all::<Metric>(&dir.join("metrics"), &mut problems);
    let catalogues = files(&dir.join("languages"), &mut problems)
        .into_iter()
        .filter_map(|path| {
            let tag = path.file_stem()?.to_string_lossy().into_owned();
            read_one::<Catalogue>(&path, &mut problems).map(|c| (tag, c))
        })
        .collect();

    let Some(header) = header else {
        return Err(problems);
    };
    if !problems.is_empty() {
        return Err(problems);
    }
    let ontology = Ontology {
        header,
        types,
        links,
        actions,
        metrics,
        catalogues,
    };
    let broken = validate::validate(&ontology);
    if broken.is_empty() {
        Ok(ontology)
    } else {
        Err(broken)
    }
}

/// Reads and parses one file, recording a problem if either fails.
fn read_one<T: DeserializeOwned>(path: &Path, problems: &mut Vec<Problem>) -> Option<T> {
    let place = path.display().to_string();
    match fs::read_to_string(path) {
        Err(e) => {
            problems.push(Problem {
                place,
                message: format!("cannot be read: {e}"),
            });
            None
        }
        Ok(text) => match parse(&text) {
            Ok(value) => Some(value),
            Err(e) => {
                problems.push(Problem {
                    place,
                    message: e.trim().to_owned(),
                });
                None
            }
        },
    }
}

/// Reads every `.toml` file in `dir`; a missing directory holds none.
fn read_all<T: DeserializeOwned>(dir: &Path, problems: &mut Vec<Problem>) -> Vec<T> {
    files(dir, problems)
        .iter()
        .filter_map(|path| read_one(path, problems))
        .collect()
}

/// The `.toml` files in `dir`, sorted so loading is deterministic.
fn files(dir: &Path, problems: &mut Vec<Problem>) -> Vec<std::path::PathBuf> {
    let Ok(entries) = fs::read_dir(dir) else {
        return Vec::new();
    };
    let mut paths: Vec<_> = entries
        .filter_map(|entry| match entry {
            Ok(e) => Some(e.path()),
            Err(e) => {
                problems.push(Problem {
                    place: dir.display().to_string(),
                    message: e.to_string(),
                });
                None
            }
        })
        .filter(|p| p.extension().is_some_and(|x| x == "toml"))
        .collect();
    paths.sort();
    paths
}
