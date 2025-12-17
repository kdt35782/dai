<template>
	<view class="notification-page">
		<!-- 标签页 -->
		<view class="tabs">
			<view 
				class="tab" 
				:class="{ active: currentTab === 'all' }"
				@click="switchTab('all')"
			>
				全部消息
			</view>
			<view 
				class="tab" 
				:class="{ active: currentTab === 'system' }"
				@click="switchTab('system')"
			>
				系统通知
			</view>
			<view 
				class="tab" 
				:class="{ active: currentTab === 'consultation' }"
				@click="switchTab('consultation')"
			>
				问诊消息
			</view>
		</view>
		
		<!-- 消息列表 -->
		<scroll-view 
			class="message-list" 
			scroll-y
			@scrolltolower="loadMore"
		>
			<view 
				class="message-item" 
				:class="{ unread: !item.isRead }"
				v-for="item in list" 
				:key="item.notificationId"
				@click="handleClick(item)"
			>
				<view class="message-icon">
					<text v-if="item.type === 'system'">📢</text>
					<text v-else-if="item.type === 'consultation'">💬</text>
					<text v-else>📬</text>
				</view>
				
				<view class="message-content">
					<text class="message-title">{{ item.title }}</text>
					<text class="message-text">{{ item.content }}</text>
					<text class="message-time">{{ item.createdAt }}</text>
				</view>
				
				<view class="unread-dot" v-if="!item.isRead"></view>
			</view>
			
			<!-- 加载状态 -->
			<view class="loading" v-if="loading">加载中...</view>
			<view class="no-more" v-if="!hasMore && list.length > 0">没有更多了</view>
			<view class="empty" v-if="!loading && list.length === 0">
				<text class="empty-icon">📭</text>
				<text class="empty-text">暂无消息</text>
			</view>
		</scroll-view>
		
		<!-- 底部操作栏 -->
		<view class="action-bar" v-if="list.length > 0">
			<button class="action-btn" @click="markAllRead">
				全部标为已读
			</button>
		</view>
	</view>
</template>

<script>
import { get, put } from '@/utils/request.js'
import { API } from '@/utils/config.js'

export default {
	data() {
		return {
			currentTab: 'all',
			list: [],
			page: 1,
			pageSize: 20,
			loading: false,
			hasMore: true
		}
	},
	
	onLoad() {
		this.loadList()
	},
	
	methods: {
		// 切换标签
		switchTab(tab) {
			this.currentTab = tab
			this.loadList(true)
		},
		
		// 加载列表
		async loadList(isRefresh = false) {
			if (this.loading) return
			
			if (isRefresh) {
				this.page = 1
				this.list = []
				this.hasMore = true
			}
			
			this.loading = true
			
			try {
				const params = {
					page: this.page,
					pageSize: this.pageSize
				}
				
				if (this.currentTab !== 'all') {
					params.type = this.currentTab
				}
				
				const res = await get(API.NOTIFICATION_LIST, params)
				
				const list = res.data.list || []
				
				if (isRefresh) {
					this.list = list
				} else {
					this.list.push(...list)
				}
				
				this.hasMore = this.list.length < res.data.total
				
			} catch (error) {
				console.error('加载消息列表失败:', error)
			} finally {
				this.loading = false
			}
		},
		
		// 加载更多
		loadMore() {
			if (this.hasMore && !this.loading) {
				this.page++
				this.loadList()
			}
		},
		
		// 点击消息
		async handleClick(item) {
			// 标记为已读
			if (!item.isRead) {
				try {
					await put(API.NOTIFICATION_MARK_READ, {
						notificationIds: [item.notificationId]
					})
					
					item.isRead = true
					
				} catch (error) {
					console.error('标记已读失败:', error)
				}
			}
			
			// 跳转相关页面
			if (item.type === 'consultation' && item.relatedId) {
				uni.navigateTo({
					url: '/pages/consultation/consultation-detail?id=' + item.relatedId
				})
			} else {
				// 显示详情
				uni.showModal({
					title: item.title,
					content: item.content,
					showCancel: false
				})
			}
		},
		
		// 全部标为已读
		async markAllRead() {
			const unreadIds = this.list
				.filter(item => !item.isRead)
				.map(item => item.notificationId)
			
			if (unreadIds.length === 0) {
				uni.showToast({
					title: '没有未读消息',
					icon: 'none'
				})
				return
			}
			
			try {
				await put(API.NOTIFICATION_MARK_READ, {
					notificationIds: unreadIds
				})
				
				this.list.forEach(item => {
					item.isRead = true
				})
				
				uni.showToast({
					title: '已全部标记为已读',
					icon: 'success'
				})
				
			} catch (error) {
				console.error('标记已读失败:', error)
			}
		}
	}
}
</script>

<style scoped>
.notification-page {
	height: 100vh;
	display: flex;
	flex-direction: column;
	background: #f5f5f5;
}

.tabs {
	display: flex;
	background: white;
	padding: 10rpx 30rpx 0;
}

.tab {
	flex: 1;
	text-align: center;
	padding: 25rpx 0;
	font-size: 28rpx;
	color: #666;
	position: relative;
}

.tab.active {
	color: #07c160;
	font-weight: bold;
}

.tab.active::after {
	content: '';
	position: absolute;
	bottom: 0;
	left: 50%;
	transform: translateX(-50%);
	width: 60rpx;
	height: 6rpx;
	background: #07c160;
	border-radius: 3rpx;
}

.message-list {
	flex: 1;
	padding-bottom: 120rpx;
}

.message-item {
	display: flex;
	align-items: flex-start;
	padding: 30rpx;
	background: white;
	margin-bottom: 2rpx;
	position: relative;
}

.message-item.unread {
	background: #f0f9ff;
}

.message-icon {
	width: 80rpx;
	height: 80rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 40rpx;
	margin-right: 20rpx;
	flex-shrink: 0;
}

.message-content {
	flex: 1;
}

.message-title {
	display: block;
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 10rpx;
}

.message-text {
	display: block;
	font-size: 26rpx;
	color: #666;
	line-height: 1.6;
	margin-bottom: 10rpx;
	overflow: hidden;
	text-overflow: ellipsis;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
}

.message-time {
	display: block;
	font-size: 22rpx;
	color: #999;
}

.unread-dot {
	width: 16rpx;
	height: 16rpx;
	background: #f56c6c;
	border-radius: 50%;
	position: absolute;
	top: 35rpx;
	right: 30rpx;
}

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

.action-bar {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: white;
	padding: 20rpx 30rpx;
	box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.action-btn {
	width: 100%;
	height: 80rpx;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
	border: none;
	border-radius: 40rpx;
	font-size: 28rpx;
}

.action-btn::after {
	border: none;
}
</style>
