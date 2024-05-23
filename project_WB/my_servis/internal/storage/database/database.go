package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	config "main.go/internal"
	model "main.go/orders_model"
)

// Connect устанавливает соединение с базой данных и возвращает объект DB.
func Connect(cfg config.DatabaseConfig) *sql.DB {
	// Формируем строку подключения к базе данных
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	// Создаем необходимые таблицы, если они еще не существуют
	createTables(db)
	return db
}

// OrderExists проверяет, существует ли заказ в базе данных.
func OrderExists(orderUID string, db *sql.DB) (bool, error) {
	var exists bool
	// Выполняем запрос к базе данных, чтобы узнать, существует ли заказ с указанным orderUID
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)", orderUID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// InsertOrderToDB вставляет заказ в базу данных
func InsertOrderToDB(order model.Order, db *sql.DB) error {
	orderID := order.OrderUID

	// Вставляем информацию о заказе
	if err := insertOrder(order, db); err != nil {
		return fmt.Errorf("ошибка вставки заказа: %v", err)
	}

	// Вставляем информацию о доставке
	if err := insertDelivery(order.Delivery, orderID, db); err != nil {
		return fmt.Errorf("ошибка вставки доставки: %v", err)
	}

	// Вставляем информацию о платеже
	if err := insertPayment(order.Payment, orderID, db); err != nil {
		return fmt.Errorf("ошибка вставки платежа: %v", err)
	}

	// Вставляем информацию о товарах
	if err := insertItems(order.Items, orderID, db); err != nil {
		return fmt.Errorf("ошибка вставки товаров: %v", err)
	}

	return nil
}

// insertOrder вставляет информацию о заказе в базу данных.
func insertOrder(order model.Order, db *sql.DB) error {
	query := `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := db.Exec(query, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SMID, order.DateCreated, order.OOFShard)
	return err
}

// insertDelivery вставляет информацию о доставке в базу данных.
func insertDelivery(delivery model.Delivery, orderID string, db *sql.DB) error {
	query := `
		INSERT INTO deliveries (name, phone, zip, city, address, region, email, order_uid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(query, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email, orderID)
	return err
}

// insertPayment вставляет информацию о платеже в базу данных.
func insertPayment(payment model.Payment, orderID string, db *sql.DB) error {
	query := `
		INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, order_uid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := db.Exec(query, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDT, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee, orderID)
	return err
}

// insertItems вставляет информацию о товарах в базу данных.
func insertItems(items []model.Item, orderID string, db *sql.DB) error {
	query := `
		INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_uid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	for _, item := range items {
		_, err := db.Exec(query, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status, orderID)
		if err != nil {
			return err
		}
	}
	return nil
}

// CacheAllOrdersFromDB кэширует все заказы из базы данных
func CacheAllOrdersFromDB(db *sql.DB) (map[string]model.Order, error) {
	orderCache := make(map[string]model.Order)
	itemsMap, err := fetchItemsMap(db)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
		       d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		       p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		JOIN deliveries d ON o.order_uid = d.order_uid
		JOIN payments p ON o.order_uid = p.order_uid`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching orders from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SMID, &order.DateCreated, &order.OOFShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
		if err != nil {
			return nil, fmt.Errorf("error scanning order row: %v", err)
		}
		order.Items = itemsMap[order.OrderUID]
		orderCache[order.OrderUID] = order
	}

	return orderCache, nil
}

// fetchItemsMap извлекает информацию о товарах из базы данных и возвращает ее в виде map
func fetchItemsMap(db *sql.DB) (map[string][]model.Item, error) {
	itemsMap := make(map[string][]model.Item)
	query := `
		SELECT order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching items from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		var orderUID string
		if err := rows.Scan(&orderUID, &item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NMID, &item.Brand, &item.Status); err != nil {
			return nil, fmt.Errorf("error scanning item row: %v", err)
		}
		itemsMap[orderUID] = append(itemsMap[orderUID], item)
	}

	return itemsMap, nil
}

// createTables создает необходимые таблицы в базе данных, если они еще не существуют.
func createTables(db *sql.DB) {
	createOrdersTable := `
	CREATE TABLE IF NOT EXISTS orders (
		order_uid VARCHAR(255) PRIMARY KEY,
		track_number VARCHAR(255),
		entry VARCHAR(255),
		locale VARCHAR(255),
		internal_signature VARCHAR(255),
		customer_id VARCHAR(255),
		delivery_service VARCHAR(255),
		shardkey VARCHAR(255),
		sm_id INT,
		date_created TIMESTAMP,
		oof_shard VARCHAR(255)
	);`

	createDeliveriesTable := `
	CREATE TABLE IF NOT EXISTS deliveries (
		order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid),
		name VARCHAR(255),
		phone VARCHAR(255),
		zip VARCHAR(255),
		city VARCHAR(255),
		address VARCHAR(255),
		region VARCHAR(255),
		email VARCHAR(255)
	);`

	createPaymentsTable := `
	CREATE TABLE IF NOT EXISTS payments (
		order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid),
		transaction VARCHAR(255),
		request_id VARCHAR(255),
		currency VARCHAR(255),
		provider VARCHAR(255),
		amount INT,
		payment_dt BIGINT,
		bank VARCHAR(255),
		delivery_cost INT,
		goods_total INT,
		custom_fee INT
	);`

	createItemsTable := `
	CREATE TABLE IF NOT EXISTS items (
		order_uid VARCHAR(255) REFERENCES orders(order_uid),
		chrt_id INT,
		track_number VARCHAR(255),
		price INT,
		rid VARCHAR(255),
		name VARCHAR(255),
		sale INT,
		size VARCHAR(255),
		total_price INT,
		nm_id INT,
		brand VARCHAR(255),
		status INT
	);`

	_, err := db.Exec(createOrdersTable)
	if err != nil {
		log.Fatalf("Error creating orders table: %v", err)
	}

	_, err = db.Exec(createDeliveriesTable)
	if err != nil {
		log.Fatalf("Error creating deliveries table: %v", err)
	}

	_, err = db.Exec(createPaymentsTable)
	if err != nil {
		log.Fatalf("Error creating payments table: %v", err)
	}

	_, err = db.Exec(createItemsTable)
	if err != nil {
		log.Fatalf("Error creating items table: %v", err)
	}
}
