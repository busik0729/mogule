package activityPack

import (
	uuid "github.com/satori/go.uuid"

	"../../database"
)

/**
Model methods
*/
func GetAll() (Activities, error) {
	var u Activities

	if err := database.GetConnection().Con.Model(&u).Select(); err != nil {
		return u, err
	}
	return u, nil
}

func GetById(id *uuid.UUID) (Activity, error) {
	u := new(Activity)

	err := database.GetConnection().Con.Model(u).Where("id = ?", id).Select()
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func GetByCardId(cardId *uuid.UUID) (Activities, error) {
	u := new(Activities)

	err := database.GetConnection().Con.Model(u).Where("card_id = ?", cardId).Select()
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func Update(a *Activity) (*Activity, error) {
	if _, err := database.GetConnection().Con.Model(a).Where("id = ?", a.Id).Update(); err != nil {
		return a, err
	}

	return a, nil
}
