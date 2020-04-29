package router

import (
	"github.com/Basic-Components/components_manager/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Basic-Components/connectproxy/pgproxy"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
)

//ComponentSchemaSource 组件schema资源
type ComponentSchemaSource struct {
}

//Get 获取全部权限列表
func (source *ComponentSchemaSource) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		component := models.Component{ID: id}
		err = pgproxy.DB.Cli.Select(&component)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(200, component.SchemaInfo(ctx.Request.Host, ctx.Request.RequestURI))
		}
	}
}

//Post 获取全部权限列表
func (source *ComponentSchemaSource) Post(ctx *gin.Context) {
	q := map[string]interface{}{}
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
				jsonLoader := gojsonschema.NewGoLoader(q)
				schemaLoader := gojsonschema.NewGoLoader(component.Schema)
				result, err := gojsonschema.Validate(schemaLoader, jsonLoader)
				if err != nil {
					ctx.JSON(200, gin.H{"result": false, "error": err.Error()})
				} else {
					if result.Valid() {
						ctx.JSON(200, gin.H{"result": true})
					} else {
						resultErrors := []string{}
						for _, desc := range result.Errors() {
							resultErrors = append(resultErrors, fmt.Sprintf("%s", desc))
						}
						ctx.JSON(200, gin.H{"result": false, "result_errors": resultErrors})
					}
				}
			}
		}
	}
}

var componentschemasource = &ComponentSchemaSource{}
