package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Подключение к серверу NATS JetStream
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Создание JetStream контекста
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	// Создание субъекта
	subject := "Json-orders"
	streamName := "Json-orders"
	_, err = js.AddStream(&nats.StreamConfig{
		Name:       streamName,
		Subjects:   []string{subject},
		Retention:  nats.WorkQueuePolicy, // Используем политику очереди работ (work queue)
		MaxAge:     1 * time.Hour,
		MaxMsgSize: 1 * 1024 * 1024,
		Storage:    nats.MemoryStorage,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created stream '%s' with subject '%s'\n", streamName, subject)
}
