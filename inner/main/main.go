package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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
var logger = common.NewLogger(cfg)

func main() {
	// создаём подключение к базе данных
	var dbWithCfg = database.ConnectDbWithCfg(cfg)
	defer func() { _ = logger.Sync() }()
	// закрываем соединение с базой данных после выхода из функции main
	defer func() {
		if err := dbWithCfg.Close(); err != nil {
			logger.Error("error closing db: %s", zap.Error(err))
		}
	}()

	var server = build(dbWithCfg)
	go func() {
		var err = server.App.Listen(":8080")
		if err != nil {
			logger.Panic("http server error: %s", zap.Error(err))
		}
	}()

	// Создаем группу для ожидания сигнала завершения работы сервера
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	// Запускаем gracefulShutdown в отдельной горутине
	go gracefulShutdown(server, wg)
	// Ожидаем сигнал от горутины gracefulShutdown, что сервер завершил работу
	wg.Wait()
	logger.Info("graceful shutdown complete")
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
	logger.Info("shutting down gracefully")
	// Контекст используется для информирования веб-сервера о том,
	// что у него есть 5 секунд на выполнение запроса, который он обрабатывает в данный момент
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.App.Shutdown(); err != nil {
		logger.Error("Server forced to shutdown with error", zap.Error(err))
	}
	logger.Info("Server exiting")
}

func build(db *sqlx.DB) *web.Server {
	server := web.NewServer()
	validate := validator.New()
	employeeRepo := employee.NewEmployeeRepository(db)
	employeeService := employee.NewService(employeeRepo, validate)
	employeeController := employee.NewController(server, employeeService, logger)
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
