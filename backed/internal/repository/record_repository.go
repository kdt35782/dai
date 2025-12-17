package repository

import (
	"sm-medical/internal/model"
	"sm-medical/pkg/database"
)

type RecordRepository struct{}

func NewRecordRepository() *RecordRepository {
	return &RecordRepository{}
}

// Create 创建病历
func (r *RecordRepository) Create(record *model.MedicalRecord) error {
	return database.GetDB().Create(record).Error
}

// FindByID 根据ID查找
func (r *RecordRepository) FindByID(id int64) (*model.MedicalRecord, error) {
	var record model.MedicalRecord
	err := database.GetDB().Where("id = ?", id).First(&record).Error
	
	// 关联查询医生信息
	if record.DoctorID != nil {
		var doctor model.User
		database.GetDB().Where("id = ?", *record.DoctorID).First(&doctor)
		record.DoctorName = doctor.Username
		record.DoctorDept = doctor.DoctorDept
	}

	return &record, err
}

// FindByPatientID 根据患者ID查询列表
func (r *RecordRepository) FindByPatientID(patientID int64, page, pageSize int, startDate, endDate string) ([]model.MedicalRecord, int64, error) {
	var records []model.MedicalRecord
	var total int64

	query := database.GetDB().Model(&model.MedicalRecord{}).Where("patient_id = ?", patientID)

	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error

	// 关联查询医生信息
	for i := range records {
		if records[i].DoctorID != nil {
			var doctor model.User
			database.GetDB().Where("id = ?", *records[i].DoctorID).First(&doctor)
			records[i].DoctorName = doctor.Username
			records[i].DoctorDept = doctor.DoctorDept
		}
	}

	return records, total, err
}

// FindByDoctorID 根据医生ID查询列表
func (r *RecordRepository) FindByDoctorID(doctorID int64, page, pageSize int, startDate, endDate string) ([]model.MedicalRecord, int64, error) {
	var records []model.MedicalRecord
	var total int64

	query := database.GetDB().Model(&model.MedicalRecord{}).Where("doctor_id = ?", doctorID)

	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error

	// 关联查询患者信息
	for i := range records {
		var patient model.User
		database.GetDB().Where("id = ?", records[i].PatientID).First(&patient)
		records[i].DoctorName = patient.Username // 使用 DoctorName 字段存储患者名，前端会处理
	}

	return records, total, err
}
