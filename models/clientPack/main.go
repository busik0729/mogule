package clientPack

import (
	"../../database"
	"../../helpers"
	"../userPack"
	"context"
	"github.com/go-pg/pg/orm"
	"github.com/satori/go.uuid"
	"reflect"
	"strings"
	"time"
)

const CACHEKEY = "Client-"

/**
Model struct
*/
type Client struct {
	tableName      struct{}   `sql:"client,alias:client"`
	Id             *uuid.UUID `valid:"uuid,optional" json:"id" sql:"type:uuid,pk" renderFor:"ALL"`
	Name           string     `valid:"ascii" json:"name" renderFor:"ALL"`
	Surname        string     `valid:"ascii" json:"surname" renderFor:"ALL"`
	Thirdname      string     `valid:"ascii" json:"thirdname" renderFor:"ALL"`
	Email          string     `valid:"ascii" json:"email" renderFor:"ALL"`
	Phone          string     `valid:"ascii" json:"phone" renderFor:"ALL"`
	VK             string     `valid:"ascii" json:"vk" renderFor:"ALL"`
	Instagram      string     `valid:"ascii" json:"instagram" renderFor:"ALL"`
	Facebook       string     `valid:"ascii" json:"facebook" renderFor:"ALL"`
	Whatsapp       string     `valid:"ascii" json:"whatsapp" renderFor:"ALL"`
	Telegram       string     `valid:"ascii" json:"telegram" renderFor:"ALL"`
	Twitter        string     `valid:"ascii" json:"twitter" renderFor:"ALL"`
	Odnoklassniki  string     `valid:"ascii" json:"odnoklassniki" renderFor:"ALL"`
	CategoryId     *uuid.UUID `valid:"uuid" json:"category_id" sql:"type:uuid" renderFor:"ALL" render_key:"category_id"`
	CategoryString string     `valid:"optional" json:"category_string" renderFor:"ALL" sql:"-"`
	LastCom        string     `valid:"ascii" json:"last_com" renderFor:"ALL" render_key:"last_com"`
	ResultCom      string     `valid:"ascii" json:"result_com" renderFor:"ALL" render_key:"result_com"`
	Manager        *uuid.UUID `valid:"uuid" json:"manager" sql:"type:uuid" renderFor:"ALL"`

	CreatedAt time.Time `valid:"-" json:"created_at" renderFor:"admin"`
	UpdatedAt time.Time `valid:"-" json:"updated_at" renderFor:"admin"`
	DeletedAt time.Time `valid:"-" json:"deleted_at" renderFor:"admin" pg:",soft_delete"`
}

type Clients []Client

func (Client) TableName() string {
	return "client"
}

func (b *Client) BeforeCreate(db orm.DB) error {
	if b.Id == nil {
		_uuid, _ := uuid.NewV4()
		b.Id = &_uuid
	}

	return nil
}

func (b *Client) BeforeUpdate(ctx context.Context) (context.Context, error) {
	b.UpdatedAt = time.Now()

	return ctx, nil
}

func (c *Client) GetMapKey() string {
	return "Client"
}

func (c *Clients) GetMapKey() string {
	return "Clients"
}

func (u *Client) SetUUID() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (u *Client) SetManager(manager userPack.User) {
	u.Manager = manager.Id
}

func (d *Client) Create() (orm.Result, error) {
	return database.GetConnection().Con.Model(d).Insert()
}

func (u *Client) SetDefault() {
	u.CreatedAt = time.Now()
}

func (u *Client) GetCacheKey() string {
	return CACHEKEY + u.Id.String()
}

func (u *Client) GetCacheExpire() time.Duration {
	return time.Hour * 24
}

func (u Clients) EmailList() []string {
	var list []string
	for _, user := range u {
		if user.Email != "" {
			list = append(list, user.Email)
		}

	}
	return list
}

func (u Client) Render(role string) map[string]interface{} {
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

func (us Clients) Render(role string) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < len(us); i++ {
		data = append(data, us[i].Render(role))
	}

	return data
}
