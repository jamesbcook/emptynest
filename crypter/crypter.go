package crypter

// Crypter is an interface for an encryption scheme.
type Crypter interface {
	Seal(key, data []byte) ([]byte, error)
	Unseal(key, data []byte) ([]byte, error)
}

var Map = map[int]Crypter{}

func init() {
	Map[0x00] = new(None)
}
