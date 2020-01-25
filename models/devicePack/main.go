package devicePack

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	uuid "github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
	"../../services"
)

const CACHEKEY = "Device-"

/**
Model struct
*/
type Device struct {
	tableName struct{}   `sql:"device,alias:device"`
	Id        *uuid.UUID `valid:"uuid,optional" json:"id" sql:"type:uuid,pk"`

	Platform   string `valid:"ascii" json:"platform"`
	DeviceUuid string `valid:"ascii" json:"device_uuid"`
	Model      string `valid:"ascii" json:"model"`
	Serial     string `valid:"ascii" json:"serial"`
	VersionOS  string `valid:"ascii" json:"version_os"`
	VersionApp string `valid:"ascii" json:"version_app"`
	WsId       string `valid:"ascii" json:"ws_id"`

	UserId       *uuid.UUID
	AccessToken  string    `valid:"ascii,optional" json:"access_token"`
	RefreshToken string    `valid:"ascii,optional" json:"refresh_token"`
	ExpiredIn    int64     `valid:"-"`
	CreatedAt    time.Time `valid:"-" json:"created_at" renderFor:"admin"`
	UpdatedAt    time.Time `valid:"-" json:"updated_at" renderFor:"admin"`
	DeletedAt    time.Time `valid:"-" json:"deleted_at" renderFor:"admin" pg:",soft_delete"`
}

type JWTTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiredIn    int64
}

type Devices []Device

func (Device) TableName() string {
	return "device"
}

func (d *Device) BeforeCreate(db orm.DB) error {
	if d.Id == nil {
		_uuid, _ := uuid.NewV4()
		d.Id = &_uuid
	}

	return nil
}

func (d *Device) SetUUID() {
	id, _ := uuid.NewV4()
	d.Id = &id
}

func (d *Device) SetWsUUID(uuid string) {
	d.WsId = uuid
}

func (d *Device) SetNullWsUUID() {
	d.WsId = ""
}

func (d *Device) GetMapKey() string {
	return "Device"
}

func (u *Devices) GetMapKey() string {
	return "Devices"
}

func (d *Device) Create() (pg.Result, error) {
	return database.GetConnection().Con.Model(d).Insert()
}

func (d *Device) GetCacheKey() string {
	return CACHEKEY + d.AccessToken
}

func (d *Device) GetCacheExpire() time.Duration {
	return time.Hour * 240
}

func (d *Device) SetAccessToken() {
	accessClaims, expireIn := helpers.GetNewAccessClaims(d.Id)
	d.AccessToken = helpers.GenerateAccessToken(accessClaims)
	d.ExpiredIn = expireIn
}

func (d *Device) SetRefreshToken() {
	refreshClaims := helpers.GetNewRefreshClaims(d.Id)
	d.RefreshToken = helpers.GenerateRefreshToken(refreshClaims)
}

func (d *Device) SetTokens() {
	services.RedisRemove(d.GetCacheKey())
	d.SetAccessToken()
	d.SetRefreshToken()
}
