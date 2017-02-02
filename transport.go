package emptynest

// ApprovalRequest ...
type ApprovalRequest struct {
	Host Host
	Chan chan ApprovalResponse
}

// ApprovalResponse ...
type ApprovalResponse struct {
	Payload Payload
	Plugin  PayloadPlugin
}

// TransportConfig ...
type TransportConfig struct {
	Addr                   string   `toml:"addr"`
	EncoderPluginLocations []string `toml:"encoder_plugin_locations"`
	CryptoPluginLocations  []string `toml:"crypto_plugin_locations"`
	KeyChain               []string `toml:"key_chain"`
}

// TransportCtx is passed to the Create function
// for a TransportPlugin.
type TransportCtx struct {
	Addr                string
	ConfigFileLocation  string
	ApprovalRequestChan chan ApprovalRequest
	DebugChan           chan string
	LogChan             chan string
	Debug               bool
	EncoderChain        []EncoderPlugin
	CryptoChain         []CryptoPlugin
	HostInfo            HostInfoPlugin
	KeyChain            [][]byte
}

// Transport is returned from Create on a plugin.
// It should implement Start() error, which should
// start a listener.
type Transport interface {
	Start() error
	Stop() error
}

// TransportPlugin ...
type TransportPlugin struct {
	Create func(TransportCtx) (Transport, error)
}
