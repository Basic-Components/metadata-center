package script

import (
	"path"
	"strings"

	"github.com/small-tk/pathlib"
	"github.com/spf13/pflag"
)

//InitFlagConfig 从命令行获取配置
func InitFlagConfig() (map[string]interface{}, error) {
	servicename := pflag.StringP("service_name", "s", "", "服务名")
	loglevel := pflag.StringP("loglevel", "l", "", "log的等级")
	address := pflag.StringP("address", "a", "", "要启动的服务器地址")
	pgurl := pflag.StringP("pg_url", "p", "", "连的pg数据库地址")
	etcdurl := pflag.StringP("etcd_url", "e", "", "连的etcd3地址")
	confPath := pflag.StringP("config", "c", "", "配置文件位置")
	pflag.Parse()
	var flagConfig = map[string]interface{}{}

	if *confPath != "" {
		p, err := pathlib.New(*confPath).Absolute()
		if err != nil {
			return flagConfig, err
		}
		if p.Exists() && p.IsFile() {
			filenameWithSuffix := path.Base(*confPath)
			fileSuffix := path.Ext(filenameWithSuffix)
			fileName := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
			dir, err := p.Parent()
			if err != nil {
				return flagConfig, err
			}
			filePaths := []string{dir.Path}
			targetfileconf, err := SetFileConfig(fileName, filePaths)
			if err != nil {
				return flagConfig, err
			}
			for k, v := range targetfileconf {
				flagConfig[k] = v
			}
		}
	}
	if *servicename != "" {
		flagConfig["SERVICE_NAME"] = *servicename
	}
	if *loglevel != "" {
		flagConfig["LOG_LEVEL"] = *loglevel
	}
	if *address != "" {
		flagConfig["ADDRESS"] = *address
	}
	if *pgurl != "" {
		flagConfig["PG_URL"] = *pgurl
	}
	if *etcdurl != "" {
		flagConfig["ETCD_URL"] = *etcdurl
	}

	return flagConfig, nil
}
