# RDR 0025 — Stage 2 prior-art record

Class question: how do declarative tools carry a declared external command so
its executed authority is reviewable statically? Instance question: what does
each peer state-machine/workflow tool do for declared command invocation?

## Accepted citations

- **PA-1 (hazard — shell-string carrier needs an after-the-fact linter).**
  awf-cli (state-machines corpus, `repos/awf-cli`), security-validator plugin
  README §Command Injection Warnings: "Detects unquoted variable references in
  commands / Warns about potential shell injection vulnerabilities"; §Timeout
  Enforcement: "Warns when command steps lack timeout configuration". The
  carrier itself (`docs/user-guide/workflow-syntax.md`) is a shell string with
  Go-template interpolation: `command: | echo "Processing {{.inputs.file}}"`.
  Found via corpus query "argv array versus shell string command declaration
  templating substitution safety" (StateMachineRes).
- **PA-2 (full-templating end of the spectrum).** ms-conductor
  (`repos/ms-conductor`), `references/yaml-schema.md` §Template Syntax: full
  Jinja2 over workflow inputs (`{{ workflow.input.param_name }}`,
  conditionals) — the executed text is a runtime function of input. Found via
  same corpus queries.
- **PA-3 (no-shell argv form, by construction).** Docker, Dockerfile reference
  (docs.docker.com/reference/dockerfile), exec form: "Using the exec form
  doesn't automatically invoke a command shell. This means that normal shell
  processing, such as variable substitution, doesn't happen." Example given:
  `RUN [ "echo", "$HOME" ]` won't substitute `$HOME`.
- **PA-4 (closed placeholder vocabulary when per-invocation data must reach a
  declared command).** git, git-difftool docs (git-scm.com/docs/git-difftool),
  `difftool.<tool>.cmd`: "The specified command is evaluated in shell with the
  following variables available: `LOCAL` is set to the name of the temporary
  file containing the contents of the diff pre-image and `REMOTE` is set to
  the name of the temporary file containing the contents of the diff
  post-image." A tool-defined, closed variable set carrying file paths — never
  open interpolation of content.
- **PA-5 (argv list + no shell + JSON string-map stdin/stdout protocol).**
  Terraform `external` data source (hashicorp/terraform-provider-external,
  `docs/data-sources/external.md`): program is "A list of strings, whose first
  element is the program to run and whose subsequent elements are optional
  command line arguments"; "Terraform does not execute the program through a
  shell"; "The program must read all of the data passed to it on `stdin`, and
  parse it as a JSON object"; "The program must then produce a valid JSON
  object on `stdout`" whose values "will always be strings".

## Instance disposition

Peer state-machine/workflow tools split two ways: shell-string templating in
the model (awf-cli, ms-conductor — PA-1/PA-2) or host-language callbacks
(xstate, temporal, restate, inngest, qmuntal-stateless in the same checkout
set — the shape RDR 0004 Alternative 5 already rejected). No peer in the
family carries a statically lint-validatable command authority bound; the
established bounded forms come from outside the family (PA-3, PA-4, PA-5).

## Demoted (could not quote from source within budget)

- Kubernetes `ExecAction` "not run inside a shell" claim: three fetch attempts
  (task page, raw `types.go`, pkg.go.dev) returned truncated content; not
  quote-confirmed, so not load-bearing here. Demoted to a Resolve assumption
  (Method: Source Search over `k8s.io/api/core/v1` `ExecAction.Command` doc
  comment). PA-3/PA-5 carry the same class claim and are quote-confirmed.

## Rejected branches

- Corpus query "SCXML invoke external service declared in state chart model"
  (StateMachineRes): top hits were format-conversion docs
  (state-machine-cat, scxmlcc) with no `<invoke>` authority-bound content —
  no useful coverage; not widened per budget.
