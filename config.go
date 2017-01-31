package emptynest

import (
	"github.com/BurntSushi/toml"
)

// Config holds configuration data for the server and menu.
type Config struct {
	ServerDebug              bool     `toml:"server_debug"`
	ServerAddr               string   `toml:"server_addr"`
	DBFile                   string   `toml:"db_file"`
	DataDir                  string   `toml:"data_dir"`
	GetLocation              string   `toml:"get_location"`
	GetParam                 string   `toml:"get_param"`
	PostLocation             string   `toml:"post_location"`
	PostParam                string   `toml:"post_param"`
	KillDate                 string   `toml:"kill_date"`
	HostInfoPlugin           string   `toml:"host_info_plugin"`
	PayloadPluginDirectories []string `toml:"payload_plugin_directories"`
	EncoderPluginChain       []string `toml:"encoder_plugin_chain"`
	CryptoPluginChain        []string `toml:"crypto_plugin_chain"`
	KeyChain                 []string `toml:"key_chain"`
}

// DecodeConfigFile returns Config by parsing a toml file.
func DecodeConfigFile(filename string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(filename, &config)
	// TODO: Validate config.
	return config, err
}
