package service

import (
	"errors"
	"fmt"
	"log"
	"sm-medical/internal/crypto"
	"sm-medical/internal/model"
	"sm-medical/internal/repository"
	"sm-medical/pkg/utils"
	"time"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// Register 用户注册
func (s *UserService) Register(username, password, email, phone string) (*model.User, error) {
	// 检查用户名是否存在
	if exists, _ := s.userRepo.ExistsByUsername(username); exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否存在（需要先加密再查询）
	encryptedEmailForCheck, _ := crypto.SM4Encrypt(email)
	if exists, _ := s.userRepo.ExistsByEmail(encryptedEmailForCheck); exists {
		return nil, errors.New("邮箱已被注册")
	}

	// SM3加密密码（前端已做一次SM3，后端再加盐哈希）
	hashedPassword := crypto.SM3HashWithSalt(password, username)

	// SM4加密敏感信息
	encryptedEmail, _ := crypto.SM4Encrypt(email)
	encryptedPhone, _ := crypto.SM4Encrypt(phone)

	user := &model.User{
		Username: username,
		Password: hashedPassword,
		Email:    encryptedEmail,
		Phone:    encryptedPhone,
		Role:     "patient",
		Status:   0,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 解密返回
	user.Email = email
	user.Phone = phone

	return user, nil
}

// RegisterDoctor 医生注册
func (s *UserService) RegisterDoctor(username, password, email, phone, realName, idCard, doctorTitle, doctorDept, specialty, introduction, certNumber, certImage string) (*model.User, error) {
	// 检查用户名是否存在
	if exists, _ := s.userRepo.ExistsByUsername(username); exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否存在（需要先加密再查询）
	encryptedEmailForCheck, _ := crypto.SM4Encrypt(email)
	if exists, _ := s.userRepo.ExistsByEmail(encryptedEmailForCheck); exists {
		return nil, errors.New("邮箱已被注册")
	}

	// SM3加密密码（前端已做一次SM3，后端再加盐哈希）
	hashedPassword := crypto.SM3HashWithSalt(password, username)

	// SM4加密敏感信息
	encryptedEmail, _ := crypto.SM4Encrypt(email)
	encryptedPhone, _ := crypto.SM4Encrypt(phone)
	encryptedRealName, _ := crypto.SM4Encrypt(realName)
	encryptedIDCard := ""
	if idCard != "" {
		encryptedIDCard, _ = crypto.SM4Encrypt(idCard)
	}

	user := &model.User{
		Username:     username,
		Password:     hashedPassword,
		Email:        encryptedEmail,
		Phone:        encryptedPhone,
		RealName:     encryptedRealName,
		IDCard:       encryptedIDCard,
		Role:         "doctor",
		Status:       0, // 正常状态
		DoctorCert:   certImage,
		DoctorTitle:  doctorTitle,
		DoctorDept:   doctorDept,
		Specialty:    specialty,
		Introduction: introduction,
		CertNumber:   certNumber,
		CertStatus:   "approved", // 直接设置为已认证
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 解密返回
	user.Email = email
	user.Phone = phone
	user.RealName = realName
	user.IDCard = idCard

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password, clientIP string) (string, map[string]interface{}, error) {
	log.Printf("[Service] 开始登录 - 用户名: %s", username)
	
	// 查询用户
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		log.Printf("[Service] 查找用户失败: %v", err)
		return "", nil, errors.New("用户名或密码错误")
	}
	log.Printf("[Service] 找到用户 - ID: %d, 用户名: %s", user.ID, user.Username)

	// 验证密码
	hashedPassword := crypto.SM3HashWithSalt(password, username)
	log.Printf("[Service] 前端发送的密码: %s", password)
	log.Printf("[Service] 后端计算的哈希: %s", hashedPassword)
	log.Printf("[Service] 数据库存储的哈希: %s", user.Password)
	
	if user.Password != hashedPassword {
		log.Printf("[Service] 密码不匹配!")
		return "", nil, errors.New("用户名或密码错误")
	}
	log.Printf("[Service] 密码验证成功")

	// 检查账号状态
	if user.Status == 1 {
		return "", nil, errors.New("账号已被禁用")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}

	// 更新最后登录时间和IP
	now := time.Now()
	encryptedIP, _ := crypto.SM4Encrypt(clientIP)
	user.LastLoginTime = &now
	user.LastLoginIP = encryptedIP
	s.userRepo.Update(user)

	// 解密敏感信息返回
	email, _ := crypto.SM4Decrypt(user.Email)
	phone, _ := crypto.SM4Decrypt(user.Phone)

	userInfo := map[string]interface{}{
		"userId":      user.ID,
		"username":    user.Username,
		"email":       email,
		"phone":       phone,
		"role":        user.Role,
		"avatar":      user.Avatar,
		"gender":      user.Gender,
		"birthDate":   user.BirthDate,
		"status":      user.Status,
		"realName":    user.RealName,
		"doctorTitle": user.DoctorTitle,
		"doctorDept":  user.DoctorDept,
		"certStatus":  user.CertStatus,
	}

	return token, userInfo, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID int64) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 解密敏感信息
	email, _ := crypto.SM4Decrypt(user.Email)
	phone, _ := crypto.SM4Decrypt(user.Phone)
	realName, _ := crypto.SM4Decrypt(user.RealName)

	userInfo := map[string]interface{}{
		"userId":       user.ID,
		"username":     user.Username,
		"email":        email,
		"phone":        phone,
		"realName":     realName,
		"role":         user.Role,
		"avatar":       user.Avatar,
		"gender":       user.Gender,
		"birthDate":    user.BirthDate,
		"status":       user.Status,
		"doctorTitle":  user.DoctorTitle,
		"doctorDept":   user.DoctorDept,
		"specialty":    user.Specialty,
		"introduction": user.Introduction,
		"certNumber":   user.CertNumber,
		"certStatus":   user.CertStatus,
		"createdAt":    user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if user.LastLoginTime != nil {
		userInfo["lastLoginTime"] = user.LastLoginTime.Format("2006-01-02 15:04:05")
	}

	return userInfo, nil
}

// UpdateProfile 更新用户信息
func (s *UserService) UpdateProfile(userID int64, avatar, realName string, gender int, birthDate, phone, email string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if avatar != "" {
		user.Avatar = avatar
	}
	if realName != "" {
		encrypted, _ := crypto.SM4Encrypt(realName)
		user.RealName = encrypted
	}
	if gender >= 0 {
		user.Gender = gender
	}
	if birthDate != "" {
		user.BirthDate = birthDate
	}
	if phone != "" {
		encrypted, _ := crypto.SM4Encrypt(phone)
		user.Phone = encrypted
	}
	if email != "" {
		encrypted, _ := crypto.SM4Encrypt(email)
		user.Email = encrypted
	}

	return s.userRepo.Update(user)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	hashedOldPassword := crypto.SM3HashWithSalt(oldPassword, user.Username)
	if user.Password != hashedOldPassword {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedNewPassword := crypto.SM3HashWithSalt(newPassword, user.Username)
	user.Password = hashedNewPassword

	return s.userRepo.Update(user)
}

// ApplyDoctor 申请成为医生
func (s *UserService) ApplyDoctor(userID int64, realName, idCard, phone, certImage, doctorTitle, doctorDept, specialty, introduction, certNumber string) (int64, error) {
	// 检查是否已是医生
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, err
	}

	if user.Role == "doctor" {
		return 0, errors.New("您已是认证医生")
	}

	// 检查是否已有待审核的申请
	existing, _ := s.userRepo.FindDoctorApplication(userID)
	if existing != nil && existing.Status == 0 {
		return 0, errors.New("您已有待审核的申请")
	}

	// 生成申请编号：DA + 时间戳 + 4位用户ID
	applicationNo := fmt.Sprintf("DA%d%04d", time.Now().Unix(), userID%10000)

	// 加密敏感信息
	encryptedRealName, _ := crypto.SM4Encrypt(realName)
	encryptedIDCard, _ := crypto.SM4Encrypt(idCard)
	encryptedPhone, _ := crypto.SM4Encrypt(phone)
	encryptedEmail, _ := crypto.SM4Encrypt(user.Email) // 使用用户的邮箱

	application := &model.DoctorApplication{
		UserID:        userID,
		ApplicationNo: applicationNo,
		RealName:      encryptedRealName,
		IDCard:        encryptedIDCard,
		Phone:         encryptedPhone,
		Email:         encryptedEmail,
		DoctorCert:    certImage,
		DoctorTitle:   doctorTitle,
		DoctorDept:    doctorDept,
		Specialty:     specialty,
		Introduction:  introduction,
		CertNumber:    certNumber,
		Status:        0,
	}

	if err := s.userRepo.CreateDoctorApplication(application); err != nil {
		return 0, err
	}

	return application.ID, nil
}

// GetDoctorApplication 获取医生申请
func (s *UserService) GetDoctorApplication(userID int64) (map[string]interface{}, error) {
	application, err := s.userRepo.FindDoctorApplication(userID)
	if err != nil {
		return nil, err
	}

	// 解密
	realName, _ := crypto.SM4Decrypt(application.RealName)

	statusText := "待审核"
	if application.Status == 1 {
		statusText = "已通过"
	} else if application.Status == 2 {
		statusText = "已拒绝"
	}

	result := map[string]interface{}{
		"applicationId": application.ID,
		"status":        application.Status,
		"statusText":    statusText,
		"realName":      realName,
		"doctorTitle":   application.DoctorTitle,
		"doctorDept":    application.DoctorDept,
		"createdAt":     application.CreatedAt.Format("2006-01-02 15:04:05"),
		"rejectReason":  application.RejectReason,
	}

	if application.ReviewTime != nil {
		result["reviewTime"] = application.ReviewTime.Format("2006-01-02 15:04:05")
	}

	return result, nil
}

// GetDoctors 获取医生列表
func (s *UserService) GetDoctors(page, pageSize int, dept, keyword string) ([]map[string]interface{}, int64, error) {
	users, total, err := s.userRepo.FindDoctors(page, pageSize, dept, keyword)
	if err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for _, user := range users {
		realName, _ := crypto.SM4Decrypt(user.RealName)
		
		result = append(result, map[string]interface{}{
			"userId":       user.ID,
			"username":     user.Username,
			"realName":     realName,
			"avatar":       user.Avatar,
			"doctorTitle":  user.DoctorTitle,
			"doctorDept":   user.DoctorDept,
			"specialty":    user.Specialty,
			"introduction": user.Introduction,
			"gender":       user.Gender,
		})
	}

	return result, total, nil
}
