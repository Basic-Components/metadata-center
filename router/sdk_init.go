package router

import (
	"strings"

	"github.com/Basic-Components/components_manager/script"

	log "github.com/Basic-Components/components_manager/logger"

	"github.com/gin-gonic/gin"
)

//SdkInitSource sdk初始化资源
type SdkInitSource struct {
}

//Get 获取sdk的初始化信息
func (source *SdkInitSource) Get(ctx *gin.Context) {
	ip := strings.Split(ctx.Request.RemoteAddr, ":")[0]
	result := gin.H{"self_ip": ip, "etcd_addresses": script.Config.ETCDURL}
	log.Info(result, "will return")
	ctx.JSON(200, result)
}

var sdkinitsource = &SdkInitSource{}
