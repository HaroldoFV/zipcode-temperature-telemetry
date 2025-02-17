// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
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
        "/validate-zipcode": {
            "post": {
                "description": "Valida se o CEP contém 8 dígitos e encaminha para o Serviço B se for válido.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Valida e encaminha um CEP",
                "parameters": [
                    {
                        "description": "CEP",
                        "name": "zipcode",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.ZipcodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Zipcode is valid and forwarded to Service B",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "invalid zipcode",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.ZipcodeRequest": {
            "type": "object",
            "properties": {
                "cep": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Zipcode Validation API",
	Description:      "Esta é uma API para validar e encaminhar CEPs.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
