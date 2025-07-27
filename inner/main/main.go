package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/role"
	"idm/inner/validator"
	"idm/inner/web"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
	go func() {
		var err = server.App.Listen(":8080")
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Создаем группу для ожидания сигнала завершения работы сервера
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	// Запускаем gracefulShutdown в отдельной горутине
	go gracefulShutdown(server, wg)
	// Ожидаем сигнал от горутины gracefulShutdown, что сервер завершил работу
	wg.Wait()
	fmt.Println("Graceful shutdown complete.")

}

// Функция "элегантного" завершения работы сервера по сигналу от операционной системы
func gracefulShutdown(server *web.Server, wg *sync.WaitGroup) {
	// Уведомить основную горутину о завершении работы
	defer wg.Done()
	// Создаём контекст, который слушает сигналы прерывания от операционной системы
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()
	// Слушаем сигнал прерывания от операционной системы
	<-ctx.Done()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")
	// Контекст используется для информирования веб-сервера о том,
	// что у него есть 5 секунд на выполнение запроса, который он обрабатывает в данный момент
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.App.Shutdown(); err != nil {
		fmt.Printf("Server forced to shutdown with error: %v\n", err)
	}
	fmt.Println("Server exiting")
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
