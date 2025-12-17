-- 检查当前表结构
DESCRIBE SM_doctor_application;

-- 添加缺失的字段
ALTER TABLE SM_doctor_application 
ADD COLUMN application_no VARCHAR(50) NOT NULL DEFAULT '' COMMENT '申请编号' AFTER id,
ADD COLUMN email VARCHAR(255) NOT NULL DEFAULT '' COMMENT '邮箱(SM4加密)' AFTER phone,
ADD COLUMN practice_cert VARCHAR(100) NULL COMMENT '执业证书编号' AFTER cert_number,
ADD COLUMN hospital_name VARCHAR(100) NULL COMMENT '所在医院' AFTER practice_cert;

-- 为 application_no 添加唯一索引
ALTER TABLE SM_doctor_application 
ADD UNIQUE INDEX idx_application_no (application_no);

-- 查看更新后的表结构
DESCRIBE SM_doctor_application;
