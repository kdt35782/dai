# API测试示例

本文档提供了主要API接口的测试示例，可以使用Postman、curl或其他API测试工具进行测试。

## 基础URL

```
http://localhost:3000/api
```

## 1. 用户模块

### 1.1 用户注册

**接口**: `POST /user/register`

**请求示例**:
```json
{
  "username": "zhangsan",
  "password": "Password123!@#",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "verifyCode": "123456"
}
```

**curl命令**:
```bash
curl -X POST http://localhost:3000/api/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "zhangsan",
    "password": "Password123!@#",
    "email": "zhangsan@example.com",
    "phone": "13800138000"
  }'
```

### 1.2 用户登录

**接口**: `POST /user/login`

**请求示例**:
```json
{
  "username": "admin",
  "password": "Admin123!@#"
}
```

**curl命令**:
```bash
curl -X POST http://localhost:3000/api/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123!@#"
  }'
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "tokenType": "Bearer",
    "expiresIn": 7200,
    "userInfo": {
      "userId": 1,
      "username": "admin",
      "role": "admin"
    }
  },
  "timestamp": 1701590400000
}
```

### 1.3 获取用户信息

**接口**: `GET /user/info`

**请求头**:
```
Authorization: Bearer <token>
```

**curl命令**:
```bash
curl -X GET http://localhost:3000/api/user/info \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 1.4 更新用户信息

**接口**: `PUT /user/profile`

**请求头**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "avatar": "https://example.com/avatar.jpg",
  "gender": 1,
  "birthDate": "1990-01-01",
  "realName": "张三",
  "phone": "13800138000"
}
```

### 1.5 修改密码

**接口**: `PUT /user/password`

**请求体**:
```json
{
  "oldPassword": "OldPassword123!@#",
  "newPassword": "NewPassword456!@#",
  "confirmPassword": "NewPassword456!@#"
}
```

### 1.6 申请成为医生

**接口**: `POST /user/apply-doctor`

**请求体**:
```json
{
  "realName": "李医生",
  "phone": "13900139000",
  "certImage": "https://example.com/cert.jpg",
  "doctorTitle": "主治医师",
  "doctorDept": "心内科",
  "specialty": "冠心病、高血压诊治",
  "introduction": "从事心内科临床工作10年",
  "certNumber": "110101199001011234"
}
```

### 1.7 获取医生列表

**接口**: `GET /user/doctors?page=1&pageSize=10&dept=心内科`

**curl命令**:
```bash
curl -X GET "http://localhost:3000/api/user/doctors?page=1&pageSize=10"
```

## 2. 管理员模块

### 2.1 审核医生申请

**接口**: `PUT /user/admin/review-doctor`

**请求头**: `Authorization: Bearer <admin_token>`

**请求体**:
```json
{
  "applicationId": 1,
  "status": 1,
  "rejectReason": ""
}
```

**说明**: status=1表示通过，status=2表示拒绝

### 2.2 获取医生申请列表

**接口**: `GET /user/admin/doctor-applications?page=1&pageSize=10&status=0`

**curl命令**:
```bash
curl -X GET "http://localhost:3000/api/user/admin/doctor-applications?page=1&pageSize=10" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 2.3 获取用户列表

**接口**: `GET /user/admin/users?page=1&pageSize=10&identify=patient`

**参数说明**:
- identify: user/admin/doctor/patient
- status: 0(正常)/1(禁用)
- keyword: 搜索关键词

### 2.4 禁用/启用用户

**接口**: `PUT /user/admin/status`

**请求体**:
```json
{
  "userId": 2,
  "status": 1
}
```

### 2.5 获取登录日志

**接口**: `GET /user/admin/login-logs?page=1&pageSize=10`

## 3. 问诊模块

### 3.1 创建问诊

**接口**: `POST /consultation/create`

**请求头**: `Authorization: Bearer <token>`

**请求体**:
```json
{
  "doctorId": 2,
  "chiefComplaint": "最近3天头晕头痛，血压偏高",
  "symptoms": {
    "age": 45,
    "gender": 1,
    "bloodPressure": "150/95",
    "heartRate": 88,
    "history": ["高血压病史5年"],
    "durationDays": 3
  },
  "needAI": true
}
```

### 3.2 获取问诊列表

**接口**: `GET /consultation/list?page=1&pageSize=10&status=0&role=patient`

**参数说明**:
- status: 0(待接诊)/1(问诊中)/2(已完成)
- role: patient(作为患者)/doctor(作为医生)

### 3.3 获取问诊详情

**接口**: `GET /consultation/detail?consultationId=1`

### 3.4 医生接诊

**接口**: `POST /consultation/accept`

**请求体**:
```json
{
  "consultationId": 1
}
```

### 3.5 完成问诊

**接口**: `POST /consultation/finish`

**请求体**:
```json
{
  "consultationId": 1,
  "diagnosis": "高血压2级，血压控制不佳",
  "prescription": "氨氯地平5mg，每日一次；低盐饮食"
}
```

## 4. 病历模块

### 4.1 获取病历列表

**接口**: `GET /record/list?page=1&pageSize=10`

### 4.2 获取病历详情

**接口**: `GET /record/detail?recordId=1`

## 5. 通知模块

### 5.1 获取通知列表

**接口**: `GET /notification/list?page=1&pageSize=10&isRead=0`

**参数说明**:
- isRead: 0(未读)/1(已读)
- notificationType: system/consultation/audit

### 5.2 获取未读数量

**接口**: `GET /notification/unread-count`

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "totalUnread": 5,
    "systemUnread": 1,
    "consultationUnread": 3,
    "auditUnread": 1
  }
}
```

### 5.3 标记已读

**接口**: `PUT /notification/mark-read`

**请求体**:
```json
{
  "notificationId": 1
}
```

或标记全部已读:
```json
{
  "markAll": true
}
```

## 6. 文件上传

**接口**: `POST /file/upload`

**请求头**: 
- `Authorization: Bearer <token>`
- `Content-Type: multipart/form-data`

**表单字段**:
- file: 文件对象
- fileType: avatar/cert/medical

**curl示例**:
```bash
curl -X POST http://localhost:3000/api/file/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/file.jpg" \
  -F "fileType=avatar"
```

## 7. 密钥管理

### 7.1 生成密钥对

**接口**: `POST /key/generate`

**请求体**:
```json
{
  "keyType": "sm2",
  "purpose": "encryption"
}
```

## Postman集合

可以创建一个Postman集合来管理这些请求：

1. 创建环境变量:
   - `baseUrl`: http://localhost:3000/api
   - `token`: (登录后自动设置)

2. 在Tests标签页添加自动保存token的脚本:
```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    if (jsonData.data && jsonData.data.token) {
        pm.environment.set("token", jsonData.data.token);
    }
}
```

3. 在需要认证的请求中添加Header:
```
Authorization: Bearer {{token}}
```

## 测试流程建议

### 基础流程
1. 管理员登录 → 获取token
2. 创建普通用户
3. 普通用户登录
4. 普通用户申请成为医生
5. 管理员审核医生申请
6. 患者创建问诊
7. 医生接诊
8. 医生完成问诊
9. 查看病历

### 完整测试清单

- [ ] 用户注册
- [ ] 用户登录
- [ ] 获取用户信息
- [ ] 更新用户信息
- [ ] 修改密码
- [ ] 申请成为医生
- [ ] 查询医生列表
- [ ] 管理员审核医生
- [ ] 管理员查看用户列表
- [ ] 创建问诊
- [ ] 查看问诊列表
- [ ] 医生接诊
- [ ] 完成问诊
- [ ] 查看病历
- [ ] 查看通知
- [ ] 文件上传
- [ ] 生成密钥

## 注意事项

1. **密码加密**: 生产环境中，密码应该在前端使用SM3加密后再发送
2. **Token过期**: Token默认2小时过期，过期后需要重新登录
3. **权限验证**: 某些接口需要特定角色权限（如管理员接口）
4. **数据加密**: 系统会自动对敏感数据进行SM4加密存储
5. **CORS**: 如果前端域名不同，需要在后端配置CORS白名单
