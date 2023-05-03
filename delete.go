package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func DeleteHandler(c *gin.Context) {
	path := c.Param("path")
	log.Printf("DELETE %s", path)

	FileCache.Delete(path)
}
