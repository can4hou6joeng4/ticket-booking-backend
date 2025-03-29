package repositorys

import (
	"context"
	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func (r *EventRepository) CreateOne(ctx context.Context, event *models.Event) (*models.Event, error) {
	res := r.db.Model(event).Create(event)
	if res.Error != nil {
		return nil, res.Error
	}
	return event, res.Error
}
func (r *EventRepository) GetOne(ctx context.Context, eventId int) (*models.Event, error) {
	event := &models.Event{}
	res := r.db.Model(event).Where("id = ?", eventId).First(event)
	if res.Error != nil {
		return nil, res.Error
	}
	return event, nil
}
func (r *EventRepository) GetMany(ctx context.Context) ([]*models.Event, error) {
	events := []*models.Event{}
	res := r.db.Model(&models.Event{}).Order("updated_at desc").Find(&events)
	if res.Error != nil {
		return nil, res.Error
	}
	return events, nil
}
func (r *EventRepository) UpdateOne(ctx context.Context, eventId int, updateData map[string]interface{}) (*models.Event, error) {
	event := &models.Event{}
	updateRes := r.db.Model(event).Where("id = ?", eventId).Updates(updateData)
	if updateRes.Error != nil {
		return nil, updateRes.Error
	}
	getRes := r.db.Model(event).Where("id = ?", eventId).First(event)
	if getRes.Error != nil {
		return nil, getRes.Error
	}
	return event, nil
}
func (r *EventRepository) DeleteOne(ctx context.Context, eventId int) error {
	event := &models.Event{}
	res := r.db.Model(event).Delete(&event, eventId)
	return res.Error
}
func NewEventRepository(db *gorm.DB) models.EventRepository {
	return &EventRepository{
		db: db,
	}
}
