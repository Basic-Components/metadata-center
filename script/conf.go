package script

import (
	"github.com/xeipuuv/gojsonschema"
)

//ConfigType 配置类型
type ConfigType map[string]interface{}

//DefaultConfig 默认配置
var DefaultConfig = ConfigType{
	"component_name":    "compontents_manager",
	"static_path":       "dist",
	"address":           "0.0.0.0:5000",
	"log_level":         "DEBUG",
	"db_url":            "postgres://postgres:postgres@localhost:5432/test?sslmode=disable",
	"etcd_url":          "172.16.1.105:12379",
	"secret":            "shared-secret",
	"default_token_exp": 1 * 24 * 60 * 60,
	"use_email":         false,
}

//schema 默认的配置样式
const schema = `{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "BUSINESS_CONFIG_CENTER_CONFIG",
    "description": "业务配置中心的配置",
    "type": "object",
    "required": [
        "component_name",
        "static_path",
        "address",
        "log_level",
        "etcd_url",
        "db_url",
        "secret",
        "use_email",
        "default_token_exp"
    ],
    "properties": {
        "component_name": {"type": "string"},
        "static_path": { "type": "string" },
        "address": { "type": "string" },
        "log_level": { "type": "string", "enum": ["DEBUG", "INFO", "WARN", "ERROR"] },
        "etcd_url": { "type": "string" },
        "db_url": { "type": "string" },
        "secret": { "type": "string" },
        "use_email": { "type": "boolean" },
        "email_sender_url": { "type": "string" },
        "default_token_exp": { "type": "integer" }
    },
}`

//VerifyConfig 验证config是否符合要求
func VerifyConfig(conf ConfigType) (bool, *gojsonschema.Result) {
	configLoader := gojsonschema.NewGoLoader(conf)
	schemaLoader := gojsonschema.NewStringLoader(schema)
	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		return false, result
	}
	if result.Valid() {
		return true, result
	}
	return false, result

}
