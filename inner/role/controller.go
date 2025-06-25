package role

import (
	"errors"
	"github.com/gofiber/fiber"
	"idm/inner/common"
	"idm/inner/validator"
	"idm/inner/web"
)

type Controller struct {
	server      *web.Server
	roleService Svc
	validate    *validator.Validator
}

// интерфейс сервиса role.Service
type Svc interface {
	FindById(id int64) (Response, error)
	CreateRole(request CreateRequest) (int64, error)
}

func NewController(server *web.Server, roleService Svc) *Controller {
	return &Controller{
		server:      server,
		roleService: roleService,
		validate:    validator.New(),
	}
}

// функция для регистрации маршрутов
func (c *Controller) RegisterRoutes() {

	// полный маршрут получится "/api/v1/role"
	c.server.GroupApiV1.Post("/role", c.CreateRole)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/role"
func (c *Controller) CreateRole(ctx *fiber.Ctx) {

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

	// вызываем метод CreateRole сервиса role.Service
	var newRoleId, err = c.roleService.CreateRole(request)
	if err != nil {
		switch {
		// если сервис возвращает ошибку RequestValidationError или AlreadyExistsError,
		// то мы возвращаем ответ с кодом 400 (BadRequest)
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		// если сервис возвращает другую ошибку, то мы возвращаем ответ с кодом 500 (InternalServerError)
		default:
			_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
		return
	}

	// функция OkResponse() формирует и направляет ответ в случае успеха
	err = common.OkResponse(ctx, newRoleId)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created role id")
		return
	}
}
