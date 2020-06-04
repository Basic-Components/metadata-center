package router

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Router 默认的路由
var Router = gin.New()

// Init 初始化路由
func Init() {
	//注册静态路由
	Router.Use(static.Serve("/", static.LocalFile("./dist", false)))

	// 注册api路由
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
	sdkInit := Router.Group("/api/sdk_init")
	{
		// 服务信息
		sdkInit.GET("/", sdkinitsource.Get)
	}
}
