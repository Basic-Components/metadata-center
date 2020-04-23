package script

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// Config 程序的配置
var Config = DefaultConfig

// InitConfig 初始化命令行传入的参数到配置,返回值为false表示要执行创建秘钥否则为启动服务
func InitConfig() error {
	log.SetFormatter(&log.JSONFormatter{})
	defaultfileConfig, err := InitFileConfig()
	if err != nil {
		log.Warn("从默认路径文件初始化配置项错误")
	} else {
		for k, v := range defaultfileConfig {
			Config[k] = v
		}
	}
	envConfig, err := InitEnvConfig()
	if err != nil {
		log.WithFields(map[string]interface{}{
			"error": err,
		}).Warn("从环境变量初始化配置项错误")
	} else {
		for k, v := range envConfig {
			Config[k] = v
		}
	}
	flagConfig, err := InitFlagConfig()
	if err != nil {
		log.WithFields(map[string]interface{}{
			"error": err,
		}).Warn("从命令行参数初始化配置项错误")
	} else {
		for k, v := range flagConfig {
			Config[k] = v
		}
	}
	flag, result := VerifyConfig(Config)
	if flag == true {
		return nil
	}
	for _, err := range result.Errors() {
		log.WithFields(map[string]interface{}{
			"error": err,
		}).Error("配置检验错误")
	}
	return errors.New("配置文件参数错误")
}
