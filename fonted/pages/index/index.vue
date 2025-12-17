<template>
	<view class="index-page">
		<!-- 顶部状态栏 -->
		<view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
		
		<!-- 顶部导航 -->
		<view class="header">
			<view class="header-content">
				<text class="logo">🏥</text>
				<view class="header-text">
					<text class="title">网上医疗问诊平台</text>
					<text class="subtitle">安全·便捷·智能</text>
				</view>
				<view class="notification" @click="goNotification">
					<text class="icon">🔔</text>
					<view class="badge" v-if="unreadCount > 0">{{ unreadCount }}</view>
				</view>
			</view>
		</view>
		
		<!-- 快捷入口 - 患者端 -->
		<view class="quick-entry" v-if="(!isLogin || userInfo?.role === 'user')">
			<view class="entry-item" @click="goConsultation">
				<view class="entry-icon">💬</view>
				<text class="entry-text">在线问诊</text>
			</view>
			<view class="entry-item" @click="goDoctors">
				<view class="entry-icon">👨‍⚕️</view>
				<text class="entry-text">找医生</text>
			</view>
			<view class="entry-item" @click="goRecords">
				<view class="entry-icon">📋</view>
				<text class="entry-text">我的病历</text>
			</view>
			<view class="entry-item" @click="goApplyDoctor">
				<view class="entry-icon">🎓</view>
				<text class="entry-text">医生入驻</text>
			</view>
		</view>
		
		<!-- 快捷入口 - 医生端 -->
		<view class="quick-entry" v-else-if="userInfo?.role === 'doctor'">
			<view class="entry-item" @click="goConsultationList">
				<view class="entry-icon doctor-icon">📋</view>
				<text class="entry-text">接诊管理</text>
			</view>
			<view class="entry-item" @click="goRecords">
				<view class="entry-icon doctor-icon">📄</view>
				<text class="entry-text">患者病历</text>
			</view>
			<view class="entry-item" @click="goNotification">
				<view class="entry-icon doctor-icon">🔔</view>
				<text class="entry-text">消息通知</text>
			</view>
			<view class="entry-item" @click="goDoctorProfile">
				<view class="entry-icon doctor-icon">⚕️</view>
				<text class="entry-text">我的资料</text>
			</view>
		</view>
		
		<!-- 快捷入口 - 管理员端 -->
		<view class="quick-entry" v-else-if="userInfo?.role === 'admin'">
			<view class="entry-item" @click="goDoctorApplications">
				<view class="entry-icon admin-icon">📋</view>
				<text class="entry-text">医生审核</text>
			</view>
			<view class="entry-item" @click="goUserManagement">
				<view class="entry-icon admin-icon">👥</view>
				<text class="entry-text">用户管理</text>
			</view>
			<view class="entry-item" @click="goSystemSettings">
				<view class="entry-icon admin-icon">⚙️</view>
				<text class="entry-text">系统设置</text>
			</view>
			<view class="entry-item" @click="goDataStatistics">
				<view class="entry-icon admin-icon">📈</view>
				<text class="entry-text">数据统计</text>
			</view>
		</view>
		
		<!-- 推荐医生 - 仅患者和未登录用户可见 -->
		<view class="section" v-if="(!isLogin || userInfo?.role === 'user') && recommendDoctors.length > 0">
			<view class="section-header">
				<text class="section-title">推荐医生</text>
				<text class="section-more" @click="goDoctors">更多 ></text>
			</view>
			
			<scroll-view class="doctor-scroll" scroll-x>
				<view class="doctor-list">
					<view 
						class="doctor-card" 
						v-for="doctor in recommendDoctors" 
						:key="doctor.userId"
						@click="goDoctorDetail(doctor.userId)"
					>
						<image class="doctor-avatar" :src="doctor.avatar || '/static/default-avatar.png'" mode="aspectFill"></image>
						<view class="doctor-info">
							<text class="doctor-name">{{ doctor.realName }}</text>
							<text class="doctor-title">{{ doctor.doctorTitle }}</text>
							<text class="doctor-dept">{{ doctor.doctorDept }}</text>
						</view>
					</view>
				</view>
			</scroll-view>
		</view>
		
		<!-- 最近问诊 - 患者端 -->
		<view class="section" v-if="isLogin && userInfo?.role !== 'doctor' && recentConsultations.length > 0">
			<view class="section-header">
				<text class="section-title">最近问诊</text>
				<text class="section-more" @click="goConsultationList">查看全部 ></text>
			</view>
			
			<view class="consultation-list">
				<view 
					class="consultation-item" 
					v-for="item in recentConsultations" 
					:key="item.consultationId"
					@click="goConsultationDetail(item.consultationId)"
				>
					<view class="consultation-header">
						<text class="doctor-name">{{ item.doctorName }} 医生</text>
						<text class="status" :class="'status-' + item.status">{{ item.statusText }}</text>
					</view>
					<text class="complaint">{{ item.chiefComplaint }}</text>
					<text class="time">{{ item.createdAt }}</text>
				</view>
			</view>
		</view>
		
		<!-- 待接诊列表 - 医生端 -->
		<view class="section" v-if="isLogin && userInfo?.role === 'doctor' && pendingConsultations.length > 0">
			<view class="section-header">
				<text class="section-title">待接诊患者</text>
				<text class="section-more" @click="goConsultationList">查看全部 ></text>
			</view>
			
			<view class="consultation-list">
				<view 
					class="consultation-item" 
					v-for="item in pendingConsultations" 
					:key="item.consultationId"
					@click="goConsultationDetail(item.consultationId)"
				>
					<view class="consultation-header">
						<text class="doctor-name">{{ item.patientName }} 患者</text>
						<text class="status status-urgent">待接诊</text>
					</view>
					<text class="complaint">主诉：{{ item.chiefComplaint }}</text>
					<text class="time">{{ item.createdAt }}</text>
				</view>
			</view>
		</view>
		
		<!-- 待审核申请列表 - 管理员端 -->
		<view class="section" v-if="isLogin && userInfo?.role === 'admin' && pendingApplications.length > 0">
			<view class="section-header">
				<text class="section-title">待审核医生申请</text>
				<text class="section-more" @click="goDoctorApplications">查看全部 ></text>
			</view>
			
			<view class="application-list">
				<view 
					class="application-item" 
					v-for="item in pendingApplications" 
					:key="item.applicationId"
					@click="goApplicationDetail(item.applicationId)"
				>
					<view class="application-header">
						<view class="applicant-info">
							<text class="applicant-name">{{ item.realName }}</text>
							<text class="applicant-title">{{ item.doctorTitle }} | {{ item.doctorDept }}</text>
						</view>
						<text class="status status-pending">待审核</text>
					</view>
					<text class="application-time">申请时间：{{ item.createdAt }}</text>
				</view>
			</view>
		</view>
		
		<!-- 系统统计数据 - 管理员端 -->
		<view class="section" v-if="isLogin && userInfo?.role === 'admin'">
			<view class="section-header">
				<text class="section-title">系统概览</text>
			</view>
			
			<view class="stats-grid">
				<view class="stat-card">
					<text class="stat-number">{{ systemStats.userCount || 0 }}</text>
					<text class="stat-label">总用户数</text>
				</view>
				<view class="stat-card">
					<text class="stat-number">{{ systemStats.doctorCount || 0 }}</text>
					<text class="stat-label">在职医生</text>
				</view>
				<view class="stat-card">
					<text class="stat-number">{{ systemStats.consultationCount || 0 }}</text>
					<text class="stat-label">总问诊数</text>
				</view>
				<view class="stat-card">
					<text class="stat-number">{{ systemStats.pendingCount || 0 }}</text>
					<text class="stat-label">待审核</text>
				</view>
			</view>
		</view>
		
		<!-- 安全提示 -->
		<view class="security-banner">
			<text class="banner-icon">🔐</text>
			<view class="banner-content">
				<text class="banner-title">国密算法保护</text>
				<text class="banner-desc">采用SM2/SM3/SM4国密算法，全程加密保护您的隐私</text>
			</view>
		</view>
	</view>
</template>

<script>
import { get } from '@/utils/request.js'
import { API, STORAGE_KEYS } from '@/utils/config.js'

export default {
		data() {
			return {
				statusBarHeight: 0,
				isLogin: false,
				userInfo: null,
				unreadCount: 0,
				recommendDoctors: [],
				recentConsultations: [],
				pendingConsultations: [],  // 医生端待接诊列表
				pendingApplications: [],   // 管理员端待审核申请
				systemStats: {             // 管理员端系统统计
					userCount: 0,
					doctorCount: 0,
					consultationCount: 0,
					pendingCount: 0
				}
			}
		},
	
	onLoad() {
		// 获取状态栏高度
		this.statusBarHeight = uni.getSystemInfoSync().statusBarHeight
		
		// 检查登录状态
		this.checkLogin()
		
		// 加载数据
		// 根据角色加载不同的内容
		if (!this.isLogin || this.userInfo?.role === 'user') {
			// 患者端：加载推荐医生
			this.loadRecommendDoctors()
		}
		
		if (this.isLogin) {
			this.loadUnreadCount()
			// 根据角色加载不同的问诊数据
			if (this.userInfo?.role === 'doctor') {
				// 医生端：加载待接诊列表
				this.loadPendingConsultations()
			} else if (this.userInfo?.role === 'admin') {
				// 管理员端：加载待审核申请和系统统计
				this.loadPendingApplications()
				this.loadSystemStats()
			} else {
				// 患者端：加载最近问诊
				this.loadRecentConsultations()
			}
		}
	},
	
	onShow() {
		// 每次显示时刷新未读消息
		if (this.isLogin) {
			this.loadUnreadCount()
		}
	},
	
	methods: {
		// 检查登录状态
		checkLogin() {
			const token = uni.getStorageSync(STORAGE_KEYS.TOKEN)
			const userInfo = uni.getStorageSync(STORAGE_KEYS.USER_INFO)
			
			this.isLogin = !!token
			this.userInfo = userInfo
		},
		
		// 加载推荐医生
		async loadRecommendDoctors() {
			try {
				const res = await get(API.USER_DOCTORS, {
					page: 1,
					pageSize: 6
				}, { noAuth: true })
				
				this.recommendDoctors = res.data.list || []
			} catch (error) {
				console.error('加载推荐医生失败:', error)
			}
		},
		
		// 加载未读消息数
		async loadUnreadCount() {
			try {
				const res = await get(API.NOTIFICATION_UNREAD_COUNT)
				this.unreadCount = res.data.totalUnread || 0
			} catch (error) {
				console.error('加载未读消息失败:', error)
			}
		},
		
		// 加载最近问诊（患者端）
		async loadRecentConsultations() {
			try {
				const res = await get(API.CONSULTATION_LIST, {
					page: 1,
					pageSize: 3,
					role: 'patient'
				})
				
				this.recentConsultations = res.data.list || []
			} catch (error) {
				console.error('加载最近问诊失败:', error)
			}
		},
		
		// 加载待接诊列表（医生端）
		async loadPendingConsultations() {
			try {
				const res = await get(API.CONSULTATION_LIST, {
					page: 1,
					pageSize: 5,
					role: 'doctor',
					status: 0  // 待接诊
				})
				
				this.pendingConsultations = res.data.list || []
			} catch (error) {
				console.error('加载待接诊列表失败:', error)
			}
		},
		
		// 加载待审核申请（管理员端）
		async loadPendingApplications() {
			try {
				const res = await get(API.USER_ADMIN_APPLICATIONS, {
					page: 1,
					pageSize: 5,
					status: 0  // 待审核
				})
				
				this.pendingApplications = res.data.list || []
			} catch (error) {
				console.error('加载待审核申请失败:', error)
			}
		},
		
		// 加载系统统计数据（管理员端）
		async loadSystemStats() {
			try {
				// 获取用户统计
				const userRes = await get(API.USER_ADMIN_USERS, {
					page: 1,
					pageSize: 1
				})
				this.systemStats.userCount = userRes.data.total || 0
				
				// 获取医生统计
				const doctorRes = await get(API.USER_DOCTORS, {
					page: 1,
					pageSize: 1
				}, { noAuth: true })
				this.systemStats.doctorCount = doctorRes.data.total || 0
				
				// 获取问诊统计
				const consultationRes = await get(API.CONSULTATION_LIST, {
					page: 1,
					pageSize: 1
				})
				this.systemStats.consultationCount = consultationRes.data.total || 0
				
				// 获取待审核统计
				const pendingRes = await get(API.USER_ADMIN_APPLICATIONS, {
					page: 1,
					pageSize: 1,
					status: 0
				})
				this.systemStats.pendingCount = pendingRes.data.total || 0
			} catch (error) {
				console.error('加载系统统计失败:', error)
			}
		},
		
		// 跳转到通知页面
		goNotification() {
			if (!this.isLogin) {
				this.goLogin()
				return
			}
			uni.navigateTo({
				url: '/pages/notification/notification'
			})
		},
		
		// 跳转到问诊
		goConsultation() {
			if (!this.isLogin) {
				this.goLogin()
				return
			}
			uni.navigateTo({
				url: '/pages/consultation/create-consultation'
			})
		},
		
		// 跳转到医生列表
		goDoctors() {
			uni.switchTab({
				url: '/pages/doctors/doctors'
			})
		},
		
		// 跳转到病历
		goRecords() {
			if (!this.isLogin) {
				this.goLogin()
				return
			}
			uni.navigateTo({
				url: '/pages/medical-record/record-list'
			})
		},
		
		// 跳转到医生申请
		goApplyDoctor() {
			if (!this.isLogin) {
				this.goLogin()
				return
			}
			uni.navigateTo({
				url: '/pages/user/apply-doctor'
			})
		},
		
		// 跳转到医生详情
		goDoctorDetail(userId) {
			uni.navigateTo({
				url: '/pages/doctor-detail/doctor-detail?userId=' + userId
			})
		},
		
		// 跳转到问诊列表
		goConsultationList() {
			uni.switchTab({
				url: '/pages/consultation/consultation-list'
			})
		},
		
		// 跳转到问诊详情
		goConsultationDetail(id) {
			uni.navigateTo({
				url: '/pages/consultation/consultation-detail?id=' + id
			})
		},
		
		// 跳转到登录
		goLogin() {
			uni.navigateTo({
				url: '/pages/login/login'
			})
		},
		
		// 跳转到医生资料页面
		goDoctorProfile() {
			uni.switchTab({
				url: '/pages/user/user'
			})
		},
		
		// 跳转到医生申请审核（管理员）
		goDoctorApplications() {
			if (!this.isLogin) {
				this.goLogin()
				return
			}
			uni.navigateTo({
				url: '/pages/admin/doctor-applications'
			})
		},
		
		// 跳转到申请详情（管理员）
		goApplicationDetail(applicationId) {
			uni.navigateTo({
				url: '/pages/admin/doctor-applications'
			})
		},
		
		// 跳转到用户管理（管理员）
		goUserManagement() {
			if (!this.isLogin) {
				this.goLogin()
				return
			}
			uni.navigateTo({
				url: '/pages/admin/user-management'
			})
		},
		
		// 跳转到系统设置（管理员）
		goSystemSettings() {
			if (!this.isLogin) {
				this.goLogin()
				return
			}
			uni.navigateTo({
				url: '/pages/admin/system-settings'
			})
		},
		
		// 跳转到数据统计（管理员）
		goDataStatistics() {
			if (!this.isLogin) {
				this.goLogin()
				return
			}
			uni.navigateTo({
				url: '/pages/admin/data-statistics'
			})
		}
	}
}
</script>

<style scoped>
.index-page {
	min-height: 100vh;
	background: #f5f5f5;
}

.status-bar {
	background: #07c160;
}

.header {
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	padding: 20rpx 30rpx 30rpx;
}

.header-content {
	display: flex;
	align-items: center;
}

.logo {
	font-size: 60rpx;
	margin-right: 20rpx;
}

.header-text {
	flex: 1;
}

.title {
	display: block;
	font-size: 32rpx;
	font-weight: bold;
	color: white;
}

.subtitle {
	display: block;
	font-size: 22rpx;
	color: rgba(255, 255, 255, 0.8);
	margin-top: 5rpx;
}

.notification {
	position: relative;
	width: 60rpx;
	height: 60rpx;
	display: flex;
	align-items: center;
	justify-content: center;
}

.notification .icon {
	font-size: 40rpx;
}

.badge {
	position: absolute;
	top: 0;
	right: 0;
	background: #f56c6c;
	color: white;
	font-size: 20rpx;
	padding: 2rpx 8rpx;
	border-radius: 20rpx;
	min-width: 30rpx;
	text-align: center;
}

.quick-entry {
	display: flex;
	background: white;
	padding: 40rpx 20rpx;
	margin-top: -20rpx;
	border-radius: 20rpx 20rpx 0 0;
}

.entry-item {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
}

.entry-icon {
	width: 100rpx;
	height: 100rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	border-radius: 20rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 50rpx;
	margin-bottom: 15rpx;
}

.doctor-icon {
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
}

.admin-icon {
	background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
}

.entry-text {
	font-size: 24rpx;
	color: #666;
}

.section {
	background: white;
	margin-top: 20rpx;
	padding: 30rpx;
}

.section-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 25rpx;
}

.section-title {
	font-size: 32rpx;
	font-weight: bold;
	color: #333;
}

.section-more {
	font-size: 26rpx;
	color: #07c160;
}

.doctor-scroll {
	white-space: nowrap;
}

.doctor-list {
	display: inline-flex;
}

.doctor-card {
	display: inline-block;
	width: 200rpx;
	margin-right: 20rpx;
}

.doctor-avatar {
	width: 100%;
	height: 200rpx;
	border-radius: 15rpx;
	margin-bottom: 15rpx;
}

.doctor-info {
	display: flex;
	flex-direction: column;
}

.doctor-name {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 5rpx;
}

.doctor-title {
	font-size: 24rpx;
	color: #07c160;
	margin-bottom: 3rpx;
}

.doctor-dept {
	font-size: 22rpx;
	color: #999;
}

.consultation-list {
	
}

.consultation-item {
	padding: 25rpx;
	background: #f9f9f9;
	border-radius: 15rpx;
	margin-bottom: 20rpx;
}

.consultation-item:last-child {
	margin-bottom: 0;
}

.consultation-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 15rpx;
}

.doctor-name {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
}

.status {
	font-size: 24rpx;
	padding: 5rpx 15rpx;
	border-radius: 20rpx;
}

.status-0 {
	background: #fff3e0;
	color: #ff9800;
}

.status-1 {
	background: #e3f2fd;
	color: #2196f3;
}

.status-2 {
	background: #e8f5e9;
	color: #4caf50;
}

.status-urgent {
	background: #fff3e0;
	color: #ff9800;
	font-weight: bold;
}

.status-pending {
	background: #fff3e0;
	color: #ff9800;
	font-size: 22rpx;
	padding: 5rpx 15rpx;
	border-radius: 20rpx;
}

/* 管理员专属样式 */
.application-list {
	
}

.application-item {
	padding: 25rpx;
	background: #f9f9f9;
	border-radius: 15rpx;
	margin-bottom: 20rpx;
}

.application-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 15rpx;
}

.applicant-info {
	flex: 1;
}

.applicant-name {
	display: block;
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 5rpx;
}

.applicant-title {
	display: block;
	font-size: 24rpx;
	color: #666;
}

.application-time {
	display: block;
	font-size: 22rpx;
	color: #999;
}

.stats-grid {
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	gap: 20rpx;
}

.stat-card {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	padding: 30rpx;
	border-radius: 15rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
}

.stat-card .stat-number {
	font-size: 48rpx;
	font-weight: bold;
	color: white;
	margin-bottom: 10rpx;
}

.stat-card .stat-label {
	font-size: 24rpx;
	color: rgba(255, 255, 255, 0.9);
}

.complaint {
	display: block;
	font-size: 26rpx;
	color: #666;
	margin-bottom: 10rpx;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.time {
	display: block;
	font-size: 22rpx;
	color: #999;
}

.security-banner {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	margin: 20rpx 30rpx;
	padding: 30rpx;
	border-radius: 20rpx;
	display: flex;
	align-items: center;
}

.banner-icon {
	font-size: 60rpx;
	margin-right: 20rpx;
}

.banner-content {
	flex: 1;
}

.banner-title {
	display: block;
	font-size: 28rpx;
	font-weight: bold;
	color: white;
	margin-bottom: 10rpx;
}

.banner-desc {
	display: block;
	font-size: 22rpx;
	color: rgba(255, 255, 255, 0.9);
	line-height: 1.5;
}
</style>
