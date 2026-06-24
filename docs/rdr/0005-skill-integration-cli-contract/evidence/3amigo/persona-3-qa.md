# Persona 3 — QA / Tester

1. I cannot write a precise pass/fail test for `flow next` because the RDR does
   not define the expected JSON data fields or text content for legal outcomes
   and conditional summaries. The validation scenario says both modes report the
   same data, but not what "same" means.

2. I cannot write complete refusal tests because the RDR requires stable
   `CLIError.Code` values but does not specify the canonical code strings or the
   exit-code group for each resolver/accessor failure class.

3. I cannot write parser tests for command inputs because the RDR leaves tag
   syntax, duplicate handling, structured input file format, and model/flow
   selection as broad prose.

4. I cannot write pass/fail tests for accessor role binding because
   `--artifact state:RDR_FILE` is illustrative only. The test cannot tell which
   malformed bindings should fail before the accessor layer and which should
   pass through as accessor errors.

5. I cannot write the `set-state` read-back test without knowing whether
   requested tags, expected read-back tags, and resolver-planned writes are all
   represented by the same CLI flag or distinct request fields.
