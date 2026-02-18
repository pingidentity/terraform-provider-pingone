package utils

import (
	"fmt"
	"net/http"
)

func ResponseErrorDetails(r *http.Response) string {
	if r == nil {
		return "HTTP response is nil"
	}
	return fmt.Sprintf("Response code: %d\nResponse content-type: %s\nFull response body: %+v", r.StatusCode, r.Header.Get("Content-Type"), r.Body)
}
