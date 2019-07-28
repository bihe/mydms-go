// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-07-28 13:54:15.773818 +0200 CEST m=+0.030949201

package docs

import (
	"bytes"
	"encoding/json"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "This is the API of the mydms application",
        "title": "mydms API",
        "contact": {},
        "license": {
            "name": "MIT License",
            "url": "https://raw.githubusercontent.com/bihe/mydms-go/master/LICENSE"
        },
        "version": "2.0"
    },
    "host": "{{.Host}}",
    "basePath": "/api/v1",
    "paths": {
        "/api/v1/appinfo": {
            "get": {
                "description": "meta-data of the application including authenticated user and version",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "appinfo"
                ],
                "summary": "provides information about the application",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/appinfo.AppInfo"
                        }
                    },
                    "401": {},
                    "403": {}
                }
            }
        },
        "/api/v1/tags": {
            "get": {
                "description": "returns all available tags in alphabetical order",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tags"
                ],
                "summary": "retrieve all tags",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/tags.Tag"
                            }
                        }
                    },
                    "401": {},
                    "403": {}
                }
            }
        },
        "/api/v1/tags/search": {
            "get": {
                "description": "returns all tags which match a given search-term",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tags"
                ],
                "summary": "search for tags",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/tags.Tag"
                            }
                        }
                    },
                    "401": {},
                    "403": {}
                }
            }
        }
    },
    "definitions": {
        "appinfo.AppInfo": {
            "type": "object",
            "properties": {
                "userInfo": {
                    "type": "object",
                    "$ref": "#/definitions/appinfo.UserInfo"
                },
                "versionInfo": {
                    "type": "object",
                    "$ref": "#/definitions/appinfo.VersionInfo"
                }
            }
        },
        "appinfo.UserInfo": {
            "type": "object",
            "properties": {
                "displayName": {
                    "description": "DisplayName of authenticated user",
                    "type": "string"
                },
                "email": {
                    "description": "Email of authenticated user",
                    "type": "string"
                },
                "roles": {
                    "description": "Roles the authenticated user possesses",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "userId": {
                    "description": "UserID of authenticated user",
                    "type": "string"
                },
                "userName": {
                    "description": "UserName of authenticated user",
                    "type": "string"
                }
            }
        },
        "appinfo.VersionInfo": {
            "type": "object",
            "properties": {
                "buildDate": {
                    "description": "BuildDate specifies the date of the build",
                    "type": "string"
                },
                "buildNumber": {
                    "description": "BuildNumber defines the specific build",
                    "type": "string"
                },
                "version": {
                    "description": "Version of the application",
                    "type": "string"
                }
            }
        },
        "tags.Tag": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{ Schemes: []string{}}

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface {}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
