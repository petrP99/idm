package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/role"
	"idm/inner/validator"
	"idm/inner/web"
)

var cfg = common.GetConfig(".env")

func main() {
	// создаём подключение к базе данных
	var dbWithCfg = database.ConnectDbWithCfg(cfg)
	// закрываем соединение с базой данных после выхода из функции main
	defer func() {
		if err := dbWithCfg.Close(); err != nil {
			fmt.Printf("error closing db: %v", err)
		}
	}()
	var server = build(dbWithCfg)
	var err = server.App.Listen(":8080")
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}

func build(db *sqlx.DB) *web.Server {
	server := web.NewServer()
	validate := validator.New()
	employeeRepo := employee.NewEmployeeRepository(db)
	employeeService := employee.NewService(employeeRepo, validate)
	employeeController := employee.NewController(server, employeeService)
	employeeController.RegisterRoutes()

	connectionService := info.NewConnectionService()
	roleRepo := role.NewRoleRepository(db)
	roleService := role.NewService(roleRepo, validate)
	roleController := role.NewController(server, roleService)
	roleController.RegisterRoutes()

	infoController := info.NewController(server, cfg, connectionService)
	infoController.RegisterRoutes()
	return server
}
