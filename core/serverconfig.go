package core

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Configuration holds the application configuration
type Configuration struct {
	Sec   Security     `json:"security"`
	DB    Database     `json:"database"`
	Log   LogConfig    `json:"logging"`
	FS    FileServer   `json:"fileServer"`
	UP    UploadConfig `json:"upload"`
	Store FileStore    `json:"filestore"`
}

// Security settings for the application
type Security struct {
	JwtIssuer     string `json:"jwtIssuer"`
	JwtSecret     string `json:"jwtSecret"`
	CookieName    string `json:"cookieName"`
	LoginRedirect string `json:"loginRedirect"`
	Claim         Claim  `json:"claim"`
	CacheDuration string `json:"cacheDuration"`
}

// Database defines the connection string
type Database struct {
	ConnStr string `json:"connectionString"`
}

// Claim defines the required claims
type Claim struct {
	Name  string   `json:"name"`
	URL   string   `json:"url"`
	Roles []string `json:"roles"`
}

// LogConfig is used to define settings for the logging process
type LogConfig struct {
	FilePath string `json:"filePath"`
	LogLevel string `json:"logLevel"`
}

// FileServer defines the settings for serving static files
type FileServer struct {
	URL          string `json:"url"`
	Path         string `json:"path"`
	SpaIndexFile string `json:"spaIndexFile"`
}

// UploadConfig defines relevant values for the upload logic
type UploadConfig struct {
	// AllowedFileTypes is a list of mime-types allowed to be uploaded
	AllowedFileTypes []string `json:"allowedFileTypes"`
	// MaxUploadSize defines the maximum permissible fiile-size
	MaxUploadSize int64 `json:"maxUploadSize"`
	// UploadPath defines a directory where uploaded files are stored
	UploadPath string `json:"uploadPath"`
}

// FileStore holds configuration settings for the backend file store
type FileStore struct {
	Region string
	Bucket string
	Key    string
	Secret string
}

// GetSettings returns application configuration values
func GetSettings(r io.Reader) (*Configuration, error) {
	var (
		c    Configuration
		cont []byte
		err  error
	)
	if cont, err = ioutil.ReadAll(r); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cont, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
