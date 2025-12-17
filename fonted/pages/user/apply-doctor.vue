<template>
	<view class="apply-page">
		<view class="form">
			<!-- 提示信息 -->
			<view class="tip-card">
				<text class="tip-icon">ℹ️</text>
				<view class="tip-content">
					<text class="tip-title">医生入驻说明</text>
					<text class="tip-text">请如实填写以下信息，我们会在1-3个工作日内完成审核</text>
				</view>
			</view>
			
			<!-- 医生信息 -->
			<view class="form-section">
				<text class="section-title">医生信息</text>
				
				<view class="form-item">
					<text class="label">真实姓名 *</text>
					<input class="input" v-model="form.realName" placeholder="请输入真实姓名"/>
				</view>
				
				<view class="form-item">
					<text class="label">身份证号</text>
					<input class="input" v-model="form.idCard" placeholder="请输入身份证号" maxlength="18"/>
				</view>
				
				<view class="form-item">
					<text class="label">手机号 *</text>
					<input class="input" v-model="form.phone" placeholder="请输入手机号" type="number" maxlength="11"/>
				</view>
				
				<view class="form-item">
					<text class="label">职称 *</text>
					<picker mode="selector" :range="titleList" @change="onTitleChange">
						<view class="picker">{{ form.doctorTitle || '请选择' }}</view>
					</picker>
				</view>
				
				<view class="form-item">
					<text class="label">科室 *</text>
					<picker mode="selector" :range="deptList" @change="onDeptChange">
						<view class="picker">{{ form.doctorDept || '请选择' }}</view>
					</picker>
				</view>
				
				<view class="form-item">
					<text class="label">擅长领域 *</text>
					<textarea 
						class="textarea" 
						v-model="form.specialty" 
						placeholder="请输入擅长的疾病治疗领域"
						maxlength="200"
					></textarea>
				</view>
				
				<view class="form-item">
					<text class="label">个人介绍</text>
					<textarea 
						class="textarea" 
						v-model="form.introduction" 
						placeholder="请输入个人简介、工作经历等"
						maxlength="500"
					></textarea>
				</view>
			</view>
			
			<!-- 执业信息 -->
			<view class="form-section">
				<text class="section-title">执业信息</text>
				
				<view class="form-item">
					<text class="label">执业证号 *</text>
					<input class="input" v-model="form.certNumber" placeholder="请输入执业证号"/>
				</view>
				
				<view class="form-item">
					<text class="label">上传执业证书 *</text>
					<view class="cert-upload">
						<view class="cert-item" v-if="certImage">
							<image class="cert-image" :src="certImage" mode="aspectFill"></image>
							<view class="delete-btn" @click="deleteCert">×</view>
						</view>
						<view class="upload-btn" @click="chooseCert" v-else>
							<text class="upload-icon">+</text>
							<text class="upload-text">上传证书照片</text>
						</view>
					</view>
				</view>
			</view>
			
			<!-- 提交按钮 -->
			<button class="submit-btn" @click="handleSubmit" :loading="loading">
				提交申请
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
			form: {
				realName: '',
				idCard: '',
				phone: '',
				doctorTitle: '',
				doctorDept: '',
				specialty: '',
				introduction: '',
				certNumber: ''
			},
			titleList: ['主治医师', '副主任医师', '主任医师', '住院医师'],
			deptList: ['内科', '外科', '儿科', '妇产科', '骨科', '皮肤科', '眼科', '耳鼻喉科', '口腔科', '中医科'],
			certImage: '',
			certUrl: '',
			loading: false
		}
	},
	
	methods: {
		// 职称选择
		onTitleChange(e) {
			this.form.doctorTitle = this.titleList[e.detail.value]
		},
		
		// 科室选择
		onDeptChange(e) {
			this.form.doctorDept = this.deptList[e.detail.value]
		},
		
		// 选择证书
		chooseCert() {
			uni.chooseImage({
				count: 1,
				sizeType: ['compressed'],
				sourceType: ['album', 'camera'],
				success: (res) => {
					this.certImage = res.tempFilePaths[0]
				}
			})
		},
		
		// 删除证书
		deleteCert() {
			this.certImage = ''
			this.certUrl = ''
		},
		
		// 上传证书
		async uploadCert() {
			if (!this.certImage) return ''
			
			try {
				const res = await uploadFile(this.certImage, 'cert')
				return res.data.fileUrl
			} catch (error) {
				console.error('上传证书失败:', error)
				throw error
			}
		},
		
		// 提交
		async handleSubmit() {
			// 验证
			if (!this.form.realName) {
				uni.showToast({ title: '请输入真实姓名', icon: 'none' })
				return
			}
			
			if (!this.form.phone) {
				uni.showToast({ title: '请输入手机号', icon: 'none' })
				return
			}
			
			// 验证手机号格式
			if (!/^1[3-9]\d{9}$/.test(this.form.phone)) {
				uni.showToast({ title: '手机号格式不正确', icon: 'none' })
				return
			}
			
			if (!this.form.doctorTitle) {
				uni.showToast({ title: '请选择职称', icon: 'none' })
				return
			}
			
			if (!this.form.doctorDept) {
				uni.showToast({ title: '请选择科室', icon: 'none' })
				return
			}
			
			if (!this.form.specialty) {
				uni.showToast({ title: '请输入擅长领域', icon: 'none' })
				return
			}
			
			if (!this.form.certNumber) {
				uni.showToast({ title: '请输入执业证号', icon: 'none' })
				return
			}
			
			if (!this.certImage) {
				uni.showToast({ title: '请上传执业证书', icon: 'none' })
				return
			}
			
			this.loading = true
			
			try {
				// 上传证书
				uni.showLoading({ title: '上传证书中...' })
				this.certUrl = await this.uploadCert()
				uni.hideLoading()
				
				// 提交申请
				await post(API.USER_APPLY_DOCTOR, {
					...this.form,
					certImage: this.certUrl
				})
				
				uni.showToast({
					title: '提交成功',
					icon: 'success'
				})
				
				setTimeout(() => {
					uni.navigateBack()
				}, 1500)
				
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
.apply-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 20rpx 30rpx;
}

.form {
	
}

.tip-card {
	background: #fff3e0;
	border-radius: 20rpx;
	padding: 25rpx;
	margin-bottom: 20rpx;
	display: flex;
}

.tip-icon {
	font-size: 40rpx;
	margin-right: 15rpx;
}

.tip-content {
	flex: 1;
}

.tip-title {
	display: block;
	font-size: 26rpx;
	font-weight: bold;
	color: #ff9800;
	margin-bottom: 8rpx;
}

.tip-text {
	display: block;
	font-size: 24rpx;
	color: #666;
	line-height: 1.6;
}

.form-section {
	background: white;
	border-radius: 20rpx;
	padding: 30rpx;
	margin-bottom: 20rpx;
}

.section-title {
	display: block;
	font-size: 30rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 30rpx;
}

.form-item {
	margin-bottom: 30rpx;
}

.form-item:last-child {
	margin-bottom: 0;
}

.label {
	display: block;
	font-size: 28rpx;
	color: #333;
	margin-bottom: 15rpx;
}

.input {
	width: 100%;
	height: 80rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 0 20rpx;
	font-size: 28rpx;
	box-sizing: border-box;
}

.picker {
	width: 100%;
	height: 80rpx;
	line-height: 80rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 0 20rpx;
	font-size: 28rpx;
	box-sizing: border-box;
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

.cert-upload {
	
}

.cert-item {
	position: relative;
	width: 100%;
	height: 400rpx;
}

.cert-image {
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

.upload-btn {
	width: 100%;
	height: 400rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	border: 2rpx dashed #ddd;
}

.upload-icon {
	font-size: 80rpx;
	color: #999;
	margin-bottom: 15rpx;
}

.upload-text {
	font-size: 26rpx;
	color: #999;
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
