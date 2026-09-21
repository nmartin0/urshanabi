# Behaviours

What the product guarantees over HTTP. This is the specification the
conformance suite implements: the suite runs against Elysium first,
must pass there, and is then pointed at Urshanabi.

Derived from Elysium at commit `f4ea94e` — 345 integration tests, of
which 328 run without a model. Every entry below cites the test that
proves it in Elysium, so each property is something a test has
already made fail, not an intention.

## How to read an entry

- **The property** is what must hold. Changing it is a breaking
  change under `RULES.md` H5.
- **Representation** is Elysium's current wire shape: status codes,
  field names, header values. Urshanabi may change a representation,
  but only deliberately, and the change is recorded here first.
- **Evidence** names the Elysium test file and test.
- **Level** says how a black-box suite can check it:
  - `http` — observable through the HTTP interface alone.
  - `oracle` — needs a second channel the suite does not yet have
    (the audit trail, a clock, a timing measurement).
  - `internal` — proven below the HTTP layer in Elysium; the suite
    needs an HTTP-level equivalent before it can claim it.

Identifiers (`AUTH-01`) are stable. A retired entry keeps its number
and is struck through.

---

## AUTH — authentication and sessions

**AUTH-01** A wrong password and an unknown username are
indistinguishable. *Representation:* both 401, identical body.
*Evidence:* `test_api.py::test_login_wrong_password_and_nonexistent_username_are_identical`.
*Level:* http.

**AUTH-02** Repeated failures lock a username out, and the locked-out
response is identical to an ordinary wrong password. *Evidence:*
`test_login_locked_out_after_max_failed_attempts`,
`test_login_lockout_response_is_identical_to_a_normal_wrong_password`.
*Level:* http.

**AUTH-03** Lockout applies to usernames that do not exist, so
throttling cannot reveal which accounts are real. *Evidence:*
`test_login_lockout_applies_to_a_nonexistent_username_too`. *Level:*
http.

**AUTH-04** Lockout is per username, not global; a success clears
prior failures. *Evidence:* `test_login_lockout_is_per_username_not_global`,
`test_login_success_clears_prior_failures`. *Level:* http.

**AUTH-05** Password verification runs even when an account is
already locked out, so response time does not reveal the lockout.
*Evidence:* `test_login_always_runs_real_password_verification_even_when_already_locked_out`
(proven with an internal spy). *Level:* oracle — needs a timing
measurement. An external review measured a 1.2 ms difference in a
163 ms operation between unknown and known usernames.

**AUTH-06** A successful login issues two distinct cookies: a session
cookie script cannot read, and an anti-forgery cookie script can.
*Representation:* 204 with an empty body; cookies `elysium_session`
and `elysium_csrf`. *Evidence:*
`test_login_success_sets_real_session_and_csrf_cookies`. *Level:* http.

**AUTH-07** State-changing requests require the anti-forgery header to
match its cookie. *Representation:* 403. *Evidence:*
`test_the_new_routes_require_a_csrf_token`. *Level:* http.

**AUTH-08** Unauthenticated requests are refused. *Representation:*
reads return 401; state-changing requests return 403, because the
anti-forgery check runs before the session check. *Evidence:*
`test_search_objects_without_token_is_rejected`,
`test_query_without_token_is_rejected`,
`test_propose_action_without_token_is_rejected`. *Level:* http.

**AUTH-09** Logging out invalidates that session; "log out
everywhere" revokes every session of the caller; an administrator may
revoke another user's sessions only with user-management rights.
*Evidence:* `test_logout_invalidates_the_token`,
`test_logout_all_revokes_every_session_for_the_caller`,
`test_admin_logout_all_for_a_target_user_requires_manage_users`,
`test_admin_logout_all_for_a_target_user_works`. *Level:* http.

**AUTH-10** Disabling a user blocks new logins with the ordinary
failure message and ends existing sessions; enabling reverses it;
deleting removes the credential and ends sessions. Acting on a user
who does not exist is a 404. *Evidence:*
`test_disable_user_blocks_new_logins_and_kills_existing_sessions`,
`test_enable_reverses_disable`,
`test_delete_user_removes_credential_and_kills_sessions`,
`test_disable_nonexistent_user_is_404`,
`test_delete_nonexistent_user_is_404`. *Level:* http.

---

## HDR — transport and headers

**HDR-01** Every response carries the security headers.
*Representation:* `x-content-type-options: nosniff`,
`x-frame-options: DENY`, `referrer-policy: same-origin`,
`content-security-policy: default-src 'self'; object-src 'none';
base-uri 'self'; frame-ancestors 'none'`. *Evidence:*
`test_security_headers_are_present_on_every_response`. *Level:* http.

**HDR-02** Every response specific to its caller forbids caching,
including write details. *Representation:*
`cache-control: no-store`. *Evidence:*
`test_me_routes_set_cache_control_no_store`,
`test_every_per_caller_route_forbids_caching`,
`test_a_write_detail_forbids_caching_too`. *Level:* http.

**HDR-03** Every response — including unauthenticated ones — says
which configuration generation answered it, and the number changes
when configuration does. *Representation:* header
`x-elysium-generation`; Urshanabi will rename it, which is a
deliberate representation change. *Evidence:*
`test_every_response_says_which_generation_answered`,
`test_the_number_changes_when_the_configuration_does`,
`test_an_unauthenticated_response_carries_it_too`. *Level:* http.

**HDR-04** No machine-readable description of the interface is
served, since it would list every administrative route to an
unauthenticated caller. *Representation:* 404. *Evidence:* the interface-description
availability test in `test_api.py`, cited by description because its
name contains product names. *Level:* http.

---

## SCHEMA — what a caller is told exists

**SCHEMA-01** A caller's profile reports their own username, role and
security value, and differs by caller. *Evidence:*
`test_my_profile_returns_the_callers_own_username_role_and_mac_value`,
`test_my_profile_differs_by_which_user_is_logged_in`. *Level:* http.

**SCHEMA-02** The visible schema is the caller's own and differs by
role; a field without a grant is absent. *Evidence:*
`test_my_visible_schema_returns_the_callers_own_view`,
`test_my_visible_schema_differs_by_role_not_a_static_response`.
*Level:* http.

**SCHEMA-03** The visible schema never carries storage internals, and
does carry what links need: target, cardinality and a link type that
is the same at both ends. *Evidence:*
`test_visible_schema_never_leaks_per_field_internals`,
`test_visible_schema_still_carries_what_links_genuinely_need`,
`test_both_ends_of_a_relationship_report_the_same_link_type`.
*Level:* http.

**SCHEMA-04** Declared display metadata reaches the caller verbatim.
The title field is reported only when the caller may read it, and is
otherwise null rather than absent. *Evidence:*
`test_declared_display_metadata_reaches_the_caller_verbatim`,
`test_my_visible_schema_shows_title_field_when_granted`,
`test_my_visible_schema_withholds_title_field_when_not_granted`.
*Level:* http.

**SCHEMA-05** The permission ladder. With discovery but not read
access, a field is named and marked unreadable and its value is
withheld; a type is listed and marked unsearchable. With no grant, a
type is absent. The agent's view omits anything unreadable entirely.
*Evidence:* `test_grant_ladder_enforcement.py` (all nine),
`test_a_discover_only_field_reports_itself_unreadable_over_http`.
*Level:* http for the interface view; internal for the agent view.

**SCHEMA-06** The list of visible applications is filtered by grant
and never reveals which permission gates each one. The notifications
application is listed for every role. *Evidence:*
`test_visible_apps_hides_admin_without_manage_users`,
`test_visible_apps_never_leaks_gating_permission`,
`test_notifications_app.py`. *Level:* http.

---

## DENY — uniform denial on reads

Each endpoint is uniform within itself. The shapes differ between
endpoints; that is recorded as the representation, not as a defect.

**DENY-01** Search returns the same body for a type the caller cannot
see, a type that does not exist, a MAC-denied match and no match.
*Representation:* 200, `results: []`, `total_matches: 0`,
`next_page_token: null`, `scan_truncated: false`, plus a
`request_id`. *Evidence:* `test_search_objects_blocks_cross_region_mac`,
`test_search_objects_unknown_type_returns_empty_results_not_error`,
`test_search_objects_no_match_returns_empty_results`. *Level:* http.
**Also verified for this specification:** a real type the caller
cannot discover and an invented one produce identical search bodies,
and both traces are empty.
**Elysium meets this only below its scan ceiling.** Above it, the
truncation flag can differ between a denied search and an empty one.
See DENY-09.

**DENY-02** Object detail returns the same shape for a nonexistent
object and a MAC-denied one: the same field keys, every value null.
For an unknown or undiscoverable type it returns empty fields and a
null `request_id`. *Evidence:*
`test_object_detail_nonexistent_id_returns_200_with_every_field_null`,
`test_object_detail_cross_region_mac_denial_is_identical_to_nonexistent`,
`test_object_detail_unknown_type_returns_200_with_empty_fields`.
*Level:* http. **Also verified:** a real undiscoverable type and an
invented one return byte-identical detail bodies.

**DENY-03** A filter on an unreadable field fails exactly like a
filter on a field that does not exist, and the message names neither.
*Representation:* 400, identical bodies. *Evidence:*
`test_a_condition_on_an_unreadable_field_is_a_400`,
`test_a_condition_on_an_unreadable_field_is_rejected_by_search_too`,
`test_error_messages_do_not_leak.py`. *Level:* http.

**DENY-04** Traversing a link returns only objects the caller could
read directly; an ungranted link field returns empty, not an error.
*Evidence:* `test_search_around_returns_only_ids_the_caller_could_read_directly`,
`test_search_around_on_an_ungranted_field_returns_empty_not_an_error`.
*Level:* http.

**DENY-05** An unreadable object's edit history is empty. *Evidence:*
`test_object_history_denies_a_caller_who_cannot_read_the_object`.
*Level:* http.

**DENY-06** Notes on an unreadable object read as an empty list;
writing one is a 404. *Evidence:*
`test_reading_notes_on_an_unreadable_object_is_empty_not_an_error`,
`test_writing_a_note_on_an_unreadable_object_is_a_404`. *Level:*
http.

**DENY-07** Link counts omit a link to an invisible type or a
withheld link field — omitted, not reported as zero. *Evidence:*
`test_link_counts.py`. *Level:* internal; the HTTP route is covered
only for serialisation (`test_link_counts_survive_serialisation`).

**DENY-08** A request trace belongs to the caller: another user's
trace, or an unknown id, is an empty list. One request's trace never
includes another's. *Evidence:*
`test_a_request_trace_is_empty_for_a_request_you_did_not_make`,
`test_another_user_cannot_read_your_trace`,
`test_one_request_s_trace_does_not_include_another_s`. *Level:* http.

**DENY-09** No signal in a response may depend on rows the caller
cannot see: not a truncation flag, not a count or total, not a page
boundary, and — as far as can be measured — not response time.
Counts, totals and truncation are derived only from rows that passed
the security check, and the security condition is applied before any
condition the caller supplied, except operators on a reviewed list of
those that cannot reveal what they are evaluated on. *Evidence:*
**Elysium does not meet this.** Measured at its commit `5ddb031`:
with the scan ceiling lowered to 2, a user who can see no
transactions at all receives `scan_truncated: true`, on both read
paths, because the ceiling bounds a scan taken before the security
check wherever security is inherited through a link. Types whose
security value is on the object itself do not leak. Precedent: a
published study found row-level security measurably leaks the size of
hidden row sets through timing whenever the caller's condition is
evaluated first; a mature database orders the security condition
first for exactly this reason. *Level:* http for the flag, with a
fixture above the ceiling; oracle for timing.

---

## QUERY — search, filters, counts and paging

**QUERY-01** Free-text search matches partially; an empty or absent
query browses everything visible. *Evidence:*
`test_search_objects_finds_a_partial_match_with_real_field_values`,
`test_search_objects_empty_query_returns_every_visible_result`,
`test_search_objects_no_query_param_at_all_also_browses_all`.
*Level:* http.

**QUERY-02** Structured conditions use the closed vocabulary; text and
conditions narrow together; the older mapping form still works;
sending both forms is refused rather than merged; malformed conditions
are a 400, never a 500. *Evidence:*
`test_count_accepts_the_full_condition_vocabulary`,
`test_text_and_conditions_narrow_together`,
`test_the_dict_form_still_works`,
`test_sending_both_forms_is_rejected_rather_than_merged`,
`test_malformed_conditions_are_a_400`. *Level:* http.

**QUERY-03** Paging returns every result exactly once, even when the
underlying order is unstable. `total_matches` counts everything, not
the page. The last page has no next token. An oversized page size is
clamped, not refused. *Evidence:*
`test_paging_returns_every_result_exactly_once`,
`test_paging_survives_an_unstable_underlying_order`,
`test_total_matches_counts_everything_not_just_the_page`,
`test_the_last_page_has_no_next_token`,
`test_an_oversized_page_size_is_clamped_not_rejected`. *Level:* http.

**QUERY-04** Page tokens are opaque. A bare offset, a malformed token
or an out-of-range token returns the first page. *Representation:*
tokens begin `v1.`. *Evidence:* `test_a_page_token_is_opaque_not_an_offset`,
`test_a_bare_offset_is_no_longer_accepted`,
`test_every_malformed_token_shape_falls_back_to_the_first_page`,
`test_history_tokens_are_opaque_too`. *Level:* http.

**QUERY-05** Results can be ordered by a field, ascending or
descending, consistently across pages; an unknown ordering field is
ignored rather than scrambling the order. *Evidence:*
`test_results_are_ordered_by_a_requested_field`,
`test_ordering_stays_consistent_across_pages`,
`test_an_unknown_order_by_field_is_ignored_rather_than_scrambling`.
*Level:* http.

**QUERY-06** A search that stops at its ceiling says so, and still
returns what it found — where the ceiling counts only rows the caller
may see, per DENY-09. A flag raised by invisible rows is a disclosure,
not a report. *Representation:* `scan_truncated`.
*Evidence:* `test_search_reports_truncation.py`. *Level:* internal —
Elysium lowers the ceiling in-process; the suite needs a fixture
large enough to reach it.

**QUERY-07** Fetching every matching id returns all matches, not a
page, honours the filter, shows only what the caller can see, and
above its ceiling **refuses** rather than truncating, naming the
ceiling and telling the caller to narrow the filter. *Evidence:*
`test_matching_ids_returns_every_match_not_a_page`,
`test_matching_ids_honours_the_filter`,
`test_matching_ids_shows_only_what_the_caller_can_see`,
`test_matching_ids_refuses_rather_than_truncating`. *Level:* http.

**QUERY-08** Object ids are strings in every response, whatever the
stored key type — **including the values of link fields.**
*Evidence:* `test_search_returns_string_ids_for_an_integer_keyed_type`,
`test_search_and_matching_ids_agree_on_representation`,
`test_search_and_detail_agree_on_representation`. *Level:* http.
**Elysium does not meet the link half of this.** On its live read
path, object detail returns link values as integers while search
returns the same objects' ids as strings; on its mirror path both are
strings. `test_object_detail_returns_every_visible_field_including_a_link`
asserts the integer form. This specification takes the string form
deliberately, and that Elysium test is the one expected divergence.

**QUERY-09** Counts are the caller's own and differ between callers.
Aggregation groups and sums; with no grouping it returns one result
under an empty key. An unknown aggregate, or one missing its field, is
a 400 naming the problem. *Evidence:*
`test_two_users_get_genuinely_different_counts`,
`test_aggregate_groups_and_sums`,
`test_aggregate_with_no_group_by_returns_one_result_under_an_empty_key`,
`test_an_unknown_aggregate_is_a_400_not_a_500`,
`test_an_aggregate_missing_its_field_is_a_400`. *Level:* http.

**QUERY-10** Route names never shadow data: an object whose id is
`count` or `search` is still reachable. *Evidence:*
`test_an_object_genuinely_named_count_is_still_reachable`,
`test_object_detail_and_search_routes_do_not_collide`. *Level:* http.

---

## ACT — actions and writes

**ACT-01** The action catalogue is the caller's own, differs by role,
and never reveals how an action changes data. *Representation:* each
entry has exactly `affected_object_types`, `parameters`,
`executable`. *Evidence:*
`test_visible_action_types_returns_the_callers_own_view`,
`test_visible_action_types_never_leaks_sub_writes_or_mutations`,
`test_visible_action_types_differs_by_role_not_a_static_response`.
*Level:* http.

**ACT-02** Discovery and execution are separate. With the catalogue
grant, every action is listed and `executable` says which may be run;
without it, only executable actions are listed. *Evidence:*
`test_visible_action_types_with_discover_grant_shows_the_whole_catalog`,
`test_visible_action_types_executable_flag_differentiates_within_one_response`,
`test_visible_action_types_executable_flag_is_always_true_without_discover_grant`.
*Level:* http.

**ACT-03** Without the catalogue grant, an unknown action, an
unauthorized one, a missing parameter and a MAC denial are
indistinguishable. *Representation:* 400 with "That action could not
be proposed. Check the action name and parameters, and that you're
authorized to perform it." *Evidence:*
`test_propose_action_unknown_action_and_real_but_unauthorized_action_are_identical`,
`test_propose_action_missing_required_parameter_returns_the_same_generic_error`,
`test_propose_action_cross_region_mac_denial_returns_the_same_generic_error`.
*Level:* http.

**ACT-04** With the catalogue grant, the two become distinguishable by
design, since the caller can already see the catalogue.
*Representation:* unknown is 400 naming the action; unauthorized is
403 naming the user and the missing grant. *Evidence:*
`test_propose_action_unknown_vs_unauthorized_are_no_longer_identical_for_a_discover_holder`.
*Level:* http.

**ACT-05** Proposing returns a pending write carrying each change and
the values it expects to replace; nothing is written until confirmed.
Approval makes the change visible through ordinary reads; rejection
leaves data unchanged. *Representation:* 202; confirm returns
`written` or `rejected`. *Evidence:*
`test_propose_action_succeeds_and_returns_a_real_pending_write`,
`test_propose_action_then_confirm_actually_changes_the_database`,
`test_propose_action_rejected_leaves_the_database_unchanged`.
*Level:* http.

**ACT-06** Confirming another user's write and confirming an unknown
id are indistinguishable. *Representation:* 404, identical body.
*Evidence:* `test_confirm_with_wrong_user_and_unknown_id_are_identical`.
*Level:* http.

**ACT-07** A question that leads to a proposal returns the pending
write rather than an answer; questions are rate-limited, and a refused
question does not count against the limit; if the caller's permissions
change while a question is processed, it is refused rather than
answered. *Representation:* 202; 429; 409. *Evidence:*
`test_query_proposing_an_action_returns_202_with_a_reference`,
`test_query_is_rejected_429_once_the_rate_limit_is_reached`,
`test_a_rejected_query_is_not_itself_recorded_as_a_new_one`,
`test_query_refuses_if_permissions_changed_during_processing`.
*Level:* http, with a scripted model.

**ACT-08** A bulk action becomes one change per object, all carrying
the same mutation; an empty list writes nothing; a list beyond the
ceiling is refused with the count and an instruction to split it.
*Evidence:* `test_bulk_actions.py`. *Level:* internal.

---

## APPR — the approvals inbox

**APPR-01** A pending write is visible to its proposer — even after
they lose the grant — and to anyone who may decide it; to anyone else
it does not exist: absent from the list, 404 on detail, 404 on
confirm, and it survives their refusal. *Evidence:*
`test_a_colleague_sees_a_write_they_did_not_propose`,
`test_a_proposer_still_sees_their_write_after_losing_the_grant`,
`test_someone_with_neither_relationship_sees_nothing`,
`test_approvals_inbox.py::TestWhatEachPersonSees`. *Level:* http.

**APPR-02** The list never carries the changed values. *Evidence:*
`test_awaiting_writes_never_returns_the_changed_values`. *Level:*
http.

**APPR-03** The detail shows a reviewer the values they may read. A
field they may not read is **redacted, not omitted**: listed as
unreadable with both values null, and the write marked as having
redactions. The proposed value is gated as tightly as the current
one. *Evidence:* `test_a_field_the_reviewer_cannot_read_is_REDACTED_not_omitted`,
`test_the_proposed_value_is_gated_as_tightly_as_the_current_one`,
`test_approvals_inbox.py::TestTheDetail`. *Level:* http.

**APPR-04** Reading a detail never consumes the write. A refused
approval leaves it for someone else; a successful one consumes it.
*Representation:* a refused approval is 409. *Evidence:*
`test_reading_the_detail_does_not_consume_the_write`,
`test_a_refused_approval_leaves_the_write_for_someone_else`,
`test_approvals_inbox.py::TestFailedDecisions`. *Level:* http.

**APPR-05** The queue lists oldest first, counts identical proposals,
reports how many tasks a reviewer may decide and how many are
approved, and flags a write whose fields the configuration no longer
declares. *Evidence:* `test_the_oldest_proposal_is_listed_first`,
`test_approvals_inbox.py::TestIdenticalProposals`,
`test_the_inbox_says_what_a_reviewer_may_decide`,
`test_the_listing_reports_a_write_the_ontology_has_outrun`. *Level:*
http.

**APPR-06** An expired write leaves the queue and cannot be
confirmed. *Representation:* 404. *Evidence:*
`test_approvals_inbox.py::TestExpiry`. *Level:* oracle — needs
control of the clock.

---

## OWN — what belongs to one user

**OWN-01** Saved views belong to their author: others never see them
listed, deleting another's is a 404, a view on an unknown type is refused (404),
filters survive the round trip, and saving under an existing name
replaces it. *Evidence:* `test_saved_views.py`. *Level:* http.

**OWN-02** Notifications belong to their recipient: others never see
them, marking another's as seen is a 404, and the unseen count is the
recipient's own. *Evidence:* `test_notifications.py`. *Level:* http,
with a way to create a notification.

**OWN-03** Notes are scoped to one object, visible to the author's
role, and an empty note is refused. *Evidence:*
`test_a_note_can_be_written_and_read_back`,
`test_notes_are_scoped_to_one_object`,
`test_a_note_is_visible_to_the_authors_ROLE`,
`test_an_empty_note_is_rejected`. *Level:* http.

---

## ADMIN — administration and health

**ADMIN-01** Listing and creating users require user-management
rights; the listing carries role, security value and disabled state,
never credentials; a created user can log in. *Evidence:*
`test_list_users_returns_non_sensitive_metadata_only`,
`test_create_user_with_manage_users_grant_succeeds_and_new_user_can_log_in`.
*Level:* http.

**ADMIN-02** An administrator can see another user's visible schema —
the deliberate, admin-only inverse of uniform denial. An unknown user
is a 404. *Evidence:*
`test_visible_schema_debug_view_shows_what_the_target_user_can_see`,
`test_visible_schema_debug_view_for_unknown_user_is_404`. *Level:*
http.

**ADMIN-03** Configuration and source status name things, never
their contents: no connection details, paths, credentials or grants.
A failing source reports the kind of failure, never its message.
*Evidence:* `test_config_route_never_discloses_connection_details`,
`test_silos_route_reports_the_failure_KIND_not_the_message`. *Level:*
http.

**ADMIN-04** Source status maps each field to the source and column
that actually back it, including renamed join keys, and omits link
fields. *Evidence:* the `test_silos_route_*` tests. *Level:* http.

**ADMIN-05** Reloading configuration needs its own grant — user
management is not enough — and advances the generation. Starting a
sync needs a grant, returns before the sync finishes, and is refused
on a deployment that reads live. *Evidence:*
`test_reload_requires_its_own_grant_not_manage_users`,
`test_reload_with_the_grant_advances_the_generation`,
`test_mirror_sync_endpoint.py`. *Level:* http.

**ADMIN-06** Metrics need the deployment-management grant and report
request count, rate and error ratio, not saturation. Latency is not
pinned by any test. *Evidence:*
`test_metrics_need_manage_deployment`,
`test_metrics_report_red_and_not_saturation`. *Level:* http.

**ADMIN-07** Health needs no login, reports sources in aggregate
rather than by name, reveals nothing about the data, and degrades
rather than failing. A missing source, or one that exists but holds
no tables, is not healthy. *Evidence:*
`test_health_reports_silos_in_aggregate_not_by_name`,
`test_health_leaks_nothing_about_the_data`,
`test_health_reports_degraded_rather_than_failing`, and the two
source-file health tests. *Level:* http for the first three;
internal for the last.

**ADMIN-08** Data freshness needs a login but no grant, and reports
whether reads are live or from the mirror, with the last sync time.
*Evidence:* the `test_data_freshness_*` tests. *Level:* http.

---

## Out of scope, and why

**The 17 model-dependent tests.** An agent is stochastic, so a single
passing run proves little. These move to the evaluation harness,
which runs each case many times and fails on a regression in rate.

**Implementation internals.** Lake portability and teardown
(`test_teardown_and_rebuild.py`), write-log batching and chunking
(`test_write_log_batching.py`), concurrent request dispatch
(`test_serve_requests.py`), the pending-write time-to-live wiring and
request-metrics storage. These describe how Elysium is built, not
what it promises.

**Audit records** (`test_bulk_read_audit.py`): one record per bulk
read, every denial its own record, each naming the security
partitions touched. These are real guarantees, but a black-box suite
needs a way to read the audit trail. That interface is an open
question for Urshanabi.

**Deployment neutrality** (`test_a_different_deployment.py`): proves
the engine assumes nothing about one deployment's names. The suite
keeps this property by running against two unrelated fixtures.

---

## What Urshanabi owes beyond Elysium

From the external review of Elysium, which found these missing:

**EXTRA-01** Login fields are length-bounded, and an oversized value
is rejected before anything is persisted. In Elysium, ten
unauthenticated requests with 900 KB usernames grew its credentials
database from 100 KB to 16.3 MB. The test must assert that no record
was written, not only the status code.

**EXTRA-02** Failed-login records and expired sessions are removed
when they expire.

**EXTRA-03** Session tokens are stored only as hashes. *Level:*
internal — not observable over HTTP.

**EXTRA-04** QUERY-08's link half: ids are strings on every read path.

**EXTRA-05** DENY-09 in full: no truncation flag, count or total
derived from rows the caller cannot see.

---

## Before the suite exists

The suite needs a fixture deployment. Elysium's integration fixture
is written in its configuration format; whether Urshanabi reads the
same format, or the suite translates it, is undecided. Everything
above holds regardless of that choice.
