package gocache

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

func (c *GoCache) GetRSSFeed() string {
	const key = "rss-feed"

	rssFeed, found := c.Cache.Get(key)
	if found {
		feed, ok := rssFeed.(string)
		if ok {
			return feed
		}
	}

	feed, err := c.Generator.RSSFeed()
	if err != nil {
		return feed
	}

	c.Cache.Set(key, feed, cache.NoExpiration)

	return feed
}

func (c *GoCache) GetTutorialsRSSFeed() string {
	const key = "tutorials-rss-feed"

	rssFeed, found := c.Cache.Get(key)
	if found {
		feed, ok := rssFeed.(string)
		if ok {
			return feed
		}
	}

	feed, err := c.Generator.TutorialsRSSFeed()
	if err != nil {
		return feed
	}

	c.Cache.Set(key, feed, cache.NoExpiration)

	return feed
}

func (c *GoCache) GetCoursesRSSFeed() string {
	const key = "courses-rss-feed"

	rssFeed, found := c.Cache.Get(key)
	if found {
		feed, ok := rssFeed.(string)
		if ok {
			return feed
		}
	}

	feed, err := c.Generator.CoursesRSSFeed()
	if err != nil {
		return feed
	}

	c.Cache.Set(key, feed, cache.NoExpiration)

	return feed
}

func (c *GoCache) GetAuthorTutorialsRSSFeed(authorSlug string) string {
	var key = fmt.Sprintf("%s-tutorials-rss-feed", authorSlug)

	rssFeed, found := c.Cache.Get(key)
	if found {
		feed, ok := rssFeed.(string)
		if ok {
			return feed
		}
	}

	feed, err := c.Generator.AuthorTutorialsRSSFeed(authorSlug)
	if err != nil {
		return feed
	}

	c.Cache.Set(key, feed, cache.NoExpiration)

	return feed
}

func (c *GoCache) GetAuthorCoursesRSSFeed(authorSlug string) string {
	var key = fmt.Sprintf("%s-courses-rss-feed", authorSlug)

	rssFeed, found := c.Cache.Get(key)
	if found {
		feed, ok := rssFeed.(string)
		if ok {
			return feed
		}
	}

	feed, err := c.Generator.AuthorCoursesRSSFeed(authorSlug)
	if err != nil {
		return feed
	}

	c.Cache.Set(key, feed, cache.NoExpiration)

	return feed
}
