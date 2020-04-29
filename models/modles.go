package models

import (
	"fmt"

	"github.com/Basic-Components/connectproxy/pgproxy"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Init 初始化模型
func Init() {
	pgproxy.DB.Regist(componentRegistCallback)
	pgproxy.DB.Regist(serviceRegistCallback)
	fmt.Println("model inited done")
}
