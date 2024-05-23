package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NMID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SMID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OOFShard          string   `json:"oof_shard"`
}

func main() {
	nc, err := nats.Connect("js://localhost:4222")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	// Создание JetStream контекста
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatalf("Error creating JetStream context: %v", err)
	}
	var wg sync.WaitGroup
	start := time.Now()
	// Отправка 10 разных ордеров
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			order := createOrder(i)
			jsonData, err := json.Marshal(order)
			if err != nil {
				log.Fatalf("Error marshaling order to JSON: %v", err)
			}

			// Отправка сообщения в NATS Streaming
			_, err = js.Publish("Json-orders", jsonData)
			if err != nil {
				log.Fatalf("Error publishing message: %v", err)
			}

			log.Printf("Sent message: %s", order.OrderUID)
		}
	}()
	wg.Add(1)
	go func() {

		defer wg.Done()
		for k := 1000; k < 2000; k++ {
			order := createOrder(k)
			jsonData, err := json.Marshal(order)
			if err != nil {
				log.Fatalf("Error marshaling order to JSON: %v", err)
			}

			// Отправка сообщения в NATS Streaming
			_, err = js.Publish("Json-orders", jsonData)
			if err != nil {
				log.Fatalf("Error publishing message: %v", err)
			}

			log.Printf("Sent message: %s", order.OrderUID)
		}
	}()
	wg.Add(1)
	go func() {

		defer wg.Done()
		for j := 2000; j < 3000; j++ {
			order := createOrder(j)
			jsonData, err := json.Marshal(order)
			if err != nil {
				log.Fatalf("Error marshaling order to JSON: %v", err)
			}

			// Отправка сообщения в NATS Streaming
			_, err = js.Publish("Json-orders", jsonData)
			if err != nil {
				log.Fatalf("Error publishing message: %v", err)
			}

			log.Printf("Sent message: %s", order.OrderUID)
		}
	}()
	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Sent all orders in %s", elapsed)
}

func createOrder(orderNum int) Order {
	return Order{
		OrderUID:    "order_" + fmt.Sprint(orderNum),
		TrackNumber: "track_" + fmt.Sprint(orderNum),
		Entry:       "entry_" + fmt.Sprint(orderNum),
		Delivery:    Delivery{Name: "Name_" + fmt.Sprint(orderNum), Phone: "Phone_" + fmt.Sprint(orderNum), Zip: "Zip_" + fmt.Sprint(orderNum), City: "City_" + fmt.Sprint(orderNum), Address: "Address_" + fmt.Sprint(orderNum), Region: "Region_" + fmt.Sprint(orderNum), Email: "Email_" + fmt.Sprint(orderNum)},
		Payment:     Payment{Transaction: "Transaction_" + fmt.Sprint(orderNum), RequestID: "RequestID_" + fmt.Sprint(orderNum), Currency: "Currency_" + fmt.Sprint(orderNum), Provider: "Provider_" + fmt.Sprint(orderNum), Amount: orderNum, PaymentDT: orderNum, Bank: "Bank_" + fmt.Sprint(orderNum), DeliveryCost: orderNum, GoodsTotal: orderNum, CustomFee: orderNum},
		Items: []Item{{ChrtID: orderNum, TrackNumber: "track_" + fmt.Sprint(orderNum), Price: orderNum, RID: "RID_" + fmt.Sprint(orderNum), Name: "Name_" + fmt.Sprint(orderNum), Sale: orderNum, Size: "Size_" + fmt.Sprint(orderNum), TotalPrice: orderNum, NMID: orderNum, Brand: "Brand_" + fmt.Sprint(orderNum), Status: orderNum},
			{ChrtID: orderNum, TrackNumber: "track_" + fmt.Sprint(orderNum), Price: orderNum, RID: "RID_" + fmt.Sprint(orderNum), Name: "Name_" + fmt.Sprint(orderNum), Sale: orderNum, Size: "Size_" + fmt.Sprint(orderNum), TotalPrice: orderNum, NMID: orderNum, Brand: "Brand_" + fmt.Sprint(orderNum), Status: orderNum}},
		Locale:            "Locale_" + fmt.Sprint(orderNum),
		InternalSignature: "InternalSignature_" + fmt.Sprint(orderNum),
		CustomerID:        "CustomerID_" + fmt.Sprint(orderNum),
		DeliveryService:   "DeliveryService_" + fmt.Sprint(orderNum),
		Shardkey:          "Shardkey_" + fmt.Sprint(orderNum),
		SMID:              orderNum,
		DateCreated:       "2021-11-26T06:22:19Z",
		OOFShard:          "OOFShard_" + fmt.Sprint(orderNum),
	}
}
