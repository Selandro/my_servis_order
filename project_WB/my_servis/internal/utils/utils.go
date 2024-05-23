package utils

import (
	"log"
	"time"
)

// ParseDuration преобразует строку в объект time.Duration.
func ParseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Fatalf("Invalid duration format: %v", err)
	}
	return d
}
