package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"idm/inner/common"
	"idm/inner/database"
	"os"
	"testing"
)

var cfg = common.GetConfig(".env")

func TestConnectionDb(t *testing.T) {

	//1. в корне проекта нет .env файла
	t.Run("not exist test.env file", func(t *testing.T) {
		_, err := os.Stat("test.env")

		assert.Error(t, err)
	})

	//2. в корне проекта есть .env  файл, но в нём нет нужных переменных и в переменных окружения их тоже нет
	t.Run("exist .env file, but no required vars and no required env-variables", func(t *testing.T) {
		defer func() {
			_ = os.Remove("test2.env")
		}()
		_ = os.WriteFile("test2.env", []byte(""), 0644)
		test2, err := os.ReadFile("test2.env")
		_ = os.Unsetenv("DB_DRIVER_NAME")
		_ = os.Unsetenv("DB_DSN")

		assert.Empty(t, "", test2)
		assert.NoError(t, err)
		assert.Equal(t, "", os.Getenv("DB_DRIVER_NAME"))
		assert.Equal(t, "", os.Getenv("DB_DSN"))
	})

	//3. в корне проекта есть .env  файл и в нём нет нужных переменных, но в переменных окружения они есть
	t.Run("exist .env file, but no required vars and exist required env-variables", func(t *testing.T) {
		defer func() {
			_ = os.Remove("test2.env")
		}()
		_ = os.WriteFile("test2.env", []byte(""), 0644)
		test2, err := os.ReadFile("test2.env")

		t.Setenv("DB_DRIVER_NAME", "driver name")
		t.Setenv("DB_DSN", "dsn")

		assert.Empty(t, "", test2)
		assert.NoError(t, err)
		assert.Equal(t, "driver name", os.Getenv("DB_DRIVER_NAME"))
		assert.Equal(t, "dsn", os.Getenv("DB_DSN"))
	})

	/* 4. в корне проекта есть .env  файл и в нём есть нужные переменные, но в переменных окружения они тоже есть
	(с другими значениями) - проверить, какие значения будут использованы приложением при подключении к базе данных*/
	t.Run("required variables in .env file conflicting env vars", func(t *testing.T) {
		defer func() {
			_ = os.Remove(".env")
		}()
		var currentDb string
		_ = os.WriteFile(".env", []byte("DB_DRIVER_NAME=postgres\nDB_DSN=host=127.0.0.1 port=5432 user=postgres password=postgres dbname=db_from_file sslmode=disable"), 0644)
		t.Setenv("DB_DRIVER_NAME", "postgres")
		t.Setenv("DB_DSN", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

		db, err := database.ConnectDb()
		require.NoError(t, err)
		err = db.QueryRow("SELECT current_database()").Scan(&currentDb)
		assert.NoError(t, err)
		assert.Equal(t, "postgres", currentDb)
	})

	//5. в корне проекта есть корректно заполненный .env файл, в переменных окружения нет конфликтующих с ним переменных
	t.Run(".env file exists and required vars no conflicting with env vars", func(t *testing.T) {
		t.Setenv("DB_DRIVER_NAME", "postgres")
		t.Setenv("DB_DSN", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

		assert.Equal(t, cfg.DbDriverName, os.Getenv(cfg.DbDriverName))
		assert.Equal(t, cfg.Dsn, os.Getenv(cfg.Dsn))
	})

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
