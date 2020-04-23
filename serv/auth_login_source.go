package serv

import (
	"github.com/gin-gonic/gin"
)

//PermissionListSource 权限资源
type AuthLoginSource struct {
}

//Post 获取全部权限列表
func (source *AuthLoginSource) Post(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "POST"})
}

var authloginsource *AuthLoginSource
