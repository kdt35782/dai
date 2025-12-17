<template>
	<view class="admin-page">
		<!-- 标题栏 -->
		<view class="header">
			<text class="title">医生申请审核</text>
		</view>
		
		<!-- 筛选栏 -->
		<view class="filter-bar">
			<view class="filter-item">
				<picker mode="selector" :range="statusList" :range-key="'label'" @change="onStatusChange">
					<view class="picker">
						<text>{{ currentStatus.label }}</text>
						<text class="arrow">▼</text>
					</view>
				</picker>
			</view>
		</view>
		
		<!-- 申请列表 -->
		<view class="application-list">
			<view v-if="loading" class="loading">
				<text>加载中...</text>
			</view>
			
			<view v-else-if="list.length === 0" class="empty">
				<text class="empty-icon">📋</text>
				<text class="empty-text">暂无申请记录</text>
			</view>
			
			<view v-else>
				<view 
					class="application-item" 
					v-for="item in list" 
					:key="item.applicationId"
					@click="showDetail(item)"
				>
					<!-- 用户信息 -->
					<view class="item-header">
						<view class="user-info">
							<text class="username">{{ item.username }}</text>
							<text class="real-name">{{ item.realName }}</text>
						</view>
						<view class="status-badge" :class="'status-' + item.status">
							{{ item.statusText }}
						</view>
					</view>
					
					<!-- 医生信息 -->
					<view class="item-body">
						<view class="info-row">
							<text class="label">职称：</text>
							<text class="value">{{ item.doctorTitle }}</text>
						</view>
						<view class="info-row">
							<text class="label">科室：</text>
							<text class="value">{{ item.doctorDept }}</text>
						</view>
						<view class="info-row">
							<text class="label">手机：</text>
							<text class="value">{{ item.phone }}</text>
						</view>
						<view class="info-row">
							<text class="label">申请时间：</text>
							<text class="value">{{ item.createdAt }}</text>
						</view>
					</view>
					
					<!-- 操作按钮 -->
					<view class="item-footer" v-if="item.status === 0">
						<button class="btn btn-reject" @click.stop="handleReject(item)">拒绝</button>
						<button class="btn btn-approve" @click.stop="handleApprove(item)">通过</button>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 分页 -->
		<view class="pagination" v-if="total > pageSize">
			<button 
				class="page-btn" 
				:disabled="page === 1"
				@click="prevPage"
			>上一页</button>
			<text class="page-info">{{ page }} / {{ totalPages }}</text>
			<button 
				class="page-btn" 
				:disabled="page >= totalPages"
				@click="nextPage"
			>下一页</button>
		</view>
		
		<!-- 详情弹窗 -->
		<view class="modal" v-if="showDetailModal" @click="closeDetail">
			<view class="modal-content" @click.stop>
				<view class="modal-header">
					<text class="modal-title">申请详情</text>
					<text class="modal-close" @click="closeDetail">×</text>
				</view>
				
				<scroll-view class="modal-body" scroll-y>
					<view class="detail-section">
						<text class="section-title">基本信息</text>
						<view class="detail-row">
							<text class="detail-label">用户名：</text>
							<text class="detail-value">{{ currentItem.username }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">真实姓名：</text>
							<text class="detail-value">{{ currentItem.realName }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">身份证号：</text>
							<text class="detail-value">{{ currentItem.idCard }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">手机号：</text>
							<text class="detail-value">{{ currentItem.phone }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">邮箱：</text>
							<text class="detail-value">{{ currentItem.email }}</text>
						</view>
					</view>
					
					<view class="detail-section">
						<text class="section-title">执业信息</text>
						<view class="detail-row">
							<text class="detail-label">职称：</text>
							<text class="detail-value">{{ currentItem.doctorTitle }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">科室：</text>
							<text class="detail-value">{{ currentItem.doctorDept }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">资格证号：</text>
							<text class="detail-value">{{ currentItem.certNumber }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">擅长领域：</text>
							<text class="detail-value">{{ currentItem.specialty }}</text>
						</view>
						<view class="detail-row">
							<text class="detail-label">个人介绍：</text>
							<text class="detail-value">{{ currentItem.introduction }}</text>
						</view>
					</view>
					
					<view class="detail-section">
						<text class="section-title">资格证书</text>
						<image 
							class="cert-image" 
							:src="currentItem.doctorCert" 
							mode="aspectFit"
							@click="previewImage(currentItem.doctorCert)"
						></image>
					</view>
				</scroll-view>
				
				<view class="modal-footer" v-if="currentItem.status === 0">
					<button class="modal-btn btn-reject" @click="handleReject(currentItem)">拒绝</button>
					<button class="modal-btn btn-approve" @click="handleApprove(currentItem)">通过</button>
				</view>
			</view>
		</view>
		
		<!-- 拒绝理由弹窗 -->
		<view class="modal" v-if="showRejectModal" @click="closeReject">
			<view class="modal-content small" @click.stop>
				<view class="modal-header">
					<text class="modal-title">拒绝申请</text>
					<text class="modal-close" @click="closeReject">×</text>
				</view>
				
				<view class="modal-body">
					<textarea 
						class="reject-input" 
						v-model="rejectReason"
						placeholder="请输入拒绝原因"
						maxlength="200"
					></textarea>
				</view>
				
				<view class="modal-footer">
					<button class="modal-btn btn-cancel" @click="closeReject">取消</button>
					<button class="modal-btn btn-confirm" @click="confirmReject">确认拒绝</button>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
import { get, put } from '@/utils/request.js'
import { API } from '@/utils/config.js'

export default {
	data() {
		return {
			list: [],
			total: 0,
			page: 1,
			pageSize: 10,
			loading: false,
			statusList: [
				{ value: null, label: '全部申请' },
				{ value: 0, label: '待审核' },
				{ value: 1, label: '已通过' },
				{ value: 2, label: '已拒绝' }
			],
			currentStatus: { value: null, label: '全部申请' },
			showDetailModal: false,
			showRejectModal: false,
			currentItem: {},
			rejectReason: ''
		}
	},
	
	computed: {
		totalPages() {
			return Math.ceil(this.total / this.pageSize)
		}
	},
	
	onLoad() {
		this.loadList()
	},
	
	methods: {
		// 加载列表
		async loadList() {
			this.loading = true
			try {
				let url = `${API.USER_ADMIN_APPLICATIONS}?page=${this.page}&pageSize=${this.pageSize}`
				if (this.currentStatus.value !== null) {
					url += `&status=${this.currentStatus.value}`
				}
				
				const res = await get(url)
				this.list = res.data.list || []
				this.total = res.data.total || 0
			} catch (error) {
				console.error('加载申请列表失败:', error)
				uni.showToast({
					title: '加载失败',
					icon: 'none'
				})
			} finally {
				this.loading = false
			}
		},
		
		// 状态筛选
		onStatusChange(e) {
			this.currentStatus = this.statusList[e.detail.value]
			this.page = 1
			this.loadList()
		},
		
		// 上一页
		prevPage() {
			if (this.page > 1) {
				this.page--
				this.loadList()
			}
		},
		
		// 下一页
		nextPage() {
			if (this.page < this.totalPages) {
				this.page++
				this.loadList()
			}
		},
		
		// 显示详情
		showDetail(item) {
			this.currentItem = item
			this.showDetailModal = true
		},
		
		// 关闭详情
		closeDetail() {
			this.showDetailModal = false
			this.currentItem = {}
		},
		
		// 预览图片
		previewImage(url) {
			uni.previewImage({
				urls: [url],
				current: url
			})
		},
		
		// 通过申请
		handleApprove(item) {
			uni.showModal({
				title: '确认操作',
				content: `确定通过 ${item.realName} 的医生申请吗？`,
				success: async (res) => {
					if (res.confirm) {
						await this.reviewApplication(item.applicationId, 1, '')
					}
				}
			})
		},
		
		// 拒绝申请
		handleReject(item) {
			this.currentItem = item
			this.rejectReason = ''
			this.showRejectModal = true
			this.showDetailModal = false
		},
		
		// 关闭拒绝弹窗
		closeReject() {
			this.showRejectModal = false
			this.rejectReason = ''
		},
		
		// 确认拒绝
		async confirmReject() {
			if (!this.rejectReason.trim()) {
				uni.showToast({
					title: '请输入拒绝原因',
					icon: 'none'
				})
				return
			}
			
			await this.reviewApplication(this.currentItem.applicationId, 2, this.rejectReason)
			this.closeReject()
		},
		
		// 审核申请
		async reviewApplication(applicationId, status, rejectReason) {
			try {
				await put(API.USER_ADMIN_REVIEW_DOCTOR, {
					applicationId,
					status,
					rejectReason
				})
				
				uni.showToast({
					title: status === 1 ? '已通过' : '已拒绝',
					icon: 'success'
				})
				
				this.closeDetail()
				this.loadList()
			} catch (error) {
				console.error('审核失败:', error)
				uni.showToast({
					title: '审核失败',
					icon: 'none'
				})
			}
		}
	}
}
</script>

<style scoped>
.admin-page {
	min-height: 100vh;
	background: #f5f5f5;
}

.header {
	background: white;
	padding: 30rpx;
	border-bottom: 1rpx solid #eee;
}

.title {
	font-size: 36rpx;
	font-weight: bold;
	color: #333;
}

.filter-bar {
	background: white;
	padding: 20rpx 30rpx;
	margin-bottom: 20rpx;
}

.picker {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 20rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
}

.arrow {
	color: #999;
	font-size: 24rpx;
}

.application-list {
	padding: 0 30rpx;
}

.loading, .empty {
	text-align: center;
	padding: 100rpx 0;
	color: #999;
}

.empty-icon {
	display: block;
	font-size: 100rpx;
	margin-bottom: 20rpx;
}

.application-item {
	background: white;
	border-radius: 20rpx;
	padding: 30rpx;
	margin-bottom: 20rpx;
}

.item-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20rpx;
	padding-bottom: 20rpx;
	border-bottom: 1rpx solid #eee;
}

.user-info {
	flex: 1;
}

.username {
	font-size: 32rpx;
	font-weight: bold;
	color: #333;
	margin-right: 15rpx;
}

.real-name {
	font-size: 26rpx;
	color: #666;
}

.status-badge {
	padding: 8rpx 20rpx;
	border-radius: 20rpx;
	font-size: 24rpx;
}

.status-0 {
	background: #fff3e0;
	color: #ff9800;
}

.status-1 {
	background: #e8f5e9;
	color: #4caf50;
}

.status-2 {
	background: #ffebee;
	color: #f44336;
}

.item-body {
	margin-bottom: 20rpx;
}

.info-row {
	display: flex;
	padding: 10rpx 0;
	font-size: 28rpx;
}

.label {
	color: #999;
	width: 150rpx;
}

.value {
	flex: 1;
	color: #333;
}

.item-footer {
	display: flex;
	gap: 20rpx;
	padding-top: 20rpx;
	border-top: 1rpx solid #eee;
}

.btn {
	flex: 1;
	height: 70rpx;
	line-height: 70rpx;
	border-radius: 10rpx;
	font-size: 28rpx;
	border: none;
}

.btn-reject {
	background: #fff;
	color: #f44336;
	border: 1rpx solid #f44336;
}

.btn-approve {
	background: #4caf50;
	color: white;
}

.pagination {
	display: flex;
	justify-content: center;
	align-items: center;
	padding: 40rpx 0;
	gap: 30rpx;
}

.page-btn {
	padding: 15rpx 40rpx;
	background: white;
	border-radius: 10rpx;
	font-size: 28rpx;
	border: 1rpx solid #ddd;
}

.page-btn:disabled {
	opacity: 0.5;
}

.page-info {
	font-size: 28rpx;
	color: #666;
}

/* 弹窗样式 */
.modal {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 1000;
}

.modal-content {
	width: 90%;
	max-height: 80%;
	background: white;
	border-radius: 20rpx;
	display: flex;
	flex-direction: column;
}

.modal-content.small {
	max-height: 400rpx;
}

.modal-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 30rpx;
	border-bottom: 1rpx solid #eee;
}

.modal-title {
	font-size: 32rpx;
	font-weight: bold;
}

.modal-close {
	font-size: 50rpx;
	color: #999;
	line-height: 1;
}

.modal-body {
	flex: 1;
	padding: 30rpx;
	overflow-y: auto;
}

.detail-section {
	margin-bottom: 30rpx;
}

.section-title {
	display: block;
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 20rpx;
	padding-bottom: 10rpx;
	border-bottom: 2rpx solid #4caf50;
}

.detail-row {
	display: flex;
	padding: 15rpx 0;
	font-size: 26rpx;
}

.detail-label {
	color: #999;
	width: 180rpx;
	flex-shrink: 0;
}

.detail-value {
	flex: 1;
	color: #333;
	word-break: break-all;
}

.cert-image {
	width: 100%;
	height: 400rpx;
	border-radius: 10rpx;
}

.reject-input {
	width: 100%;
	min-height: 200rpx;
	padding: 20rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	font-size: 28rpx;
	box-sizing: border-box;
}

.modal-footer {
	display: flex;
	gap: 20rpx;
	padding: 30rpx;
	border-top: 1rpx solid #eee;
}

.modal-btn {
	flex: 1;
	height: 80rpx;
	line-height: 80rpx;
	border-radius: 10rpx;
	font-size: 30rpx;
	border: none;
}

.btn-cancel {
	background: #f5f5f5;
	color: #666;
}

.btn-confirm {
	background: #f44336;
	color: white;
}
</style>
