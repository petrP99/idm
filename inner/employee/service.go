package employee

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo      Repo
	validator Validator
}

type ServiceStub struct {
	repo StubRepo
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{
		repo:      repo,
		validator: validator,
	}
}

type Validator interface {
	Validate(request any) error
}

type StubRepo interface {
	FindAllByIds(ids []int64) ([]Entity, error)
}

type Repo interface {
	FindAll() (listEntity []Entity, err error)
	FindById(id int64) (Entity, error)
	FindAllByIds(ids []int64) (listEntity []Entity, err error)
	Save(entity Entity, roleName string) (id int64, err error)
	Delete(id int64) error
	DeleteAllByIds(ids []int64) error
	BeginTransaction() (tx *sqlx.Tx, err error)
	FindByNameTx(tx *sqlx.Tx, name string) (isExists bool, err error)
	SaveTx(tx *sqlx.Tx, employee Entity) (employeeId int64, err error)
}

func (service *Service) FindById(id int64) (Response, error) {
	var entity, err = service.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee with id %d: %w", id, err)
	}

	return entity.toResponse(), nil
}

func (service *Service) FindAll() ([]Response, error) {
	var entity, err = service.repo.FindAll()
	if err != nil {
		return []Response{}, fmt.Errorf("error finding employees: %w", err)
	}

	return toSliceResponse(entity), nil
}

func (service *Service) FindAllByIds(ids []int64) ([]Response, error) {
	entity, err := service.repo.FindAllByIds(ids)
	if err != nil {
		return []Response{}, fmt.Errorf("error finding employees by ids: %d, %w", ids, err)
	}

	return toSliceResponse(entity), nil
}

func (service *Service) Save(entity Entity, roleName string) (int64, error) {
	var id, err = service.repo.Save(entity, roleName)
	if err != nil {
		return 0, fmt.Errorf("error saving employee name: %s: %w", entity.Name, err)
	}

	return id, nil
}

func (service *Service) Delete(id int64) error {
	err := service.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("error delete employee by id: %d: %w", id, err)
	}

	return nil
}

func (service *Service) DeleteAllByIds(ids []int64) error {
	err := service.repo.DeleteAllByIds(ids)
	if err != nil {
		return fmt.Errorf("error delete employees by ids: %d: %w", ids, err)
	}

	return nil
}

func (service *Service) SaveTx(name string) (int64, error) {
	tx, err := service.repo.BeginTransaction()
	defer func() {
		if tx != nil {
			if err != nil {
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		}
	}()
	if err != nil {
		return 0, fmt.Errorf("error save employee: error creating transaction: %w", err)
	}
	isExist, err := service.repo.FindByNameTx(tx, name)
	if err != nil {
		return 0, fmt.Errorf("error finding employee by name: %s, %w", name, err)
	}
	if isExist {
		return 0, fmt.Errorf("employee with name %s already exists", name)
	}
	entity := Entity{
		Name: name,
	}

	newEmployeeId, err := service.repo.SaveTx(tx, entity)
	if err != nil {
		err = fmt.Errorf("error creating employee with name: %s %v", name, err)
	}
	return newEmployeeId, err
}
