-- =====================================================
-- 药品数据库初始化脚本
-- 版本: 1.0
-- 创建日期: 2025-12-17
-- =====================================================

USE SM;

-- 药品基础信息表
CREATE TABLE IF NOT EXISTS SM_medicine (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '药品ID',
    medicine_code VARCHAR(50) NOT NULL UNIQUE COMMENT '药品编码',
    medicine_name VARCHAR(200) NOT NULL COMMENT '药品名称',
    common_name VARCHAR(200) COMMENT '通用名',
    english_name VARCHAR(200) COMMENT '英文名',
    medicine_type VARCHAR(50) NOT NULL COMMENT '药品类型(西药/中成药/中药饮片)',
    category VARCHAR(50) COMMENT '分类(抗生素/降压药/降糖药等)',
    specification VARCHAR(100) COMMENT '规格(如:0.5g*24片)',
    dosage_form VARCHAR(50) COMMENT '剂型(片剂/胶囊/注射液等)',
    manufacturer VARCHAR(200) COMMENT '生产厂家',
    approval_number VARCHAR(100) COMMENT '批准文号',
    price DECIMAL(10,2) DEFAULT 0.00 COMMENT '单价(元)',
    unit VARCHAR(20) DEFAULT '盒' COMMENT '单位',
    
    -- 用药信息
    indications TEXT COMMENT '适应症',
    contraindications TEXT COMMENT '禁忌症',
    usage_dosage TEXT COMMENT '用法用量',
    adverse_reactions TEXT COMMENT '不良反应',
    precautions TEXT COMMENT '注意事项',
    interactions TEXT COMMENT '药物相互作用',
    
    -- 管理信息
    prescription_type TINYINT DEFAULT 0 COMMENT '处方类型(0:普通,1:处方药,2:管制药品)',
    is_otc TINYINT DEFAULT 0 COMMENT '是否OTC(0:否,1:是)',
    is_medical_insurance TINYINT DEFAULT 1 COMMENT '是否医保(0:否,1:是)',
    insurance_code VARCHAR(50) COMMENT '医保编码',
    storage_conditions VARCHAR(200) COMMENT '贮藏条件',
    validity_period VARCHAR(50) COMMENT '有效期',
    
    -- 库存信息
    stock_quantity INT DEFAULT 0 COMMENT '库存数量',
    safety_stock INT DEFAULT 10 COMMENT '安全库存',
    status TINYINT DEFAULT 1 COMMENT '状态(0:停用,1:在用,2:缺货)',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_medicine_code (medicine_code),
    INDEX idx_medicine_name (medicine_name),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_prescription_type (prescription_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='药品基础信息表';

-- 处方表
CREATE TABLE IF NOT EXISTS SM_prescription (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '处方ID',
    prescription_no VARCHAR(50) NOT NULL UNIQUE COMMENT '处方编号',
    consultation_id BIGINT NOT NULL COMMENT '问诊ID',
    patient_id BIGINT NOT NULL COMMENT '患者ID',
    doctor_id BIGINT NOT NULL COMMENT '医生ID',
    diagnosis TEXT COMMENT '诊断(SM4加密)',
    prescription_type TINYINT DEFAULT 1 COMMENT '处方类型(1:普通,2:急诊,3:儿科,4:麻精)',
    total_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT '处方总金额',
    status TINYINT DEFAULT 0 COMMENT '状态(0:待审核,1:已审核,2:已配药,3:已取药,4:已作废)',
    
    -- 审核信息
    reviewer_id BIGINT COMMENT '审核药师ID',
    review_time TIMESTAMP NULL COMMENT '审核时间',
    review_notes TEXT COMMENT '审核备注',
    
    -- 签名和哈希
    doctor_signature VARCHAR(500) COMMENT '医生数字签名',
    data_hash VARCHAR(128) COMMENT '数据完整性哈希(SM3)',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '开具时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (consultation_id) REFERENCES SM_consultation(id) ON DELETE CASCADE,
    FOREIGN KEY (patient_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    FOREIGN KEY (doctor_id) REFERENCES SM_user(id) ON DELETE CASCADE,
    
    INDEX idx_prescription_no (prescription_no),
    INDEX idx_consultation_id (consultation_id),
    INDEX idx_patient_id (patient_id),
    INDEX idx_doctor_id (doctor_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='处方表';

-- 处方明细表
CREATE TABLE IF NOT EXISTS SM_prescription_detail (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '明细ID',
    prescription_id BIGINT NOT NULL COMMENT '处方ID',
    medicine_id BIGINT NOT NULL COMMENT '药品ID',
    medicine_name VARCHAR(200) NOT NULL COMMENT '药品名称(冗余)',
    specification VARCHAR(100) COMMENT '规格',
    quantity INT NOT NULL COMMENT '数量',
    unit VARCHAR(20) DEFAULT '盒' COMMENT '单位',
    unit_price DECIMAL(10,2) DEFAULT 0.00 COMMENT '单价',
    total_price DECIMAL(10,2) DEFAULT 0.00 COMMENT '小计',
    
    -- 用药说明
    `usage` VARCHAR(100) COMMENT '用法(口服/外用/注射等)',
    `frequency` VARCHAR(100) COMMENT '频次(每日3次/每日2次等)',
    `dosage` VARCHAR(100) COMMENT '单次剂量(1片/2粒等)',
    `duration` VARCHAR(50) COMMENT '疗程(7天/14天等)',
    `notes` TEXT COMMENT '特殊说明',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY (prescription_id) REFERENCES SM_prescription(id) ON DELETE CASCADE,
    FOREIGN KEY (medicine_id) REFERENCES SM_medicine(id) ON DELETE RESTRICT,
    
    INDEX idx_prescription_id (prescription_id),
    INDEX idx_medicine_id (medicine_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='处方明细表';

-- =====================================================
-- 插入常用药品数据
-- =====================================================

-- 抗生素类
INSERT INTO SM_medicine (medicine_code, medicine_name, common_name, medicine_type, category, specification, dosage_form, manufacturer, price, indications, usage_dosage, prescription_type, is_otc, status) VALUES
('M001', '阿莫西林胶囊', '阿莫西林', '西药', '抗生素', '0.25g*24粒', '胶囊', '华北制药', 12.50, '适用于敏感菌所致的呼吸道感染、泌尿道感染等', '口服，成人一次0.5g，每6-8小时1次', 1, 0, 1),
('M002', '头孢克肟分散片', '头孢克肟', '西药', '抗生素', '100mg*6片', '分散片', '扬子江药业', 28.00, '用于治疗敏感菌引起的呼吸道、尿路感染', '口服，成人一次100mg，每日2次', 1, 0, 1),
('M003', '阿奇霉素片', '阿奇霉素', '西药', '抗生素', '0.25g*6片', '片剂', '辉瑞制药', 35.00, '适用于敏感细菌所致的呼吸道、皮肤软组织感染', '口服，成人首日0.5g，第2-5日每日0.25g', 1, 0, 1);

-- 降压药类
INSERT INTO SM_medicine (medicine_code, medicine_name, common_name, medicine_type, category, specification, dosage_form, manufacturer, price, indications, usage_dosage, prescription_type, is_otc, status) VALUES
('M004', '硝苯地平缓释片', '硝苯地平', '西药', '降压药', '20mg*30片', '缓释片', '拜耳医药', 45.00, '用于治疗高血压、心绞痛', '口服，一次20mg，每日2次', 1, 0, 1),
('M005', '氨氯地平片', '氨氯地平', '西药', '降压药', '5mg*14片', '片剂', '辉瑞制药', 38.00, '用于治疗高血压', '口服，初始剂量5mg，每日1次', 1, 0, 1),
('M006', '缬沙坦胶囊', '缬沙坦', '西药', '降压药', '80mg*7粒', '胶囊', '北京诺华', 52.00, '用于治疗原发性高血压', '口服，初始剂量80mg，每日1次', 1, 0, 1);

-- 降糖药类
INSERT INTO SM_medicine (medicine_code, medicine_name, common_name, medicine_type, category, specification, dosage_form, manufacturer, price, indications, usage_dosage, prescription_type, is_otc, status) VALUES
('M007', '二甲双胍片', '二甲双胍', '西药', '降糖药', '0.5g*48片', '片剂', '中美上海施贵宝', 18.00, '用于治疗2型糖尿病', '口服，初始0.5g，每日2-3次，餐中或餐后服', 1, 0, 1),
('M008', '格列美脲片', '格列美脲', '西药', '降糖药', '2mg*15片', '片剂', '赛诺菲', 35.00, '用于治疗2型糖尿病', '口服，初始1mg，每日1次，早餐前或早餐中服用', 1, 0, 1),
('M009', '阿卡波糖片', '阿卡波糖', '西药', '降糖药', '50mg*30片', '片剂', '拜耳医药', 42.00, '配合饮食控制治疗2型糖尿病', '口服，初始50mg，每日3次，餐前即刻服用', 1, 0, 1);

-- 感冒药类
INSERT INTO SM_medicine (medicine_code, medicine_name, common_name, medicine_type, category, specification, dosage_form, manufacturer, price, indications, usage_dosage, prescription_type, is_otc, status) VALUES
('M010', '复方氨酚烷胺片', '复方氨酚烷胺', '西药', '感冒药', '12片', '片剂', '修正药业', 8.50, '适用于缓解感冒引起的发热、头痛等症状', '口服，成人一次1片，每日2次', 0, 1, 1),
('M011', '板蓝根颗粒', '板蓝根', '中成药', '感冒药', '10g*20袋', '颗粒', '白云山', 15.00, '清热解毒，用于感冒发热、咽喉肿痛', '开水冲服，一次5-10g，每日3-4次', 0, 1, 1),
('M012', '感冒清热颗粒', '感冒清热', '中成药', '感冒药', '12g*10袋', '颗粒', '同仁堂', 18.00, '疏风散寒，解表清热，用于风寒感冒', '开水冲服，一次1袋，每日2次', 0, 1, 1);

-- 止咳化痰类
INSERT INTO SM_medicine (medicine_code, medicine_name, common_name, medicine_type, category, specification, dosage_form, manufacturer, price, indications, usage_dosage, prescription_type, is_otc, status) VALUES
('M013', '氨溴索口服液', '氨溴索', '西药', '止咳化痰', '100ml', '口服液', '勃林格殷格翰', 28.00, '用于痰液粘稠不易咳出者', '口服，成人一次10ml，每日3次', 0, 1, 1),
('M014', '复方甘草片', '复方甘草', '中成药', '止咳化痰', '100片', '片剂', '太极集团', 6.50, '镇咳祛痰，用于咳嗽痰多', '口服，一次3-4片，每日3次', 0, 1, 1),
('M015', '川贝枇杷糖浆', '川贝枇杷', '中成药', '止咳化痰', '150ml', '糖浆', '潘高寿', 22.00, '清热宣肺，化痰止咳', '口服，一次15ml，每日3次', 0, 1, 1);

-- 消化系统药物
INSERT INTO SM_medicine (medicine_code, medicine_name, common_name, medicine_type, category, specification, dosage_form, manufacturer, price, indications, usage_dosage, prescription_type, is_otc, status) VALUES
('M016', '奥美拉唑肠溶胶囊', '奥美拉唑', '西药', '消化系统', '20mg*14粒', '肠溶胶囊', '阿斯利康', 32.00, '用于胃溃疡、十二指肠溃疡', '口服，一次20mg，每日1-2次', 1, 0, 1),
('M017', '多潘立酮片', '多潘立酮', '西药', '消化系统', '10mg*30片', '片剂', '西安杨森', 25.00, '用于消化不良、腹胀等', '口服，一次10mg，每日3次，饭前15-30分钟服用', 0, 1, 1),
('M018', '健胃消食片', '健胃消食', '中成药', '消化系统', '0.8g*36片', '片剂', '江中药业', 12.00, '健胃消食，用于脾胃虚弱所致的食积', '口服，一次3片，每日3次', 0, 1, 1);

-- 维生素类
INSERT INTO SM_medicine (medicine_code, medicine_name, common_name, medicine_type, category, specification, dosage_form, manufacturer, price, indications, usage_dosage, prescription_type, is_otc, status) VALUES
('M019', '复合维生素B片', '维生素B', '西药', '维生素', '100片', '片剂', '华润双鹤', 15.00, '用于预防和治疗B族维生素缺乏', '口服，一次1-3片，每日3次', 0, 1, 1),
('M020', '维生素C片', '维生素C', '西药', '维生素', '100mg*100片', '片剂', '华北制药', 8.00, '用于预防坏血病，增强抵抗力', '口服，一次100-200mg，每日3次', 0, 1, 1);

-- =====================================================
-- 初始化完成
-- =====================================================
-- 已插入20种常用药品
-- 包含：抗生素、降压药、降糖药、感冒药、止咳化痰药、消化系统药、维生素等
-- =====================================================
