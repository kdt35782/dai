package service

import (
	"errors"
	"log"
	"sm-medical/internal/crypto"
	"sm-medical/internal/model"
	"sm-medical/internal/repository"
)

type RecordService struct {
	repo *repository.RecordRepository
	userRepo *repository.UserRepository
}

func NewRecordService() *RecordService {
	return &RecordService{
		repo: repository.NewRecordRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

// GetList 获取病历列表（支持患者和医生）
func (s *RecordService) GetList(userID int64, page, pageSize int, startDate, endDate string) ([]map[string]interface{}, int64, error) {
	log.Printf("[RecordService.GetList] 查询病历 - 用户ID: %d, 页码: %d, 每页: %d", userID, page, pageSize)
	
	// 查询用户信息以判断角色
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[RecordService.GetList] 查询用户失败: %v", err)
		return nil, 0, err
	}
	
	var records []model.MedicalRecord
	var total int64
	
	// 根据用户角色查询不同的病历
	if user.Role == "doctor" {
		// 医生查看自己接诊的患者病历
		log.Printf("[RecordService.GetList] 医生角色 - 查询接诊病历")
		records, total, err = s.repo.FindByDoctorID(userID, page, pageSize, startDate, endDate)
	} else {
		// 患者查看自己的病历
		log.Printf("[RecordService.GetList] 患者角色 - 查询自己病历")
		records, total, err = s.repo.FindByPatientID(userID, page, pageSize, startDate, endDate)
	}
	
	if err != nil {
		log.Printf("[RecordService.GetList] 查询失败: %v", err)
		return nil, 0, err
	}
	
	log.Printf("[RecordService.GetList] 查询成功 - 找到 %d 条记录, 总数: %d", len(records), total)

	var result []map[string]interface{}
	for _, r := range records {
		// 解密
		chiefComplaint, _ := crypto.SM4Decrypt(r.ChiefComplaint)
		diagnosis, _ := crypto.SM4Decrypt(r.Diagnosis)

		// 根据角色返回不同的字段
		if user.Role == "doctor" {
			// 医生看到的是患者名
			result = append(result, map[string]interface{}{
				"recordId":       r.ID,
				"recordNo":       r.RecordNo,
				"chiefComplaint": chiefComplaint,
				"diagnosis":      diagnosis,
				"patientName":    r.DoctorName, // 注意：这里复用了 DoctorName 字段存储患者名
				"createdAt":      r.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		} else {
			// 患者看到的是医生名
			result = append(result, map[string]interface{}{
				"recordId":       r.ID,
				"recordNo":       r.RecordNo,
				"chiefComplaint": chiefComplaint,
				"diagnosis":      diagnosis,
				"doctorName":     r.DoctorName,
				"doctorDept":     r.DoctorDept,
				"createdAt":      r.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	return result, total, nil
}

// GetDetail 获取病历详情（支持患者和医生）
func (s *RecordService) GetDetail(userID, recordID int64) (map[string]interface{}, error) {
	log.Printf("[RecordService.GetDetail] 查询病历详情 - 用户ID: %d, 病历ID: %d", userID, recordID)
	
	record, err := s.repo.FindByID(recordID)
	if err != nil {
		log.Printf("[RecordService.GetDetail] 病历不存在: %v", err)
		return nil, errors.New("病历不存在")
	}

	// 查询用户信息以判断角色
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[RecordService.GetDetail] 查询用户失败: %v", err)
		return nil, err
	}

	// 权限检查：患者本人 或 接诊医生 可以查看
	isPatient := record.PatientID == userID
	isDoctor := record.DoctorID != nil && *record.DoctorID == userID
	
	if !isPatient && !isDoctor {
		log.Printf("[RecordService.GetDetail] 无权限访问 - 用户角色: %s, 患者ID: %d, 医生ID: %v", 
			user.Role, record.PatientID, record.DoctorID)
		return nil, errors.New("无权限访问")
	}
	
	log.Printf("[RecordService.GetDetail] 权限验证通过 - 患者: %v, 医生: %v", isPatient, isDoctor)

	// 解密数据
	chiefComplaint, _ := crypto.SM4Decrypt(record.ChiefComplaint)
	diagnosis, _ := crypto.SM4Decrypt(record.Diagnosis)
	treatment, _ := crypto.SM4Decrypt(record.Treatment)

	result := map[string]interface{}{
		"recordId":       record.ID,
		"recordNo":       record.RecordNo,
		"chiefComplaint": chiefComplaint,
		"diagnosis":      diagnosis,
		"treatment":      treatment,
		"doctorName":     record.DoctorName,
		"doctorDept":     record.DoctorDept,
		"aiAdvice":       record.AIAdvice,
		"createdAt":      record.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return result, nil
}
