package server

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	zlog "github.com/rs/zerolog/log"

	cmn "github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/static/styles"
)

type FiberServer struct {
	Engine     *fiber.App
	HumaRouter huma.API
	ServerHost string
	ServerPort string
}

func CreateFiber(host, port string) *FiberServer {

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/fiber", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/docs", ServFiberDoc)

	//huma related path
	config := huma.DefaultConfig("", "")
	config.DefaultFormat = "application/json"
	config.DocsPath = "/"
	serv := &FiberServer{
		Engine:     app,
		ServerHost: host,
		ServerPort: port,
	}
	serv.SetupMiddleware()

	humaRouter := humafiber.NewWithGroup(app, app, config)
	serv.HumaRouter = humaRouter

	return serv
}
func (s *FiberServer) SetupMiddleware() {
	// embeddedAssets, err := fs.Sub(EmbedAsset, "public/static/styles")
	// if err != nil {
	// 	panic(err)
	// }
	// s.Engine.Use(gin.Logger())
	// s.Engine.Use(gin.Recovery())
	s.Engine.Static("/assets", "./public/assets")
	s.Engine.Use("/static", filesystem.New(filesystem.Config{
		Root: http.FS(styles.Embedded),
		// Root: http.FS(EmbeddedAssets),
		// PathPrefix: "static",
	}))
}

func (s *FiberServer) Listen() error {

	s.Engine.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	cmn.LogTrace("server started at", fmt.Sprintf("http://127.0.0.1:%s/docs", s.ServerPort))
	err := s.Engine.Listen(s.ServerHost + ":" + s.ServerPort)

	if err != nil {
		zlog.Panic().Err(err).Msg("listen Error")
	}
	return err
}
