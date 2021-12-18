// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/block/{block}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get block by hash or id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blocks"
                ],
                "summary": "GetBlock",
                "operationId": "get-block",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Block id or hash",
                        "name": "block",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Block"
                        }
                    },
                    "404": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/blocks": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get blocks with limit and offset",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blocks"
                ],
                "summary": "GetBlocks",
                "operationId": "get-blocks",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "the amount of blocks you want to get",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "offset from the beginning of the data in database",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Block"
                            }
                        }
                    },
                    "404": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/transactions": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get transactions with limit and offset params",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transactions"
                ],
                "summary": "GetTransactions",
                "operationId": "get-transactions",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "the amount of transactions you want to get",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "offset from the beginning of the data in database",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Transaction"
                            }
                        }
                    },
                    "404": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Arg": {
            "type": "object",
            "properties": {
                "args": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Value"
                    }
                },
                "bytes": {
                    "type": "string"
                },
                "int": {
                    "type": "string"
                },
                "prim": {
                    "type": "string"
                }
            }
        },
        "models.BigMapDiff": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "big_map": {
                    "type": "string"
                },
                "key": {
                    "$ref": "#/definitions/models.Value"
                },
                "key_hash": {
                    "type": "string"
                }
            }
        },
        "models.Block": {
            "type": "object",
            "properties": {
                "chain_id": {
                    "type": "string"
                },
                "hash": {
                    "type": "string"
                },
                "header": {
                    "$ref": "#/definitions/models.Header"
                },
                "metadata": {
                    "$ref": "#/definitions/models.BlockMetadata"
                },
                "operations": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "$ref": "#/definitions/models.Operation"
                        }
                    }
                },
                "protocol": {
                    "type": "string"
                }
            }
        },
        "models.BlockMetadata": {
            "type": "object",
            "properties": {
                "baker": {
                    "type": "string"
                },
                "consumed_gas": {
                    "type": "string"
                },
                "level": {
                    "$ref": "#/definitions/models.LevelInfo"
                },
                "level_info": {
                    "$ref": "#/definitions/models.LevelInfo"
                }
            }
        },
        "models.Content": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                },
                "counter": {
                    "type": "string"
                },
                "destination": {
                    "type": "string"
                },
                "fee": {
                    "type": "string"
                },
                "gas_limit": {
                    "type": "string"
                },
                "kind": {
                    "type": "string"
                },
                "metadata": {
                    "$ref": "#/definitions/models.Metadata"
                },
                "parameters": {
                    "$ref": "#/definitions/models.Parameters"
                },
                "source": {
                    "type": "string"
                },
                "storage_limit": {
                    "type": "string"
                }
            }
        },
        "models.Diff": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "updates": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Update"
                    }
                }
            }
        },
        "models.Header": {
            "type": "object",
            "properties": {
                "baker_fee": {
                    "type": "integer"
                },
                "fitness": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "operations_hash": {
                    "type": "string"
                },
                "predecessor": {
                    "type": "string"
                },
                "priority": {
                    "type": "integer"
                },
                "proto": {
                    "type": "integer"
                },
                "signature": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                },
                "validation_pass": {
                    "type": "integer"
                }
            }
        },
        "models.InternalOperationResult": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                },
                "destination": {
                    "type": "string"
                },
                "kind": {
                    "type": "string"
                },
                "nonce": {
                    "type": "integer"
                },
                "result": {
                    "$ref": "#/definitions/models.Result"
                },
                "source": {
                    "type": "string"
                }
            }
        },
        "models.LazyStorageDiff": {
            "type": "object",
            "properties": {
                "diff": {
                    "$ref": "#/definitions/models.Diff"
                },
                "id": {
                    "type": "string"
                },
                "kind": {
                    "type": "string"
                }
            }
        },
        "models.LevelInfo": {
            "type": "object",
            "properties": {
                "cycle": {
                    "type": "integer"
                },
                "cycle_position": {
                    "type": "integer"
                },
                "level": {
                    "type": "integer"
                }
            }
        },
        "models.Metadata": {
            "type": "object",
            "properties": {
                "balance_updates": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.MetadataBalanceUpdate"
                    }
                },
                "internal_operation_results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.InternalOperationResult"
                    }
                },
                "operation_result": {
                    "$ref": "#/definitions/models.OperationResult"
                }
            }
        },
        "models.MetadataBalanceUpdate": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "change": {
                    "type": "string"
                },
                "contract": {
                    "type": "string"
                },
                "cycle": {
                    "type": "integer"
                },
                "delegate": {
                    "type": "string"
                },
                "kind": {
                    "type": "string"
                },
                "origin": {
                    "type": "string"
                }
            }
        },
        "models.Operation": {
            "type": "object",
            "properties": {
                "branch": {
                    "type": "string"
                },
                "chain_id": {
                    "type": "string"
                },
                "contents": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Content"
                    }
                },
                "hash": {
                    "type": "string"
                },
                "protocol": {
                    "type": "string"
                },
                "signature": {
                    "type": "string"
                }
            }
        },
        "models.OperationResult": {
            "type": "object",
            "properties": {
                "big_map_diff": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.BigMapDiff"
                    }
                },
                "consumed_gas": {
                    "type": "string"
                },
                "consumed_milligas": {
                    "type": "string"
                },
                "lazy_storage_diff": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.LazyStorageDiff"
                    }
                },
                "status": {
                    "type": "string"
                },
                "storage": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Storage"
                    }
                },
                "storage_size": {
                    "type": "string"
                }
            }
        },
        "models.Parameters": {
            "type": "object",
            "properties": {
                "entrypoint": {
                    "type": "string"
                },
                "value": {
                    "$ref": "#/definitions/models.Value"
                }
            }
        },
        "models.Result": {
            "type": "object",
            "properties": {
                "balance_updates": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ResultBalanceUpdate"
                    }
                },
                "consumed_gas": {
                    "type": "string"
                },
                "consumed_milligas": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.ResultBalanceUpdate": {
            "type": "object",
            "properties": {
                "change": {
                    "type": "string"
                },
                "contract": {
                    "type": "string"
                },
                "kind": {
                    "type": "string"
                },
                "origin": {
                    "type": "string"
                }
            }
        },
        "models.Storage": {
            "type": "object",
            "properties": {
                "args": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Arg"
                    }
                },
                "int": {
                    "type": "string"
                },
                "prim": {
                    "type": "string"
                }
            }
        },
        "models.Transaction": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                },
                "block_hash": {
                    "type": "string"
                },
                "branch": {
                    "type": "string"
                },
                "consumed_milligas": {
                    "type": "string"
                },
                "counter": {
                    "type": "string"
                },
                "destination": {
                    "type": "string"
                },
                "fee": {
                    "type": "string"
                },
                "gas_limit": {
                    "type": "string"
                },
                "hash": {
                    "type": "string"
                },
                "signature": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "storage_size": {
                    "type": "string"
                }
            }
        },
        "models.Update": {
            "type": "object",
            "properties": {
                "key": {
                    "$ref": "#/definitions/models.Value"
                },
                "key_hash": {
                    "type": "string"
                }
            }
        },
        "models.Value": {
            "type": "object",
            "properties": {
                "int": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
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
var SwaggerInfo = swaggerInfo{
	Version:     "0.9.1",
	Host:        "localhost",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Blockchain Explorer",
	Description: "This is a service that allows you to receive data stored in the blockchain.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
