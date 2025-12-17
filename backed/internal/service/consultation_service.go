package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sm-medical/internal/crypto"
	"sm-medical/internal/model"
	"sm-medical/internal/repository"
	"time"
)

type ConsultationService struct {
	repo              *repository.ConsultationRepository
	userRepo          *repository.UserRepository
	recordRepo        *repository.RecordRepository
	prescriptionService *PrescriptionService
	triageService     *TriageService
}

func NewConsultationService() *ConsultationService {
	return &ConsultationService{
		repo:              repository.NewConsultationRepository(),
		userRepo:          repository.NewUserRepository(),
		recordRepo:        repository.NewRecordRepository(),
		prescriptionService: NewPrescriptionService(),
		triageService:     NewTriageService(),
	}
}

// Create 创建问诊
func (s *ConsultationService) Create(patientID int64, doctorID *int64, chiefComplaint string, symptoms map[string]interface{}, needAI bool) (map[string]interface{}, error) {
	// 生成问诊编号
	consultationNo := fmt.Sprintf("CN%d", time.Now().Unix())

	// 加密主诉
	encryptedComplaint, _ := crypto.SM4Encrypt(chiefComplaint)

	// 序列化症状数据
	symptomsJSON, _ := json.Marshal(symptoms)

	consultation := &model.Consultation{
		PatientID:         patientID,
		DoctorID:          doctorID,
		ConsultationNo:    consultationNo,
		ChiefComplaint:    encryptedComplaint,
		SymptomsEncrypted: string(symptomsJSON),
		NeedAI:            needAI,
		Status:            0,
	}

	// AI智能诊断
	var aiResult *AIResult
	if needAI {
		result := s.performAIDiagnosis(chiefComplaint, symptoms)
		aiResult = &result
		consultation.AIRiskScore = &aiResult.RiskScore
		consultation.AIDiagnosis = aiResult.Diagnosis
		consultation.AISuggestions = aiResult.Suggestions
		consultation.RecommendedDept = aiResult.RecommendedDept
	}

	// 智能分诊:如果未指定医生且有AI推荐科室,自动分配医生
	if doctorID == nil && needAI && aiResult != nil && aiResult.RecommendedDept != "" {
		log.Printf("[问诊创建] 未指定医生,启动智能分诊 - 推荐科室: %s", aiResult.RecommendedDept)
		
		assignedDoctor, assignReason, err := s.triageService.AutoAssignDoctor(0, aiResult.RecommendedDept)
		if err == nil && assignedDoctor != nil {
			// 分诊成功
			consultation.DoctorID = &assignedDoctor.ID
			consultation.AutoAssigned = true
			consultation.AssignedReason = assignReason
			log.Printf("[智能分诊] 成功分配医生: %s(%s) - 原因: %s", 
				assignedDoctor.RealName, assignedDoctor.DoctorDept, assignReason)
			
			// 更新医生负载
			s.triageService.UpdateDoctorWorkload(assignedDoctor.ID, 1)
		} else {
			log.Printf("[智能分诊] 自动分配失败: %v,问诊将进入待接诊队列", err)
		}
	}

	if err := s.repo.Create(consultation); err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"consultationId": consultation.ID,
		"consultationNo": consultation.ConsultationNo,
		"status":         consultation.Status,
		"createdAt":      consultation.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 如果自动分配了医生,返回分诊信息
	if consultation.AutoAssigned && consultation.DoctorID != nil {
		doctor, _ := s.userRepo.FindByID(*consultation.DoctorID)
		if doctor != nil {
			result["autoAssigned"] = true
			result["assignedDoctor"] = map[string]interface{}{
				"doctorId":   doctor.ID,
				"doctorName": doctor.RealName,
				"doctorDept": doctor.DoctorDept,
				"doctorTitle": doctor.DoctorTitle,
			}
			result["assignedReason"] = consultation.AssignedReason
		}
	}

	if needAI && aiResult != nil && consultation.AIRiskScore != nil {
		result["aiDiagnosis"] = map[string]interface{}{
			"riskScore":        *consultation.AIRiskScore,
			"riskLevel":        s.getRiskLevel(*consultation.AIRiskScore),
			"diagnosis":        consultation.AIDiagnosis,
			"suggestions":      consultation.AISuggestions,
			"possibleDiseases": aiResult.PossibleDiseases,
			"recommendedDept":  aiResult.RecommendedDept,
			"urgencyLevel":     aiResult.UrgencyLevel,
			"detailedAnalysis": aiResult.DetailedAnalysis,
			"lifestyleAdvice":  aiResult.LifestyleAdvice,
			"followUpAdvice":   aiResult.FollowUpAdvice,
		}
	}

	return result, nil
}

// GetList 获取问诊列表
func (s *ConsultationService) GetList(userID int64, role string, page, pageSize int, status *int) ([]map[string]interface{}, int64, error) {
	consultations, total, err := s.repo.FindByUserRole(userID, role, page, pageSize, status)
	if err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for _, c := range consultations {
		chiefComplaint, _ := crypto.SM4Decrypt(c.ChiefComplaint)
		
		statusText := "待接诊"
		if c.Status == 1 {
			statusText = "问诊中"
		} else if c.Status == 2 {
			statusText = "已完成"
		}

		result = append(result, map[string]interface{}{
			"consultationId": c.ID,
			"consultationNo": c.ConsultationNo,
			"patientName":    c.PatientName,
			"doctorName":     c.DoctorName,
			"chiefComplaint": chiefComplaint,
			"status":         c.Status,
			"statusText":     statusText,
			"needAI":         c.NeedAI,
			"createdAt":      c.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, total, nil
}

// GetDetail 获取问诊详情
func (s *ConsultationService) GetDetail(userID, consultationID int64) (map[string]interface{}, error) {
	consultation, err := s.repo.FindByID(consultationID)
	if err != nil {
		return nil, errors.New("问诊不存在")
	}

	// 权限检查
	if consultation.PatientID != userID && (consultation.DoctorID == nil || *consultation.DoctorID != userID) {
		return nil, errors.New("无权限访问")
	}

	// 解密数据
	chiefComplaint, _ := crypto.SM4Decrypt(consultation.ChiefComplaint)
	
	var symptoms map[string]interface{}
	json.Unmarshal([]byte(consultation.SymptomsEncrypted), &symptoms)

	statusText := "待接诊"
	if consultation.Status == 1 {
		statusText = "问诊中"
	} else if consultation.Status == 2 {
		statusText = "已完成"
	}

	result := map[string]interface{}{
		"consultationId": consultation.ID,
		"consultationNo": consultation.ConsultationNo,
		"chiefComplaint": chiefComplaint,
		"symptoms":       symptoms,
		"status":         consultation.Status,
		"statusText":     statusText,
		"needAI":         consultation.NeedAI,
		"createdAt":      consultation.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if consultation.AIRiskScore != nil {
		// 重新执行AI诊断以获取完整结果(因为数据库只存储了部分字段)
		aiResult := s.performAIDiagnosis(chiefComplaint, symptoms)
		result["aiDiagnosis"] = map[string]interface{}{
			"riskScore":        *consultation.AIRiskScore,
			"riskLevel":        s.getRiskLevel(*consultation.AIRiskScore),
			"diagnosis":        consultation.AIDiagnosis,
			"suggestions":      consultation.AISuggestions,
			"possibleDiseases": aiResult.PossibleDiseases,
			"recommendedDept":  aiResult.RecommendedDept,
			"urgencyLevel":     aiResult.UrgencyLevel,
			"detailedAnalysis": aiResult.DetailedAnalysis,
			"lifestyleAdvice":  aiResult.LifestyleAdvice,
			"followUpAdvice":   aiResult.FollowUpAdvice,
		}
	}

	return result, nil
}

// Accept 医生接诊
func (s *ConsultationService) Accept(doctorID, consultationID int64) error {
	consultation, err := s.repo.FindByID(consultationID)
	if err != nil {
		return errors.New("问诊不存在")
	}

	if consultation.Status != 0 {
		return errors.New("问诊状态不正确")
	}

	consultation.DoctorID = &doctorID
	consultation.Status = 1

	// 更新医生负载(+1)
	if err := s.triageService.UpdateDoctorWorkload(doctorID, 1); err != nil {
		log.Printf("[接诊] 更新医生负载失败: %v", err)
	}

	return s.repo.Update(consultation)
}

// Finish 完成问诊(支持处方)
func (s *ConsultationService) Finish(doctorID, consultationID int64, diagnosis string, prescriptionData interface{}) error {
	log.Printf("[Finish] 开始处理 - 医生ID: %d, 问诊ID: %d", doctorID, consultationID)
	
	consultation, err := s.repo.FindByID(consultationID)
	if err != nil {
		log.Printf("[Finish] 查找问诊失败 - 问诊ID: %d, 错误: %v", consultationID, err)
		return errors.New("问诊不存在")
	}
	log.Printf("[Finish] 找到问诊 - ID: %d, 状态: %d, 医生ID: %v", consultation.ID, consultation.Status, consultation.DoctorID)

	if consultation.DoctorID == nil || *consultation.DoctorID != doctorID {
		log.Printf("[Finish] 权限校验失败 - 问诊医生ID: %v, 当前医生ID: %d", consultation.DoctorID, doctorID)
		return errors.New("无权限操作")
	}

	// 处理处方数据
	var prescriptionText string
	if prescriptionData != nil {
		// 解析处方数据
		var medicines []map[string]interface{}
		switch v := prescriptionData.(type) {
		case string:
			// 兼容旧版本:直接传入处方文本
			prescriptionText = v
		case []interface{}:
			// 新版本:传入药品列表,创建处方
			for _, item := range v {
				if med, ok := item.(map[string]interface{}); ok {
					medicines = append(medicines, med)
				}
			}
			
			if len(medicines) > 0 {
				// 创建处方
				log.Printf("[Finish] 创建处方 - 药品数量: %d", len(medicines))
				prescriptionResult, err := s.prescriptionService.CreatePrescription(doctorID, consultationID, diagnosis, medicines)
				if err != nil {
					log.Printf("[Finish] 创建处方失败: %v", err)
					return errors.New("处方开具失败: " + err.Error())
				}
				
				// 格式化处方为文本
				prescriptionID, ok := prescriptionResult["prescriptionId"].(int64)
				if !ok {
					log.Printf("[Finish] 处方ID类型转换失败: %v", prescriptionResult["prescriptionId"])
					return errors.New("处方ID类型错误")
				}
				prescriptionText, _ = s.prescriptionService.FormatPrescriptionForRecord(prescriptionID)
				log.Printf("[Finish] 处方创建成功 - 处方ID: %d", prescriptionID)
			}
		}
	}

	// 如果没有处方,使用空字符串
	if prescriptionText == "" {
		prescriptionText = "无处方"
	}

	// 加密诊断和处方
	encryptedDiagnosis, _ := crypto.SM4Encrypt(diagnosis)
	encryptedPrescription, _ := crypto.SM4Encrypt(prescriptionText)

	consultation.DoctorDiagnosis = encryptedDiagnosis
	consultation.Prescription = encryptedPrescription
	consultation.Status = 2
	now := time.Now()
	consultation.CompletedAt = &now

	log.Printf("[Finish] 更新问诊状态 - ID: %d, 新状态: 2", consultation.ID)
	err = s.repo.Update(consultation)
	if err != nil {
		log.Printf("[Finish] 更新问诊失败: %v", err)
		return err
	}
		
	// 更新医生负载(-1)
	if err := s.triageService.UpdateDoctorWorkload(doctorID, -1); err != nil {
		log.Printf("[Finish] 更新医生负载失败: %v", err)
	}
		
	// 自动创建病历
	log.Printf("[Finish] 开始创建病历 - 问论ID: %d", consultationID)
	err = s.createMedicalRecord(consultation, diagnosis, prescriptionText)
	if err != nil {
		log.Printf("[Finish] 创建病历失败: %v", err)
		// 病历创建失败不影响问诊完成,只记录日志
	}
		
	return nil
}

// createMedicalRecord 创建病历
func (s *ConsultationService) createMedicalRecord(consultation *model.Consultation, diagnosis, prescription string) error {
	// 生成病历编号
	recordNo := fmt.Sprintf("MR%d", time.Now().Unix())
	
	// 解密主诉
	chiefComplaint, _ := crypto.SM4Decrypt(consultation.ChiefComplaint)
	
	// 加密病历数据
	encryptedChiefComplaint, _ := crypto.SM4Encrypt(chiefComplaint)
	encryptedDiagnosis, _ := crypto.SM4Encrypt(diagnosis)
	encryptedTreatment, _ := crypto.SM4Encrypt(prescription)
	
	// 生成数据哈希
	dataForHash := fmt.Sprintf("%s|%s|%s", chiefComplaint, diagnosis, prescription)
	dataHash := crypto.SM3Hash(dataForHash)
	
	record := &model.MedicalRecord{
		RecordNo:       recordNo,
		PatientID:      consultation.PatientID,
		ConsultationID: &consultation.ID,
		RecordType:     2, // 2:在线问诊
		ChiefComplaint: encryptedChiefComplaint,
		Diagnosis:      encryptedDiagnosis,
		Treatment:      encryptedTreatment,
		DoctorID:       consultation.DoctorID,
		AIAdvice:       consultation.AIDiagnosis,
		DataHash:       dataHash,
	}
	
	err := s.recordRepo.Create(record)
	if err != nil {
		return err
	}
	
	log.Printf("[Finish] 病历创建成功 - 病历ID: %d, 病历编号: %s", record.ID, record.RecordNo)
	return nil
}

// AIResult AI诊断结果
type AIResult struct {
	RiskScore          int
	Diagnosis          string
	Suggestions        string
	PossibleDiseases   []string            // 可能疾病列表
	RecommendedDept    string              // 推荐科室
	UrgencyLevel       string              // 紧急程度: normal/attention/urgent/emergency
	DetailedAnalysis   map[string]string   // 分系统详细分析
	LifestyleAdvice    []string            // 生活指导
	FollowUpAdvice     string              // 复诊建议
	Confidence         float64             // 诊断置信度(0-1)
	DataCompleteness   float64             // 数据完整度(0-1)
}

// getRiskLevel 获取风险等级
func (s *ConsultationService) getRiskLevel(score int) string {
	if score >= 70 {
		return "high"
	} else if score >= 40 {
		return "medium"
	} else if score >= 20 {
		return "low"
	}
	return "normal"
}
