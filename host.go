package emptynest

import "plugin"

// Host is a data type to store an access request.
type Host struct {
	ID        int `storm:"id,increment"`
	IPAddress string
	Status    string
	Info      string
	Data      string
}

// HostInfoPlugin parses and presents Host information
// for requests.
type HostInfoPlugin struct {
	ArgLength    func() int
	SplitPattern func() []byte
	String       func([][]byte) string
}

// BuildHostInfoPlugin returns a HostInfoPlugin from a filename.
func BuildHostInfoPlugin(filename string) (HostInfoPlugin, error) {
	var info HostInfoPlugin
	p, err := plugin.Open(filename)
	if err != nil {
		return info, err
	}
	argfunc, err := p.Lookup("ArgLength")
	if err != nil {
		return info, err
	}
	splitfunc, err := p.Lookup("SplitPattern")
	if err != nil {
		return info, err
	}
	strfunc, err := p.Lookup("String")
	if err != nil {
		return info, err
	}
	info.ArgLength = argfunc.(func() int)
	info.SplitPattern = splitfunc.(func() []byte)
	info.String = strfunc.(func([][]byte) string)
	return info, nil
}
