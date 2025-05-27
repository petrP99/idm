package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/common"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {

	//1. в корне проекта нет .env файла
	t.Run("not exist test.env file", func(t *testing.T) {
		fileInfo, _ := os.Stat("test.env")
		emptyCfg := common.GetConfig("test.env")

		assert.Empty(t, fileInfo)
		assert.Empty(t, emptyCfg.DbDriverName)
		assert.Empty(t, emptyCfg.Dsn)
	})

	//2. в корне проекта есть .env файл, но в нём нет нужных переменных и в переменных окружения их тоже нет
	t.Run("exist .env file, but no required vars and no required env-variables", func(t *testing.T) {
		defer func() {
			_ = os.Remove("test2.env")
		}()
		_ = os.WriteFile("test2.env", []byte(""), 0644)
		_ = os.Unsetenv("DB_DRIVER_NAME")
		_ = os.Unsetenv("DB_DSN")

		fileInfo, _ := os.Stat("test.env")
		emptyCfg := common.GetConfig("test.env")

		assert.Empty(t, fileInfo)
		assert.Empty(t, emptyCfg.DbDriverName)
		assert.Empty(t, emptyCfg.Dsn)
	})

	//3. в корне проекта есть .env  файл и в нём нет нужных переменных, но в переменных окружения они есть
	t.Run("exist .env file, but no required vars and exist required env-variables", func(t *testing.T) {
		defer func() {
			_ = os.Remove("test.env")
		}()
		_ = os.WriteFile("test.env", []byte(""), 0644)
		test, _ := os.ReadFile("test.env")

		t.Setenv("DB_DRIVER_NAME", "driver name")
		t.Setenv("DB_DSN", "dsn")

		correctCfg := common.GetConfig("test.env")

		assert.Empty(t, "", test)
		assert.Equal(t, "driver name", correctCfg.DbDriverName)
		assert.Equal(t, "dsn", correctCfg.Dsn)
	})

	/* 4. в корне проекта есть .env  файл и в нём есть нужные переменные, но в переменных окружения они тоже есть
	(с другими значениями) - проверить, какие значения будут использованы приложением при подключении к базе данных*/
	t.Run("required variables in .env file conflicting env vars", func(t *testing.T) {
		defer func() {
			_ = os.Remove("test.env")
		}()
		_ = os.WriteFile("test.env", []byte("DB_DRIVER_NAME=mysql\nDB_DSN=host=localhost port=5555 user=mysql password=mysql dbname=db_from_file sslmode=disable"), 0644)
		t.Setenv("DB_DRIVER_NAME", "postgres")
		t.Setenv("DB_DSN", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

		cfg := common.GetConfig("test.env")

		assert.Equal(t, os.Getenv("DB_DRIVER_NAME"), cfg.DbDriverName)
		assert.Equal(t, os.Getenv("DB_DSN"), cfg.Dsn)
	})

	//5. в корне проекта есть корректно заполненный .env файл, в переменных окружения нет конфликтующих с ним переменных
	t.Run(".env file exists and required vars no conflicting with env vars", func(t *testing.T) {
		t.Setenv("DB_DRIVER_NAME", "postgres")
		t.Setenv("DB_DSN", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
		expectedCfg := common.Config{
			DbDriverName: "postgres",
			Dsn:          "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
		}

		actualCfg := common.GetConfig(".env")

		assert.Equal(t, expectedCfg, actualCfg)
	})
}
