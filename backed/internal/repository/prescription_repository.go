package repository

import (
	"sm-medical/internal/model"
	"sm-medical/pkg/database"
)

type PrescriptionRepository struct{}

func NewPrescriptionRepository() *PrescriptionRepository {
	return &PrescriptionRepository{}
}

// Create 创建处方
func (r *PrescriptionRepository) Create(prescription *model.Prescription) error {
	return database.DB.Create(prescription).Error
}

// CreateDetail 创建处方明细
func (r *PrescriptionRepository) CreateDetail(detail *model.PrescriptionDetail) error {
	return database.DB.Create(detail).Error
}

// GetByID 根据ID获取处方
func (r *PrescriptionRepository) GetByID(id int64) (*model.Prescription, error) {
	var prescription model.Prescription
	err := database.DB.First(&prescription, id).Error
	return &prescription, err
}

// GetByConsultationID 根据问诊ID获取处方
func (r *PrescriptionRepository) GetByConsultationID(consultationID int64) (*model.Prescription, error) {
	var prescription model.Prescription
	err := database.DB.Where("consultation_id = ?", consultationID).First(&prescription).Error
	return &prescription, err
}

// GetDetailsByPrescriptionID 获取处方明细
func (r *PrescriptionRepository) GetDetailsByPrescriptionID(prescriptionID int64) ([]model.PrescriptionDetail, error) {
	var details []model.PrescriptionDetail
	err := database.DB.Where("prescription_id = ?", prescriptionID).Find(&details).Error
	return details, err
}

// GetByPatientID 获取患者的处方列表
func (r *PrescriptionRepository) GetByPatientID(patientID int64, page, pageSize int) ([]model.Prescription, int64, error) {
	var prescriptions []model.Prescription
	var total int64

	offset := (page - 1) * pageSize
	
	err := database.DB.Model(&model.Prescription{}).
		Where("patient_id = ?", patientID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = database.DB.Where("patient_id = ?", patientID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&prescriptions).Error

	return prescriptions, total, err
}

// Update 更新处方
func (r *PrescriptionRepository) Update(prescription *model.Prescription) error {
	return database.DB.Save(prescription).Error
}
