package router

import (
	log "github.com/Basic-Components/components_manager/logger"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	ginlogrus "github.com/toorop/gin-logrus"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Router 默认的路由
var Router = gin.New()

// Init 初始化路由
func Init() {

	Router.Use(ginlogrus.Logger(log.Logger), gin.Recovery())
	// meta := Router.Group("/api/meta")
	// {
	// 	//服务的元信息
	// 	meta.GET("/etcd_url")
	// 	meta.GET("/schema/componet/post")
	// 	meta.GET("/schema/componet_id/put")
	// 	meta.GET("/schema/service/post")
	// 	meta.GET("/schema/service_id/put")
	// }

	componet := Router.Group("/api/component")
	{
		// 组件信息
		componet.GET("/", componentlistsource.Get)
		componet.POST("/", componentlistsource.Post)
		componet.GET("/:id", componentsource.Get)
		componet.PUT("/:id", componentsource.Put)
		componet.DELETE("/:id", componentsource.Delete)
		componet.GET("/:id/schema", componentschemasource.Get)
		componet.POST("/:id/schema", componentschemasource.Post)
		componet.GET("/:id/service", componentservicesource.Get)
	}

	service := Router.Group("/api/service")
	{
		// 服务信息
		service.GET("/", servicelistsource.Get)
		service.POST("/", servicelistsource.Post)
		service.GET("/:id", servicesource.Get)
		service.PUT("/:id", servicesource.Put)
		service.DELETE("/:id", servicesource.Delete)
		service.GET("/:id/config", serviceconfigsource.Get)
		service.GET("/:id/config/release", serviceconfigreleasesource.Get)
		service.POST("/:id/config/release", serviceconfigreleasesource.Post)
		service.DELETE("/:id/config/release", serviceconfigreleasesource.Delete)
	}
}
