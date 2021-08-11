package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const ServiceName = "elrond-monitor-backend"

type (
	Config struct {
		API                    API
		Postgres               Postgres
		ElasticSearch          ElasticSearch
		MarketProvider         MarketProvider
		Parser                 Parser
		Contracts              Contracts
		StakingProvidersSource string
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
	MarketProvider struct {
		Title  string
		APIKey string
	}
	Parser struct {
		Node     string
		Batch    uint64
		Fetchers uint64
	}
	ElasticSearch struct {
		Address string
	}
	Contracts struct {
		Staking           string
		DelegationManager string
		Delegation        string
		Auction           string
	}
)

func GetConfigFromFile(fileName string) (cfg Config, err error) {
	path, _ := filepath.Abs(fileName)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("ioutil.ReadFile: %s", err.Error())
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return cfg, fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	return config, config.Validate()
}

// Validate validates all Config fields.
func (config *Config) Validate() error {
	err := config.API.validate()
	if err != nil {
		return err
	}
	err = config.Parser.validate()
	if err != nil {
		return err
	}
	if config.StakingProvidersSource == "" {
		return fmt.Errorf("StakingProvidersSource is empty")
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

func (config *Parser) validate() error {
	if config.Batch == 0 {
		return fmt.Errorf("batch is zero")
	}
	if config.Fetchers == 0 {
		return fmt.Errorf("fetchers is zero")
	}
	if config.Node == "" {
		return fmt.Errorf("fetchers is empty")
	}
	return nil
}
