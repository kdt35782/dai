package handler

import (
	"github.com/gin-gonic/gin"
	"sm-medical/internal/service"
	"sm-medical/pkg/utils"
	"strconv"
)

type TriageHandler struct {
	triageService *service.TriageService
}

func NewTriageHandler() *TriageHandler {
	return &TriageHandler{
		triageService: service.NewTriageService(),
	}
}

// UpdateOnlineStatus 更新医生在线状态
func (h *TriageHandler) UpdateOnlineStatus(c *gin.Context) {
	userID := c.GetInt64("userID")
	role := c.GetString("role")

	// 只有医生可以更新在线状态
	if role != "doctor" {
		utils.Forbidden(c, "只有医生可以更新在线状态")
		return
	}

	var req struct {
		IsOnline bool `json:"isOnline" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.triageService.UpdateDoctorOnlineStatus(userID, req.IsOnline)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	statusText := "离线"
	if req.IsOnline {
		statusText = "在线"
	}

	utils.SuccessWithMessage(c, "状态更新成功", map[string]interface{}{
		"isOnline": req.IsOnline,
		"status":   statusText,
	})
}

// ManualAssignDoctor 手动触发智能分诊
func (h *TriageHandler) ManualAssignDoctor(c *gin.Context) {
	// userID := c.GetInt64("userID")  // TODO: 用于验证consultationId是否属于当前用户
	role := c.GetString("role")

	// 只有患者可以手动分诊
	if role != "patient" {
		utils.Forbidden(c, "只有患者可以使用此功能")
		return
	}

	var req struct {
		ConsultationID  int64  `json:"consultationId" binding:"required"`
		RecommendedDept string `json:"recommendedDept"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// TODO: 验证consultationId是否属于当前用户
	// TODO: 验证问诊状态是否为待接诊

	assignedDoctor, assignReason, err := h.triageService.AutoAssignDoctor(req.ConsultationID, req.RecommendedDept)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "分诊成功", map[string]interface{}{
		"doctorId":       assignedDoctor.ID,
		"doctorName":     assignedDoctor.RealName,
		"doctorDept":     assignedDoctor.DoctorDept,
		"doctorTitle":    assignedDoctor.DoctorTitle,
		"assignedReason": assignReason,
	})
}

// GetDoctorWorkload 获取医生工作负载信息
func (h *TriageHandler) GetDoctorWorkload(c *gin.Context) {
	doctorIDStr := c.Query("doctorId")
	
	var doctorID int64
	if doctorIDStr != "" {
		id, err := strconv.ParseInt(doctorIDStr, 10, 64)
		if err != nil {
			utils.BadRequest(c, "医生ID格式错误")
			return
		}
		doctorID = id
	} else {
		// 如果没有传doctorId,返回当前用户的负载信息
		userID := c.GetInt64("userID")
		role := c.GetString("role")
		
		if role != "doctor" {
			utils.Forbidden(c, "只有医生可以查看负载信息")
			return
		}
		doctorID = userID
	}

	// TODO: 实现获取医生负载详情的Repository方法
	// workload, err := h.triageService.GetDoctorWorkload(doctorID)
	
	utils.SuccessWithMessage(c, "查询成功", map[string]interface{}{
		"doctorId":      doctorID,
		"currentCount":  0, // TODO: 从数据库获取
		"maxCount":      20,
		"availableSlots": 20,
	})
}
