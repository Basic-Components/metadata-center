package router

import (
	"net/http"

	"github.com/Basic-Components/components_manager/models"

	"github.com/gin-gonic/gin"
)

//ComponentlistSource 权限资源
type ComponentlistSource struct {
}

//Get 获取全部用户
func (source *ComponentlistSource) Get(ctx *gin.Context) {
	res, err := models.ComponentGetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(200, res)
	}
}

//Post 创建新用户
func (source *ComponentlistSource) Post(ctx *gin.Context) {
	q := models.ComponentNewOptions{}
	err := ctx.ShouldBindJSON(&q)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	component, err := models.ComponentNew(&q)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(200, component.Info())
	}
}

var componentlistsource = &ComponentlistSource{}
