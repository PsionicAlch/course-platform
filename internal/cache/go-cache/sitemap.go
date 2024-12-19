package gocache

import "github.com/patrickmn/go-cache"

func (c *GoCache) GetSitemap() string {
	const key = "sitemap"

	rssFeed, found := c.Cache.Get(key)
	if found {
		feed, ok := rssFeed.(string)
		if ok {
			return feed
		}
	}

	feed, err := c.Generator.Sitemap()
	if err != nil {
		return feed
	}

	c.Cache.Set(key, feed, cache.NoExpiration)

	return feed
}
