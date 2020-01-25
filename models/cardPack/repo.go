package cardPack

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
)

/**
Model methods
*/
func GetAll() (Cards, error) {
	var u Cards

	if err := database.GetConnection().Con.Model(&u).Select(); err != nil {
		return u, err
	}
	return u, nil
}

func GetAllByDeadline() (Cards, error) {
	var u Cards

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

func GetById(id *uuid.UUID) (Card, error) {
	u := new(Card)

	err := database.GetConnection().Con.Model(u).Where("id = ?", id).Select()
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func GetByListId(listId *uuid.UUID) (Cards, error) {
	u := new(Cards)

	err := database.GetConnection().Con.Model(u).Where("list_id = ?", listId).Select()
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func Update(card *Card) (*Card, error) {
	if _, err := database.GetConnection().Con.Model(card).Where("id = ?", card.Id).Update(); err != nil {
		return card, err
	}

	return card, nil
}

func Delete(id *uuid.UUID) (*Card, error) {
	u := new(Card)
	u.Id = id

	_, err := database.GetConnection().Con.Model(u).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func DeleteByListdId(id *uuid.UUID) (*Card, error) {
	u := new(Card)

	_, err := database.GetConnection().Con.Model(u).Where("list_id = ?", id).Delete()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func DeleteAll(cls *Cards) (*Cards, error) {
	_, err := database.GetConnection().Con.Model(cls).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	return cls, nil
}
