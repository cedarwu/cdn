package main

import (
	"log"
	"net/url"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	defaultExpiration    = time.Hour * 24 * 30
	defaultRemoteTimeout = time.Second * 100
	defaultCacheControl  = "public, max-age=2592000"
)

var FileCache *cache.Cache

type CacheFile struct {
	ContentType string
	Content     []byte
}

func init() {
	FileCache = cache.New(defaultExpiration, defaultExpiration)
}

func JoinPath(base, path string) string {
	parsedBase, err := url.Parse(base)
	if err != nil {
		log.Printf("error parsing base: %v", err)
		return ""
	}

	parsedPath, err := url.Parse(path)
	if err != nil {
		log.Printf("error parsing path: %v", err)
		return ""
	}

	joinedURL := parsedBase.ResolveReference(parsedPath)

	return joinedURL.String()
}
