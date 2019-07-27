package main

// NewServerArgs creates a new instance
func NewServerArgs() *ServerArgs {
	c := new(ServerArgs)
	return c
}

// ServerArgs is uded to configure the API server
type ServerArgs struct {
	HostName   string
	Port       int
	ConfigFile string
}
