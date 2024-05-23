package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config структура, содержащая настройки приложения.
type Config struct {
	Env        string           `yaml:"env"`         // Env определяет окружение приложения.
	Database   DatabaseConfig   `yaml:"database"`    // Database содержит настройки базы данных.
	Nats       NatsConfig       `yaml:"nats"`        // Nats содержит настройки NATS.
	HTTPServer HTTPServerConfig `yaml:"http_server"` // HTTPServer содержит настройки HTTP-сервера.
}

// DatabaseConfig содержит настройки подключения к базе данных.
type DatabaseConfig struct {
	Host     string `yaml:"host"`     // Host адрес хоста базы данных.
	Port     int    `yaml:"port"`     // Port порт базы данных.
	User     string `yaml:"user"`     // User имя пользователя базы данных.
	Password string `yaml:"password"` // Password пароль пользователя базы данных.
	DBName   string `yaml:"dbname"`   // DBName имя базы данных.
	SSLMode  string `yaml:"sslmode"`  // SSLMode режим SSL подключения.
}

// NatsConfig содержит настройки подключения к NATS.
type NatsConfig struct {
	ClusterID string `yaml:"cluster_id"` // ClusterID идентификатор кластера NATS.
	ClientID  string `yaml:"client_id"`  // ClientID идентификатор клиента NATS.
	URL       string `yaml:"url"`        // URL адрес сервера NATS.
}

// HTTPServerConfig содержит настройки HTTP-сервера.
type HTTPServerConfig struct {
	Address     string `yaml:"address"`      // Address адрес, на котором запущен HTTP-сервер.
	Timeout     string `yaml:"timeout"`      // Timeout таймаут запроса к серверу.
	IdleTimeout string `yaml:"idle_timeout"` // IdleTimeout таймаут ожидания.
}

// MustLoad загружает конфигурацию из указанного файла и завершает программу с ошибкой в случае неудачи.
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg

}
