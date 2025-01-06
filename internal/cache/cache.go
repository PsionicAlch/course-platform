package cache

type Cache interface {
	// General functions.
	InvalidateCache()

	// RSS Feed functions.
	GetGeneralRSSFeed() string
	GetTutorialsRSSFeed() string
	GetTutorialRSSFeed(tutorialSlug string) string
	GetCoursesRSSFeed() string
	GetAuthorTutorialsRSSFeed(authorSlug string) string
	GetAuthorCoursesRSSFeed(authorSlug string) string
}
