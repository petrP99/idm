package tests

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"os"
	"testing"
)

var admin = "Administrator"
var developer = "Developer"
var designer = "Designer"
var db = connectDb()

func TestRepository(t *testing.T) {
	var employeeRepository = employee.NewEmployeeRepository(db)

	t.Run("find by id", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				db.MustExec("delete from employee")
			}
		}()
		err := insertTestData("testdata.sql")
		assert.NoError(t, err)
		entity, err := employeeRepository.FindById(1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), entity.Id)

	})

	t.Run("save employee", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				db.MustExec("delete from employee")
			}
		}()
		entity := employee.EmployeeEntity{
			Name: "Test user",
		}

		id, err := employeeRepository.Save(entity, admin)
		assert.NoError(t, err)
		assert.NotNil(t, id)
	})

}

func connectDb() *sqlx.DB {
	cfg := common.Config{
		DbDriverName: "postgres",
		Dsn:          "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
	}
	var db = database.ConnectDbWithCfg(cfg)
	return db
}

func executeSQLFile(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(data))
	return err
}

func insertTestData(fileName string) error {
	// Выполнение SQL из файла testdata.sql
	err := executeSQLFile(fileName)
	return err
}

//
//func () loadTestData() error {
//	// Применяем миграции
//	m, err := migrate.New(
//		"file://migrations",
//		"postgres://test:test@localhost:5432/test_db?sslmode=disable")
//	if err != nil {
//		return err
//	}
//
//	// Накатываем все миграции
//	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
//		return err
//	}
//
//	// Дополнительно заполняем тестовые данные
//	testData := `
//		INSERT INTO role (name) VALUES
//		('Developer'), ('Manager'), ('HR'), ('QA');
//
//		INSERT INTO employee (name, role_id) VALUES
//		('John Doe', 1), ('Alice Smith', 2), ('Bob Johnson', 3), ('Eva Brown', 4);
//	`
//
//	_, err = db.Exec(testData)
//	return err
//}
