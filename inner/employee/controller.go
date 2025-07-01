package employee

import (
	"errors"
	"github.com/gofiber/fiber"
	"idm/inner/common"
	"idm/inner/validator"
	"idm/inner/web"
)

type Controller struct {
	server          *web.Server
	employeeService Svc
	validate        *validator.Validator
}

// интерфейс сервиса employee.Service
type Svc interface {
	FindById(id int64) (Response, error)
	CreateEmployee(request CreateRequest) (int64, error)
}

func NewController(server *web.Server, employeeService Svc) *Controller {
	return &Controller{
		server:          server,
		employeeService: employeeService,
		validate:        validator.New(),
	}
}

// функция для регистрации маршрутов
func (c *Controller) RegisterRoutes() {

	// полный маршрут получится "/api/v1/employees"
	c.server.GroupApiV1.Post("/employees", c.CreateEmployee)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/employees"
func (c *Controller) CreateEmployee(ctx *fiber.Ctx) {

	// анмаршалим JSON body запроса в структуру CreateRequest
	var request CreateRequest
	if err := ctx.BodyParser(&request); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	if err := c.validate.Validate(request); err != nil {
		var reqErr common.RequestValidationError
		if errors.As(err, &reqErr) {
			_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
			return
		}

		// Если ошибка не RequestValidationError — InternalServerError и т.п.
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "internal validation error")
		return
	}

	// вызываем метод CreateEmployee сервиса employee.Service
	var newEmployeeId, err = c.employeeService.CreateEmployee(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		default:
			_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
		return
	}

	// функция OkResponse() формирует и направляет ответ в случае успеха
	err = common.OkResponse(ctx, newEmployeeId)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
		return
	}
}
