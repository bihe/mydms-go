package core

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, "./bookmarks.db", config.DB.ConnStr)

	assert.Equal(t, "https://login.url.com", config.Sec.LoginRedirect)
	assert.Equal(t, "bookmarks", config.Sec.Claim.Name)
	assert.Equal(t, "secret", config.Sec.JwtSecret)
	assert.Equal(t, "10m", config.Sec.CacheDuration)

	assert.Equal(t, int64(1000), config.UP.MaxUploadSize)
	assert.Equal(t, "/PATH", config.UP.UploadPath)
	assert.Equal(t, 2, len(config.UP.AllowedFileTypes))

	assert.Equal(t, "prefix", config.Log.Prefix)
	assert.Equal(t, "/temp/file", config.Log.Rolling.FilePath)
	assert.Equal(t, 4, config.Log.Rolling.MaxBackups)

	assert.Equal(t, "/tmp", config.FS.Path)
	assert.Equal(t, "/ui", config.FS.URLPath)

	assert.Equal(t, "_REGION_", config.Store.Region)
	assert.Equal(t, "_BUCKET_NAME_", config.Store.Bucket)
	assert.Equal(t, "_BUCKET_NAME_", config.Store.Bucket)
	assert.Equal(t, "key", config.Store.Key)
	assert.Equal(t, "secret", config.Store.Secret)

}
