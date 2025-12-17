-- =============================================
-- 基于国密加密的网上问诊平台 - 完整数据库
-- 数据库版本: MySQL 8.0
-- 创建日期: 2024-12-03
-- 更新日期: 2024-12-03
-- =============================================

-- 创建数据库
CREATE DATABASE IF NOT EXISTS SM DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE SM;

-- =============================================
-- 1. 用户表
-- =============================================
CREATE TABLE IF NOT EXISTS SM_user (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(128) NOT NULL COMMENT '密码(SM3加密)',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱(SM4加密)',
    phone VARCHAR(255) COMMENT '手机号(SM4加密)',
    real_name VARCHAR(255) COMMENT '真实姓名(SM4加密)',
    id_card VARCHAR(255) COMMENT '身份证号(SM4加密)',
    identify ENUM('user', 'admin', 'doctor') NOT NULL DEFAULT 'user' COMMENT '身份标识',
    avatar VARCHAR(500) COMMENT '头像URL',
    gender TINYINT DEFAULT 0 COMMENT '性别(0:未知,1:男,2:女)',
    age INT COMMENT '年龄',
    status TINYINT DEFAULT 0 COMMENT '账号状态(0:正常,1:禁用,2:待审核)',
    doctor_cert VARCHAR(500) COMMENT '医生资格证书URL',
    doctor_title VARCHAR(50) COMMENT '医生职称',
    doctor_dept VARCHAR(50) COMMENT '所属科室',
    doctor_intro TEXT COMMENT '医生简介',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    last_login_time TIMESTAMP NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(255) COMMENT '最后登录IP(SM4加密)',
    INDEX idx_identify (identify),
    INDEX idx_status (status),
    INDEX idx_doctor_dept (doctor_dept)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- =============================================
-- 2. 医生申请记录表
-- =============================================
CREATE TABLE IF NOT EXISTS SM_doctor_application (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '申请ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    real_name VARCHAR(255) NOT NULL COMMENT '真实姓名(SM4加密)',
    id_card VARCHAR(255) NOT NULL COMMENT '身份证号(SM4加密)',
    phone VARCHAR(255) NOT NULL COMMENT '手机号(SM4加密)',
    doctor_cert VARCHAR(500) NOT NULL COMMENT '医生资格证书URL',
    doctor_title VARCHAR(50) NOT NULL COMMENT '医生职称',
    doctor_dept VARCHAR(50) NOT NULL COMMENT '所属科室',
    doctor_intro TEXT COMMENT '医生简介',
    status TINYINT DEFAULT 0 COMMENT '审核状态(0:待审核,1:已通过,2:已拒绝)',
    reject_reason TEXT COMMENT '拒绝原因',
    reviewer_id BIGINT COMMENT '审核人ID',
    review_time TIMESTAMP NULL COMMENT '审核时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='医生申请记录表';

-- =============================================
-- 3. 用户登录日志表
-- =============================================
CREATE TABLE IF NOT EXISTS SM_login_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
    user_id BIGINT COMMENT '用户ID',
    username VARCHAR(50) COMMENT '用户名',
    login_ip VARCHAR(255) NOT NULL COMMENT '登录IP(SM4加密)',
    login_location VARCHAR(100) COMMENT '登录地点',
    browser VARCHAR(50) COMMENT '浏览器类型',
    os VARCHAR(50) COMMENT '操作系统',
    status TINYINT NOT NULL COMMENT '登录状态(0:失败,1:成功)',
    msg VARCHAR(255) COMMENT '提示消息',
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    INDEX idx_user_id (user_id),
    INDEX idx_login_time (login_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户登录日志表';

-- =============================================
-- 4. 问诊记录表
-- =============================================
CREATE TABLE IF NOT EXISTS SM_consultation (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '问诊ID',
    patient_id BIGINT NOT NULL COMMENT '患者ID',
    doctor_id BIGINT COMMENT '医生ID',
    consultation_no VARCHAR(50) NOT NULL UNIQUE COMMENT '问诊编号',
    chief_complaint TEXT COMMENT '主诉(SM4加密)',
    symptoms_encrypted TEXT COMMENT '症状特征(Paillier加密,用于AI)',
    ai_risk_score INT COMMENT 'AI风险评分(0-100)',
    ai_diagnosis TEXT COMMENT 'AI初步诊断建议',
    doctor_diagnosis TEXT COMMENT '医生诊断(SM4加密)',
    prescription TEXT COMMENT '处方信息(SM4加密)',
    status TINYINT DEFAULT 0 COMMENT '状态(0:待接诊,1:问诊中,2:已完成,3:已取消)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    FOREIGN KEY (patient_id) REFERENCES SM_user(id),
    FOREIGN KEY (doctor_id) REFERENCES SM_user(id),
    INDEX idx_patient_id (patient_id),
    INDEX idx_doctor_id (doctor_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='问诊记录表';

-- =============================================
-- 5. 问诊消息表
-- =============================================
CREATE TABLE IF NOT EXISTS SM_consultation_message (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '消息ID',
    consultation_id BIGINT NOT NULL COMMENT '问诊ID',
    sender_id BIGINT NOT NULL COMMENT '发送者ID',
    receiver_id BIGINT NOT NULL COMMENT '接收者ID',
    message_type TINYINT NOT NULL COMMENT '消息类型(1:文本,2:图片,3:语音,4:处方)',
    content TEXT NOT NULL COMMENT '消息内容(SM4加密)',
    file_url VARCHAR(500) COMMENT '文件URL',
    is_read TINYINT DEFAULT 0 COMMENT '是否已读(0:未读,1:已读)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
    FOREIGN KEY (consultation_id) REFERENCES SM_consultation(id) ON DELETE CASCADE,
    INDEX idx_consultation_id (consultation_id),
    INDEX idx_sender_id (sender_id),
    INDEX idx_receiver_id (receiver_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='问诊消息表';

-- =============================================
-- 6. 电子病历表
-- =============================================
CREATE TABLE IF NOT EXISTS SM_medical_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '病历ID',
    record_no VARCHAR(50) NOT NULL UNIQUE COMMENT '病历编号',
    patient_id BIGINT NOT NULL COMMENT '患者ID',
    consultation_id BIGINT COMMENT '关联问诊ID',
    record_type TINYINT NOT NULL COMMENT '病历类型(1:门诊,2:在线问诊)',
    chief_complaint TEXT COMMENT '主诉(SM4加密)',
    present_illness TEXT COMMENT '现病史(SM4加密)',
    past_history TEXT COMMENT '既往史(SM4加密)',
    diagnosis TEXT COMMENT '诊断(SM4加密)',
    treatment_plan TEXT COMMENT '治疗方案(SM4加密)',
    doctor_id BIGINT COMMENT '诊疗医生ID',
    access_policy TEXT COMMENT '访问策略(SM9属性基加密)',
    data_hash VARCHAR(128) COMMENT '数据完整性哈希(SM3)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (patient_id) REFERENCES SM_user(id),
    FOREIGN KEY (consultation_id) REFERENCES SM_consultation(id),
    FOREIGN KEY (doctor_id) REFERENCES SM_user(id),
    INDEX idx_patient_id (patient_id),
    INDEX idx_consultation_id (consultation_id),
    INDEX idx_record_no (record_no),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电子病历表';

-- =============================================
-- 7. 病历访问日志表
-- =============================================
CREATE TABLE IF NOT EXISTS SM_record_access_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
    record_id BIGINT NOT NULL COMMENT '病历ID',
    accessor_id BIGINT NOT NULL COMMENT '访问者ID',
    accessor_name VARCHAR(50) COMMENT '访问者姓名',
    access_type TINYINT NOT NULL COMMENT '访问类型(1:查看,2:编辑,3:下载)',
    access_result TINYINT NOT NULL COMMENT '访问结果(0:失败,1:成功)',
    failure_reason VARCHAR(255) COMMENT '失败原因',
    ip_address VARCHAR(255) COMMENT '访问IP(SM4加密)',
    access_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '访问时间',
    FOREIGN KEY (record_id) REFERENCES SM_medical_record(id) ON DELETE CASCADE,
    INDEX idx_record_id (record_id),
    INDEX idx_accessor_id (accessor_id),
    INDEX idx_access_time (access_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='病历访问日志表';

-- =============================================
-- 8. 系统通知表
-- =============================================
CREATE TABLE IF NOT EXISTS SM_notification (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '通知ID',
    user_id BIGINT NOT NULL COMMENT '接收用户ID',
    notification_type TINYINT NOT NULL COMMENT '通知类型(1:系统,2:问诊,3:审核,4:提醒)',
    title VARCHAR(100) NOT NULL COMMENT '通知标题',
    content TEXT COMMENT '通知内容',
    related_id BIGINT COMMENT '关联ID（问诊ID、申请ID等）',
    related_type VARCHAR(50) COMMENT '关联类型',
    is_read TINYINT DEFAULT 0 COMMENT '是否已读(0:未读,1:已读)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    read_at TIMESTAMP NULL COMMENT '读取时间',
    FOREIGN KEY (user_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统通知表';

-- =============================================
-- 9. 文件管理表
-- =============================================
CREATE TABLE IF NOT EXISTS SM_file (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '文件ID',
    file_id VARCHAR(50) NOT NULL UNIQUE COMMENT '文件编号',
    user_id BIGINT NOT NULL COMMENT '上传用户ID',
    original_name VARCHAR(255) NOT NULL COMMENT '原始文件名',
    file_type VARCHAR(50) NOT NULL COMMENT '文件类型(avatar,cert,medical)',
    file_size BIGINT NOT NULL COMMENT '文件大小(字节)',
    mime_type VARCHAR(100) COMMENT 'MIME类型',
    storage_path VARCHAR(500) NOT NULL COMMENT '存储路径',
    file_url VARCHAR(500) COMMENT '访问URL',
    is_encrypted TINYINT DEFAULT 1 COMMENT '是否加密(0:否,1:是)',
    file_hash VARCHAR(128) COMMENT '文件哈希(SM3)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
    FOREIGN KEY (user_id) REFERENCES SM_user(id),
    INDEX idx_user_id (user_id),
    INDEX idx_file_type (file_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件管理表';

-- =============================================
-- 初始化数据
-- =============================================

-- 插入默认管理员账号
-- 用户名: admin
-- 密码: Admin123!@# (需要经过SM3加密)
-- 注意: 实际使用时,password字段应存储经过SM3加密的密文
INSERT INTO SM_user (username, password, email, phone, identify, status, created_at) 
VALUES (
    'admin',
    'SM3_ENCRYPTED_PASSWORD_HERE',  -- 请替换为实际的SM3加密后的密码
    'SM4_ENCRYPTED_EMAIL_HERE',      -- 请替换为实际的SM4加密后的邮箱
    'SM4_ENCRYPTED_PHONE_HERE',      -- 请替换为实际的SM4加密后的手机号
    'admin',
    0,
    CURRENT_TIMESTAMP
);

-- 插入测试医生账号（可选）
-- INSERT INTO SM_user (username, password, email, phone, real_name, identify, doctor_title, doctor_dept, status) 
-- VALUES (
--     'doctor_test',
--     'SM3_ENCRYPTED_PASSWORD',
--     'SM4_ENCRYPTED_EMAIL',
--     'SM4_ENCRYPTED_PHONE',
--     'SM4_ENCRYPTED_NAME',
--     'doctor',
--     '主治医师',
--     '心内科',
--     0
-- );

-- =============================================
-- 查询示例
-- =============================================

-- 查询所有普通用户
-- SELECT * FROM SM_user WHERE identify = 'user' AND status = 0;

-- 查询所有医生
-- SELECT * FROM SM_user WHERE identify = 'doctor' AND status = 0;

-- 查询待审核的医生申请
-- SELECT * FROM SM_doctor_application WHERE status = 0 ORDER BY created_at DESC;

-- 查询指定用户的登录日志
-- SELECT * FROM SM_login_log WHERE user_id = 1 ORDER BY login_time DESC LIMIT 10;

-- 查询待接诊的问诊
-- SELECT * FROM SM_consultation WHERE status = 0 ORDER BY created_at DESC;

-- 查询指定患者的病历
-- SELECT * FROM SM_medical_record WHERE patient_id = 1 ORDER BY created_at DESC;

-- 查询病历访问日志
-- SELECT * FROM SM_record_access_log WHERE record_id = 1 ORDER BY access_time DESC;

-- 查询用户未读通知
-- SELECT * FROM SM_notification WHERE user_id = 1 AND is_read = 0 ORDER BY created_at DESC;

-- 查询问诊消息记录
-- SELECT * FROM SM_consultation_message WHERE consultation_id = 1 ORDER BY created_at ASC;

-- =============================================
-- 国密加密说明
-- =============================================
-- 1. SM3加密(用于密码):
--    - 密码在前端使用SM3进行哈希后传输
--    - 后端接收到密码后再次使用SM3+盐值进行哈希存储
--    - 双重加密确保密码安全
--
-- 2. SM4加密(用于敏感数据):
--    - 邮箱、手机号、身份证号、真实姓名、登录IP等敏感信息使用SM4加密存储
--    - 采用CBC模式,密钥由后端密钥管理系统统一管理
--    - 数据在存储前加密,读取后解密
--
-- 3. SM2加密(用于通信):
--    - 前后端通信中的敏感数据使用SM2非对称加密
--    - 后端提供公钥,前端使用公钥加密数据
--    - 后端使用私钥解密数据
--
-- 5. Paillier同态加密(用于AI计算):
--    - 症状特征数据使用Paillier加密
--    - 支持在加密数据上进行AI风险评分计算
--    - 保护患者症状隐私
--
-- 6. SM9属性基加密(用于访问控制):
--    - 病历访问策略使用SM9属性基加密
--    - 基于医生属性(科室、职称)进行细粒度访问控制
--    - 只有属性匹配的医生才能解密病历
--
-- 7. HTTPS + 国密SSL:
--    - 使用国密SSL证书(SM2算法)建立安全连接
--    - 确保传输层安全
--
-- 8. 完整性验证:
--    - 所有关键数据使用SM3计算哈希值
--    - 每次读取数据时验证完整性
--    - 防止数据被篡改
-- =============================================
-- 8. 完整性验证:
--    - 所有关键数据使用SM3计算哈希值
--    - 每次读取数据时验证完整性
--    - 防止数据被篡改
-- =============================================
