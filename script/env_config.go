package script

import (
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// InitEnvConfig 从环境变量获得的配置内容初始化配置
func InitEnvConfig() (ConfigType, error) {
	var _envConfig = ConfigType{}
	var envConfig = ConfigType{}
	EnvConfigViper := viper.New()
	EnvConfigViper.SetEnvPrefix("compontents_manager") // will be uppercased automatically
	EnvConfigViper.BindEnv("component_name")
	EnvConfigViper.BindEnv("static_path")
	EnvConfigViper.BindEnv("address")
	EnvConfigViper.BindEnv("log_level")
	EnvConfigViper.BindEnv("db_url")
	EnvConfigViper.BindEnv("etcd_url")
	EnvConfigViper.BindEnv("secret")
	EnvConfigViper.BindEnv("use_email")
	EnvConfigViper.BindEnv("email_sender_url")
	EnvConfigViper.BindEnv("default_token_exp")

	err := EnvConfigViper.Unmarshal(&_envConfig)
	for k, v := range _envConfig {
		switch k {
		case "use_email":
			bv := strings.ToLower(v.(string))
			switch bv {
			case "true":
				envConfig[k] = true
			case "false":
				envConfig[k] = false
			}
		case "default_token_exp":
			number, err := strconv.Atoi(v.(string))
			if err == nil {
				envConfig[k] = int(number)
			}
		default:
			envConfig[k] = v.(string)
		}
	}

	return envConfig, err
}
