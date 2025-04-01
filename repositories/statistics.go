package repositories

import (
	"context"

	"github.com/can4hou6joeng4/ticket-booking-project-v1/models"
	"gorm.io/gorm"
)

type StatisticsRepository struct {
	db *gorm.DB
}

func (r *StatisticsRepository) GetCount(ctx context.Context) (*models.Statistics, error) {
	statistics := &models.Statistics{}
	// 获取活动总数
	if err := r.db.Model(&models.Event{}).Count(&statistics.TotalEvents).Error; err != nil {
		return nil, err
	}
	// 获取票券总数
	if err := r.db.Model(&models.Ticket{}).Count(&statistics.TotalTickets).Error; err != nil {
		return nil, err
	}
	// 获取已验证的票券数量
	if err := r.db.Model(&models.Ticket{}).Where("entered = ?", true).Count(&statistics.ValidatedTickets).Error; err != nil {
		return nil, err
	}
	return statistics, nil
}

func NewStatisticsRepository(db *gorm.DB) *StatisticsRepository {
	return &StatisticsRepository{db: db}
}
