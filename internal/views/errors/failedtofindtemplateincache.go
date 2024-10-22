package errors

import "fmt"

type FailedToFindTemplateInCache struct {
	msg string
}

func CreateFailedToFindTemplateInCache(templateName, cacheName string) FailedToFindTemplateInCache {
	return FailedToFindTemplateInCache{
		msg: fmt.Sprintf("failed to find %s in %s template cache", templateName, cacheName),
	}
}

func (err FailedToFindTemplateInCache) Error() string {
	return err.msg
}
