package clientPack

import (
	"../../database"
	"../../helpers"
	"../../services"
	"../../structs/appCxt"
	"encoding/json"
	"github.com/go-pg/pg"
	"github.com/satori/go.uuid"
)

/**
Model methods
*/
func GetAll(ctx *appCxt.AppContext) (Clients, error, int) {
	var u Clients

	q := database.GetConnection().Con.Model(&u)
	q = helpers.SetFSAndPG(q, ctx.FS, ctx.Paginator)
	c, err := q.SelectAndCount()
	if err != nil {
		return u, err, c
	}

	return u, nil, c
}

func GetById(id *uuid.UUID) (Client, error) {
	u, err := GetClientByCache(id)

	if err != nil {
		u := new(Client)
		err := database.GetConnection().Con.Model(u).Where("id = ?", id).Select()
		if err != nil {
			return *u, err
		}

		return *u, err
	}

	return *u, nil
}

func GetByIds(ids []string) (Clients, error) {

	u := new(Clients)
	err := database.GetConnection().Con.Model(u).Where("id IN (?)", pg.Strings(ids)).Select()
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func Update(client *Client) (*Client, error) {
	if _, err := database.GetConnection().Con.Model(client).Where("id = ?", client.Id).Update(); err != nil {
		return client, err
	}

	services.RedisSet(client.GetCacheKey(), client)

	return client, nil
}

func GetClientByCache(id *uuid.UUID) (*Client, error) {
	u := new(Client)
	u.Id = id

	user, err := services.RedisGet(u.GetCacheKey())
	if err != nil {
		return u, err
	}

	json.Unmarshal([]byte(user), u)

	return u, err
}

func Delete(id *uuid.UUID) (*Client, error) {
	u := new(Client)
	u.Id = id

	services.RedisRemove(u.GetCacheKey())

	_, err := database.GetConnection().Con.Model(u).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func DeleteAll(cls *Clients) (*Clients, error) {

	for _, v := range *cls {
		services.RedisRemove(v.GetCacheKey())
	}

	_, err := database.GetConnection().Con.Model(cls).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	return cls, nil
}
