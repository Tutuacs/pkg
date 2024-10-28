package types

import mqtt "github.com/eclipse/paho.mqtt.golang"

type Message struct {
	Topic string      `json:"topic"`
	Data  interface{} `json:"data"`
}

type MqttMessage struct {
	mqtt.Message
}
