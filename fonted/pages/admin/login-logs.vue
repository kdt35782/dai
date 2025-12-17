<template>
	<view class="logs-page">
		<!-- 筛选栏 -->
		<view class="filter-bar">
			<view class="filter-item">
				<text class="filter-label">状态：</text>
				<picker :range="statusOptions" range-key="label" @change="onStatusChange">
					<view class="filter-value">
						{{ statusOptions[statusIndex].label }} ▼
					</view>
				</picker>
			</view>
		</view>
		
		<!-- 日志列表 -->
		<view class="log-list">
			<view class="log-item" v-for="log in logList" :key="log.logId">
				<view class="log-header">
					<view class="user-info">
						<text class="username">{{ log.username }}</text>
						<view class="status-tag" :class="log.status === 1 ? 'status-success' : 'status-fail'">
							{{ log.status === 1 ? '成功' : '失败' }}
						</view>
					</view>
					<text class="log-time">{{ log.loginTime }}</text>
				</view>
				
				<view class="log-details">
					<text class="detail-item">📍 IP: {{ log.loginIp }}</text>
					<text class="detail-item">🌐 位置: {{ log.loginLocation || '未知' }}</text>
					<text class="detail-item">💻 系统: {{ log.os }}</text>
					<text class="detail-item">🔍 浏览器: {{ log.browser }}</text>
				</view>
				
				<text class="log-msg" v-if="log.msg">{{ log.msg }}</text>
			</view>
			
			<!-- 空状态 -->
			<view class="empty-state" v-if="logList.length === 0 && !loading">
				<text class="empty-icon">📭</text>
				<text class="empty-text">暂无登录日志</text>
			</view>
		</view>
		
		<!-- 分页 -->
		<view class="pagination" v-if="total > 0">
			<button 
				class="page-btn" 
				:disabled="currentPage === 1"
				@click="prevPage"
			>上一页</button>
			<text class="page-info">{{ currentPage }} / {{ totalPages }}</text>
			<button 
				class="page-btn" 
				:disabled="currentPage === totalPages"
				@click="nextPage"
			>下一页</button>
		</view>
		
		<!-- 加载状态 -->
		<view class="loading" v-if="loading">
			<text class="loading-text">加载中...</text>
		</view>
	</view>
</template>

<script>
import { get } from '@/utils/request.js'
import { API } from '@/utils/config.js'

export default {
	data() {
		return {
			statusIndex: 0,
			statusOptions: [
				{ label: '全部状态', value: '' },
				{ label: '成功', value: 1 },
				{ label: '失败', value: 0 }
			],
			logList: [],
			currentPage: 1,
			pageSize: 10,
			total: 0,
			loading: false
		}
	},
	
	computed: {
		totalPages() {
			return Math.ceil(this.total / this.pageSize) || 1
		}
	},
	
	onLoad() {
		this.loadLogs()
	},
	
	methods: {
		// 加载日志列表
		async loadLogs() {
			this.loading = true
			try {
				const params = {
					page: this.currentPage,
					pageSize: this.pageSize
				}
				
				if (this.statusOptions[this.statusIndex].value !== '') {
					params.status = this.statusOptions[this.statusIndex].value
				}
				
				const res = await get(API.USER_ADMIN_LOGIN_LOGS, params)
				
				this.logList = res.data.list || []
				this.total = res.data.total || 0
				
			} catch (error) {
				console.error('加载日志失败:', error)
				uni.showToast({
					title: '加载失败',
					icon: 'none'
				})
			} finally {
				this.loading = false
			}
		},
		
		// 状态筛选变化
		onStatusChange(e) {
			this.statusIndex = e.detail.value
			this.currentPage = 1
			this.loadLogs()
		},
		
		// 上一页
		prevPage() {
			if (this.currentPage > 1) {
				this.currentPage--
				this.loadLogs()
			}
		},
		
		// 下一页
		nextPage() {
			if (this.currentPage < this.totalPages) {
				this.currentPage++
				this.loadLogs()
			}
		}
	}
}
</script>

<style scoped>
.logs-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 20rpx;
}

/* 筛选栏 */
.filter-bar {
	display: flex;
	gap: 30rpx;
	background: white;
	padding: 20rpx 30rpx;
	border-radius: 10rpx;
	margin-bottom: 20rpx;
}

.filter-item {
	display: flex;
	align-items: center;
}

.filter-label {
	font-size: 26rpx;
	color: #666;
	margin-right: 10rpx;
}

.filter-value {
	font-size: 26rpx;
	color: #333;
	padding: 10rpx 20rpx;
	background: #f5f5f5;
	border-radius: 8rpx;
}

/* 日志列表 */
.log-list {
	margin-bottom: 20rpx;
}

.log-item {
	background: white;
	border-radius: 15rpx;
	padding: 25rpx;
	margin-bottom: 15rpx;
}

.log-header {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	margin-bottom: 15rpx;
}

.user-info {
	flex: 1;
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
	padding: 5rpx 15rpx;
	border-radius: 8rpx;
	font-size: 22rpx;
	color: white;
}

.status-success {
	background: #4caf50;
}

.status-fail {
	background: #f44336;
}

.log-time {
	font-size: 22rpx;
	color: #999;
}

.log-details {
	display: flex;
	flex-direction: column;
	gap: 8rpx;
	margin-bottom: 10rpx;
}

.detail-item {
	font-size: 24rpx;
	color: #666;
}

.log-msg {
	display: block;
	font-size: 22rpx;
	color: #999;
	padding-top: 10rpx;
	border-top: 1rpx solid #f0f0f0;
}

/* 空状态 */
.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 100rpx 0;
}

.empty-icon {
	font-size: 100rpx;
	margin-bottom: 20rpx;
}

.empty-text {
	font-size: 28rpx;
	color: #999;
}

/* 分页 */
.pagination {
	display: flex;
	justify-content: center;
	align-items: center;
	gap: 30rpx;
	padding: 30rpx 0;
}

.page-btn {
	background: white;
	color: #333;
	border: 1rpx solid #ddd;
	border-radius: 8rpx;
	padding: 15rpx 30rpx;
	font-size: 26rpx;
}

.page-btn:disabled {
	opacity: 0.4;
}

.page-btn::after {
	border: none;
}

.page-info {
	font-size: 26rpx;
	color: #666;
}

/* 加载状态 */
.loading {
	text-align: center;
	padding: 40rpx 0;
}

.loading-text {
	font-size: 26rpx;
	color: #999;
}
</style>
