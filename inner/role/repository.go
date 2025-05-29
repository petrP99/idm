package role

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type RoleRepository struct {
	db *sqlx.DB
}

type RoleEntity struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	RoleId    int64     `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewRoleRepository(dataBase *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: dataBase}
}

func (repo *RoleRepository) FindById(id int64) (entity RoleEntity, err error) {
	err = repo.db.Get(&entity, "SELECT * FROM role WHERE id=$1", id)
	return entity, err
}

func (repo *RoleRepository) Save(entity RoleEntity) (id int64, err error) {
	err = repo.db.Get(&id, "insert into role (name) values ($1) returning id", entity.Name)
	return id, err
}

func (repo *RoleRepository) FindAll() (listEntity []RoleEntity, err error) {
	err = repo.db.Select(&listEntity, "SELECT * FROM role")
	return listEntity, err
}

func (repo *RoleRepository) FindAllByIds(ids []int64) (listEntity []RoleEntity, err error) {
	if len(ids) == 0 {
		return []RoleEntity{}, nil
	}
	err = repo.db.Select(&listEntity, "select * from role where id in (?)", ids)
	return listEntity, err
}

func (repo *RoleRepository) Delete(id int64) error {
	_, err := repo.db.Exec("delete from role where id=$1", id)
	return err
}

func (repo *RoleRepository) DeleteAllByIds(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := repo.db.Exec("delete from role where id in (?)", ids)
	return err
}
