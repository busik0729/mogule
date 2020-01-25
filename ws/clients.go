package ws

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

	"log"
)

type Clients struct {
	all map[string]*websocket.Conn
}

var cls = Clients{make(map[string]*websocket.Conn)}

func (cls Clients) AppendClient(client *websocket.Conn) string {
	id := generateUuid()
	cls.all[id] = client

	return id
}

func (cls *Clients) DeleteClient(uuid string) *Clients {
	delete(cls.all, uuid)

	return cls
}

func (cls Clients) FindById(uuid string) *websocket.Conn {
	return cls.all[uuid]
}

func GetClients() Clients {
	return cls
}

func SendMessage(wsId string, event string, message interface{}) {
	wsClient := GetClients().FindById(wsId)
	msg := Message{event, message}
	err := wsClient.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func SendBroadcast(event string, message interface{}) {
	msg := Message{event, message}
	channels := GetChannels().BroadcastCh
	channels.Ch <- msg
}

func generateUuid() string {
	id, err := uuid.NewV4()
	// out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	// return string(out[:])
	return id.String()
}
