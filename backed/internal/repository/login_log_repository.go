package repository

import (
	"sm-medical/internal/model"
	"sm-medical/pkg/database"
)

type LoginLogRepository struct{}

func NewLoginLogRepository() *LoginLogRepository {
	return &LoginLogRepository{}
}

// Create 创建登录日志
func (r *LoginLogRepository) Create(log *model.LoginLog) error {
	return database.DB.Create(log).Error
}

// FindAll 查询所有日志（分页）
func (r *LoginLogRepository) FindAll(page, pageSize int, userID *int64, status *int, startTime, endTime string) ([]*model.LoginLog, int64, error) {
	var logs []*model.LoginLog
	var total int64

	query := database.DB.Model(&model.LoginLog{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if startTime != "" {
		query = query.Where("login_time >= ?", startTime)
	}

	if endTime != "" {
		query = query.Where("login_time <= ?", endTime)
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("login_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}
