package employee

import (
	"github.com/jmoiron/sqlx"
	employee "idm/inner/role"
	"time"
)

type EmployeeRepository struct {
	db *sqlx.DB
}

type EmployeeEntity struct {
	Id        int64               `db:"id"`
	Name      string              `db:"name"`
	RoleID    employee.RoleEntity `db:"role_id"`
	CreatedAt time.Time           `db:"created_at"`
	UpdatedAt time.Time           `db:"updated_at"`
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

func (repo *EmployeeRepository) FindAll() (listEntity []EmployeeEntity, err error) {
	err = repo.db.Select(&listEntity, "SELECT * FROM employee")
	return listEntity, err
}

func (repo *EmployeeRepository) FindAllByIds(ids []int64) (listEntity []EmployeeEntity, err error) {
	if len(ids) == 0 {
		return []EmployeeEntity{}, nil
	}
	err = repo.db.Select(&listEntity, "select * from employee where id in (?)", ids)
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
	_, err := repo.db.Exec("delete from employee where id in (?)", ids)
	return err
}
