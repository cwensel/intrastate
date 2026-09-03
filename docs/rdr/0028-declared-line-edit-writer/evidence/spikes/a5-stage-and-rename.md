Model: claude-sonnet-5

# A5 spike: stage-and-rename with original file's mode copied

Assumption under test (RDR cli/0028, A5, verbatim via `rdr inspect --select
0028:A5 0028`): staging a new file beside the target and renaming over it —
with the ORIGINAL file's mode copied (unlike `flowbind.go::save`, which
writes a fixed 0600) — is acceptable for a git-tracked markdown file: git
sees a CONTENT change only, and no consumer depends on the inode or a hard
link. The assumption text also says the spike settles symlinked targets
(C3 resolves them) and xattr loss.

All mutation work ran in a scratch git repo under `/tmp/a5-spike`; nothing
under `newcoinc/` was mutated except this evidence directory.

---

## 1. Baseline: `flowbind.go::save`'s existing discipline

`/Users/cwensel/sandbox/newcoinc/intrastate/internal/cli/flowbind/flowbind.go`,
function `save` (lines 139-166):

```go
tmp, err := os.CreateTemp(filepath.Dir(path), ".flow-artifact-*")
...
if err := os.Chmod(staged, 0o600); err != nil {
    return err
}
return os.Rename(staged, path)
```

Confirmed: `save` stages a temp file beside the target
(`os.CreateTemp(filepath.Dir(path), ...)`), then unconditionally
`os.Chmod(staged, 0o600)` — a FIXED mode, never read from the original file
— then `os.Rename`s over the target. The "fixed-0600" claim in A5 is
accurate as written.

---

## 2. Git behavior: mode-copied vs fixed-0600

Scratch repo `/tmp/a5-spike`, `target.md` committed at 0644.

**RUN A/B (0644 baseline, non-executable):** stage-and-rename with
`-copy-mode` (644 -> 644) and with `-fixed-0600` (644 -> 0600) both produced
identical git output:

```
 M target.md
```
```
 target.md | 2 +-
 1 file changed, 1 insertion(+), 1 deletion(-)
```
`git diff --summary` was EMPTY in both cases — no mode line, either way.

This surfaced a nuance the assumption's phrasing elides: git's tracked mode
is binary — git mode `100644` vs `100755` (the executable bit only). Full
POSIX permission bits (0644 vs 0600) both collapse to git mode `100644`, so
for a non-executable markdown file the fixed-0600-vs-copied-mode distinction
is invisible to git regardless of whether the mode is copied.

**RUN C/D/E (0755 baseline, executable)** — re-ran the contrast starting
from an executable file to surface the real git-visible effect:

- RUN D (`-copy-mode`, 755 -> 755): `git diff --summary` empty; `git diff
  --raw` shows `:100755 100755 ... M	target.md` — mode unchanged, content
  only.
- RUN E (`-fixed-0600`, 755 -> 600): `git diff --summary` shows
  `mode change 100755 => 100644 target.md`; `git diff --raw` shows
  `:100755 100644 ... M	target.md`.

This is the real contrast the assumption is pointing at: mode-copying
preserves the executable bit and keeps git's diff content-only; a fixed
0600 (as `flowbind.go::save` does today) silently drops the executable bit
and produces a git-visible, spurious mode-change diff line WHEN the original
file was executable. For a plain non-executable markdown file, both
behaviors are git-silent — the risk A5 describes only bites executable
targets. This is worth noting to the author as a scope correction: A5 as
written should be read as "content-only for git" being true in ALL cases for
non-executable targets, and true for executable targets ONLY when the mode
is copied.

---

## 3. Inode / hard link

Before: `target.md` inode `162352375`, mode 0600, `nlink=2` after
`ln target.md hardlink-to-target.md` (hardlink also reports inode
`162352375`).

After stage-and-rename (mode copied): `target.md` inode changed to
`162357409`. The hard link still reports inode `162352375` and its content
(`cat hardlink-to-target.md`) is STALE — it shows the content from BEFORE
the rename (`# content E fixed-0600`), while `target.md` shows the new
content (`# content F hardlink-test`).

Confirms the assumption's premise directionally but sharpens the
consequence: the inode always changes under stage-and-rename (as expected —
`os.Rename` swaps the directory entry to a new inode), and any hard link to
the artifact silently orphans, frozen at the pre-rename content. The
assumption's "no consumer depends on... a hard link" is a claim that must
hold structurally (nothing in this codebase hard-links RDR/markdown
artifacts) — this spike cannot prove that universally, only that IF a hard
link existed, it would desync silently with no error raised.

---

## 4. Symlink target

**4a — resolved first (what C3 buys):** created `symlink-to-target.md ->
target.md`. Ran stage-and-rename against the symlink PATH but with
`filepath.EvalSymlinks` resolving it first (as implemented in the Go
program). Result: rename landed on the real file. `ls -la
symlink-to-target.md` after still shows `lrwxr-xr-x ... symlink-to-target.md
-> target.md` — the symlink survives untouched, and both
`symlink-to-target.md` and `target.md` read back the new content.

**4b — NOT resolved (what C3 prevents), for contrast:** renamed a staged
file directly onto the symlink PATH without resolving it (`mv $tmp
symlink-to-target.md`). Result: `test -L symlink-to-target.md` now FALSE —
the symlink is gone, replaced by a plain regular file holding the new
content. The original `target.md` is left completely untouched at its
PREVIOUS content, now orphaned from that path entirely.

This is exactly the failure mode contract C3's "symlinks resolved" clause
exists to prevent: an un-resolved rename over a symlink path destroys the
symlink and silently stops writing to the real target.

---

## 5. xattrs

Set `user.test=somevalue` on `target.md` via `xattr -w`. Confirmed present
via `xattr -l` (`user.test: somevalue`). After stage-and-rename (mode
copied): `xattr -l target.md` produced NO output — the xattr did not
survive.

Confirms the stated consequence: `os.CreateTemp` creates a fresh inode with
no xattrs, so any xattr set on the original file is unconditionally lost on
every stage-and-rename write, independent of whether the mode is copied.

---

## 6. Go program

Wrote `a5-stage-and-rename.go` (copied into this evidence directory) to
exercise the actual `os.Stat` / `os.CreateTemp` / `os.Chmod` /
`filepath.EvalSymlinks` / `os.Rename` path rather than shell `mv`. Flags
`-copy-mode` (reads the original file's mode via `os.Stat` and chmods the
staged file to match) and `-fixed-0600` (mirrors `flowbind.go::save`'s fixed
`os.Chmod(staged, 0o600)`) were both exercised above. Built cleanly with `go
build` under `go1.26.6 darwin/arm64`.

---

## Verdict

The core assumption holds for the common case (non-executable markdown
file): git sees content-only regardless of mode-copy, the inode changes
every time (expected, unavoidable with rename-based atomic writes), hard
links desync silently (structural risk, not this spike's to rule out
project-wide), symlink targets are handled correctly ONLY if resolved first
(confirmed C3's requirement is necessary, not optional), and xattrs are lost
unconditionally on every write (confirmed, independent of mode-copy).

One correction to the assumption's framing: git's mode-change visibility is
executable-bit-only. The "fixed-0600 causes a spurious mode diff" risk that
motivates copying the original mode is real ONLY when the original file is
executable — for the everyday non-executable `.md` artifact this RDR
targets, `flowbind.go::save`'s existing fixed-0600 approach is ALREADY
git-silent on mode. Copying the mode remains strictly better (it is also
correct for the rarer executable-target case and preserves permissions for
non-git consumers), but the git-diff-cleanliness argument specifically is
narrower than A5 implies.
