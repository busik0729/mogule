package config

type Database struct {
	DbHost     string
	DbName     string
	DbUser     string
	DbPassword string
	DbSsl      string
}

type Server struct {
	Host string
	Port string
}

type Config struct {
	Database Database
	Server   Server
	Redis    RdConf
}

type RdConf struct {
	Address  string
	Password string
	DB       int
}

type SMTP struct {
	Host string
	User string
	Pass string
	Port int
}

func GetDbConfig(env string) Database {
	host := "localhost"

	if env == "production" {
		host = "postgres"
	}
	return Database{
		host,
		"crm",
		"busik",
		"busik0729",
		"disable",
	}
}

func GetRdConfig(env string) RdConf {
	host := "localhost:6379"

	if env == "production" {
		host = "redis:6379"
	}

	return RdConf{
		host,
		"busik0729",
		0,
	}
}

func GetServerConfig() Server {
	return Server{
		"localhost",
		":9090",
	}
}

func GetSMTPConfig() SMTP {
	return SMTP{
		"localhost",
		"info",
		"busik0729",
		25,
	}
}

func LoadConfiguration(env string) Config {
	db := GetDbConfig(env)
	server := GetServerConfig()
	redis := GetRdConfig(env)
	return Config{
		db,
		server,
		redis,
	}
}
