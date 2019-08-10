package core

import (
	"strings"
	"testing"
)

const configString = `{
    "security": {
        "jwtIssuer": "login.binggl.net",
        "jwtSecret": "secret",
	"cookieName": "login_token",
	"loginRedirect": "https://login.url.com",
        "claim": {
            "name": "bookmarks",
            "url": "http://localhost:3000",
            "roles": ["User", "Admin"]
	},
	"cacheDuration": "10m"
    },
    "database": {
	"connectionString": "./bookmarks.db"
    },
    "upload": {
        "allowedFileTypes": ["pdf","png"],
        "maxUploadSize": 1000,
        "UploadPath": "/PATH"
    },
    "logging": {
	"logPrefix": "prefix",
	"rollingFileLogger": {
		"filePath": "/temp/file",
		"maxFileSize": 100,
		"numberOfMaxBackups": 4,
		"maxAge": 7,
		"compressFile": false
	}
    },
    "fileServer": {
	    "path": "/tmp",
	    "urlPath": "/ui"
    },
    "filestore": {
        "region": "_REGION_",
        "bucket": "_BUCKET_NAME_",
        "key": "key",
        "secret": "secret"
    }
}`

// TestConfigReader reads config settings from json
func TestConfigReader(t *testing.T) {
	reader := strings.NewReader(configString)
	config, err := GetSettings(reader)
	if err != nil {
		t.Error("Could not read.", err)
	}

	if config.Sec.JwtSecret != "secret" || config.Sec.Claim.Name != "bookmarks" || config.Sec.LoginRedirect != "https://login.url.com" || config.DB.ConnStr != "./bookmarks.db" {
		t.Error("Config values not read!")
	}
	if config.Sec.CacheDuration != "10m" {
		t.Error("Config value Sec.CacheDuration not read!")
	}

	if config.UP.MaxUploadSize != 1000 {
		t.Error("Config value UP.MaxUploadSiz not read!")
	}
	if config.UP.UploadPath != "/PATH" {
		t.Error("Config value UP.UploadPath not read!")
	}
	if len(config.UP.AllowedFileTypes) != 2 {
		t.Error("Config value UP.AllowedFileTypes not read!")
	}

	if config.Log.Prefix != "prefix" {
		t.Error("Could not read logging settings!")
	}
	if config.Log.Rolling.FilePath != "/temp/file" {
		t.Error("Could not read rolling file settings!")
	}
	if config.Log.Rolling.MaxBackups != 4 {
		t.Error("Could not read rolling file settings!")
	}

	if config.FS.Path != "/tmp" {
		t.Error("Could not read FileServer path!")
	}
	if config.FS.URLPath != "/ui" {
		t.Error("Could not read UrlPath!")
	}

	if config.Store.Region != "_REGION_" {
		t.Error("Could not read Store.Region!")
	}
	if config.Store.Bucket != "_BUCKET_NAME_" {
		t.Error("Could not read Store.Bucket!")
	}
	if config.Store.Key != "key" {
		t.Error("Could not read Store.Key!")
	}
	if config.Store.Secret != "secret" {
		t.Error("Could not read Store.Secret!")
	}

}
