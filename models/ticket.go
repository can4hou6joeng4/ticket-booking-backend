package models

import (
	"context"
	"time"
)

type Ticket struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	EventID   uint      `json:"eventId"`
	Event     Event     `json:"event" gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Entered   bool      `json:"entered"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type TicketRepository interface {
	CreateOne(ctx context.Context, ticket *Ticket) (*Ticket, error)
	GetOne(ctx context.Context, ticketId uint) (*Ticket, error)
	GetMany(ctx context.Context) ([]*Ticket, error)
	UpdateOne(ctx context.Context, ticketId uint, updateData map[string]interface{}) (*Ticket, error)
}
type ValidateTicket struct {
	TicketId uint `json:"ticketId"`
}
