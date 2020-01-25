package ws

import (
	"log"
	"net/http"

	"../../helpers"
	"../../models/devicePack"
	"../../models/eventPack"
	"../../structs/appCxt"
	"../../ws"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func WsMainConnect(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	wsClient, err := ws.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer wsClient.Close()

	clients := ws.GetClients()

	uuid := clients.AppendClient(wsClient)
	channels := ws.GetChannels().BroadcastCh
	channels.Append(uuid)

	device := appContext.CurrentDevice
	device.SetWsUUID(uuid)
	devicePack.Update(&device)
	sendEvents(device)

	for {
		// _, message, er := ws.ReadMessage()
		// if er != nil {
		// 	if websocket.IsUnexpectedCloseError(er, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 		log.Printf("error: %v", er)
		// 	}
		// 	break
		// }
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// log.Println(string(message))

		var msg ws.Message
		// Read in a new message as JSON and map it to a Message object
		err := wsClient.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			clients.DeleteClient(uuid)
			channels.Delete(uuid)

			device.SetNullWsUUID()
			devicePack.Update(&device)
			break
		}

		// Send the newly received message to the broadcast channel
		channels.Ch <- msg
	}
}

func sendEvents(device devicePack.Device) error {
	events, errGetEvents := eventPack.GetNewEventsByDevice(*device.Id)
	if errGetEvents != nil {
		helpers.LogToFile(helpers.Join(errGetEvents.Error(), ":sendEvents"))
		return errGetEvents
	}
	log.Println("..........................................................")
	log.Println(events)
	log.Println("..........................................................")

	for _, v := range events {
		// ws.SendMessage(device.WsId, v.EventName, v)
		ws.GetPersonalChannel().Ch <- ws.PersonalMessage{
			Event:    v.EventName,
			Data:     v,
			ClientId: device.WsId,
		}
		v.SetStatus("sended")
		eventPack.Update(&v)
	}

	return nil
}
