package emptynest

// Host is a data type to store an access request.
type Host struct {
	ID        int `storm:"id,increment"`
	IPAddress string
	Hostname  string
	Username  string
	Misc      string
	Status    string
	Key       []byte
}
