package config

type Config struct {
	DB    *DBConfig
	Cache *Cache
}

type Cache struct {
	Host     string
	Password string
	Db       int
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

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "postgres",
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
			Name:     "mhe",
			Charset:  "utf8",
		},
		Cache: &Cache{
			Host:     "localhost:6379",
			Password: "",
			Db:       0,
		},
	}
}
