# The ontology format

How Urshanabi's ontology is written down: its object types, properties,
links, actions, metrics and translations. This document is the format's
open specification (roadmap R-74). The one shared ontology library
implements it (R-51); no other component parses these files.

Decided by the owner on 2026-09-21: the ontology is authored in TOML,
because the only component that parses it is written in the systems
language, which has exactly one canonical, maintained TOML library
(`RULES.md` E1), while its YAML libraries are fragmented across
individual maintainers.

## Layout

One small file per thing, so every change is a readable review:

```
ontology/
  ontology.toml          the ontology's name, version and languages
  types/<name>.toml      one object type each
  links/<name>.toml      one link type each
  actions/<name>.toml    one action type each
  metrics/<name>.toml    one metric each
  languages/<tag>.toml   one translation catalogue per language
```

## Names

Every object type, link, action and metric has a stable `name` in
UpperCamelCase; every property a `name` in lower_snake_case. Names are
identifiers, never shown to people and never translated. Renaming is a
breaking change. What people see comes from the translation catalogues.

## The ontology

```toml
# ontology/ontology.toml
name      = "Sales"
version   = 3
languages = ["en", "es", "pt"]     # BCP 47 tags; the first is the source language
```

## Object types

```toml
# ontology/types/customer.toml
name        = "Customer"
source      = "curated.sales.customers"
primary_key = "customer_id"
freshness   = "near-real-time"
description = "A person or organisation that buys from us."

[[property]]
name           = "customer_id"
type           = "string"
classification = "internal"

[[property]]
name      = "credit_limit"
type      = "decimal"
editable  = true
reconcile = "edit-persists"

[[property]]
name     = "notes"
type     = "text"
editable = true
column   = ""                      # exists only through edits
```

- `source` must name a curated-layer table (R-101).
- `freshness` is `real-time`, `near-real-time` or `scheduled` (R-107).
- `column` defaults to the property's name. An empty `column` declares a
  property that exists only through edits; such properties always keep
  their edits (R-102).
- `classification` may raise a property's label above the one inherited
  through lineage, never lower it (R-76).

## Property types

| Type | Holds |
|---|---|
| `string` | Text that is an identifier or code, never language-processed |
| `text` | Natural-language text, carrying a language tag and processed for search (R-131) |
| `integer`, `decimal`, `float`, `boolean` | As named; `decimal` is exact |
| `date` | A calendar date, never a timestamp (R-128) |
| `timestamp` | An instant, always UTC (R-128) |
| `civil_time` | A future local date-time with its timezone identifier; its UTC instant is derived and recomputed when timezone rules change (R-128) |
| `array<T>` | A list of any type above |

## Editing

An editable property must declare how an edit meets fresh source data
(R-102); there is no default:

- `edit-persists` — the edit stands until edited again;
- `most-recent` — whichever changed last wins;
- `source-priority:<source>` — the named source wins over edits.

## Link types

```toml
# ontology/links/customer_orders.toml
name        = "CustomerOrders"
from        = "Customer"
to          = "Order"
cardinality = "one-to-many"
keys        = [{ from = "customer_id", to = "customer_id" }]
```

`cardinality` is `one-to-one`, `one-to-many`, `many-to-one` or
`many-to-many`; a many-to-many link names its curated join table as
`through`. A link's classification is never declared: it is the least
upper bound of the labels of the fields that establish it (R-90).

## Action types

```toml
# ontology/actions/raise_credit_limit.toml
name  = "RaiseCreditLimit"
risk  = "high"
edits = "Customer"

[[parameter]]
name = "new_limit"
type = "decimal"

[[rule]]
set = { credit_limit = "new_limit" }
```

`risk` is `low` or `high`; high-risk actions need a second person's
approval, never the proposer's own (R-42). An action may only set
editable properties. Writing back to a source is off unless a cell
enables it (R-100).

## Metrics

```toml
# ontology/metrics/total_credit.toml
name        = "TotalCredit"
of          = "Customer"
aggregation = "sum"
expression  = "credit_limit"
```

## Translations

```toml
# ontology/languages/es.toml
"Customer"              = { name = "Cliente", plural = "Clientes", description = "Una persona u organización que nos compra." }
"Customer.credit_limit" = { name = "Límite de crédito" }
"CustomerOrders"        = { name = "Pedidos del cliente" }
"RaiseCreditLimit"      = { name = "Aumentar el límite de crédito" }
```

Keys are stable names: a type, `Type.property`, a link, an action,
`Action.parameter` or a metric. Every language listed in `ontology.toml`
has a catalogue, and every key has an entry in every catalogue (R-130).

## Validation

The ontology library refuses an ontology that breaks any of these rules:

1. Every name is unique within its kind and follows the naming rules.
2. Every type's `source` is a curated table, and its `primary_key` is
   one of its properties.
3. Every property has a known type, and `civil_time` is used only for
   future local times.
4. Every editable property declares `reconcile`; an action sets only
   editable properties.
5. Every link names existing types and properties, with matching key
   types.
6. No declared classification is lower than the label inherited through
   lineage.
7. Every language has a catalogue, and every catalogue has every key.
8. Every metric names an existing type and properties.

## Exporting to the open specification

The ontology exports to the open semantic-model specification now
incubating at a software foundation (R-74): object types become entity
types identified by their primary keys, properties become relationships
to value types, links become relationships between entity types, metrics
become that specification's metrics, and each type's source becomes a
mapping to a dataset. Actions, edit rules, freshness and translations
travel in its extension mechanism. Security labels are never exported:
an exported model is not a security boundary, as that specification
itself warns.
