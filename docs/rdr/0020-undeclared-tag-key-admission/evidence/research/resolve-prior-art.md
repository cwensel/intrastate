Model: claude-sonnet-5

# A4: External prior art alignment with the carrier arm

## 1. SCXML (W3C)

Queries run (local corpora, all before web fallback):
- arc semantic, StateMachineRes: "_event.data undeclared error.execution datamodel"
- arc semantic, StateMachineRes: "SCXML data element assign undeclared location error.execution"
- arc semantic, StateMachineRes: "SCXML _event.data event payload not validated"
- arc semantic, StateMachineRes: "SCXML W3C recommendation datamodel data element error.execution"
- arc semantic, StateMachineLit: "SCXML datamodel error.execution undeclared location"

Result: negative — no quotable evidence in StateMachineRes (1189 md files) or StateMachineLit (27 PDFs). The only SCXML-adjacent hits were scxmlcc/state-machine-cat user manuals with generic attribute-reference boilerplate (`xmlns`, `datamodel` attribute enum, "Invoking scxmlcc"), none touching `<data>`/`<assign>`/`error.execution`/`_event.data` semantics. StateMachineLit returned an unrelated LLM-agent-failure paper (irrelevant, discarded).

Fell back to web (W3C SCXML Recommendation, https://www.w3.org/TR/scxml/), which fetched successfully (no 403):

Accepted citation — W3C SCXML Recommendation:
- Section 5.4 (`<assign>`): "If the location expression does not denote a valid location in the data model or if the value specified (by 'expr' or children) is not a legal value for the location specified, the SCXML Processor MUST place the error 'error.execution' in the internal event queue." — this is the CLOSED arm: assignment into a declared datamodel location is validated, and an undeclared/illegal target is rejected via `error.execution`.
- Section 5.10.1 (`_event`): "The receiving SCXML Processor SHOULD reformat this data to match its data model, but MUST NOT otherwise modify it." — the event's payload (`_event.data`) is external content the processor best-effort reformats but does not enforce a schema over; a failed reformat leaves the field blank and raises `error.execution` for the reformatting attempt itself, not for the payload's shape being wrong. The payload is carried, not validated against a declared vocabulary the way an `<assign>` target is.

Rejected branches: none opened beyond the above two fetches (first fetch already isolated both citations cleanly; second fetch to the full spec index was used only to re-confirm phrasing and section separation, not a new source).

Verdict for this system: ALIGNS — the SCXML split between declared `<data>` locations (validated, `error.execution` on violation) and the event payload (`_event.data`, carried/best-effort-reformatted, not schema-enforced) mirrors the declared-vocabulary-closed / undeclared-context-open split the carrier arm assumes.

## 2. Protobuf

Queries run:
- arc semantic, PapersFast: "protobuf unknown field preservation proto3" (crashed mid-run with a libc++ recursive_mutex error after returning results; results were still usable — a serialization-format survey paper, no unknown-field-retention text)
- arc semantic, DevRef: "protobuf unknown fields retained re-serialized proto3.5" — hit Designing Data-Intensive Applications (DDIA) ch.4 p.140-146, on field tags/schema evolution, but text says old code "can simply ignore that field" (forward-compat via ignoring), not explicit retain-and-re-emit-on-passthrough language
- arc semantic, DevRef: "protobuf forward compatibility unknown field ignored old code" — same DDIA passage, confirmed via full page read (pp.120-122 of the book); does not use "preserve" / "retain" / "unknown fields" terminology, so rejected as not a precise-enough match for the proto3.5-specific claim
- arc semantic, DevRef: "protobuf proto3 3.5 unknown field preservation restored deprecated" — same DDIA passage again, no 3.5-specific content

Result for local corpora: negative — no quotable evidence in PapersFast or DevRef precise enough for the specific "unknown-field retention" claim (DDIA describes the general tag-based forward-compatibility mechanism, which is related but does not state the proto3.5 preserve/re-emit behavior in citable terms).

Fell back to web:

Accepted citation — Protocol Buffers Language Guide (proto3), https://protobuf.dev/programming-guides/proto3/, "Unknown Fields" section: "Proto3 messages preserve unknown fields and include them during parsing and in the serialized output, which matches proto2 behavior."

Corroborating (not separately cited): WebSearch confirmed the historical detail that early proto3 releases dropped unknown fields by default and that this was reversed in the 3.5 line (restored to match proto2 behavior), consistent with the assumption's framing. GitHub issue #272 ("proto3 and unknown fields") and CHANGES.txt were checked for a directly quotable 3.5 changelog line; CHANGES.txt returned 404 at the guessed path, and issue #272's fetched content only surfaced the original post, not the comment thread with the version-specific confirmation — both rejected as non-quotable, but not needed since the current spec text above already directly asserts the behavior.

Verdict for this system: ALIGNS — unknown (undeclared-in-schema) fields are retained verbatim and passed through uninterpreted/unvalidated, while declared fields go through full type/shape handling. This is the same open/closed split as the carrier arm, applied to wire-format fields rather than CLI tag keys.

## 3. Kubernetes

Queries run:
- arc semantic, DevRef: "kubernetes labels annotations arbitrary metadata not validated by API server"
- arc semantic, DevRef: "kubernetes label value syntax 63 characters alphanumeric no semantic meaning to Kubernetes"

Result: positive on first query — hit Mastering Kubernetes (DevRefGit/Books2/dev-ops/kubernetes/mastering-kubernetes.pdf), pages 35-36 (book pages 6-7), read in full.

Accepted citation — Mastering Kubernetes (Packt), "Label" / "Annotation" sections, pp.6-7:
- Label (syntax, closed/validated): "The label key must adhere to a strict syntax. It has two parts: prefix and name. The prefix is optional. If it exists then it is separated from the name by a forward slash (/) and it must be a valid DNS sub-domain. The prefix must be 253 characters long at most. The name is mandatory and must be 63 characters long at most. Names must start and end with an alphanumeric character (a-z, A-Z, 0-9) and contain only alphanumeric characters, dots, dashes, and underscores. Values follow the same restrictions as names."
- Annotation (open/uninterpreted): "Annotations let you associate arbitrary metadata with Kubernetes objects. Kubernetes just stores the annotations and makes their metadata available. Unlike labels, they don't have strict restrictions about allowed characters and size limits."

This is a close structural match to the assumption but not an exact one: Kubernetes validates label KEY and VALUE syntax (charset/length) uniformly, and the key vocabulary itself is open (any object may declare any label key — there is no pre-declared "tag table" naming which keys are legal) while the API server never interprets what a label's value *means*. Annotations go further and drop even the syntax restrictions. In both cases, "the VALUE's meaning is uninterpreted by the API server" — matching the RDR's framing exactly — but Kubernetes does not have a declared/undeclared split at the KEY level the way the RDR's tag table does (all keys are equally undeclared from the API server's point of view; the closed vocabulary in K8s lives in each consumer controller, not in the API server's schema). So the match is on the "syntax may be checked; value semantics are never validated by the platform" axis, not on a "declared key -> validated, undeclared key -> carrier" axis.

Rejected branches: docker-orchestration.pdf and docker-cookbook-solutions-examples.pdf hits (label usage examples, `kubectl -l` selector syntax) — tutorial-level, no validation-semantics content.

Verdict for this system: ALIGNS (partially) — confirms the "syntax validated, value semantics uninterpreted by the platform" half of the carrier arm strongly; does not independently confirm a declared-vocabulary-vs-undeclared-key split, since Kubernetes' label/annotation model has no declared key table at the API-server level (that's push-down to consumers, which is actually closer to the RDR's tag-table model living in the loaded state-machine model rather than the CLI engine — an analogous layering, not an identical mechanism).

## Overall

External prior art ALIGNS with the carrier arm. SCXML gives the cleanest structural match (declared datamodel location = validated/closed, event payload = carried/open). Protobuf gives a clean match at the wire-format level (declared field = typed, unknown field = retained verbatim). Kubernetes confirms the "value semantics uninterpreted by the platform" principle but its declared/undeclared split is at a different layer (consumer controllers, not the API server) than the RDR's tag-table split (declared in the loaded model). All three citations were opened and read; no citation taken from memory. This assumption is non-fatal per the RDR — the in-repo anchors carry the design decision regardless of this finding.
