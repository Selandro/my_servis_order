package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	config "main.go/internal"
	"main.go/internal/handlers"
	"main.go/internal/interfacevivoda"
	"main.go/internal/natsstream"
	"main.go/internal/storage/cache"
	database "main.go/internal/storage/database"
	"main.go/internal/utils"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// setupLogger настраивает логгер в зависимости от окружения
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func main() {
	// Загрузка конфигурации
	cfg := config.MustLoad()

	// Настройка логгера
	log := setupLogger(cfg.Env)

	// Подключение к базе данных PostgreSQL
	db := database.Connect(cfg.Database)
	defer db.Close()

	// Инициализация кэша
	cache.InitCache()

	// Кэширование всех данных о заказах из базы данных
	allOrders, err := database.CacheAllOrdersFromDB(db)
	if err != nil {
		log.Error("Ошибка кэширования заказов из базы данных", slog.String("ошибка", err.Error()))
	}
	cache.SetCache(allOrders)

	// Подключение к NATS и JetStream
	js := natsstream.Connect(cfg.Nats)

	// Подписка на канал, где приходят JSON сообщения
	natsstream.Subscribe(js, "Json-orders", db)

	// Запуск HTTP-сервера для получения данных по id из кэша
	http.HandleFunc("/order", handlers.GetOrderFromCache)

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		ReadTimeout:  utils.ParseDuration(cfg.HTTPServer.Timeout),
		WriteTimeout: utils.ParseDuration(cfg.HTTPServer.Timeout),
		IdleTimeout:  utils.ParseDuration(cfg.HTTPServer.IdleTimeout),
	}

	log.Info("HTTP сервер запущен на", slog.String("адрес", cfg.HTTPServer.Address))

	// Запуск интерфейса вывода
	go interfacevivoda.Vivod()

	err = server.ListenAndServe()
	if err != nil {
		log.Error("Ошибка запуска сервера", slog.String("ошибка", err.Error()))
		os.Exit(1)
	}
}
