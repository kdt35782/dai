// API配置
export const API_BASE_URL = 'http://localhost:3000' // 请根据实际后端地址修改

// WebSocket配置
export const WS_BASE_URL = 'ws://localhost:3000' // WebSocket地址

// API接口路径
export const API = {
	// 用户模块
	USER_REGISTER: '/api/user/register',
	USER_LOGIN: '/api/user/login',
	USER_INFO: '/api/user/info',
	USER_UPDATE: '/api/user/profile',
	USER_PASSWORD: '/api/user/password',
	USER_LOGOUT: '/api/user/logout',
	USER_APPLY_DOCTOR: '/api/user/apply-doctor',
	USER_DOCTOR_APPLICATION: '/api/user/doctor-application',
	USER_DOCTORS: '/api/user/doctors',
	USER_DOCTOR_DETAIL: '/api/user/doctor',
	
	// 管理员接口
	USER_ADMIN_APPLICATIONS: '/api/user/admin/doctor-applications',
	USER_ADMIN_REVIEW_DOCTOR: '/api/user/admin/review-doctor',
	USER_ADMIN_USERS: '/api/user/admin/users',
	USER_ADMIN_STATUS: '/api/user/admin/status',
	USER_ADMIN_LOGIN_LOGS: '/api/user/admin/login-logs',
	
	// 国密密钥管理
	CRYPTO_PUBLIC_KEY: '/api/crypto/public-key',
	CRYPTO_GENERATE_KEYPAIR: '/api/crypto/generate-keypair',
	CRYPTO_SESSION_KEY: '/api/crypto/session-key',
	
	// 文件上传
	FILE_UPLOAD: '/api/file/upload',
	FILE_DOWNLOAD: '/api/file/download',
	
	// 问诊模块
	CONSULTATION_CREATE: '/api/consultation/create',
	CONSULTATION_LIST: '/api/consultation/list',
	CONSULTATION_DETAIL: '/api/consultation/detail',
	CONSULTATION_ACCEPT: '/api/consultation/accept',
	CONSULTATION_FINISH: '/api/consultation/finish',
	CONSULTATION_MESSAGE: '/api/consultation/{id}/message',
	CONSULTATION_MESSAGES: '/api/consultation/{id}/messages',
	
	// 病历模块
	RECORD_LIST: '/api/record/list',
	RECORD_DETAIL: '/api/record/detail',
	MEDICAL_RECORD_CREATE: '/api/medical-record/create',
	MEDICAL_RECORD_DETAIL: '/api/medical-record',
	MEDICAL_RECORD_PATIENT_LIST: '/api/medical-record/patient',
	MEDICAL_RECORD_AUTHORIZE: '/api/medical-record/{id}/authorize',
	MEDICAL_RECORD_ACCESS_LOGS: '/api/medical-record/{id}/access-logs',
	
	// 处方模块
	PRESCRIPTION_SEARCH_MEDICINES: '/api/prescription/medicines/search',
	PRESCRIPTION_RECOMMEND: '/api/prescription/medicines/recommend',
	PRESCRIPTION_DETAIL: '/api/prescription',
	
	// 聊天模块
	CHAT_WS: '/api/chat/ws',
	CHAT_SEND: '/api/chat/send',
	CHAT_MESSAGES: '/api/chat/messages',
	CHAT_UNREAD_COUNT: '/api/chat/unread-count',
	CHAT_MARK_READ: '/api/chat/mark-read',
	CHAT_ONLINE_STATUS: '/api/chat/online-status',
	CHAT_TYPING: '/api/chat/typing',
	
	// 通知模块
	NOTIFICATION_LIST: '/api/notification/list',
	NOTIFICATION_MARK_READ: '/api/notification/mark-read',
	NOTIFICATION_UNREAD_COUNT: '/api/notification/unread-count'
}

// 本地存储键名
export const STORAGE_KEYS = {
	TOKEN: 'token',
	USER_INFO: 'userInfo',
	SM2_PUBLIC_KEY: 'sm2PublicKey',
	PAILLIER_PUBLIC_KEY: 'paillierPublicKey'
}
