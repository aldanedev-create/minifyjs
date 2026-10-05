package source

// Location is the combination of a Position and the file it refers
// to. It is used by diagnostics that need to report a single point
// rather than a range (e.g. "undefined identifier" points at the
// identifier, not at the whole expression).
type Location struct {
	File     string
	Position Position
}

// IsValid reports whether the location has been initialized.
func (l Location) IsValid() bool {
	return l.Position.IsValid()
}

// String returns "file:line:col".
func (l Location) String() string {
	if l.File == "" {
		return l.Position.String()
	}
	return l.File + ":" + l.Position.String()
}