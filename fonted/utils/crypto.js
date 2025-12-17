/**
 * 国密加密工具类
 * 使用sm-crypto库实现真正的国密算法
 * 如果uni-app环境不支持，可以手动复制sm3.js到项目中
 */

/**
 * SM3哈希算法（内置实现）
 * 参考：https://github.com/JuneAndGreen/sm-crypto
 * @param {String} data 要哈希的数据
 * @returns {String} 64位十六进制哈希值
 */
export function sm3Hash(data) {
	// 使用内置的SM3实现
	return sm3(data)
}

/**
 * SM2加密
 * @param {String} data 要加密的数据
 * @param {String} publicKey SM2公钥
 * @returns {String} 加密后的数据
 */
export function sm2Encrypt(data, publicKey) {
	// TODO: 实际项目中使用真正的SM2算法
	// return sm.sm2.doEncrypt(data, publicKey, 1) // 1 - C1C3C2，0 - C1C2C3
	
	console.warn('[开发提示] 请在实际项目中集成真正的SM2加密算法')
	return btoa(data) // 示例：使用base64模拟
}

/**
 * SM2解密
 * @param {String} encryptedData 加密的数据
 * @param {String} privateKey SM2私钥
 * @returns {String} 解密后的数据
 */
export function sm2Decrypt(encryptedData, privateKey) {
	// TODO: 实际项目中使用真正的SM2算法
	// return sm.sm2.doDecrypt(encryptedData, privateKey, 1)
	
	console.warn('[开发提示] 请在实际项目中集成真正的SM2解密算法')
	try {
		return atob(encryptedData)
	} catch (e) {
		return encryptedData
	}
}

/**
 * SM4加密
 * @param {String} data 要加密的数据
 * @param {String} key SM4密钥（128位，16字节）
 * @returns {String} 加密后的数据
 */
export function sm4Encrypt(data, key) {
	// TODO: 实际项目中使用真正的SM4算法
	// return sm.sm4.encrypt(data, key, {mode: 'cbc', iv: 'your-iv-here'})
	
	console.warn('[开发提示] 请在实际项目中集成真正的SM4加密算法')
	return btoa(data) // 示例：使用base64模拟
}

/**
 * SM4解密
 * @param {String} encryptedData 加密的数据
 * @param {String} key SM4密钥
 * @returns {String} 解密后的数据
 */
export function sm4Decrypt(encryptedData, key) {
	// TODO: 实际项目中使用真正的SM4算法
	// return sm.sm4.decrypt(encryptedData, key, {mode: 'cbc', iv: 'your-iv-here'})
	
	console.warn('[开发提示] 请在实际项目中集成真正的SM4解密算法')
	try {
		return atob(encryptedData)
	} catch (e) {
		return encryptedData
	}
}

/**
 * 密码加密（SM3哈希）
 * @param {String} password 明文密码
 * @returns {String} SM3哈希值（64位十六进制）
 */
export function encryptPassword(password) {
	// 前端SM3哈希一次
	const hashedOnce = sm3Hash(password)
	// 返回哈希值（后端会再次加盐哈希）
	console.log('[SM3加密] 原始密码长度:', password.length)
	console.log('[SM3加密] 哈希结果:', hashedOnce)
	return hashedOnce
}

/**
 * 验证密码强度
 * @param {String} password 密码
 * @returns {Object} {valid: Boolean, message: String}
 */
export function validatePasswordStrength(password) {
	if (!password || password.length < 8) {
		return { valid: false, message: '密码长度至少8位' }
	}
	if (password.length > 20) {
		return { valid: false, message: '密码长度不超过20位' }
	}
	
	const hasUpperCase = /[A-Z]/.test(password)
	const hasLowerCase = /[a-z]/.test(password)
	const hasNumber = /\d/.test(password)
	const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/.test(password)
	
	if (!(hasUpperCase && hasLowerCase && hasNumber && hasSpecialChar)) {
		return { 
			valid: false, 
			message: '密码必须包含大小写字母、数字和特殊字符' 
		}
	}
	
	return { valid: true, message: '密码强度符合要求' }
}

// ========== SM3算法实现 ==========
// 基于sm-crypto的SM3实现
// 参考：https://github.com/JuneAndGreen/sm-crypto

/**
 * SM3哈希算法实现
 */
function sm3(str) {
	const msg = typeof str === 'string' ? stringToBytes(str) : str
	const m = preprocess(msg)
	const n = m.length / 16
	
	let V = [
		0x7380166f, 0x4914b2b9, 0x172442d7, 0xda8a0600,
		0xa96f30bc, 0x163138aa, 0xe38dee4d, 0xb0fb0e4e
	]
	
	for (let i = 0; i < n; i++) {
		const W = []
		const M = m.slice(i * 16, (i + 1) * 16)
		
		for (let j = 0; j < 16; j++) {
			W[j] = M[j]
		}
		
		for (let j = 16; j < 68; j++) {
			W[j] = P1(W[j - 16] ^ W[j - 9] ^ ROTL(W[j - 3], 15)) ^ ROTL(W[j - 13], 7) ^ W[j - 6]
		}
		
		const W1 = []
		for (let j = 0; j < 64; j++) {
			W1[j] = W[j] ^ W[j + 4]
		}
		
		let [A, B, C, D, E, F, G, H] = V
		
		for (let j = 0; j < 64; j++) {
			const SS1 = ROTL((ROTL(A, 12) + E + ROTL(T(j), j % 32)) & 0xffffffff, 7)
			const SS2 = SS1 ^ ROTL(A, 12)
			const TT1 = (FF(A, B, C, j) + D + SS2 + W1[j]) & 0xffffffff
			const TT2 = (GG(E, F, G, j) + H + SS1 + W[j]) & 0xffffffff
			
			D = C
			C = ROTL(B, 9)
			B = A
			A = TT1
			H = G
			G = ROTL(F, 19)
			F = E
			E = P0(TT2)
		}
		
		V = [
			V[0] ^ A, V[1] ^ B, V[2] ^ C, V[3] ^ D,
			V[4] ^ E, V[5] ^ F, V[6] ^ G, V[7] ^ H
		]
	}
	
	return V.map(v => ('00000000' + v.toString(16)).slice(-8)).join('')
}

function stringToBytes(str) {
	const bytes = []
	for (let i = 0; i < str.length; i++) {
		const code = str.charCodeAt(i)
		if (code < 0x80) {
			bytes.push(code)
		} else if (code < 0x800) {
			bytes.push(0xc0 | (code >> 6), 0x80 | (code & 0x3f))
		} else if (code < 0x10000) {
			bytes.push(0xe0 | (code >> 12), 0x80 | ((code >> 6) & 0x3f), 0x80 | (code & 0x3f))
		} else {
			bytes.push(
				0xf0 | (code >> 18),
				0x80 | ((code >> 12) & 0x3f),
				0x80 | ((code >> 6) & 0x3f),
				0x80 | (code & 0x3f)
			)
		}
	}
	return bytes
}

function preprocess(msg) {
	const len = msg.length * 8
	msg.push(0x80)
	
	const k = (448 - ((len + 1) % 512) + 512) % 512
	for (let i = 0; i < k / 8; i++) {
		msg.push(0x00)
	}
	
	for (let i = 0; i < 4; i++) {
		msg.push(0x00)
	}
	
	for (let i = 3; i >= 0; i--) {
		msg.push((len >> (i * 8)) & 0xff)
	}
	
	const result = []
	for (let i = 0; i < msg.length / 4; i++) {
		result[i] = (msg[i * 4] << 24) | (msg[i * 4 + 1] << 16) | (msg[i * 4 + 2] << 8) | msg[i * 4 + 3]
	}
	
	return result
}

function ROTL(x, n) {
	return ((x << n) | (x >>> (32 - n))) >>> 0
}

function FF(x, y, z, j) {
	return j < 16 ? x ^ y ^ z : (x & y) | (x & z) | (y & z)
}

function GG(x, y, z, j) {
	return j < 16 ? x ^ y ^ z : (x & y) | (~x & z)
}

function T(j) {
	return j < 16 ? 0x79cc4519 : 0x7a879d8a
}

function P0(x) {
	return x ^ ROTL(x, 9) ^ ROTL(x, 17)
}

function P1(x) {
	return x ^ ROTL(x, 15) ^ ROTL(x, 23)
}
