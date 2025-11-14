package core

import "time"

type Tracer struct {
	LastBlinkAt *time.Time `json:"last_blink_at" db:"last_blink_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ID          string     `json:"id" db:"id"`
	IP          string     `json:"ip" db:"ip"`
	TotalBlinks int        `json:"total_blinks" db:"total_blinks"`
}
