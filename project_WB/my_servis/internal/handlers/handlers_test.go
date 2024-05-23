package handlers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
	config "main.go/internal"
	"main.go/internal/handlers"
	"main.go/internal/storage/cache"
	database "main.go/internal/storage/database"
)

func TestGetOrderFromCache(t *testing.T) {

	cache.InitCache()
	// Предварительно загрузить данные в кэш
	preloadCache()

	// Создать запрос
	req, err := http.NewRequest("GET", "/order?id=order_1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создать запись для записи ответа
	rr := httptest.NewRecorder()

	// Создать обработчик и обработать запрос
	handler := http.HandlerFunc(handlers.GetOrderFromCache)
	handler.ServeHTTP(rr, req)

	// Проверить код ответа
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Проверить тело ответа
	expected := `{"order_uid":"order_1","track_number":"track_1","entry":"entry_1","delivery":{"name":"Name_1","phone":"Phone_1","zip":"Zip_1","city":"City_1","address":"Address_1","region":"Region_1","email":"Email_1"},"payment":{"transaction":"Transaction_1","request_id":"RequestID_1","currency":"Currency_1","provider":"Provider_1","amount":1,"payment_dt":1,"bank":"Bank_1","delivery_cost":1,"goods_total":1,"custom_fee":1},"items":[{"chrt_id":1,"track_number":"track_1","price":1,"rid":"RID_1","name":"Name_1","sale":1,"size":"Size_1","total_price":1,"nm_id":1,"brand":"Brand_1","status":1},{"chrt_id":1,"track_number":"track_1","price":1,"rid":"RID_1","name":"Name_1","sale":1,"size":"Size_1","total_price":1,"nm_id":1,"brand":"Brand_1","status":1}],"locale":"Locale_1","internal_signature":"InternalSignature_1","customer_id":"CustomerID_1","delivery_service":"DeliveryService_1","shardkey":"Shardkey_1","sm_id":1,"date_created":"2021-11-26T06:22:19Z","oof_shard":"OOFShard_1"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func preloadCache() {
	cfg := config.MustLoad()
	// Создать подключение к базе данных
	db := database.Connect(cfg.Database)
	defer db.Close()

	// Загрузить все заказы из базы данных в кэш
	allOrders, err := database.CacheAllOrdersFromDB(db)
	if err != nil {
		log.Printf("Error caching orders from database: %v", err)
	}
	cache.SetCache(allOrders)
}
