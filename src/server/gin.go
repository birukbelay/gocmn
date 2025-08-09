package server

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"

	cmn "github.com/birukbelay/gocmn/src/logger"
)

type GinServer struct {
	Engine     *gin.Engine
	HumaRouter huma.API
	ServerHost string
	ServerPort string
}

func Create(host, port string) *GinServer {
	router := gin.Default()
	//router.Use(CORSMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PATCH", "GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "http://localhost:3000"
		//},
		//MaxAge: 12 * time.Hour,
	}))
	gin.SetMode(gin.DebugMode)
	router.GET("/", func(c *gin.Context) {
		c.String(200, "holla")
	})
	serv := &GinServer{
		Engine:     router,
		ServerHost: host,
		ServerPort: port,
	}
	config := huma.DefaultConfig("", "")
	config.DefaultFormat = "application/json"
	config.DocsPath = ""

	humaRouter := humagin.New(router, config)
	router.GET("/docs", ServGinDoc)
	//huma.AutoRegister(humaRouter, As{})

	serv.SetupMiddleware()
	serv.HumaRouter = humaRouter

	// SetGinRoutes(humaRouter, &dbs)

	return serv
}
func (s *GinServer) SetupMiddleware() {
	s.Engine.Use(gin.Logger())
	s.Engine.Use(gin.Recovery())
	s.Engine.Static("/assets", "./public/assets")
	s.Engine.StaticFS("/static", http.FS(EmbeddedAssets))
}

func (s *GinServer) Listen() error {
	s.Engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	cmn.LogTrace("server started at", fmt.Sprintf("http://127.0.0.1:%s/docs", s.ServerPort))
	err := s.Engine.Run(s.ServerHost + ":" + s.ServerPort)

	if err != nil {
		zlog.Panic().Err(err).Msg("listen Error")
	}
	return err
}
