package natsstream

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	config "main.go/internal"
	"main.go/internal/storage/cache"
	database "main.go/internal/storage/database"
	"main.go/orders_model"
)

// Stream представляет поток сообщений от NATS.
type Stream struct {
	OrdersChannel chan *orders_model.Order
}

// Subscribe подписывается на поток сообщений и обрабатывает их.
func (s *Stream) Subscribe(db *sql.DB) {
	for order := range s.OrdersChannel {
		// Обработка сообщения - вставка заказа в базу данных и кэширование
		if err := database.InsertOrderToDB(*order, db); err != nil {
			fmt.Println("Ошибка при вставке заказа в базу данных:", err)
			continue
		}
		cache.CacheOrder(*order)
		fmt.Println("Заказ успешно добавлен:", order.OrderUID)
		// Отправка подтверждения обработки сообщения
		// msg.Ack() - в случае использования реального NATS
	}
}

// Connect устанавливает соединение с NATS и JetStream.
func Connect(cfg config.NatsConfig) nats.JetStreamContext {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS: %v", err)
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Ошибка при подключении к JetStream: %v", err)
	}
	return js
}

// Subscribe подписывается на указанный канал и обрабатывает полученные сообщения.
func Subscribe(js nats.JetStreamContext, subject string, db *sql.DB) {
	ackWait := 30 * time.Second
	_, err := js.Subscribe(subject, func(msg *nats.Msg) {
		var order orders_model.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			fmt.Println("Ошибка декодирования JSON:", err)
			return
		}
		available, err := database.OrderExists(order.OrderUID, db)
		if err != nil {
			fmt.Println("Ошибка при проверке существования заказа:", err)
			return
		}

		if available {
			fmt.Println("Заказ с таким же ID уже существует")
		} else {

			err := database.InsertOrderToDB(order, db)
			if err != nil {
				fmt.Println("Ошибка при вставке заказа в базу данных:", err)
				return
			}
			cache.CacheOrder(order)
			fmt.Println("Заказ успешно добавлен:", order.OrderUID)
			msg.Ack()
		}
	}, nats.AckWait(ackWait), nats.ManualAck())

	if err != nil {
		log.Fatalf("Ошибка при подписке на JetStream: %v", err)
	}
}
