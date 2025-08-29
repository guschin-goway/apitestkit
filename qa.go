package apitestkit

import (
	"testing"

	"github.com/guschin-goway/apitestkit/internal/adapter"
	"github.com/guschin-goway/apitestkit/internal/usecase"
	"github.com/guschin-goway/apitestkit/pkg/domain"
)

func New(t *testing.T, cfg *domain.Config) *usecase.Client {
	httpClient := adapter.NewHttpClient()
	reporter := adapter.NewTReporter(t)
	return usecase.NewClient(cfg, httpClient, reporter)
}

func NewConfig(baseURL string) *domain.Config {
	return domain.NewConfig(baseURL)
}

// Ассерты наружу
//var (
//	Equal    = domain.Equal
//	NotEmpty = domain.NotEmpty
//	Contains = domain.Contains
//)
