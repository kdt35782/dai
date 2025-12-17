<template>
	<view class="detail-page">
		<!-- 医生信息 -->
		<view class="doctor-card">
			<view class="doctor-header">
				<image class="avatar" :src="doctorInfo.avatar || '/static/default-avatar.png'" mode="aspectFill"></image>
				<view class="info">
					<view class="name-row">
						<text class="name">{{ doctorInfo.realName }}</text>
						<text class="title">{{ doctorInfo.doctorTitle }}</text>
					</view>
					<text class="dept">{{ doctorInfo.doctorDept }}</text>
					<view class="stats">
						<text class="stat">💬 {{ doctorInfo.consultationCount || 0 }}次</text>
						<text class="stat">⭐ {{ doctorInfo.rating || '5.0' }}分</text>
					</view>
				</view>
			</view>
			
			<view class="specialty-section" v-if="doctorInfo.specialty">
				<text class="section-title">擅长领域</text>
				<text class="specialty-text">{{ doctorInfo.specialty }}</text>
			</view>
		</view>
		
		<!-- 医生介绍 -->
		<view class="intro-card" v-if="doctorInfo.introduction">
			<text class="card-title">医生介绍</text>
			<text class="intro-text">{{ doctorInfo.introduction }}</text>
		</view>
		
		<!-- 执业信息 -->
		<view class="cert-card">
			<text class="card-title">执业信息</text>
			
			<view class="cert-item">
				<text class="cert-label">执业证号</text>
				<text class="cert-value">{{ doctorInfo.certNumber || '暂无' }}</text>
			</view>
			
			<view class="cert-item">
				<text class="cert-label">认证状态</text>
				<text class="cert-status" :class="doctorInfo.certStatus === 'approved' ? 'approved' : ''">
					{{ doctorInfo.certStatus === 'approved' ? '已认证' : '未认证' }}
				</text>
			</view>
		</view>
		
		<!-- 底部操作按钮 -->
		<view class="bottom-bar">
			<button class="consult-btn" @click="createConsultation">
				立即问诊
			</button>
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
			userId: 0,
			doctorInfo: {}
		}
	},
	
	onLoad(options) {
		if (options.userId) {
			this.userId = parseInt(options.userId)
			this.loadDoctorInfo()
		}
	},
	
	methods: {
		// 加载医生信息
		async loadDoctorInfo() {
			try {
				// 使用公开接口获取医生详情
				const res = await get(`${API.USER_DOCTOR_DETAIL}/${this.userId}`, {}, { noAuth: true })
				
				this.doctorInfo = res.data || {}
				
				// 设置标题
				uni.setNavigationBarTitle({
					title: this.doctorInfo.realName + ' 医生'
				})
				
			} catch (error) {
				console.error('加载医生信息失败:', error)
				uni.showToast({
					title: '加载失败',
					icon: 'none'
				})
			}
		},
		
		// 发起问诊
		createConsultation() {
			const token = getStorageSync(STORAGE_KEYS.TOKEN)
			
			if (!token) {
				uni.showToast({
					title: '请先登录',
					icon: 'none'
				})
				
				setTimeout(() => {
					uni.navigateTo({
						url: '/pages/login/login'
					})
				}, 1500)
				return
			}
			
			uni.navigateTo({
				url: '/pages/consultation/create-consultation?doctorId=' + this.userId + '&doctorName=' + this.doctorInfo.realName
			})
		}
	}
}
</script>

<style scoped>
.detail-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 20rpx 30rpx;
	padding-bottom: 140rpx;
}

.doctor-card {
	background: white;
	border-radius: 20rpx;
	padding: 30rpx;
	margin-bottom: 20rpx;
}

.doctor-header {
	display: flex;
	margin-bottom: 30rpx;
}

.avatar {
	width: 150rpx;
	height: 150rpx;
	border-radius: 15rpx;
	margin-right: 25rpx;
}

.info {
	flex: 1;
	display: flex;
	flex-direction: column;
	justify-content: center;
}

.name-row {
	display: flex;
	align-items: center;
	margin-bottom: 12rpx;
}

.name {
	font-size: 34rpx;
	font-weight: bold;
	color: #333;
	margin-right: 15rpx;
}

.title {
	font-size: 22rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
	padding: 5rpx 15rpx;
	border-radius: 10rpx;
}

.dept {
	font-size: 26rpx;
	color: #07c160;
	margin-bottom: 12rpx;
}

.stats {
	display: flex;
}

.stat {
	font-size: 24rpx;
	color: #999;
	margin-right: 25rpx;
}

.specialty-section {
	border-top: 1px solid #f0f0f0;
	padding-top: 25rpx;
}

.section-title {
	display: block;
	font-size: 26rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 15rpx;
}

.specialty-text {
	display: block;
	font-size: 26rpx;
	color: #666;
	line-height: 1.6;
}

.intro-card {
	background: white;
	border-radius: 20rpx;
	padding: 30rpx;
	margin-bottom: 20rpx;
}

.card-title {
	display: block;
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 20rpx;
}

.intro-text {
	display: block;
	font-size: 26rpx;
	color: #666;
	line-height: 1.8;
}

.cert-card {
	background: white;
	border-radius: 20rpx;
	padding: 30rpx;
}

.cert-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 20rpx 0;
	border-bottom: 1px solid #f0f0f0;
}

.cert-item:last-child {
	border-bottom: none;
}

.cert-label {
	font-size: 26rpx;
	color: #666;
}

.cert-value {
	font-size: 26rpx;
	color: #333;
}

.cert-status {
	font-size: 24rpx;
	padding: 5rpx 15rpx;
	border-radius: 10rpx;
	background: #f5f5f5;
	color: #999;
}

.cert-status.approved {
	background: #e8f5e9;
	color: #4caf50;
}

.bottom-bar {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: white;
	padding: 20rpx 30rpx;
	box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.consult-btn {
	width: 100%;
	height: 90rpx;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
	border: none;
	border-radius: 45rpx;
	font-size: 32rpx;
	font-weight: bold;
}

.consult-btn::after {
	border: none;
}
</style>
