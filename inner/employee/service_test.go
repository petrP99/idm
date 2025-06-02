package employee

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert" // импортируем библиотеку с ассерт-функциями
	"github.com/stretchr/testify/mock"   // импортируем пакет для создания моков
	"testing"
	"time"
)

var valueId = int64(1)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindById(id int64) (employee Entity, err error) {
	args := m.Called(id)
	return args.Get(0).(Entity), args.Error(1)
}

func (m *MockRepo) FindAll() ([]Entity, error) {
	args := m.Called()
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) FindAllByIds(ids []int64) (listEntity []Entity, err error) {
	args := m.Called(ids)
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) Save(entity Entity, roleName string) (id int64, err error) {
	args := m.Called(entity, roleName)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) DeleteAllByIds(ids []int64) error {
	args := m.Called(ids)
	return args.Error(0)
}

func TestServices(t *testing.T) {

	var a = assert.New(t)

	t.Run("should return found employee by id", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var entity = Entity{
			Id:        1,
			Name:      "John Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		var want = entity.toResponse()

		repo.On("FindById", valueId).Return(entity, nil)
		var got, err = svc.FindById(1)

		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("should return an error when not found by id", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var entity = Entity{}
		var err = errors.New("user not found")
		var want = fmt.Errorf("error finding employee with id 1: %w", err)

		repo.On("FindById", valueId).Return(entity, err)
		var response, got = svc.FindById(1)

		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("should return all found employees by ids", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var entityes = []Entity{
			{
				Id:        1,
				Name:      "John Doe",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Id:        2,
				Name:      "Doe John",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		var want = toSliceResponse(entityes)

		repo.On("FindAllByIds", []int64{1, 2}).Return(entityes, nil)
		var got, err = svc.FindAllByIds([]int64{1, 2})

		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindAllByIds", 1))
	})

	t.Run("should return all employees", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var entityes = []Entity{
			{
				Id:        1,
				Name:      "John Doe",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Id:        2,
				Name:      "Doe John",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		var want = toSliceResponse(entityes)

		repo.On("FindAll").Return(entityes, nil)
		var got, err = svc.FindAll()

		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindAll", 1))
	})

	t.Run("should delete all employees by ids", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)

		repo.On("DeleteAllByIds", []int64{1, 2}).Return(nil)
		err := svc.DeleteAllByIds([]int64{1, 2})

		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteAllByIds", 1))
	})

	t.Run("should return an error when not found by ids", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var entity = []Entity{{
			Id:        1,
			Name:      "User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}}
		var err = errors.New("users not found")
		var want = fmt.Errorf("error finding employees by ids: %d, %w", []int64{1, 3}, err)

		repo.On("FindAllByIds", []int64{1, 3}).Return(entity, err)
		var response, got = svc.FindAllByIds([]int64{1, 3})

		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindAllByIds", 1))
	})

	t.Run("should delete by id", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)

		repo.On("Delete", valueId).Return(nil)
		err := svc.Delete(1)

		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "Delete", 1))
	})

	t.Run("should return saved employee", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		roleName := "Разработчик"
		var entity = Entity{
			Name:      "User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		repo.On("Save", entity, roleName).Return(valueId, nil)
		var got, err = svc.Save(entity, roleName)

		a.Nil(err)
		a.Equal(valueId, got)
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})

	t.Run("should return error while save employee", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var entity = Entity{
			Name: "",
		}
		var err = errors.New("user not saved")
		var want = fmt.Errorf("error saving employee name: %s: %w", entity.Name, err)
		roleName := "Разработчик"

		repo.On("Save", entity, roleName).Return(int64(0), err)
		var result, got = svc.Save(entity, roleName)

		a.Equal(int64(0), result)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})

}
