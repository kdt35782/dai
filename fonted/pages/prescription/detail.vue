<template>
	<view class="detail-page">
		<!-- 处方头部 -->
		<view class="prescription-header">
			<view class="header-title">📋 电子处方</view>
			<view class="prescription-no">处方编号: {{ prescriptionInfo.prescriptionNo }}</view>
			<view class="prescription-date">开具时间: {{ prescriptionInfo.createdAt }}</view>
		</view>
		
		<!-- 医患信息 -->
		<view class="info-card">
			<view class="card-header">
				<text class="card-title">医患信息</text>
			</view>
			
			<view class="info-row">
				<text class="info-label">医生:</text>
				<text class="info-value">{{ prescriptionInfo.doctorName }}</text>
			</view>
			<view class="info-row">
				<text class="info-label">患者:</text>
				<text class="info-value">{{ prescriptionInfo.patientName }}</text>
			</view>
		</view>
		
		<!-- 诊断信息 -->
		<view class="diagnosis-card">
			<view class="card-header">
				<text class="card-title">诊断结果</text>
			</view>
			<view class="diagnosis-content">{{ prescriptionInfo.diagnosis }}</view>
		</view>
		
		<!-- 药品清单 -->
		<view class="medicines-card">
			<view class="card-header">
				<text class="card-title">药品清单</text>
			</view>
			
			<view 
				class="medicine-item" 
				v-for="(item, index) in prescriptionInfo.details" 
				:key="index"
			>
				<view class="medicine-header">
					<text class="medicine-index">{{ index + 1 }}</text>
					<text class="medicine-name">{{ item.medicineName }}</text>
				</view>
				
				<view class="medicine-spec">{{ item.specification }}</view>
				
				<view class="medicine-usage">
					<view class="usage-row">
						<text class="usage-label">用法:</text>
						<text class="usage-value">{{ item.usage || '-' }}</text>
					</view>
					<view class="usage-row">
						<text class="usage-label">频次:</text>
						<text class="usage-value">{{ item.frequency || '-' }}</text>
					</view>
					<view class="usage-row">
						<text class="usage-label">剂量:</text>
						<text class="usage-value">{{ item.dosage || '-' }}</text>
					</view>
					<view class="usage-row">
						<text class="usage-label">疗程:</text>
						<text class="usage-value">{{ item.duration || '-' }}</text>
					</view>
					<view class="usage-row">
						<text class="usage-label">数量:</text>
						<text class="usage-value">{{ item.quantity }}{{ item.unit }}</text>
					</view>
					<view class="usage-row" v-if="item.notes">
						<text class="usage-label">备注:</text>
						<text class="usage-value">{{ item.notes }}</text>
					</view>
				</view>
				
				<view class="medicine-price">
					<text class="price-label">金额:</text>
					<text class="price-value">{{ item.totalPrice.toFixed(2) }}元</text>
				</view>
			</view>
			
			<!-- 总金额 -->
			<view class="total-amount">
				<text class="total-label">总金额:</text>
				<text class="total-value">{{ prescriptionInfo.totalAmount?.toFixed(2) }}元</text>
			</view>
		</view>
		
		<!-- 用药注意事项 -->
		<view class="notice-card">
			<view class="card-header">
				<text class="card-title">⚠️ 用药注意事项</text>
			</view>
			
			<view class="notice-item">
				<text class="notice-bullet">•</text>
				<text class="notice-text">请严格按照医嘱用药，不可自行增减剂量</text>
			</view>
			<view class="notice-item">
				<text class="notice-bullet">•</text>
				<text class="notice-text">如有不适，请及时联系医生或就医</text>
			</view>
			<view class="notice-item">
				<text class="notice-bullet">•</text>
				<text class="notice-text">请注意药品的保质期和储存条件</text>
			</view>
			<view class="notice-item">
				<text class="notice-bullet">•</text>
				<text class="notice-text">孕妇、哺乳期妇女及儿童用药需遵医嘱</text>
			</view>
		</view>
		
		<!-- 数字签名 -->
		<view class="signature-card">
			<view class="signature-info">
				<text class="signature-label">数字签名:</text>
				<text class="signature-value">{{ prescriptionInfo.digitalSignature || '已验证' }}</text>
			</view>
			<view class="signature-verify">
				<text class="verify-icon">✓</text>
				<text class="verify-text">本处方已通过数字签名验证</text>
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
			prescriptionId: 0,
			prescriptionInfo: {
				details: []
			}
		}
	},
	
	onLoad(options) {
		if (options.id) {
			this.prescriptionId = parseInt(options.id)
			this.loadDetail()
		}
	},
	
	methods: {
		// 加载处方详情
		async loadDetail() {
			uni.showLoading({ title: '加载中...' })
			
			try {
				// 使用路径参数请求
				const res = await get(`${API.PRESCRIPTION_DETAIL}/${this.prescriptionId}`)
				
				this.prescriptionInfo = res.data || { details: [] }
				
			} catch (error) {
				console.error('加载处方详情失败:', error)
				uni.showToast({
					title: '加载失败',
					icon: 'none'
				})
			} finally {
				uni.hideLoading()
			}
		}
	}
}
</script>

<style scoped>
.detail-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding-bottom: 30rpx;
}

.prescription-header {
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	padding: 40rpx 30rpx;
	color: white;
}

.header-title {
	font-size: 36rpx;
	font-weight: bold;
	margin-bottom: 15rpx;
}

.prescription-no {
	font-size: 24rpx;
	opacity: 0.9;
	margin-bottom: 5rpx;
}

.prescription-date {
	font-size: 24rpx;
	opacity: 0.9;
}

.info-card,
.diagnosis-card,
.medicines-card,
.notice-card,
.signature-card {
	background: white;
	margin: 20rpx 30rpx;
	border-radius: 15rpx;
	padding: 30rpx;
}

.card-header {
	margin-bottom: 20rpx;
}

.card-title {
	font-size: 30rpx;
	font-weight: bold;
	color: #333;
}

.info-row {
	display: flex;
	margin-bottom: 15rpx;
}

.info-label {
	width: 120rpx;
	font-size: 26rpx;
	color: #666;
}

.info-value {
	flex: 1;
	font-size: 26rpx;
	color: #333;
}

.diagnosis-content {
	font-size: 26rpx;
	color: #333;
	line-height: 1.8;
	padding: 20rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
}

/* 药品清单 */
.medicine-item {
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 25rpx;
	margin-bottom: 20rpx;
}

.medicine-header {
	display: flex;
	align-items: center;
	margin-bottom: 10rpx;
}

.medicine-index {
	width: 40rpx;
	height: 40rpx;
	line-height: 40rpx;
	text-align: center;
	background: #07c160;
	color: white;
	border-radius: 50%;
	font-size: 22rpx;
	font-weight: bold;
	margin-right: 15rpx;
}

.medicine-name {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
}

.medicine-spec {
	font-size: 22rpx;
	color: #999;
	margin-bottom: 15rpx;
}

.medicine-usage {
	background: white;
	border-radius: 8rpx;
	padding: 15rpx;
	margin-bottom: 15rpx;
}

.usage-row {
	display: flex;
	margin-bottom: 10rpx;
}

.usage-row:last-child {
	margin-bottom: 0;
}

.usage-label {
	width: 100rpx;
	font-size: 24rpx;
	color: #666;
}

.usage-value {
	flex: 1;
	font-size: 24rpx;
	color: #333;
}

.medicine-price {
	display: flex;
	justify-content: flex-end;
	align-items: center;
}

.price-label {
	font-size: 24rpx;
	color: #666;
	margin-right: 10rpx;
}

.price-value {
	font-size: 28rpx;
	color: #07c160;
	font-weight: bold;
}

.total-amount {
	display: flex;
	justify-content: flex-end;
	align-items: center;
	padding-top: 20rpx;
	border-top: 1rpx solid #f0f0f0;
}

.total-label {
	font-size: 28rpx;
	color: #333;
	font-weight: bold;
	margin-right: 15rpx;
}

.total-value {
	font-size: 36rpx;
	color: #07c160;
	font-weight: bold;
}

/* 用药注意事项 */
.notice-item {
	display: flex;
	margin-bottom: 15rpx;
}

.notice-bullet {
	color: #07c160;
	margin-right: 10rpx;
	font-size: 26rpx;
}

.notice-text {
	flex: 1;
	font-size: 24rpx;
	color: #666;
	line-height: 1.6;
}

/* 数字签名 */
.signature-info {
	display: flex;
	align-items: center;
	margin-bottom: 15rpx;
}

.signature-label {
	font-size: 24rpx;
	color: #666;
	margin-right: 10rpx;
}

.signature-value {
	font-size: 22rpx;
	color: #999;
	font-family: monospace;
}

.signature-verify {
	display: flex;
	align-items: center;
	padding: 15rpx;
	background: #e8f5e9;
	border-radius: 8rpx;
}

.verify-icon {
	width: 40rpx;
	height: 40rpx;
	line-height: 40rpx;
	text-align: center;
	background: #4caf50;
	color: white;
	border-radius: 50%;
	font-size: 26rpx;
	font-weight: bold;
	margin-right: 15rpx;
}

.verify-text {
	font-size: 24rpx;
	color: #4caf50;
}
</style>
