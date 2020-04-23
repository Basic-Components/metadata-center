package serv

import (
	"github.com/gin-gonic/gin"
)

//AuthUserSource 权限资源
type AuthUserSource struct {
}

//Get 获取用户信息
func (source *AuthUserSource) Get(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "Get"})
}

//Put 修改用户信息
func (source *AuthUserSource) Put(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "Get"})
}

//delete 删除用户
func (source *AuthUserSource) Delete(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "Get"})
}

var authusersource *AuthUserSource
