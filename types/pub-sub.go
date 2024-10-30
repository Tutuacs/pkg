package types

type RedisMessage struct {
	Topic string      `json:"topic"`
	Data  interface{} `json:"data"`
}
