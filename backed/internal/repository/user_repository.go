package repository

import (
	"sm-medical/internal/model"
	"sm-medical/pkg/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return database.GetDB().Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return database.GetDB().Save(user).Error
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := database.GetDB().Where("id = ?", id).First(&user).Error
	return &user, err
}

// FindByUsername 根据用户名查找
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.GetDB().Where("username = ?", username).First(&user).Error
	return &user, err
}

// FindByRole 根据角色查找用户列表
func (r *UserRepository) FindByRole(role string) ([]model.User, error) {
	var users []model.User
	err := database.GetDB().Where("identify = ? AND status = ?", role, 0).Find(&users).Error
	return users, err
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := database.GetDB().Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := database.GetDB().Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// FindDoctors 查询医生列表
func (r *UserRepository) FindDoctors(page, pageSize int, dept, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := database.GetDB().Model(&model.User{}).Where("identify = ?", "doctor").Where("status = ?", 0)

	if dept != "" {
		query = query.Where("doctor_dept = ?", dept)
	}

	if keyword != "" {
		query = query.Where("real_name LIKE ? OR doctor_title LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

// CreateDoctorApplication 创建医生申请
func (r *UserRepository) CreateDoctorApplication(application *model.DoctorApplication) error {
	return database.GetDB().Create(application).Error
}

// FindDoctorApplication 查找区生申请
func (r *UserRepository) FindDoctorApplication(userID int64) (*model.DoctorApplication, error) {
	var application model.DoctorApplication
	err := database.GetDB().Where("user_id = ?", userID).Order("created_at DESC").First(&application).Error
	if err != nil {
		return nil, err // 没有找到时返回 nil 和 error
	}
	return &application, nil
}

// FindAll 查询所有用户（分页，管理员用）
func (r *UserRepository) FindAll(page, pageSize int, identify string, status *int, keyword string) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := database.GetDB().Model(&model.User{})

	if identify != "" {
		query = query.Where("identify = ?", identify)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}
