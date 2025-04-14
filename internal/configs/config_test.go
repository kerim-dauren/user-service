package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("LoadFromEnv", func(t *testing.T) {
		// Set environment variables
		os.Setenv("APP_NAME", "test-service")
		os.Setenv("APP_ENV", "test")
		os.Setenv("HTTP_PORT", "9090")
		os.Setenv("DB_URL", "postgres://user:password@localhost:5432/testdb")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("LOG_HANDLER", "json")
		os.Setenv("LOG_WRITER", "stderr")
		defer func() {
			os.Unsetenv("APP_NAME")
			os.Unsetenv("APP_ENV")
			os.Unsetenv("HTTP_PORT")
			os.Unsetenv("DB_URL")
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("LOG_HANDLER")
			os.Unsetenv("LOG_WRITER")
		}()

		// Load configuration
		cfg, err := LoadConfig()
		assert.NoError(t, err)

		// Assert values
		assert.Equal(t, "test-service", cfg.AppName)
		assert.Equal(t, "test", cfg.AppEnv)
		assert.Equal(t, 9090, cfg.HttpPort)
		assert.Equal(t, "postgres://user:password@localhost:5432/testdb", cfg.DbUrl)
		assert.Equal(t, "debug", cfg.Log.Level)
		assert.Equal(t, "json", cfg.Log.Handler)
		assert.Equal(t, "stderr", cfg.Log.Writer)
	})

	t.Run("LoadFromDefaults", func(t *testing.T) {
		// Ensure no environment variables are set
		os.Clearenv()

		// Load configuration
		cfg, err := LoadConfig()
		assert.NoError(t, err)

		// Assert default values
		assert.Equal(t, "user-service", cfg.AppName)
		assert.Equal(t, "dev", cfg.AppEnv)
		assert.Equal(t, 8080, cfg.HttpPort)
		assert.Empty(t, cfg.DbUrl)
		assert.Equal(t, "info", cfg.Log.Level)
		assert.Equal(t, "text", cfg.Log.Handler)
		assert.Equal(t, "stdout", cfg.Log.Writer)
	})
}
