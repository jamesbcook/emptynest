package main

import (
	"github.com/BurntSushi/toml"
	"github.com/empty-nest/emptynest"
)

type config struct {
	DBFile                   string   `toml:"db_file"`
	PayloadPluginDirectories []string `toml:"payload_plugin_directories"`
	Debug                    bool
	Transports               []emptynest.TransportConfig
}

// decodeConfigFile returns Config by parsing a toml file.
func decodeConfigFile(filename string) (config, error) {
	var conf config
	_, err := toml.DecodeFile(filename, &conf)
	// TODO: Validate config.
	return conf, err
}
