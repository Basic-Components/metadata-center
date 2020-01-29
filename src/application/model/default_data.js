// const defaultService = {
//     name: "test",
//     description: "一个测试",
//     tags: ["test"],
//     status: "test"
// }


const defaultSchemas = [
    {
        task: "card",
        version: "0.0.1",
        status: "dev",
        description: "card",
        schema: {
            "$schema": "http://json-schema.org/draft-07/schema#",
            "description": "A representation of a person, company, organization, or place",
            "type": "object",
            "required": ["familyName", "givenName"],
            "properties": {
                "fn": {
                    "description": "Formatted Name",
                    "type": "string"
                },
                "familyName": {
                    "type": "string"
                },
                "givenName": {
                    "type": "string"
                },
                "additionalName": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "honorificPrefix": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "honorificSuffix": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "nickname": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "email": {
                    "type": "object",
                    "properties": {
                        "type": {
                            "type": "string"
                        },
                        "value": {
                            "type": "string"
                        }
                    }
                },
                "tel": {
                    "type": "object",
                    "properties": {
                        "type": {
                            "type": "string"
                        },
                        "value": {
                            "type": "string"
                        }
                    }
                },
                "adr": { "$ref": "http://example.com/address.schema.json" },
                "geo": { "$ref": "http://example.com/geographical-location.schema.json" },
                "tz": {
                    "type": "string"
                },
                "photo": {
                    "type": "string"
                },
                "logo": {
                    "type": "string"
                },
                "sound": {
                    "type": "string"
                },
                "bday": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "org": {
                    "type": "object",
                    "properties": {
                        "organizationName": {
                            "type": "string"
                        },
                        "organizationUnit": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    {
        task: "address",
        version: "0.0.1",
        status: "dev",
        description: "address",
        schema: {
            "$schema": "http://json-schema.org/draft-07/schema#",
            "description": "An address similar to http://microformats.org/wiki/h-card",
            "type": "object",
            "properties": {
                "post-office-box": {
                    "type": "string"
                },
                "extended-address": {
                    "type": "string"
                },
                "street-address": {
                    "type": "string"
                },
                "locality": {
                    "type": "string"
                },
                "region": {
                    "type": "string"
                },
                "postal-code": {
                    "type": "string"
                },
                "country-name": {
                    "type": "string"
                }
            },
            "required": ["locality", "region", "country-name"],
            "dependencies": {
                "post-office-box": ["street-address"],
                "extended-address": ["street-address"]
            }
        }
    },
    {
        task: "geographical",
        version: "0.0.1",
        status: "dev",
        description: "geographical",
        schema: {
            "$schema": "http://json-schema.org/draft-07/schema#",
            "title": "Longitude and Latitude Values",
            "description": "A geographical coordinate.",
            "required": ["latitude", "longitude"],
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number",
                    "minimum": -90,
                    "maximum": 90
                },
                "longitude": {
                    "type": "number",
                    "minimum": -180,
                    "maximum": 180
                }
            }
        }
    }
]


const defaultService = {
    name: "service_manager",
    description: "管理服务的服务,包含服务数据模型,服务配置分发等",
    tags: ["内部服务", "管理工具"],
    status: "test"
}

const defaultSchemas = [
    {
        task: "query-api-service-list-get",
        version: "0.0.1",
        status: "dev",
        description: "请求service资源的api支持的URL参数",
        tags: ["query", "url_query"],
        schema: {
            "$schema": "http://json-schema.org/draft-07/schema#",
            "type": "object",
            "properties": {
                "tags": {
                    "description": "服务的标签",
                    "oneOf": [
                        {
                            "type": "string"
                        },
                        {
                            "type": "object",
                            "properties": {
                                "$contains": {
                                    "type": "array",
                                    "items": {
                                        "type": "string"
                                    },
                                    "minItems": 1
                                },
                                "$contained": {
                                    "type": "array",
                                    "items": {
                                        "type": "string"
                                    },
                                    "minItems": 1
                                },

                                "$overlap": {
                                    "type": "array",
                                    "items": {
                                        "type": "string"
                                    },
                                    "minItems": 1
                                }
                            },
                            "minProperties": 1,
                            "maxProperties": 1,
                            "additionalProperties": false
                        }
                    ]
                },
                "status": {
                    "description": "服务的状态",
                    "oneOf": [
                        {
                            "type": "string",
                            "enum": ['test', 'pre', 'online', 'deprecated']
                        },
                        {
                            "type": "object",
                            "properties": {
                                "$in": {
                                    "type": "array",
                                    "items": {
                                        "type": "string",
                                        "enum": ['test', 'pre', 'online', 'deprecated']
                                    },
                                    "minItems": 1
                                },
                                "$notIn": {
                                    "type": "array",
                                    "items": {
                                        "type": "string",
                                        "enum": ['test', 'pre', 'online', 'deprecated']
                                    },
                                    "minItems": 1
                                },
                                "$not": {
                                    "type": "string",
                                    "items": {
                                        "type": "string",
                                        "enum": ['test', 'pre', 'online', 'deprecated']
                                    },
                                    "minItems": 1
                                }
                            },
                            "minProperties": 1,
                            "maxProperties": 1,
                            "additionalProperties": false
                        }
                    ]
                }
            }
        }
    }
]

export { defaultService, defaultSchemas }