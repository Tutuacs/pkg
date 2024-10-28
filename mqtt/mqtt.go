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
		opts.SetClientID("gbase")
		opts.SetUsername("gbase_user")
		opts.SetDefaultPublishHandler(messagePubHandler)
		opts.OnConnect = connectHandler
		opts.OnConnectionLost = connectLostHandler
		opts.ConnectRetry = true

		conn := mqtt.NewClient(opts)

		if token := conn.Connect(); token.Wait() && token.Error() != nil {
			logs.ErrorLog(fmt.Sprintf("Error connecting to mqtt broker: %v", token.Error()))
		}

		client = &Mqtt{
			Addr: conf.Addr,
			conn: conn,
		}
	}

	return client, nil
}

var HandleHello mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received Hello message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var HandleUnknown mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received Unknown message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func Publish(timeout time.Duration, msg types.MqttMessage) (err error) {

	err = fmt.Errorf("mqtt client timeout")

	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if connected := client.conn.IsConnected(); !connected {
		client.conn.Connect()
	}

	token := client.conn.Publish(msg.Topic(), 0, false, msg.Payload())

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
