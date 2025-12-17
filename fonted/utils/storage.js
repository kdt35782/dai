/**
 * 统一存储工具
 * 解决uni-app H5环境下刷新页面数据丢失的问题
 * 使用运行时判断代替条件编译
 */

/**
 * 判断是否为H5环境
 */
function isH5() {
	return typeof window !== 'undefined' && typeof localStorage !== 'undefined'
}

/**
 * 设置存储（同步）
 * @param {string} key 键名
 * @param {any} value 值
 */
export function setStorageSync(key, value) {
	try {
		if (isH5()) {
			// H5环境使用localStorage确保刷新不丢失
			const data = typeof value === 'object' ? JSON.stringify(value) : String(value)
			localStorage.setItem(key, data)
			console.log(`[Storage] 保存到localStorage: ${key}`)
		} else {
			// 非H5环境使用uni-app原生API
			uni.setStorageSync(key, value)
		}
	} catch (e) {
		console.error('存储失败:', e)
	}
}

/**
 * 获取存储（同步）
 * @param {string} key 键名
 * @returns {any} 存储的值
 */
export function getStorageSync(key) {
	try {
		if (isH5()) {
			// H5环境从localStorage读取
			const data = localStorage.getItem(key)
			if (!data) return null
			
			// 尝试解析JSON
			try {
				return JSON.parse(data)
			} catch {
				return data
			}
		} else {
			// 非H5环境使用uni-app原生API
			return uni.getStorageSync(key)
		}
	} catch (e) {
		console.error('读取存储失败:', e)
		return null
	}
}

/**
 * 移除存储（同步）
 * @param {string} key 键名
 */
export function removeStorageSync(key) {
	try {
		if (isH5()) {
			localStorage.removeItem(key)
			console.log(`[Storage] 从localStorage移除: ${key}`)
		} else {
			uni.removeStorageSync(key)
		}
	} catch (e) {
		console.error('移除存储失败:', e)
	}
}

/**
 * 清空所有存储（同步）
 */
export function clearStorageSync() {
	try {
		if (isH5()) {
			localStorage.clear()
		} else {
			uni.clearStorageSync()
		}
	} catch (e) {
		console.error('清空存储失败:', e)
	}
}
