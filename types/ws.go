package types

type WsMessage struct {
	Topic string      `json:"topic"`
	Data  interface{} `json:"data"`
}
