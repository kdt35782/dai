package model

import (
	"time"
)

// User 用户表模型
type User struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"userId"`
	Username       string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password       string    `gorm:"type:varchar(128);not null" json:"-"` // SM3加密
	Email          string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"` // SM4加密
	Phone          string    `gorm:"type:varchar(255)" json:"phone"` // SM4加密
	RealName       string    `gorm:"type:varchar(255);column:real_name" json:"realName"` // SM4加密
	IDCard         string    `gorm:"type:varchar(255);column:id_card" json:"idCard"` // SM4加密
	Role           string    `gorm:"type:varchar(20);column:identify;not null;default:user" json:"role"` // user, admin, doctor
	Avatar         string    `gorm:"type:varchar(500)" json:"avatar"`
	Gender         int       `gorm:"type:tinyint;default:0" json:"gender"` // 0:未知 1:男 2:女
	BirthDate      string    `gorm:"type:varchar(20);column:birth_date" json:"birthDate"`
	Status         int       `gorm:"type:tinyint;default:0" json:"status"` // 0:正常 1:禁用 2:待审核
	IsOnline       int       `gorm:"type:tinyint;default:0;column:is_online" json:"isOnline"` // 0:离线 1:在线
	CurrentConsultationCount int `gorm:"type:int;default:0;column:current_consultation_count" json:"currentConsultationCount"` // 当前问诊数
	MaxConsultationCount int `gorm:"type:int;default:20;column:max_consultation_count" json:"maxConsultationCount"` // 最大同时问诊数
	DoctorCert     string    `gorm:"type:varchar(500);column:doctor_cert" json:"doctorCert"`
	DoctorTitle    string    `gorm:"type:varchar(50);column:doctor_title" json:"doctorTitle"`
	DoctorDept     string    `gorm:"type:varchar(50);column:doctor_dept" json:"doctorDept"`
	Specialty      string    `gorm:"type:text;column:specialty" json:"specialty"`
	Introduction   string    `gorm:"type:text;column:doctor_intro" json:"introduction"`
	CertNumber     string    `gorm:"type:varchar(100);column:cert_number" json:"certNumber"`
	CertStatus     string    `gorm:"type:varchar(20);column:cert_status" json:"certStatus"` // pending, approved, rejected
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	LastLoginTime  *time.Time `gorm:"column:last_login_time" json:"lastLoginTime"`
	LastLoginIP    string    `gorm:"type:varchar(255);column:last_login_ip" json:"-"` // SM4加密
}

func (User) TableName() string {
	return "SM_user"
}

// DoctorApplication 医生申请记录
type DoctorApplication struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"applicationId"`
	UserID        int64     `gorm:"not null;index" json:"userId"`
	ApplicationNo string    `gorm:"type:varchar(50);not null;unique;column:application_no" json:"applicationNo"` // 申请编号
	RealName      string    `gorm:"type:varchar(255);not null;column:real_name" json:"realName"`
	IDCard        string    `gorm:"type:varchar(255);not null;column:id_card" json:"idCard"`
	Phone         string    `gorm:"type:varchar(255);not null" json:"phone"`
	Email         string    `gorm:"type:varchar(255);not null" json:"email"` // 邮箱
	DoctorCert    string    `gorm:"type:varchar(500);not null;column:doctor_cert" json:"doctorCert"`
	DoctorTitle   string    `gorm:"type:varchar(50);not null;column:doctor_title" json:"doctorTitle"`
	DoctorDept    string    `gorm:"type:varchar(50);not null;column:doctor_dept" json:"doctorDept"`
	Specialty     string    `gorm:"type:text;column:specialty" json:"specialty"`
	Introduction  string    `gorm:"type:text;column:doctor_intro" json:"introduction"`
	CertNumber    string    `gorm:"type:varchar(100);column:cert_number" json:"certNumber"`
	PracticeCert  string    `gorm:"type:varchar(100);column:practice_cert" json:"practiceCert"` // 执业证号
	HospitalName  string    `gorm:"type:varchar(100);column:hospital_name" json:"hospitalName"` // 医院名称
	Status        int       `gorm:"type:tinyint;default:0;index" json:"status"` // 0:待审核 1:已通过 2:已拒绝
	StatusText    string    `gorm:"-" json:"statusText"`
	RejectReason  string    `gorm:"type:text;column:reject_reason" json:"rejectReason"`
	ReviewerID    *int64    `gorm:"column:reviewer_id" json:"reviewerId"`
	ReviewTime    *time.Time `gorm:"column:review_time" json:"reviewTime"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (DoctorApplication) TableName() string {
	return "SM_doctor_application"
}

// Consultation 问诊记录
type Consultation struct {
	ID                int64     `gorm:"primaryKey;autoIncrement;column:id" json:"consultationId"`
	PatientID         int64     `gorm:"not null;index;column:patient_id" json:"patientId"`
	DoctorID          *int64    `gorm:"index;column:doctor_id" json:"doctorId"`
	ConsultationNo    string    `gorm:"type:varchar(50);uniqueIndex;not null;column:consultation_no" json:"consultationNo"`
	ChiefComplaint    string    `gorm:"type:text;column:chief_complaint" json:"chiefComplaint"` // SM4加密
	SymptomsEncrypted string    `gorm:"type:text;column:symptoms_encrypted" json:"-"` // Paillier加密
	Symptoms          string    `gorm:"-" json:"symptoms"` // 解密后的症状
	AIRiskScore       *int      `gorm:"column:ai_risk_score" json:"aiRiskScore"`
	AIDiagnosis       string    `gorm:"type:text;column:ai_diagnosis" json:"aiDiagnosis"`
	AISuggestions     string    `gorm:"type:text;column:ai_suggestions" json:"aiSuggestions"`
	DoctorDiagnosis   string    `gorm:"type:text;column:doctor_diagnosis" json:"doctorDiagnosis"` // SM4加密
	Prescription      string    `gorm:"type:text" json:"prescription"` // SM4加密
	Status            int       `gorm:"type:tinyint;default:0;index" json:"status"` // 0:待接诊 1:问诊中 2:已完成 3:已取消
	StatusText        string    `gorm:"-" json:"statusText"`
	NeedAI            bool      `gorm:"type:tinyint;default:1;column:need_ai" json:"needAI"`
	AutoAssigned      bool      `gorm:"type:tinyint;default:0;column:auto_assigned" json:"autoAssigned"` // 是否自动分诊
	AssignedReason    string    `gorm:"type:varchar(200);column:assigned_reason" json:"assignedReason"` // 分诊原因
	RecommendedDept   string    `gorm:"type:varchar(100);column:recommended_dept" json:"recommendedDept"` // AI推荐科室
	CreatedAt         time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	CompletedAt       *time.Time `gorm:"column:completed_at" json:"completedAt"`
	
	// 关联字段
	PatientName       string    `gorm:"-" json:"patientName"`
	DoctorName        string    `gorm:"-" json:"doctorName"`
	Avatar            string    `gorm:"-" json:"avatar"`
}

func (Consultation) TableName() string {
	return "SM_consultation"
}

// MedicalRecord 电子病历
type MedicalRecord struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"recordId"`
	RecordNo        string    `gorm:"type:varchar(50);uniqueIndex;not null;column:record_no" json:"recordNo"`
	PatientID       int64     `gorm:"not null;index;column:patient_id" json:"patientId"`
	ConsultationID  *int64    `gorm:"index;column:consultation_id" json:"consultationId"`
	RecordType      int       `gorm:"type:tinyint;not null;column:record_type" json:"recordType"` // 1:门诊 2:在线问诊
	ChiefComplaint  string    `gorm:"type:text;column:chief_complaint" json:"chiefComplaint"` // SM4加密
	PresentIllness  string    `gorm:"type:text;column:present_illness" json:"presentIllness"` // SM4加密
	PastHistory     string    `gorm:"type:text;column:past_history" json:"pastHistory"` // SM4加密
	Diagnosis       string    `gorm:"type:text" json:"diagnosis"` // SM4加密
	Treatment       string    `gorm:"type:text;column:treatment_plan" json:"treatment"` // SM4加密
	DoctorID        *int64    `gorm:"column:doctor_id" json:"doctorId"`
	DoctorName      string    `gorm:"-" json:"doctorName"`
	DoctorDept      string    `gorm:"-" json:"doctorDept"`
	AIAdvice        string    `gorm:"type:text;column:ai_advice" json:"aiAdvice"`
	Symptoms        string    `gorm:"-" json:"symptoms"` // 前端需要的症状信息
	DataHash        string    `gorm:"type:varchar(128);column:data_hash" json:"hashValue"` // SM3哈希
	CreatedAt       time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (MedicalRecord) TableName() string {
	return "SM_medical_record"
}

// Notification 系统通知
type Notification struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement;column:id" json:"notificationId"`
	UserID             int64     `gorm:"not null;index;column:user_id" json:"userId"`
	NotificationType   string    `gorm:"type:varchar(20);not null;column:notification_type" json:"type"` // system, consultation, audit
	Title              string    `gorm:"type:varchar(100);not null" json:"title"`
	Content            string    `gorm:"type:text" json:"content"`
	RelatedID          *int64    `gorm:"column:related_id" json:"relatedId"`
	RelatedType        string    `gorm:"type:varchar(50);column:related_type" json:"relatedType"`
	IsRead             bool      `gorm:"type:tinyint;default:0;index;column:is_read" json:"isRead"`
	CreatedAt          time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
	ReadAt             *time.Time `gorm:"column:read_at" json:"readAt"`
}

func (Notification) TableName() string {
	return "SM_notification"
}

// LoginLog 登录日志
type LoginLog struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"logId"`
	UserID        *int64    `gorm:"column:user_id" json:"userId"`
	Username      string    `gorm:"type:varchar(50);column:username" json:"username"`
	LoginIP       string    `gorm:"type:varchar(255);not null;column:login_ip" json:"-"` // SM4加密
	LoginLocation string    `gorm:"type:varchar(100);column:login_location" json:"loginLocation"`
	Browser       string    `gorm:"type:varchar(50)" json:"browser"`
	OS            string    `gorm:"type:varchar(50);column:os" json:"os"`
	Status        int       `gorm:"type:tinyint;not null" json:"status"` // 0:失败 1:成功
	Msg           string    `gorm:"type:varchar(255)" json:"msg"`
	LoginTime     time.Time `gorm:"autoCreateTime;column:login_time" json:"loginTime"`
}

func (LoginLog) TableName() string {
	return "SM_login_log"
}

// Medicine 药品表模型
type Medicine struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement" json:"medicineId"`
	MedicineCode       string    `gorm:"type:varchar(50);uniqueIndex;not null;column:medicine_code" json:"medicineCode"`
	MedicineName       string    `gorm:"type:varchar(200);not null;column:medicine_name" json:"medicineName"`
	CommonName         string    `gorm:"type:varchar(200);column:common_name" json:"commonName"`
	MedicineType       string    `gorm:"type:varchar(50);not null;column:medicine_type" json:"medicineType"`
	Category           string    `gorm:"type:varchar(50)" json:"category"`
	Specification      string    `gorm:"type:varchar(100)" json:"specification"`
	DosageForm         string    `gorm:"type:varchar(50);column:dosage_form" json:"dosageForm"`
	Manufacturer       string    `gorm:"type:varchar(200)" json:"manufacturer"`
	Price              float64   `gorm:"type:decimal(10,2);default:0.00" json:"price"`
	Unit               string    `gorm:"type:varchar(20);default:'盒'" json:"unit"`
	Indications        string    `gorm:"type:text" json:"indications"`
	UsageDosage        string    `gorm:"type:text;column:usage_dosage" json:"usageDosage"`
	PrescriptionType   int       `gorm:"type:tinyint;default:0;column:prescription_type" json:"prescriptionType"`
	IsOTC              bool      `gorm:"type:tinyint;default:0;column:is_otc" json:"isOtc"`
	StockQuantity      int       `gorm:"type:int;default:0;column:stock_quantity" json:"stockQuantity"`
	Status             int       `gorm:"type:tinyint;default:1" json:"status"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Medicine) TableName() string {
	return "SM_medicine"
}

// Prescription 处方表模型
type Prescription struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"prescriptionId"`
	PrescriptionNo   string     `gorm:"type:varchar(50);uniqueIndex;not null;column:prescription_no" json:"prescriptionNo"`
	ConsultationID   int64      `gorm:"not null;index;column:consultation_id" json:"consultationId"`
	PatientID        int64      `gorm:"not null;index;column:patient_id" json:"patientId"`
	DoctorID         int64      `gorm:"not null;index;column:doctor_id" json:"doctorId"`
	Diagnosis        string     `gorm:"type:text" json:"diagnosis"` // SM4加密
	PrescriptionType int        `gorm:"type:tinyint;default:1;column:prescription_type" json:"prescriptionType"`
	TotalAmount      float64    `gorm:"type:decimal(10,2);default:0.00;column:total_amount" json:"totalAmount"`
	Status           int        `gorm:"type:tinyint;default:0;index" json:"status"`
	DataHash         string     `gorm:"type:varchar(128);column:data_hash" json:"dataHash"` // SM3哈希
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	
	// 关联字段
	Details          []PrescriptionDetail `gorm:"-" json:"details"`
}

func (Prescription) TableName() string {
	return "SM_prescription"
}

// PrescriptionDetail 处方明细表模型
type PrescriptionDetail struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"detailId"`
	PrescriptionID int64     `gorm:"not null;index;column:prescription_id" json:"prescriptionId"`
	MedicineID     int64     `gorm:"not null;index;column:medicine_id" json:"medicineId"`
	MedicineName   string    `gorm:"type:varchar(200);not null;column:medicine_name" json:"medicineName"`
	Specification  string    `gorm:"type:varchar(100)" json:"specification"`
	Quantity       int       `gorm:"not null" json:"quantity"`
	Unit           string    `gorm:"type:varchar(20);default:'盒'" json:"unit"`
	UnitPrice      float64   `gorm:"type:decimal(10,2);default:0.00;column:unit_price" json:"unitPrice"`
	TotalPrice     float64   `gorm:"type:decimal(10,2);default:0.00;column:total_price" json:"totalPrice"`
	Usage          string    `gorm:"type:varchar(100)" json:"usage"`
	Frequency      string    `gorm:"type:varchar(100)" json:"frequency"`
	Dosage         string    `gorm:"type:varchar(100)" json:"dosage"`
	Duration       string    `gorm:"type:varchar(50)" json:"duration"`
	Notes          string    `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (PrescriptionDetail) TableName() string {
	return "SM_prescription_detail"
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"messageId"`
	MessageNo      string     `gorm:"type:varchar(50);uniqueIndex;not null;column:message_no" json:"messageNo"`
	ConsultationID int64      `gorm:"not null;index;column:consultation_id" json:"consultationId"`
	SenderID       int64      `gorm:"not null;index;column:sender_id" json:"senderId"`
	ReceiverID     int64      `gorm:"not null;index;column:receiver_id" json:"receiverId"`
	MessageType    int        `gorm:"type:tinyint;not null;default:1;column:message_type" json:"messageType"` // 1:文本 2:图片 3:语音 4:处方 5:系统
	Content        string     `gorm:"type:text" json:"content"` // SM4加密
	FileURL        string     `gorm:"type:varchar(500);column:file_url" json:"fileUrl"`
	FileSize       int        `gorm:"type:int;column:file_size" json:"fileSize"`
	Duration       int        `gorm:"type:int" json:"duration"` // 语音时长
	ExtraData      string     `gorm:"type:text;column:extra_data" json:"extraData"` // JSON格式
	IsRead         bool       `gorm:"type:tinyint;default:0;index;column:is_read" json:"isRead"`
	ReadAt         *time.Time `gorm:"column:read_at" json:"readAt"`
	IsDeleted      bool       `gorm:"type:tinyint;default:0;column:is_deleted" json:"-"`
	CreatedAt      time.Time  `gorm:"autoCreateTime;index" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	
	// 关联字段
	SenderName     string     `gorm:"-" json:"senderName"`
	SenderAvatar   string     `gorm:"-" json:"senderAvatar"`
	SenderRole     string     `gorm:"-" json:"senderRole"`
}

func (ChatMessage) TableName() string {
	return "SM_chat_message"
}

// ChatUnreadCount 未读消息统计模型
type ChatUnreadCount struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64      `gorm:"not null;index;column:user_id" json:"userId"`
	ConsultationID  int64      `gorm:"not null;column:consultation_id" json:"consultationId"`
	UnreadCount     int        `gorm:"default:0;column:unread_count" json:"unreadCount"`
	LastMessageID   *int64     `gorm:"column:last_message_id" json:"lastMessageId"`
	LastMessageTime *time.Time `gorm:"column:last_message_time" json:"lastMessageTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (ChatUnreadCount) TableName() string {
	return "SM_chat_unread_count"
}
