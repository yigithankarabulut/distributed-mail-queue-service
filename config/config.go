package config

import (
	"errors"
	"os"
	"strconv"
)

// Config struct stores the configuration of the application
type Config struct {
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
	Port     string   `mapstructure:"port"`
	//RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
}

// Database struct stores the configuration of the database
type Database struct {
	Name    string `mapstructure:"name"`
	Host    string `mapstructure:"host"`
	Pass    string `mapstructure:"pass"`
	User    string `mapstructure:"user"`
	Port    string `mapstructure:"port"`
	Migrate bool   `mapstructure:"migrate"`
}

type Redis struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// RabbitMQ struct stores the configuration of the RabbitMQ
type RabbitMQ struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"pass"`
}

func LoadDatabase() (Database, error) {
	var db Database
	db.Name = os.Getenv("DB_NAME")
	db.Host = os.Getenv("DB_HOST")
	db.Pass = os.Getenv("DB_PASS")
	db.User = os.Getenv("DB_USER")
	db.Port = os.Getenv("DB_PORT")
	db.Migrate = os.Getenv("DB_MIGRATE") == "true"
	for _, env := range []string{"DB_NAME", "DB_HOST", "DB_PASS", "DB_USER", "DB_PORT"} {
		if os.Getenv(env) == "" {
			return db, errors.New(env + " is required")
		}
	}
	return db, nil
}

func LoadRedis() (Redis, error) {
	var redis Redis
	redis.Host = os.Getenv("REDIS_HOST")
	redis.Port = os.Getenv("REDIS_PORT")
	for _, env := range []string{"REDIS_HOST", "REDIS_PORT"} {
		if os.Getenv(env) == "" {
			return redis, errors.New(env + " is required")
		}
	}
	return redis, nil
}

func LoadRabbitMQ() (RabbitMQ, error) {
	var rmq RabbitMQ
	rmq.Host = os.Getenv("RABBITMQ_HOST")
	rmq.Port = os.Getenv("RABBITMQ_PORT")
	rmq.User = os.Getenv("RABBITMQ_USER")
	rmq.Password = os.Getenv("RABBITMQ_PASS")
	for _, env := range []string{"RABBITMQ_HOST", "RABBITMQ_PORT", "RABBITMQ_USER", "RABBITMQ_PASS"} {
		if os.Getenv(env) == "" {
			return rmq, errors.New(env + " is required")
		}
	}
	return rmq, nil
}

func LoadConfig() (*Config, error) {
	var Config Config
	db, err := LoadDatabase()
	if err != nil {
		return nil, err
	}
	//rmq, err := LoadRabbitMQ()
	//if err != nil {
	//	return nil, err
	//}
	redis, err := LoadRedis()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("PORT is required")
	}
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.New("PORT must be a number")
	}
	if portNum < 1 || portNum > 65535 {
		return nil, errors.New("PORT must be between 1 and 65535")
	}
	Config.Database = db
	Config.Redis = redis
	//Config.RabbitMQ = rmq
	Config.Port = port
	return &Config, nil
}
