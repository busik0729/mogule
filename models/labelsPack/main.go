package labelsPack

import (
	"reflect"
	"strings"

	"github.com/go-pg/pg/orm"
	uuid "github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
)

const CACHEKEY = "Labels"

/**
Model struct
*/
type Label struct {
	tableName struct{}   `sql:"label"`
	Id        *uuid.UUID `valid:"uuid,optional" json:"id" renderFor:"ALL"`
	Name      string     `valid:"ascii" json:"name" renderFor:"ALL"`
	Color     string     `valid:"ascii" json:"color" renderFor:"ALL"`
}

type Labels []Label

func (Label) TableName() string {
	return "label"
}

func (b *Label) BeforeCreate(db orm.DB) error {
	if b.Id == nil {
		_uuid, _ := uuid.NewV4()
		b.Id = &_uuid
	}

	return nil
}

func (u *Label) SetUUID() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (u *Label) GetCacheKey() string {
	return CACHEKEY
}

func (u *Labels) GetCacheKey() string {
	return CACHEKEY
}

func (c *Label) GetMapKey() string {
	return "Label"
}

func (c *Labels) GetMapKey() string {
	return "Labels"
}

func (d *Label) Create() (orm.Result, error) {
	res, err := database.GetConnection().Con.Model(d).Insert()
	if err == nil {
		RefreshCache()
	}
	return res, err
}

func (u Label) Render(role string) map[string]interface{} {
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

func (us Labels) Render(role string) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < len(us); i++ {
		data = append(data, us[i].Render(role))
	}

	return data
}
