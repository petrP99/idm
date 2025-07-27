package role

import (
	"errors"
	"fmt"
	"github.com/78bits/go-sqlmock-sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"idm/inner/validator"
	"testing"
	"time"
)

var valueId = int64(1)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindById(id int64) (role Entity, err error) {
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

func (m *MockRepo) Save(entity Entity) (id int64, err error) {
	args := m.Called(entity)
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

func (m *MockRepo) BeginTransaction() (*sqlx.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (m *MockRepo) FindByNameTx(tx *sqlx.Tx, name string) (bool, error) {
	args := m.Called(tx, name)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockRepo) SaveTx(tx *sqlx.Tx, entity Entity) (int64, error) {
	args := m.Called(tx, entity)
	return args.Get(0).(int64), args.Error(1)
}

var val = validator.New()

func TestServiceSaveTxSuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close database: %v", err)
		}
	}()

	sqlxDB := sqlx.NewDb(db, "postgres")
	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery("select exists(select 1 from role where name = $1)").
		WithArgs("test").
		WillReturnRows(rows)

	insertRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("insert into role (name) values ($1) returning id").
		WithArgs("test").
		WillReturnRows(insertRows)

	repo := &Repository{db: sqlxDB}
	service := NewService(repo, val)

	id, err := service.SaveTx("test")
	mock.ExpectCommit()

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestServiceSaveTxBeginTxError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("failed to create mock database")
	}
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := &Repository{db: sqlxDB}
	service := NewService(repo, val)

	mock.ExpectBegin().WillReturnError(fmt.Errorf("tx begin error"))

	id, err := service.SaveTx("test")
	assert.Error(t, err)
	assert.Zero(t, id)
	assert.Contains(t, err.Error(), "error creating transaction")
}

func TestServiceSaveTxFindByNameError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal("failed to create mock database")
	}
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := &Repository{db: sqlxDB}
	service := NewService(repo, val)

	mock.ExpectBegin()

	mock.ExpectQuery("select exists(select 1 from role where name = $1)").
		WithArgs("test").
		WillReturnError(fmt.Errorf("find by name error"))

	mock.ExpectRollback()

	id, err := service.SaveTx("test")
	assert.Error(t, err)
	assert.Zero(t, id)
	assert.Contains(t, err.Error(), "error finding role by name")
}

func TestServiceSaveTxroleAlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal("failed to create mock database")
	}
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := &Repository{db: sqlxDB}
	service := NewService(repo, val)

	mock.ExpectBegin()

	mock.ExpectQuery("select exists(select 1 from role where name = $1)").
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectRollback()

	id, err := service.SaveTx("test")
	assert.Error(t, err)
	assert.Zero(t, id)
	assert.Contains(t, err.Error(), "already exists")
}

func TestServiceSaveTxSaveError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal("failed to create mock database")
	}
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := &Repository{db: sqlxDB}
	service := NewService(repo, val)

	mock.ExpectBegin()

	mock.ExpectQuery("select exists(select 1 from role where name = $1)").
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectQuery("insert into role (name) values ($1) returning id").
		WithArgs("test").
		WillReturnError(fmt.Errorf("save error"))

	mock.ExpectRollback()

	id, err := service.SaveTx("test")
	assert.Error(t, err)
	assert.Zero(t, id)
	assert.Contains(t, err.Error(), "error creating role")
}

func TestServices(t *testing.T) {
	var a = assert.New(t)
	t.Run("should return found role by id", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)
		var entity = Entity{
			Id:        1,
			Name:      "Разработчик",
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
		var svc = NewService(repo, val)
		var entity = Entity{}
		var err = errors.New("user not found")
		var want = fmt.Errorf("error finding role with id 1: %w", err)

		repo.On("FindById", valueId).Return(entity, err)
		var response, got = svc.FindById(1)

		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("should return all found roles by ids", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)
		var entityes = []Entity{
			{
				Id:        1,
				Name:      "Разработчик",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Id:        2,
				Name:      "Админ",
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

	t.Run("should return an error when not found by ids", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)
		var entity = []Entity{{
			Id:        1,
			Name:      "Разработчик",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}}
		var err = errors.New("users not found")
		var want = fmt.Errorf("error finding roles by ids: %d, %w", []int64{1, 3}, err)

		repo.On("FindAllByIds", []int64{1, 3}).Return(entity, err)
		var response, got = svc.FindAllByIds([]int64{1, 3})

		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindAllByIds", 1))
	})

	t.Run("should return all roles", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)
		var entityes = []Entity{
			{
				Name:      "Разработчик",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Name:      "Админ",
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

	t.Run("should delete all roles by ids", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)

		repo.On("DeleteAllByIds", []int64{1, 2}).Return(nil)
		err := svc.DeleteAllByIds([]int64{1, 2})

		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteAllByIds", 1))
	})

	t.Run("should delete by id", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)

		repo.On("Delete", valueId).Return(nil)
		err := svc.Delete(1)

		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "Delete", 1))
	})

	t.Run("should return saved role", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)
		var entity = Entity{
			Name:      "User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		repo.On("Save", entity).Return(valueId, nil)
		var got, err = svc.Save(entity)

		a.Nil(err)
		a.Equal(valueId, got)
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})

	t.Run("should return error while save role", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)
		var entity = Entity{
			Name: "",
		}
		var err = errors.New("user not saved")
		var want = fmt.Errorf("error saving role name: %s: %w", entity.Name, err)

		repo.On("Save", entity).Return(int64(0), err)
		var result, got = svc.Save(entity)

		a.Equal(int64(0), result)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})
}
