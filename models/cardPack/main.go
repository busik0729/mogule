package cardPack

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-pg/pg/orm"
	uuid "github.com/satori/go.uuid"

	"../../database"
	"../../helpers"
	"../activityPack"
	"../commentPack"
)

var DAEMON_STATUS = map[int]string{1: "np", 2: "parse"}

/**
Model struct
*/
type Card struct {
	tableName         struct{}                `sql:"card"`
	Id                *uuid.UUID              `json:"id" renderFor:"ALL"`
	Name              string                  `json:"name" renderFor:"ALL"`
	Description       string                  `json:"description" renderFor:"ALL"`
	IdAttachmentCover string                  `json:"id_attachment_cover" render_key:"id_attachment_cover" renderFor:"ALL"`
	Subscribed        bool                    `json:"subscribed" renderFor:"ALL"`
	Checklists        []Checklist             `json:"checklists" renderFor:"ALL"`
	CheckItems        int                     `json:"check_items" renderFor:"ALL" render_key:"check_items"`
	CheckItemsChecked int                     `json:"check_items_checked" renderFor:"ALL" render_key:"check_items_checked"`
	Due               time.Time               `json:"due" renderFor:"ALL"`
	ListId            *uuid.UUID              `json:"list_id" renderFor:"ALL" render_key:"list_id"`
	IdMembers         []string                `json:"idMembers" render_key:"idMembers" renderFor:"ALL" sql:"idmembers,array"`
	IdLabels          []string                `json:"idLabels" render_key:"idLabels" renderFor:"ALL" sql:"idlabels,array"`
	Activities        activityPack.Activities `json:"activities" render_key:"activities" renderFor:"ALL"`
	Comments          commentPack.Comments    `json:"comments" render_key:"comments" renderFor:"ALL"`
	Bid               int                     `sql:"-"`
	Position          float64                 `json:"position" renderFor:"ALL"`
	DeletedAt         time.Time               `valid:"-" json:"deleted_at" renderFor:"admin" pg:",soft_delete"`
	DaemonStatus      int                     `json:"daemon_status" renderFor:"NOBODY" render_key:"daemon_status"`
}

type Cards []Card

type Checklist struct {
	Id                string `json:"id" renderFor:"ALL"`
	Name              string `json:"name" renderFor:"ALL"`
	CheckItemsChecked int    `json:"checkItemsChecked" renderFor:"ALL"`
	CheckItems        []Item `json:"checkItems" renderFor:"ALL"`
}

type Item struct {
	Name    string `json:"name" renderFor:"ALL"`
	Checked bool   `json:"checked" renderFor:"ALL"`
}

func (Card) TableName() string {
	return "card"
}

func (b *Card) BeforeCreate(db orm.DB) error {
	if b.Id == nil {
		_uuid, _ := uuid.NewV4()
		b.Id = &_uuid
	}

	return nil
}

func (u *Card) SetUUID() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (c *Card) GetMapKey() string {
	return "Card"
}

func (c *Cards) GetMapKey() string {
	return "Cards"
}

func (u *Card) SetDefault() {
	u.SetDaemonStatus("np")
}

func (u *Card) SetDaemonStatus(status string) {
	indexOfStatus := helpers.IndexOf(status, DAEMON_STATUS)
	if indexOfStatus != -1 {
		u.DaemonStatus = indexOfStatus
	}
}

func (d *Card) Create() (orm.Result, error) {
	return database.GetConnection().Con.Model(d).Insert()
}

func (u Card) Render(role string) map[string]interface{} {
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

func (us Cards) Render(role string) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < len(us); i++ {
		data = append(data, us[i].Render(role))
	}

	return data
}
