package cache

type Cache interface {
	// General functions.
	InvalidateCache()

	// RSS Feed functions.
	GetRSSFeed() string
	GetTutorialsRSSFeed() string
	GetCoursesRSSFeed() string
	GetAuthorTutorialsRSSFeed(authorSlug string) string
	GetAuthorCoursesRSSFeed(authorSlug string) string
}
