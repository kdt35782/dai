package handler

import (
	"github.com/gin-gonic/gin"
	"sm-medical/internal/service"
	"sm-medical/pkg/utils"
	"strconv"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		adminService: service.NewAdminService(),
	}
}

// ReviewDoctorApplication 审核医生申请
func (h *AdminHandler) ReviewDoctorApplication(c *gin.Context) {
	adminID := c.GetInt64("userID")
	role := c.GetString("role")

	if role != "admin" {
		utils.Forbidden(c, "无权限访问")
		return
	}

	var req struct {
		ApplicationID int64  `json:"applicationId" binding:"required"`
		Status        int    `json:"status" binding:"required"` // 1:通过 2:拒绝
		RejectReason  string `json:"rejectReason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if req.Status == 2 && req.RejectReason == "" {
		utils.BadRequest(c, "拒绝时必须填写拒绝原因")
		return
	}

	if err := h.adminService.ReviewDoctorApplication(adminID, req.ApplicationID, req.Status, req.RejectReason); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "审核成功", nil)
}

// GetDoctorApplications 获取医生申请列表
func (h *AdminHandler) GetDoctorApplications(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		utils.Forbidden(c, "无权限访问")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	statusStr := c.Query("status")

	var status *int
	if statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	list, total, err := h.adminService.GetDoctorApplications(page, pageSize, status)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"list":     list,
	})
}

// GetUsers 获取用户列表
func (h *AdminHandler) GetUsers(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		utils.Forbidden(c, "无权限访问")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	identify := c.Query("identify")
	statusStr := c.Query("status")
	keyword := c.Query("keyword")

	var status *int
	if statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	list, total, err := h.adminService.GetUsers(page, pageSize, identify, status, keyword)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"list":     list,
	})
}

// UpdateUserStatus 禁用/启用用户
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		utils.Forbidden(c, "无权限访问")
		return
	}

	var req struct {
		UserID int64 `json:"userId" binding:"required"`
		Status int   `json:"status" binding:"required"` // 0:启用 1:禁用
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.adminService.UpdateUserStatus(req.UserID, req.Status); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "操作成功", nil)
}

// GetLoginLogs 获取登录日志
func (h *AdminHandler) GetLoginLogs(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		utils.Forbidden(c, "无权限访问")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	userIDStr := c.Query("userId")
	statusStr := c.Query("status")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	var userID *int64
	if userIDStr != "" {
		id, _ := strconv.ParseInt(userIDStr, 10, 64)
		userID = &id
	}

	var status *int
	if statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	list, total, err := h.adminService.GetLoginLogs(page, pageSize, userID, status, startTime, endTime)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"list":     list,
	})
}
