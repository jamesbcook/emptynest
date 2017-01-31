package emptynest

import "plugin"

// EncoderPlugin is used to decode and encode byte slices.
type EncoderPlugin struct {
	Name   func() string
	Encode func([]byte) ([]byte, error)
	Decode func([]byte) ([]byte, error)
}

// BuildEncoderChain takes a slice of filenames for
// plugins and returns a chain that can be used for execution.
func BuildEncoderChain(files []string) ([]EncoderPlugin, error) {
	var chain []EncoderPlugin
	for _, f := range files {
		p, err := plugin.Open(f)
		if err != nil {
			return chain, err
		}
		namefunc, err := p.Lookup("Name")
		if err != nil {
			return chain, err
		}
		encodefunc, err := p.Lookup("Encode")
		if err != nil {
			return chain, err
		}
		decodefunc, err := p.Lookup("Decode")
		if err != nil {
			return chain, err
		}
		chain = append(chain, EncoderPlugin{
			Name:   namefunc.(func() string),
			Encode: encodefunc.(func([]byte) ([]byte, error)),
			Decode: decodefunc.(func([]byte) ([]byte, error)),
		})
	}
	return chain, nil
}
