package services

import (
	"../database/structs"
	"../helpers"
	"encoding/json"
	"log"
	"reflect"
	"time"
)

const EXPIRE = time.Minute * 0
const EXPIRED = time.Second * 1

const CACHE_EXPIRE_METHOD_NAME = "GetCacheExpire"

func getRedis(rd structs.Redis) structs.Redis {
	return rd
}

func RedisSet(key string, value interface{}) {
	rd, _ := helpers.GetDI().DIInjector.Invoke(getRedis)
	redis := rd[0].Interface().(structs.Redis)
	val, _ := json.Marshal(value)

	exp := EXPIRE

	if helpers.HasMethod(&value, CACHE_EXPIRE_METHOD_NAME) {
		l := reflect.ValueOf(&value).MethodByName(CACHE_EXPIRE_METHOD_NAME)
		e := l.Call([]reflect.Value{})
		valueField := e[0]

		exp = valueField.Interface().(time.Duration)
	}

	err := redis.Con.Set(key, val, exp).Err()
	if err != nil {
		log.Fatalln(err)
	}
}

func RedisRemove(key string) {
	rd, _ := helpers.GetDI().DIInjector.Invoke(getRedis)
	redis := rd[0].Interface().(structs.Redis)

	err := redis.Con.Set(key, "", EXPIRED).Err()
	if err != nil {
		log.Fatalln(err)
	}
}

func RedisGet(key string) (string, error) {
	rd, _ := helpers.GetDI().DIInjector.Invoke(getRedis)
	redis := rd[0].Interface().(structs.Redis)

	return redis.Con.Get(key).Result()
}
