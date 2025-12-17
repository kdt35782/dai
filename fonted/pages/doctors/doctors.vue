<template>
	<view class="doctors-page">
		<!-- 搜索栏 -->
		<view class="search-bar">
			<view class="search-box">
				<text class="search-icon">🔍</text>
				<input 
					class="search-input" 
					v-model="keyword" 
					placeholder="搜索医生姓名或擅长领域"
					@confirm="handleSearch"
				/>
			</view>
		</view>
		
		<!-- 科室分类 -->
		<scroll-view class="dept-tabs" scroll-x>
			<view 
				class="dept-item" 
				:class="{ active: currentDept === item }"
				v-for="item in deptList" 
				:key="item"
				@click="selectDept(item)"
			>
				{{ item }}
			</view>
		</scroll-view>
		
		<!-- 医生列表 -->
		<scroll-view 
			class="doctor-list" 
			scroll-y
			@scrolltolower="loadMore"
		>
			<view 
				class="doctor-item" 
				v-for="doctor in doctorList" 
				:key="doctor.userId"
				@click="goDoctorDetail(doctor.userId)"
			>
				<image class="avatar" :src="doctor.avatar || '/static/default-avatar.png'" mode="aspectFill"></image>
				
				<view class="info">
					<view class="name-row">
						<text class="name">{{ doctor.realName }}</text>
						<text class="title">{{ doctor.doctorTitle }}</text>
					</view>
					
					<text class="dept">{{ doctor.doctorDept }}</text>
					
					<text class="specialty">擅长：{{ doctor.specialty || '暂无' }}</text>
					
					<view class="stats">
						<text class="stat-item">💬 {{ doctor.consultationCount || 0 }}次问诊</text>
						<text class="stat-item">⭐ {{ doctor.rating || '5.0' }}分</text>
					</view>
				</view>
				
				<view class="action">
					<view class="consult-btn" @click.stop="createConsultation(doctor)">
						问诊
					</view>
				</view>
			</view>
			
			<!-- 加载状态 -->
			<view class="loading" v-if="loading">加载中...</view>
			<view class="no-more" v-if="!hasMore && doctorList.length > 0">没有更多了</view>
			<view class="empty" v-if="!loading && doctorList.length === 0">
				<text class="empty-icon">👨‍⚕️</text>
				<text class="empty-text">暂无医生</text>
			</view>
		</scroll-view>
	</view>
</template>

<script>
import { get } from '@/utils/request.js'
import { API, STORAGE_KEYS } from '@/utils/config.js'

export default {
	data() {
		return {
			keyword: '',
			currentDept: '全部',
			deptList: ['全部', '内科', '外科', '儿科', '妇产科', '骨科', '皮肤科', '眼科', '耳鼻喉科', '口腔科', '中医科'],
			doctorList: [],
			page: 1,
			pageSize: 10,
			loading: false,
			hasMore: true
		}
	},
	
	onLoad() {
		this.loadDoctors()
	},
	
	methods: {
		// 加载医生列表
		async loadDoctors(isRefresh = false) {
			if (this.loading) return
			
			if (isRefresh) {
				this.page = 1
				this.doctorList = []
				this.hasMore = true
			}
			
			this.loading = true
			
			try {
				const params = {
					page: this.page,
					pageSize: this.pageSize
				}
				
				if (this.currentDept !== '全部') {
					params.dept = this.currentDept
				}
				
				if (this.keyword) {
					params.keyword = this.keyword
				}
				
				const res = await get(API.USER_DOCTORS, params, { noAuth: true })
				
				const list = res.data.list || []
				
				if (isRefresh) {
					this.doctorList = list
				} else {
					this.doctorList.push(...list)
				}
				
				this.hasMore = this.doctorList.length < res.data.total
				
			} catch (error) {
				console.error('加载医生列表失败:', error)
			} finally {
				this.loading = false
			}
		},
		
		// 搜索
		handleSearch() {
			this.loadDoctors(true)
		},
		
		// 选择科室
		selectDept(dept) {
			this.currentDept = dept
			this.loadDoctors(true)
		},
		
		// 加载更多
		loadMore() {
			if (this.hasMore && !this.loading) {
				this.page++
				this.loadDoctors()
			}
		},
		
		// 跳转医生详情
		goDoctorDetail(userId) {
			uni.navigateTo({
				url: '/pages/doctor-detail/doctor-detail?userId=' + userId
			})
		},
		
		// 发起问诊
		createConsultation(doctor) {
			const token = uni.getStorageSync(STORAGE_KEYS.TOKEN)
			
			if (!token) {
				uni.showToast({
					title: '请先登录',
					icon: 'none'
				})
				
				setTimeout(() => {
					uni.navigateTo({
						url: '/pages/login/login'
					})
				}, 1500)
				return
			}
			
			uni.navigateTo({
				url: `/pages/consultation/create-consultation?doctorId=${doctor.userId}&doctorName=${doctor.realName}&doctorDept=${doctor.doctorDept || ''}`
			})
		}
	}
}
</script>

<style scoped>
.doctors-page {
	height: 100vh;
	display: flex;
	flex-direction: column;
	background: #f5f5f5;
}

.search-bar {
	background: white;
	padding: 20rpx 30rpx;
}

.search-box {
	display: flex;
	align-items: center;
	background: #f5f5f5;
	border-radius: 50rpx;
	padding: 15rpx 25rpx;
}

.search-icon {
	font-size: 32rpx;
	margin-right: 15rpx;
}

.search-input {
	flex: 1;
	font-size: 28rpx;
}

.dept-tabs {
	background: white;
	white-space: nowrap;
	padding: 20rpx 30rpx;
	border-top: 1px solid #f0f0f0;
}

.dept-item {
	display: inline-block;
	padding: 10rpx 25rpx;
	margin-right: 20rpx;
	background: #f5f5f5;
	border-radius: 30rpx;
	font-size: 26rpx;
	color: #666;
}

.dept-item.active {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
}

.doctor-list {
	flex: 1;
	padding: 20rpx 30rpx;
}

.doctor-item {
	display: flex;
	background: white;
	border-radius: 20rpx;
	padding: 25rpx;
	margin-bottom: 20rpx;
}

.avatar {
	width: 140rpx;
	height: 140rpx;
	border-radius: 15rpx;
	margin-right: 20rpx;
	flex-shrink: 0;
}

.info {
	flex: 1;
	display: flex;
	flex-direction: column;
	justify-content: space-between;
}

.name-row {
	display: flex;
	align-items: center;
	margin-bottom: 8rpx;
}

.name {
	font-size: 30rpx;
	font-weight: bold;
	color: #333;
	margin-right: 15rpx;
}

.title {
	font-size: 22rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	color: white;
	padding: 4rpx 12rpx;
	border-radius: 10rpx;
}

.dept {
	font-size: 24rpx;
	color: #07c160;
	margin-bottom: 8rpx;
}

.specialty {
	font-size: 24rpx;
	color: #666;
	margin-bottom: 8rpx;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.stats {
	display: flex;
}

.stat-item {
	font-size: 22rpx;
	color: #999;
	margin-right: 20rpx;
}

.action {
	display: flex;
	align-items: center;
}

.consult-btn {
	width: 100rpx;
	height: 60rpx;
	line-height: 60rpx;
	text-align: center;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
	border-radius: 30rpx;
	font-size: 26rpx;
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
