package repository

import (
	"sm-medical/internal/model"
	"sm-medical/pkg/database"
)

type DoctorApplicationRepository struct{}

func NewDoctorApplicationRepository() *DoctorApplicationRepository {
	return &DoctorApplicationRepository{}
}

// Create 创建医生申请
func (r *DoctorApplicationRepository) Create(application *model.DoctorApplication) error {
	return database.DB.Create(application).Error
}

// Update 更新医生申请
func (r *DoctorApplicationRepository) Update(application *model.DoctorApplication) error {
	return database.DB.Save(application).Error
}

// FindByID 根据ID查询
func (r *DoctorApplicationRepository) FindByID(id int64) (*model.DoctorApplication, error) {
	var application model.DoctorApplication
	err := database.DB.Where("id = ?", id).First(&application).Error
	return &application, err
}

// FindByUserID 根据用户ID查询
func (r *DoctorApplicationRepository) FindByUserID(userID int64) (*model.DoctorApplication, error) {
	var application model.DoctorApplication
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").First(&application).Error
	return &application, err
}

// ExistsPendingByUserID 检查用户是否有待审核的申请
func (r *DoctorApplicationRepository) ExistsPendingByUserID(userID int64) (bool, error) {
	var count int64
	err := database.DB.Model(&model.DoctorApplication{}).
		Where("user_id = ? AND status = 0", userID).
		Count(&count).Error
	return count > 0, err
}

// FindAll 查询所有申请（分页）
func (r *DoctorApplicationRepository) FindAll(page, pageSize int, status *int) ([]*model.DoctorApplication, int64, error) {
	var applications []*model.DoctorApplication
	var total int64

	query := database.DB.Model(&model.DoctorApplication{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&applications).Error

	return applications, total, err
}
