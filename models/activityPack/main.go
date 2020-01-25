package activityPack

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-pg/pg/orm"
	uuid "github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
)

/**
Model struct
*/
type Activity struct {
	tableName struct{}   `sql:"activity"`
	Id        *uuid.UUID `valid:"uuid,optional" json:"id" renderFor:"ALL"`
	Message   string     `valid:"ascii" json:"message" renderFor:"ALL"`
	Time      *time.Time `valid:"ascii" json:"time" renderFor:"ALL"`
	MemberId  *uuid.UUID `json:"member_id" renderFor:"ALL" render_key:"member_id"`
	CardId    *uuid.UUID `json:"card_id" renderFor:"ALL" render_key:"card_id"`
}

type Activities []Activity

func (Activity) TableName() string {
	return "member"
}

func (b *Activity) BeforeCreate(db orm.DB) error {
	if b.Id == nil {
		_uuid, _ := uuid.NewV4()
		b.Id = &_uuid
	}

	return nil
}

func (u *Activity) SetUUID() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (d *Activity) Create() (orm.Result, error) {
	return database.GetConnection().Con.Model(d).Insert()
}

func (u Activity) Render(role string) map[string]interface{} {
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

func (us Activities) Render(role string) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < len(us); i++ {
		data = append(data, us[i].Render(role))
	}

	return data
}
