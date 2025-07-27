package employee

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"idm/inner/common"
	"idm/inner/web"
	"net/http/httptest"
	"testing"
)

// Объявляем структуру мока сервиса employee.Service
type MockService struct {
	mock.Mock
}

// Реализуем функции мок-сервиса
func (svc *MockService) FindById(id int64) (Response, error) {
	args := svc.Called(id)
	return args.Get(0).(Response), args.Error(1)
}

func (svc *MockService) FindAll() ([]Response, error) {
	args := svc.Called()
	return args.Get(0).([]Response), args.Error(1)
}

func (svc *MockService) CreateEmployee(request CreateRequest) (int64, error) {
	args := svc.Called(request)
	return args.Get(0).(int64), args.Error(1)
}

func getTestRequestBody(req CreateRequest) *bytes.Buffer {
	body, _ := json.Marshal(req)
	return bytes.NewBuffer(body)
}

func TestController(t *testing.T) {
	roleId := int64(2)
	var cfg = common.Config{
		LogDevelopMode: true,
		LogLevel:       "debug",
	}
	mockService := new(MockService)
	logger := common.NewLogger(cfg)
	server := web.NewServer()
	controller := NewController(server, mockService, logger)
	controller.RegisterRoutes()

	t.Run("CreateSuccess", func(t *testing.T) {
		req := CreateRequest{
			Name:   "John Doe",
			RoleId: &roleId,
		}
		body := getTestRequestBody(req)
		request := httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		request.Header.Set("Content-Type", "application/json")

		mockService.On("CreateEmployee", req).Return(int64(123), nil)
		actualResp, err := server.App.Test(request)
		assert.NoError(t, err)

		var response common.Response[int64]
		err = json.NewDecoder(actualResp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, actualResp.StatusCode)
		assert.True(t, response.Success)
		assert.Equal(t, int64(123), response.Data)
		mockService.AssertExpectations(t)
	})

	t.Run("ValidationFailed", func(t *testing.T) {
		req := CreateRequest{
			Name: "John Doe",
		}

		body := getTestRequestBody(req)
		request := httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		request.Header.Set("Content-Type", "application/json")

		mockService.On("CreateEmployee", req).Return(int64(123), nil)
		resp, err := server.App.Test(request)
		assert.NoError(t, err)

		var response common.Response[any]
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.False(t, response.Success)
		assert.Contains(t, response.Message, "RoleId: is required")
	})

	t.Run("AlreadyExistsError", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		req := CreateRequest{
			Name:   "John Doe",
			RoleId: &roleId,
		}

		body := getTestRequestBody(req)
		request := httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		request.Header.Set("Content-Type", "application/json")

		mockService.On("CreateEmployee", req).Return(int64(2), common.AlreadyExistsError{})
		resp, err := server.App.Test(request)

		assert.NoError(t, err)

		var response common.Response[any]
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.False(t, response.Success)
		assert.Contains(t, response.Message, "already exists")

	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		req := CreateRequest{
			Name:   "John Doe",
			RoleId: &roleId,
		}

		body := getTestRequestBody(req)
		request := httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		request.Header.Set("Content-Type", "application/json")

		mockService.On("CreateEmployee", req).Return(int64(1), &common.InternalServerError{})

		resp, err := server.App.Test(request)
		assert.NoError(t, err)

		var response common.Response[any]
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		assert.False(t, response.Success)
		assert.Contains(t, response.Message, "Internal server error")
	})
}
