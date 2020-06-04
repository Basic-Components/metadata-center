package router

import (
	"net/http"
	"strconv"

	"github.com/Basic-Components/components_manager/models"

	"github.com/Basic-Components/connectproxy/pgproxy"
	"github.com/gin-gonic/gin"
)

//ServiceConfigSource 组件schema资源
type ServiceConfigSource struct {
}

//Get 获取全部权限列表
func (source *ServiceConfigSource) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		service := models.Service{}
		err = pgproxy.DB.Cli.Model(&service).Where("service.id = ?", id).Column("service.*").Relation("Component").Select()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(200, service.Config)
		}
	}
}

var serviceconfigsource = &ServiceConfigSource{}
