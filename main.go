package main

import (
	"strings"

	log "github.com/Basic-Components/components_manager/logger"
	script "github.com/Basic-Components/components_manager/script"
	serv "github.com/Basic-Components/components_manager/serv"

	conn "github.com/Basic-Components/components_manager/connects"
)

func main() {
	config, _ := script.InitConfig()
	address := config["address"].(string)
	appname := config["app_name"].(string)
	log.Init("INFO", map[string]interface{}{"app_name": appname})
	log.Info(map[string]interface{}{"config": config}, "config init")
	etcdAddresses := strings.Split(config["etcd_url"].(string), ",")
	log.Debug(map[string]interface{}{"value": etcdAddresses}, "etcdAddresses")
	err := conn.Etcd.Init(etcdAddresses, 3)
	if err != nil {
		log.Error(map[string]interface{}{"error": err}, "etcd init error")
		return
	} else {
		log.Info(map[string]interface{}{"address": etcdAddresses}, "etcd init ok")
	}
	err = conn.DB.Init(config["db_url"], 10, 20, 5)
	if err != nil {
		log.Error(map[string]interface{}{"error": err}, "db init error")
		return
	} else {
		log.Info(map[string]interface{}{"address": config["db_url"]}, "db init ok")
	}
	serv.RunServer(address)
}
