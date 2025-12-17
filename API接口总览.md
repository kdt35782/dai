# 基于国密加密的网上问诊平台 - API接口总览

## 📋 项目概述

本项目是一个基于国密算法的在线医疗问诊平台，采用国密S系列算法（SM2/SM3/SM4/SM9）对通信内容和数据库数据进行加密，同时集成Paillier同态加密实现隐私保护的AI智能诊断。

---

## 🔐 核心安全特性

### 1. 国密算法应用
- **SM2** - 非对称加密，用于身份认证和敏感数据传输
- **SM3** - 哈希算法，用于密码加密和数据完整性验证
- **SM4** - 对称加密，用于敏感数据存储（邮箱、手机号、病历等）
- **SM9** - 属性基加密，用于病历访问控制

### 2. 高级加密技术
- **Paillier同态加密** - 支持在加密数据上进行AI风险评分计算
- **国密SSL** - 使用SM2算法的SSL证书保障传输层安全

### 3. 隐私保护机制
- 端到端加密通信
- 细粒度访问控制（基于医生属性）
- 完整性验证（防篡改）
- 操作审计日志

---

## 📊 数据库设计（9张表）

| 表名 | 说明 | 加密字段 |
|------|------|----------|
| SM_user | 用户表 | password(SM3), email(SM4), phone(SM4), real_name(SM4), id_card(SM4) |
| SM_doctor_application | 医生申请记录表 | real_name(SM4), id_card(SM4), phone(SM4) |
| SM_login_log | 用户登录日志表 | login_ip(SM4) |
| SM_consultation | 问诊记录表 | chief_complaint(SM4), symptoms(Paillier), diagnosis(SM4), prescription(SM4) |
| SM_consultation_message | 问诊消息表 | content(SM4) |
| SM_medical_record | 电子病历表 | 所有病历内容(SM4), access_policy(SM9), data_hash(SM3) |
| SM_record_access_log | 病历访问日志表 | ip_address(SM4) |
| SM_notification | 系统通知表 | - |
| SM_file | 文件管理表 | 文件内容(可选加密), file_hash(SM3) |

---

## 🔌 API接口分类（共35+个接口）

### 一、用户模块（15个接口）

#### 普通用户接口（10个）
1. ✅ **POST** `/api/user/register` - 用户注册
2. ✅ **POST** `/api/user/login` - 用户登录
3. ✅ **GET** `/api/user/info` - 获取当前用户信息
4. ✅ **PUT** `/api/user/info` - 更新用户信息
5. ✅ **PUT** `/api/user/password` - 修改密码
6. ✅ **POST** `/api/user/apply-doctor` - 申请成为医生
7. ✅ **GET** `/api/user/doctor-application` - 查询医生申请状态
8. ✅ **GET** `/api/user/doctors` - 获取医生列表
9. ✅ **GET** `/api/user/doctor/:userId` - 获取医生详情
10. ✅ **POST** `/api/user/logout` - 用户退出登录

#### 管理员接口（5个）
11. ✅ **PUT** `/api/user/admin/review-doctor` - 审核医生申请
12. ✅ **GET** `/api/user/admin/doctor-applications` - 获取医生申请列表
13. ✅ **GET** `/api/user/admin/users` - 获取用户列表
14. ✅ **PUT** `/api/user/admin/status` - 禁用/启用用户
15. ✅ **GET** `/api/user/admin/login-logs` - 获取登录日志

---

### 二、国密密钥管理模块（3个接口）

1. ✅ **GET** `/api/crypto/public-key` - 获取系统公钥（SM2/Paillier）
2. ✅ **POST** `/api/crypto/generate-keypair` - 生成用户密钥对
3. ✅ **POST** `/api/crypto/session-key` - 生成会话密钥

**功能说明**：
- 前端通过此模块获取加密所需的公钥
- 支持SM2和Paillier两种密钥类型
- 会话密钥用于临时数据加密

---

### 三、文件上传模块（2个接口）

1. ✅ **POST** `/api/file/upload` - 上传文件（头像/证书/病历附件）
2. ✅ **GET** `/api/file/download/:fileId` - 下载文件

**安全特性**：
- 文件自动SM4加密存储
- SM3完整性验证
- 权限控制下载

---

### 四、问诊模块（7个接口）

1. ✅ **POST** `/api/consultation/create` - 创建问诊（患者发起）
2. ✅ **GET** `/api/consultation/list` - 获取问诊列表
3. ✅ **GET** `/api/consultation/:consultationId` - 获取问诊详情
4. ✅ **POST** `/api/consultation/:consultationId/accept` - 医生接诊
5. ✅ **POST** `/api/consultation/:consultationId/message` - 发送问诊消息
6. ✅ **GET** `/api/consultation/:consultationId/messages` - 获取问诊消息记录
7. ✅ **POST** `/api/consultation/:consultationId/complete` - 完成问诊（医生）

**核心功能**：
- 🤖 **AI智能诊断** - Paillier同态加密保护下的风险评分
- 💬 **实时消息** - 支持文本、图片、语音、处方
- 🔒 **端到端加密** - 所有内容SM4加密存储

---

### 五、病历模块（5个接口）

1. ✅ **POST** `/api/medical-record/create` - 创建病历（医生）
2. ✅ **GET** `/api/medical-record/:recordId` - 查看病历
3. ✅ **GET** `/api/medical-record/patient/:patientId` - 获取患者病历列表
4. ✅ **POST** `/api/medical-record/:recordId/authorize` - 授权病历访问（患者）
5. ✅ **GET** `/api/medical-record/:recordId/access-logs` - 查看病历访问日志

**安全亮点**：
- 🛡️ **SM9属性基加密** - 基于医生属性（科室、职称）的访问控制
- ✅ **完整性验证** - SM3哈希确保病历未被篡改
- 📝 **审计追溯** - 所有访问操作记录可查

---

### 六、消息通知模块（4个接口）

1. ✅ **GET** `/api/notification/list` - 获取通知列表
2. ✅ **PUT** `/api/notification/:notificationId/read` - 标记通知已读
3. ✅ **PUT** `/api/notification/read-all` - 标记全部已读
4. ✅ **GET** `/api/notification/unread-count` - 获取未读数量

**通知类型**：
- 系统通知
- 问诊消息通知
- 审核结果通知
- 提醒通知

---

## 🔄 典型业务流程

### 流程1：患者在线问诊
```
1. 患者注册/登录 → 获取JWT Token
2. 患者创建问诊 → 系统AI智能诊断（Paillier加密计算）
3. 系统分配医生 / 患者选择医生
4. 医生接诊 → 实时加密消息交流
5. 医生开具处方 → 完成问诊
6. 系统自动创建加密病历
```

### 流程2：医生申请与审核
```
1. 用户注册为普通用户
2. 提交医生资格申请 → 上传资格证书
3. 管理员审核申请
4. 审核通过 → 用户身份变更为医生
5. 医生可以接诊和开具处方
```

### 流程3：病历访问控制
```
1. 医生创建病历 → 设置访问策略（允许的科室和职称）
2. 病历使用SM9属性基加密存储
3. 其他医生访问时 → 验证属性是否匹配
4. 属性匹配 → 解密查看 | 不匹配 → 拒绝访问
5. 所有访问记录在审计日志
```

---

## 🎯 技术实现要点

### 前端加密处理
```javascript
// 1. 获取系统公钥
const sm2PublicKey = await getPublicKey('sm2');

// 2. 密码前端SM3加密
const hashedPassword = sm3Hash(rawPassword);

// 3. 敏感数据SM2加密传输
const encryptedData = sm2Encrypt(sensitiveData, sm2PublicKey);

// 4. 症状数据Paillier加密（用于AI）
const paillierKey = await getPublicKey('paillier');
const encryptedSymptoms = paillierEncrypt(symptoms, paillierKey);
```

### 后端加密处理
```python
# 1. 密码SM3+盐值双重加密
password_hash = sm3_hash(sm3_hash(password) + salt)

# 2. 敏感数据SM4加密存储
encrypted_email = sm4_encrypt(email, master_key)

# 3. Paillier同态计算
ai_score = paillier_compute(encrypted_features, model_weights)

# 4. SM9属性基加密
encrypted_record = sm9_encrypt(record, access_policy)
```

---

## 📈 系统扩展建议

### 近期优化
1. ✨ 实现WebSocket实时消息推送
2. ✨ 集成邮箱/短信验证码服务
3. ✨ 添加图片/语音消息支持
4. ✨ 完善AI诊断模型训练

### 中期规划
1. 🚀 区块链存证集成（操作哈希上链）
2. 🚀 多方安全计算（MPC）增强隐私保护
3. 🚀 联邦学习实现多医院协同AI训练
4. 🚀 支付模块（问诊费用）

### 长期展望
1. 🌟 移动端APP开发
2. 🌟 智能健康档案管理
3. 🌟 远程会诊功能
4. 🌟 健康大数据分析（隐私保护下）

---

## 🔧 开发注意事项

### 安全规范
1. ✅ 所有密码必须SM3加密后传输和存储
2. ✅ 敏感字段（邮箱、手机号、病历）必须SM4加密
3. ✅ 关键操作必须记录审计日志
4. ✅ API接口必须进行权限校验
5. ✅ 数据传输必须使用HTTPS + 国密SSL

### 测试要点
1. 🧪 加密/解密功能测试
2. 🧪 属性基访问控制测试
3. 🧪 完整性验证测试
4. 🧪 并发访问测试
5. 🧪 性能压力测试

---

## 📞 联系方式

- 项目文档：见 `开发文档.md`
- 数据库脚本：见 `database/SM_user_module.sql`
- 系统规划：见 `规划.md`

---

**最后更新时间**: 2024-12-03
**文档版本**: v1.0
