package cache

import (
	"sync"

	model "main.go/orders_model"
)

var (
	OrderCache     map[string]model.Order
	orderCacheLock sync.RWMutex
)

// InitCache инициализирует кэш заказов.
func InitCache() {
	OrderCache = make(map[string]model.Order)
}

// CacheOrder добавляет заказ в кэш.
func CacheOrder(order model.Order) {
	orderCacheLock.Lock()
	defer orderCacheLock.Unlock()
	OrderCache[order.OrderUID] = order
}

// GetOrderFromCache получает заказ из кэша по его идентификатору.
func GetOrderFromCache(orderUID string) (model.Order, bool) {
	orderCacheLock.RLock()
	defer orderCacheLock.RUnlock()
	order, exists := OrderCache[orderUID]
	return order, exists
}

// SetCache кэширует все заказы.
func SetCache(allOrders map[string]model.Order) {
	orderCacheLock.Lock()
	defer orderCacheLock.Unlock()
	for k, v := range allOrders {
		OrderCache[k] = v
	}
}
