<template>
	<view class="logs-page">
		<!-- 日志类型切换 -->
		<view class="log-tabs">
			<view 
				class="log-tab" 
				:class="{ active: currentTab === 'login' }"
				@click="switchTab('login')"
			>
				<text class="tab-icon">🔐</text>
				<text class="tab-text">登录日志</text>
			</view>
			<view 
				class="log-tab" 
				:class="{ active: currentTab === 'operation' }"
				@click="switchTab('operation')"
			>
				<text class="tab-icon">⚙️</text>
				<text class="tab-text">操作日志</text>
			</view>
			<view 
				class="log-tab" 
				:class="{ active: currentTab === 'error' }"
				@click="switchTab('error')"
			>
				<text class="tab-icon">⚠️</text>
				<text class="tab-text">错误日志</text>
				<view class="error-badge" v-if="errorCount > 0">{{ errorCount }}</view>
			</view>
		</view>
		
		<!-- 筛选栏 -->
		<view class="filter-bar" v-if="currentTab === 'login'">
			<scroll-view class="status-tabs" scroll-x>
				<view 
					class="status-tab" 
					:class="{ active: currentStatus === item.value }"
					v-for="item in statusOptions" 
					:key="item.value"
					@click="selectStatus(item.value)"
				>
					{{ item.label }}
				</view>
			</scroll-view>
		</view>
		
		<!-- 登录日志列表 -->
		<scroll-view 
			class="log-list" 
			scroll-y
			@scrolltolower="loadMore"
			v-if="currentTab === 'login'"
		>
			<view 
				class="log-item" 
				v-for="item in logList" 
				:key="item.logId"
			>
				<view class="log-header">
					<view class="user-info">
						<text class="username">{{ item.username }}</text>
						<text class="status-tag" :class="'status-' + item.status">
							{{ item.statusText }}
						</text>
					</view>
					<text class="log-time">{{ formatTime(item.loginTime) }}</text>
				</view>
				
				<view class="log-details">
					<view class="detail-row">
						<text class="detail-label">IP地址：</text>
						<text class="detail-value">{{ item.loginIp || '未知' }}</text>
					</view>
					<view class="detail-row">
						<text class="detail-label">位置：</text>
						<text class="detail-value">{{ item.loginLocation || '未知' }}</text>
					</view>
					<view class="detail-row">
						<text class="detail-label">系统：</text>
						<text class="detail-value">{{ item.os || '未知' }}</text>
					</view>
					<view class="detail-row">
						<text class="detail-label">浏览器：</text>
						<text class="detail-value">{{ item.browser || '未知' }}</text>
					</view>
					<view class="detail-row" v-if="item.msg">
						<text class="detail-label">信息：</text>
						<text class="detail-value">{{ item.msg }}</text>
					</view>
				</view>
			</view>
			
			<!-- 加载状态 -->
			<view class="loading" v-if="loading">加载中...</view>
			<view class="no-more" v-if="!hasMore && logList.length > 0">没有更多了</view>
			<view class="empty" v-if="!loading && logList.length === 0">
				<text class="empty-icon">📋</text>
				<text class="empty-text">暂无日志记录</text>
			</view>
		</scroll-view>
		
		<!-- 操作日志列表 -->
		<scroll-view 
			class="log-list" 
			scroll-y
			v-if="currentTab === 'operation'"
		>
			<view class="log-item" v-for="item in operationLogs" :key="item.id">
				<view class="log-header">
					<view class="user-info">
						<text class="username">{{ item.username }}</text>
						<text class="operation-type">{{ item.operationType }}</text>
					</view>
					<text class="log-time">{{ item.createTime }}</text>
				</view>
				
				<view class="log-content">
					<text class="content-text">{{ item.content }}</text>
				</view>
				
				<view class="log-details" v-if="item.details">
					<view class="detail-row">
						<text class="detail-label">模块：</text>
						<text class="detail-value">{{ item.module }}</text>
					</view>
					<view class="detail-row">
						<text class="detail-label">IP：</text>
						<text class="detail-value">{{ item.ip }}</text>
					</view>
				</view>
			</view>
			
			<view class="empty" v-if="operationLogs.length === 0">
				<text class="empty-icon">⚙️</text>
				<text class="empty-text">暂无操作日志</text>
			</view>
		</scroll-view>
		
		<!-- 错误日志列表 -->
		<scroll-view 
			class="log-list" 
			scroll-y
			v-if="currentTab === 'error'"
		>
			<view class="log-item error-log" v-for="item in errorLogs" :key="item.id">
				<view class="log-header">
					<view class="user-info">
						<text class="error-level" :class="'level-' + item.level">
							{{ item.levelText }}
						</text>
						<text class="error-module">{{ item.module }}</text>
					</view>
					<text class="log-time">{{ item.createTime }}</text>
				</view>
				
				<view class="error-message">
					<text class="message-text">{{ item.message }}</text>
				</view>
				
				<view class="error-stack" v-if="item.stack">
					<text class="stack-title">堆栈信息：</text>
					<text class="stack-text">{{ item.stack }}</text>
				</view>
			</view>
			
			<view class="empty" v-if="errorLogs.length === 0">
				<text class="empty-icon">✅</text>
				<text class="empty-text">暂无错误日志</text>
			</view>
		</scroll-view>
	</view>
</template>

<script>
import { get } from '@/utils/request.js'
import { API } from '@/utils/config.js'

export default {
	data() {
		return {
			currentTab: 'login',
			currentStatus: '',
			statusOptions: [
				{ label: '全部', value: '' },
				{ label: '成功', value: 1 },
				{ label: '失败', value: 0 }
			],
			logList: [],
			operationLogs: [],
			errorLogs: [],
			errorCount: 0,
			page: 1,
			pageSize: 10,
			loading: false,
			hasMore: true
		}
	},
	
	onLoad(options) {
		// 从参数中获取要显示的tab
		if (options.tab) {
			this.currentTab = options.tab
		}
		
		this.loadLogs()
	},
	
	methods: {
		// 切换标签页
		switchTab(tab) {
			this.currentTab = tab
			this.page = 1
			this.loadLogs()
		},
		
		// 选择状态
		selectStatus(status) {
			this.currentStatus = status
			this.page = 1
			this.loadLogs()
		},
		
		// 加载日志
		async loadLogs() {
			if (this.currentTab === 'login') {
				await this.loadLoginLogs()
			} else if (this.currentTab === 'operation') {
				await this.loadOperationLogs()
			} else if (this.currentTab === 'error') {
				await this.loadErrorLogs()
			}
		},
		
		// 加载登录日志
		async loadLoginLogs() {
			if (this.loading) return
			
			this.loading = true
			
			try {
				const params = {
					page: this.page,
					pageSize: this.pageSize
				}
				
				if (this.currentStatus !== '') {
					params.status = this.currentStatus
				}
				
				const res = await get(API.USER_ADMIN_LOGIN_LOGS, params)
				
				if (this.page === 1) {
					this.logList = res.data.list || []
				} else {
					this.logList.push(...(res.data.list || []))
				}
				
				this.hasMore = this.logList.length < res.data.total
				
			} catch (error) {
				console.error('加载登录日志失败:', error)
				uni.showToast({
					title: '加载失败',
					icon: 'none'
				})
			} finally {
				this.loading = false
			}
		},
		
		// 加载操作日志（模拟数据）
		async loadOperationLogs() {
			// TODO: 后端API开发后替换为真实接口
			this.operationLogs = [
				{
					id: 1,
					username: 'admin',
					operationType: '用户管理',
					content: '禁用了用户 zhangsan',
					module: '用户模块',
					ip: '192.168.1.100',
					createTime: '2024-12-05 10:30:25'
				},
				{
					id: 2,
					username: 'admin',
					operationType: '医生审核',
					content: '通过了医生申请 DR20241205001',
					module: '审核模块',
					ip: '192.168.1.100',
					createTime: '2024-12-05 09:15:10'
				},
				{
					id: 3,
					username: 'admin',
					operationType: '系统设置',
					content: '修改了Token有效期为2小时',
					module: '系统设置',
					ip: '192.168.1.100',
					createTime: '2024-12-05 08:45:33'
				}
			]
		},
		
		// 加载错误日志（模拟数据）
		async loadErrorLogs() {
			// TODO: 后端API开发后替换为真实接口
			this.errorLogs = [
				{
					id: 1,
					level: 2,
					levelText: '错误',
					module: '数据库',
					message: '数据库连接超时',
					stack: 'Error: connect ETIMEDOUT\n  at Connection.connect\n  at Database.query',
					createTime: '2024-12-04 23:58:12'
				},
				{
					id: 2,
					level: 1,
					levelText: '警告',
					module: 'API',
					message: 'API请求响应时间过长 (>3s)',
					stack: null,
					createTime: '2024-12-04 18:22:45'
				}
			]
			
			this.errorCount = this.errorLogs.filter(log => log.level >= 2).length
		},
		
		// 加载更多
		loadMore() {
			if (this.hasMore && !this.loading && this.currentTab === 'login') {
				this.page++
				this.loadLogs()
			}
		},
		
		// 格式化时间
		formatTime(timeStr) {
			if (!timeStr) return ''
			
			const now = new Date()
			const time = new Date(timeStr)
			const diff = now - time
			
			// 一分钟内
			if (diff < 60 * 1000) {
				return '刚刚'
			}
			
			// 一小时内
			if (diff < 60 * 60 * 1000) {
				return Math.floor(diff / (60 * 1000)) + '分钟前'
			}
			
			// 今天
			if (now.toDateString() === time.toDateString()) {
				return '今天 ' + timeStr.split(' ')[1].substring(0, 5)
			}
			
			// 昨天
			const yesterday = new Date(now)
			yesterday.setDate(yesterday.getDate() - 1)
			if (yesterday.toDateString() === time.toDateString()) {
				return '昨天 ' + timeStr.split(' ')[1].substring(0, 5)
			}
			
			// 其他
			return timeStr
		}
	}
}
</script>

<style scoped>
.logs-page {
	height: 100vh;
	display: flex;
	flex-direction: column;
	background: #f5f5f5;
}

/* 日志类型切换 */
.log-tabs {
	display: flex;
	background: white;
	padding: 20rpx 30rpx;
	box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.log-tab {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 20rpx 0;
	position: relative;
}

.log-tab.active .tab-icon {
	transform: scale(1.1);
}

.log-tab.active .tab-text {
	color: #ff6b6b;
	font-weight: bold;
}

.log-tab.active::after {
	content: '';
	position: absolute;
	bottom: 0;
	left: 50%;
	transform: translateX(-50%);
	width: 60rpx;
	height: 6rpx;
	background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
	border-radius: 3rpx;
}

.tab-icon {
	font-size: 48rpx;
	margin-bottom: 10rpx;
	transition: transform 0.3s;
}

.tab-text {
	font-size: 24rpx;
	color: #666;
	transition: all 0.3s;
}

.error-badge {
	position: absolute;
	top: 15rpx;
	right: 20%;
	background: #ff4444;
	color: white;
	font-size: 20rpx;
	padding: 4rpx 10rpx;
	border-radius: 20rpx;
	min-width: 30rpx;
	text-align: center;
}

/* 筛选栏 */
.filter-bar {
	background: white;
	padding: 20rpx 0;
	border-top: 1rpx solid #f0f0f0;
}

.status-tabs {
	white-space: nowrap;
	padding: 0 30rpx;
}

.status-tab {
	display: inline-block;
	padding: 10rpx 25rpx;
	margin-right: 20rpx;
	background: #f5f5f5;
	border-radius: 30rpx;
	font-size: 26rpx;
	color: #666;
	transition: all 0.3s;
}

.status-tab.active {
	background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
	color: white;
}

/* 日志列表 */
.log-list {
	flex: 1;
	padding: 20rpx 30rpx;
}

.log-item {
	background: white;
	border-radius: 15rpx;
	padding: 25rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.log-item.error-log {
	border-left: 6rpx solid #ff4444;
}

.log-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20rpx;
}

.user-info {
	display: flex;
	align-items: center;
	gap: 15rpx;
}

.username {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
}

.status-tag {
	font-size: 22rpx;
	padding: 5rpx 15rpx;
	border-radius: 20rpx;
}

.status-0 {
	background: #ffebee;
	color: #f44336;
}

.status-1 {
	background: #e8f5e9;
	color: #4caf50;
}

.operation-type {
	font-size: 24rpx;
	padding: 5rpx 15rpx;
	border-radius: 20rpx;
	background: #e3f2fd;
	color: #2196f3;
}

.log-time {
	font-size: 22rpx;
	color: #999;
}

.log-details {
	margin-top: 15rpx;
}

.detail-row {
	display: flex;
	margin-bottom: 10rpx;
	font-size: 24rpx;
}

.detail-row:last-child {
	margin-bottom: 0;
}

.detail-label {
	color: #999;
	min-width: 120rpx;
}

.detail-value {
	color: #666;
	flex: 1;
}

.log-content {
	margin-bottom: 15rpx;
}

.content-text {
	font-size: 26rpx;
	color: #333;
	line-height: 1.5;
}

/* 错误日志特殊样式 */
.error-level {
	font-size: 22rpx;
	padding: 5rpx 15rpx;
	border-radius: 20rpx;
	font-weight: bold;
}

.level-1 {
	background: #fff3e0;
	color: #ff9800;
}

.level-2 {
	background: #ffebee;
	color: #f44336;
}

.level-3 {
	background: #fce4ec;
	color: #e91e63;
}

.error-module {
	font-size: 24rpx;
	color: #666;
}

.error-message {
	margin-bottom: 15rpx;
	padding: 15rpx;
	background: #fff5f5;
	border-radius: 10rpx;
}

.message-text {
	font-size: 26rpx;
	color: #d32f2f;
	line-height: 1.5;
}

.error-stack {
	padding: 15rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
}

.stack-title {
	display: block;
	font-size: 24rpx;
	color: #999;
	margin-bottom: 10rpx;
}

.stack-text {
	font-size: 22rpx;
	color: #666;
	font-family: monospace;
	line-height: 1.6;
	word-break: break-all;
}

/* 加载状态 */
.loading {
	text-align: center;
	padding: 30rpx 0;
	font-size: 26rpx;
	color: #999;
}

.no-more {
	text-align: center;
	padding: 30rpx 0;
	font-size: 26rpx;
	color: #999;
}

.empty {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 150rpx 0;
}

.empty-icon {
	font-size: 120rpx;
	margin-bottom: 30rpx;
}

.empty-text {
	font-size: 28rpx;
	color: #999;
}
</style>
