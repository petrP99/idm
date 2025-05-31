package role

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRoleRepository(dataBase *sqlx.DB) *Repository {
	return &Repository{db: dataBase}
}

func (repo *Repository) FindById(id int64) (entity Entity, err error) {
	err = repo.db.Get(&entity, "SELECT * FROM role WHERE id=$1", id)
	return entity, err
}

func (repo *Repository) Save(entity Entity) (id int64, err error) {
	err = repo.db.Get(&id, "insert into role (name) values ($1) returning id", entity.Name)
	return id, err
}

func (repo *Repository) FindAll() (listEntity []Entity, err error) {
	err = repo.db.Select(&listEntity, "SELECT * FROM role")
	return listEntity, err
}

func (repo *Repository) FindAllByIds(ids []int64) (listEntity []Entity, err error) {
	if len(ids) == 0 {
		return []Entity{}, nil
	}
	query, args, err := sqlx.In("SELECT * FROM role WHERE id IN (?)", ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build IN query: %w", err)
	}
	query = repo.db.Rebind(query)
	err = repo.db.Select(&listEntity, query, args...)
	return listEntity, err
}

func (repo *Repository) Delete(id int64) error {
	_, err := repo.db.Exec("delete from role where id=$1", id)
	return err
}

func (repo *Repository) DeleteAllByIds(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	query, args, err := sqlx.In("DELETE  FROM role  WHERE ID IN (?)", ids)
	if err != nil {
		return fmt.Errorf("failed to build IN query: %w", err)
	}
	query = repo.db.Rebind(query)
	_, err = repo.db.Exec(query, args...)
	return err
}
