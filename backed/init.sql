-- =====================================================
-- 基于国密加密的网上看诊系统数据库
-- 版本: 1.0
-- 创建日期: 2024-12-04
-- =====================================================

-- 创建数据库
CREATE DATABASE IF NOT EXISTS SM DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE SM;

-- =====================================================
-- 核心业务表
-- =====================================================

-- 用户表
-- 说明：存储所有用户信息（患者、医生、管理员）
-- 加密字段：email, phone, real_name, id_card, last_login_ip 使用SM4加密
-- 密码加密：双重SM3哈希（前端SM3 + 后端加盐SM3）
CREATE TABLE IF NOT EXISTS SM_user (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名(4-20字符)',
    password VARCHAR(128) NOT NULL COMMENT '密码(SM3双重哈希)',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱(SM4加密)',
    phone VARCHAR(255) COMMENT '手机号(SM4加密)',
    real_name VARCHAR(255) COMMENT '真实姓名(SM4加密)',
    id_card VARCHAR(255) COMMENT '身份证号(SM4加密)',
    identify ENUM('patient', 'admin', 'doctor') NOT NULL DEFAULT 'patient' COMMENT '身份标识',
    avatar VARCHAR(500) DEFAULT '/static/avatar/default.png' COMMENT '头像URL',
    gender TINYINT DEFAULT 0 COMMENT '性别(0:未知,1:男,2:女)',
    birth_date DATE COMMENT '出生日期',
    age INT COMMENT '年龄(冗余字段，便于统计)',
    status TINYINT DEFAULT 0 COMMENT '账号状态(0:正常,1:禁用,2:待审核,3:已注销)',
    
    -- 医生相关字段
    doctor_cert VARCHAR(500) COMMENT '医生资格证书URL',
    doctor_title VARCHAR(50) COMMENT '医生职称(如:主治医师、副主任医师)',
    doctor_dept VARCHAR(50) COMMENT '所属科室',
    specialty TEXT COMMENT '专业特长(多个特长用逗号分隔)',
    doctor_intro TEXT COMMENT '医生简介',
    cert_number VARCHAR(100) COMMENT '医师资格证号',
    practice_cert VARCHAR(100) COMMENT '医师执业证号',
    hospital_name VARCHAR(100) COMMENT '执业医院名称',
    consultation_price DECIMAL(10,2) DEFAULT 0.00 COMMENT '问诊价格(元)',
    total_consultations INT DEFAULT 0 COMMENT '累计问诊次数',
    rating DECIMAL(3,2) DEFAULT 5.00 COMMENT '医生评分(1-5)',
    
    -- 时间戳字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    last_login_time TIMESTAMP NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(255) COMMENT '最后登录IP(SM4加密)',
    deleted_at TIMESTAMP NULL COMMENT '删除时间(软删除)',
    
    -- 索引
    INDEX idx_identify (identify),
    INDEX idx_status (status),
    INDEX idx_doctor_dept (doctor_dept),
    INDEX idx_username (username),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 医生申请记录表
-- 说明：记录用户申请成为医生的审核流程
-- 审核通过后会更新SM_user表中的医生相关字段
CREATE TABLE IF NOT EXISTS SM_doctor_application (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '申请ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    application_no VARCHAR(50) NOT NULL UNIQUE COMMENT '申请编号',
    real_name VARCHAR(255) NOT NULL COMMENT '真实姓名(SM4加密)',
    id_card VARCHAR(255) NOT NULL COMMENT '身份证号(SM4加密)',
    phone VARCHAR(255) NOT NULL COMMENT '手机号(SM4加密)',
    email VARCHAR(255) NOT NULL COMMENT '邮箱(SM4加密)',
    doctor_cert VARCHAR(500) NOT NULL COMMENT '医生资格证书URL',
    doctor_title VARCHAR(50) NOT NULL COMMENT '医生职称',
    doctor_dept VARCHAR(50) NOT NULL COMMENT '所属科室',
    specialty TEXT COMMENT '专业特长',
    doctor_intro TEXT COMMENT '医生简介',
    cert_number VARCHAR(100) NOT NULL COMMENT '医师资格证号',
    practice_cert VARCHAR(100) COMMENT '医师执业证号',
    hospital_name VARCHAR(100) COMMENT '执业医院名称',
    status TINYINT DEFAULT 0 COMMENT '审核状态(0:待审核,1:已通过,2:已拒绝)',
    reject_reason TEXT COMMENT '拒绝原因',
    reviewer_id BIGINT COMMENT '审核人ID',
    review_time TIMESTAMP NULL COMMENT '审核时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewer_id) REFERENCES SM_user(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_application_no (application_no),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='医生申请记录表';

-- 用户登录日志表
-- 说明：记录所有用户的登录行为，用于安全审计和异常检测
CREATE TABLE IF NOT EXISTS SM_login_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
    user_id BIGINT COMMENT '用户ID',
    username VARCHAR(50) COMMENT '用户名',
    login_ip VARCHAR(255) NOT NULL COMMENT '登录IP(SM4加密)',
    login_location VARCHAR(100) COMMENT '登录地点',
    browser VARCHAR(50) COMMENT '浏览器类型',
    os VARCHAR(50) COMMENT '操作系统',
    device_type VARCHAR(20) COMMENT '设备类型(PC/Mobile/Tablet)',
    status TINYINT NOT NULL COMMENT '登录状态(0:失败,1:成功)',
    msg VARCHAR(255) COMMENT '提示消息',
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    INDEX idx_user_id (user_id),
    INDEX idx_login_time (login_time),
    INDEX idx_status (status),
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户登录日志表';

-- =====================================================
-- 问诊业务表
-- =====================================================

-- 问诊记录表
-- 说明：记录患者和医生的在线问诊信息，支持AI辅助诊断
CREATE TABLE IF NOT EXISTS SM_consultation (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '问诊ID',
    patient_id BIGINT NOT NULL COMMENT '患者ID',
    doctor_id BIGINT COMMENT '医生ID',
    consultation_no VARCHAR(50) NOT NULL UNIQUE COMMENT '问诊编号(格式:C+时间戳+随机数)',
    chief_complaint TEXT COMMENT '主诉(SM4加密)',
    symptoms_desc TEXT COMMENT '症状描述(SM4加密)',
    symptoms_duration VARCHAR(50) COMMENT '症状持续时间',
    symptoms_encrypted TEXT COMMENT '症状特征(加密用于AI分析)',
    images VARCHAR(1000) COMMENT '相关图片URLs(多个用逗号分隔)',
    ai_risk_score INT COMMENT 'AI风险评分(0-100)',
    ai_diagnosis TEXT COMMENT 'AI初步诊断建议',
    ai_suggestions TEXT COMMENT 'AI建议',
    doctor_diagnosis TEXT COMMENT '医生诊断(SM4加密)',
    prescription TEXT COMMENT '处方信息(SM4加密)',
    doctor_advice TEXT COMMENT '医生建议(SM4加密)',
    need_ai TINYINT DEFAULT 1 COMMENT '是否需要AI诊断(0:否,1:是)',
    status TINYINT DEFAULT 0 COMMENT '状态(0:待接诊,1:问诊中,2:已完成,3:已取消,4:已超时)',
    price DECIMAL(10,2) DEFAULT 0.00 COMMENT '问诊费用',
    rating TINYINT COMMENT '患者评分(1-5星)',
    rating_comment TEXT COMMENT '患者评价',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    accepted_at TIMESTAMP NULL COMMENT '医生接诊时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    cancelled_at TIMESTAMP NULL COMMENT '取消时间',
    FOREIGN KEY (patient_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    FOREIGN KEY (doctor_id) REFERENCES SM_user(id) ON DELETE SET NULL,
    INDEX idx_patient_id (patient_id),
    INDEX idx_doctor_id (doctor_id),
    INDEX idx_status (status),
    INDEX idx_consultation_no (consultation_no),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='问诊记录表';

-- 问诊消息表
-- 说明：记录问诊过程中的实时聊天消息
-- 支持文本、图片、语音、处方等多种消息类型
CREATE TABLE IF NOT EXISTS SM_consultation_message (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '消息ID',
    consultation_id BIGINT NOT NULL COMMENT '问诊ID',
    sender_id BIGINT NOT NULL COMMENT '发送者ID',
    receiver_id BIGINT NOT NULL COMMENT '接收者ID',
    message_type TINYINT NOT NULL COMMENT '消息类型(1:文本,2:图片,3:语音,4:处方,5:系统消息)',
    content TEXT NOT NULL COMMENT '消息内容(SM4加密)',
    content_plain TEXT COMMENT '纯文本内容(用于搜索)',
    file_url VARCHAR(500) COMMENT '文件URL',
    file_size INT COMMENT '文件大小(字节)',
    duration INT COMMENT '语音时长(秒)',
    is_read TINYINT DEFAULT 0 COMMENT '是否已读(0:未读,1:已读)',
    read_at TIMESTAMP NULL COMMENT '阅读时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间(软删除)',
    FOREIGN KEY (consultation_id) REFERENCES SM_consultation(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    INDEX idx_consultation_id (consultation_id),
    INDEX idx_sender_id (sender_id),
    INDEX idx_receiver_id (receiver_id),
    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='问诊消息表';

-- =====================================================
-- 病历管理表
-- =====================================================

-- 电子病历表
-- 说明：存储患者的完整电子病历信息
-- 安全特性：使用SM3哈希值验证数据完整性
CREATE TABLE IF NOT EXISTS SM_medical_record (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '病历ID',
    record_no VARCHAR(50) NOT NULL UNIQUE COMMENT '病历编号(格式:MR+时间戳+随机数)',
    patient_id BIGINT NOT NULL COMMENT '患者ID',
    consultation_id BIGINT COMMENT '关联问诊ID',
    record_type TINYINT NOT NULL COMMENT '病历类型(1:门诊,2:在线问诊,3:住院)',
    chief_complaint TEXT COMMENT '主诉(SM4加密)',
    present_illness TEXT COMMENT '现病史(SM4加密)',
    past_history TEXT COMMENT '既往史(SM4加密)',
    allergy_history TEXT COMMENT '过敏史(SM4加密)',
    family_history TEXT COMMENT '家族史(SM4加密)',
    physical_examination TEXT COMMENT '体格检查(SM4加密)',
    auxiliary_examination TEXT COMMENT '辅助检查(SM4加密)',
    diagnosis TEXT COMMENT '诊断(SM4加密)',
    treatment_plan TEXT COMMENT '治疗方案(SM4加密)',
    prescription TEXT COMMENT '处方(SM4加密)',
    doctor_id BIGINT COMMENT '诊疗医生ID',
    ai_advice TEXT COMMENT 'AI建议',
    data_hash VARCHAR(128) COMMENT '数据完整性哈希(SM3)',
    signature VARCHAR(500) COMMENT '医生数字签名',
    is_locked TINYINT DEFAULT 0 COMMENT '是否锁定(0:否,1:是)',
    lock_time TIMESTAMP NULL COMMENT '锁定时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (patient_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    FOREIGN KEY (consultation_id) REFERENCES SM_consultation(id) ON DELETE SET NULL,
    FOREIGN KEY (doctor_id) REFERENCES SM_user(id) ON DELETE SET NULL,
    INDEX idx_patient_id (patient_id),
    INDEX idx_consultation_id (consultation_id),
    INDEX idx_record_no (record_no),
    INDEX idx_doctor_id (doctor_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='电子病历表';

-- 病历访问日志表
-- 说明：记录所有对电子病历的访问操作，确保数据安全和可追溯性
CREATE TABLE IF NOT EXISTS SM_record_access_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
    record_id BIGINT NOT NULL COMMENT '病历ID',
    accessor_id BIGINT NOT NULL COMMENT '访问者ID',
    accessor_name VARCHAR(50) COMMENT '访问者姓名',
    accessor_role VARCHAR(20) COMMENT '访问者角色',
    access_type TINYINT NOT NULL COMMENT '访问类型(1:查看,2:编辑,3:下载,4:打印,5:删除)',
    access_result TINYINT NOT NULL COMMENT '访问结果(0:失败,1:成功)',
    failure_reason VARCHAR(255) COMMENT '失败原因',
    ip_address VARCHAR(255) COMMENT '访问IP(SM4加密)',
    user_agent VARCHAR(500) COMMENT '用户代理',
    access_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '访问时间',
    FOREIGN KEY (record_id) REFERENCES SM_medical_record(id) ON DELETE CASCADE,
    FOREIGN KEY (accessor_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    INDEX idx_record_id (record_id),
    INDEX idx_accessor_id (accessor_id),
    INDEX idx_access_time (access_time),
    INDEX idx_access_result (access_result)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='病历访问日志表';

-- =====================================================
-- 系统管理表
-- =====================================================

-- 系统通知表
-- 说明：存储系统通知、消息推送等信息
CREATE TABLE IF NOT EXISTS SM_notification (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '通知ID',
    user_id BIGINT NOT NULL COMMENT '接收用户ID',
    notification_type VARCHAR(20) NOT NULL COMMENT '通知类型(system:系统,consultation:问诊,audit:审核,payment:支付)',
    title VARCHAR(100) NOT NULL COMMENT '通知标题',
    content TEXT COMMENT '通知内容',
    related_id BIGINT COMMENT '关联ID',
    related_type VARCHAR(50) COMMENT '关联类型(consultation/application/record)',
    priority TINYINT DEFAULT 0 COMMENT '优先级(0:普通,1:重要,2:紧急)',
    is_read TINYINT DEFAULT 0 COMMENT '是否已读(0:未读,1:已读)',
    is_popup TINYINT DEFAULT 0 COMMENT '是否弹窗提示(0:否,1:是)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    read_at TIMESTAMP NULL COMMENT '读取时间',
    expired_at TIMESTAMP NULL COMMENT '过期时间',
    FOREIGN KEY (user_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_is_read (is_read),
    INDEX idx_notification_type (notification_type),
    INDEX idx_priority (priority),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统通知表';

-- =====================================================
-- 插入初始数据
-- =====================================================

-- 插入管理员账号
-- 用户名: admin
-- 密码: Admin123!@# 
-- 密码经过双重SM3哈希：前端SM3 + 后端加盐(用户名)再SM3
INSERT INTO SM_user (
    username, 
    password, 
    email, 
    phone, 
    identify, 
    status,
    gender,
    birth_date,
    avatar,
    created_at
) VALUES (
    'admin',
    '137678542f3473058bac95a1c5bb30b9d708f5cef150717308430b5d4a944378',
    'admin@sm-medical.com',
    '18800000000',
    'admin',
    0,
    0,
    '1994-01-01',
    '/static/avatar/admin.png',
    NOW()
);

-- 插入测试患者账号
-- 用户名: patient01
-- 密码: Patient123!@#
-- 注意：实际密码需要通过注册接口生成正确的哈希值
INSERT INTO SM_user (
    username,
    password,
    email,
    phone,
    real_name,
    identify,
    status,
    gender,
    birth_date,
    avatar,
    created_at
) VALUES (
    'patient01',
    'patient_password_hash_placeholder',
    'patient01@example.com',
    '13800138000',
    '张三',
    'patient',
    0,
    1,
    '1996-03-15',
    '/static/avatar/default.png',
    NOW()
);

-- 插入测试医生账号
-- 用户名: doctor01
-- 密码: Doctor123!@#
INSERT INTO SM_user (
    username,
    password,
    email,
    phone,
    real_name,
    identify,
    status,
    gender,
    birth_date,
    avatar,
    doctor_title,
    doctor_dept,
    specialty,
    doctor_intro,
    doctor_cert,
    cert_number,
    practice_cert,
    hospital_name,
    consultation_price,
    total_consultations,
    rating,
    created_at
) VALUES (
    'doctor01',
    'doctor_password_hash_placeholder',
    'doctor01@sm-medical.com',
    '13900139000',
    '李医生',
    'doctor',
    0,
    1,
    '1986-05-20',
    '/static/avatar/doctor01.png',
    '主治医师',
    '内科',
    '高血压,糖尿病,冠心病',
    '从事内科临床工作15年，擅长心血管疾病、内分泌疾病的诊断和治疗。',
    '/static/cert/doctor01_cert.jpg',
    '110123456789012345',
    '110987654321098765',
    '北京市人民医院',
    50.00,
    0,
    5.00,
    NOW()
);

-- 插入第二位测试医生
INSERT INTO SM_user (
    username,
    password,
    email,
    phone,
    real_name,
    identify,
    status,
    gender,
    birth_date,
    avatar,
    doctor_title,
    doctor_dept,
    specialty,
    doctor_intro,
    doctor_cert,
    cert_number,
    practice_cert,
    hospital_name,
    consultation_price,
    total_consultations,
    rating,
    created_at
) VALUES (
    'doctor02',
    'doctor_password_hash_placeholder',
    'doctor02@sm-medical.com',
    '13900139001',
    '王医生',
    'doctor',
    0,
    2,
    '1982-08-10',
    '/static/avatar/doctor02.png',
    '副主任医师',
    '儿科',
    '儿童呼吸道疾病,儿童消化系统疾病',
    '从事儿科临床工作20年，擅长儿童常见病、多发病的诊治。',
    '/static/cert/doctor02_cert.jpg',
    '110123456789012346',
    '110987654321098766',
    '北京儿童医院',
    80.00,
    0,
    5.00,
    NOW()
);

-- =====================================================
-- 数据库初始化完成
-- =====================================================
-- 
-- 系统信息：
-- 数据库名: SM
-- 字符集: utf8mb4
-- 排序规则: utf8mb4_unicode_ci
-- 
-- 表统计 (共8张表)：
-- 1. SM_user - 用户表 (核心表)
-- 2. SM_doctor_application - 医生申请记录表
-- 3. SM_login_log - 用户登录日志表
-- 4. SM_consultation - 问诊记录表
-- 5. SM_consultation_message - 问诊消息表
-- 6. SM_medical_record - 电子病历表
-- 7. SM_record_access_log - 病历访问日志表
-- 8. SM_notification - 系统通知表
-- 
-- 初始用户 (共4个)：
-- 1. admin - 管理员账号 (密码: Admin123!@#) ✓
-- 2. patient01 - 测试患者账号 (密码需注册时生成)
-- 3. doctor01 - 测试医生账号 (密码需注册时生成)
-- 4. doctor02 - 测试医生账号 (密码需注册时生成)
-- 
-- 安全特性：
-- ✓ 敏感信息SM4加密 (邮箱、手机、姓名、身份证、IP)
-- ✓ 密码双重SM3哈希 (前端+后端加盐)
-- ✓ 病历数据完整性SM3校验
-- ✓ 外键约束确保数据一致性
-- ✓ 索引优化查询性能
-- ✓ 软删除支持数据恢复
-- 
-- 注意事项：
-- 1. 测试账号密码需要通过前端注册接口生成正确的哈希值
-- 2. 生产环境建议修改管理员密码
-- 3. 定期备份数据库
-- 4. 监控日志表大小，定期归档
-- =====================================================
