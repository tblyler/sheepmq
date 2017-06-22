package lease

// Leaser denotes whether a given lease is still valid
type Leaser interface {
	// Update and check the leaser
	Valid() bool

	// Check the leaser without any sort of updates
	Check() bool
}
