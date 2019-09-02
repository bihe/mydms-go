// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-09-02 19:39:40.8310167 +0200 CEST m=+0.071986201

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
    "basePath": "{{.BasePath}}",
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
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    }
                }
            }
        },
        "/api/v1/documents/search": {
            "get": {
                "description": "use filters to search for docments. the result is a paged set",
                "tags": [
                    "documents"
                ],
                "summary": "search for documents",
                "parameters": [
                    {
                        "type": "string",
                        "description": "title search",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "tag search",
                        "name": "tag",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "sender search",
                        "name": "sender",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "start date",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end date",
                        "name": "to",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "limit max results",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "skip N results",
                        "name": "skip",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/documents.PagedDcoument"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    }
                }
            }
        },
        "/api/v1/documents/{id}": {
            "get": {
                "description": "use the supplied id to lookup the document from the store",
                "tags": [
                    "documents"
                ],
                "summary": "get a document by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "document ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/documents.Document"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    }
                }
            },
            "delete": {
                "description": "use the supplied id to delete the document from the store",
                "tags": [
                    "documents"
                ],
                "summary": "delete a document by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "document ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/documents.Result"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    }
                }
            }
        },
        "/api/v1/file": {
            "get": {
                "description": "use a base64 encoded path to fetch the binary payload of a file from the store",
                "tags": [
                    "filestore"
                ],
                "summary": "get a file from the backend store",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Path",
                        "name": "path",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    }
                }
            }
        },
        "/api/v1/senders": {
            "get": {
                "description": "returns all available senders in alphabetical order",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "senders"
                ],
                "summary": "retrieve all senders",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/senders.Sender"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    }
                }
            }
        },
        "/api/v1/senders/search": {
            "get": {
                "description": "returns all senders which match a given search-term",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "senders"
                ],
                "summary": "search for senders",
                "parameters": [
                    {
                        "type": "string",
                        "description": "SearchString",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/senders.Sender"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    }
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
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    }
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
                "parameters": [
                    {
                        "type": "string",
                        "description": "SearchString",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
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
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    }
                }
            }
        },
        "/api/v1/uploads/file": {
            "post": {
                "description": "temporarily stores a file and creates a item in the repository",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "upload"
                ],
                "summary": "upload a document",
                "parameters": [
                    {
                        "type": "file",
                        "description": "file to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/upload.Result"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/core.ProblemDetail"
                        }
                    }
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
        "core.ProblemDetail": {
            "type": "object",
            "properties": {
                "detail": {
                    "description": "Detail is a human-readable explanation specific to this occurrence of the problem",
                    "type": "string"
                },
                "instance": {
                    "description": "Instance is a URI reference that identifies the specific occurrence of the problem",
                    "type": "string"
                },
                "status": {
                    "description": "Status is the HTTP status code",
                    "type": "integer"
                },
                "title": {
                    "description": "Title is a short, human-readable summary of the problem type",
                    "type": "string"
                },
                "type": {
                    "description": "Type is a URI reference [RFC3986] that identifies the\nproblem type.  This specification encourages that, when\ndereferenced, it provide human-readable documentation for the problem",
                    "type": "string"
                }
            }
        },
        "documents.Document": {
            "type": "object",
            "properties": {
                "alternativeId": {
                    "type": "string"
                },
                "amount": {
                    "type": "number"
                },
                "created": {
                    "type": "string"
                },
                "fileName": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "modified": {
                    "type": "string"
                },
                "previewLink": {
                    "type": "string"
                },
                "senders": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                },
                "uploadFileToken": {
                    "type": "string"
                }
            }
        },
        "documents.PagedDcoument": {
            "type": "object",
            "properties": {
                "documents": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/documents.Document"
                    }
                },
                "totalEntries": {
                    "type": "integer"
                }
            }
        },
        "documents.Result": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "result": {
                    "type": "integer"
                }
            }
        },
        "senders.Sender": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
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
        },
        "upload.Result": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "token": {
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
