package employee

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(dataBase *sqlx.DB) *Repository {
	return &Repository{db: dataBase}
}

func (repo *Repository) FindById(id int64) (entity Entity, err error) {
	err = repo.db.Get(&entity, "SELECT * FROM employee WHERE id=$1", id)
	return entity, err
}

func (repo *Repository) Save(entity Entity, roleName string) (id int64, err error) {
	query := "insert into employee (name, role_id) values ($1,(select id from role where name = $2)) returning id"
	err = repo.db.Get(&id, query, entity.Name, roleName)
	return id, err
}

func (repo *Repository) FindAllByIds(ids []int64) (listEntity []Entity, err error) {
	if len(ids) == 0 {
		return []Entity{}, nil
	}
	query, args, err := sqlx.In("SELECT * FROM employee WHERE id IN (?)", ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build IN query: %w", err)
	}
	query = repo.db.Rebind(query)
	err = repo.db.Select(&listEntity, query, args...)
	return listEntity, err
}

func (repo *Repository) FindAll() (listEntity []Entity, err error) {
	err = repo.db.Select(&listEntity, "SELECT * FROM employee")
	return listEntity, err
}

func (repo *Repository) Delete(id int64) error {
	_, err := repo.db.Exec("delete from employee where id=$1", id)
	return err
}

func (repo *Repository) DeleteAllByIds(ids []int64) error {
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

func (repo *Repository) FindByName(name string) (entity Entity, err error) {
	err = repo.db.Get(&entity, "SELECT * FROM employee WHERE name=$1", name)
	return entity, err
}

func (repo *Repository) BeginTransaction() (tx *sqlx.Tx, err error) {
	return repo.db.Beginx()
}

func (repo *Repository) FindByNameTx(tx *sqlx.Tx, name string) (isExists bool, err error) {
	err = tx.Get(
		&isExists,
		"select exists(select 1 from employee where name = $1)",
		name,
	)
	return isExists, err
}

func (repo *Repository) SaveTx(tx *sqlx.Tx, employee Entity) (employeeId int64, err error) {
	err = tx.Get(
		&employeeId,
		"insert into employee (name) values ($1) returning id",
		employee.Name,
	)
	return employeeId, err
}
