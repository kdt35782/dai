package handler

import (
	"github.com/gin-gonic/gin"
	"sm-medical/internal/service"
	"sm-medical/pkg/utils"
	"strconv"
)

type ConsultationHandler struct {
	service *service.ConsultationService
}

func NewConsultationHandler() *ConsultationHandler {
	return &ConsultationHandler{
		service: service.NewConsultationService(),
	}
}

// Create 创建问诊
func (h *ConsultationHandler) Create(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req struct {
		DoctorID       *int64                 `json:"doctorId"`
		ChiefComplaint string                 `json:"chiefComplaint" binding:"required"`
		Symptoms       map[string]interface{} `json:"symptoms" binding:"required"`
		NeedAI         bool                   `json:"needAI"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	consultation, err := h.service.Create(userID, req.DoctorID, req.ChiefComplaint, req.Symptoms, req.NeedAI)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "问诊创建成功", consultation)
}

// GetList 获取问诊列表
func (h *ConsultationHandler) GetList(c *gin.Context) {
	userID := c.GetInt64("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	role := c.DefaultQuery("role", "patient")
	statusStr := c.Query("status")

	var status *int
	if statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	list, total, err := h.service.GetList(userID, role, page, pageSize, status)
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

// GetDetail 获取问诊详情
func (h *ConsultationHandler) GetDetail(c *gin.Context) {
	userID := c.GetInt64("userID")
	consultationIDStr := c.Query("consultationId")
	
	consultationID, err := strconv.ParseInt(consultationIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "问诊ID格式错误")
		return
	}

	detail, err := h.service.GetDetail(userID, consultationID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, detail)
}

// Accept 医生接诊
func (h *ConsultationHandler) Accept(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req struct {
		ConsultationID int64 `json:"consultationId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.Accept(userID, req.ConsultationID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "接诊成功", nil)
}

// Finish 完成问诊
func (h *ConsultationHandler) Finish(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req struct {
		ConsultationID int64                   `json:"consultationId" binding:"required"`
		Diagnosis      string                 `json:"diagnosis" binding:"required"`
		Prescription   interface{}            `json:"prescription"` // 可以是字符串或药品列表
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.Finish(userID, req.ConsultationID, req.Diagnosis, req.Prescription); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "问诊已完成", nil)
}
