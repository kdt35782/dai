package repository

import (
	"sm-medical/internal/model"
	"sm-medical/pkg/database"
)

type MedicineRepository struct{}

func NewMedicineRepository() *MedicineRepository {
	return &MedicineRepository{}
}

// SearchMedicines 搜索药品
func (r *MedicineRepository) SearchMedicines(keyword string, category string, page, pageSize int) ([]model.Medicine, int64, error) {
	var medicines []model.Medicine
	var total int64

	db := database.DB.Model(&model.Medicine{}).Where("status = 1")

	if keyword != "" {
		db = db.Where("medicine_name LIKE ? OR common_name LIKE ? OR medicine_code = ?", 
			"%"+keyword+"%", "%"+keyword+"%", keyword)
	}

	if category != "" {
		db = db.Where("category = ?", category)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Order("medicine_name ASC").Find(&medicines).Error
	
	return medicines, total, err
}

// GetByID 根据ID获取药品
func (r *MedicineRepository) GetByID(id int64) (*model.Medicine, error) {
	var medicine model.Medicine
	err := database.DB.Where("id = ? AND status = 1", id).First(&medicine).Error
	return &medicine, err
}

// GetByCategory 根据分类获取药品
func (r *MedicineRepository) GetByCategory(category string) ([]model.Medicine, error) {
	var medicines []model.Medicine
	err := database.DB.Where("category = ? AND status = 1", category).Find(&medicines).Error
	return medicines, err
}

// GetRecommendedByAI 根据AI诊断推荐药品
func (r *MedicineRepository) GetRecommendedByAI(categories []string) ([]model.Medicine, error) {
	var medicines []model.Medicine
	err := database.DB.Where("category IN ? AND status = 1", categories).
		Limit(20).
		Find(&medicines).Error
	return medicines, err
}
