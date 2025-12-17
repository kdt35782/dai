-- 插入测试病历数据
-- 注意：这里的加密数据需要使用实际的SM4加密结果
-- 为了测试，这里直接插入明文，实际应该通过API创建

-- 1. 先查看现有用户和问诊数据
SELECT id, username, identify FROM SM_user WHERE identify IN ('patient', 'doctor') LIMIT 5;

-- 2. 查看已完成的问诊
SELECT id, patient_id, doctor_id, consultation_no, status FROM SM_consultation WHERE status = 2 LIMIT 5;

-- 3. 如果有已完成但未生成病历的问诊，可以通过以下方式手动创建病历
-- （注意：实际环境中应该通过完成问诊的API自动创建）

-- 示例：为已完成的问诊补充病历（假设consultation_id=1, patient_id=3, doctor_id=2）
/*
INSERT INTO SM_medical_record (
    record_no,
    patient_id,
    consultation_id,
    record_type,
    chief_complaint,
    diagnosis,
    treatment_plan,
    doctor_id,
    ai_advice,
    data_hash,
    created_at
) VALUES (
    CONCAT('MR', UNIX_TIMESTAMP()),  -- 病历编号
    3,  -- 患者ID（需要替换为实际ID）
    1,  -- 问诊ID（需要替换为实际ID）
    2,  -- 病历类型：2=在线问诊
    '头痛头晕3天',  -- 主诉（实际应该SM4加密）
    '高血压1级',  -- 诊断（实际应该SM4加密）
    '建议休息，注意饮食，定期监测血压',  -- 治疗方案（实际应该SM4加密）
    2,  -- 医生ID（需要替换为实际ID）
    'AI初步分析：可能存在 高血压1级（轻度）',  -- AI建议
    'test_hash',  -- 数据哈希（实际应该SM3哈希）
    NOW()
);
*/

-- 4. 检查插入结果
SELECT * FROM SM_medical_record ORDER BY created_at DESC LIMIT 5;
