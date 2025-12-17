<template>
	<view class="password-page">
		<!-- 提示信息 -->
		<view class="tip-card">
			<text class="tip-icon">🔒</text>
			<view class="tip-content">
				<text class="tip-title">密码安全提示</text>
				<text class="tip-text">• 密码长度8-20位</text>
				<text class="tip-text">• 必须包含大小写字母、数字和特殊字符</text>
				<text class="tip-text">• 建议定期修改密码以保障账户安全</text>
			</view>
		</view>
		
		<!-- 表单 -->
		<view class="form-section">
			<view class="form-item">
				<text class="form-label">旧密码</text>
				<input 
					class="form-input" 
					type="password"
					v-model="formData.oldPassword"
					placeholder="请输入当前密码"
					placeholder-class="input-placeholder"
				/>
			</view>
			
			<view class="form-item">
				<text class="form-label">新密码</text>
				<input 
					class="form-input" 
					type="password"
					v-model="formData.newPassword"
					placeholder="请输入新密码"
					placeholder-class="input-placeholder"
				/>
			</view>
			
			<view class="form-item">
				<text class="form-label">确认密码</text>
				<input 
					class="form-input" 
					type="password"
					v-model="formData.confirmPassword"
					placeholder="请再次输入新密码"
					placeholder-class="input-placeholder"
				/>
			</view>
		</view>
		
		<!-- 密码强度指示器 -->
		<view class="strength-section" v-if="formData.newPassword">
			<text class="strength-label">密码强度：</text>
			<view class="strength-bar">
				<view 
					class="strength-item" 
					:class="{ active: passwordStrength >= 1, weak: passwordStrength === 1 }"
				></view>
				<view 
					class="strength-item" 
					:class="{ active: passwordStrength >= 2, medium: passwordStrength === 2 }"
				></view>
				<view 
					class="strength-item" 
					:class="{ active: passwordStrength >= 3, strong: passwordStrength === 3 }"
				></view>
			</view>
			<text class="strength-text" :class="'strength-' + passwordStrength">
				{{ getStrengthText(passwordStrength) }}
			</text>
		</view>
		
		<!-- 提交按钮 -->
		<view class="button-section">
			<button class="submit-btn" @click="handleSubmit" :disabled="submitting">
				{{ submitting ? '提交中...' : '确认修改' }}
			</button>
		</view>
	</view>
</template>

<script>
import { put } from '@/utils/request.js'
import { API } from '@/utils/config.js'
import { encryptPassword } from '@/utils/crypto.js'

export default {
	data() {
		return {
			formData: {
				oldPassword: '',
				newPassword: '',
				confirmPassword: ''
			},
			submitting: false
		}
	},
	
	computed: {
		// 密码强度（1:弱 2:中 3:强）
		passwordStrength() {
			const pwd = this.formData.newPassword
			if (!pwd) return 0
			
			let strength = 0
			
			// 长度检查
			if (pwd.length >= 8) strength++
			
			// 复杂度检查
			const hasLower = /[a-z]/.test(pwd)
			const hasUpper = /[A-Z]/.test(pwd)
			const hasNumber = /\d/.test(pwd)
			const hasSpecial = /[!@#$%^&*(),.?":{}|<>]/.test(pwd)
			
			const complexityCount = [hasLower, hasUpper, hasNumber, hasSpecial].filter(Boolean).length
			
			if (complexityCount >= 2) strength++
			if (complexityCount >= 4) strength++
			
			return strength
		}
	},
	
	methods: {
		// 获取强度文本
		getStrengthText(strength) {
			const map = {
				0: '无',
				1: '弱',
				2: '中',
				3: '强'
			}
			return map[strength] || '无'
		},
		
		// 验证表单
		validateForm() {
			const { oldPassword, newPassword, confirmPassword } = this.formData
			
			if (!oldPassword) {
				uni.showToast({
					title: '请输入旧密码',
					icon: 'none'
				})
				return false
			}
			
			if (!newPassword) {
				uni.showToast({
					title: '请输入新密码',
					icon: 'none'
				})
				return false
			}
			
			// 密码长度检查
			if (newPassword.length < 8 || newPassword.length > 20) {
				uni.showToast({
					title: '密码长度应为8-20位',
					icon: 'none'
				})
				return false
			}
			
			// 密码强度检查
			const hasLower = /[a-z]/.test(newPassword)
			const hasUpper = /[A-Z]/.test(newPassword)
			const hasNumber = /\d/.test(newPassword)
			const hasSpecial = /[!@#$%^&*(),.?":{}|<>]/.test(newPassword)
			
			if (!(hasLower && hasUpper && hasNumber && hasSpecial)) {
				uni.showToast({
					title: '密码必须包含大小写字母、数字和特殊字符',
					icon: 'none',
					duration: 3000
				})
				return false
			}
			
			if (newPassword !== confirmPassword) {
				uni.showToast({
					title: '两次密码输入不一致',
					icon: 'none'
				})
				return false
			}
			
			if (oldPassword === newPassword) {
				uni.showToast({
					title: '新密码不能与旧密码相同',
					icon: 'none'
				})
				return false
			}
			
			return true
		},
		
		// 提交修改
		async handleSubmit() {
			if (!this.validateForm()) {
				return
			}
			
			this.submitting = true
			
			try {
				// 使用SM3加密密码
				const oldPasswordHash = encryptPassword(this.formData.oldPassword)
				const newPasswordHash = encryptPassword(this.formData.newPassword)
				
				await put(API.USER_PASSWORD, {
					oldPassword: oldPasswordHash,
					newPassword: newPasswordHash,
					confirmPassword: newPasswordHash
				})
				
				uni.showToast({
					title: '密码修改成功',
					icon: 'success'
				})
				
				// 1.5秒后返回并清空表单
				setTimeout(() => {
					this.formData = {
						oldPassword: '',
						newPassword: '',
						confirmPassword: ''
					}
					
					// 提示重新登录
					uni.showModal({
						title: '提示',
						content: '密码已修改，请重新登录',
						showCancel: false,
						success: () => {
							// 清除本地存储
							uni.removeStorageSync('token')
							uni.removeStorageSync('userInfo')
							
							// 跳转到登录页
							uni.reLaunch({
								url: '/pages/login/login'
							})
						}
					})
				}, 1500)
				
			} catch (error) {
				console.error('修改密码失败:', error)
				uni.showToast({
					title: error.message || '修改失败',
					icon: 'none'
				})
			} finally {
				this.submitting = false
			}
		}
	}
}
</script>

<style scoped>
.password-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 30rpx;
}

.tip-card {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	border-radius: 20rpx;
	padding: 40rpx;
	margin-bottom: 30rpx;
	display: flex;
	align-items: flex-start;
}

.tip-icon {
	font-size: 60rpx;
	margin-right: 25rpx;
}

.tip-content {
	flex: 1;
}

.tip-title {
	display: block;
	font-size: 30rpx;
	font-weight: bold;
	color: white;
	margin-bottom: 15rpx;
}

.tip-text {
	display: block;
	font-size: 24rpx;
	color: rgba(255, 255, 255, 0.9);
	line-height: 1.6;
	margin-bottom: 8rpx;
}

.form-section {
	background: white;
	border-radius: 20rpx;
	padding: 20rpx 0;
	margin-bottom: 30rpx;
}

.form-item {
	padding: 30rpx 30rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.form-item:last-child {
	border-bottom: none;
}

.form-label {
	display: block;
	font-size: 26rpx;
	color: #666;
	margin-bottom: 15rpx;
}

.form-input {
	width: 100%;
	height: 70rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 0 20rpx;
	font-size: 28rpx;
	color: #333;
}

.input-placeholder {
	color: #999;
}

.strength-section {
	background: white;
	border-radius: 20rpx;
	padding: 30rpx;
	margin-bottom: 30rpx;
	display: flex;
	align-items: center;
}

.strength-label {
	font-size: 26rpx;
	color: #666;
	margin-right: 15rpx;
}

.strength-bar {
	flex: 1;
	display: flex;
	gap: 10rpx;
	margin-right: 15rpx;
}

.strength-item {
	flex: 1;
	height: 8rpx;
	background: #e0e0e0;
	border-radius: 4rpx;
	transition: all 0.3s;
}

.strength-item.active.weak {
	background: #f56c6c;
}

.strength-item.active.medium {
	background: #e6a23c;
}

.strength-item.active.strong {
	background: #67c23a;
}

.strength-text {
	font-size: 24rpx;
	font-weight: bold;
}

.strength-1 {
	color: #f56c6c;
}

.strength-2 {
	color: #e6a23c;
}

.strength-3 {
	color: #67c23a;
}

.button-section {
	margin-top: 60rpx;
}

.submit-btn {
	width: 100%;
	height: 90rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
	border: none;
	border-radius: 45rpx;
	font-size: 30rpx;
	font-weight: bold;
}

.submit-btn[disabled] {
	opacity: 0.6;
}

.submit-btn::after {
	border: none;
}
</style>
