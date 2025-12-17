-- =====================================================
-- 聊天消息数据库脚本
-- 版本: 1.0
-- 创建日期: 2025-12-17
-- =====================================================

USE SM;

-- 聊天消息表
CREATE TABLE IF NOT EXISTS SM_chat_message (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '消息ID',
    message_no VARCHAR(50) NOT NULL UNIQUE COMMENT '消息编号',
    consultation_id BIGINT NOT NULL COMMENT '问诊ID',
    sender_id BIGINT NOT NULL COMMENT '发送者ID',
    receiver_id BIGINT NOT NULL COMMENT '接收者ID',
    
    -- 消息内容
    message_type TINYINT NOT NULL DEFAULT 1 COMMENT '消息类型(1:文本,2:图片,3:语音,4:处方,5:系统)',
    content TEXT COMMENT '消息内容(SM4加密)',
    file_url VARCHAR(500) COMMENT '文件URL(图片/语音)',
    file_size INT COMMENT '文件大小(字节)',
    duration INT COMMENT '语音时长(秒)',
    
    -- 扩展数据(JSON格式,存储处方ID等)
    extra_data TEXT COMMENT '扩展数据',
    
    -- 状态
    is_read TINYINT DEFAULT 0 COMMENT '是否已读(0:未读,1:已读)',
    read_at TIMESTAMP NULL COMMENT '已读时间',
    is_deleted TINYINT DEFAULT 0 COMMENT '是否删除(0:否,1:是)',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    -- 索引
    FOREIGN KEY (consultation_id) REFERENCES SM_consultation(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    
    INDEX idx_message_no (message_no),
    INDEX idx_consultation_id (consultation_id),
    INDEX idx_sender_id (sender_id),
    INDEX idx_receiver_id (receiver_id),
    INDEX idx_created_at (created_at),
    INDEX idx_is_read (is_read)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天消息表';

-- 未读消息统计表(优化查询性能)
CREATE TABLE IF NOT EXISTS SM_chat_unread_count (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    consultation_id BIGINT NOT NULL COMMENT '问诊ID',
    unread_count INT DEFAULT 0 COMMENT '未读消息数',
    last_message_id BIGINT COMMENT '最后一条消息ID',
    last_message_time TIMESTAMP NULL COMMENT '最后消息时间',
    
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    UNIQUE KEY uk_user_consultation (user_id, consultation_id),
    INDEX idx_user_id (user_id),
    INDEX idx_last_message_time (last_message_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='未读消息统计表';

-- =====================================================
-- 初始化完成
-- =====================================================
