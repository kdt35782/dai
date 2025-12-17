-- 为问诊表添加 ai_suggestions 字段

-- 1. 检查当前表结构
DESCRIBE SM_consultation;

-- 2. 添加 ai_suggestions 字段
ALTER TABLE SM_consultation 
ADD COLUMN ai_suggestions TEXT COMMENT 'AI就医建议' AFTER ai_diagnosis;

-- 3. 查看更新后的表结构
DESCRIBE SM_consultation;

-- 4. 验证字段已添加
SELECT COLUMN_NAME, COLUMN_TYPE, COLUMN_COMMENT 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'SM_consultation' 
  AND TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'ai_suggestions';
