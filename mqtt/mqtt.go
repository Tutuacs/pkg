package mqtt

import (
	"context"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/Tutuacs/pkg/config"
	"github.com/Tutuacs/pkg/logs"
	"github.com/Tutuacs/pkg/types"
)

type Mqtt struct {
	Addr string
	conn mqtt.Client
}

var client *Mqtt

func init() {
	client = nil
}

func UseMqtt() (*Mqtt, error) {

	if client == nil {

		conf := config.GetMqtt()

		opts := mqtt.NewClientOptions()
		opts.AddBroker(conf.Addr)
		opts.SetClientID("go_mqtt_client")
		opts.SetUsername("gbase_username")
		opts.SetDefaultPublishHandler(messagePubHandler)
		opts.OnConnect = connectHandler
		opts.OnConnectionLost = connectLostHandler
		opts.ConnectRetry = true
		mqttClient := mqtt.NewClient(opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		client = &Mqtt{
			Addr: conf.Addr,
			conn: mqttClient,
		}
	}

	return client, nil
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	logs.MessageLog(fmt.Sprintf("Received from topic: %s message: %s\n", msg.Topic(), msg.Payload()))
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	logs.OkLog("MQTT client connected...")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	logs.WarnLog("MQTT client lost connection...")
}

func (mqtt *Mqtt) Publish(timeout time.Duration, msg types.MqttMessage) (err error) {

	err = fmt.Errorf("mqtt client timeout")

	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if connected := client.conn.IsConnected(); !connected {
		client.conn.Connect()
	}

	token := client.conn.Publish(msg.Topic, 0, false, msg.Payload)

	token.Wait()
	err = token.Error()
	return
}

func (mqtt *Mqtt) Subscribe(topic string, handler mqtt.MessageHandler) (err error) {

	if connected := client.conn.IsConnected(); !connected {
		mqtt.conn.Connect()
	}

	token := client.conn.Subscribe(topic, 0, handler)

	token.Wait()
	err = token.Error()
	return
}

// TODO: Create your handlers

var HandleHello mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received from topic: %s message: %s\n", msg.Topic(), msg.Payload())
	fmt.Println("Look at Hello")
}

var HandleUnknown mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received from topic: %s message: %s\n", msg.Topic(), msg.Payload())
	fmt.Println("Look at unknown")
}
