package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brian926/Bytez/server/controllers"
	"github.com/brian926/Bytez/server/db"
	"github.com/brian926/Bytez/server/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	load := godotenv.Load()
	if load != nil {
		log.Fatal("Error loading .env file")
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.SetTrustedProxies([]string{"127.0.0.1"})

	url := new(controllers.UrlController)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hey Weclome to Bytez! The URL Shortener API",
		})
	})

	r.POST("/create-short-url", func(c *gin.Context) {
		url.CreateShortUrl(c)
	})

	r.GET("/:shortUrl", func(c *gin.Context) {
		url.HandleShortUrlRedirect(c)
	})

	r.GET("/pong", func(c *gin.Context) {
		controllers.Pong(c)
	})

	storeModel := new(store.UrlModel)
	storeModel.InitializeStore()

	//db.Init()
	db.Init()

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	err := r.Run(":9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}
