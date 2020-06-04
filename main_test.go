package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bgithub.com/Basic-Components/components_manager/script"
	log "github.com/Basic-Components/components_manager/logger"
	"github.com/Basic-Components/components_manager/models"
	"github.com/Basic-Components/components_manager/router"

	"github.com/Basic-Components/connectproxy/etcd3proxy"
	"github.com/Basic-Components/connectproxy/pgproxy"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	ginlogrus "github.com/toorop/gin-logrus"
	"go.etcd.io/etcd/clientv3"
	//"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	err := script.Init()
	if err != nil {
		log.Logger.Fatalln("script init error ", err)
	}
	log.Init(script.Config.LogLevel, map[string]interface{}{
		"component_name": script.Config.ComponentName,
	})
	log.Info(map[string]interface{}{"config": script.Config}, "config inited")
	models.Init()
	log.Info(nil, "models inited")
	err = pgproxy.DB.InitFromURL(script.Config.PGURL)
	if err != nil {
		log.Logger.Fatalln("pgproxy.DB.InitFromURL error ", err)
	}
	log.Info(nil, "pgproxy.DB inited")
	defer pgproxy.DB.Close()

	// 初始化etcd3
	etcdaddresses := strings.Split(script.Config.ETCDURL, ",")
	err = etcd3proxy.Etcd.InitFromOptions(&clientv3.Config{Endpoints: etcdaddresses, DialTimeout: 5 * time.Second})
	if err != nil {
		log.Logger.Fatalln("etcd3proxy.Etcd.InitFromOptions error ", err)
	}
	defer etcd3proxy.Etcd.Close()

	router.Router.Use(ginlogrus.Logger(log.Logger), gin.Recovery())
	router.Init()

	router := router.Router
	// //测试创建
	// w := httptest.NewRecorder()
	// jsonStr := `{
	// 	"name":"test",
	// 	"version":"0.0.1",
	// 	"schema":{
	// 		"$schema": "http://json-schema.org/draft-07/schema#",
	// 		"type": "object",
	//         "required": [
	//             "app",
	//             "bussiness",
	//             "experiment_id"
	//         ],
	//         "properties": {
	//             "app": {
	//                 "type": "string",
	//                 "description": "配置项服务的应用名",
	//                 "enum": [
	//                     "samh",
	//                     "mkz"
	//                 ]
	//             },
	//             "bussiness": {
	//                 "type": "string",
	//                 "description": "业务名"
	//             },
	//             "experiment_id": {
	//                 "type": "integer",
	//                 "description": "实验名"
	//             }
	//         },
	//         "additionalProperties": false
	// 	}
	// }`
	// req, _ := http.NewRequest("POST", "/api/component/", bytes.NewBuffer([]byte(jsonStr)))
	// req.Header.Set("Content-Type", "application/json")
	// router.ServeHTTP(w, req)
	// assert.Equal(t, 200, w.Code)
	// assert.Equal(t, `{"desc":"","id":1,"image":"","name":"test","version":"0.0.1"}`, w.Body.String())

	// 测试获取全部
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/component/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `[{"desc":"","id":1,"image":"","name":"test","version":"0.0.1"}]`, w.Body.String())

	// 测试使用id获取数据
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/component/1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"desc":"","id":1,"image":"","name":"test","version":"0.0.1"}`, w.Body.String())

	// 测试使用id取不存在的组件
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/component/10", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, `{"error":"pg: no rows in result set"}`, w.Body.String())

	// 测试使用非法id取组件信息
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/component/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"strconv.Atoi: parsing \"test\": invalid syntax"}`, w.Body.String())

	// 测试使用id获取schema
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/component/1/schema", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"$id":"","$schema":"http://json-schema.org/draft-07/schema#","additionalProperties":false,"properties":{"app":{"description":"配置项服务的应用名","enum":["samh","mkz"],"type":"string"},"bussiness":{"description":"业务名","type":"string"},"experiment_id":{"description":"实验名","type":"integer"}},"required":["app","bussiness","experiment_id"],"type":"object"}`, w.Body.String())

	// 测试使用id验证schema
	configStr := `{
		"app": "mkz",
		"bussiness":"banner",
		"experiment_id":1234
	}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/component/1/schema", bytes.NewBuffer([]byte(configStr)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"result":true}`, w.Body.String())

	// // 测试创建service
	// queryStr := `{
	// 	"name":"testservice1",
	// 	"version":"0.0.1",
	// 	"component_id": 1,
	// 	"config":{
	// 		"app": "mkz",
	// 		"bussiness":"banner",
	// 		"experiment_id":1234
	// 	}
	// }
	// `
	// w = httptest.NewRecorder()
	// req, _ = http.NewRequest("POST", "/api/service/", bytes.NewBuffer([]byte(queryStr)))
	// req.Header.Set("Content-Type", "application/json")
	// router.ServeHTTP(w, req)
	// assert.Equal(t, 200, w.Code)
	// assert.Equal(t, `{"component":{"desc":"","id":1,"image":"","name":"test","version":"0.0.1"},"desc":"","name":"testservice1","version":"0.0.1"}`, w.Body.String())

	// 测试使用id获取service
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/service/1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"component":{"desc":"","id":1,"image":"","name":"test","version":"0.0.1"},"desc":"","name":"testservice1","version":"0.0.1"}`, w.Body.String())
	// 测试使用id获取component的service
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/component/1/service", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `[{"component":{"desc":"","id":1,"image":"","name":"test","version":"0.0.1"},"desc":"","name":"testservice1","version":"0.0.1"}]`, w.Body.String())

	// 测试使用id获取service的config
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/service/1/config", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"app":"mkz","bussiness":"banner","experiment_id":1234}`, w.Body.String())

	// 测试获取全部service
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/service/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `[{"component":{"desc":"","id":1,"image":"","name":"test","version":"0.0.1"},"desc":"","name":"testservice1","version":"0.0.1"}]`, w.Body.String())

	// 测试获取service 1 的上线状态
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/service/1/release", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"result":true}`, w.Body.String())

	// // 测试更改service
	// updateStr := `{
	// 	"config":{
	// 		"app": "mkz",
	// 		"bussiness":"banner",
	// 		"experiment_id":12345
	// 	}
	// }
	// `
	// w = httptest.NewRecorder()
	// req, _ = http.NewRequest("PUT", "/api/service/1", bytes.NewBuffer([]byte(updateStr)))
	// req.Header.Set("Content-Type", "application/json")
	// router.ServeHTTP(w, req)
	// assert.Equal(t, 200, w.Code)
	// assert.Equal(t, `{"result":"ok"}`, w.Body.String())
}
