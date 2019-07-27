package config

// App contains the "global" components that are
// passed around, especially through HTTP handlers.
type App struct {
	//DB        *sqlx.DB
	// V provides version information about the application
	V VersionInfo
}

// VersionInfo provides application meta-data
type VersionInfo struct {
	Build     string
	Version   string
	BuildDate string
}
