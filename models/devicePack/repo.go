package devicePack

import (
	"github.com/go-pg/pg"

	"../../database"
	"../../services"
	"encoding/json"
	"log"

	"github.com/satori/go.uuid"
)

/**
Model methods
*/
func GetAll() (Devices, error) {
	var u Devices

	if err := database.GetConnection().Con.Select(&u); err != nil {
		return u, err
	}
	return u, nil
}

func GetById(id *uuid.UUID) (Device, error) {

	u := new(Device)
	if err := database.GetConnection().Con.Model(u).Where("id = ?", id).Select(); err != nil {
		return *u, err
	}

	return *u, nil
}

func GetByUserId(id *uuid.UUID) (Device, error) {

	u := new(Device)
	if err := database.GetConnection().Con.Model(u).Where("user_id = ?", id).Select(); err != nil {
		return *u, err
	}

	return *u, nil
}

func GetByUserIds(ids []string) (Devices, error) {

	u := new(Devices)
	if err := database.GetConnection().Con.Model(u).Where("user_id IN (?)", pg.In(ids)).Select(); err != nil {
		return *u, err
	}

	return *u, nil
}

func GetByAccessToken(AT string) (Device, error) {

	cdev, err := GetDeviceByCache(AT)
	if err != nil {
		u := new(Device)
		err := database.GetConnection().Con.
			Model(u).
			Where("access_token = ?", AT).
			Select()
		if err != nil {
			return *u, err
		}

		return *u, err
	}

	return *cdev, nil
}

func GetByRefreshToken(RT string) (Device, error) {
	u := Device{}

	if err := database.GetConnection().Con.Select(&Device{RefreshToken: RT}); err != nil {
		return u, err
	}

	return u, nil
}

func GetByHeader(device []string, userId *uuid.UUID) (Device, error) {
	u := new(Device)
	if err := database.GetConnection().Con.
		Model(u).
		Where("platform = ?", device[0]).
		Where("device_uuid = ?", device[1]).
		Where("model = ?", device[2]).
		Where("serial = ?", device[3]).
		Where("user_id = ?", userId).
		Select(); err != nil {
		return *u, err
	}
	return *u, nil
}

func GetByDeviceAndRefreshToken(device []string, refreshToken string) (Device, error) {
	u := new(Device)
	log.Println("refresh")
	if err := database.GetConnection().Con.
		Model(u).
		Where("platform = ?", device[0]).
		Where("device_uuid = ?", device[1]).
		Where("model = ?", device[2]).
		Where("serial = ?", device[3]).
		Where("refresh_token = ?", refreshToken).
		Select(); err != nil {
		return *u, err
	}

	return *u, nil
}

func GetJWTTokens(d Device) JWTTokens {
	return JWTTokens{AccessToken: d.AccessToken, RefreshToken: d.RefreshToken, ExpiredIn: d.ExpiredIn}
}

func Update(device *Device) (*Device, error) {
	if err := database.GetConnection().Con.Update(device); err != nil {
		return device, err
	}

	services.RedisSet(device.GetCacheKey(), device)

	return device, nil
}

func GetDeviceByCache(AT string) (*Device, error) {
	u := new(Device)
	u.AccessToken = AT

	device, err := services.RedisGet(u.GetCacheKey())
	if err != nil || len(device) < 1 {
		return u, err
	}

	json.Unmarshal([]byte(device), u)

	return u, err
}
