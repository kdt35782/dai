-- 测试插入药品数据
USE SM;

-- 先删除可能存在的数据
DELETE FROM SM_prescription_detail;
DELETE FROM SM_prescription;
DELETE FROM SM_medicine;

-- 插入测试数据(不指定id列,让其自动生成)
INSERT INTO SM_medicine (
    medicine_code, 
    medicine_name, 
    common_name, 
    medicine_type, 
    category, 
    specification, 
    dosage_form, 
    manufacturer, 
    price, 
    indications, 
    usage_dosage, 
    prescription_type, 
    is_otc, 
    status
) VALUES
('M001', '阿莫西林胶囊', '阿莫西林', '西药', '抗生素', '0.25g*24粒', '胶囊', '华北制药', 12.50, 
 '适用于敏感菌所致的呼吸道感染、泌尿道感染等', '口服，成人一次0.5g，每6-8小时1次', 1, 0, 1),
 
('M002', '999感冒灵颗粒', '999感冒灵', '中成药', '感冒药', '10g*9袋', '颗粒', '华润三九', 15.00,
 '清热解毒，用于感冒引起的头痛发热', '开水冲服，一次1袋，每日3次', 0, 1, 1),
 
('M003', '板蓝根颗粒', '板蓝根', '中成药', '感冒药', '10g*20袋', '颗粒', '白云山', 15.00,
 '清热解毒，用于感冒发热、咽喉肿痛', '开水冲服，一次5-10g，每日3-4次', 0, 1, 1);

-- 验证插入结果
SELECT * FROM SM_medicine;
SELECT COUNT(*) as total FROM SM_medicine;
