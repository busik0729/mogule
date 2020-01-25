package userPack

import (
	"../../database"
	"../../helpers"
	"../../services"
	"../../structs/fs"
	"../../structs/paginator"
	"encoding/json"
	"github.com/go-pg/pg"
	"github.com/satori/go.uuid"
)

/**
Model methods
*/
func GetAll() (Users, error) {
	var u Users

	if err := database.GetConnection().Con.Model(&u).Where("role <> 1").Select(); err != nil {
		return u, err
	}
	return u, nil
}

func GetAllWithFS(fs fs.FS, p paginator.Paginator) (Users, error, int) {
	var u Users

	q := database.GetConnection().Con.Model(&u).Where("role <> 1")
	q = helpers.SetFSAndPG(q, fs, p)
	c, err := q.SelectAndCount()
	if err != nil {
		return u, err, c
	}

	return u, nil, c
}

func GetById(id *uuid.UUID) (User, error) {
	u, err := GetUserByCache(id)

	if err != nil {
		u := new(User)
		err := database.GetConnection().Con.Model(u).Where("id = ?", id).Select()
		if err != nil {
			return *u, err
		}

		return *u, err
	}

	return *u, nil
}

func GetByIds(ids []*uuid.UUID) (Users, error) {

	u := new(Users)
	err := database.GetConnection().Con.Model(u).Select(pg.Array(&ids))
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func GetByUsername(username string) (User, error) {
	u := new(User)
	if err := database.GetConnection().Con.Model(u).Where("username = ?", username).Select(); err != nil {
		return *u, err
	}

	return *u, nil
}

func Update(user *User) (*User, error) {
	if _, err := database.GetConnection().Con.Model(user).Where("id = ?", user.Id).Update(); err != nil {
		return user, err
	}

	services.RedisSet(user.GetCacheKey(), user)

	return user, nil
}

func GetUserByCache(id *uuid.UUID) (*User, error) {
	u := new(User)
	u.Id = id

	user, err := services.RedisGet(u.GetCacheKey())
	if err != nil {
		return u, err
	}

	json.Unmarshal([]byte(user), u)

	return u, err
}

func GetRoles() []Role {
	roles := []Role{}
	for i, v := range USER_ROLES {
		roles = append(roles, Role{i, v})
	}
	return roles
}

func GetRUSRoles() []Role {
	roles := []Role{}
	for i, v := range USER_ROLES_RUS {
		roles = append(roles, Role{i, v})
	}
	return roles
}

func CreateDefaultUsers() (User, error) {
	u := new(User)
	u.CreateDefaultAdmin()

	if _, err := database.GetConnection().Con.Model(u).Insert(); err != nil {
		return *u, err
	}

	return *u, nil
}
