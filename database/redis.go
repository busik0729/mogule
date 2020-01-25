package database

import (
	"../config"
	"../helpers"
	"./structs"
	"github.com/go-redis/redis"
)

var rd structs.Redis = structs.Redis{}

func GetConnRedis() structs.Redis {
	return rd
}

func InitializeRedis(env string) structs.Redis {
	cf := config.LoadConfiguration(env)

	address := cf.Redis.Address
	password := cf.Redis.Password
	DB := cf.Redis.DB

	cl := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       DB,       // use default DB
	})

	_, err := cl.Ping().Result()
	if err != nil {
		panic("Redis connect is refused")
	}

	rd.Con = cl
	helpers.GetDI().DIInjector.Map(rd)

	return rd
}
