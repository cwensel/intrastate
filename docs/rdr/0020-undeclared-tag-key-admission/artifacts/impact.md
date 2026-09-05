# Impact — 0020-undeclared-tag-key-admission

families: 2
rows: 6
records: 0005 0008
literals: `looksArray` `is not set-valued` `was given an empty value` `this is not the place that makes it` `An undeclared tag whose value is an array is refused` `every set-valued key the caller passes must be declared`
repo: /Users/cwensel/sandbox/newcoinc/intrastate/.claude/worktrees/rdr-0020 @d97f2c5
convention: go

## TestAdv0008 (5 tests, 2 files)
| test | file | arm |
|---|---|---|
| TestAdv0008_LintCarriesEveryLoadCategoryAsAFinding | internal/cli/reserved_key_adv_0008_test.go | record:0008 |
| TestAdv0008_LintCarriesReservedKeyPayloadOnTheWire | internal/cli/reserved_key_adv_0008_test.go | record:0008 |
| TestAdv0008_LintNeverNamesTheReservedKeyAsTheRequiredName | internal/cli/reserved_key_adv_0008_test.go | record:0008 |
| TestAdv0008_LintSurfacesTheNearMissAdvisory | internal/cli/reserved_key_adv_0008_test.go | record:0008 |
| TestAdv0008_DoublyBreachingModelReportsOneStableDirection | internal/table/reserved_key_adv_0008_test.go | record:0008 |

## TestMVV0008 (1 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestMVV0008_ReservedKeyOwnershipEndToEnd | internal/table/mvv_0008_test.go | record:0008 |
