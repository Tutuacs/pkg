package types

type MqttMessage struct {
	Topic    string
	Qos      byte
	Retained bool
	Payload  []byte
}
