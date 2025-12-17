<template>
	<view class="list-page">
		<!-- 管理员提示 -->
		<view class="admin-notice" v-if="isAdmin">
			<view class="notice-card">
				<text class="notice-icon">🛡️</text>
				<text class="notice-title">管理员身份</text>
				<text class="notice-text">管理员不需要使用问诊功能</text>
				<text class="notice-desc">您可以通过系统管理、用户管理和数据统计查看相关数据</text>
				<button class="back-btn" @click="goHome">返回首页</button>
			</view>
		</view>
		
		<!-- 原有问诊页面内容 -->
		<view v-else>
			<!-- 角色切换 -->
			<view class="role-tabs">
				<view 
					class="role-tab" 
					:class="{ active: currentRole === 'patient' }"
					@click="switchRole('patient')"
				>
					我的问诊
				</view>
				<view 
					class="role-tab" 
					:class="{ active: currentRole === 'doctor' }"
					@click="switchRole('doctor')"
					v-if="isDoctorAuth"
				>
					接诊记录
				</view>
			</view>
			
			<!-- 状态筛选 -->
			<scroll-view class="status-tabs" scroll-x>
				<view 
					class="status-tab" 
					:class="{ active: currentStatus === item.value }"
					v-for="item in statusList" 
					:key="item.value"
					@click="selectStatus(item.value)"
				>
					{{ item.label }}
				</view>
			</scroll-view>
			
			<!-- 问诊列表 -->
			<scroll-view 
				class="consultation-list" 
				scroll-y
				@scrolltolower="loadMore"
			>
				<view 
					class="consultation-item" 
					v-for="item in list" 
					:key="item.consultationId"
					@click="goDetail(item.consultationId)"
				>
					<view class="item-header">
						<view class="doctor-info">
							<image class="avatar" :src="item.avatar || '/static/default-avatar.png'" mode="aspectFill"></image>
							<view class="info">
								<text class="name">{{ currentRole === 'patient' ? item.doctorName : item.patientName }}</text>
								<text class="role">{{ currentRole === 'patient' ? '医生' : '患者' }}</text>
							</view>
						</view>
						<text class="status" :class="'status-' + item.status">{{ item.statusText }}</text>
					</view>
					
					<text class="complaint">主诉：{{ item.chiefComplaint }}</text>
					
					<view class="item-footer">
						<text class="time">{{ item.createdAt }}</text>
						<text class="ai-tag" v-if="item.needAI">🤖 AI辅助</text>
					</view>
				</view>
				
				<!-- 加载状态 -->
				<view class="loading" v-if="loading">加载中...</view>
				<view class="no-more" v-if="!hasMore && list.length > 0">没有更多了</view>
				<view class="empty" v-if="!loading && list.length === 0">
					<text class="empty-icon">💬</text>
					<text class="empty-text">暂无问诊记录</text>
					<button class="start-btn" @click="goCreate" v-if="currentRole === 'patient'">发起问诊</button>
				</view>
			</scroll-view>
			
			<!-- 悬浮按钮 -->
			<view class="fab" @click="goCreate" v-if="currentRole === 'patient'">
				<text class="fab-icon">+</text>
			</view>
		</view>
	</view>
</template>

<script>
import { get } from '@/utils/request.js'
import { API, STORAGE_KEYS } from '@/utils/config.js'
import { getStorageSync } from '@/utils/storage.js'

export default {
	data() {
		return {
			isAdmin: false,
			currentRole: 'patient',
			currentStatus: '',
			isDoctorAuth: false,
			statusList: [
				{ label: '全部', value: '' },
				{ label: '待接诊', value: 0 },
				{ label: '问诊中', value: 1 },
				{ label: '已完成', value: 2 }
			],
			list: [],
			page: 1,
			pageSize: 10,
			loading: false,
			hasMore: true
		}
	},
	
	onLoad(options) {
		// 检查是否为管理员
		const userInfo = getStorageSync(STORAGE_KEYS.USER_INFO)
		console.log('[问诊列表] 用户信息:', userInfo)
		this.isAdmin = userInfo && userInfo.role === 'admin'
		
		// 如果是管理员，不加载问诊数据
		if (this.isAdmin) {
			return
		}
		
		// 检查是否为认证医生
		this.isDoctorAuth = userInfo && userInfo.role === 'doctor' && userInfo.certStatus === 'approved'
		console.log('[问诊列表] 是否为认证医生:', this.isDoctorAuth, '角色:', userInfo?.role, '认证状态:', userInfo?.certStatus)
		
		// 如果是医生，默认显示接诊记录
		if (this.isDoctorAuth) {
			this.currentRole = 'doctor'
			console.log('[问诊列表] 设置为医生角色，默认显示接诊记录')
			// 如果从用户页面点击"待接诊"进入，默认筛选待接诊状态
			if (options.status !== undefined) {
				this.currentStatus = parseInt(options.status)
			} else {
				// 否则默认显示待接诊
				this.currentStatus = 0
			}
			console.log('[问诊列表] 当前状态筛选:', this.currentStatus)
		}
		
		this.loadList()
	},
	
	onShow() {
		// 如果是管理员，不刷新列表
		if (this.isAdmin) {
			return
		}
		
		// 重新检查医生身份（防止用户信息更新）
		const userInfo = getStorageSync(STORAGE_KEYS.USER_INFO)
		this.isDoctorAuth = userInfo && userInfo.role === 'doctor' && userInfo.certStatus === 'approved'
		
		// 如果是医生且当前角色不是医生，切换到医生角色
		if (this.isDoctorAuth && this.currentRole !== 'doctor') {
			this.currentRole = 'doctor'
			this.currentStatus = 0  // 默认显示待接诊
		}
		
		// 每次显示刷新列表
		this.loadList(true)
	},
	
	methods: {
		// 返回首页
		goHome() {
			uni.switchTab({
				url: '/pages/index/index'
			})
		},
		
		// 切换角色
		switchRole(role) {
			this.currentRole = role
			this.currentStatus = ''
			this.loadList(true)
		},
		
		// 选择状态
		selectStatus(status) {
			this.currentStatus = status
			this.loadList(true)
		},
		
		// 加载列表
		async loadList(isRefresh = false) {
			if (this.loading) return
			
			if (isRefresh) {
				this.page = 1
				this.list = []
				this.hasMore = true
			}
			
			this.loading = true
			
			try {
				const params = {
					page: this.page,
					pageSize: this.pageSize,
					role: this.currentRole
				}
				
				if (this.currentStatus !== '') {
					params.status = this.currentStatus
				}
				
				const res = await get(API.CONSULTATION_LIST, params)
				
				const list = res.data.list || []
				
				if (isRefresh) {
					this.list = list
				} else {
					this.list.push(...list)
				}
				
				this.hasMore = this.list.length < res.data.total
				
			} catch (error) {
				console.error('加载问诊列表失败:', error)
			} finally {
				this.loading = false
			}
		},
		
		// 加载更多
		loadMore() {
			if (this.hasMore && !this.loading) {
				this.page++
				this.loadList()
			}
		},
		
		// 跳转详情
		goDetail(id) {
			uni.navigateTo({
				url: '/pages/consultation/consultation-detail?id=' + id
			})
		},
		
		// 发起问诊
		goCreate() {
			uni.navigateTo({
				url: '/pages/consultation/create-consultation'
			})
		}
	}
}
</script>

<style scoped>
.list-page {
	height: 100vh;
	display: flex;
	flex-direction: column;
	background: #f5f5f5;
}

/* 管理员提示 */
.admin-notice {
	flex: 1;
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 40rpx;
}

.notice-card {
	background: white;
	border-radius: 30rpx;
	padding: 80rpx 60rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	box-shadow: 0 10rpx 40rpx rgba(0, 0, 0, 0.08);
}

.notice-icon {
	font-size: 120rpx;
	margin-bottom: 30rpx;
}

.notice-title {
	font-size: 36rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 20rpx;
}

.notice-text {
	font-size: 28rpx;
	color: #666;
	margin-bottom: 15rpx;
}

.notice-desc {
	font-size: 24rpx;
	color: #999;
	text-align: center;
	line-height: 1.6;
	margin-bottom: 50rpx;
}

.back-btn {
	width: 300rpx;
	height: 80rpx;
	background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
	color: white;
	border: none;
	border-radius: 40rpx;
	font-size: 28rpx;
}

.back-btn::after {
	border: none;
}

.role-tabs {
	display: flex;
	background: white;
	padding: 10rpx 30rpx 0;
}

.role-tab {
	flex: 1;
	text-align: center;
	padding: 25rpx 0;
	font-size: 28rpx;
	color: #666;
	position: relative;
}

.role-tab.active {
	color: #07c160;
	font-weight: bold;
}

.role-tab.active::after {
	content: '';
	position: absolute;
	bottom: 0;
	left: 50%;
	transform: translateX(-50%);
	width: 60rpx;
	height: 6rpx;
	background: #07c160;
	border-radius: 3rpx;
}

.status-tabs {
	background: white;
	white-space: nowrap;
	padding: 20rpx 30rpx;
	border-top: 1px solid #f0f0f0;
}

.status-tab {
	display: inline-block;
	padding: 10rpx 25rpx;
	margin-right: 20rpx;
	background: #f5f5f5;
	border-radius: 30rpx;
	font-size: 26rpx;
	color: #666;
}

.status-tab.active {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
}

.consultation-list {
	flex: 1;
	padding: 20rpx 30rpx;
}

.consultation-item {
	background: white;
	border-radius: 20rpx;
	padding: 25rpx;
	margin-bottom: 20rpx;
}

.item-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20rpx;
}

.doctor-info {
	display: flex;
	align-items: center;
}

.avatar {
	width: 80rpx;
	height: 80rpx;
	border-radius: 50%;
	margin-right: 20rpx;
}

.info {
	display: flex;
	flex-direction: column;
}

.name {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 5rpx;
}

.role {
	font-size: 22rpx;
	color: #999;
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

.complaint {
	display: block;
	font-size: 26rpx;
	color: #666;
	margin-bottom: 15rpx;
	overflow: hidden;
	text-overflow: ellipsis;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
}

.item-footer {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.time {
	font-size: 22rpx;
	color: #999;
}

.ai-tag {
	font-size: 22rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
	padding: 5rpx 12rpx;
	border-radius: 10rpx;
}

.loading {
	text-align: center;
	padding: 30rpx 0;
	font-size: 26rpx;
	color: #999;
}

.no-more {
	text-align: center;
	padding: 30rpx 0;
	font-size: 26rpx;
	color: #999;
}

.empty {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 150rpx 0;
}

.empty-icon {
	font-size: 120rpx;
	margin-bottom: 30rpx;
}

.empty-text {
	font-size: 28rpx;
	color: #999;
	margin-bottom: 40rpx;
}

.start-btn {
	width: 300rpx;
	height: 80rpx;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
	border: none;
	border-radius: 40rpx;
	font-size: 28rpx;
}

.start-btn::after {
	border: none;
}

.fab {
	position: fixed;
	right: 30rpx;
	bottom: 100rpx;
	width: 100rpx;
	height: 100rpx;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	border-radius: 50%;
	box-shadow: 0 4rpx 20rpx rgba(7, 193, 96, 0.4);
	display: flex;
	align-items: center;
	justify-content: center;
}

.fab-icon {
	font-size: 60rpx;
	color: white;
	font-weight: 300;
}
</style>
