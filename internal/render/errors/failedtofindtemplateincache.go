package errors

import "fmt"

type FailedToFindTemplateInCache struct {
	msg string
}

func CreateFailedToFindTemplateInCache(templateName, templateDir string) FailedToFindTemplateInCache {
	return FailedToFindTemplateInCache{
		msg: fmt.Sprintf("failed to find %s inside %s", templateName, templateDir),
	}
}

func (err FailedToFindTemplateInCache) Error() string {
	return err.msg
}
