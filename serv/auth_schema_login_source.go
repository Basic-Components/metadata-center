package serv

import (
	"github.com/gin-gonic/gin"
)

type AuthSchemaLoginSource struct {
}

//Get 获取全部权限列表
func (source *AuthSchemaLoginSource) Get(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "GET"})
}

//Post 获取全部权限列表
func (source *AuthSchemaLoginSource) Post(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "POST"})
}

var authschemaloginsource *AuthSchemaLoginSource
