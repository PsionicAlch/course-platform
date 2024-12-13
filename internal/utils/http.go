package utils

import (
	"net/http"
	"strconv"
)

func Redirect(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	var statusCode int
	if len(status) > 0 {
		statusCode = status[0]
	} else {
		statusCode = http.StatusFound
	}

	// Check if the request is an HTMX request
	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("HX-Redirect", url)
		return
	}

	// Standard browser redirect for non-HTMX requests
	http.Redirect(w, r, url, statusCode)
}

func IsHTMX(r *http.Request) bool {
	hxRequest, err := strconv.ParseBool(r.Header.Get("HX-Request"))
	return err == nil && hxRequest
}
