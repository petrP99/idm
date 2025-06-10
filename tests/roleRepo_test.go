package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/employee"
	"idm/inner/role"
	"testing"
)

func TestRoleTransaction(t *testing.T) {

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

func TestRepositoryRole(t *testing.T) {

	t.Run("save role", func(t *testing.T) {
		fixture := NewFixture()

		entity := role.Entity{Name: "Директор"}
		id, err := fixture.RoleRepo.Save(entity)
		result, _ := fixture.RoleRepo.FindById(id)

		assert.NoError(t, err)
		assert.Equal(t, "Директор", result.Name)
	})

	t.Run("find role by id", func(t *testing.T) {
		fixture := NewFixture()

		result, err := fixture.RoleRepo.FindById(1)

		assert.NoError(t, err)
		assert.Equal(t, "Администратор", result.Name)
		assert.Equal(t, int64(1), result.Id)
	})

	t.Run("find all role", func(t *testing.T) {
		fixture := NewFixture()

		result, err := fixture.RoleRepo.FindAll()

		assert.NoError(t, err)
		assert.Equal(t, 3, len(result))
	})

	t.Run("find all role by ids", func(t *testing.T) {
		fixture := NewFixture()

		result, err := fixture.RoleRepo.FindAllByIds([]int64{1, 2})

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Equal(t, 2, len(result))
	})

	t.Run("delete all by ids", func(t *testing.T) {
		fixture := NewFixture()

		err := fixture.RoleRepo.DeleteAllByIds([]int64{1, 2, 3})
		entity, _ := fixture.RoleRepo.FindAll()

		assert.NoError(t, err)
		assert.Empty(t, 0, len(entity))
	})

	t.Run("delete by id", func(t *testing.T) {
		fixture := NewFixture()

		errDelete := fixture.RoleRepo.Delete(2)
		_, errFind := fixture.RoleRepo.FindById(2)

		assert.NoError(t, errDelete)
		assert.Error(t, errFind)
	})
}
