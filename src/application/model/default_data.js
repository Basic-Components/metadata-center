const selfService = {
    name: "service_manager",
    tag: "test-0.0.1",
    description: "管理服务的服务,包含服务数据模型,服务配置分发等",
    labels: ["内部服务", "管理工具"]
}

const selfSchemas = [
    {
        task: "operation-array-string",
        version: "0.0.1",
        env: ['dev', 'test', 'produce', 'release'],
        description: "请求资源的api查找符合条件的array类型数据的可选操作",
        direction: "ref",
        labels: ["api", "operation"],
        schema: {
            "$schema": "http://json-schema.org/draft-07/schema#",
            "title": "operation-array-string",
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
    },
    {
        task: "operation-enum-env",
        version: "0.0.1",
        env: ['dev', 'test', 'produce', 'release'],
        description: "请求资源的api查找符合条件的array类型数据的可选操作",
        direction: "ref",
        labels: ["api", "operation"],
        schema: {
            "$schema": "http://json-schema.org/draft-07/schema#",
            "title": "operation-enum-env",
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
    },
    {
        task: "response-self",
        version: "0.0.1",
        env: ['dev', 'test', 'produce', 'release'],
        description: "请求service资源的api支持的URL参数",
        direction: "ref",
        labels: ["api"],
        schema: {
            "$schema": "http://json-schema.org/draft-07/schema#",
            "title": "response-self",
            "type": "object",
            "required": ["src", "url", "description"],
            "properties": {
                "src": {
                    "type": "string",
                    "description": "完整的请求路径"
                },
                "url": {
                    "type": "string",
                    "description": "去掉schema和host信息的请求路径"
                },
                "description": {
                    "type": "string",
                    "description": "描述文本"
                },
                "query": {
                    "type": "object",
                    "description": "请求对象"
                }
            }
        }
    },
    {
        task: "service-list-get",
        version: "0.0.1",
        env: ['dev', 'test', 'produce', 'release'],
        description: "请求service资源的api获取的数据模式",
        direction: "response",
        labels: ["api"],
        schema: {
            "$schema": "http://json-schema.org/draft-07/schema#",
            "title":"service-list-get-response",
            "type": "object",
            "properties": {
                "self":{
                    "$ref":"http://localhost:5000/api/service/service_manager/test-0.0.1/?task=response-self&version=0.0.1&direction=ref"
                },
                ""
            }
        }
    },
]

export { selfService, selfSchemas }