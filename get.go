package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context) {
	path := c.Param("path")
	if len(path) <= 1 {
		c.JSON(404, gin.H{
			"message": "invalid path",
		})
		return
	}

	log.Printf("GET %s", path)

	data, ok := FileCache.Get(path)
	if ok {
		log.Printf("cache hit %s", path)
		cf := data.(CacheFile)
		c.Writer.Header().Set("Cache-Control", defaultCacheControl)
		c.Data(200, cf.ContentType, cf.Content)
		return
	}

	host := c.Request.Host
	if len(host) == 0 {
		log.Printf("can not get host")
		c.JSON(400, gin.H{
			"message": "invalid host",
		})
		return
	}
	domain, _, err := net.SplitHostPort(host)
	if err != nil {
		domain = host
	}

	base := strings.TrimSpace(c.GetHeader("Origin"))
	if len(base) == 0 {
		base = strings.TrimSpace(c.GetHeader("Referer"))
	}
	if len(base) == 0 {
		c.JSON(404, gin.H{
			"message": "no refer or origin",
		})
		return
	}
	if strings.Contains(base, domain) {
		c.JSON(400, gin.H{
			"message": "invalid refer/origin",
		})
		return
	}

	joinedPath := JoinPath(base, path)
	if len(joinedPath) == 0 {
		c.JSON(400, gin.H{
			"message": "path invalid",
		})
		return
	}

	log.Printf("downloading from %s", joinedPath)

	client := http.Client{
		Timeout: defaultRemoteTimeout,
	}
	resp, err := client.Get(joinedPath)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var cf CacheFile
	cf.ContentType = resp.Header.Get("Content-Type")
	cf.Content, err = io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	if resp.StatusCode != 200 {
		c.Data(resp.StatusCode, cf.ContentType, cf.Content)
		return
	}

	FileCache.SetDefault(path, cf)
	log.Printf("save to cache %s", path)

	c.Writer.Header().Set("Cache-Control", defaultCacheControl)
	c.Data(200, cf.ContentType, cf.Content)
}
