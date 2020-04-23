package serv

import (
	log "github.com/Basic-Components/components_manager/logger"
	"github.com/delicb/gstring"
	"github.com/gin-gonic/gin"
	ginlogrus "github.com/toorop/gin-logrus"
)

type server struct {
	Config map[string]interface{}
	router *gin.Engine
}

func NewServ(config map[string]interface{}) *server {
	var s *server = new(server)
	s.Config = config
	if s.Config["log_level"] != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}
	s.router = gin.New()

	s.router.Use(ginlogrus.Logger(log.Logger), gin.Recovery())

	auth := s.router.Group("/api/auth")
	{
		// 格式校验
		auth.GET("/schema/create", authschemacreatesource.Get)
		auth.POST("/schema/create", authschemacreatesource.Post)
		auth.GET("/schema/update", authschemaupdatesource.Get)
		auth.POST("/schema/update", authschemaupdatesource.Post)
		auth.GET("/schema/login", authschemaloginsource.Get)
		auth.POST("/schema/login", authschemaloginsource.Post)

		// 用户注册
		auth.GET("/user", authuserlistsource.Get)
		auth.POST("/user", authuserlistsource.Post)
		auth.GET("/user/:id", authusersource.Get)
		auth.PUT("/user/:id", authusersource.Put)
		auth.DELETE("/user/:id", authusersource.Delete)
		auth.POST("/login", authloginsource.Post)
		auth.GET("/verify/:token", authverifysource.Get)
	}
	return s
}

func (s *server) Run() {
	log.Info(map[string]interface{}{"config": s.Config}, gstring.Sprintm("servrt start @ {address}", map[string]interface{}{"address": s.Config["address"]}))
	s.router.Run(s.Config["address"].(string))
}
