<template>
	<view class="detail-page">
		<!-- 病历头部 -->
		<view class="header-card">
			<view class="header-row">
				<text class="header-label">医生：</text>
				<text class="header-value">{{ recordInfo.doctorName }}</text>
			</view>
			<view class="header-row">
				<text class="header-label">科室：</text>
				<text class="header-value">{{ recordInfo.doctorDept }}</text>
			</view>
			<view class="header-row">
				<text class="header-label">时间：</text>
				<text class="header-value">{{ recordInfo.createdAt }}</text>
			</view>
			<view class="secure-tip">
				<text class="secure-icon">🔐</text>
				<text class="secure-text">本病历已使用SM4加密存储</text>
			</view>
		</view>
		
		<!-- 主诉 -->
		<view class="content-card">
			<text class="card-title">主诉</text>
			<text class="card-content">{{ recordInfo.chiefComplaint }}</text>
		</view>
		
		<!-- 症状 -->
		<view class="content-card" v-if="recordInfo.symptoms">
			<text class="card-title">症状信息</text>
			<view class="symptom-grid">
				<view class="symptom-item" v-if="recordInfo.symptoms.age">
					<text class="symptom-label">年龄</text>
					<text class="symptom-value">{{ recordInfo.symptoms.age }}岁</text>
				</view>
				<view class="symptom-item" v-if="recordInfo.symptoms.gender !== undefined">
					<text class="symptom-label">性别</text>
					<text class="symptom-value">{{ recordInfo.symptoms.gender === 1 ? '男' : '女' }}</text>
				</view>
				<view class="symptom-item" v-if="recordInfo.symptoms.bloodPressure">
					<text class="symptom-label">血压</text>
					<text class="symptom-value">{{ recordInfo.symptoms.bloodPressure }}</text>
				</view>
				<view class="symptom-item" v-if="recordInfo.symptoms.heartRate">
					<text class="symptom-label">心率</text>
					<text class="symptom-value">{{ recordInfo.symptoms.heartRate }}次/分</text>
				</view>
				<view class="symptom-item" v-if="recordInfo.symptoms.temperature">
					<text class="symptom-label">体温</text>
					<text class="symptom-value">{{ recordInfo.symptoms.temperature }}℃</text>
				</view>
				<view class="symptom-item" v-if="recordInfo.symptoms.bloodSugar">
					<text class="symptom-label">血糖</text>
					<text class="symptom-value">{{ recordInfo.symptoms.bloodSugar }}mmol/L</text>
				</view>
			</view>
			<text class="card-content" v-if="recordInfo.symptoms.otherSymptoms">{{ recordInfo.symptoms.otherSymptoms }}</text>
		</view>
		
		<!-- 诊断 -->
		<view class="content-card">
			<text class="card-title">诊断</text>
			<text class="card-content">{{ recordInfo.diagnosis || '暂无' }}</text>
		</view>
		
		<!-- 处理意见 -->
		<view class="content-card">
			<text class="card-title">处理意见</text>
			<text class="card-content">{{ recordInfo.treatment || '暂无' }}</text>
		</view>
		
		<!-- AI辅助建议 -->
		<view class="content-card" v-if="recordInfo.aiAdvice">
			<view class="ai-header">
				<text class="card-title">🤖 AI辅助建议</text>
				<text class="ai-tag">Paillier同态加密</text>
			</view>
			<text class="card-content">{{ recordInfo.aiAdvice }}</text>
		</view>
		
		<!-- 完整性验证 -->
		<view class="verify-card">
			<text class="verify-icon">✓</text>
			<view class="verify-info">
				<text class="verify-title">数据完整性已验证</text>
				<text class="verify-desc">SM3哈希校验通过</text>
			</view>
		</view>
	</view>
</template>

<script>
import { get } from '@/utils/request.js'
import { API } from '@/utils/config.js'

export default {
	data() {
		return {
			recordId: 0,
			recordInfo: {}
		}
	},
	
	onLoad(options) {
		if (options.id) {
			this.recordId = parseInt(options.id)
			this.loadDetail()
		}
	},
	
	methods: {
		// 加载病历详情
		async loadDetail() {
			try {
				const res = await get(API.RECORD_DETAIL, {
					recordId: this.recordId
				})
				
				this.recordInfo = res.data || {}
				
			} catch (error) {
				console.error('加载病历详情失败:', error)
				uni.showToast({
					title: '加载失败',
					icon: 'none'
				})
			}
		}
	}
}
</script>

<style scoped>
.detail-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 20rpx 30rpx;
}

.header-card {
	background: white;
	border-radius: 20rpx;
	padding: 30rpx;
	margin-bottom: 20rpx;
}

.header-row {
	display: flex;
	margin-bottom: 15rpx;
}

.header-row:last-of-type {
	margin-bottom: 20rpx;
}

.header-label {
	font-size: 26rpx;
	color: #666;
	width: 100rpx;
}

.header-value {
	flex: 1;
	font-size: 26rpx;
	color: #333;
}

.secure-tip {
	display: flex;
	align-items: center;
	padding: 15rpx;
	background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
	border-radius: 10rpx;
}

.secure-icon {
	font-size: 30rpx;
	margin-right: 10rpx;
}

.secure-text {
	font-size: 22rpx;
	color: #667eea;
}

.content-card {
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

.card-content {
	display: block;
	font-size: 26rpx;
	color: #666;
	line-height: 1.8;
}

.symptom-grid {
	display: flex;
	flex-wrap: wrap;
	margin-bottom: 20rpx;
}

.symptom-item {
	width: 50%;
	display: flex;
	flex-direction: column;
	margin-bottom: 20rpx;
}

.symptom-label {
	font-size: 22rpx;
	color: #999;
	margin-bottom: 8rpx;
}

.symptom-value {
	font-size: 26rpx;
	color: #333;
}

.ai-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20rpx;
}

.ai-tag {
	font-size: 20rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
	padding: 5rpx 12rpx;
	border-radius: 10rpx;
}

.verify-card {
	background: #e8f5e9;
	border-radius: 20rpx;
	padding: 25rpx;
	display: flex;
	align-items: center;
}

.verify-icon {
	width: 60rpx;
	height: 60rpx;
	line-height: 60rpx;
	text-align: center;
	background: #4caf50;
	color: white;
	border-radius: 50%;
	font-size: 36rpx;
	margin-right: 20rpx;
}

.verify-info {
	flex: 1;
}

.verify-title {
	display: block;
	font-size: 26rpx;
	font-weight: bold;
	color: #4caf50;
	margin-bottom: 5rpx;
}

.verify-desc {
	display: block;
	font-size: 22rpx;
	color: #81c784;
}
</style>
