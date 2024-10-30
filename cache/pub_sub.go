package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"

	"github.com/Tutuacs/pkg/config"
	"github.com/Tutuacs/pkg/logs"
	"github.com/Tutuacs/pkg/types"
)

type RedisPubSub struct {
	client      *redis.Client
	subscribers map[string]func(ctx context.Context, msg types.RedisMessage)
}

var pubsub *RedisPubSub

func init() {
	pubsub = nil
}

// return always the same instance of RedisPubSub
func UseRedisPubSub() (*RedisPubSub, error) {
	if pubsub == nil {

		conf := config.GetRedis()

		client, err := NewRedisClient(conf.Addr)
		if err != nil {
			return nil, err
		}

		pubsub = &RedisPubSub{
			client:      client.conn,
			subscribers: make(map[string]func(ctx context.Context, msg types.RedisMessage)),
		}
	}

	logs.MessageLog("Initializing Redis PubSub...")

	return pubsub, nil
}

// Subscribe the client on a topic and respective function handler
func (r *RedisPubSub) Subscribe(topic string, handler func(ctx context.Context, msg types.RedisMessage)) error {
	r.subscribers[topic] = handler
	sub := r.client.Subscribe(context.Background(), topic)
	_, err := sub.Receive(context.Background())
	if err != nil {
		return fmt.Errorf("error subscribing to topic %s: %w", topic, err)
	}
	return nil
}

// Publish a message to a topic
func (r *RedisPubSub) Publish(ctx context.Context, msg types.RedisMessage) error {
	messageBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	return r.client.Publish(ctx, msg.Topic, messageBytes).Err()
}

// Configure on the cmd/pub-sub service
func (r *RedisPubSub) Listen() {
	for topic := range r.subscribers {
		go func(topic string) {
			sub := r.client.Subscribe(context.Background(), topic)
			defer sub.Close()

			for {
				msg, err := sub.ReceiveMessage(ctx)
				if err != nil {
					log.Printf("Error receiving message from topic %s: %v", topic, err)
					break
				}

				// Decodifica a mensagem e invoca o handler associado ao tópico
				var message types.RedisMessage
				if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
					log.Printf("Error unmarshalling message from topic %s: %v", topic, err)
					continue
				}

				// Executa o handler em uma goroutine
				if handler, exists := r.subscribers[topic]; exists {
					handler(ctx, message)
				}
			}
		}(topic)
	}
}

// Exemplo de uso: criando funções para lidar com mensagens de diferentes tópicos
func HandleHello(ctx context.Context, msg types.RedisMessage) {
	fmt.Println("Handling /hello with data:", msg)
	// Processar a mensagem para o tópico "/hello"
}

func HandleUnknown(ctx context.Context, msg types.RedisMessage) {
	fmt.Println("Handling /unknown with data:", msg)
	fmt.Println("Oia onde chegou")
	// Processar a mensagem para o tópico "/unknown"
}
