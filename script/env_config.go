package script

import (
	"strings"

	"github.com/spf13/viper"
)

// InitEnvConfig 从环境变量获得的配置内容初始化配置
func InitEnvConfig() (map[string]interface{}, error) {
	var envConfig = map[string]interface{}{}
	EnvConfigViper := viper.New()
	EnvConfigViper.SetEnvPrefix("config_center") // will be uppercased automatically
	EnvConfigViper.BindEnv("SERVICE_NAME")
	EnvConfigViper.BindEnv("ADDRESS")
	EnvConfigViper.BindEnv("LOG_LEVEL")
	EnvConfigViper.BindEnv("PG_URL")
	EnvConfigViper.BindEnv("ETCD_URL")
	err := EnvConfigViper.Unmarshal(&envConfig)

	res := map[string]interface{}{}
	for key, value := range envConfig {
		res[strings.ToUpper(key)] = value
		// switch key {
		// default:
		// }
	}

	return res, err
}
