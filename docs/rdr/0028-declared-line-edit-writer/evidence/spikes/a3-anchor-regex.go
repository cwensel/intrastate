//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"regexp/syntax"
	"sort"
	"strings"
)

// Candidate anchors under test.
// Anchor A: RDR record Metadata Status bullet (top-level, zero indent).
//
//	Distinguishes from nested "- **Status**:" lines inside Assumption blocks
//	(which are indented by at least one space) by pinning ^- at column 0.
var statusAnchor = regexp.MustCompile(`^- \*\*Status\*\*: (.+)$`)

// Anchor B (generic): README index row, any record number.
var readmeRowGeneric = regexp.MustCompile(`^\| \[\d{4}\]\([^)]+\)[^|]*\| [^|]*\| ([^|]+) \|`)

// buildReadmeRowForID returns the per-record anchor, built the way the RDR
// tool would build it from {tag.nnnn} substituted regexp-quoted into a
// template.
func buildReadmeRowForID(id string) *regexp.Regexp {
	pat := `^\| \[` + regexp.QuoteMeta(id) + `\]\([^)]+\)[^|]*\| [^|]*\| ([^|]+) \|`
	return regexp.MustCompile(pat)
}

func hasLookaround(pat string) (bool, string) {
	// RE2/Go regexp/syntax will fail to parse lookaround since RE2 doesn't
	// support it; but as an explicit documented check, scan for the literal
	// PCRE lookaround tokens that a human might have been tempted to write,
	// and confirm syntax.Parse succeeds under POSIX-compatible RE2 flags.
	for _, tok := range []string{"(?=", "(?!", "(?<=", "(?<!"} {
		if strings.Contains(pat, tok) {
			return true, tok
		}
	}
	_, err := syntax.Parse(pat, syntax.Perl)
	if err != nil {
		return true, "parse-error: " + err.Error()
	}
	return false, ""
}

type fileResult struct {
	path    string
	count   int
	lines   []string
	lineNos []int
}

func scanFile(path string, anchor *regexp.Regexp) (fileResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return fileResult{}, err
	}
	lines := strings.Split(string(data), "\n")
	res := fileResult{path: path}
	for i, ln := range lines {
		if anchor.MatchString(ln) {
			res.count++
			res.lines = append(res.lines, ln)
			res.lineNos = append(res.lineNos, i+1)
		}
	}
	return res, nil
}

func main() {
	rdrDir := os.Args[1] // docs/rdr directory
	readmePath := filepath.Join(rdrDir, "README.md")

	fmt.Println("--- Lookaround check ---")
	for name, pat := range map[string]string{
		"statusAnchor":     statusAnchor.String(),
		"readmeRowGeneric": readmeRowGeneric.String(),
	} {
		la, tok := hasLookaround(pat)
		fmt.Printf("%s: lookaround=%v tok=%q pattern=%s\n", name, la, tok, pat)
	}

	fmt.Println("\n--- Status anchor scan over all *.md record files ---")
	var mdFiles []string
	entries, err := os.ReadDir(rdrDir)
	if err != nil {
		fmt.Println("ERROR reading dir:", err)
		os.Exit(1)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".md") {
			mdFiles = append(mdFiles, filepath.Join(rdrDir, e.Name()))
		}
	}
	sort.Strings(mdFiles)

	scanned := 0
	var badFiles []fileResult
	var zeroFiles []string
	for _, f := range mdFiles {
		base := filepath.Base(f)
		if base == "README.md" {
			continue // handled separately below
		}
		scanned++
		res, err := scanFile(f, statusAnchor)
		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}
		fmt.Printf("%-70s count=%d\n", base, res.count)
		if res.count == 0 {
			zeroFiles = append(zeroFiles, base)
		}
		if res.count != 1 {
			badFiles = append(badFiles, res)
		}
	}

	fmt.Printf("\nTotal record files scanned: %d\n", scanned)
	fmt.Printf("Files where Status-anchor count != 1: %d\n", len(badFiles))
	for _, r := range badFiles {
		fmt.Printf("  BAD FILE: %s count=%d\n", r.path, r.count)
		for i, ln := range r.lines {
			fmt.Printf("    line %d: %q\n", r.lineNos[i], ln)
		}
	}
	if len(zeroFiles) > 0 {
		fmt.Println("Files with ZERO matches (informational, not necessarily a defect if not yet a full record):")
		for _, z := range zeroFiles {
			fmt.Println("  ", z)
		}
	}

	fmt.Println("\n--- README generic row anchor scan ---")
	data, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Println("ERROR reading README:", err)
		os.Exit(1)
	}
	readmeLines := strings.Split(string(data), "\n")
	genCount := 0
	for i, ln := range readmeLines {
		if readmeRowGeneric.MatchString(ln) {
			genCount++
			m := readmeRowGeneric.FindStringSubmatch(ln)
			fmt.Printf("  line %d matched, status-capture=%q :: %q\n", i+1, strings.TrimSpace(m[1]), ln)
		}
	}
	fmt.Printf("Generic row anchor total matches in README: %d (expect == number of linked rows)\n", genCount)

	fmt.Println("\n--- README per-record anchor scan (parameterized nnnn) ---")
	ids := []string{"0001", "0002", "0006", "0026", "0027", "0028", "9999"}
	for _, id := range ids {
		anchor := buildReadmeRowForID(id)
		res, err := scanFile(readmePath, anchor)
		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}
		fmt.Printf("id=%s count=%d\n", id, res.count)
		for i, ln := range res.lines {
			fmt.Printf("    line %d: %q\n", res.lineNos[i], ln)
		}
		if res.count != 1 {
			fmt.Printf("  ANOMALY: id=%s expected count=1, got %d\n", id, res.count)
		}
	}

	fmt.Println("\n--- Qualifier capture + rewrite demo ---")
	// Use the RDR engine's own testdata fixture which has a real
	// bracket-qualified Status line, since no live intrastate docs/rdr
	// record currently carries a top-level qualified Status bullet.
	qualFile := os.Args[2]
	qdata, err := os.ReadFile(qualFile)
	if err != nil {
		fmt.Println("ERROR reading qualifier fixture:", err)
		os.Exit(1)
	}
	qlines := strings.Split(string(qdata), "\n")
	fmt.Println("Fixture:", qualFile)
	for i, ln := range qlines {
		if statusAnchor.MatchString(ln) {
			fmt.Printf("  BEFORE (line %d): %q\n", i+1, ln)
			// Anchor with two groups: value-without-qualifier or full tail,
			// and a separate qualifier-only capture for demonstration.
			qualAnchor := regexp.MustCompile(`^- \*\*Status\*\*: (\S+)( \[.*)?$`)
			m := qualAnchor.FindStringSubmatch(ln)
			if m == nil {
				fmt.Println("  NO MATCH on qualifier-split anchor (qualifier likely wraps to next line — see note)")
				continue
			}
			fmt.Printf("  matched: word=%q qualifierTail=%q\n", m[1], m[2])
			after := qualAnchor.ReplaceAllString(ln, `- **Status**: Final${2}`)
			fmt.Printf("  AFTER  (regexp.Expand via ReplaceAllString): %q\n", after)
		}
	}

	fmt.Println("\n--- Single-line synthetic qualifier demo (regexp.Expand) ---")
	synthetic := "- **Status**: Draft [joint decision → JDR 0003 §D1]"
	fmt.Printf("  BEFORE: %q\n", synthetic)
	qualAnchor := regexp.MustCompile(`^- \*\*Status\*\*: (\S+)( \[.*\])?$`)
	m := qualAnchor.FindStringSubmatchIndex(synthetic)
	if m == nil {
		fmt.Println("  NO MATCH")
	} else {
		dst := qualAnchor.Expand(nil, []byte(`- **Status**: Final${2}`), []byte(synthetic), m)
		fmt.Printf("  AFTER  (regexp.Expand, ${2} preserves qualifier): %q\n", string(dst))
	}
}
