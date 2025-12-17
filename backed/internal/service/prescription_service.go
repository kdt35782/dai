package service

import (
	"errors"
	"fmt"
	"sm-medical/internal/crypto"
	"sm-medical/internal/model"
	"sm-medical/internal/repository"
	"strings"
	"time"
)

type PrescriptionService struct {
	prescriptionRepo *repository.PrescriptionRepository
	medicineRepo     *repository.MedicineRepository
	consultationRepo *repository.ConsultationRepository
}

func NewPrescriptionService() *PrescriptionService {
	return &PrescriptionService{
		prescriptionRepo: repository.NewPrescriptionRepository(),
		medicineRepo:     repository.NewMedicineRepository(),
		consultationRepo: repository.NewConsultationRepository(),
	}
}

// CreatePrescription 创建处方
func (s *PrescriptionService) CreatePrescription(doctorID, consultationID int64, diagnosis string, medicines []map[string]interface{}) (map[string]interface{}, error) {
	// 验证问诊权限
	consultation, err := s.consultationRepo.FindByID(consultationID)
	if err != nil {
		return nil, errors.New("问诊不存在")
	}

	if consultation.DoctorID == nil || *consultation.DoctorID != doctorID {
		return nil, errors.New("无权限操作")
	}

	if len(medicines) == 0 {
		return nil, errors.New("处方药品不能为空")
	}

	// 生成处方编号
	prescriptionNo := fmt.Sprintf("RX%d", time.Now().Unix())

	// 加密诊断
	encryptedDiagnosis, _ := crypto.SM4Encrypt(diagnosis)

	// 先创建处方主表(不计算总金额,后面更新)
	prescription := &model.Prescription{
		PrescriptionNo:   prescriptionNo,
		ConsultationID:   consultationID,
		PatientID:        consultation.PatientID,
		DoctorID:         doctorID,
		Diagnosis:        encryptedDiagnosis,
		PrescriptionType: 1,
		TotalAmount:      0, // 先设为0,添加明细后再更新
		Status:           1, // 已审核(简化流程)
	}

	// 生成数据哈希(先用临时值)
	dataForHash := fmt.Sprintf("%d|%s|0", consultationID, diagnosis)
	prescription.DataHash = crypto.SM3Hash(dataForHash)

	if err := s.prescriptionRepo.Create(prescription); err != nil {
		return nil, err
	}

	// 创建处方明细并计算总金额
	var details []map[string]interface{}
	var totalAmount float64
	for _, med := range medicines {
		medicineID := int64(med["medicineId"].(float64))
		quantity := int(med["quantity"].(float64))
		
		// 获取药品信息
		medicine, err := s.medicineRepo.GetByID(medicineID)
		if err != nil {
			continue
		}

		totalPrice := float64(quantity) * medicine.Price
		totalAmount += totalPrice

		detail := &model.PrescriptionDetail{
			PrescriptionID: prescription.ID,
			MedicineID:     medicineID,
			MedicineName:   medicine.MedicineName,
			Specification:  medicine.Specification,
			Quantity:       quantity,
			Unit:           medicine.Unit,
			UnitPrice:      medicine.Price,
			TotalPrice:     totalPrice,
			Usage:          med["usage"].(string),
			Frequency:      med["frequency"].(string),
			Dosage:         med["dosage"].(string),
			Duration:       med["duration"].(string),
		}

		if notes, ok := med["notes"].(string); ok {
			detail.Notes = notes
		}

		if err := s.prescriptionRepo.CreateDetail(detail); err == nil {
			details = append(details, map[string]interface{}{
				"medicineName":   detail.MedicineName,
				"specification":  detail.Specification,
				"quantity":       detail.Quantity,
				"unit":           detail.Unit,
				"usage":          detail.Usage,
				"frequency":      detail.Frequency,
				"dosage":         detail.Dosage,
				"duration":       detail.Duration,
				"totalPrice":     detail.TotalPrice,
			})
		}
	}
	
	// 更新处方总金额
	prescription.TotalAmount = totalAmount
	if err := s.prescriptionRepo.Update(prescription); err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"prescriptionId": prescription.ID,
		"prescriptionNo": prescription.PrescriptionNo,
		"totalAmount":    prescription.TotalAmount,
		"details":        details,
		"createdAt":      prescription.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return result, nil
}

// GetPrescriptionDetail 获取处方详情
func (s *PrescriptionService) GetPrescriptionDetail(userID, prescriptionID int64) (map[string]interface{}, error) {
	prescription, err := s.prescriptionRepo.GetByID(prescriptionID)
	if err != nil {
		return nil, errors.New("处方不存在")
	}

	// 权限检查
	if prescription.PatientID != userID && prescription.DoctorID != userID {
		return nil, errors.New("无权限访问")
	}

	// 解密诊断
	diagnosis, _ := crypto.SM4Decrypt(prescription.Diagnosis)

	// 获取处方明细
	details, _ := s.prescriptionRepo.GetDetailsByPrescriptionID(prescriptionID)

	var detailList []map[string]interface{}
	for _, detail := range details {
		detailList = append(detailList, map[string]interface{}{
			"medicineName":  detail.MedicineName,
			"specification": detail.Specification,
			"quantity":      detail.Quantity,
			"unit":          detail.Unit,
			"unitPrice":     detail.UnitPrice,
			"totalPrice":    detail.TotalPrice,
			"usage":         detail.Usage,
			"frequency":     detail.Frequency,
			"dosage":        detail.Dosage,
			"duration":      detail.Duration,
			"notes":         detail.Notes,
		})
	}

	result := map[string]interface{}{
		"prescriptionId":   prescription.ID,
		"prescriptionNo":   prescription.PrescriptionNo,
		"consultationId":   prescription.ConsultationID,
		"diagnosis":        diagnosis,
		"prescriptionType": prescription.PrescriptionType,
		"totalAmount":      prescription.TotalAmount,
		"status":           prescription.Status,
		"details":          detailList,
		"createdAt":        prescription.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return result, nil
}

// SearchMedicines 搜索药品
func (s *PrescriptionService) SearchMedicines(keyword, category string, page, pageSize int) ([]map[string]interface{}, int64, error) {
	medicines, total, err := s.medicineRepo.SearchMedicines(keyword, category, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for _, med := range medicines {
		result = append(result, map[string]interface{}{
			"medicineId":       med.ID,
			"medicineCode":     med.MedicineCode,
			"medicineName":     med.MedicineName,
			"commonName":       med.CommonName,
			"medicineType":     med.MedicineType,
			"category":         med.Category,
			"specification":    med.Specification,
			"dosageForm":       med.DosageForm,
			"manufacturer":     med.Manufacturer,
			"price":            med.Price,
			"unit":             med.Unit,
			"indications":      med.Indications,
			"usageDosage":      med.UsageDosage,
			"prescriptionType": med.PrescriptionType,
			"isOtc":            med.IsOTC,
			"stockQuantity":    med.StockQuantity,
		})
	}

	return result, total, nil
}

// GetRecommendedMedicines 根据AI诊断推荐药品
func (s *PrescriptionService) GetRecommendedMedicines(aiDiagnosis string) ([]map[string]interface{}, error) {
	// 解析AI诊断,提取可能的疾病类型
	categories := s.extractMedicineCategories(aiDiagnosis)
	
	if len(categories) == 0 {
		// 如果没有匹配的分类,返回常用药
		categories = []string{"感冒药", "止咳化痰", "消化系统", "维生素"}
	}

	medicines, err := s.medicineRepo.GetRecommendedByAI(categories)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, med := range medicines {
		result = append(result, map[string]interface{}{
			"medicineId":    med.ID,
			"medicineName":  med.MedicineName,
			"specification": med.Specification,
			"price":         med.Price,
			"unit":          med.Unit,
			"category":      med.Category,
			"indications":   med.Indications,
			"usageDosage":   med.UsageDosage,
		})
	}

	return result, nil
}

// extractMedicineCategories 从AI诊断中提取药品分类
func (s *PrescriptionService) extractMedicineCategories(aiDiagnosis string) []string {
	var categories []string
	
	diagnosisLower := strings.ToLower(aiDiagnosis)
	
	// 疾病关键词到药品分类的映射
	keywordMap := map[string]string{
		"感冒":   "感冒药",
		"发热":   "感冒药",
		"咳嗽":   "止咳化痰",
		"咽痛":   "感冒药",
		"高血压": "降压药",
		"血压":   "降压药",
		"糖尿病": "降糖药",
		"血糖":   "降糖药",
		"消化":   "消化系统",
		"胃痛":   "消化系统",
		"腹泻":   "消化系统",
		"感染":   "抗生素",
		"炎症":   "抗生素",
	}

	for keyword, category := range keywordMap {
		if strings.Contains(diagnosisLower, keyword) {
			categories = append(categories, category)
		}
	}

	return categories
}

// FormatPrescriptionForRecord 将处方格式化为病历文本
func (s *PrescriptionService) FormatPrescriptionForRecord(prescriptionID int64) (string, error) {
	prescription, err := s.prescriptionRepo.GetByID(prescriptionID)
	if err != nil {
		return "", err
	}

	details, err := s.prescriptionRepo.GetDetailsByPrescriptionID(prescriptionID)
	if err != nil {
		return "", err
	}

	var prescriptionText strings.Builder
	prescriptionText.WriteString(fmt.Sprintf("处方编号: %s\n", prescription.PrescriptionNo))
	prescriptionText.WriteString(fmt.Sprintf("开方时间: %s\n\n", prescription.CreatedAt.Format("2006-01-02 15:04:05")))
	prescriptionText.WriteString("药品清单:\n")

	for i, detail := range details {
		prescriptionText.WriteString(fmt.Sprintf("%d. %s %s\n", i+1, detail.MedicineName, detail.Specification))
		prescriptionText.WriteString(fmt.Sprintf("   数量: %d%s\n", detail.Quantity, detail.Unit))
		prescriptionText.WriteString(fmt.Sprintf("   用法: %s\n", detail.Usage))
		prescriptionText.WriteString(fmt.Sprintf("   用量: %s %s\n", detail.Dosage, detail.Frequency))
		prescriptionText.WriteString(fmt.Sprintf("   疗程: %s\n", detail.Duration))
		if detail.Notes != "" {
			prescriptionText.WriteString(fmt.Sprintf("   备注: %s\n", detail.Notes))
		}
		prescriptionText.WriteString("\n")
	}

	prescriptionText.WriteString(fmt.Sprintf("总金额: %.2f元\n", prescription.TotalAmount))

	return prescriptionText.String(), nil
}
