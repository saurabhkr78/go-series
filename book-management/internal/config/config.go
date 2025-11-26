package config

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func LoadConfig() (Config, error) {
	return Config{
		DBUser:     "admin",
		DBPassword: "admin",
		DBHost:     "localhost",
		DBPort:     "5435",
		DBName:     "bookstore",
	}, nil
}
