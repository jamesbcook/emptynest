package main

import (
	"strings"
)

// Name is a descriptive name of the
func Name() string {
	return "command"
}

// ID returns a unique integer across plugins
func ID() int {
	return 2
}

// Help returns documentation for the plugin.
func Help() string {
	return "command <command>"
}

// String returns a string fiendly version of the command.
func String(data []byte) string {
	return string(data)
}

// Generate generates a payload based on the arguments provided
func Generate(data []byte) ([]byte, error) {
	return data, nil
}

// Process processes arguments and creates the stored payload.
func Process(args []string) ([]byte, error) {
	return []byte(strings.Join(args, " ")), nil
}
