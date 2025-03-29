package models

import (
	"context"
	"time"
)

type Event struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type EventRepository interface {
	CreateOne(ctx context.Context, event *Event) (*Event, error)
	GetOne(ctx context.Context, eventId int) (*Event, error)
	GetMany(ctx context.Context) ([]*Event, error)
	UpdateOne(ctx context.Context, eventId int, updateData map[string]interface{}) (*Event, error)
	DeleteOne(ctx context.Context, eventId int) error
}
