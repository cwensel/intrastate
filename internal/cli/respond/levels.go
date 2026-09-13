package respond

// Levels returns the stderr advisory `level` vocabulary (`0029:C4`). The
// set is frozen: no member is added or removed within a major.
//
// It lives in its own file so the enumeration seam's members are readable
// as the accessor's own literals, unmixed with any other vocabulary's.
func Levels() []string { return []string{"note", "warning"} }
