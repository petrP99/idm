package info

import (
	"encoding/json"
	"idm/inner/common"
	"idm/inner/web"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (srv *MockService) CheckDbConnection(cfg common.Config) bool {
	args := srv.Called(cfg)
	return args.Get(0).(bool)
}

func TestInternalApiHealth(t *testing.T) {
	a := assert.New(t)

	t.Run("internal/health - should return Ok", func(t *testing.T) {
		server := web.NewServer()
		cfg := common.GetConfig(".env_info")
		srv := new(MockService)
		controller := NewController(server, cfg, srv)
		controller.RegisterRoutes()

		req := httptest.NewRequest(fiber.MethodGet, "/internal/health", nil)
		srv.On("CheckDbConnection", mock.AnythingOfType("common.Config")).Return(true)

		resp, err := server.App.Test(req)

		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		body := string(bytesData)
		a.NotEmpty(body)
		a.Equal("Ok", body)
	})

	t.Run("internal/health - should return database is unreachable", func(t *testing.T) {
		server := web.NewServer()
		cfg := common.GetConfig(".env_info")
		srv := new(MockService)
		controller := NewController(server, cfg, srv)
		controller.RegisterRoutes()

		req := httptest.NewRequest(fiber.MethodGet, "/internal/health", nil)
		srv.On("CheckDbConnection", mock.AnythingOfType("common.Config")).Return(false)

		resp, err := server.App.Test(req)

		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusServiceUnavailable, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		body := string(bytesData)
		a.NotEmpty(body)
		a.Equal("database is unreachable", body)
	})
}

func TestInternalApiInfo(t *testing.T) {
	a := assert.New(t)

	t.Run("internal/info should return equal true", func(t *testing.T) {
		server := web.NewServer()
		cfg := common.GetConfig(".env_info")
		srv := new(MockService)
		controller := NewController(server, cfg, srv)
		controller.RegisterRoutes()

		req := httptest.NewRequest(fiber.MethodGet, "/internal/info", nil)

		resp, err := server.App.Test(req)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody InfoResponse
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal("idm", responseBody.Name)
		a.Equal("0.0.0", responseBody.Version)
	})
}
