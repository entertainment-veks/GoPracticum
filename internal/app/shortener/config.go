package shortener

import (
	"flag"
	"os"
)

const SERVER_ADDRESS_KEY = "SERVER_ADDRESS"
const BASE_URL_KEY = "BASE_URL_KEY"
const FILE_STORAGE_PATH_KEY = "FILE_STORAGE_PATH"
const DATABASE_DSN_KEY = "DATABASE_DSN"

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseURL     string
}

func NewConfig() *Config {
	return &Config{
		ServerAddress:   ":8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "file",
		DatabaseURL:     "host=localhost dbname=shortener_db sslmode=disable user=postgres password=postgres",
	}
}

func (c *Config) ConfigureViaEnv() {
	if len(os.Getenv(SERVER_ADDRESS_KEY)) != 0 {
		c.ServerAddress = os.Getenv(SERVER_ADDRESS_KEY)
	}

	if len(os.Getenv(BASE_URL_KEY)) != 0 {
		c.BaseURL = os.Getenv(BASE_URL_KEY)
	}

	if len(os.Getenv(FILE_STORAGE_PATH_KEY)) != 0 {
		c.FileStoragePath = os.Getenv(FILE_STORAGE_PATH_KEY)
	}

	if len(os.Getenv(DATABASE_DSN_KEY)) != 0 {
		c.DatabaseURL = os.Getenv(DATABASE_DSN_KEY)
	}
}

func (c *Config) ConfigureViaFlags() {
	flag.Func("a", "Server address", func(s string) error {
		c.ServerAddress = s
		return nil
	})

	flag.Func("b", "Base url", func(s string) error {
		c.BaseURL = s
		return nil
	})

	flag.Func("f", "File storage path", func(s string) error {
		c.FileStoragePath = s
		return nil
	})

	flag.Func("d", "Database path", func(s string) error {
		c.DatabaseURL = s
		return nil
	})

	flag.Parse()
}
