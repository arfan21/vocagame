// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.synapsis.id"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/products": {
            "get": {
                "description": "Get Products",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Product"
                ],
                "summary": "Get Products",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of product",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Owner ID",
                        "name": "owner_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "product_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "allOf": [
                                                {
                                                    "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.PaginationResponse-array_github_com_arfan21_vocagame_internal_model_GetProductResponse"
                                                },
                                                {
                                                    "type": "object",
                                                    "properties": {
                                                        "data": {
                                                            "type": "array",
                                                            "items": {
                                                                "$ref": "#/definitions/github_com_arfan21_vocagame_internal_model.GetProductResponse"
                                                            }
                                                        }
                                                    }
                                                }
                                            ]
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create Product",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Product"
                ],
                "summary": "Create Product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "With the bearer started",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Payload Create Product Request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_internal_model.ProductCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    },
                    "400": {
                        "description": "Error validation field",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.ErrValidationResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/products/:productId": {
            "put": {
                "description": "Update Product",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Product"
                ],
                "summary": "Update Product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "With the bearer started",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Payload Update Product Request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_internal_model.ProductUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    },
                    "400": {
                        "description": "Error validation field",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.ErrValidationResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete Product",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Product"
                ],
                "summary": "Delete Product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "With the bearer started",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "productId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/users/login": {
            "post": {
                "description": "Login user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "Payload user Login Request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_internal_model.UserLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/github_com_arfan21_vocagame_internal_model.UserLoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Error validation field",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.ErrValidationResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/users/logout": {
            "post": {
                "description": "Logout user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Logout user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "With the bearer started",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Payload user Logout Request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_internal_model.UserLogoutRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    },
                    "400": {
                        "description": "Error validation field",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.ErrValidationResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/users/refresh-token": {
            "post": {
                "description": "Refresh Token user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Refresh Token user",
                "parameters": [
                    {
                        "description": "Payload user Refresh Token Request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_internal_model.UserRefreshTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/github_com_arfan21_vocagame_internal_model.UserLoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Error validation field",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.ErrValidationResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/users/register": {
            "post": {
                "description": "Register user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "Payload user Register Request",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_internal_model.UserRegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    },
                    "400": {
                        "description": "Error validation field",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.ErrValidationResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_arfan21_vocagame_internal_model.GetProductResponse": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner_id": {
                    "type": "string"
                },
                "owner_name": {
                    "type": "string"
                },
                "price": {
                    "type": "string"
                },
                "stok": {
                    "type": "integer"
                }
            }
        },
        "github_com_arfan21_vocagame_internal_model.ProductCreateRequest": {
            "type": "object",
            "required": [
                "description",
                "name",
                "price",
                "stok",
                "user_id"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "string"
                },
                "stok": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "github_com_arfan21_vocagame_internal_model.ProductUpdateRequest": {
            "type": "object",
            "required": [
                "description",
                "name",
                "price",
                "stok"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "string"
                },
                "stok": {
                    "type": "integer"
                }
            }
        },
        "github_com_arfan21_vocagame_internal_model.UserLoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 8
                }
            }
        },
        "github_com_arfan21_vocagame_internal_model.UserLoginResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expires_in": {
                    "type": "integer"
                },
                "expires_in_refresh_token": {
                    "type": "integer"
                },
                "refresh_token": {
                    "type": "string"
                },
                "token_type": {
                    "type": "string"
                }
            }
        },
        "github_com_arfan21_vocagame_internal_model.UserLogoutRequest": {
            "type": "object",
            "required": [
                "refresh_token"
            ],
            "properties": {
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "github_com_arfan21_vocagame_internal_model.UserRefreshTokenRequest": {
            "type": "object",
            "required": [
                "refresh_token"
            ],
            "properties": {
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "github_com_arfan21_vocagame_internal_model.UserRegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "fullname",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 8
                }
            }
        },
        "github_com_arfan21_vocagame_pkg_pkgutil.ErrValidationResponse": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_arfan21_vocagame_pkg_pkgutil.HTTPResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "data": {},
                "errors": {
                    "type": "array",
                    "items": {}
                },
                "message": {
                    "type": "string",
                    "example": "Success"
                },
                "status": {
                    "type": "string",
                    "example": "OK"
                }
            }
        },
        "github_com_arfan21_vocagame_pkg_pkgutil.PaginationResponse-array_github_com_arfan21_vocagame_internal_model_GetProductResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_arfan21_vocagame_internal_model.GetProductResponse"
                    }
                },
                "limit": {
                    "type": "integer",
                    "example": 10
                },
                "page": {
                    "type": "integer",
                    "example": 1
                },
                "total_data": {
                    "type": "integer",
                    "example": 1
                },
                "total_page": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Voca Game API",
	Description:      "This is a sample server cell for Voca Game Test API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
