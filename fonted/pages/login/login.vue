<template>
	<view class="login-container">
		<view class="login-box">
			<view class="logo-section">
				<image class="logo" src="/static/logo.png" mode="aspectFit"></image>
				<text class="app-name">网上医疗问诊平台</text>
				<text class="app-desc">基于国密算法的安全医疗平台</text>
			</view>
			
			<view class="form-section">
				<view class="input-item">
					<text class="icon">📱</text>
					<input 
						class="input" 
						v-model="form.username" 
						placeholder="请输入用户名/手机号/邮箱"
						placeholder-style="color: #999"
					/>
				</view>
				
				<view class="input-item">
					<text class="icon">🔒</text>
					<input 
						class="input" 
						v-model="form.password" 
						type="password"
						placeholder="请输入密码"
						placeholder-style="color: #999"
					/>
				</view>
				
				<view class="captcha-item">
					<view class="captcha-input-wrapper">
						<text class="icon">🔢</text>
						<input 
							class="input" 
							v-model="form.captcha" 
							type="number"
							maxlength="4"
							placeholder="请输入验证码"
							placeholder-style="color: #999"
						/>
					</view>
					<canvas 
						canvas-id="captchaCanvas" 
						id="captchaCanvas"
						class="captcha-canvas"
						@click="refreshCaptcha"
					></canvas>
				</view>
				
				<button class="login-btn" @click="handleLogin" :loading="loading">登录</button>
				
				<view class="links">
					<text class="link" @click="goRegister">立即注册</text>
					<text class="link">忘记密码?</text>
				</view>
			</view>
			
			<view class="security-tip">
				<text class="tip-icon">🔐</text>
				<text class="tip-text">采用国密SM2/SM3/SM4算法保障数据安全</text>
			</view>
		</view>
	</view>
</template>

<script>
import { post } from '@/utils/request.js'
import { API, STORAGE_KEYS } from '@/utils/config.js'
import { encryptPassword } from '@/utils/crypto.js'
import { setStorageSync, getStorageSync } from '@/utils/storage.js'

export default {
	data() {
		return {
			form: {
				username: '',
				password: '',
				captcha: ''
			},
			loading: false,
			captchaCode: '', // 真实的验证码
			canvasContext: null
		}
	},
	
	onReady() {
		// 页面渲染完成后生成验证码
		this.initCaptcha()
	},
	methods: {
		// 初始化验证码
		initCaptcha() {
			this.canvasContext = uni.createCanvasContext('captchaCanvas', this)
			this.refreshCaptcha()
		},
		
		// 生成随机验证码
		generateCaptcha() {
			let code = ''
			for (let i = 0; i < 4; i++) {
				code += Math.floor(Math.random() * 10)
			}
			return code
		},
		
		// 刷新验证码
		refreshCaptcha() {
			if (!this.canvasContext) {
				this.canvasContext = uni.createCanvasContext('captchaCanvas', this)
			}
			
			// 生成新的验证码
			this.captchaCode = this.generateCaptcha()
			this.form.captcha = '' // 清空输入
			
			const ctx = this.canvasContext
			const width = 120
			const height = 40
			
			// 清空画布
			ctx.clearRect(0, 0, width, height)
			
			// 绘制背景
			ctx.setFillStyle('#f0f0f0')
			ctx.fillRect(0, 0, width, height)
			
			// 绘制干扰线
			for (let i = 0; i < 3; i++) {
				ctx.setStrokeStyle(this.randomColor(100, 200))
				ctx.beginPath()
				ctx.moveTo(Math.random() * width, Math.random() * height)
				ctx.lineTo(Math.random() * width, Math.random() * height)
				ctx.stroke()
			}
			
			// 绘制验证码文字
			for (let i = 0; i < this.captchaCode.length; i++) {
				const char = this.captchaCode[i]
				const x = 20 + i * 25
				const y = 25 + Math.random() * 5 - 2.5
				const rotate = (Math.random() - 0.5) * 0.3
				
				ctx.save()
				ctx.translate(x, y)
				ctx.rotate(rotate)
				ctx.setFillStyle(this.randomColor(50, 150))
				ctx.setFontSize(28)
				ctx.setTextAlign('center')
				ctx.fillText(char, 0, 0)
				ctx.restore()
			}
			
			// 绘制噪点
			for (let i = 0; i < 30; i++) {
				ctx.setFillStyle(this.randomColor(0, 255))
				ctx.beginPath()
				ctx.arc(Math.random() * width, Math.random() * height, 1, 0, 2 * Math.PI)
				ctx.fill()
			}
			
			ctx.draw()
		},
		
		// 生成随机颜色
		randomColor(min, max) {
			const r = Math.floor(Math.random() * (max - min) + min)
			const g = Math.floor(Math.random() * (max - min) + min)
			const b = Math.floor(Math.random() * (max - min) + min)
			return `rgb(${r},${g},${b})`
		},
		
		async handleLogin() {
			// 表单验证
			if (!this.form.username) {
				uni.showToast({
					title: '请输入用户名',
					icon: 'none'
				})
				return
			}
			
			if (!this.form.password) {
				uni.showToast({
					title: '请输入密码',
					icon: 'none'
				})
				return
			}
			
			// 验证码验证
			if (!this.form.captcha) {
				uni.showToast({
					title: '请输入验证码',
					icon: 'none'
				})
				return
			}
			
			if (this.form.captcha !== this.captchaCode) {
				uni.showToast({
					title: '验证码错误',
					icon: 'none'
				})
				this.refreshCaptcha() // 刷新验证码
				return
			}
			
			this.loading = true
			
			try {
				// 密码SM3加密
				const encryptedPassword = encryptPassword(this.form.password)
				
				// 调用登录API
				const res = await post(API.USER_LOGIN, {
					username: this.form.username,
					password: encryptedPassword,
					loginType: 'account'
				}, { noAuth: true })
				
				// 保存token和用户信息
				console.log('[登录响应数据]', res.data)
				setStorageSync(STORAGE_KEYS.TOKEN, res.data.token)
				setStorageSync(STORAGE_KEYS.USER_INFO, res.data.userInfo)
							
				// 验证token是否保存成功
				const savedToken = getStorageSync(STORAGE_KEYS.TOKEN)
				console.log('[Token保存验证]', savedToken ? '保存成功' : '保存失败')
				console.log('[Token内容]', savedToken)
				
				uni.showToast({
					title: '登录成功',
					icon: 'success'
				})
				
				// 跳转到首页
				setTimeout(() => {
					uni.switchTab({
						url: '/pages/index/index'
					})
				}, 1500)
				
			} catch (error) {
				// 登录失败后刷新验证码
				this.refreshCaptcha()
				console.error('登录失败:', error)
				uni.showToast({
					title: error.message || '登录失败，请检查用户名和密码',
					icon: 'none',
					duration: 2000
				})
			} finally {
				this.loading = false
			}
		},
		
		goRegister() {
			uni.navigateTo({
				url: '/pages/register/register'
			})
		}
	}
}
</script>

<style scoped>
.login-container {
	min-height: 100vh;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 40rpx;
}

.login-box {
	width: 100%;
	max-width: 600rpx;
	background: white;
	border-radius: 20rpx;
	padding: 60rpx 40rpx;
	box-shadow: 0 10rpx 40rpx rgba(0, 0, 0, 0.1);
}

.logo-section {
	text-align: center;
	margin-bottom: 60rpx;
}

.logo {
	width: 120rpx;
	height: 120rpx;
	margin-bottom: 20rpx;
}

.app-name {
	display: block;
	font-size: 36rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 10rpx;
}

.app-desc {
	display: block;
	font-size: 24rpx;
	color: #999;
}

.form-section {
	margin-bottom: 40rpx;
}

.input-item {
	display: flex;
	align-items: center;
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 0 20rpx;
	margin-bottom: 30rpx;
}

.icon {
	font-size: 40rpx;
	margin-right: 15rpx;
}

.input {
	flex: 1;
	height: 90rpx;
	font-size: 28rpx;
}

.captcha-item {
	display: flex;
	align-items: center;
	margin-bottom: 30rpx;
	gap: 20rpx;
}

.captcha-input-wrapper {
	flex: 1;
	display: flex;
	align-items: center;
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 0 20rpx;
}

.captcha-canvas {
	width: 240rpx;
	height: 80rpx;
	border-radius: 10rpx;
	border: 2rpx solid #e0e0e0;
	background: white;
}

.login-btn {
	width: 100%;
	height: 90rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
	border: none;
	border-radius: 10rpx;
	font-size: 32rpx;
	font-weight: bold;
	margin-top: 20rpx;
}

.login-btn::after {
	border: none;
}

.links {
	display: flex;
	justify-content: space-between;
	margin-top: 30rpx;
}

.link {
	font-size: 26rpx;
	color: #667eea;
}

.security-tip {
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 20rpx;
	background: #f0f9ff;
	border-radius: 10rpx;
}

.tip-icon {
	font-size: 32rpx;
	margin-right: 10rpx;
}

.tip-text {
	font-size: 22rpx;
	color: #0ea5e9;
}
</style>
