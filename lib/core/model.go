package core

import "time"

type Tracer struct {
	LastBlinkAt *time.Time `json:"last_blink_at" db:"last_blink_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ID          string     `json:"id" db:"id"`
	Nickname    string     `json:"nickname" db:"nickname"`
	TotalBlinks int        `json:"total_blinks" db:"total_blinks"`
}

type Blink struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ID        string    `json:"id" db:"id"`
	TracerID  string    `json:"tracer_id" db:"tracer_id"`
}
