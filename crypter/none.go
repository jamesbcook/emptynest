package crypter

// None doesn't do any crypto at all!
type None struct {
}

func (n *None) Seal(key, data []byte) ([]byte, error) {
	return data, nil
}

func (n *None) Unseal(key, data []byte) ([]byte, error) {
	return data, nil
}
