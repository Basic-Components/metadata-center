package router

import (
	"net/http"
	"strconv"

	"github.com/Basic-Components/components_manager/models"

	"github.com/Basic-Components/connectproxy/pgproxy"
	"github.com/gin-gonic/gin"
)

//ComponentSource 权限资源
type ComponentSource struct {
}

//Get 获取组件信息
func (source *ComponentSource) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		component := models.Component{ID: id}
		err = pgproxy.DB.Cli.Select(&component)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(200, component.InfoWithSchema(ctx.Request.Host, ctx.Request.RequestURI))
		}
	}
}

//Put 修改组件信息
func (source *ComponentSource) Put(ctx *gin.Context) {
	q := models.ComponentUpdateOptions{}
	err := ctx.ShouldBindJSON(&q)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			component := models.Component{ID: id}
			err = pgproxy.DB.Cli.Select(&component)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				err := component.Update(&q)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				} else {
					ctx.JSON(200, gin.H{"result": "ok"})
				}
			}
		}
	}
}

//Delete 删除组件
func (source *ComponentSource) Delete(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"id":     ctx.Param("id"),
		"method": "Delete"})
}

var componentsource = &ComponentSource{}
