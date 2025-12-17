<template>
	<view class="edit-page">
		<view class="form">
			<!-- 头像 -->
			<view class="form-item">
				<text class="label">头像</text>
				<view class="avatar-box" @click="chooseAvatar">
					<image class="avatar" :src="form.avatar || '/static/default-avatar.png'" mode="aspectFill"></image>
					<text class="change-text">点击更换</text>
				</view>
			</view>
			
			<!-- 真实姓名 -->
			<view class="form-item">
				<text class="label">真实姓名</text>
				<input class="input" v-model="form.realName" placeholder="请输入真实姓名"/>
			</view>
			
			<!-- 性别 -->
			<view class="form-item">
				<text class="label">性别</text>
				<picker mode="selector" :range="genderList" :value="form.gender" @change="onGenderChange">
					<view class="picker">{{ genderList[form.gender] || '请选择' }}</view>
				</picker>
			</view>
			
			<!-- 出生日期 -->
			<view class="form-item">
				<text class="label">出生日期</text>
				<picker mode="date" :value="form.birthDate" @change="onDateChange">
					<view class="picker">{{ form.birthDate || '请选择' }}</view>
				</picker>
			</view>
			
			<!-- 手机号 -->
			<view class="form-item">
				<text class="label">手机号</text>
				<input class="input" v-model="form.phone" type="number" placeholder="请输入手机号" maxlength="11"/>
			</view>
			
			<!-- 邮箱 -->
			<view class="form-item">
				<text class="label">邮箱</text>
				<input class="input" v-model="form.email" placeholder="请输入邮箱地址"/>
			</view>
			
			<!-- 修改密码 -->
			<view class="form-item">
				<text class="label">修改密码</text>
				<text class="change-password" @click="showPasswordDialog">点击修改</text>
			</view>
			
			<!-- 提交按钮 -->
			<button class="submit-btn" @click="handleSubmit" :loading="loading">
				保存
			</button>
		</view>
	</view>
</template>

<script>
import { get, put } from '@/utils/request.js'
import { API, STORAGE_KEYS } from '@/utils/config.js'
import { uploadFile } from '@/utils/request.js'
import { encryptPassword, validatePasswordStrength } from '@/utils/crypto.js'
import { getStorageSync, setStorageSync } from '@/utils/storage.js'

export default {
	data() {
		return {
			form: {
				avatar: '',
				realName: '',
				gender: 0,
				birthDate: '',
				phone: '',
				email: ''
			},
			genderList: ['男', '女'],
			loading: false
		}
	},
	
	onLoad() {
		this.loadUserInfo()
	},
	
	methods: {
		// 加载用户信息
		async loadUserInfo() {
			try {
				const res = await get(API.USER_INFO)
				
				const userInfo = res.data || {}
				this.form = {
					avatar: userInfo.avatar || '',
					realName: userInfo.realName || '',
					gender: userInfo.gender || 0,
					birthDate: userInfo.birthDate || '',
					phone: userInfo.phone || '',
					email: userInfo.email || ''
				}
				
			} catch (error) {
				console.error('加载用户信息失败:', error)
			}
		},
		
		// 选择头像
		chooseAvatar() {
			uni.chooseImage({
				count: 1,
				sizeType: ['compressed'],
				sourceType: ['album', 'camera'],
				success: async (res) => {
					const filePath = res.tempFilePaths[0]
					
					uni.showLoading({ title: '上传中...' })
					
					try {
						const uploadRes = await uploadFile(filePath, 'avatar')
						console.log('[头像上传] 后端返回:', uploadRes)
						console.log('[头像上传] fileUrl:', uploadRes.data.fileUrl)
						this.form.avatar = uploadRes.data.fileUrl
						
						uni.showToast({
							title: '上传成功',
							icon: 'success'
						})
						
					} catch (error) {
						console.error('上传头像失败:', error)
					} finally {
						uni.hideLoading()
					}
				}
			})
		},
		
		// 性别选择
		onGenderChange(e) {
			this.form.gender = parseInt(e.detail.value)
		},
		
		// 日期选择
		onDateChange(e) {
			this.form.birthDate = e.detail.value
		},
		
		// 显示修改密码对话框
		showPasswordDialog() {
			uni.navigateTo({
				url: '/pages/user/change-password'
			})
		},
		
		// 提交
		async handleSubmit() {
			// 验证
			if (!this.form.realName) {
				uni.showToast({
					title: '请输入真实姓名',
					icon: 'none'
				})
				return
			}
			
			if (!this.form.phone) {
				uni.showToast({
					title: '请输入手机号',
					icon: 'none'
				})
				return
			}
			
			if (!/^1[3-9]\d{9}$/.test(this.form.phone)) {
				uni.showToast({
					title: '手机号格式不正确',
					icon: 'none'
				})
				return
			}
			
			if (!this.form.email) {
				uni.showToast({
					title: '请输入邮箱',
					icon: 'none'
				})
				return
			}
			
			if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.form.email)) {
				uni.showToast({
					title: '邮箱格式不正确',
					icon: 'none'
				})
				return
			}
			
			this.loading = true
			
			try {
				console.log('[保存资料] 提交数据:', this.form)
				const res = await put(API.USER_UPDATE, this.form)
				console.log('[保存资料] 后端响应:', res)
				
				// 更新本地用户信息
				const userInfo = getStorageSync(STORAGE_KEYS.USER_INFO)
				Object.assign(userInfo, this.form)
				setStorageSync(STORAGE_KEYS.USER_INFO, userInfo)
				
				uni.showToast({
					title: '保存成功',
					icon: 'success'
				})
				
				setTimeout(() => {
					uni.navigateBack()
				}, 1500)
				
			} catch (error) {
				console.error('保存失败:', error)
			} finally {
				this.loading = false
			}
		}
	}
}
</script>

<style scoped>
.edit-page {
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
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.label {
	font-size: 28rpx;
	color: #333;
	width: 150rpx;
}

.input {
	flex: 1;
	font-size: 28rpx;
	text-align: right;
}

.picker {
	font-size: 28rpx;
	color: #333;
}

.avatar-box {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.avatar {
	width: 150rpx;
	height: 150rpx;
	border-radius: 50%;
	margin-bottom: 15rpx;
}

.change-text {
	font-size: 24rpx;
	color: #07c160;
}

.change-password {
	font-size: 28rpx;
	color: #07c160;
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
	margin-top: 40rpx;
}

.submit-btn::after {
	border: none;
}
</style>
