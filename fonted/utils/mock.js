/**
 * Mock数据模拟系统
 * 用于前端开发测试，无需后端即可查看页面效果
 */

import { STORAGE_KEYS } from './config.js'

// 模拟用户数据
export const mockUsers = {
	// 普通用户（患者）
	patient: {
		userId: 1001,
		username: 'testuser',
		realName: '张三',
		role: 'patient',
		gender: 1,
		birthDate: '1990-05-15',
		phone: '13800138000',
		email: 'zhangsan@example.com',
		avatar: 'https://via.placeholder.com/150',
		certStatus: null
	},
	
	// 医生用户
	doctor: {
		userId: 2001,
		username: 'doctor_li',
		realName: '李医生',
		role: 'doctor',
		gender: 1,
		age: 45,
		birthDate: '1980-03-20',
		phone: '13900139000',
		email: 'doctor.li@example.com',
		avatar: 'https://via.placeholder.com/150',
		doctorTitle: '主任医师',
		doctorDept: '内科',
		specialty: '心血管疾病、高血压、糖尿病的诊治',
		doctorIntro: '从事临床工作20余年，擅长心血管疾病的诊断和治疗。曾在多家三甲医院工作，具有丰富的临床经验。多次参加国内外学术交流，发表SCI论朇10余篇。',
		workYears: 20,
		doctorCert: 'https://via.placeholder.com/600x800?text=Doctor+Certificate',
		certNumber: 'DOC2020001',
		certStatus: 'approved',
		consultationCount: 328,
		rating: '4.9'
	},
	
	// 管理员用户
	admin: {
		userId: 9001,
		username: 'admin',
		realName: '系统管理员',
		role: 'admin',
		gender: 1,
		age: 35,
		birthDate: '1989-01-01',
		phone: '13700137000',
		email: 'admin@example.com',
		avatar: 'https://via.placeholder.com/150',
		certStatus: null
	}
}

// 模拟医生列表
export const mockDoctors = [
	{
		userId: 2001,
		realName: '李医生',
		doctorTitle: '主任医师',
		doctorDept: '内科',
		specialty: '心血管疾病、高血压、糖尿病',
		avatar: 'https://via.placeholder.com/150',
		consultationCount: 328,
		rating: '4.9'
	},
	{
		userId: 2002,
		realName: '王医生',
		doctorTitle: '副主任医师',
		doctorDept: '外科',
		specialty: '普通外科、微创手术',
		avatar: 'https://via.placeholder.com/150',
		consultationCount: 215,
		rating: '4.8'
	},
	{
		userId: 2003,
		realName: '刘医生',
		doctorTitle: '主治医师',
		doctorDept: '儿科',
		specialty: '儿童常见病、呼吸系统疾病',
		avatar: 'https://via.placeholder.com/150',
		consultationCount: 156,
		rating: '4.7'
	},
	{
		userId: 2004,
		realName: '陈医生',
		doctorTitle: '主任医师',
		doctorDept: '妇产科',
		specialty: '妇科疾病、产前检查',
		avatar: 'https://via.placeholder.com/150',
		consultationCount: 289,
		rating: '4.9'
	},
	{
		userId: 2005,
		realName: '赵医生',
		doctorTitle: '主治医师',
		doctorDept: '骨科',
		specialty: '骨折、关节疾病',
		avatar: 'https://via.placeholder.com/150',
		consultationCount: 167,
		rating: '4.6'
	}
]

// 模拟问诊列表
export const mockConsultations = [
	{
		consultationId: 3001,
		doctorId: 2001,
		doctorName: '李医生',
		doctorDept: '内科',
		patientId: 1001,
		patientName: '张三',
		avatar: 'https://via.placeholder.com/150',
		chiefComplaint: '头晕头痛3天，血压偏高',
		symptoms: {
			age: 35,
			gender: 1,
			bloodPressure: '150/95',
			heartRate: 88,
			otherSymptoms: '最近工作压力大，睡眠质量不好'
		},
		status: 2,
		statusText: '已完成',
		needAI: true,
		createdAt: '2024-12-01 10:30',
		updatedAt: '2024-12-01 12:45'
	},
	{
		consultationId: 3002,
		doctorId: 2003,
		doctorName: '刘医生',
		doctorDept: '儿科',
		patientId: 1001,
		patientName: '张三',
		avatar: 'https://via.placeholder.com/150',
		chiefComplaint: '孩子发热咳嗽2天',
		symptoms: {
			age: 5,
			gender: 0,
			temperature: '38.5',
			otherSymptoms: '有轻微咳嗽，食欲不振'
		},
		status: 1,
		statusText: '问诊中',
		needAI: false,
		createdAt: '2024-12-03 09:15',
		updatedAt: '2024-12-03 09:20'
	},
	{
		consultationId: 3003,
		doctorId: 2002,
		doctorName: '王医生',
		doctorDept: '外科',
		patientId: 1001,
		patientName: '张三',
		avatar: 'https://via.placeholder.com/150',
		chiefComplaint: '右腿膝盖疼痛1周',
		symptoms: {
			age: 35,
			gender: 1,
			otherSymptoms: '上楼梯时疼痛明显，休息后缓解'
		},
		status: 0,
		statusText: '待接诊',
		needAI: false,
		createdAt: '2024-12-03 14:20',
		updatedAt: '2024-12-03 14:20'
	}
]

// 模拟病历列表
export const mockRecords = [
	{
		recordId: 4001,
		consultationId: 3001,
		doctorId: 2001,
		doctorName: '李医生',
		doctorDept: '内科',
		patientId: 1001,
		chiefComplaint: '头晕头痛3天，血压偏高',
		symptoms: {
			age: 35,
			gender: 1,
			bloodPressure: '150/95',
			heartRate: 88,
			otherSymptoms: '最近工作压力大，睡眠质量不好'
		},
		diagnosis: '高血压（1级）',
		treatment: '1. 注意休息，避免劳累\n2. 低盐低脂饮食\n3. 规律作息，保证睡眠\n4. 建议定期监测血压\n5. 必要时药物治疗',
		aiAdvice: 'AI分析：患者血压偏高，建议进一步检查心血管功能，排除器质性病变。生活方式调整是首要措施。',
		createdAt: '2024-12-01 12:45',
		hashValue: 'a1b2c3d4e5f6...'
	},
	{
		recordId: 4002,
		consultationId: 3002,
		doctorId: 2003,
		doctorName: '刘医生',
		doctorDept: '儿科',
		patientId: 1001,
		chiefComplaint: '孩子发热咳嗽2天',
		symptoms: {
			age: 5,
			gender: 0,
			temperature: '38.5',
			otherSymptoms: '有轻微咳嗽，食欲不振'
		},
		diagnosis: '急性上呼吸道感染',
		treatment: '1. 多喝水，注意休息\n2. 物理降温\n3. 清淡饮食\n4. 观察病情变化',
		aiAdvice: null,
		createdAt: '2024-12-02 16:30',
		hashValue: 'f6e5d4c3b2a1...'
	}
]

// 模拟消息通知
export const mockNotifications = [
	{
		notificationId: 5001,
		type: 'consultation',
		title: '问诊消息',
		content: '李医生已回复您的问诊，请查看',
		relatedId: 3002,
		isRead: false,
		createdAt: '2024-12-03 10:30'
	},
	{
		notificationId: 5002,
		type: 'system',
		title: '系统通知',
		content: '您的病历已生成，可在病历列表中查看',
		relatedId: null,
		isRead: false,
		createdAt: '2024-12-02 16:35'
	},
	{
		notificationId: 5003,
		type: 'system',
		title: '系统通知',
		content: '欢迎使用国密问诊平台',
		relatedId: null,
		isRead: true,
		createdAt: '2024-12-01 09:00'
	}
]

/**
 * Mock API响应
 */
export function mockApiResponse(url, method, data) {
	console.log('[Mock API]', method, url, data)
	
	// 统一返回格式
	const success = (responseData) => ({
		code: 200,
		message: '操作成功',
		data: responseData,
		timestamp: Date.now()
	})
	
	const error = (message) => ({
		code: 400,
		message: message,
		data: null,
		timestamp: Date.now()
	})
	
	// 登录
	if (url.includes('/api/user/login')) {
		// 根据用户名返回不同角色的用户
		let user
		if (data.username === 'doctor' || data.username === 'doctor_li' || data.username === 'xlf_20') {
			// 医生账号
			if (data.username === 'xlf_20') {
				// xlf_20 专属医生账号
				user = {
					userId: 2020,
					username: 'xlf_20',
					realName: 'xlf医生',
					role: 'doctor',
					gender: 1,
					age: 35,
					birthDate: '1989-05-20',
					phone: '13900139020',
					email: 'xlf_20@example.com',
					avatar: 'https://via.placeholder.com/150',
					doctorTitle: '副主任医师',
					doctorDept: '外科',
					specialty: '骨科疾病、运动损伤',
					doctorIntro: '从事骨科临床工作10余年，擅长骨折、关节疾病的诊治。',
					workYears: 10,
					doctorCert: 'https://via.placeholder.com/600x800?text=Doctor+Certificate',
					certNumber: 'DOC2020020',
					certStatus: 'approved',
					consultationCount: 156,
					rating: '4.8'
				}
			} else {
				user = mockUsers.doctor
			}
		} else if (data.username === 'admin') {
			user = mockUsers.admin
		} else {
			// 为其他用户名动态生成用户信息
			user = {
				userId: 1000 + Math.floor(Math.random() * 1000),
				username: data.username,
				realName: data.username, // 使用用户名作为真实姓名
				role: 'user',
				gender: 1,
				age: 30,
				birthDate: '1994-01-01',
				phone: '138' + String(Math.floor(Math.random() * 100000000)).padStart(8, '0'),
				email: data.username + '@example.com',
				avatar: 'https://via.placeholder.com/150',
				certStatus: null
			}
		}
		
		console.log('👤 Mock 登录成功:', user.realName, '(角色:', user.role + ')')
		
		return success({
			token: 'mock_token_' + user.role + '_' + Date.now(),
			userInfo: user
		})
	}
	
	// 注册
	if (url.includes('/api/user/register')) {
		return success({
			token: 'mock_token_' + Date.now(),
			userInfo: mockUsers.patient
		})
	}
	
	// 获取用户信息
	if (url.includes('/api/user/info')) {
		const currentUser = uni.getStorageSync(STORAGE_KEYS.USER_INFO) || mockUsers.patient
		
		// 如果有userId参数，返回对应医生信息
		if (data && data.userId) {
			const doctor = mockDoctors.find(d => d.userId === parseInt(data.userId))
			return success(doctor || mockUsers.doctor)
		}
		
		return success(currentUser)
	}
	
	// 更新用户资料
	if ((url.includes('/api/user/info') || url.includes('/api/user/profile')) && method === 'PUT') {
		const currentUser = uni.getStorageSync(STORAGE_KEYS.USER_INFO)
		const updatedUser = { ...currentUser, ...data }
		uni.setStorageSync(STORAGE_KEYS.USER_INFO, updatedUser)
		return success(updatedUser)
	}
	
	// 申请成为医生
	if (url.includes('/api/user/apply-doctor')) {
		return success({ message: '申请已提交，等待审核' })
	}
	
	// 获取医生列表
	if (url.includes('/api/user/doctors')) {
		let list = [...mockDoctors]
		
		// 科室筛选
		if (data && data.dept) {
			list = list.filter(d => d.doctorDept === data.dept)
		}
		
		// 关键词搜索
		if (data && data.keyword) {
			list = list.filter(d => 
				d.realName.includes(data.keyword) || 
				d.specialty.includes(data.keyword)
			)
		}
		
		return success({
			list: list,
			total: list.length,
			page: data?.page || 1,
			pageSize: data?.pageSize || 10
		})
	}
	
	// 创建问诊
	if (url.includes('/api/consultation/create')) {
		const newConsultation = {
			consultationId: 3000 + Math.floor(Math.random() * 1000),
			...data,
			status: 0,
			statusText: '待接诊',
			createdAt: new Date().toLocaleString('zh-CN'),
			updatedAt: new Date().toLocaleString('zh-CN')
		}
		return success(newConsultation)
	}
	
	// 获取问诊列表
	if (url.includes('/api/consultation/list')) {
		let list = [...mockConsultations]
		
		// 状态筛选
		if (data && data.status !== undefined && data.status !== '') {
			list = list.filter(c => c.status === data.status)
		}
		
		return success({
			list: list,
			total: list.length,
			page: data?.page || 1,
			pageSize: data?.pageSize || 10
		})
	}
	
	// 获取问诊详情
	if (url.includes('/api/consultation/detail')) {
		const consultation = mockConsultations.find(c => c.consultationId === data.consultationId)
		return success(consultation || mockConsultations[0])
	}
	
	// 接诊
	if (url.includes('/api/consultation/accept')) {
		return success({ message: '接诊成功' })
	}
	
	// 完成问诊
	if (url.includes('/api/consultation/finish')) {
		return success({ message: '问诊已完成' })
	}
	
	// 获取病历列表
	if (url.includes('/api/record/list')) {
		let list = [...mockRecords]
		
		// 日期筛选
		if (data && data.startDate) {
			// 简单模拟，实际需要日期比较
			list = list.filter(r => r.createdAt >= data.startDate)
		}
		
		return success({
			list: list,
			total: list.length,
			page: data?.page || 1,
			pageSize: data?.pageSize || 10
		})
	}
	
	// 获取病历详情
	if (url.includes('/api/record/detail')) {
		const record = mockRecords.find(r => r.recordId === data.recordId)
		return success(record || mockRecords[0])
	}
	
	// 获取消息列表
	if (url.includes('/api/notification/list')) {
		let list = [...mockNotifications]
		
		// 类型筛选
		if (data && data.type) {
			list = list.filter(n => n.type === data.type)
		}
		
		return success({
			list: list,
			total: list.length,
			page: data?.page || 1,
			pageSize: data?.pageSize || 10
		})
	}
	
	// 获取未读消息数
	if (url.includes('/api/notification/unread-count')) {
		const unreadCount = mockNotifications.filter(n => !n.isRead).length
		return success({
			totalUnread: unreadCount,
			systemUnread: mockNotifications.filter(n => !n.isRead && n.type === 'system').length,
			consultationUnread: mockNotifications.filter(n => !n.isRead && n.type === 'consultation').length
		})
	}
	
	// 标记已读
	if (url.includes('/api/notification/mark-read')) {
		return success({ message: '已标记为已读' })
	}
	
	// 生成密钥
	if (url.includes('/api/key/generate')) {
		return success({
			publicKey: 'mock_sm2_public_key_' + Date.now(),
			message: '密钥已生成'
		})
	}
	
	// 文件上传
	if (url.includes('/api/file/upload') || url.includes('/api/upload')) {
		// 根据文件类型返回不同的 Mock URL
		const fileType = data.fileType || 'image'
		let mockUrl = ''
		
		if (fileType === 'avatar') {
			mockUrl = 'https://via.placeholder.com/150?text=Avatar'
		} else if (fileType === 'cert') {
			mockUrl = 'https://via.placeholder.com/600x800?text=Certificate'
		} else {
			mockUrl = 'https://via.placeholder.com/400?text=Uploaded+Image'
		}
		
		return success({
			fileUrl: mockUrl,
			fileName: 'mock_file_' + Date.now() + '.png'
		})
	}
	
	// 默认返回成功
	return success({ message: '操作成功' })
}

/**
 * 启用Mock模式
 */
export function enableMockMode() {
	console.log('✅ Mock模式已启用 - 无需后端即可测试')
	
	// 自动登录模拟用户（默认普通用户）
	uni.setStorageSync(STORAGE_KEYS.TOKEN, 'mock_token_' + Date.now())
	uni.setStorageSync(STORAGE_KEYS.USER_INFO, mockUsers.patient)
	uni.setStorageSync(STORAGE_KEYS.SM2_PUBLIC_KEY, 'mock_sm2_public_key')
	
	console.log('👤 已自动登录模拟用户:', mockUsers.patient.username)
	console.log('💡 提示:')
	console.log('  - 普通用户: 用户名输入任意名称（如 testuser、lisi、wangwu）')
	console.log('  - 医生角色: 用户名输入 "doctor" 或 "doctor_li"')
	console.log('  - 管理员: 用户名输入 "admin"')
	console.log('  - 每个不同的用户名会生成不同的账号信息')
}

/**
 * 快速切换为医生角色（测试用）
 */
export function switchToDoctorMode() {
	console.log('👨‍⚕️ 切换到医生模式')
	uni.setStorageSync(STORAGE_KEYS.TOKEN, 'mock_token_doctor_' + Date.now())
	uni.setStorageSync(STORAGE_KEYS.USER_INFO, mockUsers.doctor)
	console.log('✅ 已切换为医生角色:', mockUsers.doctor.realName)
	return mockUsers.doctor
}

/**
 * 快速切换为管理员角色（测试用）
 */
export function switchToAdminMode() {
	console.log('👨‍💼 切换到管理员模式')
	uni.setStorageSync(STORAGE_KEYS.TOKEN, 'mock_token_admin_' + Date.now())
	uni.setStorageSync(STORAGE_KEYS.USER_INFO, mockUsers.admin)
	console.log('✅ 已切换为管理员角色:', mockUsers.admin.realName)
	return mockUsers.admin
}

/**
 * 快速切换为普通用户（测试用）
 */
export function switchToPatientMode() {
	console.log('👤 切换到普通用户模式')
	uni.setStorageSync(STORAGE_KEYS.TOKEN, 'mock_token_patient_' + Date.now())
	uni.setStorageSync(STORAGE_KEYS.USER_INFO, mockUsers.patient)
	console.log('✅ 已切换为普通用户:', mockUsers.patient.realName)
	return mockUsers.patient
}

/**
 * 检查是否启用Mock模式
 * 前后端联调时设置为 false
 */
export function isMockEnabled() {
	// 前后端联调：关闭Mock模式，连接真实后端
	// 开发调试：设置为 true 使用Mock数据
	return false  // ❌ 已关闭 Mock 模式，使用真实后端
}
