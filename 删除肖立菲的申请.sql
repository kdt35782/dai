-- 先查看肖立菲的用户ID和申请记录
SELECT u.id, u.username, u.real_name 
FROM SM_user u 
WHERE u.username LIKE '%肖立菲%' OR u.real_name LIKE '%肖立菲%';

-- 查看该用户的申请记录
SELECT a.* 
FROM SM_doctor_application a
INNER JOIN SM_user u ON a.user_id = u.id
WHERE u.username LIKE '%肖立菲%' OR u.real_name LIKE '%肖立菲%';

-- 删除肖立菲的医生申请记录
DELETE a FROM SM_doctor_application a
INNER JOIN SM_user u ON a.user_id = u.id
WHERE u.username LIKE '%肖立菲%' OR u.real_name LIKE '%肖立菲%';

-- 确认删除成功
SELECT COUNT(*) as remaining_count 
FROM SM_doctor_application a
INNER JOIN SM_user u ON a.user_id = u.id
WHERE u.username LIKE '%肖立菲%' OR u.real_name LIKE '%肖立菲%';
