package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

const (
	configJsonPathKay  = "CONFIG_JSON_PATH"
	serverAddressKey   = "SERVER_ADDRESS"
	baseURLKey         = "BASE_URL_KEY"
	fileStoragePathKey = "FILE_STORAGE_PATH"
	databaseDSNKey     = "DATABASE_DSN"
	enableHttpsKey     = "ENABLE_HTTPS"
)

const (
	defServerAddress   = ":8080"
	defBaseURL         = "http://127.0.0.1:8080"
	defFileStoragePath = "file"
	defDatabaseURL     = "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable"
	defEnableHttps     = false
)

var (
	configuredByFlags sync.Once
	configJsonPath    string
)

type Config struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseURL     string `json:"database_url"`
	EnableHTTPS     bool   `json:"enable_https"`
}

var (
	configJsonPathFlagsValue  string
	ServerAddressFlagsValue   string
	BaseURLFlagsValue         string
	FileStoragePathFlagsValue string
	DatabaseURLFlagsValue     string
	EnableHTTPSFlagsValue     bool
)

func NewConfig() *Config {
	c := &Config{
		ServerAddress:   defServerAddress,
		BaseURL:         defBaseURL,
		FileStoragePath: defFileStoragePath,
		DatabaseURL:     defDatabaseURL,
		EnableHTTPS:     defEnableHttps,
	}

	c.configureViaEnv()
	c.configureViaFlags()
	c.configureViaJson() //it should be last

	return c
}

func (c *Config) configureViaJson() {
	configJsonPath = os.Getenv(serverAddressKey)

	data, err := ioutil.ReadFile(configJsonPath)
	if err != nil {
		fmt.Println("Cannot read json config. Ignore if you don't create it", err)
		return
	}

	var newCfg Config
	if err := json.Unmarshal(data, &newCfg); err != nil {
		fmt.Println("Cannot unmarshal json config. Ignore if you don't create it", err)
		return
	}

	if c.ServerAddress == defServerAddress {
		c.ServerAddress = newCfg.ServerAddress
	}
	if c.BaseURL == defBaseURL {
		c.BaseURL = newCfg.BaseURL
	}
	if c.FileStoragePath == defFileStoragePath {
		c.FileStoragePath = newCfg.FileStoragePath
	}
	if c.DatabaseURL == defDatabaseURL {
		c.DatabaseURL = newCfg.DatabaseURL
	}
	if c.EnableHTTPS == defEnableHttps {
		c.EnableHTTPS = newCfg.EnableHTTPS
	}
}

func (c *Config) configureViaEnv() {
	if val := os.Getenv(configJsonPathKay); len(val) != 0 {
		configJsonPath = val
	}

	if val := os.Getenv(serverAddressKey); len(val) != 0 {
		c.ServerAddress = val
	}

	if val := os.Getenv(baseURLKey); len(val) != 0 {
		c.BaseURL = val
	}

	if val := os.Getenv(fileStoragePathKey); len(val) != 0 {
		c.FileStoragePath = val
	}

	if val := os.Getenv(databaseDSNKey); len(val) != 0 {
		c.DatabaseURL = val
	}

	if val := os.Getenv(enableHttpsKey); len(val) != 0 {
		c.EnableHTTPS = true
	}
}

func (c *Config) configureViaFlags() {
	configuredByFlags.Do(func() {
		flag.Func("c", "Config path", func(s string) error {
			configJsonPathFlagsValue = s
			return nil
		})

		flag.Func("a", "Server address", func(s string) error {
			ServerAddressFlagsValue = s
			return nil
		})

		flag.Func("b", "Base url", func(s string) error {
			BaseURLFlagsValue = s
			return nil
		})

		flag.Func("f", "File storage path", func(s string) error {
			FileStoragePathFlagsValue = s
			return nil
		})

		flag.Func("d", "Database path", func(s string) error {
			DatabaseURLFlagsValue = s
			return nil
		})

		flag.Func("s", "Enable HTTPS", func(s string) error {
			EnableHTTPSFlagsValue = true
			return nil
		})

		flag.Parse()
	})

	if len(configJsonPathFlagsValue) != 0 {
		configJsonPath = configJsonPathFlagsValue
	}
	if len(ServerAddressFlagsValue) != 0 {
		c.ServerAddress = ServerAddressFlagsValue
	}
	if len(BaseURLFlagsValue) != 0 {
		c.BaseURL = BaseURLFlagsValue
	}
	if len(FileStoragePathFlagsValue) != 0 {
		c.FileStoragePath = FileStoragePathFlagsValue
	}
	if len(DatabaseURLFlagsValue) != 0 {
		c.DatabaseURL = DatabaseURLFlagsValue
	}
	if EnableHTTPSFlagsValue {
		c.EnableHTTPS = EnableHTTPSFlagsValue
	}
}
