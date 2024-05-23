package handlers

import (
	"encoding/json"
	"net/http"

	cache "main.go/internal/storage/cache"
)

// GetOrderFromCache обрабатывает запрос на получение данных о заказе из кэша по его идентификатору.
func GetOrderFromCache(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор заказа из параметров запроса
	orderUID := r.URL.Query().Get("id")
	if orderUID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest) // Возвращаем ошибку, если идентификатор заказа отсутствует в запросе.
		return
	}

	// Получаем заказ из кэша
	order, exists := cache.GetOrderFromCache(orderUID)
	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound) // Возвращаем ошибку, если заказ не найден в кэше.
		return
	}

	// Преобразуем данные заказа в формат JSON
	responseData, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Error marshaling response data", http.StatusInternalServerError) // Возвращаем ошибку, если возникла ошибка при преобразовании данных в JSON.
		return
	}

	// Устанавливаем заголовок Content-Type в application/json
	w.Header().Set("Content-Type", "application/json")

	// Отправляем данные заказа в ответ на запрос
	w.Write(responseData)
}
