<template>
	<view class="statistics-page">
		<!-- 时间筛选 -->
		<view class="time-filter">
			<view 
				class="time-btn" 
				v-for="(item, index) in timeOptions" 
				:key="index"
				:class="{ active: timeIndex === index }"
				@click="changeTimeRange(index)"
			>
				{{ item.label }}
			</view>
		</view>
		
		<!-- 核心数据卡片 -->
		<view class="stats-overview">
			<view class="stat-card" @click="viewDetail('users')">
				<view class="stat-icon user-icon">👥</view>
				<view class="stat-info">
					<text class="stat-number">{{ statistics.totalUsers }}</text>
					<text class="stat-label">总用户数</text>
					<view class="stat-trend">
						<text class="trend-icon">↗</text>
						<text class="trend-text">+{{ statistics.newUsersToday }} 今日新增</text>
					</view>
				</view>
			</view>
			
			<view class="stat-card" @click="viewDetail('doctors')">
				<view class="stat-icon doctor-icon">⚕️</view>
				<view class="stat-info">
					<text class="stat-number">{{ statistics.totalDoctors }}</text>
					<text class="stat-label">在职医生</text>
					<view class="stat-trend">
						<text class="trend-icon">↗</text>
						<text class="trend-text">+{{ statistics.newDoctorsThisMonth }} 本月新增</text>
					</view>
				</view>
			</view>
			
			<view class="stat-card" @click="viewDetail('consultations')">
				<view class="stat-icon consult-icon">💬</view>
				<view class="stat-info">
					<text class="stat-number">{{ statistics.totalConsultations }}</text>
					<text class="stat-label">总问诊数</text>
					<view class="stat-trend">
						<text class="trend-icon">↗</text>
						<text class="trend-text">+{{ statistics.consultationsToday }} 今日</text>
					</view>
				</view>
			</view>
			
			<view class="stat-card" @click="viewDetail('records')">
				<view class="stat-icon record-icon">📋</view>
				<view class="stat-info">
					<text class="stat-number">{{ statistics.totalRecords }}</text>
					<text class="stat-label">电子病历</text>
					<view class="stat-trend">
						<text class="trend-icon">↗</text>
						<text class="trend-text">+{{ statistics.recordsToday }} 今日</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 图表区域（简化版） -->
		<view class="chart-section">
			<view class="section-title">📈 问诊趋势</view>
			<view class="chart-card">
				<view class="bar-chart">
					<view 
						class="bar-item" 
						v-for="(item, index) in chartData" 
						:key="index"
					>
						<view class="bar-wrapper">
							<view 
								class="bar" 
								:style="{ height: (item.value / maxValue * 100) + '%' }"
							></view>
						</view>
						<text class="bar-value">{{ item.value }}</text>
						<text class="bar-label">{{ item.label }}</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 部门统计 -->
		<view class="dept-section">
			<view class="section-title">🏥 科室统计</view>
			<view class="dept-card">
				<view class="dept-item" v-for="dept in deptStats" :key="dept.name">
					<view class="dept-header">
						<text class="dept-name">{{ dept.name }}</text>
						<text class="dept-count">{{ dept.doctorCount }} 位医生</text>
					</view>
					<view class="dept-bar">
						<view 
							class="dept-progress" 
							:style="{ width: (dept.consultationCount / maxConsultation * 100) + '%' }"
						></view>
					</view>
					<text class="dept-consult">问诊：{{ dept.consultationCount }} 次</text>
				</view>
			</view>
		</view>
		
		<!-- 热门医生 -->
		<view class="doctor-section">
			<view class="section-title">⭐ 热门医生 TOP 5</view>
			<view class="doctor-card">
				<view class="doctor-rank-item" v-for="(doctor, index) in topDoctors" :key="doctor.userId">
					<view class="rank-badge" :class="'rank-' + (index + 1)">
						{{ index + 1 }}
					</view>
					<view class="doctor-info">
						<text class="doctor-name">{{ doctor.realName }}</text>
						<text class="doctor-dept">{{ doctor.doctorDept }} · {{ doctor.doctorTitle }}</text>
					</view>
					<view class="doctor-stats">
						<text class="consult-count">{{ doctor.consultationCount }} 次</text>
						<text class="rating">⭐ {{ doctor.rating }}</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 系统健康 -->
		<view class="health-section">
			<view class="section-title">💊 系统健康度</view>
			<view class="health-card">
				<view class="health-item">
					<text class="health-label">数据库连接</text>
					<view class="health-status health-good">正常</view>
				</view>
				<view class="health-item">
					<text class="health-label">API 响应时间</text>
					<view class="health-status health-good">{{ systemHealth.apiResponseTime }}ms</view>
				</view>
				<view class="health-item">
					<text class="health-label">加密服务</text>
					<view class="health-status health-good">运行中</view>
				</view>
				<view class="health-item">
					<text class="health-label">存储空间</text>
					<view class="health-status" :class="systemHealth.storageUsage > 80 ? 'health-warning' : 'health-good'">
						{{ systemHealth.storageUsage }}%
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
import { get } from '@/utils/request.js'
import { API } from '@/utils/config.js'

export default {
	data() {
		return {
			timeIndex: 0,
			timeOptions: [
				{ label: '今日', value: 'today' },
				{ label: '本周', value: 'week' },
				{ label: '本月', value: 'month' },
				{ label: '本年', value: 'year' }
			],
			statistics: {
				totalUsers: 0,
				newUsersToday: 0,
				totalDoctors: 0,
				newDoctorsThisMonth: 0,
				totalConsultations: 0,
				consultationsToday: 0,
				totalRecords: 0,
				recordsToday: 0
			},
			chartData: [
				{ label: '周一', value: 45 },
				{ label: '周二', value: 52 },
				{ label: '周三', value: 38 },
				{ label: '周四', value: 67 },
				{ label: '周五', value: 58 },
				{ label: '周六', value: 42 },
				{ label: '周日', value: 35 }
			],
			deptStats: [
				{ name: '内科', doctorCount: 12, consultationCount: 328 },
				{ name: '外科', doctorCount: 8, consultationCount: 215 },
				{ name: '儿科', doctorCount: 6, consultationCount: 156 },
				{ name: '妇产科', doctorCount: 5, consultationCount: 189 },
				{ name: '骨科', doctorCount: 4, consultationCount: 98 }
			],
			topDoctors: [
				{ userId: 2001, realName: '李医生', doctorDept: '内科', doctorTitle: '主任医师', consultationCount: 328, rating: 4.9 },
				{ userId: 2002, realName: '王医生', doctorDept: '外科', doctorTitle: '副主任医师', consultationCount: 215, rating: 4.8 },
				{ userId: 2003, realName: '刘医生', doctorDept: '儿科', doctorTitle: '主治医师', consultationCount: 156, rating: 4.7 },
				{ userId: 2004, realName: '陈医生', doctorDept: '妇产科', doctorTitle: '主任医师', consultationCount: 189, rating: 4.9 },
				{ userId: 2005, realName: '赵医生', doctorDept: '骨科', doctorTitle: '主治医师', consultationCount: 98, rating: 4.6 }
			],
			systemHealth: {
				apiResponseTime: 45,
				storageUsage: 35
			}
		}
	},
	
	computed: {
		maxValue() {
			return Math.max(...this.chartData.map(item => item.value))
		},
		maxConsultation() {
			return Math.max(...this.deptStats.map(item => item.consultationCount))
		}
	},
	
	onLoad() {
		this.loadStatistics()
	},
	
	methods: {
		// 加载统计数据
		async loadStatistics() {
			try {
				// 加载用户统计
				const userRes = await get(API.USER_ADMIN_USERS, {
					page: 1,
					pageSize: 1
				})
				this.statistics.totalUsers = userRes.data.total || 0
				
				// 加载医生统计
				const doctorRes = await get(API.USER_DOCTORS, {
					page: 1,
					pageSize: 1
				}, { noAuth: true })
				this.statistics.totalDoctors = doctorRes.data.total || 0
				
				// 加载问诊统计
				const consultationRes = await get(API.CONSULTATION_LIST, {
					page: 1,
					pageSize: 1
				})
				this.statistics.totalConsultations = consultationRes.data.total || 0
				
				// 加载病历统计
				const recordRes = await get(API.RECORD_LIST, {
					page: 1,
					pageSize: 1
				})
				this.statistics.totalRecords = recordRes.data.total || 0
				
				// 模拟新增数据
				this.statistics.newUsersToday = Math.floor(Math.random() * 10) + 1
				this.statistics.newDoctorsThisMonth = Math.floor(Math.random() * 5) + 1
				this.statistics.consultationsToday = Math.floor(Math.random() * 20) + 5
				this.statistics.recordsToday = Math.floor(Math.random() * 15) + 3
				
			} catch (error) {
				console.error('加载统计数据失败:', error)
			}
		},
		
		// 改变时间范围
		changeTimeRange(index) {
			this.timeIndex = index
			// 这里可以根据时间范围重新加载数据
			uni.showToast({
				title: '正在加载' + this.timeOptions[index].label + '数据',
				icon: 'none'
			})
		},
		
		// 查看详情
		viewDetail(type) {
			const typeMap = {
				'users': '用户',
				'doctors': '医生',
				'consultations': '问诊',
				'records': '病历'
			}
			uni.showToast({
				title: typeMap[type] + '详情功能开发中',
				icon: 'none'
			})
		}
	}
}
</script>

<style scoped>
.statistics-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 20rpx;
	padding-bottom: 40rpx;
}

/* 时间筛选 */
.time-filter {
	display: flex;
	gap: 15rpx;
	margin-bottom: 20rpx;
}

.time-btn {
	flex: 1;
	background: white;
	border-radius: 10rpx;
	padding: 20rpx;
	text-align: center;
	font-size: 26rpx;
	color: #666;
	transition: all 0.3s;
}

.time-btn.active {
	background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
	color: white;
	font-weight: bold;
}

/* 统计概览 */
.stats-overview {
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	gap: 15rpx;
	margin-bottom: 20rpx;
}

.stat-card {
	background: white;
	border-radius: 15rpx;
	padding: 25rpx;
	display: flex;
	gap: 20rpx;
}

.stat-icon {
	width: 80rpx;
	height: 80rpx;
	border-radius: 15rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 40rpx;
}

.user-icon {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.doctor-icon {
	background: linear-gradient(135deg, #07c160 0%, #05a64c 100%);
}

.consult-icon {
	background: linear-gradient(135deg, #ffa940 0%, #fa8c16 100%);
}

.record-icon {
	background: linear-gradient(135deg, #5cdbd3 0%, #13c2c2 100%);
}

.stat-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	justify-content: center;
}

.stat-number {
	font-size: 40rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 5rpx;
}

.stat-label {
	font-size: 22rpx;
	color: #999;
	margin-bottom: 10rpx;
}

.stat-trend {
	display: flex;
	align-items: center;
	gap: 5rpx;
}

.trend-icon {
	font-size: 20rpx;
	color: #4caf50;
}

.trend-text {
	font-size: 20rpx;
	color: #4caf50;
}

/* 图表区域 */
.chart-section {
	margin-bottom: 20rpx;
}

.section-title {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	padding: 0 10rpx 15rpx;
}

.chart-card {
	background: white;
	border-radius: 15rpx;
	padding: 30rpx;
}

.bar-chart {
	display: flex;
	align-items: flex-end;
	justify-content: space-between;
	height: 250rpx;
	padding-bottom: 50rpx;
}

.bar-item {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 10rpx;
}

.bar-wrapper {
	flex: 1;
	width: 100%;
	display: flex;
	align-items: flex-end;
	justify-content: center;
	padding: 0 10rpx;
}

.bar {
	width: 100%;
	background: linear-gradient(to top, #ff6b6b 0%, #ee5a6f 100%);
	border-radius: 8rpx 8rpx 0 0;
	min-height: 20rpx;
}

.bar-value {
	font-size: 20rpx;
	font-weight: bold;
	color: #333;
}

.bar-label {
	font-size: 20rpx;
	color: #999;
}

/* 科室统计 */
.dept-section {
	margin-bottom: 20rpx;
}

.dept-card {
	background: white;
	border-radius: 15rpx;
	padding: 25rpx;
}

.dept-item {
	margin-bottom: 25rpx;
}

.dept-item:last-child {
	margin-bottom: 0;
}

.dept-header {
	display: flex;
	justify-content: space-between;
	margin-bottom: 10rpx;
}

.dept-name {
	font-size: 26rpx;
	font-weight: bold;
	color: #333;
}

.dept-count {
	font-size: 22rpx;
	color: #999;
}

.dept-bar {
	height: 15rpx;
	background: #f0f0f0;
	border-radius: 10rpx;
	overflow: hidden;
	margin-bottom: 8rpx;
}

.dept-progress {
	height: 100%;
	background: linear-gradient(90deg, #ff6b6b 0%, #ee5a6f 100%);
	border-radius: 10rpx;
}

.dept-consult {
	font-size: 22rpx;
	color: #666;
}

/* 热门医生 */
.doctor-section {
	margin-bottom: 20rpx;
}

.doctor-card {
	background: white;
	border-radius: 15rpx;
	padding: 20rpx;
}

.doctor-rank-item {
	display: flex;
	align-items: center;
	gap: 20rpx;
	padding: 20rpx 15rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.doctor-rank-item:last-child {
	border-bottom: none;
}

.rank-badge {
	width: 50rpx;
	height: 50rpx;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 24rpx;
	font-weight: bold;
	color: white;
}

.rank-1 {
	background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
}

.rank-2 {
	background: linear-gradient(135deg, #c0c0c0 0%, #e8e8e8 100%);
}

.rank-3 {
	background: linear-gradient(135deg, #cd7f32 0%, #d4a76a 100%);
}

.rank-4, .rank-5 {
	background: #e0e0e0;
	color: #666;
}

.doctor-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	gap: 5rpx;
}

.doctor-name {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
}

.doctor-dept {
	font-size: 22rpx;
	color: #999;
}

.doctor-stats {
	display: flex;
	flex-direction: column;
	align-items: flex-end;
	gap: 5rpx;
}

.consult-count {
	font-size: 24rpx;
	color: #ff6b6b;
	font-weight: bold;
}

.rating {
	font-size: 20rpx;
	color: #ffa940;
}

/* 系统健康 */
.health-section {
	margin-bottom: 20rpx;
}

.health-card {
	background: white;
	border-radius: 15rpx;
	padding: 25rpx;
}

.health-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 20rpx 0;
	border-bottom: 1rpx solid #f0f0f0;
}

.health-item:last-child {
	border-bottom: none;
}

.health-label {
	font-size: 26rpx;
	color: #333;
}

.health-status {
	padding: 8rpx 20rpx;
	border-radius: 8rpx;
	font-size: 22rpx;
	font-weight: bold;
}

.health-good {
	background: #e8f5e9;
	color: #4caf50;
}

.health-warning {
	background: #fff3e0;
	color: #ff9800;
}
</style>
