package eventPack

import (
	"../../database"
	"../../helpers"
	"github.com/go-pg/pg/orm"
	"github.com/satori/go.uuid"
	"log"
	"reflect"
	"strings"
	"time"
)

var TYPE = map[int]string{1: "broadcast", 2: "personal"}
var STATUS = map[int]string{1: "new", 2: "readed", 3: "sended"}

/**
Model struct
*/
type Event struct {
	tableName struct{}    `sql:"event,alias:event"`
	Id        *uuid.UUID  `valid:"uuid,optional" json:"id" sql:"type:uuid,pk" renderFor:"ALL"`
	EventName string      `valid:"ascii" json:"event_name" renderFor:"ALL"`
	Data      interface{} `valid:"ascii" json:"data" renderFor:"ALL"`
	DeviceId  *uuid.UUID  `valid:"ascii" json:"device_id" renderFor:"ALL"`
	TypeEvent int         `valid:"ascii" json:"type_event" renderFor:"NOBODY"`
	Status    int         `valid:"ascii" json:"status" renderFor:"NOBODY"`
	Bid       int         `sql:"-"`

	CreatedAt time.Time `valid:"-" json:"created_at" renderFor:"admin"`
	UpdatedAt time.Time `valid:"-" json:"updated_at" renderFor:"admin"`
	DeletedAt time.Time `valid:"-" json:"deleted_at" renderFor:"admin" pg:",soft_delete"`
}

type Events []Event

func (Event) TableName() string {
	return "event"
}

func (b *Event) BeforeCreate(db orm.DB) error {
	if b.Id == nil {
		_uuid, _ := uuid.NewV4()
		b.Id = &_uuid
	}

	return nil
}

func (u *Event) SetUUID() {
	id, _ := uuid.NewV4()
	u.Id = &id
}

func (u *Event) GetMapKey() string {
	return "Event"
}

func (u *Events) GetMapKey() string {
	return "Events"
}

func (d *Event) Create() error {
	return database.GetConnection().Con.Insert(d)
}

func (u *Event) SetDefault() {
	u.SetEventType(TYPE[1])
	u.SetStatus(STATUS[1])
	u.CreatedAt = time.Now()
}

func (u *Event) SetEventType(typeS string) {
	indexOfType := helpers.IndexOf(typeS, TYPE)
	if indexOfType != -1 {
		log.Println("___________________________")
		log.Println(indexOfType)
		log.Println("___________________________")
		u.TypeEvent = indexOfType
	}
}

func (u *Event) SetStatus(status string) {
	indexOfStatus := helpers.IndexOf(status, STATUS)
	if indexOfStatus != -1 {
		u.Status = indexOfStatus
	}
}

func (u Event) GetTypeEvent() int {
	return u.TypeEvent
}

func GetPersonalTypeEvent() int {
	return 2
}

func (u Event) GetTypeEventString() string {
	return TYPE[u.TypeEvent]
}

func (u Event) GetStatus() int {
	return u.Status
}

func GetSendedStatus() int {
	return 3
}

func GetReadedStatus() int {
	return 2
}

func GetNewStatus() int {
	return 1
}

func (u Event) Render(role string) map[string]interface{} {
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

func (us Events) Render(role string) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < len(us); i++ {
		data = append(data, us[i].Render(role))
	}

	return data
}
