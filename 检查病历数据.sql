-- 检查病历表数据
SELECT 
    id,
    record_no,
    patient_id,
    consultation_id,
    record_type,
    doctor_id,
    created_at
FROM SM_medical_record
ORDER BY created_at DESC
LIMIT 10;

-- 检查病历总数
SELECT COUNT(*) as total FROM SM_medical_record;

-- 检查问诊完成记录（应该会生成病历）
SELECT 
    id,
    consultation_no,
    patient_id,
    doctor_id,
    status,
    completed_at,
    created_at
FROM SM_consultation
WHERE status = 2
ORDER BY created_at DESC
LIMIT 10;

-- 检查是否有已完成但未生成病历的问诊
SELECT 
    c.id,
    c.consultation_no,
    c.patient_id,
    c.doctor_id,
    c.status,
    c.completed_at,
    mr.id as record_id
FROM SM_consultation c
LEFT JOIN SM_medical_record mr ON c.id = mr.consultation_id
WHERE c.status = 2
ORDER BY c.created_at DESC;
