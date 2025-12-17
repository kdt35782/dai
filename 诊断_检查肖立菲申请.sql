-- 1. 查找肖立菲的用户ID
SELECT id, username, real_name 
FROM SM_user 
WHERE username = '肖立菲' OR username LIKE '%肖立菲%';

-- 2. 直接查询所有医生申请（看是否有残留）
SELECT * FROM SM_doctor_application;

-- 3. 查询肖立菲的申请记录（假设user_id已知，请替换为实际ID）
-- 如果上面查询到肖立菲的user_id是某个值，比如5，则：
SELECT * FROM SM_doctor_application WHERE user_id = 5;

-- 4. 查看表总记录数
SELECT COUNT(*) as total FROM SM_doctor_application;

-- 5. 检查是否有status=0的待审核记录
SELECT * FROM SM_doctor_application WHERE status = 0;
