package serv

import (
	"github.com/gin-gonic/gin"
)

//AuthVerifySource 权限资源
type AuthVerifySource struct {
}

//Post 获取全部权限列表
func (source *AuthVerifySource) Get(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "Get"})
}

var authverifysource *AuthVerifySource
