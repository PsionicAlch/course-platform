package sitemapper

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type URL struct {
	Link        string
	Checksum    string
	LastChanged time.Time
}

type Crawler struct {
	Domain   string
	URLMutex sync.Mutex
	Visited  map[string]URL
	Links    map[string]URL
}

func NewCrawler(domain string) *Crawler {
	return &Crawler{
		Domain:  domain,
		Visited: make(map[string]URL),
		Links:   make(map[string]URL),
	}
}

func (crawler *Crawler) Crawl(url string) {
	crawler.URLMutex.Lock()
	defer crawler.URLMutex.Unlock()

	crawler.Visited = make(map[string]URL)

	normalizedURL, ok := crawler.NormalizeURL(url)
	if !ok {
		return
	}

	queue := []string{normalizedURL}

	for len(queue) > 0 {
		currentURL := queue[0]
		queue = queue[1:]

		if _, has := crawler.Visited[currentURL]; has {
			continue
		}

		fmt.Println("Crawling: ", currentURL)

		resp, err := http.Get(currentURL)
		if err != nil {
			fmt.Printf("Error fetching \"%s\": %s\n", currentURL, err)
			continue
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err)
			continue
		}

		links := crawler.ExtractLinks(bytes.NewReader(bodyBytes))

		for _, link := range links {
			if _, has := crawler.Visited[link]; !has {
				queue = append(queue, link)
			}
		}

		hasher := sha256.New()
		hasher.Write(bodyBytes)

		url := URL{
			Link:        currentURL,
			Checksum:    hex.EncodeToString(hasher.Sum(nil)),
			LastChanged: time.Now(),
		}

		crawler.Visited[currentURL] = url
	}

	newLinks := make(map[string]URL)

	for linkVisited, urlVisited := range crawler.Visited {
		if oldUrl, has := crawler.Links[linkVisited]; has {
			if urlVisited.Checksum != oldUrl.Checksum {
				newLinks[linkVisited] = urlVisited
			} else {
				newLinks[linkVisited] = oldUrl
			}
		} else {
			newLinks[linkVisited] = urlVisited
		}
	}

	crawler.Links = newLinks
}

func (crawler *Crawler) GetLinks() []URL {
	crawler.URLMutex.Lock()
	defer crawler.URLMutex.Unlock()

	return slices.Collect(maps.Values(crawler.Links))
}

func (crawler *Crawler) ExtractLinks(r io.Reader) []string {
	links := []string{}
	tokenizer := html.NewTokenizer(r)

	for {
		tt := tokenizer.Next()

		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			// Check anchor tags for href attributes.
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := attr.Val
						if normalized, ok := crawler.NormalizeURL(link); ok {
							links = append(links, normalized)
						}
					}
				}
			} else {
				// Check all other elements for hx-get attributes.
				for _, attr := range token.Attr {
					if attr.Key == "hx-get" {
						link := attr.Val
						if normalized, ok := crawler.NormalizeURL(link); ok {
							links = append(links, normalized)
						}
					}
				}
			}
		case html.ErrorToken:
			return links
		}
	}
}

func (crawler *Crawler) NormalizeURL(href string) (string, bool) {
	parsedURL, err := url.Parse(href)
	if err != nil || parsedURL.Scheme == "javascript" {
		return "", false
	}

	if !parsedURL.IsAbs() {
		baseURL, err := url.Parse(crawler.Domain)
		if err != nil {
			return "", false
		}
		parsedURL = baseURL.ResolveReference(parsedURL)
	}

	parsedURL.Fragment = ""
	normalized := strings.TrimRight(parsedURL.String(), "/")

	if strings.HasPrefix(normalized, crawler.Domain) {
		return normalized, true
	}

	return "", false
}

func EnsureTrailingSlash(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	if parsedURL.Path == "" || !strings.HasSuffix(parsedURL.Path, "/") {
		// Skip if it's a file
		if !strings.Contains(parsedURL.Path, ".") {
			parsedURL.Path += "/"
		}
	}

	return parsedURL.String()
}
