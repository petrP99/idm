package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/employee"
	"testing"
)

func TestRepositoryEmployee(t *testing.T) {

	t.Run("save employee", func(t *testing.T) {
		fixture := NewFixture()

		entity := employee.EmployeeEntity{Name: "Test user"}
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
