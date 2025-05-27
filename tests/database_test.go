package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"idm/inner/common"
	"idm/inner/database"
	"testing"
)

func TestConnectionDb(t *testing.T) {

	//6.приложение не может подключиться к базе данных с некорректным конфигом (например, неправильно указан: хост, порт, имя базы данных, логин или пароль)
	t.Run("invalid data, couldn't connect", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("The code did not panic")
			} else {
				assert.Contains(t, r.(error).Error(), "password authentication failed")
			}
		}()

		wrongCfg := common.Config{
			DbDriverName: "postgres",
			Dsn:          "host=127.0.0.1 port=5432 user=postgres password=7 dbname=postgres sslmode=disable",
		}

		_ = database.ConnectDbWithCfg(wrongCfg) // Должен вызвать панику
	})

	//7.приложение может подключиться к базе данных с корректным конфигом
	t.Run("successful connection to db with correct configs", func(t *testing.T) {
		t.Setenv("DB_DRIVER_NAME", "postgres")
		t.Setenv("DB_DSN", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

		var user string
		db, err := database.ConnectDb()
		require.NoError(t, err)

		err = db.QueryRow("SELECT current_user").Scan(&user)
		assert.NoError(t, err)
		assert.Equal(t, "postgres", user)
	})
}
