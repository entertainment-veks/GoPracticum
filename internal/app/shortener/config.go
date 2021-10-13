package shortener

import (
	"flag"
	"os"
)

const (
	serverAddressKey   = "SERVER_ADDRESS"
	baseURLKey         = "BASE_URL_KEY"
	fileStoragePathKey = "FILE_STORAGE_PATH"
	databaseDSNKey     = "DATABASE_DSN"
)

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
	if len(os.Getenv(serverAddressKey)) != 0 {
		c.ServerAddress = os.Getenv(serverAddressKey)
	}

	if len(os.Getenv(baseURLKey)) != 0 {
		c.BaseURL = os.Getenv(baseURLKey)
	}

	if len(os.Getenv(fileStoragePathKey)) != 0 {
		c.FileStoragePath = os.Getenv(fileStoragePathKey)
	}

	if len(os.Getenv(databaseDSNKey)) != 0 {
		c.DatabaseURL = os.Getenv(databaseDSNKey)
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
