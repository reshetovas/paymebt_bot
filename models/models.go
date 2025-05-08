package models

import "time"

type Payment struct {
	ID       uint
	UserID   int64
	Category string
	Amount   float64
	Date     time.Time
	CretedAt time.Time
}

// enam for State
type State string

const (
	Idle             State = "idle"
	AwaitingCategory State = "awaiting_category"
	AwaitingAmount   State = "awaiting_amount"
	AwaitingDate     State = "awaiting_date"
)
