<template>
	<view class="list-page">
		<!-- 筛选栏 -->
		<view class="filter-bar">
			<picker mode="date" :value="startDate" @change="onStartDateChange">
				<view class="date-picker">
					<text class="date-label">开始日期</text>
					<text class="date-value">{{ startDate || '请选择' }}</text>
				</view>
			</picker>
			
			<text class="date-separator">-</text>
			
			<picker mode="date" :value="endDate" @change="onEndDateChange">
				<view class="date-picker">
					<text class="date-label">结束日期</text>
					<text class="date-value">{{ endDate || '请选择' }}</text>
				</view>
			</picker>
		</view>
		
		<!-- 病历列表 -->
		<scroll-view 
			class="record-list" 
			scroll-y
			@scrolltolower="loadMore"
		>
			<view 
				class="record-item" 
				v-for="item in list" 
				:key="item.recordId"
				@click="goDetail(item.recordId)"
			>
				<view class="record-header">
					<view class="doctor-info">
						<!-- 如果有医生名，显示医生信息（患者端） -->
						<template v-if="item.doctorName">
							<text class="doctor-name">{{ item.doctorName }} 医生</text>
							<text class="dept">{{ item.doctorDept }}</text>
						</template>
						<!-- 如果有患者名，显示患者信息（医生端） -->
						<template v-else-if="item.patientName">
							<text class="doctor-name">🧑‍⚕️ {{ item.patientName }}</text>
							<text class="dept">患者</text>
						</template>
					</view>
					<text class="secure-tag">🔐 加密</text>
				</view>
				
				<view class="record-content">
					<view class="info-row">
						<text class="info-label">诊断：</text>
						<text class="info-value">{{ item.diagnosis || '暂无' }}</text>
					</view>
					
					<view class="info-row">
						<text class="info-label">主诉：</text>
						<text class="info-value">{{ item.chiefComplaint }}</text>
					</view>
				</view>
				
				<view class="record-footer">
					<text class="time">{{ item.createdAt }}</text>
					<text class="arrow">></text>
				</view>
			</view>
			
			<!-- 加载状态 -->
			<view class="loading" v-if="loading">加载中...</view>
			<view class="no-more" v-if="!hasMore && list.length > 0">没有更多了</view>
			<view class="empty" v-if="!loading && list.length === 0">
				<text class="empty-icon">📋</text>
				<text class="empty-text">暂无病历记录</text>
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
			startDate: '',
			endDate: '',
			list: [],
			page: 1,
			pageSize: 10,
			loading: false,
			hasMore: true
		}
	},
	
	onLoad() {
		this.loadList()
	},
	
	methods: {
		// 开始日期选择
		onStartDateChange(e) {
			this.startDate = e.detail.value
			this.loadList(true)
		},
		
		// 结束日期选择
		onEndDateChange(e) {
			this.endDate = e.detail.value
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
				
				if (this.startDate) {
					params.startDate = this.startDate
				}
				
				if (this.endDate) {
					params.endDate = this.endDate
				}
				
				const res = await get(API.RECORD_LIST, params)
				
				const list = res.data.list || []
				
				if (isRefresh) {
					this.list = list
				} else {
					this.list.push(...list)
				}
				
				this.hasMore = this.list.length < res.data.total
				
			} catch (error) {
				console.error('加载病历列表失败:', error)
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
		
		// 跳转详情
		goDetail(id) {
			uni.navigateTo({
				url: '/pages/medical-record/record-detail?id=' + id
			})
		}
	}
}
</script>

<style scoped>
.list-page {
	height: 100vh;
	display: flex;
	flex-direction: column;
	background: #f5f5f5;
}

.filter-bar {
	background: white;
	padding: 25rpx 30rpx;
	display: flex;
	align-items: center;
	border-bottom: 1px solid #f0f0f0;
}

.date-picker {
	flex: 1;
	display: flex;
	flex-direction: column;
}

.date-label {
	font-size: 22rpx;
	color: #999;
	margin-bottom: 8rpx;
}

.date-value {
	font-size: 26rpx;
	color: #333;
}

.date-separator {
	margin: 0 20rpx;
	font-size: 26rpx;
	color: #999;
}

.record-list {
	flex: 1;
	padding: 20rpx 30rpx;
}

.record-item {
	background: white;
	border-radius: 20rpx;
	padding: 25rpx;
	margin-bottom: 20rpx;
}

.record-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20rpx;
	padding-bottom: 20rpx;
	border-bottom: 1px solid #f0f0f0;
}

.doctor-info {
	display: flex;
	flex-direction: column;
}

.doctor-name {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 5rpx;
}

.dept {
	font-size: 22rpx;
	color: #07c160;
}

.secure-tag {
	font-size: 22rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
	padding: 5rpx 12rpx;
	border-radius: 10rpx;
}

.record-content {
	margin-bottom: 15rpx;
}

.info-row {
	display: flex;
	margin-bottom: 10rpx;
}

.info-row:last-child {
	margin-bottom: 0;
}

.info-label {
	font-size: 26rpx;
	color: #666;
	width: 100rpx;
	flex-shrink: 0;
}

.info-value {
	flex: 1;
	font-size: 26rpx;
	color: #333;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.record-footer {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.time {
	font-size: 22rpx;
	color: #999;
}

.arrow {
	font-size: 28rpx;
	color: #999;
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
</style>
