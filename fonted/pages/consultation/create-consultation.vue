<template>
	<view class="create-page">
		<view class="form">
			<!-- 选择医生 -->
			<view class="form-item">
				<text class="label">选择医生</text>
				<view class="doctor-selector" @click="selectDoctor">
					<text class="doctor-name" v-if="selectedDoctor">{{ selectedDoctor.realName }} | {{ selectedDoctor.doctorDept }}</text>
					<text class="placeholder" v-else-if="!form.autoAssign">请选择医生</text>
					<text class="placeholder ai-assign" v-else>🤖 AI智能分诊</text>
					<text class="arrow">></text>
				</view>
				
				<!-- 智能分诊选项 -->
				<view class="auto-assign-option" v-if="!selectedDoctor">
					<view class="auto-text">
						<text class="auto-title">✨ 智能自动分诊</text>
						<text class="auto-desc">根据您的症状自动匹配最合适的在线医生</text>
					</view>
					<switch :checked="form.autoAssign" @change="onAutoAssignChange" color="#1890ff"/>
				</view>
			</view>
			
			<!-- 主诉 -->
			<view class="form-item">
				<text class="label">主诉 *</text>
				<textarea 
					class="textarea" 
					v-model="form.chiefComplaint" 
					placeholder="请简要描述您的主要症状，如：头痛、发热等"
					maxlength="200"
				></textarea>
				<text class="word-count">{{ form.chiefComplaint.length }}/200</text>
			</view>
			
			<!-- 详细症状 -->
			<view class="form-item">
				<text class="label">详细症状信息</text>
			</view>
			
			<view class="symptom-row">
				<view class="symptom-item">
					<text class="symptom-label">年龄</text>
					<input class="symptom-input" v-model.number="form.symptoms.age" type="number" placeholder="岁"/>
				</view>
				<view class="symptom-item">
					<text class="symptom-label">性别</text>
					<picker mode="selector" :range="genderList" @change="onGenderChange">
						<view class="symptom-input">{{ genderList[form.symptoms.gender] || '请选择' }}</view>
					</picker>
				</view>
			</view>
			
			<view class="symptom-row">
				<view class="symptom-item">
					<text class="symptom-label">血压</text>
					<input class="symptom-input" v-model="form.symptoms.bloodPressure" placeholder="如：120/80"/>
				</view>
				<view class="symptom-item">
					<text class="symptom-label">心率</text>
					<input class="symptom-input" v-model.number="form.symptoms.heartRate" type="number" placeholder="次/分"/>
				</view>
			</view>
			
			<view class="symptom-row">
				<view class="symptom-item">
					<text class="symptom-label">体温</text>
					<input class="symptom-input" v-model.number="form.symptoms.temperature" type="digit" placeholder="℃"/>
				</view>
				<view class="symptom-item">
					<text class="symptom-label">血糖</text>
					<input class="symptom-input" v-model.number="form.symptoms.bloodSugar" type="digit" placeholder="mmol/L"/>
				</view>
			</view>
			
			<!-- 其他症状 -->
			<view class="form-item">
				<text class="label">其他症状说明</text>
				<textarea 
					class="textarea" 
					v-model="form.symptoms.otherSymptoms" 
					placeholder="请详细描述其他相关症状、持续时间、是否用药等"
					maxlength="500"
				></textarea>
				<text class="word-count">{{ (form.symptoms.otherSymptoms || '').length }}/500</text>
			</view>
			
			<!-- 上传图片 -->
			<view class="form-item">
				<text class="label">上传相关图片（可选）</text>
				<view class="image-list">
					<view class="image-item" v-for="(img, index) in imageList" :key="index">
						<image class="image" :src="img" mode="aspectFill"></image>
						<view class="delete-btn" @click="deleteImage(index)">×</view>
					</view>
					<view class="add-btn" @click="chooseImage" v-if="imageList.length < 9">
						<text class="add-icon">+</text>
						<text class="add-text">上传图片</text>
					</view>
				</view>
			</view>
			
			<!-- AI辅助 -->
			<view class="form-item">
				<view class="ai-option">
					<view class="ai-text">
						<text class="ai-title">🤖 AI智能辅助诊断</text>
						<text class="ai-desc">基于隐私计算的AI辅助分析</text>
					</view>
					<switch :checked="form.needAI" @change="onAIChange" color="#07c160"/>
				</view>
			</view>
			
			<!-- 提交按钮 -->
			<button class="submit-btn" @click="handleSubmit" :loading="loading">
				提交问诊
			</button>
		</view>
	</view>
</template>

<script>
import { post } from '@/utils/request.js'
import { API } from '@/utils/config.js'
import { uploadFile } from '@/utils/request.js'

export default {
	data() {
		return {
			selectedDoctor: null,
			form: {
				chiefComplaint: '',
				symptoms: {
					age: '',
					gender: 0,
					bloodPressure: '',
					heartRate: '',
					temperature: '',
					bloodSugar: '',
					otherSymptoms: ''
				},
				needAI: true,
				autoAssign: false  // 新增:是否启用智能分诊
			},
			genderList: ['男', '女'],
			imageList: [],
			uploadedUrls: [],
			loading: false
		}
	},
	
	onLoad(options) {
		if (options.doctorId && options.doctorName) {
			this.selectedDoctor = {
				userId: parseInt(options.doctorId),
				realName: options.doctorName,
				doctorDept: options.doctorDept || ''
			}
		}
	},
	
	methods: {
		// 选择医生
		selectDoctor() {
			if (this.form.autoAssign) {
				uni.showToast({
					title: '已启用智能分诊,无需手动选择医生',
					icon: 'none'
				})
				return
			}
			uni.switchTab({
				url: '/pages/doctors/doctors'
			})
		},
		
		// 性别选择
		onGenderChange(e) {
			this.form.symptoms.gender = parseInt(e.detail.value)
		},
		
		// AI选择
		onAIChange(e) {
			this.form.needAI = e.detail.value
			// 如果关闭AI,也关闭智能分诊
			if (!e.detail.value) {
				this.form.autoAssign = false
			}
		},
		
		// 智能分诊开关
		onAutoAssignChange(e) {
			this.form.autoAssign = e.detail.value
			// 如果启用智能分诊,自动启用AI
			if (e.detail.value) {
				this.form.needAI = true
				this.selectedDoctor = null  // 清除已选医生
			}
		},
		
		// 选择图片
		chooseImage() {
			uni.chooseImage({
				count: 9 - this.imageList.length,
				sizeType: ['compressed'],
				sourceType: ['album', 'camera'],
				success: (res) => {
					this.imageList.push(...res.tempFilePaths)
				}
			})
		},
		
		// 删除图片
		deleteImage(index) {
			this.imageList.splice(index, 1)
			this.uploadedUrls.splice(index, 1)
		},
		
		// 上传图片
		async uploadImages() {
			const urls = []
			
			for (let i = 0; i < this.imageList.length; i++) {
				try {
					const res = await uploadFile(this.imageList[i], 'consultation')
					urls.push(res.data.fileUrl)
				} catch (error) {
					console.error('图片上传失败:', error)
				}
			}
			
			return urls
		},
		
		// 提交
		async handleSubmit() {
			// 验证:如果未启用智能分诊,必须选择医生
			if (!this.selectedDoctor && !this.form.autoAssign) {
				uni.showToast({
					title: '请选择医生或启用智能分诊',
					icon: 'none'
				})
				return
			}
			
			if (!this.form.chiefComplaint) {
				uni.showToast({
					title: '请填写主诉',
					icon: 'none'
				})
				return
			}
			
			this.loading = true
			
			try {
				// 上传图片
				if (this.imageList.length > 0) {
					uni.showLoading({ title: '上传图片中...' })
					this.uploadedUrls = await this.uploadImages()
					uni.hideLoading()
				}
				
				// 整理症状数据
				const symptoms = {
					...this.form.symptoms,
					images: this.uploadedUrls
				}
				
				// 调用API
				const res = await post(API.CONSULTATION_CREATE, {
					doctorId: this.form.autoAssign ? null : this.selectedDoctor.userId,  // 智能分诊时不传doctorId
					chiefComplaint: this.form.chiefComplaint,
					symptoms: symptoms,
					needAI: this.form.needAI
				})
				
				// 显示分诊结果
				if (res.data.autoAssigned) {
					const doctor = res.data.assignedDoctor
					uni.showModal({
						title: '智能分诊成功',
						content: `已为您分配医生:
${doctor.doctorName} | ${doctor.doctorDept}
职称:${doctor.doctorTitle}

分配原因:${res.data.assignedReason}`,
						showCancel: false,
						confirmText: '查看详情',
						success: (modalRes) => {
							if (modalRes.confirm) {
								uni.redirectTo({
									url: '/pages/consultation/consultation-detail?id=' + res.data.consultationId
								})
							}
						}
					})
				} else {
					uni.showToast({
						title: '提交成功',
						icon: 'success'
					})
					
					// 跳转到问诊详情
					setTimeout(() => {
						uni.redirectTo({
							url: '/pages/consultation/consultation-detail?id=' + res.data.consultationId
						})
					}, 1500)
				}
				
			} catch (error) {
				console.error('提交失败:', error)
			} finally {
				this.loading = false
			}
		}
	}
}
</script>

<style scoped>
.create-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 20rpx 30rpx;
}

.form {
	
}

.form-item {
	background: white;
	border-radius: 20rpx;
	padding: 30rpx;
	margin-bottom: 20rpx;
}

.label {
	display: block;
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 20rpx;
}

.doctor-selector {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 25rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
}

.doctor-name {
	font-size: 28rpx;
	color: #333;
}

.placeholder {
	font-size: 28rpx;
	color: #999;
}

.arrow {
	font-size: 28rpx;
	color: #999;
}

.textarea {
	width: 100%;
	min-height: 150rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 20rpx;
	font-size: 28rpx;
	box-sizing: border-box;
}

.word-count {
	display: block;
	text-align: right;
	font-size: 24rpx;
	color: #999;
	margin-top: 10rpx;
}

.symptom-row {
	display: flex;
	margin-bottom: 20rpx;
}

.symptom-row:last-child {
	margin-bottom: 0;
}

.symptom-item {
	flex: 1;
	margin-right: 20rpx;
}

.symptom-item:last-child {
	margin-right: 0;
}

.symptom-label {
	display: block;
	font-size: 24rpx;
	color: #666;
	margin-bottom: 10rpx;
}

.symptom-input {
	width: 100%;
	height: 70rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 0 20rpx;
	font-size: 26rpx;
	box-sizing: border-box;
	line-height: 70rpx;
}

.image-list {
	display: flex;
	flex-wrap: wrap;
}

.image-item {
	position: relative;
	width: 200rpx;
	height: 200rpx;
	margin-right: 20rpx;
	margin-bottom: 20rpx;
}

.image {
	width: 100%;
	height: 100%;
	border-radius: 10rpx;
}

.delete-btn {
	position: absolute;
	top: -10rpx;
	right: -10rpx;
	width: 50rpx;
	height: 50rpx;
	background: #f56c6c;
	color: white;
	border-radius: 50%;
	text-align: center;
	line-height: 50rpx;
	font-size: 40rpx;
}

.add-btn {
	width: 200rpx;
	height: 200rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
}

.add-icon {
	font-size: 60rpx;
	color: #999;
	margin-bottom: 10rpx;
}

.add-text {
	font-size: 24rpx;
	color: #999;
}

.ai-option {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.ai-text {
	flex: 1;
}

.ai-title {
	display: block;
	font-size: 28rpx;
	color: #333;
	margin-bottom: 5rpx;
}

.ai-desc {
	display: block;
	font-size: 24rpx;
	color: #999;
}

.auto-assign-option {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-top: 20rpx;
	padding: 20rpx;
	background: linear-gradient(135deg, #e6f7ff 0%, #f0f9ff 100%);
	border-radius: 10rpx;
	border: 1px solid #91d5ff;
}

.auto-text {
	flex: 1;
}

.auto-title {
	display: block;
	font-size: 28rpx;
	color: #1890ff;
	margin-bottom: 5rpx;
	font-weight: bold;
}

.auto-desc {
	display: block;
	font-size: 24rpx;
	color: #666;
}

.ai-assign {
	color: #1890ff;
	font-weight: bold;
}

.submit-btn {
	width: 100%;
	height: 90rpx;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
	border: none;
	border-radius: 45rpx;
	font-size: 32rpx;
	font-weight: bold;
	margin-top: 20rpx;
}

.submit-btn::after {
	border: none;
}
</style>
