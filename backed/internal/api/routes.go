package api

import (
	"sm-medical/internal/api/handler"
	"sm-medical/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API根路径
	api := r.Group("/api")

	// 用户模块
	userHandler := handler.NewUserHandler()
	user := api.Group("/user")
	{
		// 公开接口
		user.POST("/register", userHandler.Register)
		user.POST("/login", userHandler.Login)
		user.GET("/doctors", userHandler.GetDoctors)
		user.GET("/doctor/:userId", userHandler.GetDoctorDetail)

		// 需要认证的接口
		auth := user.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/info", userHandler.GetUserInfo) // 获取当前用户信息，需要登录
			auth.PUT("/profile", userHandler.UpdateProfile)
			auth.PUT("/password", userHandler.ChangePassword)
			// 以下接口已废弃，现在区分医生和患者在注册页面完成
			// auth.POST("/apply-doctor", userHandler.ApplyDoctor)
			// auth.GET("/doctor-application", userHandler.GetDoctorApplication)
			auth.POST("/logout", userHandler.Logout)
		}
	}

	// 问诊模块
	consultationHandler := handler.NewConsultationHandler()
	consultation := api.Group("/consultation")
	consultation.Use(middleware.AuthMiddleware())
	{
		consultation.POST("/create", consultationHandler.Create)
		consultation.GET("/list", consultationHandler.GetList)
		consultation.GET("/detail", consultationHandler.GetDetail)
		consultation.POST("/accept", consultationHandler.Accept)
		consultation.POST("/finish", consultationHandler.Finish)
	}

	// 病历模块
	recordHandler := handler.NewRecordHandler()
	record := api.Group("/record")
	record.Use(middleware.AuthMiddleware())
	{
		record.GET("/list", recordHandler.GetList)
		record.GET("/detail", recordHandler.GetDetail)
	}

	// 通知模块
	notificationHandler := handler.NewNotificationHandler()
	notification := api.Group("/notification")
	notification.Use(middleware.AuthMiddleware())
	{
		notification.GET("/list", notificationHandler.GetList)
		notification.GET("/unread-count", notificationHandler.GetUnreadCount)
		notification.PUT("/mark-read", notificationHandler.MarkRead)
	}

	// 密钥管理
	keyHandler := handler.NewKeyHandler()
	key := api.Group("/key")
	{
		key.POST("/generate", keyHandler.Generate).Use(middleware.AuthMiddleware())
	}

	// 文件上传
	fileHandler := handler.NewFileHandler()
	file := api.Group("/file")
	file.Use(middleware.AuthMiddleware())
	{
		file.POST("/upload", fileHandler.Upload)
	}

	// 处方和药品模块
	prescriptionHandler := handler.NewPrescriptionHandler()
	prescription := api.Group("/prescription")
	prescription.Use(middleware.AuthMiddleware())
	{
		prescription.GET("/medicines/search", prescriptionHandler.SearchMedicines)           // 搜索药品
		prescription.POST("/medicines/recommend", prescriptionHandler.GetRecommendedMedicines) // AI推荐药品
		prescription.GET("/:prescriptionId", prescriptionHandler.GetPrescriptionDetail)       // 获取处方详情
	}

	// 聊天模块
	handler.InitChatHandler() // 初始化聊天服务
	chat := api.Group("/chat")
	{
		// WebSocket连接(不需要认证中间件,在连接时验证)
		chat.GET("/ws", handler.WebSocketConnect)
		
		// 需要认证的接口
		authChat := chat.Group("")
		authChat.Use(middleware.AuthMiddleware())
		{
			authChat.POST("/send", handler.SendMessage)              // 发送消息
			authChat.GET("/messages", handler.GetMessageList)        // 获取消息列表
			authChat.GET("/unread-count", handler.GetUnreadCount)    // 获取未读数量
			authChat.PUT("/mark-read", handler.MarkAsRead)           // 标记已读
			authChat.GET("/online-status", handler.GetOnlineStatus)  // 获取在线状态
			authChat.POST("/typing", handler.SendTypingStatus)       // 发送正在输入状态
		}
	}

	// 管理员模块
	adminHandler := handler.NewAdminHandler()
	admin := api.Group("/user/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		// 以下接口已废弃，现在医生注册直接生效，不需要审核
		// admin.PUT("/review-doctor", adminHandler.ReviewDoctorApplication)
		// admin.GET("/doctor-applications", adminHandler.GetDoctorApplications)
		admin.GET("/users", adminHandler.GetUsers)
		admin.PUT("/status", adminHandler.UpdateUserStatus)
		admin.GET("/login-logs", adminHandler.GetLoginLogs)
	}

	// 智能分诊模块
	triageHandler := handler.NewTriageHandler()
	triage := api.Group("/triage")
	triage.Use(middleware.AuthMiddleware())
	{
		triage.PUT("/online-status", triageHandler.UpdateOnlineStatus)      // 更新医生在线状态
		triage.POST("/auto-assign", triageHandler.ManualAssignDoctor)       // 手动触发智能分诊
		triage.GET("/workload", triageHandler.GetDoctorWorkload)           // 获取医生负载信息
	}
}
