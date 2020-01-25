package labelsPack

import (
	"github.com/helloeave/json"
	uuid "github.com/satori/go.uuid"

	"../../database"
	"../../services"
)

/**
Model methods
*/
func GetAll() (Labels, error) {
	ccs, e := GetLabelsByCache()
	if e != nil || ccs == nil {
		var u Labels
		err := database.GetConnection().Con.Model(&u).Select()
		if err != nil {
			return u, err
		}

		services.RedisSet(u.GetCacheKey(), u)

		return u, err
	}

	return *ccs, nil

}

func GetAllWithoutCache() (Labels, error) {
	var u Labels

	if err := database.GetConnection().Con.Model(&u).Select(); err != nil {
		return u, err
	}
	return u, nil
}

func GetById(id *uuid.UUID) (Label, error) {
	u := new(Label)

	err := database.GetConnection().Con.Model(u).Where("id = ?", id).Select()
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func Update(label *Label) (*Label, error) {
	if _, err := database.GetConnection().Con.Model(label).Where("id = ?", label.Id).Update(); err != nil {
		return label, err
	}

	RefreshCache()

	return label, nil
}

func RefreshCache() {
	labels, e := GetAllWithoutCache()
	if e == nil {
		services.RedisSet(labels.GetCacheKey(), labels)
	}
}

func GetLabelsByCache() (*Labels, error) {
	u := new(Labels)

	user, err := services.RedisGet(u.GetCacheKey())
	if err != nil {
		return u, err
	}

	er := json.Unmarshal([]byte(user), u)
	if er != nil {
		return u, err
	}

	return u, err
}

func Delete(id *uuid.UUID) (*Label, error) {
	u := new(Label)
	u.Id = id

	services.RedisRemove(u.GetCacheKey())

	_, err := database.GetConnection().Con.Model(u).WherePK().Delete()
	if err != nil {
		return nil, err
	}

	RefreshCache()

	return u, nil
}
