<template>
	<view class="profile-page">
		<!-- 头部背景 -->
		<view class="header-bg"></view>
		
		<!-- 医生信息卡片 -->
		<view class="doctor-card">
			<image class="doctor-avatar" :src="doctorInfo.avatar || '/static/default-avatar.png'" mode="aspectFill"></image>
			<view class="doctor-info">
				<text class="doctor-name">{{ doctorInfo.realName || doctorInfo.username }}</text>
				<text class="doctor-title">{{ doctorInfo.doctorTitle || '医生' }}</text>
				<text class="doctor-dept">{{ doctorInfo.doctorDept || '暂无科室' }}</text>
			</view>
		</view>
		
		<!-- 认证状态 -->
		<view class="status-card">
			<view class="status-item">
				<text class="status-icon">✓</text>
				<text class="status-text">已认证医生</text>
			</view>
			<view class="status-item">
				<text class="status-label">认证时间：</text>
				<text class="status-value">{{ doctorInfo.createdAt || '未知' }}</text>
			</view>
		</view>
		
		<!-- 详细信息 -->
		<view class="section">
			<view class="section-title">基本信息</view>
			<view class="info-card">
				<view class="info-item">
					<text class="info-label">用户名</text>
					<text class="info-value">{{ doctorInfo.username }}</text>
				</view>
				<view class="info-item">
					<text class="info-label">真实姓名</text>
					<text class="info-value">{{ doctorInfo.realName || '未填写' }}</text>
				</view>
				<view class="info-item">
					<text class="info-label">性别</text>
					<text class="info-value">{{ getGenderText(doctorInfo.gender) }}</text>
				</view>
				<view class="info-item">
					<text class="info-label">联系电话</text>
					<text class="info-value">{{ doctorInfo.phone || '未填写' }}</text>
				</view>
				<view class="info-item">
					<text class="info-label">邮箱</text>
					<text class="info-value">{{ doctorInfo.email || '未填写' }}</text>
				</view>
			</view>
		</view>
		
		<!-- 执业信息 -->
		<view class="section">
			<view class="section-title">执业信息</view>
			<view class="info-card">
				<view class="info-item">
					<text class="info-label">职称</text>
					<text class="info-value">{{ doctorInfo.doctorTitle || '未填写' }}</text>
				</view>
				<view class="info-item">
					<text class="info-label">科室</text>
					<text class="info-value">{{ doctorInfo.doctorDept || '未填写' }}</text>
				</view>
				<view class="info-item" v-if="doctorInfo.doctorCert">
					<text class="info-label">资格证书</text>
					<view class="cert-wrapper" @click="previewCert">
						<image class="cert-thumb" :src="doctorInfo.doctorCert" mode="aspectFill"></image>
						<text class="view-text">点击查看</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 个人简介 -->
		<view class="section" v-if="doctorInfo.doctorIntro">
			<view class="section-title">个人简介</view>
			<view class="intro-card">
				<text class="intro-text">{{ doctorInfo.doctorIntro }}</text>
			</view>
		</view>
		
		<!-- 操作按钮 -->
		<view class="action-section">
			<button class="action-btn" @click="editProfile">
				编辑资料
			</button>
		</view>
	</view>
</template>

<script>
import { get } from '@/utils/request.js'
import { API, STORAGE_KEYS } from '@/utils/config.js'

export default {
	data() {
		return {
			doctorInfo: {}
		}
	},
	
	onLoad() {
		this.loadDoctorInfo()
	},
	
	methods: {
		// 加载医生信息
		async loadDoctorInfo() {
			try {
				const res = await get(API.USER_INFO)
				this.doctorInfo = res.data
				
				// 如果不是医生角色，提示并返回
				if (this.doctorInfo.role !== 'doctor') {
					uni.showToast({
						title: '仅医生可访问',
						icon: 'none'
					})
					setTimeout(() => {
						uni.navigateBack()
					}, 1500)
				}
				
			} catch (error) {
				console.error('加载医生信息失败:', error)
				// 401错误会被request.js自动处理，这里不再显示toast
				// 其他错误才显示提示
				if (error.code !== 401) {
					uni.showToast({
						title: error.message || '加载失败',
						icon: 'none'
					})
				}
			}
		},
		
		// 性别文本转换
		getGenderText(gender) {
			const map = {
				0: '未知',
				1: '男',
				2: '女'
			}
			return map[gender] || '未知'
		},
		
		// 预览证书
		previewCert() {
			if (this.doctorInfo.doctorCert) {
				uni.previewImage({
					urls: [this.doctorInfo.doctorCert],
					current: this.doctorInfo.doctorCert
				})
			}
		},
		
		// 编辑资料
		editProfile() {
			uni.navigateTo({
				url: '/pages/user/edit-profile'
			})
		}
	}
}
</script>

<style scoped>
.profile-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding-bottom: 40rpx;
}

.header-bg {
	height: 300rpx;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
}

.doctor-card {
	margin: -180rpx 30rpx 20rpx;
	background: white;
	border-radius: 20rpx;
	padding: 40rpx;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);
	display: flex;
	align-items: center;
}

.doctor-avatar {
	width: 140rpx;
	height: 140rpx;
	border-radius: 50%;
	border: 6rpx solid white;
	margin-right: 30rpx;
}

.doctor-info {
	flex: 1;
}

.doctor-name {
	display: block;
	font-size: 36rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 10rpx;
}

.doctor-title {
	display: block;
	font-size: 26rpx;
	color: #07c160;
	margin-bottom: 5rpx;
}

.doctor-dept {
	display: block;
	font-size: 24rpx;
	color: #999;
}

.status-card {
	background: white;
	margin: 0 30rpx 20rpx;
	border-radius: 15rpx;
	padding: 30rpx;
}

.status-item {
	display: flex;
	align-items: center;
	margin-bottom: 15rpx;
}

.status-item:last-child {
	margin-bottom: 0;
}

.status-icon {
	width: 40rpx;
	height: 40rpx;
	background: #07c160;
	color: white;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 24rpx;
	font-weight: bold;
	margin-right: 15rpx;
}

.status-text {
	font-size: 28rpx;
	color: #07c160;
	font-weight: bold;
}

.status-label {
	font-size: 24rpx;
	color: #999;
}

.status-value {
	font-size: 24rpx;
	color: #666;
}

.section {
	margin: 0 30rpx 20rpx;
}

.section-title {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 15rpx;
	padding-left: 10rpx;
}

.info-card {
	background: white;
	border-radius: 15rpx;
	overflow: hidden;
}

.info-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 30rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.info-item:last-child {
	border-bottom: none;
}

.info-label {
	font-size: 26rpx;
	color: #666;
}

.info-value {
	font-size: 26rpx;
	color: #333;
	text-align: right;
}

.cert-wrapper {
	display: flex;
	align-items: center;
}

.cert-thumb {
	width: 120rpx;
	height: 80rpx;
	border-radius: 10rpx;
	margin-right: 15rpx;
}

.view-text {
	font-size: 24rpx;
	color: #07c160;
}

.intro-card {
	background: white;
	border-radius: 15rpx;
	padding: 30rpx;
}

.intro-text {
	font-size: 26rpx;
	color: #666;
	line-height: 1.8;
}

.action-section {
	padding: 30rpx;
}

.action-btn {
	width: 100%;
	height: 90rpx;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
	border: none;
	border-radius: 45rpx;
	font-size: 30rpx;
	font-weight: bold;
}

.action-btn::after {
	border: none;
}
</style>
