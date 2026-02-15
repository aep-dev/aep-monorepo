package utils

import (
	"fmt"
	"strings"
)

type ListResponse struct {
	Resources     []map[string]interface{} `json:"-"`
	NextPageToken string                   `json:"next_page_token"`
}

type Lister interface {
	List(url string) (*ListResponse, error)
}

func FetchList(lister Lister, baseURL string, pageToken string, maxPageSize int) (*ListResponse, error) {
	var params []string
	if pageToken != "" {
		params = append(params, fmt.Sprintf("page_token=%s", pageToken))
	}
	if maxPageSize > 0 {
		params = append(params, fmt.Sprintf("max_page_size=%d", maxPageSize))
	}

	url := baseURL
	if len(params) > 0 {
		if strings.Contains(url, "?") {
			url += "&"
		} else {
			url += "?"
		}
		url += strings.Join(params, "&")
	}

	return lister.List(url)
}
