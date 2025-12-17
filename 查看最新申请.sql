-- 查看最新的医生申请记录
SELECT 
    id,
    user_id,
    application_no,
    doctor_title,
    doctor_dept,
    cert_number,
    status,
    created_at
FROM SM_doctor_application 
ORDER BY created_at DESC 
LIMIT 5;

-- 查看待审核的申请
SELECT 
    a.id,
    a.application_no,
    u.username,
    a.doctor_title,
    a.doctor_dept,
    a.status,
    a.created_at
FROM SM_doctor_application a
LEFT JOIN SM_user u ON a.user_id = u.id
WHERE a.status = 0
ORDER BY a.created_at DESC;
