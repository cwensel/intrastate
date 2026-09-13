Model: claude-sonnet-5

# Prior art: in-band vs out-of-band abstraction/unknown markers

Corpora searched: DevRef (arc search text/semantic), PapersFast (arc search semantic), StateMachineLit, StateMachineRes. Web sources excluded per task instructions (corpora are the source of record; publisher sites 403).

---

## 1. Hazard of in-band sentinels colliding with legitimate data

**SQL Antipatterns (Bill Karwin), Chapter 14 "Fear of the Unknown," section 14.2 "Antipattern: Use Null as an Ordinary Value, or Vice Versa"** — DevRefGit/Database-Books/db(33).pdf, pp. 163-167.

- p.163-164: "Null is not the same as zero... Null is not the same as a string of zero length... Null is not the same as false." Establishes NULL as a distinct marker rather than a member of the value domain (foreshadows sub-question 4: presence as its own type, not a distinguished value).
- p.165-166 (worked example, "Avoiding the Issue"): shows a team disallowing NULL and instead choosing `-1` as a sentinel for "unspecified" in a numeric column. The book walks through the resulting failures: `-1` throws off `SUM()`/`AVG()` and must be excluded via special-case `WHERE hours <> -1` clauses; the sentinel choice becomes column-specific and undocumented, adding "a lot of meticulous and unnecessary work"; and a foreign-key column has no safe non-null sentinel at all without inventing a fake "unassigned" row in the referenced table.
- **p.167, "How to Recognize the Antipattern," second bullet — the direct, citable statement of the general hazard**: *"It's common that we can't use the string we've been using to represent unknown in the [table], so we need to have a meeting to discuss what new special value we can use... This is a likely consequence of assigning a special flag value that could be a legitimate value in your column's domain. Eventually, you may find you need to use that value for its literal meaning instead of its flag meaning."* This is the canonical statement of sub-question 1's hazard: an in-band sentinel drawn from the same domain as legitimate data will eventually collide with a real value that needs that literal meaning.
- p.167 sidebar "Are Nulls Relational?": notes Codd introduced NULL specifically to signify missing data as a marker outside the value domain, precisely to avoid this collision class.

**SQL for Smarties, 5th ed. (Joe Celko)**, DevRefGit/pdf, pp. 295-296 — corroborates from the language-design side: SQL's NULL is deliberately "no data type" in relational theory (must be given storage/type accommodation in implementation) and drives three-valued logic (`TRUE`/`FALSE`/`UNKNOWN`) rather than folding "unknown" into the two-valued domain as a special string or number.

**SQL Database Programmer's Handbook**, DevRefGit/Database-Books/db(40).pdf, p.47 (Ch.5 "SQL NULL," "The Null of It All") — the maximally quotable framing for sub-question 1's title claim: *"A NULL is not a value; it is a marker for a value that is missing."* Confirms searches (`2 < NULL`, `assigned_to = NULL`) as well as comparisons return UNKNOWN, not a truth value, per db(4).pdf p.124 (three-valued logic on NULL-involving expressions).

## 2. When is an in-band sentinel acceptable?

No corpus document states the "closed domain + reserved/escaped at every ingress" condition explicitly as a named rule. The closest material is negative evidence pointing the same direction: SQL Antipatterns p.166-167 shows that even a single homogeneous numeric/FK column (about as "closed" a domain as exists) fails to make an in-band sentinel safe, because the domain is not actually reserved — nothing prevents `-1` or the placeholder FK row from later being a legitimate value, and nothing enforces exclusion at every read site (aggregate functions, equality predicates) uniformly. The corpora support the negative claim (in-band sentinels are unsafe absent enforced reservation) but do not supply an explicit positive statement of the condition under which they are sound. Treating sub-question 2 as **not found** rather than stretching the p.166-167 material into an affirmative rule.

## 3. Abstract interpretation TOP (⊤) element serialization

Searched PapersFast with multiple phrasings ("abstract domain top element unknown value representation," "widening operator top element lattice fixpoint," "abstract interpretation Cousot abstract domain unknown over-approximation"). All returns were irrelevant: table-of-contents/index pages, unrelated papers (LLM alignment, requirements engineering, sketch-based query monitoring), and page-number noise with no discussion of abstract-domain TOP-element representation or serialization. StateMachineLit and StateMachineRes were also searched ("top element unknown state abstraction serialization") and returned only unrelated state-machine library docs (JavaScript State Machine, xgrammar, ragel-colm-suite) with no lattice/abstract-interpretation content.

**Result: INCOMPLETE for this sub-question.** PapersFast does not appear to contain primary abstract-interpretation literature (Cousot & Cousot, or lattice-theoretic program analysis papers) that is retrievable by these queries, and StateMachineLit/StateMachineRes are scoped to state-machine engineering rather than abstract interpretation. No claim is made here — an honest gap, not an analogy.

## 4. Protocol/schema evolution: presence vs. sentinel

**Designing Data-Intensive Applications (Martin Kleppmann), Chapter 4 "Encoding and Evolution"** — DevRefGit/Books/Designing Data Intensive Applications.pdf (duplicated as Database-Books/db(46).pdf), pp.120-124.

- p.120-121: Protobuf/Thrift field presence is carried **out-of-band** in the wire format itself — "If a field value is not set, it is simply omitted from the encoded record" — rather than by writing an in-band marker into the value's own domain. `required` vs `optional` makes no difference to encoding ("nothing in the binary data indicates whether a field was required"); presence/absence is structural (tag appears or it doesn't), not a distinguished value written into the field's type.
- p.121: schema evolution rule for backward/forward compatibility explicitly ties safety to fields being **optional or carrying a default**, i.e., absence is a first-class structural state, not something the value domain has to represent by convention.
- p.122-124 (Avro): reinforces the same lesson from a different angle — Avro's `union { null, long } favoriteNumber = null` makes "unknown/absent" an explicit union branch (a distinct type alternative), not a magic long value; the reader/writer schema-resolution process (Figure 4-6) fills genuinely missing fields from a declared default rather than relying on the writer having encoded a sentinel long. This directly supports sub-question 4: the general lesson is to represent absence/unknown as a distinct type (a union arm, an omitted tag) rather than a distinguished in-domain value.

No explicit JSON Schema `null`-vs-absent discussion was found in the corpora under these queries; the Avro/Protobuf material stood in as the schema-evolution evidence base for this sub-question.

---

## Summary of what was and wasn't found

Found, with citations: (1) the general collision hazard, strongly supported (SQL Antipatterns, SQL for Smarties, SQL Database Programmer's Handbook); (4) the presence-vs-sentinel lesson from protocol/schema evolution, strongly supported (DDIA Ch.4, Protobuf and Avro). Not found: (2) an explicit statement of the "closed domain + reserved/escaped everywhere" sufficient condition — only negative/circumstantial support. Not found: (3) abstract-interpretation TOP-element serialization treatment — PapersFast and StateMachineLit/Res returned no relevant material under multiple query phrasings.
