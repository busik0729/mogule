package ws

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type PersonalMessage struct {
	Event    string      `json:"event"`
	Data     interface{} `json:"data"`
	ClientId string
}
