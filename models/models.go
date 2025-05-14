package models

import (
	"time"
)

type Payment struct {
	ID        uint
	UserID    int64
	Category  string
	Amount    float64
	Date      time.Time
	CreatedAt time.Time
}

type CategoryReport struct {
	Category string
	Amount   float64
}

// enam for State
type State string

const (
	Idle                 State = "idle"
	AwaitingCategory     State = "awaiting_category"
	AwaitingAmount       State = "awaiting_amount"
	AwaitingDate         State = "awaiting_date"
	AwaitingCustomReport State = "awaiting_custom_report"
	AwaitingExport       State = "awaiting_export"
)
