package service

import (
	"errors"
	"sm-medical/internal/crypto"
	"sm-medical/internal/repository"
	"time"
)

type AdminService struct {
	userRepo        *repository.UserRepository
	applicationRepo *repository.DoctorApplicationRepository
	loginLogRepo    *repository.LoginLogRepository
}

func NewAdminService() *AdminService {
	return &AdminService{
		userRepo:        repository.NewUserRepository(),
		applicationRepo: repository.NewDoctorApplicationRepository(),
		loginLogRepo:    repository.NewLoginLogRepository(),
	}
}

// ReviewDoctorApplication 审核医生申请
func (s *AdminService) ReviewDoctorApplication(adminID, applicationID int64, status int, rejectReason string) error {
	// 查询申请
	application, err := s.applicationRepo.FindByID(applicationID)
	if err != nil {
		return errors.New("申请不存在")
	}

	if application.Status != 0 {
		return errors.New("申请已被审核")
	}

	// 更新申请状态
	application.Status = status
	application.RejectReason = rejectReason
	application.ReviewerID = &adminID
	now := time.Now()
	application.ReviewTime = &now

	if err := s.applicationRepo.Update(application); err != nil {
		return err
	}

	// 如果通过，更新用户角色和信息
	if status == 1 {
		user, err := s.userRepo.FindByID(application.UserID)
		if err != nil {
			return err
		}

		user.Role = "doctor"
		user.Status = 0
		user.RealName = application.RealName
		user.IDCard = application.IDCard
		user.Phone = application.Phone
		user.DoctorCert = application.DoctorCert
		user.DoctorTitle = application.DoctorTitle
		user.DoctorDept = application.DoctorDept
		user.Specialty = application.Specialty
		user.Introduction = application.Introduction
		user.CertNumber = application.CertNumber
		user.CertStatus = "approved"

		if err := s.userRepo.Update(user); err != nil {
			return err
		}
	}

	return nil
}

// GetDoctorApplications 获取医生申请列表
func (s *AdminService) GetDoctorApplications(page, pageSize int, status *int) ([]map[string]interface{}, int64, error) {
	applications, total, err := s.applicationRepo.FindAll(page, pageSize, status)
	if err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for _, app := range applications {
		// 解密敏感信息
		realName, _ := crypto.SM4Decrypt(app.RealName)
		idCard, _ := crypto.SM4Decrypt(app.IDCard)
		phone, _ := crypto.SM4Decrypt(app.Phone)
		email, _ := crypto.SM4Decrypt(app.Email)

		// 获取用户信息
		user, _ := s.userRepo.FindByID(app.UserID)
		username := ""
		if user != nil {
			username = user.Username
		}

		statusText := "待审核"
		if app.Status == 1 {
			statusText = "已通过"
		} else if app.Status == 2 {
			statusText = "已拒绝"
		}

		item := map[string]interface{}{
			"applicationId":  app.ID,
			"applicationNo":  app.ApplicationNo,
			"userId":         app.UserID,
			"username":       username,
			"realName":       realName,
			"idCard":         idCard,
			"phone":          phone,
			"email":          email,
			"doctorCert":     app.DoctorCert,
			"doctorTitle":    app.DoctorTitle,
			"doctorDept":     app.DoctorDept,
			"specialty":      app.Specialty,
			"introduction":   app.Introduction,
			"certNumber":     app.CertNumber,
			"practiceCert":   app.PracticeCert,
			"hospitalName":   app.HospitalName,
			"status":         app.Status,
			"statusText":     statusText,
			"createdAt":      app.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if app.ReviewTime != nil {
			item["reviewTime"] = app.ReviewTime.Format("2006-01-02 15:04:05")
		}
		if app.RejectReason != "" {
			item["rejectReason"] = app.RejectReason
		}

		result = append(result, item)
	}

	return result, total, nil
}

// GetUsers 获取用户列表
func (s *AdminService) GetUsers(page, pageSize int, identify string, status *int, keyword string) ([]map[string]interface{}, int64, error) {
	users, total, err := s.userRepo.FindAll(page, pageSize, identify, status, keyword)
	if err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for _, user := range users {
		// 解密敏感信息
		email, _ := crypto.SM4Decrypt(user.Email)
		phone, _ := crypto.SM4Decrypt(user.Phone)

		item := map[string]interface{}{
			"userId":    user.ID,
			"username":  user.Username,
			"email":     email,
			"phone":     phone,
			"role":      user.Role,
			"status":    user.Status,
			"createdAt": user.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if user.LastLoginTime != nil {
			item["lastLoginTime"] = user.LastLoginTime.Format("2006-01-02 15:04:05")
		}

		result = append(result, item)
	}

	return result, total, nil
}

// UpdateUserStatus 更新用户状态
func (s *AdminService) UpdateUserStatus(userID int64, status int) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.Status = status
	return s.userRepo.Update(user)
}

// GetLoginLogs 获取登录日志
func (s *AdminService) GetLoginLogs(page, pageSize int, userID *int64, status *int, startTime, endTime string) ([]map[string]interface{}, int64, error) {
	logs, total, err := s.loginLogRepo.FindAll(page, pageSize, userID, status, startTime, endTime)
	if err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for _, log := range logs {
		// 解密IP
		loginIP, _ := crypto.SM4Decrypt(log.LoginIP)

		statusText := "失败"
		if log.Status == 1 {
			statusText = "成功"
		}

		result = append(result, map[string]interface{}{
			"logId":         log.ID,
			"userId":        log.UserID,
			"username":      log.Username,
			"loginIp":       loginIP,
			"loginLocation": log.LoginLocation,
			"browser":       log.Browser,
			"os":            log.OS,
			"status":        log.Status,
			"statusText":    statusText,
			"msg":           log.Msg,
			"loginTime":     log.LoginTime.Format("2006-01-02 15:04:05"),
		})
	}

	return result, total, nil
}
