package config

import "github.com/bihe/mydms/persistence"

// App contains the "global" components that are
// passed around, especially through HTTP handlers.
type App struct {
	// DB is used to interact with the persistence layer
	DB persistence.Repository
	// V provides version information about the application
	V VersionInfo
}

// VersionInfo provides application meta-data
type VersionInfo struct {
	Build     string
	Version   string
	BuildDate string
}
