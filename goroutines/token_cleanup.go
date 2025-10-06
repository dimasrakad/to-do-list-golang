package goroutines

import (
	"log"
	"time"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
)

func StartTokenCleanup() {
	ticker := time.NewTicker(1 * time.Hour) // Run every hour
	go func() {
		for range ticker.C {
			now := time.Now()

			if err := config.DB.Where("expires_at < ?", now).Delete(&models.RefreshToken{}).Error; err != nil {
				log.Println("Error cleaning refresh tokens:", err)
			}

			if err := config.DB.Where("expires_at < ?", now).Delete(&models.RevokedToken{}).Error; err != nil {
				log.Println("Error cleaning revoked tokens:", err)
			}
		}
	}()
}
