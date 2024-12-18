package gocache

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/patrickmn/go-cache"
)

type GoCache struct {
	utils.Loggers
	Cache *cache.Cache
}

func SetupGoCache() *GoCache {
	loggers := utils.CreateLoggers("CACHE")

	c := cache.New(time.Hour*24*7, time.Hour*24)

	return &GoCache{
		Loggers: loggers,
		Cache:   c,
	}
}
