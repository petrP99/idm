package employee

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

type stubRepository struct {
	entities []Entity
	err      error
}

func (r *stubRepository) FindAllByIds([]int64) ([]Entity, error) {
	return r.entities, r.err
}

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
	mock.ExpectQuery("select exists(select 1 from employee where name = $1)").
		WithArgs("test").
		WillReturnRows(rows)

	insertRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("insert into employee (name) values ($1) returning id").
		WithArgs("test").
		WillReturnRows(insertRows)

	repo := &Repository{db: sqlxDB}
	v := validator.New()
	service := NewService(repo, v)

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
	v := validator.New()
	service := NewService(repo, v)

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
	v := validator.New()
	service := NewService(repo, v)

	mock.ExpectBegin()

	mock.ExpectQuery("select exists(select 1 from employee where name = $1)").
		WithArgs("test").
		WillReturnError(fmt.Errorf("find by name error"))

	mock.ExpectRollback()

	id, err := service.SaveTx("test")
	assert.Error(t, err)
	assert.Zero(t, id)
	assert.Contains(t, err.Error(), "error finding employee by name")
}

func TestServiceSaveTxEmployeeAlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal("failed to create mock database")
	}
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := &Repository{db: sqlxDB}
	v := validator.New()
	service := NewService(repo, v)

	mock.ExpectBegin()

	mock.ExpectQuery("select exists(select 1 from employee where name = $1)").
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
	v := validator.New()
	service := NewService(repo, v)

	mock.ExpectBegin()

	mock.ExpectQuery("select exists(select 1 from employee where name = $1)").
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectQuery("insert into employee (name) values ($1) returning id").
		WithArgs("test").
		WillReturnError(fmt.Errorf("save error"))

	mock.ExpectRollback()

	id, err := service.SaveTx("test")
	assert.Error(t, err)
	assert.Zero(t, id)
	assert.Contains(t, err.Error(), "error creating employee")
}

func TestFindAllByIds(t *testing.T) {
	t.Run("success with stub", func(t *testing.T) {
		// Подготовка данных
		stub := &stubRepository{
			entities: []Entity{
				{
					Name: "Alice",
				},
				{
					Name: "Bob",
				},
			},
			err: nil,
		}
		service := &ServiceStub{repo: stub}

		ids := []int64{1, 2}
		result, err := service.repo.FindAllByIds(ids)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Alice", result[0].Name)
		assert.Equal(t, "Bob", result[1].Name)
	})

	t.Run("error with stub", func(t *testing.T) {
		var err = errors.New("users not found")
		ids := []int64{1, 2}
		var want = fmt.Errorf("error finding employees by ids: %d, %w", ids, err)
		stub := &stubRepository{
			entities: nil,
			err:      want,
		}

		service := &ServiceStub{repo: stub}
		_, expectedErr := service.repo.FindAllByIds(ids)

		assert.Error(t, expectedErr)
		assert.Equal(t, want, expectedErr)
	})
}

func TestServices(t *testing.T) {

	var a = assert.New(t)
	var val = validator.New()

	t.Run("should return found employee by id", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)
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
		var svc = NewService(repo, val)
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
		var svc = NewService(repo, val)
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
		var svc = NewService(repo, val)
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
		var svc = NewService(repo, val)

		repo.On("DeleteAllByIds", []int64{1, 2}).Return(nil)
		err := svc.DeleteAllByIds([]int64{1, 2})

		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "DeleteAllByIds", 1))
	})

	t.Run("should return an error when not found by ids", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)
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
		var svc = NewService(repo, val)

		repo.On("Delete", valueId).Return(nil)
		err := svc.Delete(1)

		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "Delete", 1))
	})

	t.Run("should return saved employee", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, val)
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
		var svc = NewService(repo, val)
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
