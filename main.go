package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/praadit/dating-apps/config"
	"github.com/praadit/dating-apps/controllers"
	"github.com/praadit/dating-apps/middleware"
)

func main() {
	conf, _ := config.LoadEnv()

	bunDB := config.DB(conf)
	memCache := cache.New(1*time.Minute, 10*time.Minute)
	validate := validator.New()

	controller := controllers.NewController(bunDB, memCache, validate)

	corsConf := CorsSetup()

	r := gin.New()
	r.Use(cors.New(corsConf))
	r.Use(gin.Recovery())

	RouteSetup(r, controller)

	if err := r.Run(conf.ServerAddr); err != nil {
		log.Fatal("cannot run the server:", err)
	}
}

func CorsSetup() cors.Config {
	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = []string{"*"}
	corsConf.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConf.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"}
	corsConf.AllowCredentials = true

	return corsConf
}

func RouteSetup(r *gin.Engine, controller *controllers.Controller) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{
			"status": "healthy",
		})
	})
	r.POST("/login", controller.Login)
	r.POST("/signup", controller.Signup)

	r.GET("/explore", middleware.AuthenticatedOnly, controller.Explore)
	r.POST("/swipe", middleware.AuthenticatedOnly, controller.Swipe)

	r.GET("/package", middleware.AuthenticatedOnly, controller.Packages)
	r.GET("/package/:id", middleware.AuthenticatedOnly, controller.Package)
	r.POST("/buy-pack", middleware.AuthenticatedOnly, controller.Buy)
}
