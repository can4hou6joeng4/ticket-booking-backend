package models

import "context"

type Statistics struct {
	TotalEvents      int64 `json:"eventCount"`
	TotalTickets     int64 `json:"ticketCount"`
	ValidatedTickets int64 `json:"validationCount"`
}
type StatisticsRepository interface {
	GetCount(ctx context.Context) (*Statistics, error)
}
