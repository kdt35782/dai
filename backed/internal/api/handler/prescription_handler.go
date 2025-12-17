package handler

import (
	"github.com/gin-gonic/gin"
	"sm-medical/internal/service"
	"sm-medical/pkg/utils"
	"strconv"
)

type PrescriptionHandler struct {
	prescriptionService *service.PrescriptionService
}

func NewPrescriptionHandler() *PrescriptionHandler {
	return &PrescriptionHandler{
		prescriptionService: service.NewPrescriptionService(),
	}
}

// SearchMedicines 搜索药品
func (h *PrescriptionHandler) SearchMedicines(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	list, total, err := h.prescriptionService.SearchMedicines(keyword, category, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetRecommendedMedicines 获取AI推荐药品
func (h *PrescriptionHandler) GetRecommendedMedicines(c *gin.Context) {
	var req struct {
		AIDiagnosis string `json:"aiDiagnosis" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	medicines, err := h.prescriptionService.GetRecommendedMedicines(req.AIDiagnosis)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, medicines)
}

// GetPrescriptionDetail 获取处方详情
func (h *PrescriptionHandler) GetPrescriptionDetail(c *gin.Context) {
	userID := c.GetInt64("userID")
	prescriptionIDStr := c.Param("prescriptionId")
	prescriptionID, err := strconv.ParseInt(prescriptionIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "处方ID格式错误")
		return
	}

	detail, err := h.prescriptionService.GetPrescriptionDetail(userID, prescriptionID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, detail)
}
