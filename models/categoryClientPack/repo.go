package categoryClientPack

import (
	"../../database"
	"../../services"
	"encoding/json"
	"errors"
	"log"

	"github.com/satori/go.uuid"
)

/**
Model methods
*/
func GetAll() (CategoryClients, error) {
	ccs, e := GetCategoryClientsByCache()
	if e != nil || ccs == nil {
		var u CategoryClients
		err := database.GetConnection().Con.Model(&u).Select()
		if err != nil {
			return u, err
		}

		services.RedisSet(CACHEKEY_COLLECTION, u)

		return u, err
	}

	return *ccs, nil
}

func GetAllWithoutCache() (CategoryClients, error) {
	var u CategoryClients
	err := database.GetConnection().Con.Model(&u).Select()
	if err != nil {
		return u, err
	}

	return u, err
}

func GetById(id *uuid.UUID) (CategoryClient, error) {
	u, err := GetCategoryClientByCache(id)

	if err != nil {
		u := new(CategoryClient)
		err := database.GetConnection().Con.Model(u).Where("id = ?", id).Select()
		if err != nil {
			return *u, err
		}

		return *u, err
	}

	return *u, nil
}

func Update(client *CategoryClient) (*CategoryClient, error) {
	if _, err := database.GetConnection().Con.Model(client).Where("id = ?", client.Id).Update(); err != nil {
		return client, err
	}

	services.RedisSet(client.GetCacheKey(), client)
	RefreshCacheByALL()

	return client, nil
}

func GetCategoryClientByCache(id *uuid.UUID) (*CategoryClient, error) {
	u := new(CategoryClient)
	u.Id = id

	user, err := services.RedisGet(u.GetCacheKey())
	if err != nil {
		return u, err
	}

	json.Unmarshal([]byte(user), u)

	return u, err
}

func GetCategoryClientsByCache() (*CategoryClients, error) {
	u := new(CategoryClients)

	ccs, err := services.RedisGet(CACHEKEY_COLLECTION)
	if err != nil || len(ccs) < 1 {
		return u, errors.New("Not Found")
	}

	json.Unmarshal([]byte(ccs), u)

	return u, err
}

func RefreshCacheByALL() {
	cls, e := GetAllWithoutCache()
	log.Println(cls)
	if e != nil {
		// write log
		return
	}

	services.RedisSet(CACHEKEY_COLLECTION, cls)
}
