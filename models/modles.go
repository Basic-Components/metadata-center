package models

import (
	"github.com/Basic-Components/components_manager/connects"
)

// Init 初始化模型
func Init() {
	connects.DB.Regist(userRegistCallback)
}
