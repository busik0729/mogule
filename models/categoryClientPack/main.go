package categoryClientPack

import (
	"reflect"
	"strings"

	"github.com/go-pg/pg/orm"
	"github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
)

const CACHEKEY = "CategoryClient-"
const CACHEKEY_COLLECTION = "CategoryClients"

/**
Model struct
*/
type CategoryClient struct {
	tableName struct{}   `sql:"category_client,alias:category_client"`
	Id        *uuid.UUID `valid:"uuid,optional" json:"id" sql:"type:uuid,pk" renderFor:"ALL"`
	Title     string     `valid:"ascii" json:"title" renderFor:"ALL"`
}

type CategoryClients []CategoryClient

func (CategoryClient) TableName() string {
	return "category_client"
}

func (b *CategoryClient) BeforeCreate(db orm.DB) error {
	if b.Id == nil {
		_uuid, _ := uuid.NewV4()
		b.Id = &_uuid
	}

	return nil
}

func (u *CategoryClient) GetMapKey() string {
	return "CategoryClient"
}

func (c *CategoryClients) GetMapKey() string {
	return "CategoryClients"
}

func (u *CategoryClient) SetUUID() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (d *CategoryClient) Create() (orm.Result, error) {
	res, err := database.GetConnection().Con.Model(d).Insert()
	RefreshCacheByALL()
	return res, err
}

func (u *CategoryClient) GetCacheKey() string {
	return CACHEKEY + u.Id.String()
}

func (u CategoryClient) Render(role string) map[string]interface{} {
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

func (us CategoryClients) Render(role string) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < len(us); i++ {
		data = append(data, us[i].Render(role))
	}

	return data
}
