package serv

import (
	"github.com/gin-gonic/gin"
)

//AuthUserSource 权限资源
type AuthUserlistSource struct {
}

//Get 获取全部用户
func (source *AuthUserlistSource) Get(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "Get"})
}

//Post 创建新用户
func (source *AuthUserlistSource) Post(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "Get"})
}

var authuserlistsource *AuthUserlistSource
