-- 检查问诊表数据
SELECT 
    id,
    consultation_no,
    patient_id,
    doctor_id,
    status,
    created_at,
    updated_at
FROM SM_consultation
ORDER BY created_at DESC
LIMIT 10;

-- 检查问诊总数
SELECT COUNT(*) as total FROM SM_consultation;

-- 检查各状态的问诊数量
SELECT 
    status,
    CASE status
        WHEN 0 THEN '待接诊'
        WHEN 1 THEN '问诊中'
        WHEN 2 THEN '已完成'
        ELSE '未知'
    END as status_name,
    COUNT(*) as count
FROM SM_consultation
GROUP BY status;

-- 查看最新的一条问诊记录的详细信息
SELECT * FROM SM_consultation ORDER BY created_at DESC LIMIT 1;
