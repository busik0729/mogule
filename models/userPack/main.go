package userPack

import (
	"../../database"
	"../../helpers"
	"github.com/go-pg/pg/orm"
	"github.com/satori/go.uuid"
	"reflect"
	"strings"
	"time"
)

const CACHEKEY = "User-"

var USER_ROLES = map[int]string{1: "admin", 2: "leadership", 3: "caller", 4: "pm", 5: "manager", 6: "copywriter", 7: "guest"}
var USER_ROLES_RUS = map[int]string{1: "Администратор", 2: "Руководитель", 3: "Отдел продаж", 4: "Менеджер проекта", 5: "Менеджер", 6: "Копирайтер", 7: "Гость"}
var USER_ROLES_INT = []int{1, 2, 3, 4, 5, 6, 7}

/**
Model struct
*/
type User struct {
	Id            *uuid.UUID `valid:"uuid,optional" json:"id" sql:"type:uuid,pk" renderFor:"ALL"`
	Username      string     `valid:"ascii" json:"username" renderFor:"ALL"`
	Name          string     `valid:"ascii" json:"name" renderFor:"ALL"`
	Surname       string     `valid:"ascii" json:"surname" renderFor:"ALL"`
	Password      string     `valid:"ascii" json:"password" renderFor:"NOBODY"`
	Role          int        `valid:"optional" json:"role" renderFor:"ALL"`
	RoleString    string     `sql:"-"`
	RoleStringRus string     `valid:"optional" json:"role_string" renderFor:"ALL" sql:"-" render_key:"role_string"`
	Avatar        string     `json:"avatar" renderFor:"ALL"`

	CreatedAt time.Time `valid:"-" json:"created_at" renderFor:"admin"`
	UpdatedAt time.Time `valid:"-" json:"updated_at" renderFor:"admin"`
	DeletedAt time.Time `valid:"-" json:"deleted_at" renderFor:"admin" pg:",soft_delete"`
}

type Users []User

type Role struct {
	Id    int    `json:"id" renderFor:"ALL"`
	Title string `json:"title" renderFor:"ALL"`
}

func (User) TableName() string {
	return "users"
}

func (b *User) BeforeCreate(db orm.DB) error {
	if b.Id == nil {
		_uuid, _ := uuid.NewV4()
		b.Id = &_uuid
	}

	return nil
}

func (u *User) SetUUID() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (u *User) GetMapKey() string {
	return "User"
}

func (u *Users) GetMapKey() string {
	return "Users"
}

func (d *User) Create() error {
	return database.GetConnection().Con.Insert(d)
}

func (u *User) CreateDefaultAdmin() {
	id, _ := uuid.NewV4()
	u.Id = &id
	u.Username = "admin"
	u.Name = "admin"
	u.Surname = "admin"
	u.Password = "admin"
	u.CryptPassword()
	u.SetDefault()
	u.SetRole(USER_ROLES[1])
}

func (u *User) CryptPassword() {
	u.Password, _ = helpers.HashPassword(u.Password)
}

func (u *User) CheckHashPassword(pass string) bool {
	return helpers.CheckPasswordHash(pass, u.Password)
}

func (u *User) SetDefault() {
	u.SetRole(USER_ROLES[7])
	u.CreatedAt = time.Now()
	u.Avatar = "assets/images/avatars/profile.jpg"
}

func (u *User) GetCacheKey() string {
	return CACHEKEY + u.Id.String()
}

func (u User) IsAdmin() bool {
	return u.Role == 1
}

func (u User) IsLeadership() bool {
	return u.Role == 2
}

func (u User) IsCaller() bool {
	return u.Role == 3
}

func (u User) IsPM() bool {
	return u.Role == 4
}

func (u User) IsManager() bool {
	return u.Role == 5
}

func (u User) IsCopywriter() bool {
	return u.Role == 6
}

func (u User) IsGuest() bool {
	return u.Role == 7
}

func (u *User) SetRole(role string) {
	indexOfRole := helpers.IndexOf(role, USER_ROLES)
	if indexOfRole != -1 {
		u.Role = indexOfRole
	}
}

func (u User) GetRole() int {
	return u.Role
}

func (u User) GetRoleString() string {
	return USER_ROLES[u.Role]
}

func (u User) GetRoleStringRus() string {
	return USER_ROLES_RUS[u.Role]
}

func (u User) GetCacheExpire() time.Duration {
	return time.Hour * 240
}

func (u User) Can(roles []string) bool {
	if len(roles) == 0 {
		return true
	}

	if roles[0] == "ALL" {
		return true
	}

	b, _ := helpers.InArray(u.GetRoleString(), roles)

	return b
}

func (u User) Render(role string) map[string]interface{} {
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

func (us Users) Render(role string) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < len(us); i++ {
		data = append(data, us[i].Render(role))
	}

	return data
}
