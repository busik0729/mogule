package structs

import (
	"github.com/go-redis/redis"
)

type Redis struct {
	Con *redis.Client
}
