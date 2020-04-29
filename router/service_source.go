package router

import (
	"github.com/Basic-Components/components_manager/models"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Basic-Components/connectproxy/pgproxy"
	"github.com/gin-gonic/gin"
)

//ServiceSource 权限资源
type ServiceSource struct {
}

//Get 获取组件信息
func (source *ServiceSource) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		service := models.Service{}
		err = pgproxy.DB.Cli.Model(&service).Where("service.id = ?", id).Column("service.*").Relation("Component").Select()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(200, service.Info())
		}
	}
}

//Put 修改组件信息
func (source *ServiceSource) Put(ctx *gin.Context) {
	q := models.ServiceUpdateOptions{}
	qq := map[string]interface{}{}
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = json.Unmarshal(data, &qq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		config, ok1 := qq["config"]
		desc, ok2 := qq["desc"]
		if !ok1 && !ok2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "need config or desc"})
		} else {
			if ok1 {
				q.Config = config.(map[string]interface{})
			}
			if !ok2 {
				q.Desc = ""
			} else {
				q.Desc = desc.(string)
			}
			id, err := strconv.Atoi(ctx.Param("id"))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				service := models.Service{}
				err = pgproxy.DB.Cli.Model(&service).Where("service.id = ?", id).Column("service.*").Relation("Component").Select()
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				} else {
					err := service.Update(&q)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					} else {
						ctx.JSON(200, gin.H{"result": "ok"})
					}
				}
			}
		}
	}
}

//Delete 删除组件
func (source *ServiceSource) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		service := models.Service{ID: id}
		_, err := pgproxy.DB.Cli.Model(&service).WherePK().Delete()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(200, gin.H{"result": "ok"})
		}
	}
}

var servicesource = &ServiceSource{}
