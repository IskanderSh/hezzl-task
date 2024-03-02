package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string        `yaml:"env"`
	LogLevel      string        `yaml:"log_level"`
	Application   Application   `yaml:"application"`
	Storage       Storage       `yaml:"storage"`
	Cache         Cache         `yaml:"cache"`
	MessageBroker MessageBroker `yaml:"broker"`
	LogStorage    LogStorage    `yaml:"log_storage"`
}

type Application struct {
	Port int `yaml:"port"`
}

type Storage struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Cache struct {
	Host string        `yaml:"host"`
	Port int           `yaml:"port"`
	TTL  time.Duration `yaml:"ttl"`
}

type MessageBroker struct {
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
	Subject string `yaml:"subject"'`
}

type LogStorage struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("file is not exists")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config file")
	}

	return &cfg
}

func fetchConfigPath() string {
	print("fetching config path")
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
