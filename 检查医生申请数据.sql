-- 检查医生申请数据
USE SM;

-- 1. 查看所有医生申请记录
SELECT 
    id AS '申请ID',
    user_id AS '用户ID',
    application_no AS '申请编号',
    doctor_title AS '职称',
    doctor_dept AS '科室',
    cert_number AS '证书号',
    status AS '状态',
    created_at AS '申请时间'
FROM SM_doctor_application 
ORDER BY created_at DESC;

-- 2. 查看待审核的医生申请
SELECT 
    id AS '申请ID',
    user_id AS '用户ID',
    application_no AS '申请编号',
    doctor_title AS '职称',
    doctor_dept AS '科室',
    status AS '状态(0=待审核)',
    created_at AS '申请时间'
FROM SM_doctor_application 
WHERE status = 0
ORDER BY created_at DESC;

-- 3. 查询申请对应的用户信息
SELECT 
    a.id AS '申请ID',
    a.application_no AS '申请编号',
    u.username AS '用户名',
    u.identify AS '用户身份',
    a.doctor_title AS '职称',
    a.doctor_dept AS '科室',
    a.status AS '审核状态',
    a.created_at AS '申请时间'
FROM SM_doctor_application a
LEFT JOIN SM_user u ON a.user_id = u.id
ORDER BY a.created_at DESC;

-- 4. 统计申请数量
SELECT 
    status AS '状态',
    COUNT(*) AS '数量',
    CASE 
        WHEN status = 0 THEN '待审核'
        WHEN status = 1 THEN '已通过'
        WHEN status = 2 THEN '已拒绝'
        ELSE '未知'
    END AS '状态说明'
FROM SM_doctor_application 
GROUP BY status;

-- 5. 查看表结构（确认字段是否都存在）
DESCRIBE SM_doctor_application;
