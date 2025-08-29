package examples

import (
	"testing"

	"github.com/guschin-goway/apitestkit"
	"github.com/guschin-goway/apitestkit/pkg/assertion"
)

// структура для распарсенного ответа
type MonitoringHealth struct {
	Name    string `json:"name"`
	Service string `json:"service"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

func TestGetUser(t *testing.T) {
	// создаём конфиг
	cfg := apitestkit.NewConfig("http://localhost:8082").
		WithHeader("Content-Type", "application/json")

	var mh MonitoringHealth

	// сам тест
	apitestkit.New(t, cfg).
		GET("/monitoring/health/").
		Expect().
		Code(200).
		JSONSchema(MonitoringHealth{}).
		JSONObject(&mh).
		Assert(
			assertion.Required(
				assertion.NotEmpty(mh.Status),
				assertion.Equal("UP", mh.Status),
				assertion.NotEmpty(mh.Service),
				assertion.NotEmpty(mh.Version),
				assertion.NotEmpty(mh.Name),
			),
			assertion.Optional(
				assertion.Equal("ASSET_MANAGER", mh.Service),
				assertion.Equal("v1.0.0", mh.Version),
				assertion.Equal("asset-manager", mh.Name),
			),
		)
}
