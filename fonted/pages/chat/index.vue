<template>
	<view class="chat-container">
		<!-- 顶部导航栏 -->
		<view class="chat-header">
			<view class="header-left" @click="goBack">
				<text class="iconfont icon-back">←</text>
			</view>
			<view class="header-center">
				<view class="doctor-info" v-if="consultationInfo">
					<text class="doctor-name">{{ otherUserName }}</text>
					<text class="online-status" :class="{ 'online': isOnline }">
						{{ isOnline ? '在线' : '离线' }}
					</text>
				</view>
			</view>
			<view class="header-right" @click="showMenu">
				<text class="iconfont">⋮</text>
			</view>
		</view>

		<!-- 消息列表 -->
		<scroll-view 
			class="message-list" 
			scroll-y 
			:scroll-into-view="scrollToView"
			:scroll-with-animation="true"
		>
			<view class="message-wrapper">
				<view 
					v-for="(msg, index) in messages" 
					:key="msg.messageId"
					:id="'msg-' + msg.messageId"
					class="message-item"
					:class="{ 'message-mine': msg.senderId === currentUserId }"
				>
					<!-- 时间分割线 -->
					<view v-if="showTimeLabel(index)" class="time-divider">
						<text>{{ formatTime(msg.createdAt) }}</text>
					</view>

					<!-- 消息内容 -->
					<view class="message-content">
						<!-- 对方头像 -->
						<image 
							v-if="msg.senderId !== currentUserId" 
							class="avatar" 
							:src="msg.senderAvatar || '/static/default-avatar.png'"
						></image>

						<!-- 消息气泡 -->
						<view class="message-bubble">
							<!-- 文本消息 -->
							<view v-if="msg.messageType === 1" class="text-message">
								{{ msg.content }}
							</view>

							<!-- 图片消息 -->
							<image 
								v-else-if="msg.messageType === 2" 
								class="image-message"
								:src="msg.fileUrl"
								mode="widthFix"
								@click="previewImage(msg.fileUrl)"
							></image>

							<!-- 语音消息 -->
							<view v-else-if="msg.messageType === 3" class="voice-message" @click="playVoice(msg)">
								<text class="voice-icon">🔊</text>
								<text class="voice-duration">{{ msg.duration }}"</text>
							</view>

							<!-- 处方消息 -->
							<view v-else-if="msg.messageType === 4" class="prescription-message" @click="viewPrescription(msg)">
								<text class="prescription-icon">📋</text>
								<text class="prescription-text">医生开具了处方,点击查看</text>
							</view>

							<!-- 系统消息 -->
							<view v-else-if="msg.messageType === 5" class="system-message">
								<text>{{ msg.content }}</text>
							</view>
						</view>

						<!-- 我的头像 -->
						<image 
							v-if="msg.senderId === currentUserId" 
							class="avatar avatar-mine" 
							:src="userInfo.avatar || '/static/default-avatar.png'"
						></image>
					</view>
				</view>

				<!-- 正在输入提示 -->
				<view v-if="otherUserTyping" class="typing-indicator">
					<text>对方正在输入...</text>
				</view>
			</view>
		</scroll-view>

		<!-- 底部输入栏 -->
		<view class="chat-input">
			<view class="input-toolbar">
				<view class="toolbar-btn" @click="chooseImage">
					<text class="iconfont">🖼️</text>
				</view>
				<view class="toolbar-btn" @click="recordVoice">
					<text class="iconfont">🎤</text>
				</view>
			</view>

			<view class="input-box">
				<textarea 
					v-model="inputMessage"
					class="input-area"
					placeholder="请输入消息..."
					:auto-height="true"
					:maxlength="500"
					@focus="onInputFocus"
					@blur="onInputBlur"
					@input="onTyping"
				></textarea>
			</view>

			<view class="send-btn" :class="{ 'active': canSend }" @click="sendMessage">
				发送
			</view>
		</view>
	</view>
</template>

<script>
import { API_BASE_URL, WS_BASE_URL, STORAGE_KEYS } from '@/utils/config.js';
import { getStorageSync } from '@/utils/storage.js';

export default {
	data() {
		return {
			baseUrl: API_BASE_URL,
			wsUrl: WS_BASE_URL,
			consultationId: 0,
			currentUserId: 0,
			userInfo: {},
			consultationInfo: null,
			otherUserName: '',
			isOnline: false,
			
			messages: [],
			inputMessage: '',
			scrollToView: '',
			otherUserTyping: false,
			
			// WebSocket
			socketTask: null,
			reconnectTimer: null,
			heartbeatTimer: null,
			typingTimer: null,
			
			// 分页
			page: 1,
			pageSize: 50,
			hasMore: true,
		};
	},
	
	computed: {
		canSend() {
			return this.inputMessage.trim().length > 0;
		}
	},
	
	onLoad(options) {
		this.consultationId = parseInt(options.consultationId || 0);
		
		// 从存储中获取用户信息
		this.userInfo = getStorageSync(STORAGE_KEYS.USER_INFO) || {};
		
		// 调试信息
		console.log('[聊天室] consultationId:', this.consultationId);
		console.log('[聊天室] userInfo:', this.userInfo);
		
		// 获取userId，兼容多种字段名
		this.currentUserId = this.userInfo.userId || this.userInfo.user_id || this.userInfo.id || 0;
		
		console.log('[聊天室] currentUserId:', this.currentUserId);
		
		if (!this.consultationId) {
			uni.showToast({ 
				title: '缺少问诚ID参数', 
				icon: 'none',
				duration: 2000
			});
			setTimeout(() => uni.navigateBack(), 2000);
			return;
		}
		
		if (!this.currentUserId) {
			uni.showToast({ 
				title: '用户信息失效，请重新登录', 
				icon: 'none',
				duration: 2000
			});
			setTimeout(() => {
				uni.redirectTo({ url: '/pages/login/login' });
			}, 2000);
			return;
		}
		
		this.loadConsultationInfo();
		this.loadMessages();
		this.connectWebSocket();
		this.checkOnlineStatus();
	},
	
	onUnload() {
		this.closeWebSocket();
	},
	
	methods: {
		// 加载问诊信息
		async loadConsultationInfo() {
			try {
				const token = getStorageSync(STORAGE_KEYS.TOKEN);
				const res = await uni.request({
					url: `${this.baseUrl}/api/consultation/detail`,
					method: 'GET',
					data: { consultationId: this.consultationId },
					header: { 'Authorization': `Bearer ${token}` }
				});
				
				if (res.data.code === 200) {
					this.consultationInfo = res.data.data;
					// 确定对方用户名
					if (this.currentUserId === this.consultationInfo.patientId) {
						this.otherUserName = this.consultationInfo.doctorName;
					} else {
						this.otherUserName = this.consultationInfo.patientName;
					}
				}
			} catch (error) {
				console.error('加载问诊信息失败:', error);
			}
		},
		
		// 加载消息列表
		async loadMessages(loadMore = false) {
			try {
				const token = getStorageSync(STORAGE_KEYS.TOKEN);
				const res = await uni.request({
					url: `${this.baseUrl}/api/chat/messages`,
					method: 'GET',
					data: {
						consultationId: this.consultationId,
						userId: this.currentUserId,
						page: this.page,
						pageSize: this.pageSize
					},
					header: { 'Authorization': `Bearer ${token}` }
				});
				
				if (res.data.code === 200) {
					const { messages, total } = res.data.data;
					
					if (loadMore) {
						this.messages = [...messages, ...this.messages];
					} else {
						this.messages = messages;
						// 滚动到底部
						this.$nextTick(() => {
							if (messages.length > 0) {
								this.scrollToBottom();
							}
						});
					}
					
					this.hasMore = this.messages.length < total;
				}
			} catch (error) {
				console.error('加载消息失败:', error);
			}
		},
		
		// 连接WebSocket
		connectWebSocket() {
			const wsUrl = `${this.wsUrl}/api/chat/ws?userId=${this.currentUserId}&consultationId=${this.consultationId}`;
			
			this.socketTask = uni.connectSocket({
				url: wsUrl,
				success: () => {
					console.log('WebSocket连接成功');
				},
				fail: (err) => {
					console.error('WebSocket连接失败:', err);
					this.scheduleReconnect();
				}
			});
			
			this.socketTask.onOpen(() => {
				console.log('WebSocket已打开');
				this.startHeartbeat();
			});
			
			this.socketTask.onMessage((res) => {
				this.handleWebSocketMessage(res.data);
			});
			
			this.socketTask.onError((err) => {
				console.error('WebSocket错误:', err);
			});
			
			this.socketTask.onClose(() => {
				console.log('WebSocket已关闭');
				this.stopHeartbeat();
				this.scheduleReconnect();
			});
		},
		
		// 处理WebSocket消息
		handleWebSocketMessage(data) {
			try {
				const message = JSON.parse(data);
				
				switch (message.type) {
					case 'connected':
						console.log('WebSocket连接确认');
						break;
						
					case 'chat':
						// 新消息
						this.messages.push(message.data);
						this.$nextTick(() => this.scrollToBottom());
						break;
						
					case 'status':
						// 在线状态变化
						this.isOnline = message.data.status === 'online';
						break;
						
					case 'typing':
						// 对方正在输入
						if (message.data.userId !== this.currentUserId) {
							this.otherUserTyping = true;
							clearTimeout(this.typingTimer);
							this.typingTimer = setTimeout(() => {
								this.otherUserTyping = false;
							}, 3000);
						}
						break;
						
					case 'pong':
						// 心跳响应
						break;
				}
			} catch (error) {
				console.error('解析WebSocket消息失败:', error);
			}
		},
		
		// 发送消息
		async sendMessage() {
			if (!this.canSend) return;
			
			const content = this.inputMessage.trim();
			this.inputMessage = '';
			
			try {
				const token = getStorageSync(STORAGE_KEYS.TOKEN);
				const res = await uni.request({
					url: `${this.baseUrl}/api/chat/send`,
					method: 'POST',
					data: {
						consultationId: this.consultationId,
						senderId: this.currentUserId,
						messageType: 1, // 文本消息
						content: content
					},
					header: { 'Authorization': `Bearer ${token}` }
				});
				
				if (res.data.code === 200) {
					// 消息已通过WebSocket推送,不需要手动添加
				} else {
					uni.showToast({ title: res.data.message || '发送失败', icon: 'none' });
				}
			} catch (error) {
				console.error('发送消息失败:', error);
				uni.showToast({ title: '发送失败', icon: 'none' });
			}
		},
		
		// 选择图片
		chooseImage() {
			uni.chooseImage({
				count: 1,
				success: (res) => {
					this.uploadImage(res.tempFilePaths[0]);
				}
			});
		},
		
		// 上传图片
		async uploadImage(filePath) {
			uni.showLoading({ title: '上传中...' });
			
			try {
				const token = getStorageSync(STORAGE_KEYS.TOKEN);
				const uploadRes = await uni.uploadFile({
					url: `${this.baseUrl}/api/file/upload`,
					filePath: filePath,
					name: 'file',
					header: { 'Authorization': `Bearer ${token}` }
				});
				
				const result = JSON.parse(uploadRes.data);
				if (result.code === 200) {
					// 发送图片消息
					await uni.request({
						url: `${this.baseUrl}/api/chat/send`,
						method: 'POST',
						data: {
							consultationId: this.consultationId,
							senderId: this.currentUserId,
							messageType: 2, // 图片消息
							fileUrl: result.data.url,
							fileSize: result.data.size
						},
						header: { 'Authorization': `Bearer ${token}` }
					});
				}
			} catch (error) {
				console.error('上传图片失败:', error);
				uni.showToast({ title: '上传失败', icon: 'none' });
			} finally {
				uni.hideLoading();
			}
		},
		
		// 预览图片
		previewImage(url) {
			uni.previewImage({
				urls: [url],
				current: url
			});
		},
		
		// 查看处方
		viewPrescription(msg) {
			try {
				const extraData = JSON.parse(msg.extraData || '{}');
				if (extraData.prescriptionId) {
					uni.navigateTo({
						url: `/pages/prescription/detail?prescriptionId=${extraData.prescriptionId}`
					});
				}
			} catch (error) {
				console.error('解析处方数据失败:', error);
			}
		},
		
		// 正在输入
		onTyping() {
			if (this.socketTask) {
				this.socketTask.send({
					data: JSON.stringify({
						type: 'typing',
						data: true
					})
				});
			}
		},
		
		// 检查在线状态
		async checkOnlineStatus() {
			try {
				const token = uni.getStorageSync('token');
				const res = await uni.request({
					url: `${this.baseUrl}/api/chat/online-status`,
					method: 'GET',
					data: { consultationId: this.consultationId },
					header: { 'Authorization': `Bearer ${token}` }
				});
				
				if (res.data.code === 200) {
					const onlineUsers = res.data.data.onlineUsers || [];
					this.isOnline = onlineUsers.some(id => id !== this.currentUserId);
				}
			} catch (error) {
				console.error('检查在线状态失败:', error);
			}
		},
		
		// 心跳
		startHeartbeat() {
			this.heartbeatTimer = setInterval(() => {
				if (this.socketTask) {
					this.socketTask.send({
						data: JSON.stringify({ type: 'ping' })
					});
				}
			}, 30000);
		},
		
		stopHeartbeat() {
			if (this.heartbeatTimer) {
				clearInterval(this.heartbeatTimer);
				this.heartbeatTimer = null;
			}
		},
		
		// 重连
		scheduleReconnect() {
			if (this.reconnectTimer) return;
			
			this.reconnectTimer = setTimeout(() => {
				console.log('尝试重连WebSocket...');
				this.reconnectTimer = null;
				this.connectWebSocket();
			}, 5000);
		},
		
		// 关闭WebSocket
		closeWebSocket() {
			if (this.socketTask) {
				this.socketTask.close();
				this.socketTask = null;
			}
			this.stopHeartbeat();
			if (this.reconnectTimer) {
				clearTimeout(this.reconnectTimer);
				this.reconnectTimer = null;
			}
		},
		
		// 滚动到底部
		scrollToBottom() {
			if (this.messages.length > 0) {
				const lastMsg = this.messages[this.messages.length - 1];
				this.scrollToView = 'msg-' + lastMsg.messageId;
			}
		},
		
		// 显示时间标签
		showTimeLabel(index) {
			if (index === 0) return true;
			
			const current = new Date(this.messages[index].createdAt);
			const previous = new Date(this.messages[index - 1].createdAt);
			const diff = (current - previous) / 1000 / 60; // 分钟
			
			return diff > 5; // 超过5分钟显示时间
		},
		
		// 格式化时间
		formatTime(dateStr) {
			const date = new Date(dateStr);
			const now = new Date();
			const diff = (now - date) / 1000;
			
			if (diff < 60) return '刚刚';
			if (diff < 3600) return Math.floor(diff / 60) + '分钟前';
			if (diff < 86400) return Math.floor(diff / 3600) + '小时前';
			
			const month = date.getMonth() + 1;
			const day = date.getDate();
			const hour = date.getHours().toString().padStart(2, '0');
			const minute = date.getMinutes().toString().padStart(2, '0');
			
			return `${month}月${day}日 ${hour}:${minute}`;
		},
		
		onInputFocus() {
			// 可以在这里添加逻辑
		},
		
		onInputBlur() {
			// 可以在这里添加逻辑
		},
		
		recordVoice() {
			uni.showToast({ title: '语音功能开发中', icon: 'none' });
		},
		
		playVoice(msg) {
			uni.showToast({ title: '播放语音', icon: 'none' });
		},
		
		showMenu() {
			uni.showActionSheet({
				itemList: ['查看问诊详情', '清空聊天记录'],
				success: (res) => {
					if (res.tapIndex === 0) {
						uni.navigateTo({
							url: `/pages/consultation/detail?consultationId=${this.consultationId}`
						});
					}
				}
			});
		},
		
		goBack() {
			uni.navigateBack();
		}
	}
};
</script>

<style scoped>
.chat-container {
	display: flex;
	flex-direction: column;
	height: 100vh;
	background-color: #f5f5f5;
}

.chat-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	height: 88rpx;
	padding: 0 30rpx;
	background-color: #fff;
	border-bottom: 1px solid #e5e5e5;
}

.header-left, .header-right {
	width: 80rpx;
}

.header-center {
	flex: 1;
	text-align: center;
}

.doctor-info {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.doctor-name {
	font-size: 32rpx;
	font-weight: bold;
	color: #333;
}

.online-status {
	font-size: 24rpx;
	color: #999;
	margin-top: 4rpx;
}

.online-status.online {
	color: #07c160;
}

.message-list {
	flex: 1;
	overflow-y: auto;
}

.message-wrapper {
	padding: 20rpx 30rpx;
}

.time-divider {
	text-align: center;
	margin: 20rpx 0;
}

.time-divider text {
	display: inline-block;
	padding: 8rpx 24rpx;
	font-size: 24rpx;
	color: #999;
	background-color: rgba(0, 0, 0, 0.05);
	border-radius: 8rpx;
}

.message-item {
	margin-bottom: 30rpx;
}

.message-content {
	display: flex;
	align-items: flex-end;
}

.message-mine .message-content {
	flex-direction: row-reverse;
}

.avatar {
	width: 80rpx;
	height: 80rpx;
	border-radius: 8rpx;
	flex-shrink: 0;
}

.avatar-mine {
	margin-left: 20rpx;
}

.message-bubble {
	max-width: 500rpx;
	padding: 20rpx 24rpx;
	margin: 0 20rpx;
	background-color: #fff;
	border-radius: 8rpx;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
}

.message-mine .message-bubble {
	background-color: #95ec69;
}

.text-message {
	font-size: 28rpx;
	line-height: 1.6;
	color: #333;
	word-wrap: break-word;
}

.image-message {
	max-width: 400rpx;
	border-radius: 8rpx;
}

.system-message {
	padding: 0;
	background-color: transparent;
	box-shadow: none;
	text-align: center;
	font-size: 24rpx;
	color: #999;
}

.typing-indicator {
	text-align: center;
	padding: 20rpx 0;
	font-size: 24rpx;
	color: #999;
}

.chat-input {
	display: flex;
	align-items: flex-end;
	padding: 20rpx 30rpx;
	background-color: #fff;
	border-top: 1px solid #e5e5e5;
}

.input-toolbar {
	display: flex;
	margin-right: 20rpx;
}

.toolbar-btn {
	width: 60rpx;
	height: 60rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-right: 10rpx;
	font-size: 40rpx;
}

.input-box {
	flex: 1;
	min-height: 60rpx;
	max-height: 200rpx;
	padding: 10rpx 20rpx;
	background-color: #f5f5f5;
	border-radius: 8rpx;
}

.input-area {
	width: 100%;
	font-size: 28rpx;
	line-height: 1.6;
}

.send-btn {
	width: 120rpx;
	height: 60rpx;
	margin-left: 20rpx;
	background-color: #e0e0e0;
	color: #999;
	text-align: center;
	line-height: 60rpx;
	border-radius: 8rpx;
	font-size: 28rpx;
}

.send-btn.active {
	background-color: #07c160;
	color: #fff;
}
</style>
