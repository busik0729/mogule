package listPack

import (
	"encoding/json"
	"errors"

	"github.com/satori/go.uuid"

	"../../database"
	"../../services"
	"../cardPack"
)

/**
Model methods
*/
func GetAll() (Lists, error) {
	var u Lists

	if err := database.GetConnection().Con.Model(&u).Select(); err != nil {
		return u, err
	}
	return u, nil
}

func GetById(id *uuid.UUID) (List, error) {
	u, err := GetBoardByCache(id)

	if err != nil {
		u := new(List)
		err := database.GetConnection().Con.Model(u).Where("id = ?", id).Select()
		if err != nil {
			return *u, err
		}

		return *u, err
	}

	return *u, nil
}

func GetByBoardId(bid *uuid.UUID) (Lists, error) {
	u := new(Lists)

	err := database.GetConnection().Con.Model(u).
		Column("list.*").
		Relation("Cards").
		Relation("Cards.Members").
		Relation("Cards.Labels").
		Relation("Cards.Activities").
		Where("board_id = ?", bid).
		Select()
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func Update(list *List) (*List, error) {
	if _, err := database.GetConnection().Con.Model(list).Where("id = ?", list.Id).Update(); err != nil {
		return list, err
	}

	return list, nil
}

func GetBoardByCache(id *uuid.UUID) (*List, error) {
	u := new(List)
	u.Id = id

	user, err := services.RedisGet(u.GetCacheKey())
	if err != nil || len(user) < 1 {
		return u, errors.New("Not Found")
	}

	json.Unmarshal([]byte(user), u)

	return u, err
}

func Delete(id *uuid.UUID) (*List, error) {
	u := new(List)
	u.Id = id

	_, err := database.GetConnection().Con.Model(u).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	_, errList := cardPack.DeleteByListdId(id)
	if errList != nil {
		return nil, err
	}

	return u, nil
}

func DeleteByBoardId(id *uuid.UUID) (*List, error) {
	u := new(List)

	_, err := database.GetConnection().Con.Model(u).Where("board_id = ?", id).Returning("id").OnConflict("").Delete()
	if err != nil {
		return nil, err
	}

	cardPack.DeleteByListdId(u.Id)

	return u, nil
}

func DeleteAll(cls *Lists) (*Lists, error) {
	_, err := database.GetConnection().Con.Model(cls).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	return cls, nil
}
