package mailPack

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
)

/**
Model struct
*/
type Mail struct {
	tableName  struct{}   `sql:"mail"`
	Id         *uuid.UUID `json:"id" renderFor:"ALL"`
	Text       string     `json:"text" renderFor:"ALL"`
	TemplateId *uuid.UUID `json:"template_id" renderFor:"ALL"`
	ClientIds  []string   `json:"client_ids" renderFor:"ALL" sql:"client_ids,array"`
	Bid        int        `sql:"-"`
	DeletedAt  time.Time  `json:"deleted_at" renderFor:"admin" pg:",soft_delete"`
}

type Mails []Mail

func (Mail) TableName() string {
	return "mail"
}

func (b *Mail) BeforeCreate(db orm.DB) error {
	if b.Id == nil {
		_uuid, _ := uuid.NewV4()
		b.Id = &_uuid
	}

	return nil
}

func (b *Mail) BeforeUpdate(ctx context.Context) (context.Context, error) {

	return ctx, nil
}

func (b *Mail) AfterCreate(ctx context.Context) (context.Context, error) {

	return ctx, nil
}

func (c *Mail) GetMapKey() string {
	return "Mail"
}

func (c *Mails) GetMapKey() string {
	return "Mails"
}

func (u *Mail) SetUUID() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (d *Mail) Create() (orm.Result, error) {
	return database.GetConnection().Con.Model(d).Insert()
}

func (u Mail) Render(role string) map[string]interface{} {
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

func (us Mails) Render(role string) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < len(us); i++ {
		data = append(data, us[i].Render(role))
	}

	return data
}
