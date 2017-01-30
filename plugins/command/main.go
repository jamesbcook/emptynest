package main

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

// Generate generates a payload based on the arguments provided
func Generate(args []string) ([]byte, error) {

}
