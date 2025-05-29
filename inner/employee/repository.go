package employee

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type EmployeeRepository struct {
	db *sqlx.DB
}

type EmployeeEntity struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	RoleID    *int64    `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewEmployeeRepository(dataBase *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{db: dataBase}
}

func (repo *EmployeeRepository) FindById(id int64) (entity EmployeeEntity, err error) {
	err = repo.db.Get(&entity, "SELECT * FROM employee WHERE id=$1", id)
	return entity, err
}

func (repo *EmployeeRepository) Save(entity EmployeeEntity, roleName string) (id int64, err error) {
	query := "insert into employee (name, role_id) values ($1,(select id from role where name = $2)) returning id"
	err = repo.db.Get(&id, query, entity.Name, roleName)
	return id, err
}

func (repo *EmployeeRepository) FindAllByIds(ids []int64) (listEntity []EmployeeEntity, err error) {
	if len(ids) == 0 {
		return []EmployeeEntity{}, nil
	}
	query, args, err := sqlx.In("SELECT * FROM employee WHERE id IN (?)", ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build IN query: %w", err)
	}
	query = repo.db.Rebind(query)
	err = repo.db.Select(&listEntity, query, args...)
	return listEntity, err
}

func (repo *EmployeeRepository) Delete(id int64) error {
	_, err := repo.db.Exec("delete from employee where id=$1", id)
	return err
}

func (repo *EmployeeRepository) DeleteAllByIds(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	query, args, err := sqlx.In("DELETE FROM employee WHERE ID IN (?)", ids)
	if err != nil {
		return fmt.Errorf("failed to build IN query: %w", err)
	}
	query = repo.db.Rebind(query)
	_, err = repo.db.Exec(query, args...)
	return err
}

func (repo *EmployeeRepository) FindAll() (listEntity []EmployeeEntity, err error) {
	err = repo.db.Select(&listEntity, "SELECT * FROM employee")
	return listEntity, err
}
