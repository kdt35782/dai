<template>
	<view class="register-container">
		<view class="register-box">
			<view class="header">
				<text class="title">欢迎注册</text>
				<text class="subtitle">国密问诊平台</text>
			</view>
			
			<view class="form">
				<!-- 用户类型选择 -->
				<view class="input-item">
					<text class="label">注册类型</text>
					<view class="role-tabs">
						<view 
							class="role-tab" 
							:class="{ active: form.role === 'patient' }"
							@click="form.role = 'patient'"
						>
							<text class="role-icon">👤</text>
							<text class="role-name">患者</text>
						</view>
						<view 
							class="role-tab" 
							:class="{ active: form.role === 'doctor' }"
							@click="form.role = 'doctor'"
						>
							<text class="role-icon">👨‍⚕️</text>
							<text class="role-name">医生</text>
						</view>
					</view>
				</view>
				
				<view class="input-item">
					<text class="label">用户名</text>
					<input 
						class="input" 
						v-model="form.username" 
						placeholder="4-20个字符，字母数字下划线"
						maxlength="20"
					/>
				</view>
				
				<view class="input-item">
					<text class="label">手机号</text>
					<input 
						class="input" 
						v-model="form.phone" 
						type="number"
						placeholder="请输入手机号"
						maxlength="11"
					/>
				</view>
				
				<view class="input-item">
					<text class="label">邮箱</text>
					<input 
						class="input" 
						v-model="form.email" 
						placeholder="请输入邮箱地址"
					/>
				</view>
				
				<view class="input-item">
					<text class="label">密码</text>
					<input 
						class="input" 
						v-model="form.password" 
						type="password"
						placeholder="8-20位，含大小写字母数字特殊字符"
					/>
					<text class="strength" :class="strengthClass">{{ strengthText }}</text>
				</view>
				
				<view class="input-item">
					<text class="label">确认密码</text>
					<input 
						class="input" 
						v-model="form.confirmPassword" 
						type="password"
						placeholder="请再次输入密码"
					/>
				</view>
				
				<!-- 医生专属字段 -->
				<template v-if="form.role === 'doctor'">
					<view class="doctor-section-title">
						<text class="section-icon">🏥</text>
						<text>医生专业信息</text>
					</view>
					
					<view class="input-item">
						<text class="label">真实姓名 *</text>
						<input 
							class="input" 
							v-model="form.realName" 
							placeholder="请输入真实姓名"
						/>
					</view>
					
					<view class="input-item">
						<text class="label">身份证号</text>
						<input 
							class="input" 
							v-model="form.idCard" 
							placeholder="请输入身份证号"
							maxlength="18"
						/>
					</view>
					
					<view class="input-item">
						<text class="label">职称 *</text>
						<picker mode="selector" :range="titleList" @change="onTitleChange">
							<view class="picker">{{ form.doctorTitle || '请选择职称' }}</view>
						</picker>
					</view>
					
					<view class="input-item">
						<text class="label">科室 *</text>
						<picker mode="selector" :range="deptList" @change="onDeptChange">
							<view class="picker">{{ form.doctorDept || '请选择科室' }}</view>
						</picker>
					</view>
					
					<view class="input-item">
						<text class="label">擅长领域 *</text>
						<textarea 
							class="textarea" 
							v-model="form.specialty" 
							placeholder="请输入擅长的疾病治疗领域"
							maxlength="200"
						></textarea>
					</view>
					
					<view class="input-item">
						<text class="label">个人介绍</text>
						<textarea 
							class="textarea" 
							v-model="form.introduction" 
							placeholder="请输入个人简介、工作经历等"
							maxlength="500"
						></textarea>
					</view>
					
					<view class="input-item">
						<text class="label">执业证号 *</text>
						<input 
							class="input" 
							v-model="form.certNumber" 
							placeholder="请输入执业证号"
						/>
					</view>
					
					<view class="input-item">
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
				</template>
				
				<view class="agreement">
					<checkbox-group @change="onAgreementChange">
						<checkbox value="agree" :checked="agreed" />
					</checkbox-group>
					<text class="agreement-text">
						我已阅读并同意
						<text class="link">《用户协议》</text>
						和
						<text class="link">《隐私政策》</text>
					</text>
				</view>
				
				<button class="register-btn" @click="handleRegister" :loading="loading">
					注册
				</button>
				
				<view class="footer">
					<text class="tip">已有账号？</text>
					<text class="link" @click="goLogin">立即登录</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
import { post } from '@/utils/request.js'
import { API } from '@/utils/config.js'
import { encryptPassword, validatePasswordStrength } from '@/utils/crypto.js'

export default {
	data() {
		return {
			form: {
				role: 'patient', // 默认患者
				username: '',
				phone: '',
				email: '',
				password: '',
				confirmPassword: '',
				// 医生专属字段
				realName: '',
				idCard: '',
				doctorTitle: '',
				doctorDept: '',
				specialty: '',
				introduction: '',
				certNumber: ''
			},
			titleList: ['主治医师', '副主任医师', '主任医师', '住院医师'],
			deptList: ['内科', '外科', '儿科', '妇产科', '骨科', '皮肤科', '眼科', '耳鼻喉科', '口腔科', '中医科'],
			certImage: '',
			agreed: false,
			loading: false
		}
	},
	computed: {
		passwordStrength() {
			return validatePasswordStrength(this.form.password)
		},
		strengthClass() {
			if (!this.form.password) return ''
			return this.passwordStrength.valid ? 'strong' : 'weak'
		},
		strengthText() {
			if (!this.form.password) return ''
			return this.passwordStrength.message
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
		},
		
		// 上传证书
		async uploadCert() {
			if (!this.certImage) return ''
			
			try {
				// 这里使用base64编码作为临时方案
				// 实际项目中应该调用文件上传接口
				return this.certImage
			} catch (error) {
				console.error('上传证书失败:', error)
				throw error
			}
		},
		
		// 同意协议
		onAgreementChange(e) {
			this.agreed = e.detail.value.includes('agree')
		},
		
		// 注册
		async handleRegister() {
			// 基础表单验证
			if (!this.form.username) {
				uni.showToast({ title: '请输入用户名', icon: 'none' })
				return
			}
			
			if (!/^[a-zA-Z0-9_]{4,20}$/.test(this.form.username)) {
				uni.showToast({ 
					title: '用户名格式不正确', 
					icon: 'none' 
				})
				return
			}
			
			if (!this.form.phone) {
				uni.showToast({ title: '请输入手机号', icon: 'none' })
				return
			}
			
			if (!/^1[3-9]\d{9}$/.test(this.form.phone)) {
				uni.showToast({ title: '手机号格式不正确', icon: 'none' })
				return
			}
			
			if (!this.form.email) {
				uni.showToast({ title: '请输入邮箱', icon: 'none' })
				return
			}
			
			if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.form.email)) {
				uni.showToast({ title: '邮箱格式不正确', icon: 'none' })
				return
			}
			
			if (!this.passwordStrength.valid) {
				uni.showToast({ 
					title: this.passwordStrength.message, 
					icon: 'none' 
				})
				return
			}
			
			if (this.form.password !== this.form.confirmPassword) {
				uni.showToast({ title: '两次密码不一致', icon: 'none' })
				return
			}
			
			// 医生专属字段验证
			if (this.form.role === 'doctor') {
				if (!this.form.realName) {
					uni.showToast({ title: '请输入真实姓名', icon: 'none' })
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
			}
			
			if (!this.agreed) {
				uni.showToast({ title: '请同意用户协议', icon: 'none' })
				return
			}
			
			this.loading = true
			
			try {
				// 密码SM3加密
				const encryptedPassword = encryptPassword(this.form.password)
				
				// 准备注册数据
				const registerData = {
					username: this.form.username,
					password: encryptedPassword,
					email: this.form.email,
					phone: this.form.phone,
					role: this.form.role
				}
				
				// 如果是医生注册，添加医生信息
				if (this.form.role === 'doctor') {
					// 上传证书
					uni.showLoading({ title: '上传证书中...' })
					const certUrl = await this.uploadCert()
					uni.hideLoading()
					
					registerData.realName = this.form.realName
					registerData.idCard = this.form.idCard
					registerData.doctorTitle = this.form.doctorTitle
					registerData.doctorDept = this.form.doctorDept
					registerData.specialty = this.form.specialty
					registerData.introduction = this.form.introduction
					registerData.certNumber = this.form.certNumber
					registerData.certImage = certUrl
				}
				
				// 调用注册API
				const res = await post(API.USER_REGISTER, registerData, { noAuth: true })
				
				uni.showToast({
					title: '注册成功',
					icon: 'success'
				})
				
				// 跳转到登录页
				setTimeout(() => {
					uni.navigateBack()
				}, 1500)
				
			} catch (error) {
				console.error('注册失败:', error)
			} finally {
				this.loading = false
			}
		},
		
		goLogin() {
			uni.navigateBack()
		}
	}
}
</script>

<style scoped>
.register-container {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 40rpx;
}

.register-box {
	background: white;
	border-radius: 20rpx;
	padding: 40rpx;
}

.header {
	text-align: center;
	margin-bottom: 40rpx;
}

.title {
	display: block;
	font-size: 40rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 10rpx;
}

.subtitle {
	display: block;
	font-size: 26rpx;
	color: #999;
}

.form {
	
}

.input-item {
	margin-bottom: 30rpx;
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

.strength {
	display: block;
	font-size: 24rpx;
	margin-top: 10rpx;
}

.strength.weak {
	color: #f56c6c;
}

.strength.strong {
	color: #67c23a;
}

.code-item {
	position: relative;
}

.code-input {
	width: calc(100% - 200rpx);
}

.code-btn {
	position: absolute;
	right: 0;
	bottom: 0;
	width: 180rpx;
	height: 80rpx;
	line-height: 80rpx;
	background: #07c160;
	color: white;
	border: none;
	border-radius: 10rpx;
	font-size: 24rpx;
	padding: 0;
}

.code-btn::after {
	border: none;
}

.code-btn[disabled] {
	background: #e0e0e0;
	color: #999;
}

.agreement {
	display: flex;
	align-items: center;
	margin: 30rpx 0;
}

.agreement-text {
	font-size: 24rpx;
	color: #666;
	margin-left: 10rpx;
}

.link {
	color: #07c160;
}

.register-btn {
	width: 100%;
	height: 90rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
	border: none;
	border-radius: 10rpx;
	font-size: 32rpx;
	font-weight: bold;
	margin: 30rpx 0;
}

.register-btn::after {
	border: none;
}

.footer {
	text-align: center;
	font-size: 26rpx;
}

.tip {
	color: #999;
}

.role-tabs {
	display: flex;
	gap: 20rpx;
}

.role-tab {
	flex: 1;
	height: 120rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	border: 2rpx solid transparent;
	transition: all 0.3s;
}

.role-tab.active {
	background: #e8f4ff;
	border-color: #409eff;
}

.role-icon {
	font-size: 40rpx;
	margin-bottom: 5rpx;
}

.role-name {
	font-size: 26rpx;
	color: #666;
}

.role-tab.active .role-name {
	color: #409eff;
	font-weight: bold;
}

.doctor-section-title {
	display: flex;
	align-items: center;
	font-size: 28rpx;
	font-weight: bold;
	color: #409eff;
	margin: 30rpx 0 20rpx 0;
	padding-bottom: 15rpx;
	border-bottom: 2rpx solid #e0e0e0;
}

.section-icon {
	margin-right: 10rpx;
	font-size: 32rpx;
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
	min-height: 120rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 20rpx;
	font-size: 28rpx;
	box-sizing: border-box;
}

.cert-upload {
	margin-top: 10rpx;
}

.cert-item {
	position: relative;
	width: 100%;
	height: 300rpx;
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
	height: 300rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	border: 2rpx dashed #ddd;
}

.upload-icon {
	font-size: 60rpx;
	color: #999;
	margin-bottom: 10rpx;
}

.upload-text {
	font-size: 24rpx;
	color: #999;
}
</style>
