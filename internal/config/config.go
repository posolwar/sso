package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

const (
	defaultConfigPath = "config/config.yml"

	EnvConfigPathName  = "CONFIG_PATH"
	FlagConfigPathName = "config"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	StoragePath    string     `yaml:"storage_path" env-required:"true"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func Load() (*Config, error) {
	configPath := fetchConfigPath()
	if configPath == "" {
		return nil, errors.New("Ошибка получения пути к конфигу")
	}

	// Check file exists
	_, err := os.Stat(configPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Ошибка получения файла по указанному пути %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, errors.Wrap(err, "config path is empty")
	}

	return &cfg, nil
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("Ошибка получения пути к конфигу")
	}

	// Check file exists
	_, err := os.Stat(configPath)
	if err != nil {
		panic("Ошибка получения файла по указанному пути: " + configPath + ", ошибка: " + err.Error())
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty, err: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, FlagConfigPathName, defaultConfigPath, "path to config file")
	flag.Parse()

	envConfigPath := os.Getenv(EnvConfigPathName)

	if envConfigPath != "" && res == defaultConfigPath {
		res = os.Getenv(EnvConfigPathName)
	}

	return res
}
