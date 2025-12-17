package repository

import (
	"sm-medical/internal/model"
	"sm-medical/pkg/database"
)

type ConsultationRepository struct{}

func NewConsultationRepository() *ConsultationRepository {
	return &ConsultationRepository{}
}

// Create 创建问诊
func (r *ConsultationRepository) Create(consultation *model.Consultation) error {
	return database.GetDB().Create(consultation).Error
}

// Update 更新问诊
func (r *ConsultationRepository) Update(consultation *model.Consultation) error {
	return database.GetDB().Save(consultation).Error
}

// FindByID 根据ID查找
func (r *ConsultationRepository) FindByID(id int64) (*model.Consultation, error) {
	var consultation model.Consultation
	err := database.GetDB().Where("id = ?", id).First(&consultation).Error
	return &consultation, err
}

// FindByUserRole 根据用户角色查询列表
func (r *ConsultationRepository) FindByUserRole(userID int64, role string, page, pageSize int, status *int) ([]model.Consultation, int64, error) {
	var consultations []model.Consultation
	var total int64

	query := database.GetDB().Model(&model.Consultation{})

	if role == "patient" {
		query = query.Where("patient_id = ?", userID)
	} else if role == "doctor" {
		// 医生可以看到：
		// 1. 已经指定给自己的问诊（doctor_id = 自己的ID）
		// 2. 如果没有筛选状态，也显示待接诊的问诊（doctor_id IS NULL 且 status = 0）
		if status == nil {
			// 没有状态筛选：显示指定给自己的 + 所有待接诊的
			query = query.Where("doctor_id = ? OR (doctor_id IS NULL AND status = 0)", userID)
		} else if *status == 0 {
			// 筛选待接诊：只显示指定给自己且待接诊的
			query = query.Where("doctor_id = ? AND status = 0", userID)
		} else {
			// 筛选其他状态：只显示指定给自己的
			query = query.Where("doctor_id = ?", userID)
		}
	}

	if status != nil && role != "doctor" {
		query = query.Where("status = ?", *status)
	} else if status != nil && role == "doctor" && *status != 0 {
		// 医生筛选非待接诊状态
		query = query.Where("status = ?", *status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&consultations).Error

	// 关联查询用户信息（简化版）
	for i := range consultations {
		var patient model.User
		database.GetDB().Where("id = ?", consultations[i].PatientID).First(&patient)
		consultations[i].PatientName = patient.Username

		if consultations[i].DoctorID != nil {
			var doctor model.User
			database.GetDB().Where("id = ?", *consultations[i].DoctorID).First(&doctor)
			consultations[i].DoctorName = doctor.Username
		}
	}

	return consultations, total, err
}
