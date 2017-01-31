package emptynest

import (
	"github.com/BurntSushi/toml"
)

// Config holds configuration data for the server and menu.
type Config struct {
	ServerDebug              bool
	ServerAddr               string
	DBFile                   string
	DataDir                  string
	GetLocation              string
	GetParam                 string
	PostLocation             string
	PostParam                string
	KillDate                 string
	HostInfoPlugin           string
	PayloadPluginDirectories []string
	EncoderPluginChain       []string
	CryptoPluginChain        []string
	KeyChain                 []string
}

// DecodeConfigFile returns Config by parsing a toml file.
func DecodeConfigFile(filename string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(filename, &config)
	// TODO: Validate config.
	return config, err
}
