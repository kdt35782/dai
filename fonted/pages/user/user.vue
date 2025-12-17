<template>
	<view class="user-page">
		<!-- 用户信息卡片 -->
		<view class="user-card">
			<image class="avatar" :src="getAvatarUrl(userInfo.avatar)" mode="aspectFill" @click="editProfile"></image>
			<view class="user-info">
				<text class="username">{{ userInfo.username || '未登录' }}</text>
				<text class="real-name" v-if="userInfo.realName">{{ userInfo.realName }}</text>
				<view class="role-tag" v-if="userInfo.role === 'doctor'">
					<text>{{ userInfo.doctorTitle || '医生' }}</text>
				</view>
			</view>
			<text class="edit-btn" @click="editProfile">编辑 ></text>
		</view>
		
		<!-- 快捷统计 -->
		<view class="stat-card">
			<!-- 患者端统计 -->
			<template v-if="userInfo.role === 'user'">
				<view class="stat-item" @click="goConsultationList">
					<text class="stat-number">{{ stats.consultationCount || 0 }}</text>
					<text class="stat-label">我的问诊</text>
				</view>
				<view class="stat-item" @click="goRecordList">
					<text class="stat-number">{{ stats.recordCount || 0 }}</text>
					<text class="stat-label">我的病历</text>
				</view>
				<view class="stat-item" @click="goNotification">
					<text class="stat-number">{{ stats.unreadCount || 0 }}</text>
					<text class="stat-label">未读消息</text>
				</view>
			</template>
			
			<!-- 医生端统计 -->
			<template v-else-if="userInfo.role === 'doctor'">
				<view class="stat-item" @click="goConsultationList">
					<text class="stat-number">{{ stats.pendingCount || 0 }}</text>
					<text class="stat-label">待接诊</text>
				</view>
				<view class="stat-item" @click="goRecordList">
					<text class="stat-number">{{ stats.recordCount || 0 }}</text>
					<text class="stat-label">患者病历</text>
				</view>
				<view class="stat-item" @click="goNotification">
					<text class="stat-number">{{ stats.unreadCount || 0 }}</text>
					<text class="stat-label">未读消息</text>
				</view>
			</template>
			
			<!-- 管理员端统计 -->
			<template v-else-if="userInfo.role === 'admin'">
				<!-- 医生申请审核功能已废弃 -->
				<!-- <view class="stat-item" @click="goDoctorApplications">
					<text class="stat-number admin-number">{{ stats.pendingApplications || 0 }}</text>
					<text class="stat-label">待审核</text>
				</view> -->
				<view class="stat-item" @click="goUserManagement">
					<text class="stat-number admin-number">{{ stats.totalUsers || 0 }}</text>
					<text class="stat-label">总用户</text>
				</view>
				<view class="stat-item" @click="goSystemLogs">
					<text class="stat-number admin-number">-</text>
					<text class="stat-label">系统日志</text>
				</view>
				<view class="stat-item" @click="goNotification">
					<text class="stat-number admin-number">{{ stats.unreadCount || 0 }}</text>
					<text class="stat-label">未读消息</text>
				</view>
			</template>
		</view>
		
		<!-- 功能列表 -->
		<view class="menu-section">
			<!-- 管理员专属功能 -->
			<template v-if="userInfo.role === 'admin'">
				<!-- 医生申请审核功能已废弃，现在医生注册直接生效 -->
				<!-- <view class="menu-item" @click="goDoctorApplications">
					<text class="menu-icon">📋</text>
					<text class="menu-text">医生申请审核</text>
					<view class="badge" v-if="stats.pendingApplications > 0">{{ stats.pendingApplications }}</view>
					<text class="menu-arrow">></text>
				</view> -->
				
				<view class="menu-item" @click="goUserManagement">
					<text class="menu-icon">👥</text>
					<text class="menu-text">用户管理</text>
					<text class="menu-arrow">></text>
				</view>
				
				<view class="menu-item" @click="goSystemLogs">
					<text class="menu-icon">📄</text>
					<text class="menu-text">系统日志</text>
					<text class="menu-arrow">></text>
				</view>
			</template>
			
			<!-- 医生专属功能 -->
			<template v-else-if="userInfo.role === 'doctor'">
				<view class="menu-item" @click="goDoctorProfile">
					<text class="menu-icon">🎖️</text>
					<text class="menu-text">医生资质</text>
					<text class="menu-arrow">></text>
				</view>
				
				<view class="menu-item" @click="goRecordList">
					<text class="menu-icon">📋</text>
					<text class="menu-text">患者病历</text>
					<text class="menu-arrow">></text>
				</view>
			</template>
			
			<!-- 患者功能 -->
			<template v-else>
				<view class="menu-item" @click="goRecordList">
					<text class="menu-icon">📋</text>
					<text class="menu-text">我的病历</text>
					<text class="menu-arrow">></text>
				</view>
				
				<!-- 申请成为医生功能已废弃，现在注册时直接选择 -->
				<!-- <view class="menu-item" @click="goApplyDoctor">
					<text class="menu-icon">🎓</text>
					<text class="menu-text">申请成为医生</text>
					<text class="menu-arrow">></text>
				</view> -->
			</template>
			
			<!-- 通用功能 -->
			<view class="menu-item" @click="showKeyManagement">
				<text class="menu-icon">🔐</text>
				<text class="menu-text">密钥管理</text>
				<text class="menu-arrow">></text>
			</view>
			
			<view class="menu-item" @click="goNotification">
				<text class="menu-icon">🔔</text>
				<text class="menu-text">消息通知</text>
				<view class="badge" v-if="stats.unreadCount > 0">{{ stats.unreadCount }}</view>
				<text class="menu-arrow">></text>
			</view>
		</view>
		
		<!-- 设置 -->
		<view class="menu-section">
			<view class="menu-item" @click="changePassword">
				<text class="menu-icon">🔑</text>
				<text class="menu-text">修改密码</text>
				<text class="menu-arrow">></text>
			</view>
			
			<view class="menu-item" @click="showAbout">
				<text class="menu-icon">ℹ️</text>
				<text class="menu-text">关于系统</text>
				<text class="menu-arrow">></text>
			</view>
		</view>
		
		<!-- 退出登录 -->
		<view class="logout-section">
			<button class="logout-btn" @click="handleLogout">
				退出登录
			</button>
		</view>
	</view>
</template>

<script>
import { get, post } from '@/utils/request.js'
import { API, STORAGE_KEYS } from '@/utils/config.js'
import { getStorageSync, setStorageSync, removeStorageSync } from '@/utils/storage.js'

export default {
	data() {
		return {
			userInfo: {},
			stats: {
				consultationCount: 0,
				pendingCount: 0,  // 医生待接诊数量
				recordCount: 0,
				unreadCount: 0,
				pendingApplications: 0,  // 管理员待审核数量
				totalUsers: 0  // 管理员总用户数
			}
		}
	},
	
	onShow() {
		this.loadUserInfo()
		this.loadStats()
	},
	
	methods: {
		// 加载用户信息
		async loadUserInfo() {
			// 先从本地存储加载，快速显示
			const localUserInfo = getStorageSync(STORAGE_KEYS.USER_INFO)
			if (localUserInfo) {
				this.userInfo = localUserInfo
			}
			
			// 如果未登录，跳转到登录页
			if (!localUserInfo) {
				uni.navigateTo({
					url: '/pages/login/login'
				})
				return
			}
			
			// 从后端获取最新用户信息（包括头像）
			try {
				const res = await get(API.USER_INFO)
				if (res.data) {
					this.userInfo = res.data
					// 更新本地存储
					setStorageSync(STORAGE_KEYS.USER_INFO, res.data)
				}
			} catch (error) {
				console.error('加载用户信息失败:', error)
				// 如果401错误，request.js会自动处理跳转
			}
		},
		
		// 获取头像URL（处理相对路径）
		getAvatarUrl(avatar) {
			if (!avatar) {
				// 返回默认头像（使用网络图片或base64）
				return 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iI2RkZCIvPjx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBmb250LXNpemU9IjQwIiBmaWxsPSIjOTk5IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBkeT0iLjNlbSI+8J+RpDwvdGV4dD48L3N2Zz4='
			}
			
			// 如果是完整URL，直接返回
			if (avatar.startsWith('http://') || avatar.startsWith('https://') || avatar.startsWith('data:')) {
				return avatar
			}
			
			// 如果是相对路径，拼接完整URL
			const API_BASE_URL = 'http://localhost:3000'
			return API_BASE_URL + avatar
		},
		
		// 加载统计数据
		async loadStats() {
			try {
				const userRole = this.userInfo.role
				
				if (userRole === 'doctor') {
					// 医生端：加载待接诊数量
					const pendingRes = await get(API.CONSULTATION_LIST, {
						page: 1,
						pageSize: 1,
						role: 'doctor',
						status: 0  // 待接诊
					})
					this.stats.pendingCount = pendingRes.data.total || 0
					
					// 加载病历数量
					const recordRes = await get(API.RECORD_LIST, {
						page: 1,
						pageSize: 1
					})
					this.stats.recordCount = recordRes.data.total || 0
				} else if (userRole === 'admin') {
					// 管理员端：加载待审核申请数量
					const applicationRes = await get(API.USER_ADMIN_APPLICATIONS, {
						page: 1,
						pageSize: 1,
						status: 0  // 待审核
					})
					this.stats.pendingApplications = applicationRes.data.total || 0
					
					// 加载总用户数
					const userRes = await get(API.USER_ADMIN_USERS, {
						page: 1,
						pageSize: 1
					})
					this.stats.totalUsers = userRes.data.total || 0
				} else {
					// 患者端：加载问诊数量
					const consultationRes = await get(API.CONSULTATION_LIST, {
						page: 1,
						pageSize: 1,
						role: 'patient'
					})
					this.stats.consultationCount = consultationRes.data.total || 0
					
					// 加载病历数量
					const recordRes = await get(API.RECORD_LIST, {
						page: 1,
						pageSize: 1
					})
					this.stats.recordCount = recordRes.data.total || 0
				}
				
				// 加载未读消息数（所有角色）
				const notificationRes = await get(API.NOTIFICATION_UNREAD_COUNT)
				this.stats.unreadCount = notificationRes.data.totalUnread || 0
				
			} catch (error) {
				console.error('加载统计数据失败:', error)
			}
		},
		
		// 编辑资料
		editProfile() {
			uni.navigateTo({
				url: '/pages/user/edit-profile'
			})
		},
		
		// 跳转问诊列表
		goConsultationList() {
			uni.switchTab({
				url: '/pages/consultation/consultation-list'
			})
		},
		
		// 跳转病历列表
		goRecordList() {
			uni.navigateTo({
				url: '/pages/medical-record/record-list'
			})
		},
		
		// 跳转消息通知
		goNotification() {
			uni.navigateTo({
				url: '/pages/notification/notification'
			})
		},
		
		// 申请成为医生
		goApplyDoctor() {
			uni.navigateTo({
				url: '/pages/user/apply-doctor'
			})
		},
		
		// 医生申请审核（管理员）
		goDoctorApplications() {
			uni.navigateTo({
				url: '/pages/admin/doctor-applications'
			})
		},
		
		// 医生资质页面
		goDoctorProfile() {
			uni.navigateTo({
				url: '/pages/user/doctor-profile'
			})
		},
		
		// 用户管理（管理员）
		goUserManagement() {
			uni.navigateTo({
				url: '/pages/admin/user-management'
			})
		},
		
		// 系统日志（管理员）
		goSystemLogs() {
			uni.navigateTo({
				url: '/pages/admin/system-logs'
			})
		},
		
		// 密钥管理
		showKeyManagement() {
			const sm2PublicKey = getStorageSync(STORAGE_KEYS.SM2_PUBLIC_KEY)
			
			uni.showModal({
				title: '密钥管理',
				content: `SM2公钥：${sm2PublicKey || '未设置'}\n\n提示：密钥用于加密敏感数据，请妥善保管`,
				confirmText: '重新生成',
				success: async (res) => {
					if (res.confirm) {
						try {
							// 调用重新生成密钥API
							const result = await post(API.KEY_GENERATE)
							
							setStorageSync(STORAGE_KEYS.SM2_PUBLIC_KEY, result.data.publicKey)
							
							uni.showToast({
								title: '密钥已更新',
								icon: 'success'
							})
						} catch (error) {
							console.error('生成密钥失败:', error)
						}
					}
				}
			})
		},
		
		// 修改密码
		changePassword() {
			uni.navigateTo({
				url: '/pages/user/change-password'
			})
		},
		
		// 关于系统
		showAbout() {
			uni.showModal({
				title: '关于系统',
				content: '基于国密加密的网上问诊平台\n\n版本：v1.0.0\n\n采用SM2/SM3/SM4国密算法\n全程加密保护您的隐私',
				showCancel: false
			})
		},
		
		// 退出登录
		handleLogout() {
			uni.showModal({
				title: '提示',
				content: '确认退出登录吗？',
				success: (res) => {
					if (res.confirm) {
						// 清除本地数据
						removeStorageSync(STORAGE_KEYS.TOKEN)
						removeStorageSync(STORAGE_KEYS.USER_INFO)
						
						uni.showToast({
							title: '已退出登录',
							icon: 'success'
						})
						
						// 跳转到登录页
						setTimeout(() => {
							uni.reLaunch({
								url: '/pages/login/login'
							})
						}, 1500)
					}
				}
			})
		}
	}
}
</script>

<style scoped>
.user-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding-bottom: 40rpx;
}

.user-card {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	padding: 50rpx 30rpx;
	display: flex;
	align-items: center;
}

.avatar {
	width: 120rpx;
	height: 120rpx;
	border-radius: 50%;
	margin-right: 25rpx;
	border: 4rpx solid rgba(255, 255, 255, 0.3);
}

.user-info {
	flex: 1;
}

.username {
	display: block;
	font-size: 32rpx;
	font-weight: bold;
	color: white;
	margin-bottom: 8rpx;
}

.real-name {
	display: block;
	font-size: 24rpx;
	color: rgba(255, 255, 255, 0.8);
	margin-bottom: 8rpx;
}

.role-tag {
	display: inline-block;
	background: rgba(255, 255, 255, 0.2);
	color: white;
	font-size: 22rpx;
	padding: 5rpx 15rpx;
	border-radius: 20rpx;
}

.edit-btn {
	font-size: 26rpx;
	color: white;
}

.stat-card {
	display: flex;
	background: white;
	margin: -30rpx 30rpx 20rpx;
	border-radius: 20rpx;
	padding: 30rpx 0;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.stat-item {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	border-right: 1px solid #f0f0f0;
}

.stat-item:last-child {
	border-right: none;
}

.stat-number {
	font-size: 36rpx;
	font-weight: bold;
	color: #07c160;
	margin-bottom: 10rpx;
}

.admin-number {
	color: #ff6b6b;
}

.stat-label {
	font-size: 24rpx;
	color: #999;
}

.menu-section {
	background: white;
	margin: 20rpx 30rpx;
	border-radius: 20rpx;
	overflow: hidden;
}

.menu-item {
	display: flex;
	align-items: center;
	padding: 30rpx 25rpx;
	border-bottom: 1px solid #f0f0f0;
	position: relative;
}

.menu-item:last-child {
	border-bottom: none;
}

.menu-icon {
	font-size: 40rpx;
	margin-right: 20rpx;
}

.menu-text {
	flex: 1;
	font-size: 28rpx;
	color: #333;
}

.badge {
	background: #f56c6c;
	color: white;
	font-size: 20rpx;
	padding: 4rpx 12rpx;
	border-radius: 20rpx;
	margin-right: 15rpx;
}

.menu-arrow {
	font-size: 28rpx;
	color: #999;
}

.logout-section {
	padding: 0 30rpx;
	margin-top: 40rpx;
}

.logout-btn {
	width: 100%;
	height: 90rpx;
	background: white;
	color: #f56c6c;
	border: none;
	border-radius: 45rpx;
	font-size: 32rpx;
}

.logout-btn::after {
	border: none;
}
</style>
