package tests

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/role"
	"log"
	"os"
)

type Fixture struct {
	DB            *sqlx.DB
	EmployeesRepo *employee.EmployeeRepository
	RoleRepo      *role.RoleRepository
}

func NewFixture() *Fixture {
	db := SetupDB()

	return &Fixture{
		DB:            db,
		EmployeesRepo: employee.NewEmployeeRepository(db),
		RoleRepo:      role.NewRoleRepository(db),
	}
}

func SetupDB() *sqlx.DB {
	db, err := database.ConnectDb()
	if err != nil {
		log.Fatalln("Failed to connect to DB:", err)
	}

	resetDB(db)
	initSql(db)
	return db
}

func resetDB(db *sqlx.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS employee, role CASCADE")
	if err != nil {
		log.Fatalln("Failed to drop tables:", err)
	}
}

func initSql(db *sqlx.DB) {
	sql, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatalln("Cannot read sql file:", err)
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalln("initialization sql failed:", err)
	}
}
