package errors

import "fmt"

type FailedToCompileTemplate struct {
	msg string
}

func CreateFailedToCompileTemplate(templateName, err string) FailedToCompileTemplate {
	return FailedToCompileTemplate{
		msg: fmt.Sprintf("failed to compile %s: %s", templateName, err),
	}
}

func (err FailedToCompileTemplate) Error() string {
	return err.msg
}
