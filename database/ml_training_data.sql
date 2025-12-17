-- =============================================
-- 机器学习训练数据表
-- 用于收集真实诊断数据，训练AI模型
-- =============================================

-- 1. AI训练数据表
CREATE TABLE IF NOT EXISTS SM_ai_training_data (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '训练数据ID',
    consultation_id BIGINT NOT NULL COMMENT '关联问诊ID',
    patient_id BIGINT NOT NULL COMMENT '患者ID(去标识化)',
    
    -- 输入特征（脱敏后）
    age INT COMMENT '年龄',
    gender TINYINT COMMENT '性别(0:女,1:男)',
    systolic_bp INT COMMENT '收缩压(mmHg)',
    diastolic_bp INT COMMENT '舒张压(mmHg)',
    heart_rate INT COMMENT '心率(次/分)',
    temperature DECIMAL(4,1) COMMENT '体温(℃)',
    blood_sugar DECIMAL(5,2) COMMENT '血糖(mmol/L)',
    bmi DECIMAL(5,2) COMMENT 'BMI指数',
    
    -- 症状关键词(JSON格式)
    symptom_keywords TEXT COMMENT '症状关键词列表["头痛","发热"...]',
    symptom_duration VARCHAR(50) COMMENT '症状持续时间',
    symptom_severity TINYINT COMMENT '症状严重程度(1-10)',
    
    -- 既往史特征
    has_hypertension TINYINT DEFAULT 0 COMMENT '是否有高血压史',
    has_diabetes TINYINT DEFAULT 0 COMMENT '是否有糖尿病史',
    has_heart_disease TINYINT DEFAULT 0 COMMENT '是否有心脏病史',
    smoking_status TINYINT COMMENT '吸烟状态(0:不吸,1:曾吸,2:现吸)',
    drinking_status TINYINT COMMENT '饮酒状态(0:不饮,1:偶尔,2:经常)',
    
    -- AI预测结果
    ai_risk_score INT COMMENT 'AI预测风险分数(0-100)',
    ai_predicted_disease VARCHAR(200) COMMENT 'AI预测疾病',
    ai_recommended_dept VARCHAR(100) COMMENT 'AI推荐科室',
    
    -- 医生确诊结果（标签）
    doctor_diagnosis TEXT COMMENT '医生确诊疾病(标准化后)',
    doctor_dept VARCHAR(100) COMMENT '实际就诊科室',
    diagnosis_icd10 VARCHAR(20) COMMENT 'ICD-10疾病编码',
    treatment_result TINYINT COMMENT '治疗结果(1:治愈,2:好转,3:无效,4:恶化)',
    
    -- 数据质量标记
    data_quality TINYINT DEFAULT 1 COMMENT '数据质量(1:高,2:中,3:低)',
    is_verified TINYINT DEFAULT 0 COMMENT '是否医生复核(0:否,1:是)',
    verified_by BIGINT COMMENT '复核医生ID',
    verified_at TIMESTAMP NULL COMMENT '复核时间',
    
    -- 模型训练标记
    is_used_for_training TINYINT DEFAULT 0 COMMENT '是否用于训练(0:否,1:是)',
    training_set_type TINYINT COMMENT '数据集类型(1:训练集,2:验证集,3:测试集)',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    FOREIGN KEY (consultation_id) REFERENCES SM_consultation(id) ON DELETE CASCADE,
    INDEX idx_consultation_id (consultation_id),
    INDEX idx_patient_id (patient_id),
    INDEX idx_doctor_diagnosis (doctor_diagnosis(100)),
    INDEX idx_is_verified (is_verified),
    INDEX idx_is_used_for_training (is_used_for_training),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI训练数据表';

-- 2. 模型版本管理表
CREATE TABLE IF NOT EXISTS SM_ai_model_version (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '模型版本ID',
    model_name VARCHAR(100) NOT NULL COMMENT '模型名称',
    version VARCHAR(50) NOT NULL COMMENT '版本号(如v1.0.0)',
    model_type VARCHAR(50) NOT NULL COMMENT '模型类型(逻辑回归/随机森林/神经网络)',
    
    -- 训练参数
    training_samples INT COMMENT '训练样本数',
    validation_samples INT COMMENT '验证样本数',
    test_samples INT COMMENT '测试样本数',
    feature_count INT COMMENT '特征数量',
    training_duration INT COMMENT '训练时长(秒)',
    
    -- 模型性能指标
    accuracy DECIMAL(5,4) COMMENT '准确率',
    precision_score DECIMAL(5,4) COMMENT '精确率',
    recall_score DECIMAL(5,4) COMMENT '召回率',
    f1_score DECIMAL(5,4) COMMENT 'F1分数',
    auc_score DECIMAL(5,4) COMMENT 'AUC值',
    
    -- 混淆矩阵(JSON)
    confusion_matrix TEXT COMMENT '混淆矩阵',
    
    -- 模型文件路径
    model_file_path VARCHAR(500) COMMENT '模型文件路径',
    model_file_size BIGINT COMMENT '模型文件大小(字节)',
    model_hash VARCHAR(128) COMMENT '模型文件SM3哈希',
    
    -- 部署信息
    is_deployed TINYINT DEFAULT 0 COMMENT '是否部署(0:否,1:是)',
    deployed_at TIMESTAMP NULL COMMENT '部署时间',
    deployed_by BIGINT COMMENT '部署人ID',
    
    -- 备注
    description TEXT COMMENT '版本说明',
    training_config TEXT COMMENT '训练配置(JSON)',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    UNIQUE KEY uk_model_version (model_name, version),
    INDEX idx_model_name (model_name),
    INDEX idx_is_deployed (is_deployed),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI模型版本管理表';

-- 3. 模型预测日志表
CREATE TABLE IF NOT EXISTS SM_ai_prediction_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '预测日志ID',
    consultation_id BIGINT NOT NULL COMMENT '问诊ID',
    model_version_id BIGINT COMMENT '模型版本ID',
    
    -- 输入特征
    input_features TEXT COMMENT '输入特征(JSON)',
    
    -- 预测结果
    predicted_disease VARCHAR(200) COMMENT '预测疾病',
    predicted_dept VARCHAR(100) COMMENT '预测科室',
    risk_score INT COMMENT '风险评分',
    confidence_score DECIMAL(5,4) COMMENT '置信度(0-1)',
    
    -- 性能指标
    prediction_time INT COMMENT '预测耗时(毫秒)',
    
    -- 实际结果（后续回填）
    actual_disease VARCHAR(200) COMMENT '实际诊断',
    is_correct TINYINT COMMENT '预测是否正确(0:否,1:是)',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY (consultation_id) REFERENCES SM_consultation(id) ON DELETE CASCADE,
    FOREIGN KEY (model_version_id) REFERENCES SM_ai_model_version(id),
    INDEX idx_consultation_id (consultation_id),
    INDEX idx_model_version_id (model_version_id),
    INDEX idx_is_correct (is_correct),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI预测日志表';

-- 4. 疾病标准化映射表
CREATE TABLE IF NOT EXISTS SM_disease_mapping (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '映射ID',
    disease_name VARCHAR(200) NOT NULL COMMENT '疾病名称',
    disease_alias VARCHAR(500) COMMENT '疾病别名(多个用逗号分隔)',
    icd10_code VARCHAR(20) COMMENT 'ICD-10编码',
    disease_category VARCHAR(100) COMMENT '疾病分类',
    recommended_dept VARCHAR(100) COMMENT '推荐科室',
    severity_level TINYINT COMMENT '严重程度(1:轻,2:中,3:重,4:危重)',
    
    description TEXT COMMENT '疾病描述',
    typical_symptoms TEXT COMMENT '典型症状(JSON)',
    
    is_enabled TINYINT DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    UNIQUE KEY uk_disease_name (disease_name),
    INDEX idx_icd10_code (icd10_code),
    INDEX idx_disease_category (disease_category),
    INDEX idx_is_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='疾病标准化映射表';

-- 插入常见疾病映射数据
INSERT INTO SM_disease_mapping (disease_name, disease_alias, icd10_code, disease_category, recommended_dept, severity_level, typical_symptoms) VALUES
('高血压1级', '轻度高血压', 'I10', '心血管疾病', '心内科', 2, '["头晕","头痛","心悸"]'),
('高血压2级', '中度高血压', 'I10', '心血管疾病', '心血管内科', 3, '["头晕","头痛","心悸","胸闷"]'),
('高血压3级', '重度高血压', 'I10', '心血管疾病', '心血管内科', 4, '["头晕","头痛","视物模糊","胸痛"]'),
('低血压', '血压偏低', 'I95', '心血管疾病', '心内科', 2, '["头晕","乏力","眼前发黑"]'),
('心动过速', '心率过快', 'R00.0', '心血管症状', '心内科', 2, '["心慌","胸闷","气短"]'),
('心动过缓', '心率过慢', 'R00.1', '心血管症状', '心内科', 2, '["头晕","乏力","黑朦"]'),
('糖尿病', '高血糖', 'E11', '内分泌疾病', '内分泌科', 3, '["多饮","多尿","多食","消瘦"]'),
('低血糖', '血糖偏低', 'E16.2', '内分泌症状', '内分泌科', 3, '["心慌","出汗","饥饿感","头晕"]'),
('发热', '体温升高', 'R50', '全身症状', '呼吸内科', 2, '["发热","乏力","头痛"]'),
('高热', '高烧', 'R50.9', '全身症状', '感染科', 3, '["高热","寒战","头痛","乏力"]'),
('感冒', '上呼吸道感染', 'J06.9', '呼吸系统疾病', '呼吸内科', 1, '["咳嗽","流涕","咽痛","发热"]'),
('急性胃炎', '胃炎', 'K29.1', '消化系统疾病', '消化内科', 2, '["腹痛","恶心","呕吐"]'),
('偏头痛', '头痛', 'G43', '神经系统疾病', '神经内科', 2, '["头痛","恶心","畏光"]');

-- 创建触发器：问诊完成后自动收集训练数据
DELIMITER //
CREATE TRIGGER after_consultation_complete
AFTER UPDATE ON SM_consultation
FOR EACH ROW
BEGIN
    -- 当问诊状态变为已完成(2)且有医生诊断时，自动创建训练数据
    IF NEW.status = 2 AND OLD.status != 2 AND NEW.doctor_diagnosis IS NOT NULL AND NEW.doctor_diagnosis != '' THEN
        INSERT INTO SM_ai_training_data (
            consultation_id,
            patient_id,
            ai_risk_score,
            ai_predicted_disease,
            ai_recommended_dept,
            doctor_diagnosis,
            data_quality,
            is_verified,
            verified_by,
            verified_at
        ) VALUES (
            NEW.id,
            NEW.patient_id,
            NEW.ai_risk_score,
            NEW.ai_diagnosis,
            NULL, -- 需要从ai_suggestions中提取
            NEW.doctor_diagnosis,
            2, -- 默认中等质量
            1, -- 已医生确认
            NEW.doctor_id,
            NEW.completed_at
        );
    END IF;
END//
DELIMITER ;

-- 创建视图：高质量训练数据集
CREATE OR REPLACE VIEW v_ai_training_dataset AS
SELECT 
    t.id,
    t.age,
    t.gender,
    t.systolic_bp,
    t.diastolic_bp,
    t.heart_rate,
    t.temperature,
    t.blood_sugar,
    t.bmi,
    t.symptom_keywords,
    t.has_hypertension,
    t.has_diabetes,
    t.has_heart_disease,
    t.smoking_status,
    t.drinking_status,
    t.doctor_diagnosis AS label,
    t.diagnosis_icd10,
    m.disease_category,
    m.severity_level,
    t.ai_risk_score AS ai_prediction,
    CASE 
        WHEN t.doctor_diagnosis LIKE CONCAT('%', t.ai_predicted_disease, '%') THEN 1
        ELSE 0
    END AS ai_correct,
    t.training_set_type,
    t.created_at
FROM SM_ai_training_data t
LEFT JOIN SM_disease_mapping m ON t.diagnosis_icd10 = m.icd10_code
WHERE t.is_verified = 1 
  AND t.data_quality IN (1, 2) -- 高质量或中等质量
  AND t.doctor_diagnosis IS NOT NULL
  AND t.doctor_diagnosis != '';

-- 创建数据统计视图
CREATE OR REPLACE VIEW v_ai_training_stats AS
SELECT 
    COUNT(*) AS total_samples,
    SUM(CASE WHEN is_verified = 1 THEN 1 ELSE 0 END) AS verified_samples,
    SUM(CASE WHEN is_used_for_training = 1 THEN 1 ELSE 0 END) AS training_samples,
    SUM(CASE WHEN training_set_type = 1 THEN 1 ELSE 0 END) AS train_set,
    SUM(CASE WHEN training_set_type = 2 THEN 1 ELSE 0 END) AS val_set,
    SUM(CASE WHEN training_set_type = 3 THEN 1 ELSE 0 END) AS test_set,
    COUNT(DISTINCT doctor_diagnosis) AS unique_diseases,
    AVG(data_quality) AS avg_quality,
    MIN(created_at) AS earliest_date,
    MAX(created_at) AS latest_date
FROM SM_ai_training_data;
