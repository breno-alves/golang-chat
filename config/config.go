package config

import (
	"errors"
	"os"
)

type Config struct {
	DB        *DBConfig
	Cache     *CacheConfig
	Exchanger *ExchangerConfig
}

type ExchangerConfig struct {
	Host string
}

func (q *ExchangerConfig) Validate() error {
	if q.Host == "" {
		return errors.New("queue host is required")
	}
	return nil
}

type CacheConfig struct {
	Host     string
	Password string
	Db       int
}

func (c *CacheConfig) Validate() error {
	if c.Host == "" {
		return errors.New("cache host is required")
	}
	return nil
}

type DBConfig struct {
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	Charset  string
}

func (c *DBConfig) Validate() error {
	if c.Host == "" {
		return errors.New("DB host is required")
	}

	if c.Port == 0 {
		return errors.New("DB port is required")
	}

	if c.Username == "" {
		return errors.New("DB username is required")
	}
	if c.Password == "" {
		return errors.New("DB password is required")
	}

	if c.Name == "" {
		return errors.New("DB name is required")
	}

	if c.Charset == "" {
		return errors.New("DB charset is required")
	}
	return nil
}

func GetConfig() *Config {
	dbConfig := &DBConfig{
		Dialect:  "postgres",
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		Charset:  "utf8",
	}
	//err := dbConfig.Validate()
	//if err != nil {
	//	panic(err)
	//}

	cacheConfig := &CacheConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		Db:       0,
	}
	//err = cacheConfig.Validate()
	//if err != nil {
	//	panic(err)
	//}

	exchangerConfig := &ExchangerConfig{
		Host: os.Getenv("EXCHANGER_HOST"),
	}
	//err = exchangerConfig.Validate()
	//if err != nil {
	//	panic(err)
	//}
	return &Config{
		DB:        dbConfig,
		Cache:     cacheConfig,
		Exchanger: exchangerConfig,
	}
}
