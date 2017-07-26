package main

// Name is a descriptive name of the
func Name() string {
	return "proxy"
}

// ID returns a unique integer across plugins
func ID() int {
	return 3
}

// Help returns documentation for the plugin.
func Help() string {
	return "proxy <url>"
}

// Generate generates a payload based on the arguments provided
func Generate(args []string) ([]byte, error) {

}
