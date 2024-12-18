package gocache

import (
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/patrickmn/go-cache"
)

type Generators struct {
	RSSFeed                func() (string, error)
	TutorialsRSSFeed       func() (string, error)
	CoursesRSSFeed         func() (string, error)
	AuthorTutorialsRSSFeed func(authorSlug string) (string, error)
	AuthorCoursesRSSFeed   func(authorSlug string) (string, error)
}

type GoCache struct {
	utils.Loggers
	Cache     *cache.Cache
	Generator *Generators
}

func SetupGoCache(generators *Generators) *GoCache {
	loggers := utils.CreateLoggers("GO CACHE")

	// We don't want the cache to expire. We'll manually delete the cache
	// as and when is necessary.
	c := cache.New(cache.NoExpiration, cache.NoExpiration)

	return &GoCache{
		Loggers:   loggers,
		Cache:     c,
		Generator: generators,
	}
}

func (c *GoCache) InvalidateCache() {
	c.Cache.Flush()
}
