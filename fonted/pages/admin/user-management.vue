<template>
	<view class="management-page">
		<!-- 顶部搜索栏 -->
		<view class="search-bar">
			<input 
				class="search-input" 
				v-model="searchKeyword" 
				placeholder="搜索用户名或邮箱"
				@confirm="handleSearch"
			/>
			<button class="search-btn" @click="handleSearch">🔍 搜索</button>
		</view>
		
		<!-- 筛选条件 -->
		<view class="filter-bar">
			<view class="filter-item">
				<text class="filter-label">身份：</text>
				<picker :range="identityOptions" range-key="label" @change="onIdentityChange">
					<view class="filter-value">
						{{ identityOptions[identityIndex].label }} ▼
					</view>
				</picker>
			</view>
			<view class="filter-item">
				<text class="filter-label">状态：</text>
				<picker :range="statusOptions" range-key="label" @change="onStatusChange">
					<view class="filter-value">
						{{ statusOptions[statusIndex].label }} ▼
					</view>
				</picker>
			</view>
		</view>
		
		<!-- 用户列表 -->
		<view class="user-list">
			<view class="user-item" v-for="user in userList" :key="user.userId">
				<view class="user-header">
					<view class="user-basic">
						<text class="user-name">{{ user.username }}</text>
						<view class="role-tag" :class="'role-' + user.identify">
							{{ getRoleName(user.identify) }}
						</view>
						<view class="status-tag" :class="user.status === 0 ? 'status-active' : 'status-disabled'">
							{{ user.status === 0 ? '正常' : '禁用' }}
						</view>
					</view>
					<button 
						class="action-btn" 
						:class="user.status === 0 ? 'btn-danger' : 'btn-success'"
						@click="toggleUserStatus(user)"
					>
						{{ user.status === 0 ? '禁用' : '启用' }}
					</button>
				</view>
				
				<view class="user-info">
					<text class="info-item">📧 {{ user.email || '未填写' }}</text>
					<text class="info-item">📱 {{ user.phone || '未填写' }}</text>
				</view>
				
				<view class="user-footer">
					<text class="time-text">注册时间：{{ user.createdAt }}</text>
					<text class="time-text">最后登录：{{ user.lastLoginTime || '从未登录' }}</text>
				</view>
			</view>
			
			<!-- 空状态 -->
			<view class="empty-state" v-if="userList.length === 0 && !loading">
				<text class="empty-icon">📭</text>
				<text class="empty-text">暂无用户数据</text>
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
import { get, put } from '@/utils/request.js'
import { API } from '@/utils/config.js'

export default {
	data() {
		return {
			searchKeyword: '',
			identityIndex: 0,
			statusIndex: 0,
			identityOptions: [
				{ label: '全部身份', value: '' },
				{ label: '普通用户', value: 'user' },
				{ label: '医生', value: 'doctor' },
				{ label: '管理员', value: 'admin' }
			],
			statusOptions: [
				{ label: '全部状态', value: '' },
				{ label: '正常', value: 0 },
				{ label: '禁用', value: 1 }
			],
			userList: [],
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
		this.loadUsers()
	},
	
	methods: {
		// 加载用户列表
		async loadUsers() {
			this.loading = true
			try {
				const params = {
					page: this.currentPage,
					pageSize: this.pageSize
				}
				
				if (this.searchKeyword) {
					params.keyword = this.searchKeyword
				}
				
				if (this.identityOptions[this.identityIndex].value) {
					params.identify = this.identityOptions[this.identityIndex].value
				}
				
				if (this.statusOptions[this.statusIndex].value !== '') {
					params.status = this.statusOptions[this.statusIndex].value
				}
				
				const res = await get(API.USER_ADMIN_USERS, params)
				
				this.userList = res.data.list || []
				this.total = res.data.total || 0
				
			} catch (error) {
				console.error('加载用户列表失败:', error)
				uni.showToast({
					title: '加载失败',
					icon: 'none'
				})
			} finally {
				this.loading = false
			}
		},
		
		// 搜索
		handleSearch() {
			this.currentPage = 1
			this.loadUsers()
		},
		
		// 身份筛选变化
		onIdentityChange(e) {
			this.identityIndex = e.detail.value
			this.currentPage = 1
			this.loadUsers()
		},
		
		// 状态筛选变化
		onStatusChange(e) {
			this.statusIndex = e.detail.value
			this.currentPage = 1
			this.loadUsers()
		},
		
		// 获取角色名称
		getRoleName(identify) {
			const roleMap = {
				'user': '普通用户',
				'doctor': '医生',
				'admin': '管理员'
			}
			return roleMap[identify] || '未知'
		},
		
		// 切换用户状态
		toggleUserStatus(user) {
			const newStatus = user.status === 0 ? 1 : 0
			const actionText = newStatus === 1 ? '禁用' : '启用'
			
			uni.showModal({
				title: '确认操作',
				content: `确定要${actionText}用户 ${user.username} 吗？`,
				success: async (res) => {
					if (res.confirm) {
						await this.updateUserStatus(user.userId, newStatus)
					}
				}
			})
		},
		
		// 更新用户状态
		async updateUserStatus(userId, status) {
			try {
				await put(API.USER_ADMIN_STATUS, {
					userId,
					status
				})
				
				uni.showToast({
					title: '操作成功',
					icon: 'success'
				})
				
				// 刷新列表
				this.loadUsers()
				
			} catch (error) {
				console.error('更新用户状态失败:', error)
				uni.showToast({
					title: '操作失败',
					icon: 'none'
				})
			}
		},
		
		// 上一页
		prevPage() {
			if (this.currentPage > 1) {
				this.currentPage--
				this.loadUsers()
			}
		},
		
		// 下一页
		nextPage() {
			if (this.currentPage < this.totalPages) {
				this.currentPage++
				this.loadUsers()
			}
		}
	}
}
</script>

<style scoped>
.management-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 20rpx;
}

/* 搜索栏 */
.search-bar {
	display: flex;
	gap: 15rpx;
	margin-bottom: 20rpx;
}

.search-input {
	flex: 1;
	background: white;
	border-radius: 10rpx;
	padding: 20rpx 30rpx;
	font-size: 28rpx;
}

.search-btn {
	background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
	color: white;
	border: none;
	border-radius: 10rpx;
	padding: 20rpx 30rpx;
	font-size: 28rpx;
}

.search-btn::after {
	border: none;
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

/* 用户列表 */
.user-list {
	margin-bottom: 20rpx;
}

.user-item {
	background: white;
	border-radius: 15rpx;
	padding: 25rpx;
	margin-bottom: 15rpx;
}

.user-header {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	margin-bottom: 20rpx;
}

.user-basic {
	flex: 1;
	display: flex;
	flex-wrap: wrap;
	gap: 10rpx;
	align-items: center;
}

.user-name {
	font-size: 30rpx;
	font-weight: bold;
	color: #333;
}

.role-tag {
	padding: 5rpx 15rpx;
	border-radius: 8rpx;
	font-size: 22rpx;
	color: white;
}

.role-user {
	background: #667eea;
}

.role-doctor {
	background: #07c160;
}

.role-admin {
	background: #ff6b6b;
}

.status-tag {
	padding: 5rpx 15rpx;
	border-radius: 8rpx;
	font-size: 22rpx;
	color: white;
}

.status-active {
	background: #4caf50;
}

.status-disabled {
	background: #999;
}

.action-btn {
	padding: 10rpx 25rpx;
	border-radius: 8rpx;
	font-size: 24rpx;
	color: white;
	border: none;
}

.action-btn::after {
	border: none;
}

.btn-danger {
	background: #ff6b6b;
}

.btn-success {
	background: #4caf50;
}

.user-info {
	display: flex;
	flex-direction: column;
	gap: 10rpx;
	margin-bottom: 15rpx;
}

.info-item {
	font-size: 24rpx;
	color: #666;
}

.user-footer {
	display: flex;
	justify-content: space-between;
	padding-top: 15rpx;
	border-top: 1rpx solid #f0f0f0;
}

.time-text {
	font-size: 22rpx;
	color: #999;
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
