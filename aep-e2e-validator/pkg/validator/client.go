package validator

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type Header struct {
	Key   string
	Value string
}

type RequestLog struct {
	Method   string `json:"method"`
	URL      string `json:"url"`
	ReqBody  string `json:"request_body,omitempty"`
	ReqType  string `json:"request_content_type,omitempty"`
	RespCode int    `json:"response_code,omitempty"`
	RespBody string `json:"response_body,omitempty"`
	RespType string `json:"response_content_type,omitempty"`
}

type extendedClient struct {
	inner   *http.Client
	headers []Header
	logs    []RequestLog
	logger  *log.Logger
}

func (c *extendedClient) clearLogs() {
	c.logs = nil
}

func prettyPrintBody(body string, contentType string) string {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return ""
	}
	if strings.Contains(strings.ToLower(contentType), "json") {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, []byte(trimmed), "     ", "  "); err == nil {
			return prettyJSON.String()
		}
	}
	return trimmed
}

func (c *extendedClient) printLogs() {
	if len(c.logs) == 0 {
		return
	}
	c.logger.Println("   --- Request/Response Logs ---")
	for i, l := range c.logs {
		c.logger.Printf("   Request %d:\n", i+1)
		c.logger.Printf("     %s %s\n", l.Method, l.URL)
		if l.ReqBody != "" {
			c.logger.Printf("     Body:\n     %s\n", prettyPrintBody(l.ReqBody, l.ReqType))
		}
		if l.RespCode != 0 {
			c.logger.Printf("   Response %d:\n", i+1)
			c.logger.Printf("     Status: %d\n", l.RespCode)
			if l.RespBody != "" {
				c.logger.Printf("     Body:\n     %s\n", prettyPrintBody(l.RespBody, l.RespType))
			}
		} else {
			c.logger.Printf("   Response %d: (No response)\n", i+1)
		}
	}
	c.logger.Println("   -----------------------------")
}

func (c *extendedClient) Do(req *http.Request) (*http.Response, error) {
	for _, h := range c.headers {
		req.Header.Add(h.Key, h.Value)
	}

	var reqBody string
	reqType := req.Header.Get("Content-Type")
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		reqBody = string(bodyBytes)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	resp, err := c.inner.Do(req)

	var respBody string
	var respCode int
	var respType string
	if resp != nil {
		respCode = resp.StatusCode
		respType = resp.Header.Get("Content-Type")
		if resp.Body != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			respBody = string(bodyBytes)
			resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	c.logs = append(c.logs, RequestLog{
		Method:   req.Method,
		URL:      req.URL.String(),
		ReqBody:  reqBody,
		ReqType:  reqType,
		RespCode: respCode,
		RespBody: respBody,
		RespType: respType,
	})

	return resp, err
}
