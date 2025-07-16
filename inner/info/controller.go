package info

import (
	"github.com/gofiber/fiber"
	"idm/inner/common"
	"idm/inner/web"
)

type Controller struct {
	server            *web.Server
	cfg               common.Config
	connectionService Srv
}

type Srv interface {
	CheckDbConnection(cfg common.Config) bool
}

func NewController(server *web.Server, cfg common.Config, connectionService Srv) *Controller {
	return &Controller{
		server:            server,
		cfg:               cfg,
		connectionService: connectionService,
	}
}

type InfoResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *Controller) RegisterRoutes() {
	// полный путь будет "/internal/info"
	c.server.GroupInternal.Get("/info", c.GetInfo)
	// полный путь будет "/internal/health"
	c.server.GroupInternal.Get("/health", c.GetHealth)
}

// GetInfo получение информации о приложении
func (c *Controller) GetInfo(ctx *fiber.Ctx) {
	var err = ctx.Status(fiber.StatusOK).JSON(&InfoResponse{
		Name:    c.cfg.AppName,
		Version: c.cfg.AppVersion,
	})
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning info")
		return
	}
}

// GetHealth проверка работоспособности приложения
func (c *Controller) GetHealth(ctx *fiber.Ctx) {
	checkDb := c.connectionService.CheckDbConnection(c.cfg)
	if checkDb {
		ctx.Status(fiber.StatusOK).SendString("Ok")
	} else {
		ctx.Status(fiber.StatusServiceUnavailable).SendString("database is unreachable")
	}
}
