package schedulers

import (
	"log"
	"sync"
	"time"
	"to-do-list-golang/config"
	"to-do-list-golang/enums"
	"to-do-list-golang/models"
	"to-do-list-golang/utils"
)

func StartDueNotifier() {
	ticker := time.NewTicker(1 * time.Minute) // Run every minute

	go func() {
		for range ticker.C {
			var wg sync.WaitGroup

			for reminderName, minutes := range enums.TodoReminders {
				wg.Add(1)
				go func(name string, mins int) {
					defer wg.Done()
					checkAndSendNotifications(name, mins)
				}(reminderName, int(minutes.Minutes()))
			}

			wg.Wait()
		}
	}()
}

func checkAndSendNotifications(notificationType string, minutes int) {
	var todos []models.Todo

	now := time.Now()
	query := `
		status <> 'done'
		AND due > ?
		AND TIMESTAMPDIFF(MINUTE, ?, due) BETWEEN 0 AND ?
		AND (last_notified_type IS NULL OR last_notified_type <> ?)
	`

	if err := config.DB.Where(query, now, now, minutes, notificationType).Preload("Assignees").Find(&todos).Error; err != nil {
		log.Println("Error fetching due todo(s):", err)
		return
	}

	if len(todos) == 0 {
		return
	}

	sem := make(chan struct{}, 5) // Limit to 5 concurrent email sends
	var wg sync.WaitGroup

	for _, todo := range todos {
		wg.Add(1)
		sem <- struct{}{} // take a slot

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			tx := config.DB.Begin()
			if err := utils.SendNotificationEmail(notificationType, todo); err != nil {
				log.Println("Failed to send email:", err)
				tx.Rollback()
				return
			}

			if err := tx.Model(&todo).Updates(map[string]any{
				"last_notified_at":   time.Now(),
				"last_notified_type": notificationType,
			}).Error; err != nil {
				log.Println("Failed to update notification status")
				tx.Rollback()
			}
			tx.Commit()
		}()
	}

	wg.Wait()
}
