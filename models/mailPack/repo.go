package mailPack

import (
	"github.com/satori/go.uuid"

	"../../database"
)

/**
Model methods
*/
func GetAll() (Mails, error) {
	var u Mails

	if err := database.GetConnection().Con.Model(&u).Select(); err != nil {
		return u, err
	}
	return u, nil
}

func GetById(id *uuid.UUID) (Mail, error) {
	u := new(Mail)

	err := database.GetConnection().Con.
		Model(u).
		Where("id = ?", id).
		Select()
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func Update(board *Mail) (*Mail, error) {
	if _, err := database.GetConnection().Con.Model(board).Where("id = ?", board.Id).Update(); err != nil {
		return board, err
	}

	return board, nil
}

func Delete(id *uuid.UUID) (*Mail, error) {
	u := new(Mail)
	u.Id = id

	_, err := database.GetConnection().Con.Model(u).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func DeleteAll(cls *Mails) (*Mails, error) {
	_, err := database.GetConnection().Con.Model(cls).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	return cls, nil
}
