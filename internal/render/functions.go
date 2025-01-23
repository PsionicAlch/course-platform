package render

import "net/http"

// GetStatusCode returns either the first status code provided or a default status code if none were provided.
func GetStatusCode(status ...int) int {
	var statusCode int

	if len(status) > 0 {
		statusCode = status[0]
	} else {
		statusCode = http.StatusOK
	}

	return statusCode
}
