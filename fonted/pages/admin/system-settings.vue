<template>
	<view class="settings-page">
		<!-- 系统信息 -->
		<view class="section">
			<view class="section-title">📊 系统信息</view>
			<view class="setting-card">
				<view class="setting-item">
					<text class="setting-label">系统名称</text>
					<text class="setting-value">基于国密加密的网上看诊系统</text>
				</view>
				<view class="setting-item">
					<text class="setting-label">系统版本</text>
					<text class="setting-value">v1.0.0</text>
				</view>
				<view class="setting-item">
					<text class="setting-label">数据库版本</text>
					<text class="setting-value">MySQL 8.0</text>
				</view>
				<view class="setting-item">
					<text class="setting-label">运行状态</text>
					<view class="status-indicator">
						<view class="status-dot status-online"></view>
						<text class="status-text">运行中</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 加密配置 -->
		<view class="section">
			<view class="section-title">🔐 国密加密配置</view>
			<view class="setting-card">
				<view class="setting-item">
					<text class="setting-label">SM2 非对称加密</text>
					<view class="status-indicator">
						<view class="status-dot status-online"></view>
						<text class="status-text">已启用</text>
					</view>
				</view>
				<view class="setting-item">
					<text class="setting-label">SM3 哈希算法</text>
					<view class="status-indicator">
						<view class="status-dot status-online"></view>
						<text class="status-text">已启用</text>
					</view>
				</view>
				<view class="setting-item">
					<text class="setting-label">SM4 对称加密</text>
					<view class="status-indicator">
						<view class="status-dot status-online"></view>
						<text class="status-text">已启用</text>
					</view>
				</view>
				<view class="setting-item">
					<text class="setting-label">Paillier 同态加密</text>
					<view class="status-indicator">
						<view class="status-dot status-online"></view>
						<text class="status-text">已启用</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 安全设置 -->
		<view class="section">
			<view class="section-title">🛡️ 安全设置</view>
			<view class="setting-card">
				<view class="setting-item clickable" @click="showFeature('token')">
					<text class="setting-label">Token 有效期</text>
					<view class="setting-right">
						<text class="setting-value">2 小时</text>
						<text class="arrow">›</text>
					</view>
				</view>
				<view class="setting-item clickable" @click="showFeature('password')">
					<text class="setting-label">密码强度要求</text>
					<view class="setting-right">
						<text class="setting-value">高</text>
						<text class="arrow">›</text>
					</view>
				</view>
				<view class="setting-item clickable" @click="showFeature('login')">
					<text class="setting-label">登录失败锁定</text>
					<view class="setting-right">
						<switch :checked="loginLockEnabled" @change="toggleLoginLock" color="#ff6b6b" />
					</view>
				</view>
				<view class="setting-item clickable" @click="showFeature('log')">
					<text class="setting-label">操作日志记录</text>
					<view class="setting-right">
						<switch :checked="true" disabled color="#ff6b6b" />
					</view>
				</view>
			</view>
		</view>
		
		<!-- 数据管理 -->
		<view class="section">
			<view class="section-title">💾 数据管理</view>
			<view class="setting-card">
				<view class="setting-item">
					<text class="setting-label">数据存储量</text>
					<text class="setting-value">{{ formatSize(dataSize) }}</text>
				</view>
				<view class="setting-item">
					<text class="setting-label">备份频率</text>
					<text class="setting-value">每天 02:00</text>
				</view>
				<view class="setting-item clickable" @click="viewBackupLogs">
					<text class="setting-label">最近备份</text>
					<view class="setting-right">
						<text class="setting-value">2024-12-05 02:00</text>
						<text class="arrow">›</text>
					</view>
				</view>
				<view class="setting-item clickable" @click="confirmBackup">
					<text class="setting-label">立即备份</text>
					<view class="setting-right">
						<text class="action-text">执行</text>
						<text class="arrow">›</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 系统维护 -->
		<view class="section">
			<view class="section-title">🔧 系统维护</view>
			<view class="setting-card">
				<view class="setting-item clickable" @click="clearCache">
					<text class="setting-label">清理缓存</text>
					<view class="setting-right">
						<text class="action-text">清理</text>
						<text class="arrow">›</text>
					</view>
				</view>
				<view class="setting-item clickable" @click="viewSystemLogs">
					<text class="setting-label">系统日志</text>
					<view class="setting-right">
						<text class="setting-value">查看</text>
						<text class="arrow">›</text>
					</view>
				</view>
				<view class="setting-item clickable" @click="viewErrorLogs">
					<text class="setting-label">错误日志</text>
					<view class="setting-right">
						<view class="badge-dot" v-if="errorCount > 0"></view>
						<text class="setting-value">{{ errorCount }} 条</text>
						<text class="arrow">›</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 通知设置 -->
		<view class="section">
			<view class="section-title">🔔 通知设置</view>
			<view class="setting-card">
				<view class="setting-item clickable">
					<text class="setting-label">新用户注册通知</text>
					<switch :checked="notifySettings.newUser" @change="toggleNotify('newUser')" color="#ff6b6b" />
				</view>
				<view class="setting-item clickable">
					<text class="setting-label">医生申请通知</text>
					<switch :checked="notifySettings.doctorApp" @change="toggleNotify('doctorApp')" color="#ff6b6b" />
				</view>
				<view class="setting-item clickable">
					<text class="setting-label">系统异常通知</text>
					<switch :checked="notifySettings.systemError" @change="toggleNotify('systemError')" color="#ff6b6b" />
				</view>
			</view>
		</view>
		
		<!-- 关于 -->
		<view class="section">
			<view class="section-title">ℹ️ 关于</view>
			<view class="setting-card">
				<view class="setting-item clickable" @click="viewDocs">
					<text class="setting-label">开发文档</text>
					<text class="arrow">›</text>
				</view>
				<view class="setting-item clickable" @click="viewLicense">
					<text class="setting-label">许可协议</text>
					<text class="arrow">›</text>
				</view>
				<view class="setting-item">
					<text class="setting-label">技术支持</text>
					<text class="setting-value">admin@example.com</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
export default {
	data() {
		return {
			loginLockEnabled: true,
			dataSize: 1024 * 1024 * 256, // 256MB
			errorCount: 3,
			notifySettings: {
				newUser: true,
				doctorApp: true,
				systemError: true
			}
		}
	},
	
	methods: {
		// 格式化文件大小
		formatSize(bytes) {
			if (bytes < 1024) return bytes + ' B'
			if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
			if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
			return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
		},
		
		// 显示功能（待开发）
		showFeature(type) {
			uni.showToast({
				title: '功能开发中',
				icon: 'none'
			})
		},
		
		// 切换登录锁定
		toggleLoginLock(e) {
			this.loginLockEnabled = e.detail.value
			uni.showToast({
				title: this.loginLockEnabled ? '已启用登录锁定' : '已关闭登录锁定',
				icon: 'success'
			})
		},
		
		// 切换通知
		toggleNotify(type) {
			this.notifySettings[type] = !this.notifySettings[type]
		},
		
		// 查看备份日志
		viewBackupLogs() {
			uni.showToast({
				title: '功能开发中',
				icon: 'none'
			})
		},
		
		// 确认备份
		confirmBackup() {
			uni.showModal({
				title: '确认备份',
				content: '确定要立即备份数据库吗？',
				success: (res) => {
					if (res.confirm) {
						this.executeBackup()
					}
				}
			})
		},
		
		// 执行备份
		executeBackup() {
			uni.showLoading({
				title: '备份中...'
			})
			
			setTimeout(() => {
				uni.hideLoading()
				uni.showToast({
					title: '备份成功',
					icon: 'success'
				})
			}, 2000)
		},
		
		// 清理缓存
		clearCache() {
			uni.showModal({
				title: '确认清理',
				content: '确定要清理系统缓存吗？',
				success: (res) => {
					if (res.confirm) {
						uni.showToast({
							title: '缓存已清理',
							icon: 'success'
						})
					}
				}
			})
		},
		
		// 查看系统日志
		viewSystemLogs() {
			uni.navigateTo({
				url: '/pages/admin/system-logs?tab=login'
			})
		},
		
		// 查看错误日志
		viewErrorLogs() {
			uni.navigateTo({
				url: '/pages/admin/system-logs?tab=error'
			})
		},
		
		// 查看文档
		viewDocs() {
			uni.showToast({
				title: '功能开发中',
				icon: 'none'
			})
		},
		
		// 查看许可
		viewLicense() {
			uni.showToast({
				title: '功能开发中',
				icon: 'none'
			})
		}
	}
}
</script>

<style scoped>
.settings-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 20rpx;
	padding-bottom: 40rpx;
}

.section {
	margin-bottom: 30rpx;
}

.section-title {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	padding: 0 10rpx 15rpx;
}

.setting-card {
	background: white;
	border-radius: 15rpx;
	overflow: hidden;
}

.setting-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 30rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.setting-item:last-child {
	border-bottom: none;
}

.setting-item.clickable {
	cursor: pointer;
}

.setting-label {
	font-size: 28rpx;
	color: #333;
}

.setting-value {
	font-size: 26rpx;
	color: #999;
}

.setting-right {
	display: flex;
	align-items: center;
	gap: 15rpx;
}

.arrow {
	font-size: 32rpx;
	color: #ccc;
}

.action-text {
	font-size: 26rpx;
	color: #ff6b6b;
}

/* 状态指示器 */
.status-indicator {
	display: flex;
	align-items: center;
	gap: 10rpx;
}

.status-dot {
	width: 16rpx;
	height: 16rpx;
	border-radius: 50%;
}

.status-online {
	background: #4caf50;
	box-shadow: 0 0 10rpx rgba(76, 175, 80, 0.5);
}

.status-offline {
	background: #999;
}

.status-text {
	font-size: 26rpx;
	color: #4caf50;
}

/* 徽章点 */
.badge-dot {
	width: 12rpx;
	height: 12rpx;
	background: #ff6b6b;
	border-radius: 50%;
}
</style>
