package sitemapper

import (
	"fmt"
	"time"
)

type SiteMapper struct {
	Spider        *Crawler
	RecrawlSignal chan bool
	Domain        string
}

func NewSiteMapper(domain string, crawlInterval time.Duration, startingURL ...string) *SiteMapper {
	var url string
	if len(startingURL) > 0 {
		url = startingURL[0]
	} else {
		url = "/"
	}

	mapper := &SiteMapper{
		Spider:        NewCrawler(domain),
		RecrawlSignal: make(chan bool),
		Domain:        domain,
	}

	go func() {
		time.Sleep(3 * time.Second)

		fmt.Println("Initial site crawl")
		mapper.Spider.Crawl(url)

		ticker := time.NewTicker(crawlInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				mapper.Spider.Crawl(url)
			case <-mapper.RecrawlSignal:
				mapper.Spider.Crawl(url)
			}
		}
	}()

	return mapper
}

func (mapper *SiteMapper) RecrawlSite() {
	mapper.RecrawlSignal <- true
}
