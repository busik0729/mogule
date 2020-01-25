package listPack

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
	"../cardPack"
)

const CACHEKEY = "List-"

var TRACKING = map[int]string{10: "no", 1: "kpi", 2: "proccess", 3: "done"}

/**
Model struct
*/
type List struct {
	tableName struct{}       `sql:"list"`
	Id        *uuid.UUID     `valid:"uuid,optional" json:"id" renderFor:"ALL"`
	Name      string         `valid:"ascii" json:"name" renderFor:"ALL"`
	BoardId   *uuid.UUID     `valid:"uuid" json:"board_id" render_key:"board_id" renderFor:"ALL"`
	Cards     cardPack.Cards `json:"cards" renderFor:"ALL"`
	Bid       int            `sql:"-"`
	Tracking  int            `json:"tracking" renderFor:"ALL"`
	Position  float64        `json:"position" renderFor:"ALL"`
	DeletedAt time.Time      `valid:"-" json:"deleted_at" renderFor:"admin" pg:",soft_delete"`
}

type Lists []List

func (List) TableName() string {
	return "list"
}

func (b *List) BeforeCreate(db orm.DB) error {
	if b.Id == nil {
		_uuid, _ := uuid.NewV4()
		b.Id = &_uuid
	}

	return nil
}

func (b *List) BeforeUpdate(ctx context.Context) (context.Context, error) {

	return ctx, nil
}

func (b *List) AfterCreate(ctx context.Context) (context.Context, error) {

	return ctx, nil
}

func (c *List) GetMapKey() string {
	return "List"
}

func (c *Lists) GetMapKey() string {
	return "Lists"
}

func (u *List) SetUUID() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (u *List) SetBoardId(uuid *uuid.UUID) {
	u.BoardId = uuid
}

func (d *List) Create() (orm.Result, error) {
	return database.GetConnection().Con.Model(d).Returning("id").Insert()
}

func (u *List) GetCacheKey() string {
	return CACHEKEY + u.Id.String()
}

func (u *List) GetCacheExpire() time.Duration {
	return time.Hour * 24
}

func (u List) Render(role string) map[string]interface{} {
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

func (us Lists) Render(role string) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < len(us); i++ {
		data = append(data, us[i].Render(role))
	}

	return data
}
