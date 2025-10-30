package enums

import "time"

var (
	TodoStatuses   = []string{"pending", "in progress", "done"}
	TodoPriorities = []string{"low", "medium", "high"}
	TodoReminders  = map[string]time.Duration{
		"1 day":    24 * time.Hour,
		"12 hours": 12 * time.Hour,
		"1 hour":   time.Hour,
	}
)
