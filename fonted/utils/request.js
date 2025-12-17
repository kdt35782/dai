import { API_BASE_URL, STORAGE_KEYS } from './config.js'
import { isMockEnabled, mockApiResponse } from './mock.js'
import { getStorageSync } from './storage.js'

/**
 * 统一请求封装
 */
export function request(options) {
	return new Promise((resolve, reject) => {
		// Mock模式拦截
		if (isMockEnabled()) {
			setTimeout(() => {
				try {
					const mockData = mockApiResponse(options.url, options.method || 'GET', options.data || {})
					resolve(mockData)
				} catch (error) {
					console.error('Mock数据错误:', error)
					reject(error)
				}
			}, 300) // 模拟网络延迟
			return
		}
		
		// 获取token
		const token = getStorageSync(STORAGE_KEYS.TOKEN)
		
		// 调试日志
		if (!options.noAuth) {
			console.log('[请求URL]', options.url)
			console.log('[Token状态]', token ? `存在(${token.substring(0, 20)}...)` : '不存在')
		}
		
		// 设置请求头
		const header = {
			'Content-Type': 'application/json',
			...options.header
		}
		
		// 如果有token，添加到请求头
		if (token && !options.noAuth) {
			header['Authorization'] = `Bearer ${token}`
		}
		
		// 完整URL
		const url = options.url.startsWith('http') ? options.url : API_BASE_URL + options.url
		
		// 发起请求
		uni.request({
			url,
			method: options.method || 'GET',
			data: options.data || {},
			header,
			success: (res) => {
				const data = res.data
				
				// 根据后端返回的统一格式处理
				if (data.code === 200) {
					resolve(data)
				} else {
					// 错误处理
					handleError(data)
					reject(data)
				}
			},
			fail: (err) => {
				console.error('请求失败:', err)
				uni.showToast({
					title: '网络请求失败',
					icon: 'none'
				})
				reject(err)
			}
		})
	})
}

/**
 * GET请求
 */
export function get(url, data = {}, options = {}) {
	return request({
		url,
		method: 'GET',
		data,
		...options
	})
}

/**
 * POST请求
 */
export function post(url, data = {}, options = {}) {
	return request({
		url,
		method: 'POST',
		data,
		...options
	})
}

/**
 * PUT请求
 */
export function put(url, data = {}, options = {}) {
	return request({
		url,
		method: 'PUT',
		data,
		...options
	})
}

/**
 * DELETE请求
 */
export function del(url, data = {}, options = {}) {
	return request({
		url,
		method: 'DELETE',
		data,
		...options
	})
}

/**
 * 文件上传
 */
export function uploadFile(filePath, fileType, options = {}) {
	return new Promise((resolve, reject) => {
		// Mock模式拦截
		if (isMockEnabled()) {
			setTimeout(() => {
				const mockData = mockApiResponse('/api/file/upload', 'POST', { fileType })
				resolve(mockData)
			}, 500) // 模拟上传延迟
			return
		}
		
		const token = getStorageSync(STORAGE_KEYS.TOKEN)
		
		uni.uploadFile({
			url: API_BASE_URL + '/api/file/upload',
			filePath,
			name: 'file',
			formData: {
				fileType,
				encrypt: true
			},
			header: {
				'Authorization': `Bearer ${token}`
			},
			success: (res) => {
				try {
					const data = JSON.parse(res.data)
					if (data.code === 200) {
						resolve(data)
					} else {
						handleError(data)
						reject(data)
					}
				} catch (e) {
					reject(e)
				}
			},
			fail: (err) => {
				console.error('上传失败:', err)
				uni.showToast({
					title: '文件上传失败',
					icon: 'none'
				})
				reject(err)
			}
		})
	})
}

/**
 * 统一错误处理
 */
function handleError(data) {
	const code = data.code
	const message = data.message || '请求失败'
	
	switch (code) {
		case 401:
			// 未认证，跳转到登录页
			uni.showToast({
				title: '请先登录',
				icon: 'none'
			})
			setTimeout(() => {
				uni.reLaunch({
					url: '/pages/login/login'
				})
			}, 1500)
			break
		case 403:
			uni.showToast({
				title: '无权限访问',
				icon: 'none'
			})
			break
		case 404:
			uni.showToast({
				title: '资源不存在',
				icon: 'none'
			})
			break
		case 500:
			uni.showToast({
				title: '服务器错误',
				icon: 'none'
			})
			break
		default:
			uni.showToast({
				title: message,
				icon: 'none'
			})
	}
}
