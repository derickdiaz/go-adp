package adp

import (
	"net/http"
	"os"
)

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		return false
	}
	return true
}

func isValidResponseStatusCode(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
