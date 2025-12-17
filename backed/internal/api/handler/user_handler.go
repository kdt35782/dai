package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"sm-medical/internal/model"
	"sm-medical/internal/service"
	"sm-medical/pkg/utils"
	"strconv"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username     string `json:"username" binding:"required"`
		Password     string `json:"password" binding:"required"`
		Email        string `json:"email" binding:"required,email"`
		Phone        string `json:"phone" binding:"required"`
		Role         string `json:"role"` // patient 或 doctor
		// 医生专属字段
		RealName     string `json:"realName"`
		IDCard       string `json:"idCard"`
		DoctorTitle  string `json:"doctorTitle"`
		DoctorDept   string `json:"doctorDept"`
		Specialty    string `json:"specialty"`
		Introduction string `json:"introduction"`
		CertNumber   string `json:"certNumber"`
		CertImage    string `json:"certImage"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 默认角色为患者
	if req.Role == "" {
		req.Role = "patient"
	}

	// 如果是医生注册，验证必填字段
	if req.Role == "doctor" {
		if req.RealName == "" || req.DoctorTitle == "" || req.DoctorDept == "" || 
			req.Specialty == "" || req.CertNumber == "" || req.CertImage == "" {
			utils.BadRequest(c, "医生注册信息不完整")
			return
		}
	}

	// 调用service注册
	var user *model.User
	var err error
	
	if req.Role == "doctor" {
		// 医生注册
		user, err = h.userService.RegisterDoctor(
			req.Username, req.Password, req.Email, req.Phone,
			req.RealName, req.IDCard, req.DoctorTitle, req.DoctorDept,
			req.Specialty, req.Introduction, req.CertNumber, req.CertImage,
		)
	} else {
		// 患者注册
		user, err = h.userService.Register(req.Username, req.Password, req.Email, req.Phone)
	}
	
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "注册成功", gin.H{
		"userId":   user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		LoginType string `json:"loginType"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 添加调试日志
	log.Printf("[登录请求] 用户名: %s", req.Username)
	log.Printf("[登录请求] 密码长度: %d", len(req.Password))
	log.Printf("[登录请求] 密码前16位: %s...", req.Password[:16])

	// 调用service登录
	token, userInfo, err := h.userService.Login(req.Username, req.Password, c.ClientIP())
	if err != nil {
		log.Printf("[登录失败] 用户名: %s, 错误: %v", req.Username, err)
		utils.Unauthorized(c, err.Error())
		return
	}

	log.Printf("[登录成功] 用户名: %s", req.Username)
	utils.SuccessWithMessage(c, "登录成功", gin.H{
		"token":     token,
		"tokenType": "Bearer",
		"expiresIn": 7200,
		"userInfo":  userInfo,
	})
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 支持查询其他用户（医生信息）
	userIDStr := c.Query("userId")
	
	var userID int64
	if userIDStr != "" {
		var err error
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			utils.BadRequest(c, "用户ID格式错误")
			return
		}
	} else {
		// 查询当前用户
		val, exists := c.Get("userID")
		if !exists {
			utils.Unauthorized(c, "请先登录")
			return
		}
		userID = val.(int64)
	}

	userInfo, err := h.userService.GetUserInfo(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, userInfo)
}

// UpdateProfile 更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req struct {
		Avatar    string `json:"avatar"`
		Gender    int    `json:"gender"`
		BirthDate string `json:"birthDate"`
		RealName  string `json:"realName"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.userService.UpdateProfile(userID, req.Avatar, req.RealName, req.Gender, req.BirthDate, req.Phone, req.Email); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req struct {
		OldPassword     string `json:"oldPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		utils.BadRequest(c, "两次密码不一致")
		return
	}

	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "密码修改成功，请重新登录", nil)
}

// ApplyDoctor 申请成为医生
func (h *UserHandler) ApplyDoctor(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req struct {
		RealName     string `json:"realName" binding:"required"`
		IDCard       string `json:"idCard"`
		Phone        string `json:"phone" binding:"required"`
		DoctorCert   string `json:"doctorCert"`
		CertImage    string `json:"certImage" binding:"required"`
		DoctorTitle  string `json:"doctorTitle" binding:"required"`
		DoctorDept   string `json:"doctorDept" binding:"required"`
		Specialty    string `json:"specialty"`
		Introduction string `json:"introduction"`
		CertNumber   string `json:"certNumber"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	applicationID, err := h.userService.ApplyDoctor(userID, req.RealName, req.IDCard, req.Phone,
		req.CertImage, req.DoctorTitle, req.DoctorDept, req.Specialty, req.Introduction, req.CertNumber)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "申请已提交，请等待管理员审核", gin.H{
		"applicationId": applicationID,
	})
}

// GetDoctorApplication 查询医生申请状态
func (h *UserHandler) GetDoctorApplication(c *gin.Context) {
	userID := c.GetInt64("userID")

	application, err := h.userService.GetDoctorApplication(userID)
	if err != nil {
		utils.NotFound(c, "未找到申请记录")
		return
	}

	utils.Success(c, application)
}

// GetDoctors 获取医生列表
func (h *UserHandler) GetDoctors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	dept := c.Query("dept")
	keyword := c.Query("keyword")

	list, total, err := h.userService.GetDoctors(page, pageSize, dept, keyword)
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

// GetDoctorDetail 获取医生详情
func (h *UserHandler) GetDoctorDetail(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "用户ID格式错误")
		return
	}

	userInfo, err := h.userService.GetUserInfo(userID)
	if err != nil {
		utils.NotFound(c, "医生不存在")
		return
	}

	utils.Success(c, userInfo)
}

// Logout 退出登录
func (h *UserHandler) Logout(c *gin.Context) {
	// JWT无状态，客户端删除token即可
	utils.SuccessWithMessage(c, "退出成功", nil)
}
