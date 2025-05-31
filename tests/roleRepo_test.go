package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/role"
	"testing"
)

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
