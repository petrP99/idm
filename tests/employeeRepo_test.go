package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/employee"
	"testing"
)

func TestEmployeeTransaction(t *testing.T) {

	t.Run("successful transaction", func(t *testing.T) {
		fixture := NewFixture()
		repo := employee.NewEmployeeRepository(fixture.DB)
		entity := employee.Entity{Name: "Test user"}
		tx, err := repo.BeginTransaction()
		assert.NoError(t, err, "BeginTransaction should not return error")

		exists, err := repo.FindByNameTx(tx, "test_employee")
		assert.NoError(t, err, "FindByNameTx should not return error")
		assert.False(t, exists, "Employee should exist within the transaction before rollback")

		id, err := fixture.EmployeesRepo.SaveTx(tx, entity)
		assert.NoError(t, err)
		assert.True(t, id > 0)

		exists, err = repo.FindByNameTx(tx, "Test user")
		assert.NoError(t, err)
		assert.True(t, exists)

		err = tx.Commit()
		assert.NoError(t, err)
	})

	t.Run("rollback on error", func(t *testing.T) {
		fixture := NewFixture()
		repo := employee.NewEmployeeRepository(fixture.DB)

		tx, err := repo.BeginTransaction()
		assert.NoError(t, err, "BeginTransaction should not return error")

		_, err = repo.SaveTx(tx, employee.Entity{Name: "Test user"})
		assert.NoError(t, err, "SaveTx should not return error")

		exists, err := repo.FindByNameTx(tx, "Test user")
		assert.NoError(t, err, "FindByNameTx should not return error")
		assert.True(t, exists, "Employee should exist within the transaction before rollback")

		err = tx.Rollback()
		assert.NoError(t, err, "Rollback should not return error")

		result, err := repo.FindByName("Test user")
		assert.Error(t, err, "FindByName should return error after rollback")
		assert.Equal(t, int64(0), result.Id, "ID should be 0 for non-existent employee")
	})
}

func TestRepositoryEmployee(t *testing.T) {

	t.Run("save employee", func(t *testing.T) {
		fixture := NewFixture()

		entity := employee.Entity{Name: "Test user"}
		id, err := fixture.EmployeesRepo.Save(entity, "Разработчик")
		result, _ := fixture.EmployeesRepo.FindById(id)

		assert.NoError(t, err)
		assert.Equal(t, "Test user", result.Name)
	})

	t.Run("find employee by id", func(t *testing.T) {
		fixture := NewFixture()

		result, err := fixture.EmployeesRepo.FindById(1)

		assert.NoError(t, err)
		assert.Equal(t, "Иванов Петр", result.Name)
		assert.Equal(t, int64(1), result.Id)
	})

	t.Run("find all employee", func(t *testing.T) {
		fixture := NewFixture()

		result, err := fixture.EmployeesRepo.FindAll()

		assert.NoError(t, err)
		assert.Equal(t, 4, len(result))
	})

	t.Run("find all employee by ids", func(t *testing.T) {
		fixture := NewFixture()

		result, err := fixture.EmployeesRepo.FindAllByIds([]int64{1, 2})

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Equal(t, 2, len(result))
	})

	t.Run("delete all by ids", func(t *testing.T) {
		fixture := NewFixture()

		err := fixture.EmployeesRepo.DeleteAllByIds([]int64{1, 2})
		entity, _ := fixture.EmployeesRepo.FindAll()

		assert.NoError(t, err)
		assert.Equal(t, 2, len(entity))
	})

	t.Run("delete by id", func(t *testing.T) {
		fixture := NewFixture()

		errDelete := fixture.EmployeesRepo.Delete(2)
		_, errFind := fixture.EmployeesRepo.FindById(2)

		assert.NoError(t, errDelete)
		assert.Error(t, errFind)
	})
}
