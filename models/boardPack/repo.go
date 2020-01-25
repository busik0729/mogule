package boardPack

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
	"../../services"
	"../labelsPack"
	"../listPack"
	"../userPack"
)

/**
Model methods
*/
func GetAll() (Boards, error) {
	var u Boards

	if err := database.GetConnection().Con.Model(&u).Order("bid ASC").Select(); err != nil {
		return u, err
	}
	return u, nil
}

func GetById(id *uuid.UUID) (Board, error) {
	u := new(Board)
	u, er := GetBoardByCache(id)

	if er != nil {
		err := database.GetConnection().Con.
			Model(u).
			Relation("Lists", func(q *orm.Query) (*orm.Query, error) {
				return q.Order("position ASC"), nil
			}).
			Relation("Lists.Cards", func(q *orm.Query) (*orm.Query, error) {
				return q.Order("position ASC"), nil
			}).
			Relation("Lists.Cards.Activities").
			Relation("Lists.Cards.Comments").
			Where("id = ?", id).
			Order("bid ASC").
			Select()
		if err != nil {
			return *u, err
		}

		users, _ := userPack.GetAll()
		u.Members = users
		labels, _ := labelsPack.GetAll()
		u.Labels = labels
	}

	return *u, nil
}

func GetAllByDeadline() (Boards, error) {
	var u Boards

	deadlineTwoDayInSec := time.Now().Unix() + 172800
	deadlineTwoDayTime := helpers.GetDateWithZeroHour(deadlineTwoDayInSec)
	deadline := time.Now().Unix()
	deadlineTime := helpers.GetDateWithZeroHour(deadline)

	if err := database.GetConnection().Con.
		Model(&u).
		Where("due BETWEEN ? AND ?", deadlineTime, deadlineTwoDayTime).
		Where("daemon_status = ? OR daemon_status IS NULL", 1).
		Select(); err != nil {
		return u, err
	}
	return u, nil
}

func Update(board *Board) (*Board, error) {
	if _, err := database.GetConnection().Con.Model(board).Where("id = ?", board.Id).Update(); err != nil {
		return board, err
	}

	return board, nil
}

func GetBoardByCache(id *uuid.UUID) (*Board, error) {
	u := new(Board)
	u.Id = id

	user, err := services.RedisGet(u.GetCacheKey())
	if err != nil || len(user) < 1 {
		return u, errors.New("Not Found")
	}

	json.Unmarshal([]byte(user), u)

	return u, err
}

func Delete(id *uuid.UUID) (*Board, error) {
	u := new(Board)
	u.Id = id

	_, err := database.GetConnection().Con.Model(u).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	listPack.DeleteByBoardId(id)

	return u, nil
}

func DeleteAll(cls *Boards) (*Boards, error) {
	_, err := database.GetConnection().Con.Model(cls).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	return cls, nil
}
