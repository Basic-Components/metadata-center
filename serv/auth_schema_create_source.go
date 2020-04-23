package serv

import (
	"github.com/gin-gonic/gin"
)

//PermissionListSource 权限资源
type AuthSchemaCreateSource struct {
}

//Get 获取全部权限列表
func (source *AuthSchemaCreateSource) Get(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "GET"})
}

//Post 获取全部权限列表
func (source *AuthSchemaCreateSource) Post(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "POST"})
}

var authschemacreatesource *AuthSchemaCreateSource
