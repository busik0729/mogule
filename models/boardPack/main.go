package boardPack

import (
	"context"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
	"../labelsPack"
	"../listPack"
	"../userPack"
)

const CACHEKEY = "Board-"

var DAEMON_STATUS = map[int]string{1: "np", 2: "parse"}

/**
Model struct
*/
type Board struct {
	tableName    struct{}           `sql:"board,alias:board"`
	Id           *uuid.UUID         `valid:"uuid,optional" json:"id" renderFor:"ALL"`
	Name         string             `valid:"ascii" json:"name" renderFor:"ALL"`
	Uri          string             `valid:"ascii" json:"uri" renderFor:"ALL"`
	Settings     Settings           `json:"settings" renderFor:"ALL"`
	Lists        listPack.Lists     `json:"lists" renderFor:"ALL" ifNull:"[]"`
	Pm           *uuid.UUID         `valid:"uuid" json:"pm" renderFor:"ALL"`
	Bid          int                `sql:"-"`
	DeletedAt    time.Time          `valid:"-" json:"deleted_at" renderFor:"admin" pg:",soft_delete"`
	Members      []userPack.User    `sql:"-" json:"members" renderFor:"ALL"`
	Labels       []labelsPack.Label `sql:"-" json:"labels" renderFor:"ALL"`
	Due          time.Time          `json:"due" renderFor:"ALL"`
	ClientId     *uuid.UUID         `valid:"uuid" json:"client_id" renderFor:"ALL" render_key:"client_id"`
	DaemonStatus int                `json:"daemon_status" renderFor:"NOBODY" render_key:"daemon_status"`
}

type Settings struct {
	Color           string `json:"color" renderFor:"ALL"`
	Subscribed      bool   `json:"subscribed" renderFor:"ALL"`
	CardCoverImages bool   `json:"cardCoverImages" renderFor:"ALL"`
}

type Boards []Board

func (Board) TableName() string {
	return "board"
}

func (b *Board) BeforeCreate(db orm.DB) error {
	if b.Id == nil {
		_uuid, _ := uuid.NewV4()
		b.Id = &_uuid
	}

	return nil
}

func (b *Board) BeforeUpdate(ctx context.Context) (context.Context, error) {
	log.Println("............................................................")
	log.Println(ctx)
	log.Println(b)
	log.Println("............................................................")
	return ctx, nil
}

func (b *Board) AfterCreate(ctx context.Context) (context.Context, error) {

	return ctx, nil
}

func (c *Board) GetMapKey() string {
	return "Board"
}

func (c *Boards) GetMapKey() string {
	return "Boards"
}

func (u *Board) SetPM(pm userPack.User) {
	u.Pm = pm.Id
}

func (u *Board) SetClient(id *uuid.UUID) {
	u.Pm = id
}

func (u *Board) SetUUID() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (d *Board) Create() (orm.Result, error) {
	return database.GetConnection().Con.Model(d).Insert()
}

func (u *Board) SetDefault() {
	u.Settings = Settings{Color: "fuse-dark", Subscribed: false, CardCoverImages: true}
	u.SetDaemonStatus("np")
	u.Due = time.Now()
}

func (u *Board) SetDaemonStatus(status string) {
	indexOfStatus := helpers.IndexOf(status, DAEMON_STATUS)
	if indexOfStatus != -1 {
		u.DaemonStatus = indexOfStatus
	}
}

func (u *Board) SetParseStatus() {
	indexOfStatus := helpers.IndexOf("parse", DAEMON_STATUS)
	if indexOfStatus != -1 {
		u.DaemonStatus = indexOfStatus
	}
}

func (u *Board) GetCacheKey() string {
	return CACHEKEY + u.Id.String()
}

func (u *Board) GetCacheExpire() time.Duration {
	return time.Hour * 24
}

func (u Board) Render(role string) map[string]interface{} {
	var data = make(map[string]interface{})
	val := reflect.ValueOf(&u).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		typeFieldName := typeField.Name
		keyName, ok := typeField.Tag.Lookup("render_key")
		if ok != true {
			keyName = strings.ToLower(typeFieldName)
		}
		methodName := "Get" + strings.Title(typeFieldName)

		valueField := val.Field(i)

		c := valueField.Convert(valueField.Type())
		if c.CanInterface() {
			l := c.Interface()
			if helpers.HasMethod(l, "Render") {
				l := reflect.ValueOf(l).MethodByName("Render")
				e := l.Call([]reflect.Value{reflect.ValueOf(role)})
				valueField = e[0]
			} else if helpers.HasMethod(u, methodName) {
				l := reflect.ValueOf(&u).MethodByName(methodName)
				e := l.Call([]reflect.Value{})
				valueField = e[0]
			}
		}

		tag := typeField.Tag.Get("renderFor")
		aAllows := strings.Split(tag, ",")
		b, _ := helpers.InArray(role, aAllows)
		if tag == "ALL" || b {
			data[keyName] = valueField.Interface()
		}
	}

	return data
}

func (us Boards) Render(role string) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < len(us); i++ {
		data = append(data, us[i].Render(role))
	}

	return data
}
