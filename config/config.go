package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type (
	Config struct {
		API      API
		Postgres Postgres
	}
	API struct {
		ListenOnPort       uint16
		CORSAllowedOrigins []string
	}
	Postgres struct {
		Host     string
		User     string
		Password string
		Database string
		SSLMode  string
	}
)

// load loads config data from any reader or from ENV
func load(cfg interface{}, source io.Reader) error {
	decoder := json.NewDecoder(source)
	err := decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigFromFile(fileName string) (config Config, err error) {
	path, err := filepath.Abs(fileName)
	if err != nil {
		return config, err
	}
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()
	err = load(config, file)
	if err != nil {
		return config, err
	}

	return config, config.Validate()
}

// Validate validates all Config fields.
func (config *Config) Validate() error {
	err := config.API.validate()
	if err != nil {
		return err
	}

	return config.Postgres.validate()
}

func (config *API) validate() error {
	if len(config.CORSAllowedOrigins) == 0 {
		return fmt.Errorf("CORSAllowedOrigins are empty")
	}

	return nil
}

func (config *Postgres) validate() error {
	if len(config.Host) == 0 {
		return fmt.Errorf("host is empty")
	}
	if len(config.User) == 0 {
		return fmt.Errorf("user is empty")
	}
	if len(config.Password) == 0 {
		return fmt.Errorf("password is empty")
	}
	if len(config.Database) == 0 {
		return fmt.Errorf("database is empty")
	}
	if len(config.SSLMode) == 0 {
		return fmt.Errorf("SSLMode is empty")
	}

	return nil
}
