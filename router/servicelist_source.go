package router

import (
	"net/http"

	"github.com/Basic-Components/components_manager/models"

	"github.com/gin-gonic/gin"
)

//ServicelistSource 权限资源
type ServicelistSource struct {
}

//Get 获取全部用户
func (source *ServicelistSource) Get(ctx *gin.Context) {
	res, err := models.ServiceGetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(200, res)
	}
}

//Post 创建新用户
func (source *ServicelistSource) Post(ctx *gin.Context) {
	q := models.ServiceNewOptions{}
	err := ctx.ShouldBindJSON(&q)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	service, err := models.ServiceNew(&q)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(200, service.Info())
	}
}

var servicelistsource = &ServicelistSource{}
