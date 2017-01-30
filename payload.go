package emptynest

import (
	"io/ioutil"
	"path/filepath"
	"plugin"
)

// Payload is executed on approved hosts.
type Payload struct {
	ID   int    `storm:"id,increment"`
	Name string `storm:"unique"`
	Kind string
	Data []byte
}

// PayloadPlugin is used to generate a payload.
type PayloadPlugin struct {
	ID       func() int
	Generate func([]string) ([]byte, error)
	Help     func() string
	Name     func() string
	String   func([]byte) string
}

// PayloadMap ...
func PayloadMap(directory string) (map[string]PayloadPlugin, error) {
	payloadMap := make(map[string]PayloadPlugin)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return payloadMap, err
	}

	for _, f := range files {
		p, err := plugin.Open(filepath.Join(directory, f.Name()))
		if err != nil {
			return payloadMap, err
		}
		idfunc, err := p.Lookup("ID")
		if err != nil {
			return payloadMap, err
		}
		namefunc, err := p.Lookup("Name")
		if err != nil {
			return payloadMap, err
		}
		helpfunc, err := p.Lookup("Help")
		if err != nil {
			return payloadMap, err
		}
		genfunc, err := p.Lookup("Generate")
		if err != nil {
			return payloadMap, err
		}
		strfunc, err := p.Lookup("String")
		if err != nil {
			return payloadMap, err
		}
		payloadMap[namefunc.(func() string)()] = PayloadPlugin{
			ID:       idfunc.(func() int),
			Name:     namefunc.(func() string),
			Help:     helpfunc.(func() string),
			Generate: genfunc.(func([]string) ([]byte, error)),
			String:   strfunc.(func([]byte) string),
		}
	}
	return payloadMap, nil
}
