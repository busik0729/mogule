package ws

import (
	"log"

	uuid "github.com/satori/go.uuid"

	"../models/eventPack"
)

type Messenger interface {
	Init()
	Delete()
}

type Broadcast struct {
	Ch chan Message
	Cl map[string]bool
}

type Personal struct {
	Ch chan PersonalMessage
}

type Channels struct {
	BroadcastCh Broadcast
	PersonalCh  Personal
}

var ch = Channels{
	Broadcast{
		make(chan Message),
		make(map[string]bool),
	},
	Personal{
		make(chan PersonalMessage),
	},
}

func Init() {
	log.Println("Channels init")
	for {
		// Grab the next message from the broadcast channel
		msg := <-ch.BroadcastCh.Ch
		log.Println("broadcast: ", msg)
		// msg.Parse()
		// Send it out to every client that is currently connected
		for clientId := range ch.BroadcastCh.Cl {
			client := GetClients().FindById(clientId)
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				er := client.Close()
				if er != nil {
					log.Printf("error: %v", er)
				}
				ch.BroadcastCh.Delete(clientId)
			}
		}
	}
}

func InitPersonalChannel() {
	log.Println("Personal channel is init......")
	for {
		// Grab the next message from the broadcast channel
		msg := <-ch.PersonalCh.Ch
		log.Println("personal: ", msg)
		msg.Parse()

		SendMessage(msg.ClientId, msg.Event, msg.Data)
	}
}

func (m PersonalMessage) Parse() {
	if m.Event == "event-read" {
		mp := m.Data.(map[string]interface{})
		log.Println(mp["id"])
		id, _ := uuid.FromString(mp["id"].(string))
		event, errGetEvent := eventPack.GetById(&id)
		if errGetEvent != nil {
			log.Println(errGetEvent)
			return
		}

		event.SetStatus("readed")
		_, errEventUpdate := eventPack.Update(&event)
		if errEventUpdate != nil {
			log.Println(errEventUpdate)
		}
	}
}

func (b Broadcast) Delete(id string) {
	delete(b.Cl, id)
}

func (b Broadcast) Append(id string) {
	b.Cl[id] = true
}

func GetChannels() Channels {
	return ch
}

func GetPersonalChannel() Personal {
	return ch.PersonalCh
}
