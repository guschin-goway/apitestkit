package usecase

import (
	"github.com/guschin-goway/apitestkit/internal/adapter"
	"github.com/guschin-goway/apitestkit/pkg/domain"
)

type HttpClient interface {
	DoRequest(method, url string, body any, headers map[string]string) (int, []byte, error)
}

type Client struct {
	cfg    *domain.Config
	http   HttpClient
	report adapter.TestReporter
}

func NewClient(cfg *domain.Config, http HttpClient, report adapter.TestReporter) *Client {
	return &Client{cfg: cfg, http: http, report: report}
}

func (c *Client) GET(path string) *Response {
	status, body, err := c.http.DoRequest("GET", c.cfg.BaseURL+path, nil, c.cfg.Headers)
	if err != nil {
		c.report.Fatal(err)
	}
	return &Response{
		status: status,
		body:   body,
		report: c.report,
	}
}
