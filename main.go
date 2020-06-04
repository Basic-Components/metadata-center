package main //import "github.com/Basic-Components/components_manager"

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	log "github.com/Basic-Components/components_manager/logger"
	"github.com/Basic-Components/components_manager/models"
	"github.com/Basic-Components/components_manager/router"
	"github.com/Basic-Components/components_manager/script"

	"github.com/Basic-Components/connectproxy/etcd3proxy"
	"github.com/Basic-Components/connectproxy/pgproxy"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	ginlogrus "github.com/toorop/gin-logrus"
	"go.etcd.io/etcd/clientv3"
)

func run(addr string, handler *gin.Engine) {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(map[string]interface{}{"error": err}, "listen error")
			os.Exit(2)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info(nil, "Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(map[string]interface{}{"error": err}, "server shutdown error")
		os.Exit(2)
	}
	log.Info(nil, "Server exiting")
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}
func main() {
	// 初始化配置
	err := script.Init()
	if err != nil {
		log.Logger.Fatalln("script init error ", err)
	}
	// 初始化log
	log.Init(script.Config.LogLevel, map[string]interface{}{
		"component_name": script.Config.ComponentName,
		"service_name":   script.Config.ServiceName,
	})
	log.Info(map[string]interface{}{"config": script.Config}, "config inited")
	// 初始化数据模型
	models.Init()
	log.Info(nil, "models inited")
	// 初始化pg数据库
	err = pgproxy.DB.InitFromURL(script.Config.PGURL)
	if err != nil {
		log.Logger.Fatalln("pgproxy.DB.InitFromURL error ", err)
	}
	log.Info(nil, "pgproxy.DB inited")
	defer pgproxy.DB.Close()
	// 初始化etcd3
	etcdaddresses := strings.Split(script.Config.ETCDURL, ",")
	err = etcd3proxy.Etcd.InitFromOptions(&clientv3.Config{Endpoints: etcdaddresses, DialTimeout: 5 * time.Second})
	if err != nil {
		log.Logger.Fatalln("etcd3proxy.Etcd.InitFromOptions error ", err)
	}
	defer etcd3proxy.Etcd.Close()

	// 初始化gin的router
	router.Router.Use(cors.Default())
	router.Router.Use(ginlogrus.Logger(log.Logger), gin.Recovery())
	router.Init()
	// 启动服务
	log.Info(map[string]interface{}{"address": script.Config.Address}, "servrt start")
	run(script.Config.Address, router.Router)
}
