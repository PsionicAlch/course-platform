package render

import "net/http"

func GetStatusCode(status ...int) int {
	var statusCode int

	if len(status) > 0 {
		statusCode = status[0]
	} else {
		statusCode = http.StatusOK
	}

	return statusCode
}
