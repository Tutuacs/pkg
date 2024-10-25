package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/Tutuacs/pkg/logs"
	"github.com/Tutuacs/pkg/types"
)

type RedisPubSub struct {
	client      *redis.Client
	subscribers map[string]func(ctx context.Context, msg types.Message)
}

var pubsub *RedisPubSub

// Inicia o RedisPubSub com o cliente Redis
func NewRedisPubSub(addr string) (*RedisPubSub, error) {
	client, err := NewRedisClient(addr)
	if err != nil {
		return nil, err
	}

	pubsub = &RedisPubSub{
		client:      client.conn,
		subscribers: make(map[string]func(ctx context.Context, msg types.Message)),
	}

	logs.MessageLog("Initializing Redis PubSub...")
	go pubsub.listen()

	return pubsub, nil
}

// Subscreve-se a um tópico com uma função de manipulação
func (r *RedisPubSub) Subscribe(topic string, handler func(ctx context.Context, msg types.Message)) {
	r.subscribers[topic] = handler
	r.client.Subscribe(context.Background(), topic)
}

// Publica uma mensagem em um tópico
func (r *RedisPubSub) Publish(ctx context.Context, topic string, msg types.Message) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("falha ao codificar mensagem: %w", err)
	}

	return r.client.Publish(ctx, topic, messageBytes).Err()
}

// Função de listener para receber mensagens de tópicos assinados
func (r *RedisPubSub) listen() {
	pubsub := r.client.Subscribe(context.Background())
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			fmt.Println("Erro ao receber mensagem:", err)
			continue
		}

		// Decodifica a mensagem e invoca a função associada ao tópico
		var message types.Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			fmt.Println("Erro ao decodificar mensagem:", err)
			continue
		}

		if handler, exists := r.subscribers[msg.Channel]; exists {
			go handler(context.Background(), message) // Chama o manipulador em uma goroutine
		}
	}
}

// Exemplo de uso: criando funções para lidar com mensagens de diferentes tópicos
func HandleHello(ctx context.Context, msg types.Message) {
	fmt.Println("Handling /hello with data:", msg)
	// Processar a mensagem para o tópico "/hello"
}

func HandleUnknown(ctx context.Context, msg types.Message) {
	fmt.Println("Handling /unknown with data:", msg)
	// Processar a mensagem para o tópico "/unknown"
}
