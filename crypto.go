package emptynest

import "plugin"

// CryptoPlugin is used to decode and encode byte slices.
type CryptoPlugin struct {
	Name func() string
	Open func(key, data []byte) ([]byte, error)
	Seal func(key, data []byte) ([]byte, error)
}

// BuildCryptoChain takes a slice of filenames for
// plugins and returns a chain that can be used for execution.
func BuildCryptoChain(files []string) ([]CryptoPlugin, error) {
	var chain []CryptoPlugin
	for _, f := range files {
		p, err := plugin.Open(f)
		if err != nil {
			return chain, err
		}
		namefunc, err := p.Lookup("Name")
		if err != nil {
			return chain, err
		}
		openfunc, err := p.Lookup("Open")
		if err != nil {
			return chain, err
		}
		sealfunc, err := p.Lookup("Seal")
		if err != nil {
			return chain, err
		}
		chain = append(chain, CryptoPlugin{
			Name: namefunc.(func() string),
			Open: openfunc.(func([]byte, []byte) ([]byte, error)),
			Seal: sealfunc.(func([]byte, []byte) ([]byte, error)),
		})
	}
	return chain, nil
}
