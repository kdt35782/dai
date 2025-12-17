<template>
	<view class="detail-page">
		<!-- 问诊信息卡片 -->
		<view class="info-card">
			<view class="card-header">
				<text class="card-title">问诊信息</text>
				<text class="status" :class="'status-' + consultationInfo.status">
					{{ consultationInfo.statusText }}
				</text>
			</view>
			
			<view class="info-item">
				<text class="info-label">主诉：</text>
				<text class="info-value">{{ consultationInfo.chiefComplaint }}</text>
			</view>
			
			<view class="info-item" v-if="consultationInfo.symptoms">
				<text class="info-label">症状信息：</text>
				<view class="symptom-detail">
					<text v-if="consultationInfo.symptoms.age">年龄：{{ consultationInfo.symptoms.age }}岁　</text>
					<text v-if="consultationInfo.symptoms.gender !== undefined">性别：{{ consultationInfo.symptoms.gender === 1 ? '男' : '女' }}　</text>
					<text v-if="consultationInfo.symptoms.bloodPressure">血压：{{ consultationInfo.symptoms.bloodPressure }}　</text>
					<text v-if="consultationInfo.symptoms.heartRate">心率：{{ consultationInfo.symptoms.heartRate }}次/分　</text>
					<text v-if="consultationInfo.symptoms.temperature">体温：{{ consultationInfo.symptoms.temperature }}℃　</text>
					<text v-if="consultationInfo.symptoms.bloodSugar">血糖：{{ consultationInfo.symptoms.bloodSugar }}mmol/L</text>
				</view>
			</view>
			
			<view class="info-item" v-if="consultationInfo.symptoms && consultationInfo.symptoms.otherSymptoms">
				<text class="info-label">其他症状：</text>
				<text class="info-value">{{ consultationInfo.symptoms.otherSymptoms }}</text>
			</view>
			
			<view class="info-item" v-if="consultationInfo.needAI">
				<text class="ai-tag">🤖 已启用AI辅助诊断</text>
			</view>
		</view>
		
		<!-- AI诊断结果卡片 -->
		<view class="ai-card" v-if="consultationInfo.aiDiagnosis">
			<view class="card-header">
				<text class="card-title">🤖 AI智能诊断</text>
				<view class="risk-badge" :class="'risk-' + consultationInfo.aiDiagnosis.riskLevel">
					<text class="risk-text">{{ getRiskText(consultationInfo.aiDiagnosis.riskLevel) }}</text>
				</view>
			</view>
			
			<view class="ai-content">
				<view class="score-section">
					<text class="score-label">风险评分：</text>
					<text class="score-value" :class="'score-' + consultationInfo.aiDiagnosis.riskLevel">
						{{ consultationInfo.aiDiagnosis.riskScore }}
					</text>
					<text class="score-max">/100</text>
				</view>
				
				<view class="diagnosis-section">
					<text class="section-label">初步分析：</text>
					<text class="section-value">{{ consultationInfo.aiDiagnosis.diagnosis }}</text>
				</view>
				
				<view class="suggestion-section">
					<text class="section-label">就医建议：</text>
					<text class="section-value">{{ consultationInfo.aiDiagnosis.suggestions }}</text>
				</view>
				
				<!-- 分系统详细分析 -->
				<view class="analysis-section" v-if="consultationInfo.aiDiagnosis.detailedAnalysis">
					<text class="section-label">详细分析：</text>
					<view class="analysis-item" v-for="(analysis, system) in consultationInfo.aiDiagnosis.detailedAnalysis" :key="system">
						<text class="analysis-system">{{ system }}：</text>
						<text class="analysis-content">{{ analysis }}</text>
					</view>
				</view>
				
				<!-- 生活方式建议 -->
				<view class="lifestyle-section" v-if="consultationInfo.aiDiagnosis.lifestyleAdvice && consultationInfo.aiDiagnosis.lifestyleAdvice.length > 0">
					<text class="section-label">生活建议：</text>
					<view class="lifestyle-item" v-for="(advice, index) in consultationInfo.aiDiagnosis.lifestyleAdvice" :key="index">
						<text class="advice-bullet">•</text>
						<text class="advice-content">{{ advice }}</text>
					</view>
				</view>
				
				<!-- 复诊建议 -->
				<view class="followup-section" v-if="consultationInfo.aiDiagnosis.followUpAdvice">
					<text class="section-label">复诊建议：</text>
					<text class="section-value">{{ consultationInfo.aiDiagnosis.followUpAdvice }}</text>
				</view>
				
				<!-- 推荐科室 -->
				<view class="dept-section" v-if="consultationInfo.aiDiagnosis.recommendedDept && consultationInfo.aiDiagnosis.recommendedDept !== '全科'">
					<text class="section-label">推荐科室：</text>
					<text class="section-value dept-highlight">{{ consultationInfo.aiDiagnosis.recommendedDept }}</text>
				</view>
			</view>
		</view>
		
		<!-- 聊天区域 -->
		<scroll-view 
			class="chat-area" 
			scroll-y
			:scroll-into-view="scrollToView"
			scroll-with-animation
		>
			<view 
				class="message-item" 
				:class="msg.role === 'patient' ? 'right' : 'left'"
				v-for="(msg, index) in messages" 
				:key="index"
				:id="'msg-' + index"
			>
				<image class="avatar" :src="msg.avatar || '/static/default-avatar.png'" mode="aspectFill"></image>
				
				<view class="message-content">
					<text class="message-text" v-if="msg.type === 'text'">{{ msg.content }}</text>
					
					<image class="message-image" v-if="msg.type === 'image'" :src="msg.content" mode="aspectFill" @click="previewImage(msg.content)"></image>
					
					<view class="prescription-card" v-if="msg.type === 'prescription'" @click="viewPrescription(msg.prescriptionId)">
						<text class="prescription-icon">📋</text>
						<view class="prescription-info">
							<text class="prescription-title">电子处方</text>
							<text class="prescription-desc">点击查看详情</text>
						</view>
					</view>
					
					<text class="message-time">{{ msg.createdAt }}</text>
				</view>
			</view>
		</scroll-view>
		
		<!-- 输入栏 -->
		<view class="input-bar" v-if="consultationInfo.status === 1">
			<view class="input-box">
				<text class="add-btn" @click="showActionSheet">+</text>
				<input 
					class="input" 
					v-model="inputText" 
					placeholder="输入消息..."
					@confirm="sendMessage"
				/>
				<button class="send-btn" @click="sendMessage" :disabled="!inputText.trim()">
					发送
				</button>
			</view>
		</view>
		
		<!-- 操作按钮 -->
		<view class="action-bar" v-if="isDoctor && consultationInfo.status === 0">
			<button class="accept-btn" @click="acceptConsultation">
				接诊
			</button>
		</view>
		
		<view class="action-bar" v-if="consultationInfo.status === 1">
			<button class="chat-btn" @click="enterChat">
				💬 进入聊天室
			</button>
			<button class="finish-btn" v-if="isDoctor" @click="showFinishOptions">
				完成问诊
			</button>
		</view>
	</view>
</template>

<script>
import { get, post } from '@/utils/request.js'
import { API, STORAGE_KEYS } from '@/utils/config.js'
import { uploadFile } from '@/utils/request.js'
import { getStorageSync } from '@/utils/storage.js'

export default {
	data() {
		return {
			consultationId: 0,
			consultationInfo: {},
			messages: [],
			inputText: '',
			scrollToView: '',
			isDoctor: false,
			refreshTimer: null,
			finishForm: {
				diagnosis: ''
			}
		}
	},
	
	onLoad(options) {
		if (options.id) {
			this.consultationId = parseInt(options.id)
			this.checkRole()
			this.loadDetail()
			
			// 启动定时刷新（每5秒）
			this.refreshTimer = setInterval(() => {
				this.loadMessages(false)
			}, 5000)
		}
	},
	
	onUnload() {
		if (this.refreshTimer) {
			clearInterval(this.refreshTimer)
		}
	},
	
	methods: {
		// 检查角色
		checkRole() {
			const userInfo = getStorageSync(STORAGE_KEYS.USER_INFO)
			this.isDoctor = userInfo && userInfo.role === 'doctor'
			console.log('[问诊详情] 用户角色:', userInfo?.role, '是否为医生:', this.isDoctor)
		},
		
		// 加载问诊详情
		async loadDetail() {
			try {
				const res = await get(API.CONSULTATION_DETAIL, {
					consultationId: this.consultationId
				})
				
				this.consultationInfo = res.data || {}
				this.loadMessages()
				
			} catch (error) {
				console.error('加载问诊详情失败:', error)
			}
		},
		
		// 加载消息
		async loadMessages(showLoading = true) {
			try {
				// TODO: 调用获取消息列表API
				// 这里暂时使用模拟数据
				// const res = await get(API.CONSULTATION_MESSAGES, {
				//   consultationId: this.consultationId
				// })
				// this.messages = res.data.list || []
				
				// 滚动到底部
				if (this.messages.length > 0) {
					this.scrollToView = 'msg-' + (this.messages.length - 1)
				}
				
			} catch (error) {
				console.error('加载消息失败:', error)
			}
		},
		
		// 显示操作菜单
		showActionSheet() {
			uni.showActionSheet({
				itemList: ['发送图片'],
				success: (res) => {
					if (res.tapIndex === 0) {
						this.chooseImage()
					}
				}
			})
		},
		
		// 选择图片
		chooseImage() {
			uni.chooseImage({
				count: 1,
				sizeType: ['compressed'],
				sourceType: ['album', 'camera'],
				success: async (res) => {
					const filePath = res.tempFilePaths[0]
					
					uni.showLoading({ title: '发送中...' })
					
					try {
						const uploadRes = await uploadFile(filePath, 'chat')
						
						// TODO: 发送图片消息
						// await this.sendImageMessage(uploadRes.data.fileUrl)
						
						this.loadMessages()
						
					} catch (error) {
						console.error('发送图片失败:', error)
					} finally {
						uni.hideLoading()
					}
				}
			})
		},
		
		// 预览图片
		previewImage(url) {
			uni.previewImage({
				urls: [url],
				current: url
			})
		},
		
		// 查看处方
		viewPrescription(prescriptionId) {
			uni.navigateTo({
				url: `/pages/prescription/detail?id=${prescriptionId}`
			})
		},
		
		// 进入聊天室
		enterChat() {
			// 调试信息
			const userInfo = getStorageSync(STORAGE_KEYS.USER_INFO);
			const token = getStorageSync(STORAGE_KEYS.TOKEN);
					
			console.log('[问诊详情] 准备进入聊天室');
			console.log('[问诊详情] consultationId:', this.consultationId);
			console.log('[问诊详情] token:', token ? '存在' : '不存在');
			console.log('[问诊详情] userInfo:', userInfo);
					
			if (!userInfo || !userInfo.userId) {
				uni.showToast({
					title: '用户信息失效，请重新登录',
					icon: 'none'
				});
				setTimeout(() => {
					uni.redirectTo({ url: '/pages/login/login' });
				}, 1500);
				return;
			}
					
			uni.navigateTo({
				url: `/pages/chat/index?consultationId=${this.consultationId}`
			});
		},
		
		// 发送消息
		async sendMessage() {
			if (!this.inputText.trim()) return
			
			const content = this.inputText.trim()
			this.inputText = ''
			
			try {
				// TODO: 调用发送消息API
				// await post(API.CONSULTATION_SEND_MESSAGE, {
				//   consultationId: this.consultationId,
				//   content: content,
				//   type: 'text'
				// })
				
				this.loadMessages()
				
			} catch (error) {
				console.error('发送消息失败:', error)
			}
		},
		
		// 接诊
		async acceptConsultation() {
			try {
				await post(API.CONSULTATION_ACCEPT, {
					consultationId: this.consultationId
				})
				
				uni.showToast({
					title: '接诊成功',
					icon: 'success'
				})
				
				this.loadDetail()
				
			} catch (error) {
				console.error('接诊失败:', error)
			}
		},
		
		// 完成问诊 - 显示选项
		showFinishOptions() {
			uni.showActionSheet({
				title: '请选择完成方式',
				itemList: ['仅填写诊断意见', '填写诊断并开处方'],
				success: (res) => {
					if (res.tapIndex === 0) {
						// 仅诊断
						this.inputDiagnosis(false)
					} else if (res.tapIndex === 1) {
						// 诊断+处方
						this.inputDiagnosis(true)
					}
				}
			})
		},
		
		// 输入诊断意见
		inputDiagnosis(needPrescription) {
			uni.showModal({
				title: '诊断意见',
				editable: true,
				placeholderText: '请输入诊断意见',
				success: async (res) => {
					if (res.confirm) {
						const diagnosis = res.content?.trim()
						
						if (!diagnosis) {
							uni.showToast({
								title: '请输入诊断意见',
								icon: 'none'
							})
							return
						}
						
						// 保存诊断
						this.finishForm.diagnosis = diagnosis
						
						if (needPrescription) {
							// 需要开处方,跳转到处方页面
							uni.navigateTo({
								url: `/pages/prescription/create?consultationId=${this.consultationId}&diagnosis=${encodeURIComponent(diagnosis)}`
							})
						} else {
							// 不需要处方,直接完成
							this.submitFinish(diagnosis, null)
						}
					}
				}
			})
		},
		
		// 提交完成问诊
		async submitFinish(diagnosis, medicines) {
			try {
				await post(API.CONSULTATION_FINISH, {
					consultationId: this.consultationId,
					diagnosis: diagnosis,
					prescription: medicines
				})
				
				uni.showToast({
					title: '问诊已完成',
					icon: 'success'
				})
				
				this.loadDetail()
				
			} catch (error) {
				console.error('完成问诊失败:', error)
				uni.showToast({
					title: error.message || '操作失败',
					icon: 'none'
				})
			}
		},
		
		// 获取风险等级文本
		getRiskText(level) {
			const riskMap = {
				'high': '高风险',
				'medium': '中等风险',
				'low': '低风险',
				'normal': '正常'
			}
			return riskMap[level] || '未知'
		}
	}
}
</script>

<style scoped>
.detail-page {
	height: 100vh;
	display: flex;
	flex-direction: column;
	background: #f5f5f5;
}

.info-card {
	background: white;
	padding: 25rpx 30rpx;
	margin-bottom: 10rpx;
}

.card-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20rpx;
}

.card-title {
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

.info-item {
	margin-bottom: 15rpx;
}

.info-item:last-child {
	margin-bottom: 0;
}

.info-label {
	font-size: 26rpx;
	color: #666;
}

.info-value {
	font-size: 26rpx;
	color: #333;
}

.symptom-detail {
	font-size: 24rpx;
	color: #666;
	line-height: 1.8;
	margin-top: 5rpx;
}

.ai-tag {
	display: inline-block;
	font-size: 24rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
	padding: 8rpx 20rpx;
	border-radius: 20rpx;
}

/* AI诊断卡片样式 */
.ai-card {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	margin: 10rpx 30rpx;
	border-radius: 20rpx;
	padding: 30rpx;
	box-shadow: 0 8rpx 20rpx rgba(102, 126, 234, 0.3);
}

.ai-card .card-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 25rpx;
}

.ai-card .card-title {
	font-size: 30rpx;
	font-weight: bold;
	color: white;
}

.risk-badge {
	padding: 8rpx 20rpx;
	border-radius: 30rpx;
	font-size: 24rpx;
}

.risk-high {
	background: rgba(244, 67, 54, 0.9);
}

.risk-medium {
	background: rgba(255, 152, 0, 0.9);
}

.risk-low {
	background: rgba(76, 175, 80, 0.9);
}

.risk-normal {
	background: rgba(33, 150, 243, 0.9);
}

.risk-text {
	color: white;
	font-weight: bold;
}

.ai-content {
	background: rgba(255, 255, 255, 0.95);
	border-radius: 15rpx;
	padding: 25rpx;
}

.score-section {
	display: flex;
	align-items: baseline;
	margin-bottom: 20rpx;
	padding-bottom: 20rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.score-label {
	font-size: 26rpx;
	color: #666;
}

.score-value {
	font-size: 60rpx;
	font-weight: bold;
	margin: 0 10rpx;
}

.score-high {
	color: #f44336;
}

.score-medium {
	color: #ff9800;
}

.score-low {
	color: #4caf50;
}

.score-normal {
	color: #2196f3;
}

.score-max {
	font-size: 28rpx;
	color: #999;
}

.diagnosis-section, .suggestion-section, .analysis-section, .lifestyle-section, .followup-section, .dept-section {
	margin-bottom: 15rpx;
}

.section-label {
	display: block;
	font-size: 26rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 8rpx;
}

.section-value {
	font-size: 26rpx;
	color: #666;
	line-height: 1.8;
}

.dept-highlight {
	color: #07c160;
	font-weight: bold;
}

.analysis-item {
	margin-bottom: 10rpx;
}

.analysis-system {
	font-weight: bold;
	color: #333;
}

.analysis-content {
	color: #666;
}

.lifestyle-item {
	display: flex;
	margin-bottom: 8rpx;
}

.advice-bullet {
	color: #07c160;
	margin-right: 10rpx;
}

.advice-content {
	color: #666;
	flex: 1;
}

.chat-area {
	flex: 1;
	padding: 20rpx 30rpx;
}

.message-item {
	display: flex;
	margin-bottom: 30rpx;
}

.message-item.left {
	flex-direction: row;
}

.message-item.right {
	flex-direction: row-reverse;
}

.avatar {
	width: 70rpx;
	height: 70rpx;
	border-radius: 50%;
	flex-shrink: 0;
}

.message-content {
	max-width: 500rpx;
	margin: 0 20rpx;
}

.message-item.left .message-content {
	margin-left: 20rpx;
	margin-right: 0;
}

.message-item.right .message-content {
	margin-right: 20rpx;
	margin-left: 0;
	display: flex;
	flex-direction: column;
	align-items: flex-end;
}

.message-text {
	display: inline-block;
	padding: 20rpx;
	background: white;
	border-radius: 10rpx;
	font-size: 28rpx;
	color: #333;
	line-height: 1.6;
}

.message-item.right .message-text {
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
}

.message-image {
	width: 300rpx;
	height: 300rpx;
	border-radius: 10rpx;
}

.prescription-card {
	display: flex;
	align-items: center;
	padding: 20rpx;
	background: white;
	border-radius: 10rpx;
	border: 2rpx solid #07c160;
}

.prescription-icon {
	font-size: 50rpx;
	margin-right: 15rpx;
}

.prescription-info {
	display: flex;
	flex-direction: column;
}

.prescription-title {
	font-size: 28rpx;
	color: #333;
	margin-bottom: 5rpx;
}

.prescription-desc {
	font-size: 22rpx;
	color: #999;
}

.message-time {
	display: block;
	font-size: 22rpx;
	color: #999;
	margin-top: 10rpx;
}

.input-bar {
	background: white;
	padding: 20rpx 30rpx;
	border-top: 1px solid #f0f0f0;
}

.input-box {
	display: flex;
	align-items: center;
}

.add-btn {
	width: 60rpx;
	height: 60rpx;
	line-height: 60rpx;
	text-align: center;
	font-size: 40rpx;
	color: #666;
	margin-right: 15rpx;
}

.input {
	flex: 1;
	height: 70rpx;
	background: #f5f5f5;
	border-radius: 35rpx;
	padding: 0 25rpx;
	font-size: 28rpx;
}

.send-btn {
	width: 120rpx;
	height: 70rpx;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
	border: none;
	border-radius: 35rpx;
	font-size: 28rpx;
	margin-left: 15rpx;
	padding: 0;
	line-height: 70rpx;
}

.send-btn::after {
	border: none;
}

.send-btn[disabled] {
	background: #e0e0e0;
	color: #999;
}

.action-bar {
	background: white;
	padding: 20rpx 30rpx;
	border-top: 1px solid #f0f0f0;
	display: flex;
	gap: 20rpx;
}

.accept-btn,
.chat-btn,
.prescription-btn,
.finish-btn {
	flex: 1;
	height: 90rpx;
	color: white;
	border: none;
	border-radius: 45rpx;
	font-size: 32rpx;
	font-weight: bold;
}

.accept-btn,
.finish-btn {
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
}

.chat-btn {
	background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
}

.prescription-btn {
	background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
}

.accept-btn::after,
.chat-btn::after,
.prescription-btn::after,
.finish-btn::after {
	border: none;
}
</style>
