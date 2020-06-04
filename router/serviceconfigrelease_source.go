package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Basic-Components/components_manager/models"

	"github.com/Basic-Components/connectproxy/pgproxy"
	"github.com/gin-gonic/gin"
)

//ServiceConfigReleaseSource 服务的上线资源
type ServiceConfigReleaseSource struct {
}

//Get 获取服务配置是否上线的状态
func (source *ServiceConfigReleaseSource) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		service := models.Service{}
		err = pgproxy.DB.Cli.Model(&service).Where("service.id = ?", id).Column("service.*").Relation("Component").Select()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			res, err := service.ConfigReleased(time.Second)
			if err != nil {
				ctx.JSON(200, gin.H{"error": err.Error(), "result": res})
			} else {
				ctx.JSON(200, gin.H{"result": res})
			}
		}
	}
}

//Post 上线服务配置
func (source *ServiceConfigReleaseSource) Post(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		service := models.Service{}
		err = pgproxy.DB.Cli.Model(&service).Where("service.id = ?", id).Column("service.*").Relation("Component").Select()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			err = service.ReleaseConfig(time.Duration(2) * time.Second)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, gin.H{"result": "ok"})
			}
		}
	}
}

//Delete 下线服务配置
func (source *ServiceConfigReleaseSource) Delete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		service := models.Service{}
		err = pgproxy.DB.Cli.Model(&service).Where("service.id = ?", id).Column("service.*").Relation("Component").Select()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			err = service.UnreleaseConfig(time.Second)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, gin.H{"result": "ok"})
			}
		}
	}
}

var serviceconfigreleasesource = &ServiceConfigReleaseSource{}
