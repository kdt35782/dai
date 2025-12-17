-- 检查药品相关表是否存在
USE SM;

-- 查看所有表
SHOW TABLES;

-- 检查药品表
SELECT COUNT(*) as medicine_count FROM SM_medicine;

-- 查看部分药品数据
SELECT medicine_id, medicine_name, category, price, unit 
FROM SM_medicine 
LIMIT 5;

-- 检查处方表
SELECT COUNT(*) as prescription_count FROM SM_prescription;

-- 检查处方明细表
SELECT COUNT(*) as detail_count FROM SM_prescription_detail;
