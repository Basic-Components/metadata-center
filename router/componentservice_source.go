package router

import (
	"github.com/Basic-Components/components_manager/models"
	"net/http"
	"strconv"

	"github.com/Basic-Components/connectproxy/pgproxy"
	"github.com/gin-gonic/gin"
)

//ComponentServiceSource 组件schema资源
type ComponentServiceSource struct {
}

//Get 获取全部权限列表
func (source *ComponentServiceSource) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		services := []models.Service{}
		err = pgproxy.DB.Cli.Model(&services).Where("service.component_id = ?", id).Column("service.*").Relation("Component").Select()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			res := []map[string]interface{}{}
			for _, service := range services {
				res = append(res, service.InfoWithConfigReleaseStatus())
			}
			ctx.JSON(200, res)
		}
	}
}

var componentservicesource = &ComponentServiceSource{}
