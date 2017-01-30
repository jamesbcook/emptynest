package emptynest

// Key stores an encryption key by ID.
type Key struct {
	ID  int `storm:"id,increment"`
	Key []byte
}
