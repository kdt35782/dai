-- 智能分诊功能数据库增强脚本
-- 为问诊表添加分诊相关字段

USE SM;

-- 1. 添加分诊相关字段
ALTER TABLE SM_consultation 
ADD COLUMN auto_assigned TINYINT DEFAULT 0 COMMENT '是否自动分诊(0:否,1:是)' AFTER need_ai,
ADD COLUMN assigned_reason VARCHAR(200) COMMENT '分诊原因' AFTER auto_assigned,
ADD COLUMN recommended_dept VARCHAR(100) COMMENT 'AI推荐科室' AFTER assigned_reason;

-- 2. 为医生表添加在线状态和负载字段(如果不存在)
ALTER TABLE SM_user 
ADD COLUMN is_online TINYINT DEFAULT 0 COMMENT '是否在线(0:否,1:是)' AFTER status,
ADD COLUMN current_consultation_count INT DEFAULT 0 COMMENT '当前问诊数量' AFTER is_online,
ADD COLUMN max_consultation_count INT DEFAULT 20 COMMENT '最大同时问诊数' AFTER current_consultation_count;

-- 3. 创建医生工作负载统计视图
CREATE OR REPLACE VIEW v_doctor_workload AS
SELECT 
    u.id AS doctor_id,
    u.username,
    u.real_name AS doctor_name,
    u.doctor_dept,
    u.doctor_title,
    u.is_online,
    u.current_consultation_count,
    u.max_consultation_count,
    COUNT(DISTINCT c.id) AS pending_count,
    (u.max_consultation_count - u.current_consultation_count) AS available_slots
FROM SM_user u
LEFT JOIN SM_consultation c ON c.doctor_id = u.id AND c.status IN (0, 1)
WHERE u.identify = 'doctor' AND u.status = 0
GROUP BY u.id
ORDER BY u.is_online DESC, pending_count ASC;

-- 4. 为分诊相关字段添加索引
CREATE INDEX idx_auto_assigned ON SM_consultation(auto_assigned, status);
CREATE INDEX idx_recommended_dept ON SM_consultation(recommended_dept);
CREATE INDEX idx_doctor_online ON SM_user(identify, is_online, doctor_dept);
