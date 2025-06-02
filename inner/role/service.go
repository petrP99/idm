package role

import "fmt"

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

type Repo interface {
	FindAll() (listEntity []Entity, err error)
	FindById(id int64) (Entity, error)
	FindAllByIds(ids []int64) (listEntity []Entity, err error)
	Save(entity Entity) (id int64, err error)
	Delete(id int64) error
	DeleteAllByIds(ids []int64) error
}

func (service *Service) FindById(id int64) (Response, error) {
	var entity, err = service.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding role with id %d: %w", id, err)
	}

	return entity.toResponse(), nil
}

func (service *Service) FindAll() ([]Response, error) {
	var entity, err = service.repo.FindAll()
	if err != nil {
		return []Response{}, fmt.Errorf("error finding roles: %w", err)
	}

	return toSliceResponse(entity), nil
}

func (service *Service) FindAllByIds(ids []int64) ([]Response, error) {
	entity, err := service.repo.FindAllByIds(ids)
	if err != nil {
		return []Response{}, fmt.Errorf("error finding roles by ids: %d, %w", ids, err)
	}

	return toSliceResponse(entity), nil
}

func (service *Service) Save(entity Entity) (int64, error) {
	var id, err = service.repo.Save(entity)
	if err != nil {
		return 0, fmt.Errorf("error saving role name: %s: %w", entity.Name, err)
	}

	return id, nil
}

func (service *Service) Delete(id int64) error {
	err := service.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("error delete role by id: %d: %w", id, err)
	}

	return nil
}

func (service *Service) DeleteAllByIds(ids []int64) error {
	err := service.repo.DeleteAllByIds(ids)
	if err != nil {
		return fmt.Errorf("error delete roles by ids: %d: %w", ids, err)
	}

	return nil
}
