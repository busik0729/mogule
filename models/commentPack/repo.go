package commentPack

import (
	uuid "github.com/satori/go.uuid"

	"../../database"
)

/**
Model methods
*/
func GetAll() (Comments, error) {
	var u Comments

	if err := database.GetConnection().Con.Model(&u).Select(); err != nil {
		return u, err
	}
	return u, nil
}

func GetById(id *uuid.UUID) (Comment, error) {
	u := new(Comment)

	err := database.GetConnection().Con.Model(u).Where("id = ?", id).Select()
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func GetByCardId(listId *uuid.UUID) (Comments, error) {
	u := new(Comments)

	err := database.GetConnection().Con.Model(u).Where("card_id = ?", listId).Select()
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func Update(comment *Comment) (*Comment, error) {
	if _, err := database.GetConnection().Con.Model(comment).Where("id = ?", comment.Id).Update(); err != nil {
		return comment, err
	}

	return comment, nil
}
