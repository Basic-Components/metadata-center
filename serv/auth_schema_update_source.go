package serv

import (
	"github.com/gin-gonic/gin"
)

type AuthSchemaUpdateSource struct {
}

//Get 获取全部权限列表
func (source *AuthSchemaUpdateSource) Get(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "GET"})
}

//Post 获取全部权限列表
func (source *AuthSchemaUpdateSource) Post(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "POST"})
}

var authschemaupdatesource *AuthSchemaUpdateSource
