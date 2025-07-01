package role

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

// Объявляем структуру мока сервиса role.Service
type MockService struct {
	mock.Mock
}

// Реализуем функции мок-сервиса
func (svc *MockService) FindById(id int64) (Response, error) {
	args := svc.Called(id)
	return args.Get(0).(Response), args.Error(1)
}

func (svc *MockService) CreateRole(request CreateRequest) (int64, error) {
	args := svc.Called(request)
	return args.Get(0).(int64), args.Error(1)
}

func getTestRequestBody(req CreateRequest) *bytes.Buffer {
	body, _ := json.Marshal(req)
	return bytes.NewBuffer(body)
}

func TestController(t *testing.T) {
	mockService := new(MockService)
	server := web.NewServer()
	controller := NewController(server, mockService)
	controller.RegisterRoutes()

	t.Run("CreateSuccess", func(t *testing.T) {
		req := CreateRequest{
			Name: "John Doe",
		}
		body := getTestRequestBody(req)
		request := httptest.NewRequest(fiber.MethodPost, "/api/v1/role", body)
		request.Header.Set("Content-Type", "application/json")

		mockService.On("CreateRole", req).Return(int64(123), nil)
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
			Name: "",
		}

		body := getTestRequestBody(req)
		request := httptest.NewRequest(fiber.MethodPost, "/api/v1/role", body)
		request.Header.Set("Content-Type", "application/json")

		mockService.On("CreateRole", req).Return(int64(123), nil)
		resp, err := server.App.Test(request)
		assert.NoError(t, err)

		var response common.Response[any]
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.False(t, response.Success)
		assert.Contains(t, response.Message, "Name: is required")
	})

	t.Run("AlreadyExistsError", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		req := CreateRequest{
			Name: "John Doe",
		}

		body := getTestRequestBody(req)
		request := httptest.NewRequest(fiber.MethodPost, "/api/v1/role", body)
		request.Header.Set("Content-Type", "application/json")

		mockService.On("CreateRole", req).Return(int64(2), common.AlreadyExistsError{})
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
			Name: "John Doe",
		}

		body := getTestRequestBody(req)
		request := httptest.NewRequest(fiber.MethodPost, "/api/v1/role", body)
		request.Header.Set("Content-Type", "application/json")

		mockService.On("CreateRole", req).Return(int64(1), &common.InternalServerError{})

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
