package interfacevivoda

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"main.go/internal/storage/cache"
	model "main.go/orders_model"
)

// Структура Delivery представляет информацию о доставке.
type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

// Структура Payment представляет информацию о платеже.
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

// Структура Item представляет информацию о товаре.
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

// Структура Order представляет информацию о заказе.
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

// displayOrder выводит подробности о заказе.
func displayOrder(order model.Order) {
	fmt.Println("Order ID:", order.OrderUID)
	fmt.Println("Track Number:", order.TrackNumber)
	fmt.Println("Entry:", order.Entry)
	fmt.Println("Locale:", order.Locale)
	fmt.Println("Internal Signature:", order.InternalSignature)
	fmt.Println("Customer ID:", order.CustomerID)
	fmt.Println("Delivery Service:", order.DeliveryService)
	fmt.Println("Shardkey:", order.Shardkey)
	fmt.Println("SMID:", order.SMID)
	fmt.Println("Date Created:", order.DateCreated)
	fmt.Println("OOF Shard:", order.OOFShard)

	fmt.Println("Delivery:")
	fmt.Println("  Name:", order.Delivery.Name)
	fmt.Println("  Phone:", order.Delivery.Phone)
	fmt.Println("  Zip:", order.Delivery.Zip)
	fmt.Println("  City:", order.Delivery.City)
	fmt.Println("  Address:", order.Delivery.Address)
	fmt.Println("  Region:", order.Delivery.Region)
	fmt.Println("  Email:", order.Delivery.Email)

	fmt.Println("Payment:")
	fmt.Println("  Transaction:", order.Payment.Transaction)
	fmt.Println("  Request ID:", order.Payment.RequestID)
	fmt.Println("  Currency:", order.Payment.Currency)
	fmt.Println("  Provider:", order.Payment.Provider)
	fmt.Println("  Amount:", order.Payment.Amount)
	fmt.Println("  Payment Date:", order.Payment.PaymentDT)
	fmt.Println("  Bank:", order.Payment.Bank)
	fmt.Println("  Delivery Cost:", order.Payment.DeliveryCost)
	fmt.Println("  Goods Total:", order.Payment.GoodsTotal)
	fmt.Println("  Custom Fee:", order.Payment.CustomFee)

	fmt.Println("Items:")
	for _, item := range order.Items {
		fmt.Println("  Chart ID:", item.ChrtID)
		fmt.Println("  Track Number:", item.TrackNumber)
		fmt.Println("  Price:", item.Price)
		fmt.Println("  RID:", item.RID)
		fmt.Println("  Name:", item.Name)
		fmt.Println("  Sale:", item.Sale)
		fmt.Println("  Size:", item.Size)
		fmt.Println("  Total Price:", item.TotalPrice)
		fmt.Println("  NMID:", item.NMID)
		fmt.Println("  Brand:", item.Brand)
		fmt.Println("  Status:", item.Status)
	}
}

// Vivod запускает интерфейс вывода данных о заказах.
func Vivod() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Введите ID заказа для отображения его подробностей (или введите 'exit', чтобы выйти):")

	for {
		fmt.Print("Order ID: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			fmt.Println("Выход...")
			break
		}

		order, found := cache.OrderCache[input]
		if !found {
			fmt.Println("Заказ с ID", input, "не найден.")
			continue
		}

		displayOrder(order)

		fmt.Println()
	}
}
