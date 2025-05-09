package repositories

import (
	"context"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func (r *TicketRepository) CreateOne(ctx context.Context, userId uint, ticket *models.Ticket) (*models.Ticket, error) {
	ticket.UserID = userId
	// 插入数据，ID 被回填
	res := r.db.Model(ticket).Create(ticket)
	if res.Error != nil {
		return nil, res.Error
	}
	// 此时 ticket.ID 已经是数据库中的新 ID
	return r.GetOne(ctx, userId, ticket.ID)
}

func (r TicketRepository) GetOne(ctx context.Context, userId uint, ticketId uint) (*models.Ticket, error) {
	ticket := &models.Ticket{}
	res := r.db.Model(ticket).Where("id = ?", ticketId).Where("user_id = ?", userId).Preload("Event").First(ticket)
	if res.Error != nil {
		return nil, res.Error
	}
	return ticket, nil
}

func (r TicketRepository) GetMany(ctx context.Context, userId uint) ([]*models.Ticket, error) {
	tickets := []*models.Ticket{}
	// 预加载关联的 Event 数据
	res := r.db.Model(&tickets).Where("user_id = ?", userId).Preload("Event").Order("updated_at DESC").Find(&tickets)
	if res.Error != nil {
		return nil, res.Error
	}
	return tickets, nil
}

func (r TicketRepository) UpdateOne(ctx context.Context, userId uint, ticketId uint, updateData map[string]interface{}) (*models.Ticket, error) {
	ticket := &models.Ticket{}
	res := r.db.Model(ticket).Where("id = ?", ticketId).Updates(updateData)
	if res.Error != nil {
		return nil, res.Error
	}
	return r.GetOne(ctx, userId, ticketId)
}

func NewTicketRepository(db *gorm.DB) models.TicketRepository {
	return &TicketRepository{
		db: db,
	}
}
