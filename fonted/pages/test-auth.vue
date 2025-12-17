<template>
	<view class="test-page">
		<view class="test-card">
			<text class="title">认证调试工具</text>
			
			<view class="section">
				<text class="section-title">1. 本地存储检查</text>
				<view class="info-item">
					<text class="label">Token存在:</text>
					<text class="value">{{ hasToken ? '✅ 是' : '❌ 否' }}</text>
				</view>
				<view class="info-item">
					<text class="label">Token内容:</text>
					<text class="value small">{{ tokenPreview }}</text>
				</view>
				<view class="info-item">
					<text class="label">用户信息:</text>
					<text class="value small">{{ userInfoPreview }}</text>
				</view>
			</view>
			
			<view class="section">
				<text class="section-title">2. API测试</text>
				<button class="test-btn" @click="testUserInfo">测试 GET /api/user/info</button>
				<view class="result" v-if="testResult">
					<text class="result-title">测试结果:</text>
					<text class="result-content">{{ testResult }}</text>
				</view>
			</view>
			
			<view class="section">
				<text class="section-title">3. 操作</text>
				<button class="action-btn" @click="clearStorage">清除本地存储</button>
				<button class="action-btn primary" @click="goLogin">重新登录</button>
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
			hasToken: false,
			tokenPreview: '',
			userInfoPreview: '',
			testResult: ''
		}
	},
	
	onLoad() {
		this.checkStorage()
	},
	
	methods: {
		// 检查本地存储
		checkStorage() {
			const token = uni.getStorageSync(STORAGE_KEYS.TOKEN)
			const userInfo = uni.getStorageSync(STORAGE_KEYS.USER_INFO)
			
			this.hasToken = !!token
			this.tokenPreview = token ? `${token.substring(0, 30)}...` : '未找到'
			this.userInfoPreview = userInfo ? JSON.stringify(userInfo) : '未找到'
			
			console.log('[存储检查] Token:', token)
			console.log('[存储检查] UserInfo:', userInfo)
		},
		
		// 测试用户信息接口
		async testUserInfo() {
			this.testResult = '测试中...'
			
			try {
				const res = await get(API.USER_INFO)
				this.testResult = '✅ 成功!\n' + JSON.stringify(res, null, 2)
				
				uni.showToast({
					title: '接口调用成功',
					icon: 'success'
				})
			} catch (error) {
				this.testResult = '❌ 失败!\n' + JSON.stringify(error, null, 2)
				
				console.error('[测试失败]', error)
			}
		},
		
		// 清除本地存储
		clearStorage() {
			uni.showModal({
				title: '确认操作',
				content: '确定要清除所有本地存储吗？',
				success: (res) => {
					if (res.confirm) {
						uni.removeStorageSync(STORAGE_KEYS.TOKEN)
						uni.removeStorageSync(STORAGE_KEYS.USER_INFO)
						
						uni.showToast({
							title: '已清除',
							icon: 'success'
						})
						
						setTimeout(() => {
							this.checkStorage()
						}, 1000)
					}
				}
			})
		},
		
		// 跳转到登录页
		goLogin() {
			uni.reLaunch({
				url: '/pages/login/login'
			})
		}
	}
}
</script>

<style scoped>
.test-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 40rpx;
}

.test-card {
	background: white;
	border-radius: 20rpx;
	padding: 40rpx;
}

.title {
	display: block;
	font-size: 36rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 40rpx;
	text-align: center;
}

.section {
	margin-bottom: 40rpx;
	padding-bottom: 40rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.section:last-child {
	border-bottom: none;
	margin-bottom: 0;
}

.section-title {
	display: block;
	font-size: 28rpx;
	font-weight: bold;
	color: #666;
	margin-bottom: 20rpx;
}

.info-item {
	display: flex;
	justify-content: space-between;
	padding: 20rpx 0;
	border-bottom: 1rpx solid #f5f5f5;
}

.info-item:last-child {
	border-bottom: none;
}

.label {
	font-size: 26rpx;
	color: #999;
}

.value {
	font-size: 26rpx;
	color: #333;
	font-weight: bold;
	flex: 1;
	text-align: right;
}

.value.small {
	font-size: 22rpx;
	word-break: break-all;
	font-weight: normal;
}

.test-btn {
	width: 100%;
	height: 80rpx;
	background: #07c160;
	color: white;
	border: none;
	border-radius: 10rpx;
	font-size: 28rpx;
	margin-top: 20rpx;
}

.test-btn::after {
	border: none;
}

.result {
	margin-top: 20rpx;
	padding: 20rpx;
	background: #f8f8f8;
	border-radius: 10rpx;
}

.result-title {
	display: block;
	font-size: 24rpx;
	color: #666;
	margin-bottom: 10rpx;
}

.result-content {
	display: block;
	font-size: 22rpx;
	color: #333;
	white-space: pre-wrap;
	word-break: break-all;
}

.action-btn {
	width: 100%;
	height: 80rpx;
	background: #f0f0f0;
	color: #666;
	border: none;
	border-radius: 10rpx;
	font-size: 28rpx;
	margin-top: 20rpx;
}

.action-btn.primary {
	background: #667eea;
	color: white;
}

.action-btn::after {
	border: none;
}
</style>
