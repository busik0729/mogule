package eventPack

import (
	"github.com/go-pg/pg"
	"github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
	"../../structs/fs"
	"../../structs/paginator"
)

/**
Model methods
*/
func GetAll() (Events, error) {
	var u Events

	if err := database.GetConnection().Con.Model(&u).Select(); err != nil {
		return u, err
	}
	return u, nil
}

func GetAllWithFS(fs fs.FS, p paginator.Paginator) (Events, error, int) {
	var u Events

	q := database.GetConnection().Con.Model(&u)
	q = helpers.SetFSAndPG(q, fs, p)
	c, err := q.SelectAndCount()
	if err != nil {
		return u, err, c
	}

	return u, nil, c
}

func GetById(id *uuid.UUID) (Event, error) {
	u := new(Event)
	err := database.GetConnection().Con.Model(u).Where("id = ?", id).Select()
	if err != nil {
		return *u, err
	}

	return *u, err
}

func GetByIds(ids []*uuid.UUID) (Events, error) {
	u := new(Events)
	err := database.GetConnection().Con.Model(u).Select(pg.Array(&ids))
	if err != nil {
		return *u, err
	}

	return *u, nil
}

func GetNewEventsByDevice(deviceId uuid.UUID) (Events, error) {
	u := new(Events)
	if err := database.GetConnection().Con.Model(u).
		Where("device_id = ?", deviceId).
		Where("status = ? OR status = ?", GetNewStatus(), GetSendedStatus()).
		Where("type_event = ?", GetPersonalTypeEvent()).
		Select(); err != nil {
		return *u, err
	}

	return *u, nil
}

func Update(event *Event) (*Event, error) {
	if _, err := database.GetConnection().Con.Model(event).Where("id = ?", event.Id).Update(); err != nil {
		return event, err
	}

	return event, nil
}
