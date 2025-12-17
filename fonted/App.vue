<script>
import { enableMockMode } from './utils/mock.js'

export default {
	onLaunch: function() {
		console.log('App Launch')
		
		// 启用Mock模式（开发测试用）
		enableMockMode()
		
		console.log('🎉 应用启动成功')
		console.log('💡 提示：当前为Mock模式，使用模拟数据')
		console.log('👤 已自动登录为：患者账号')
		console.log('🔄 登录时使用 "doctor" 可切换为医生账号')
	},
	onShow: function() {
		console.log('App Show')
	},
	onHide: function() {
		console.log('App Hide')
	}
}
</script>

<style>
	/*每个页面公共css */
</style>
