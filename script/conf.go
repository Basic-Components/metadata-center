package script

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// ConfigType 配置类型
type ConfigType struct {
	Address          string `json:"ADDRESS"`
	ComponentName    string `json:"COMPONENT_NAME"`
	ComponentVersion string `json:"COMPONENT_VERSION"`
	ServiceName      string `json:"SERVICE_NAME"`
	LogLevel         string `json:"LOG_LEVEL"`
	PGURL            string `json:"PG_URL"`
	ETCDURL          string `json:"ETCD_URL"`
}

//DefaultConfig 默认配置
var DefaultConfig = map[string]interface{}{
	"COMPONENT_NAME":    "config-center",
	"COMPONENT_VERSION": "0.0.0",
	"SERVICE_NAME":      "test-config",
	"ADDRESS":           "0.0.0.0:5000",
	"LOG_LEVEL":         "DEBUG",
	"PG_URL":            "postgres://agens:postgres@localhost:5432/test?sslmode=disable",
	"ETCD_URL":          "localhost:12379",
}

//Schema 默认的配置样式
const Schema = `{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "config-center",
    "description": "配置中心",
    "type": "object",
    "required": [
		"COMPONENT_NAME",
		"COMPONENT_VERSION",
		"SERVICE_NAME",
		"ADDRESS",
		"LOG_LEVEL",
		"PG_URL",
		"ETCD_URL"
    ],
    "properties": {
		"COMPONENT_NAME": { "type": "string" },
		"COMPONENT_VERSION": { "type": "string" },
		"SERVICE_NAME": { "type": "string" },
        "ADDRESS": {"type": "string"},
        "LOG_LEVEL": { "type": "string", "enum": ["DEBUG", "INFO", "WARN", "ERROR"] },
		"PG_URL": { "type": "string" },
		"ETCD_URL": { "type": "string" }
    }
}`

//VerifyConfig 验证config是否符合要求
func VerifyConfig(conf ConfigType) (bool, *gojsonschema.Result) {
	configLoader := gojsonschema.NewGoLoader(conf)
	schemaLoader := gojsonschema.NewStringLoader(Schema)
	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		fmt.Println(err)
		return false, result
	}
	if result.Valid() {
		return true, result
	}
	return false, result
}
