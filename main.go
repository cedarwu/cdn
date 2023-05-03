package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/*path", GetHandler)
	router.DELETE("/*path", DeleteHandler)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := router.Run()
	if err != nil {
		log.Println("run error:", err)
	}
}
