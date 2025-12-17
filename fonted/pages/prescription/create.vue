<template>
	<view class="prescription-page">
		<!-- 诊断信息 -->
		<view class="diagnosis-card">
			<view class="card-header">
				<text class="card-title">诊断信息</text>
			</view>
			<textarea 
				class="diagnosis-input" 
				v-model="form.diagnosis" 
				placeholder="请输入诊断结果（必填）"
				maxlength="500"
			></textarea>
			<view class="char-count">{{ form.diagnosis.length }}/500</view>
		</view>
		
		<!-- AI推荐药品 -->
		<view class="recommend-card" v-if="recommendedMedicines.length > 0">
			<view class="card-header">
				<text class="card-title">🤖 AI推荐药品</text>
				<text class="card-tip">基于AI诊断分析</text>
			</view>
			
			<scroll-view class="recommend-list" scroll-x>
				<view 
					class="recommend-item" 
					v-for="medicine in recommendedMedicines" 
					:key="medicine.medicineId"
					@click="addMedicine(medicine)"
				>
					<view class="medicine-name">{{ medicine.medicineName }}</view>
					<view class="medicine-spec">{{ medicine.specification }}</view>
					<view class="medicine-price">{{ medicine.price }}元</view>
					<view class="add-icon">+</view>
				</view>
			</scroll-view>
		</view>
		
		<!-- 药品搜索 -->
		<view class="search-card">
			<view class="search-box">
				<input 
					class="search-input" 
					v-model="searchKeyword" 
					placeholder="搜索药品名称"
					@confirm="searchMedicines"
				/>
				<button class="search-btn" @click="searchMedicines">搜索</button>
			</view>
			
			<!-- 分类选择 -->
			<scroll-view class="category-tabs" scroll-x>
				<view 
					class="category-item" 
					:class="{ active: selectedCategory === '' }"
					@click="selectCategory('')"
				>
					全部
				</view>
				<view 
					class="category-item" 
					:class="{ active: selectedCategory === category }"
					v-for="category in categories" 
					:key="category"
					@click="selectCategory(category)"
				>
					{{ category }}
				</view>
			</scroll-view>
		</view>
		
		<!-- 搜索结果 -->
		<view class="search-results" v-if="searchResults.length > 0">
			<view class="result-header">
				<text class="result-title">搜索结果</text>
				<text class="result-count">共{{ searchResults.length }}种药品</text>
			</view>
			
			<view 
				class="medicine-item" 
				v-for="medicine in searchResults" 
				:key="medicine.medicineId"
				@click="addMedicine(medicine)"
			>
				<view class="medicine-info">
					<view class="medicine-name">{{ medicine.medicineName }}</view>
					<view class="medicine-detail">
						<text class="medicine-category">{{ medicine.category }}</text>
						<text class="medicine-spec">{{ medicine.specification }}</text>
					</view>
					<view class="medicine-manufacturer">{{ medicine.manufacturer }}</view>
				</view>
				<view class="medicine-right">
					<view class="medicine-price">{{ medicine.price }}元/{{ medicine.unit }}</view>
					<view class="add-btn">添加</view>
				</view>
			</view>
		</view>
		
		<!-- 已选药品 -->
		<view class="selected-card" v-if="selectedMedicines.length > 0">
			<view class="card-header">
				<text class="card-title">已选药品 ({{ selectedMedicines.length }})</text>
				<text class="total-price">总计: {{ totalPrice.toFixed(2) }}元</text>
			</view>
			
			<view 
				class="selected-item" 
				v-for="(item, index) in selectedMedicines" 
				:key="index"
			>
				<view class="selected-info">
					<view class="selected-name">{{ item.medicineName }}</view>
					<view class="selected-spec">{{ item.specification }}</view>
					
					<!-- 用药详情 -->
					<view class="dosage-form">
						<view class="form-row">
							<text class="form-label">用法:</text>
							<input 
								class="form-input" 
								v-model="item.usage" 
								placeholder="如: 口服"
							/>
						</view>
						<view class="form-row">
							<text class="form-label">频次:</text>
							<input 
								class="form-input" 
								v-model="item.frequency" 
								placeholder="如: 每日3次"
							/>
						</view>
						<view class="form-row">
							<text class="form-label">剂量:</text>
							<input 
								class="form-input" 
								v-model="item.dosage" 
								placeholder="如: 1片"
							/>
						</view>
						<view class="form-row">
							<text class="form-label">疗程:</text>
							<input 
								class="form-input" 
								v-model="item.duration" 
								placeholder="如: 7天"
							/>
						</view>
						<view class="form-row">
							<text class="form-label">数量:</text>
							<view class="quantity-control">
								<text class="quantity-btn" @click="decreaseQuantity(index)">-</text>
								<text class="quantity-value">{{ item.quantity }}</text>
								<text class="quantity-btn" @click="increaseQuantity(index)">+</text>
							</view>
						</view>
						<view class="form-row">
							<text class="form-label">备注:</text>
							<input 
								class="form-input" 
								v-model="item.notes" 
								placeholder="特殊说明(可选)"
							/>
						</view>
					</view>
				</view>
				
				<view class="selected-actions">
					<text class="item-price">{{ (item.unitPrice * item.quantity).toFixed(2) }}元</text>
					<text class="remove-btn" @click="removeMedicine(index)">删除</text>
				</view>
			</view>
		</view>
		
		<!-- 提交按钮 -->
		<view class="submit-bar">
			<button 
				class="submit-btn" 
				@click="submitPrescription"
				:disabled="!canSubmit"
			>
				开具处方
			</button>
		</view>
	</view>
</template>

<script>
import { get, post } from '@/utils/request.js'
import { API } from '@/utils/config.js'

export default {
	data() {
		return {
			consultationId: 0,
			aiDiagnosis: '',
			
			form: {
				diagnosis: ''
			},
			
			// 药品分类
			categories: ['感冒药', '止咳化痰', '消化系统', '心血管', '降压药', '降糖药', '其他'],
			selectedCategory: '',
			
			// 搜索
			searchKeyword: '',
			searchResults: [],
			recommendedMedicines: [],
			
			// 已选药品
			selectedMedicines: []
		}
	},
	
	computed: {
		totalPrice() {
			return this.selectedMedicines.reduce((sum, item) => {
				return sum + (item.unitPrice * item.quantity)
			}, 0)
		},
		
		canSubmit() {
			return this.form.diagnosis.trim() && this.selectedMedicines.length > 0
		}
	},
	
	onLoad(options) {
		if (options.consultationId) {
			this.consultationId = parseInt(options.consultationId)
		}
		if (options.aiDiagnosis) {
			this.aiDiagnosis = decodeURIComponent(options.aiDiagnosis)
			this.loadRecommendedMedicines()
		}
		
		// 页面加载时自动搜索所有药品
		this.searchMedicines()
	},
	
	methods: {
		// 加载AI推荐药品
		async loadRecommendedMedicines() {
			if (!this.aiDiagnosis) return
			
			try {
				const res = await post(API.PRESCRIPTION_RECOMMEND, {
					aiDiagnosis: this.aiDiagnosis
				})
				
				this.recommendedMedicines = res.data || []
				
			} catch (error) {
				console.error('加载推荐药品失败:', error)
			}
		},
		
		// 选择分类
		selectCategory(category) {
			this.selectedCategory = category
			this.searchMedicines()
		},
		
		// 搜索药品
		async searchMedicines() {
			uni.showLoading({ title: '搜索中...' })
			
			try {
				console.log('开始搜索药品', {
					keyword: this.searchKeyword,
					category: this.selectedCategory
				})
				
				const res = await get(API.PRESCRIPTION_SEARCH_MEDICINES, {
					keyword: this.searchKeyword,
					category: this.selectedCategory,
					page: 1,
					pageSize: 20
				})
				
				console.log('搜索结果:', res)
				
				if (res && res.data) {
					this.searchResults = res.data.list || []
					
					if (this.searchResults.length === 0) {
						uni.showToast({
							title: '未找到相关药品',
							icon: 'none'
						})
					}
				} else {
					this.searchResults = []
				}
				
			} catch (error) {
				console.error('搜索药品失败:', error)
				
				uni.showToast({
					title: error.message || '搜索失败,请检查后端服务',
					icon: 'none',
					duration: 3000
				})
				
				this.searchResults = []
			} finally {
				uni.hideLoading()
			}
		},
		
		// 添加药品
		addMedicine(medicine) {
			// 检查是否已添加
			const exists = this.selectedMedicines.find(item => 
				item.medicineId === medicine.medicineId
			)
			
			if (exists) {
				uni.showToast({
					title: '该药品已添加',
					icon: 'none'
				})
				return
			}
			
			// 添加到已选列表
			this.selectedMedicines.push({
				medicineId: medicine.medicineId,
				medicineName: medicine.medicineName,
				specification: medicine.specification,
				unit: medicine.unit || '盒',
				unitPrice: medicine.price,
				quantity: 1,
				usage: '口服',
				frequency: '每日3次',
				dosage: '1片',
				duration: '7天',
				notes: ''
			})
			
			uni.showToast({
				title: '已添加',
				icon: 'success',
				duration: 1000
			})
		},
		
		// 删除药品
		removeMedicine(index) {
			uni.showModal({
				title: '确认删除',
				content: '确定要删除这个药品吗？',
				success: (res) => {
					if (res.confirm) {
						this.selectedMedicines.splice(index, 1)
					}
				}
			})
		},
		
		// 增加数量
		increaseQuantity(index) {
			this.selectedMedicines[index].quantity++
		},
		
		// 减少数量
		decreaseQuantity(index) {
			if (this.selectedMedicines[index].quantity > 1) {
				this.selectedMedicines[index].quantity--
			}
		},
		
		// 提交处方
		async submitPrescription() {
			if (!this.canSubmit) return
			
			uni.showLoading({ title: '提交中...' })
			
			try {
				// 准备处方数据
				const medicines = this.selectedMedicines.map(item => ({
					medicineId: item.medicineId,
					quantity: item.quantity,
					usage: item.usage,
					frequency: item.frequency,
					dosage: item.dosage,
					duration: item.duration,
					notes: item.notes
				}))
				
				// 调用完成问诊API，包含处方数据
				await post(API.CONSULTATION_FINISH, {
					consultationId: this.consultationId,
					diagnosis: this.form.diagnosis,
					prescription: medicines
				})
				
				uni.hideLoading()
				
				uni.showToast({
					title: '处方已开具',
					icon: 'success'
				})
				
				setTimeout(() => {
					uni.navigateBack()
				}, 1500)
				
			} catch (error) {
				uni.hideLoading()
				console.error('开具处方失败:', error)
				uni.showToast({
					title: error.message || '开具失败',
					icon: 'none'
				})
			}
		}
	}
}
</script>

<style scoped>
.prescription-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding-bottom: 120rpx;
}

.diagnosis-card,
.recommend-card,
.search-card,
.selected-card {
	background: white;
	margin: 20rpx 30rpx;
	border-radius: 15rpx;
	padding: 30rpx;
}

.card-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20rpx;
}

.card-title {
	font-size: 30rpx;
	font-weight: bold;
	color: #333;
}

.card-tip {
	font-size: 22rpx;
	color: #999;
}

.diagnosis-input {
	width: 100%;
	min-height: 150rpx;
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 20rpx;
	font-size: 28rpx;
	color: #333;
}

.char-count {
	text-align: right;
	font-size: 22rpx;
	color: #999;
	margin-top: 10rpx;
}

/* 推荐药品 */
.recommend-list {
	white-space: nowrap;
}

.recommend-item {
	display: inline-block;
	width: 200rpx;
	padding: 20rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	border-radius: 10rpx;
	margin-right: 20rpx;
	position: relative;
}

.recommend-item .medicine-name {
	font-size: 26rpx;
	color: white;
	font-weight: bold;
	margin-bottom: 10rpx;
}

.recommend-item .medicine-spec {
	font-size: 22rpx;
	color: rgba(255, 255, 255, 0.8);
	margin-bottom: 10rpx;
}

.recommend-item .medicine-price {
	font-size: 28rpx;
	color: #fff;
	font-weight: bold;
}

.recommend-item .add-icon {
	position: absolute;
	top: 10rpx;
	right: 10rpx;
	width: 40rpx;
	height: 40rpx;
	line-height: 40rpx;
	text-align: center;
	background: rgba(255, 255, 255, 0.3);
	border-radius: 50%;
	color: white;
	font-size: 30rpx;
}

/* 搜索 */
.search-box {
	display: flex;
	align-items: center;
	margin-bottom: 20rpx;
}

.search-input {
	flex: 1;
	height: 70rpx;
	background: #f5f5f5;
	border-radius: 35rpx;
	padding: 0 25rpx;
	font-size: 28rpx;
	margin-right: 15rpx;
}

.search-btn {
	width: 120rpx;
	height: 70rpx;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
	border: none;
	border-radius: 35rpx;
	font-size: 28rpx;
	padding: 0;
	line-height: 70rpx;
}

.search-btn::after {
	border: none;
}

.category-tabs {
	white-space: nowrap;
}

.category-item {
	display: inline-block;
	padding: 10rpx 25rpx;
	background: #f5f5f5;
	border-radius: 30rpx;
	font-size: 24rpx;
	color: #666;
	margin-right: 15rpx;
}

.category-item.active {
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
}

/* 搜索结果 */
.search-results {
	margin: 20rpx 30rpx;
}

.result-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 15rpx;
}

.result-title {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
}

.result-count {
	font-size: 24rpx;
	color: #999;
}

.medicine-item {
	background: white;
	border-radius: 15rpx;
	padding: 25rpx;
	margin-bottom: 15rpx;
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.medicine-info {
	flex: 1;
}

.medicine-name {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 10rpx;
}

.medicine-detail {
	font-size: 22rpx;
	color: #999;
	margin-bottom: 5rpx;
}

.medicine-category {
	margin-right: 15rpx;
}

.medicine-manufacturer {
	font-size: 22rpx;
	color: #999;
}

.medicine-right {
	display: flex;
	flex-direction: column;
	align-items: flex-end;
}

.medicine-price {
	font-size: 28rpx;
	color: #07c160;
	font-weight: bold;
	margin-bottom: 10rpx;
}

.add-btn {
	padding: 5rpx 20rpx;
	background: linear-gradient(135deg, #07c160 0%, #05a04e 100%);
	color: white;
	border-radius: 20rpx;
	font-size: 22rpx;
}

/* 已选药品 */
.total-price {
	font-size: 28rpx;
	color: #07c160;
	font-weight: bold;
}

.selected-item {
	background: #f5f5f5;
	border-radius: 10rpx;
	padding: 20rpx;
	margin-bottom: 15rpx;
}

.selected-name {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 5rpx;
}

.selected-spec {
	font-size: 22rpx;
	color: #999;
	margin-bottom: 15rpx;
}

.dosage-form {
	margin-bottom: 10rpx;
}

.form-row {
	display: flex;
	align-items: center;
	margin-bottom: 15rpx;
}

.form-label {
	width: 100rpx;
	font-size: 24rpx;
	color: #666;
}

.form-input {
	flex: 1;
	height: 60rpx;
	background: white;
	border-radius: 8rpx;
	padding: 0 15rpx;
	font-size: 24rpx;
}

.quantity-control {
	display: flex;
	align-items: center;
	background: white;
	border-radius: 8rpx;
	overflow: hidden;
}

.quantity-btn {
	width: 60rpx;
	height: 60rpx;
	line-height: 60rpx;
	text-align: center;
	font-size: 30rpx;
	color: #07c160;
}

.quantity-value {
	width: 80rpx;
	text-align: center;
	font-size: 26rpx;
	color: #333;
}

.selected-actions {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-top: 15rpx;
}

.item-price {
	font-size: 26rpx;
	color: #07c160;
	font-weight: bold;
}

.remove-btn {
	padding: 5rpx 15rpx;
	background: #ff3b30;
	color: white;
	border-radius: 15rpx;
	font-size: 22rpx;
}

/* 提交按钮 */
.submit-bar {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: white;
	padding: 20rpx 30rpx;
	border-top: 1px solid #f0f0f0;
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
}

.submit-btn::after {
	border: none;
}

.submit-btn[disabled] {
	background: #e0e0e0;
	color: #999;
}
</style>
