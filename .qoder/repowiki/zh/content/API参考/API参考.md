# API参考

<cite>
**本文档引用文件**  
- [API接口总览.md](file://API接口总览.md)
- [routes.go](file://backed/internal/api/routes.go)
- [user_handler.go](file://backed/internal/api/handler/user_handler.go)
- [key_handler.go](file://backed/internal/api/handler/key_handler.go)
- [file_handler.go](file://backed/internal/api/handler/file_handler.go)
- [consultation_handler.go](file://backed/internal/api/handler/consultation_handler.go)
- [record_handler.go](file://backed/internal/api/handler/record_handler.go)
- [notification_handler.go](file://backed/internal/api/handler/notification_handler.go)
- [admin_handler.go](file://backed/internal/api/handler/admin_handler.go)
- [auth.go](file://backed/internal/middleware/auth.go)
- [response.go](file://backed/pkg/utils/response.go)
- [jwt.go](file://backed/pkg/utils/jwt.go)
- [config.yaml](file://backed/config/config.yaml)
- [开发文档.md](file://开发文档.md)
</cite>

## 目录
1. [简介](#简介)
2. [通用规范](#通用规范)
3. [用户模块](#用户模块)
4. [国密密钥管理](#国密密钥管理)
5. [文件上传](#文件上传)
6. [问诊模块](#问诊模块)
7. [病历模块](#病历模块)
8. [消息通知](#消息通知)
9. [管理员模块](#管理员模块)
10. [错误码说明](#错误码说明)

## 简介
本API参考文档详细描述了基于国密加密的网上看诊系统的RESTful API接口。系统采用国密S系列算法（SM2/SM3/SM4/SM9）和Paillier同态加密技术，确保用户数据在传输和存储过程中的安全性。文档涵盖了用户管理、密钥管理、文件上传、问诊、病历和消息通知等核心功能模块。

所有API均通过JWT进行身份认证，并遵循统一的响应格式。前端在调用需要认证的接口时，必须在请求头中包含有效的Authorization令牌。

**文档来源**
- [API接口总览.md](file://API接口总览.md#L1-L259)

## 通用规范

### 基础信息
- **API基础URL**: `http://localhost:3000/api`（开发环境）
- **请求格式**: JSON
- **响应格式**: 统一JSON格式
- **字符编码**: UTF-8

### 认证机制
系统采用JWT（JSON Web Token）进行身份认证。用户登录成功后将获得一个token，后续请求需在Header中携带：

```
Authorization: Bearer <token>
```

token有效期为2小时（7200秒），过期后需重新登录获取。

### 响应格式
所有API响应遵循统一格式：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {},
  "timestamp": 1701590400000
}
```

**字段说明**：
- `code`: 状态码（200表示成功）
- `message`: 响应消息
- `data`: 返回数据（成功时）或null（失败时）
- `timestamp`: 时间戳（毫秒）

### 分页参数
大多数列表接口支持分页查询，通用分页参数如下：

| 参数名 | 类型 | 必填 | 说明 | 默认值 |
|--------|------|------|------|--------|
| page | int | 否 | 页码 | 1 |
| pageSize | int | 否 | 每页数量 | 10 |

**通用来源**
- [response.go](file://backed/pkg/utils/response.go#L8-L68)
- [jwt.go](file://backed/pkg/utils/jwt.go#L10-L51)
- [auth.go](file://backed/internal/middleware/auth.go#L9-L42)
- [config.yaml](file://backed/config/config.yaml#L21-L23)

## 用户模块

### 用户注册
创建新用户账户。

**接口信息**
- **HTTP方法**: POST
- **URL**: `/api/user/register`
- **认证要求**: 无需认证

**请求头**
```
Content-Type: application/json
```

**请求体**
```json
{
  "username": "string",
  "password": "string",
  "email": "string",
  "phone": "string"
}
```

**字段说明**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码（前端需SM3加密） |
| email | string | 是 | 邮箱地址 |
| phone | string | 是 | 手机号码 |

**成功响应**
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "userId": 1001,
    "username": "zhangsan",
    "email": "zhangsan@example.com"
  }
}
```

**错误码**
- 400: 参数错误或用户已存在

**示例请求**
```bash
curl -X POST http://localhost:3000/api/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "zhangsan",
    "password": "e10adc3949ba59abbe56e057f20f883e",
    "email": "zhangsan@example.com",
    "phone": "13800138000"
  }'
```

**来源**
- [user_handler.go](file://backed/internal/api/handler/user_handler.go#L21-L47)
- [routes.go](file://backed/internal/api/routes.go#L19)

### 用户登录
用户身份验证并获取访问令牌。

**接口信息**
- **HTTP方法**: POST
- **URL**: `/api/user/login`
- **认证要求**: 无需认证

**请求头**
```
Content-Type: application/json
```

**请求体**
```json
{
  "username": "string",
  "password": "string"
}
```

**字段说明**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码（前端需SM3加密） |

**成功响应**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "tokenType": "Bearer",
    "expiresIn": 7200,
    "userInfo": {
      "userId": 1001,
      "username": "zhangsan",
      "email": "zhangsan@example.com",
      "identify": "user",
      "avatar": "https://example.com/avatar/1001.jpg",
      "status": 0
    }
  }
}
```

**错误码**
- 401: 用户名或密码错误
- 403: 账号已被禁用

**示例请求**
```bash
curl -X POST http://localhost:3000/api/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "zhangsan",
    "password": "e10adc3949ba59abbe56e057f20f883e"
  }'
```

**来源**
- [user_handler.go](file://backed/internal/api/handler/user_handler.go#L49-L82)
- [routes.go](file://backed/internal/api/routes.go#L20)
- [开发文档.md](file://开发文档.md#L243-L274)

### 获取用户信息
获取当前登录用户或指定用户的信息。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/user/info`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | int64 | 否 | 用户ID（查询其他用户时使用） |

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "userId": 1001,
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "realName": "张三",
    "gender": 1,
    "birthDate": "1990-01-01",
    "avatar": "https://example.com/avatar/1001.jpg",
    "identify": "user",
    "status": 0,
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

**错误码**
- 401: 未登录或token无效
- 404: 用户不存在

**示例请求**
```bash
curl -X GET http://localhost:3000/api/user/info \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [user_handler.go](file://backed/internal/api/handler/user_handler.go#L84-L113)
- [routes.go](file://backed/internal/api/routes.go#L28)

### 更新用户信息
修改当前用户的个人信息。

**接口信息**
- **HTTP方法**: PUT
- **URL**: `/api/user/profile`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "avatar": "string",
  "realName": "string",
  "gender": 0,
  "birthDate": "string",
  "phone": "string",
  "email": "string"
}
```

**字段说明**
| 字段 | 类型 | 必填 | 说明 | 取值范围 |
|------|------|------|------|--------|
| avatar | string | 否 | 头像URL | |
| realName | string | 否 | 真实姓名 | |
| gender | int | 否 | 性别 | 0:未知, 1:男, 2:女 |
| birthDate | string | 否 | 出生日期 | YYYY-MM-DD |
| phone | string | 否 | 手机号码 | |
| email | string | 否 | 邮箱地址 | |

**成功响应**
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

**错误码**
- 400: 参数错误
- 500: 服务器内部错误

**示例请求**
```bash
curl -X PUT http://localhost:3000/api/user/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "realName": "张三",
    "gender": 1,
    "birthDate": "1990-01-01",
    "phone": "13800138000",
    "email": "zhangsan@example.com"
  }'
```

**来源**
- [user_handler.go](file://backed/internal/api/handler/user_handler.go#L116-L139)
- [routes.go](file://backed/internal/api/routes.go#L29)

### 修改密码
更改当前用户的登录密码。

**接口信息**
- **HTTP方法**: PUT
- **URL**: `/api/user/password`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "oldPassword": "string",
  "newPassword": "string",
  "confirmPassword": "string"
}
```

**字段说明**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| oldPassword | string | 是 | 旧密码（前端需SM3加密） |
| newPassword | string | 是 | 新密码（前端需SM3加密） |
| confirmPassword | string | 是 | 确认新密码 |

**成功响应**
```json
{
  "code": 200,
  "message": "密码修改成功，请重新登录",
  "data": null
}
```

**错误码**
- 400: 旧密码错误或两次新密码不一致
- 401: 旧密码错误

**示例请求**
```bash
curl -X PUT http://localhost:3000/api/user/password \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "oldPassword": "e10adc3949ba59abbe56e057f20f883e",
    "newPassword": "25d55ad283aa400af464c76d713c07ad",
    "confirmPassword": "25d55ad283aa400af464c76d713c07ad"
  }'
```

**来源**
- [user_handler.go](file://backed/internal/api/handler/user_handler.go#L142-L167)
- [routes.go](file://backed/internal/api/routes.go#L30)
- [开发文档.md](file://开发文档.md#L411-L417)

### 申请成为医生
普通用户申请成为认证医生。

**接口信息**
- **HTTP方法**: POST
- **URL**: `/api/user/apply-doctor`
- **认证要求**: 需要认证（普通用户）

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "realName": "string",
  "idCard": "string",
  "phone": "string",
  "certImage": "string",
  "doctorTitle": "string",
  "doctorDept": "string",
  "specialty": "string",
  "introduction": "string",
  "certNumber": "string"
}
```

**字段说明**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| realName | string | 是 | 真实姓名 |
| idCard | string | 否 | 身份证号码 |
| phone | string | 是 | 手机号码 |
| certImage | string | 是 | 资格证书图片URL |
| doctorTitle | string | 是 | 医生职称 |
| doctorDept | string | 是 | 所属科室 |
| specialty | string | 否 | 专业特长 |
| introduction | string | 否 | 个人简介 |
| certNumber | string | 否 | 证书编号 |

**成功响应**
```json
{
  "code": 200,
  "message": "申请已提交，请等待管理员审核",
  "data": {
    "applicationId": 101
  }
}
```

**错误码**
- 400: 参数错误或已提交申请
- 500: 服务器内部错误

**示例请求**
```bash
curl -X POST http://localhost:3000/api/user/apply-doctor \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "realName": "张医生",
    "phone": "13800138001",
    "certImage": "/uploads/cert/doctor_123456789.jpg",
    "doctorTitle": "主任医师",
    "doctorDept": "内科",
    "specialty": "心血管疾病",
    "introduction": "10年临床经验"
  }'
```

**来源**
- [user_handler.go](file://backed/internal/api/handler/user_handler.go#L170-L202)
- [routes.go](file://backed/internal/api/routes.go#L31)

### 查询医生申请状态
查看当前用户的医生申请状态。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/user/doctor-application`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
```

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "applicationId": 101,
    "userId": 1001,
    "realName": "张医生",
    "idCard": "123456789012345678",
    "phone": "13800138001",
    "certImage": "/uploads/cert/doctor_123456789.jpg",
    "doctorTitle": "主任医师",
    "doctorDept": "内科",
    "specialty": "心血管疾病",
    "introduction": "10年临床经验",
    "certNumber": "ZY123456",
    "status": 0,
    "rejectReason": "",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-02T00:00:00Z"
  }
}
```

**状态说明**
| 状态值 | 说明 |
|--------|------|
| 0 | 待审核 |
| 1 | 已通过 |
| 2 | 已拒绝 |

**错误码**
- 404: 未找到申请记录

**示例请求**
```bash
curl -X GET http://localhost:3000/api/user/doctor-application \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [user_handler.go](file://backed/internal/api/handler/user_handler.go#L204-L214)
- [routes.go](file://backed/internal/api/routes.go#L32)

### 获取医生列表
获取所有认证医生的列表。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/user/doctors`
- **认证要求**: 无需认证

**请求头**
```
Content-Type: application/json
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 | 1 |
| pageSize | int | 否 | 每页数量 | 10 |
| dept | string | 否 | 科室筛选 |
| keyword | string | 否 | 关键词搜索 |

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "list": [
      {
        "userId": 1002,
        "username": "zhangdoctor",
        "realName": "张医生",
        "avatar": "https://example.com/avatar/1002.jpg",
        "doctorTitle": "主任医师",
        "doctorDept": "内科",
        "specialty": "心血管疾病",
        "introduction": "10年临床经验",
        "consultationCount": 150,
        "rating": 4.8
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  }
}
```

**错误码**
- 500: 服务器内部错误

**示例请求**
```bash
curl -X GET "http://localhost:3000/api/user/doctors?page=1&pageSize=10&dept=内科" \
  -H "Content-Type: application/json"
```

**来源**
- [user_handler.go](file://backed/internal/api/handler/user_handler.go#L217-L235)
- [routes.go](file://backed/internal/api/routes.go#L21)

### 获取医生详情
获取指定医生的详细信息。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/user/doctor/:userId`
- **认证要求**: 无需认证

**请求头**
```
Content-Type: application/json
```

**路径参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | int64 | 是 | 医生用户ID |

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "userId": 1002,
    "username": "zhangdoctor",
    "realName": "张医生",
    "avatar": "https://example.com/avatar/1002.jpg",
    "gender": 1,
    "birthDate": "1980-01-01",
    "phone": "13800138001",
    "email": "zhangdoctor@example.com",
    "doctorTitle": "主任医师",
    "doctorDept": "内科",
    "specialty": "心血管疾病",
    "introduction": "10年临床经验",
    "certImage": "/uploads/cert/doctor_123456789.jpg",
    "certNumber": "ZY123456",
    "consultationCount": 150,
    "rating": 4.8,
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

**错误码**
- 400: 用户ID格式错误
- 404: 医生不存在

**示例请求**
```bash
curl -X GET http://localhost:3000/api/user/doctor/1002 \
  -H "Content-Type: application/json"
```

**来源**
- [user_handler.go](file://backed/internal/api/handler/user_handler.go#L238-L253)
- [routes.go](file://backed/internal/api/routes.go#L22)

### 用户退出登录
用户退出当前会话。

**接口信息**
- **HTTP方法**: POST
- **URL**: `/api/user/logout`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
```

**成功响应**
```json
{
  "code": 200,
  "message": "退出成功",
  "data": null
}
```

**说明**
由于JWT是无状态的，服务器端无法主动使token失效。此接口仅作为客户端清理本地存储的提示，实际的token失效依赖于过期时间。

**示例请求**
```bash
curl -X POST http://localhost:3000/api/user/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [user_handler.go](file://backed/internal/api/handler/user_handler.go#L256-L259)
- [routes.go](file://backed/internal/api/routes.go#L33)

## 国密密钥管理

### 获取系统公钥
获取用于前端加密的系统公钥。

**接口信息**
- **HTTP方法**: POST
- **URL**: `/api/key/generate`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
空

**成功响应**
```json
{
  "code": 200,
  "message": "密钥已生成",
  "data": {
    "publicKey": "048d3a5c6b7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d"
  }
}
```

**字段说明**
- `publicKey`: SM2算法的公钥（十六进制格式）

**用途**
前端使用此公钥对敏感数据进行SM2加密后再传输到服务器。

**错误码**
- 401: 未登录或token无效
- 500: 服务器内部错误

**示例请求**
```bash
curl -X POST http://localhost:3000/api/key/generate \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

**来源**
- [key_handler.go](file://backed/internal/api/handler/key_handler.go#L15-L22)
- [routes.go](file://backed/internal/api/routes.go#L72)

## 文件上传

### 上传文件
上传用户头像、医生资格证书等文件。

**接口信息**
- **HTTP方法**: POST
- **URL**: `/api/file/upload`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
```

**表单参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| file | file | 是 | 要上传的文件 |
| fileType | string | 否 | 文件类型（avatar, cert, record等） | general |

**限制**
- 最大文件大小：10MB
- 允许类型：image/jpeg, image/png, image/jpg

**成功响应**
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "fileUrl": "/uploads/avatar_1701590400000.jpg",
    "fileName": "avatar.jpg",
    "fileSize": 10240,
    "uploadTime": "2024-01-01 12:00:00"
  }
}
```

**错误码**
- 400: 未选择文件或文件大小超限
- 500: 文件保存失败

**示例请求**
```bash
curl -X POST http://localhost:3000/api/file/upload \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -F "file=@avatar.jpg" \
  -F "fileType=avatar"
```

**来源**
- [file_handler.go](file://backed/internal/api/handler/file_handler.go#L18-L56)
- [routes.go](file://backed/internal/api/routes.go#L80)
- [config.yaml](file://backed/config/config.yaml#L30-L36)

## 问诊模块

### 创建问诊
患者发起新的问诊请求。

**接口信息**
- **HTTP方法**: POST
- **URL**: `/api/consultation/create`
- **认证要求**: 需要认证（患者）

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "doctorId": 1002,
  "chiefComplaint": "持续胸痛3天",
  "symptoms": {
    "duration": "3天",
    "frequency": "持续",
    "severity": 8,
    "location": "胸部中央",
    "accompanyingSymptoms": ["呼吸困难", "出汗"]
  },
  "needAI": true
}
```

**字段说明**
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| doctorId | int64 | 否 | 指定医生ID（null表示系统分配） |
| chiefComplaint | string | 是 | 主诉 |
| symptoms | object | 是 | 症状详情（前端需Paillier加密） |
| needAI | boolean | 否 | 是否需要AI智能诊断 | false |

**成功响应**
```json
{
  "code": 200,
  "message": "问诊创建成功",
  "data": {
    "consultationId": 2001,
    "patientId": 1001,
    "doctorId": 1002,
    "status": 0,
    "chiefComplaint": "持续胸痛3天",
    "aiScore": 0.75,
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

**状态说明**
| 状态值 | 说明 |
|--------|------|
| 0 | 待接诊 |
| 1 | 进行中 |
| 2 | 已完成 |

**错误码**
- 400: 参数错误
- 500: 服务器内部错误

**示例请求**
```bash
curl -X POST http://localhost:3000/api/consultation/create \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "doctorId": 1002,
    "chiefComplaint": "持续胸痛3天",
    "symptoms": {"duration": "3天", "severity": 8},
    "needAI": true
  }'
```

**来源**
- [consultation_handler.go](file://backed/internal/api/handler/consultation_handler.go#L20-L42)
- [routes.go](file://backed/internal/api/routes.go#L42)

### 获取问诊列表
获取当前用户的问诊记录列表。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/consultation/list`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 | 1 |
| pageSize | int | 否 | 每页数量 | 10 |
| role | string | 否 | 角色 | patient |
| status | int | 否 | 状态筛选 |

**角色说明**
- `patient`: 患者视角
- `doctor`: 医生视角

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "list": [
      {
        "consultationId": 2001,
        "patientId": 1001,
        "doctorId": 1002,
        "patientName": "张三",
        "doctorName": "张医生",
        "chiefComplaint": "持续胸痛3天",
        "status": 2,
        "aiScore": 0.75,
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T01:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  }
}
```

**错误码**
- 500: 服务器内部错误

**示例请求**
```bash
curl -X GET "http://localhost:3000/api/consultation/list?role=patient&page=1&pageSize=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [consultation_handler.go](file://backed/internal/api/handler/consultation_handler.go#L45-L70)
- [routes.go](file://backed/internal/api/routes.go#L43)

### 获取问诊详情
获取指定问诊的详细信息。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/consultation/detail`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| consultationId | int64 | 是 | 问诊ID |

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "consultationId": 2001,
    "patientId": 1001,
    "doctorId": 1002,
    "patientName": "张三",
    "doctorName": "张医生",
    "chiefComplaint": "持续胸痛3天",
    "symptoms": {"duration": "3天", "severity": 8},
    "diagnosis": "心绞痛",
    "prescription": "硝酸甘油片 0.5mg",
    "status": 2,
    "aiScore": 0.75,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T01:00:00Z"
  }
}
```

**错误码**
- 400: 问诊ID格式错误
- 404: 问诊不存在

**示例请求**
```bash
curl -X GET "http://localhost:3000/api/consultation/detail?consultationId=2001" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [consultation_handler.go](file://backed/internal/api/handler/consultation_handler.go#L73-L90)
- [routes.go](file://backed/internal/api/routes.go#L44)

### 医生接诊
医生接受问诊请求，开始诊疗。

**接口信息**
- **HTTP方法**: POST
- **URL**: `/api/consultation/accept`
- **认证要求**: 需要认证（医生）

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "consultationId": 2001
}
```

**成功响应**
```json
{
  "code": 200,
  "message": "接诊成功",
  "data": null
}
```

**错误码**
- 400: 参数错误或无权接诊
- 500: 服务器内部错误

**示例请求**
```bash
curl -X POST http://localhost:3000/api/consultation/accept \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{"consultationId": 2001}'
```

**来源**
- [consultation_handler.go](file://backed/internal/api/handler/consultation_handler.go#L93-L111)
- [routes.go](file://backed/internal/api/routes.go#L45)

### 完成问诊
医生完成问诊，填写诊断和处方。

**接口信息**
- **HTTP方法**: POST
- **URL**: `/api/consultation/finish`
- **认证要求**: 需要认证（医生）

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "consultationId": 2001,
  "diagnosis": "心绞痛",
  "prescription": "硝酸甘油片 0.5mg"
}
```

**成功响应**
```json
{
  "code": 200,
  "message": "问诊已完成",
  "data": null
}
```

**错误码**
- 400: 参数错误或无权操作
- 500: 服务器内部错误

**示例请求**
```bash
curl -X POST http://localhost:3000/api/consultation/finish \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "consultationId": 2001,
    "diagnosis": "心绞痛",
    "prescription": "硝酸甘油片 0.5mg"
  }'
```

**来源**
- [consultation_handler.go](file://backed/internal/api/handler/consultation_handler.go#L114-L134)
- [routes.go](file://backed/internal/api/routes.go#L46)

## 病历模块

### 获取病历列表
获取当前用户的病历记录列表。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/record/list`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 | 1 |
| pageSize | int | 否 | 每页数量 | 10 |
| startDate | string | 否 | 开始日期 | YYYY-MM-DD |
| endDate | string | 否 | 结束日期 | YYYY-MM-DD |

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "list": [
      {
        "recordId": 3001,
        "patientId": 1001,
        "patientName": "张三",
        "doctorId": 1002,
        "doctorName": "张医生",
        "diagnosis": "心绞痛",
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T01:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  }
}
```

**错误码**
- 500: 服务器内部错误

**示例请求**
```bash
curl -X GET "http://localhost:3000/api/record/list?page=1&pageSize=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [record_handler.go](file://backed/internal/api/handler/record_handler.go#L20-L39)
- [routes.go](file://backed/internal/api/routes.go#L54)

### 获取病历详情
获取指定病历的详细信息。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/record/detail`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| recordId | int64 | 是 | 病历ID |

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "recordId": 3001,
    "patientId": 1001,
    "patientName": "张三",
    "doctorId": 1002,
    "doctorName": "张医生",
    "diagnosis": "心绞痛",
    "prescription": "硝酸甘油片 0.5mg",
    "symptoms": {"duration": "3天", "severity": 8},
    "aiSuggestions": "建议进行心电图检查",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T01:00:00Z"
  }
}
```

**错误码**
- 400: 病历ID格式错误
- 404: 病历不存在或无权访问

**示例请求**
```bash
curl -X GET "http://localhost:3000/api/record/detail?recordId=3001" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [record_handler.go](file://backed/internal/api/handler/record_handler.go#L42-L59)
- [routes.go](file://backed/internal/api/routes.go#L55)

## 消息通知

### 获取通知列表
获取当前用户的通知列表。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/notification/list`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 | 1 |
| pageSize | int | 否 | 每页数量 | 10 |
| type | string | 否 | 通知类型筛选 |

**通知类型**
- `system`: 系统通知
- `consultation`: 问诊消息
- `review`: 审核结果
- `reminder`: 提醒通知

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "list": [
      {
        "notificationId": 4001,
        "title": "问诊已完成",
        "content": "您的问诊已由张医生完成，请查看诊断结果。",
        "type": "consultation",
        "read": false,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  }
}
```

**错误码**
- 500: 服务器内部错误

**示例请求**
```bash
curl -X GET "http://localhost:3000/api/notification/list?page=1&pageSize=10&type=consultation" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [notification_handler.go](file://backed/internal/api/handler/notification_handler.go#L20-L38)
- [routes.go](file://backed/internal/api/routes.go#L63)

### 获取未读数量
获取当前用户的未读通知数量。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/notification/unread-count`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
```

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": 3
}
```

**错误码**
- 500: 服务器内部错误

**示例请求**
```bash
curl -X GET http://localhost:3000/api/notification/unread-count \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [notification_handler.go](file://backed/internal/api/handler/notification_handler.go#L41-L51)
- [routes.go](file://backed/internal/api/routes.go#L64)

### 标记已读
将指定通知标记为已读。

**接口信息**
- **HTTP方法**: PUT
- **URL**: `/api/notification/mark-read`
- **认证要求**: 需要认证

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "notificationIds": [4001, 4002]
}
```

**成功响应**
```json
{
  "code": 200,
  "message": "已标记为已读",
  "data": null
}
```

**错误码**
- 400: 参数错误
- 500: 服务器内部错误

**示例请求**
```bash
curl -X PUT http://localhost:3000/api/notification/mark-read \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{"notificationIds": [4001, 4002]}'
```

**来源**
- [notification_handler.go](file://backed/internal/api/handler/notification_handler.go#L54-L72)
- [routes.go](file://backed/internal/api/routes.go#L65)

## 管理员模块

### 审核医生申请
管理员审核用户的医生资格申请。

**接口信息**
- **HTTP方法**: PUT
- **URL**: `/api/user/admin/review-doctor`
- **认证要求**: 需要认证（管理员）

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "applicationId": 101,
  "status": 1,
  "rejectReason": ""
}
```

**字段说明**
| 字段 | 类型 | 必填 | 说明 | 取值范围 |
|------|------|------|------|--------|
| applicationId | int64 | 是 | 申请ID | |
| status | int | 是 | 审核状态 | 1:通过, 2:拒绝 |
| rejectReason | string | 否 | 拒绝原因（拒绝时必填） | |

**成功响应**
```json
{
  "code": 200,
  "message": "审核成功",
  "data": null
}
```

**错误码**
- 400: 参数错误或拒绝原因为空
- 403: 无权限访问
- 500: 服务器内部错误

**示例请求**
```bash
curl -X PUT http://localhost:3000/api/user/admin/review-doctor \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "applicationId": 101,
    "status": 1
  }'
```

**来源**
- [admin_handler.go](file://backed/internal/api/handler/admin_handler.go#L20-L51)
- [routes.go](file://backed/internal/api/routes.go#L88)

### 获取医生申请列表
获取所有医生申请的列表。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/user/admin/doctor-applications`
- **认证要求**: 需要认证（管理员）

**请求头**
```
Authorization: Bearer <token>
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 | 1 |
| pageSize | int | 否 | 每页数量 | 10 |
| status | int | 否 | 状态筛选 | 0:待审核, 1:已通过, 2:已拒绝 |

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "total": 1,
    "page": 1,
    "pageSize": 10,
    "list": [
      {
        "applicationId": 101,
        "userId": 1001,
        "realName": "张医生",
        "idCard": "123456789012345678",
        "phone": "13800138001",
        "certImage": "/uploads/cert/doctor_123456789.jpg",
        "doctorTitle": "主任医师",
        "doctorDept": "内科",
        "specialty": "心血管疾病",
        "status": 0,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

**错误码**
- 403: 无权限访问
- 500: 服务器内部错误

**示例请求**
```bash
curl -X GET "http://localhost:3000/api/user/admin/doctor-applications?page=1&pageSize=10&status=0" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [admin_handler.go](file://backed/internal/api/handler/admin_handler.go#L54-L83)
- [routes.go](file://backed/internal/api/routes.go#L89)

### 获取用户列表
获取所有用户的列表。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/user/admin/users`
- **认证要求**: 需要认证（管理员）

**请求头**
```
Authorization: Bearer <token>
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 | 1 |
| pageSize | int | 否 | 每页数量 | 10 |
| identify | string | 否 | 身份筛选 | user, doctor, admin |
| status | int | 否 | 状态筛选 | 0:启用, 1:禁用 |
| keyword | string | 否 | 关键词搜索 |

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "total": 2,
    "page": 1,
    "pageSize": 10,
    "list": [
      {
        "userId": 1001,
        "username": "zhangsan",
        "email": "zhangsan@example.com",
        "phone": "13800138000",
        "identify": "user",
        "status": 0,
        "createdAt": "2024-01-01T00:00:00Z"
      },
      {
        "userId": 1002,
        "username": "zhangdoctor",
        "email": "zhangdoctor@example.com",
        "phone": "13800138001",
        "identify": "doctor",
        "status": 0,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

**错误码**
- 403: 无权限访问
- 500: 服务器内部错误

**示例请求**
```bash
curl -X GET "http://localhost:3000/api/user/admin/users?page=1&pageSize=10&identify=doctor" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [admin_handler.go](file://backed/internal/api/handler/admin_handler.go#L86-L117)
- [routes.go](file://backed/internal/api/routes.go#L90)

### 禁用/启用用户
管理员禁用或启用用户账户。

**接口信息**
- **HTTP方法**: PUT
- **URL**: `/api/user/admin/status`
- **认证要求**: 需要认证（管理员）

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "userId": 1001,
  "status": 1
}
```

**字段说明**
| 字段 | 类型 | 必填 | 说明 | 取值范围 |
|------|------|------|------|--------|
| userId | int64 | 是 | 用户ID | |
| status | int | 是 | 状态 | 0:启用, 1:禁用 |

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": null
}
```

**错误码**
- 400: 参数错误
- 403: 无权限访问
- 500: 服务器内部错误

**示例请求**
```bash
curl -X PUT http://localhost:3000/api/user/admin/status \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1001,
    "status": 1
  }'
```

**来源**
- [admin_handler.go](file://backed/internal/api/handler/admin_handler.go#L120-L143)
- [routes.go](file://backed/internal/api/routes.go#L91)

### 获取登录日志
获取用户的登录日志记录。

**接口信息**
- **HTTP方法**: GET
- **URL**: `/api/user/admin/login-logs`
- **认证要求**: 需要认证（管理员）

**请求头**
```
Authorization: Bearer <token>
```

**查询参数**
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 | 1 |
| pageSize | int | 否 | 每页数量 | 10 |
| userId | int64 | 否 | 用户ID筛选 |
| status | int | 否 | 登录状态 | 0:成功, 1:失败 |
| startTime | string | 否 | 开始时间 | YYYY-MM-DD HH:MM:SS |
| endTime | string | 否 | 结束时间 | YYYY-MM-DD HH:MM:SS |

**成功响应**
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "total": 1,
    "page": 1,
    "pageSize": 10,
    "list": [
      {
        "logId": 5001,
        "userId": 1001,
        "username": "zhangsan",
        "loginIP": "192.168.1.100",
        "userAgent": "Mozilla/5.0...",
        "status": 0,
        "message": "登录成功",
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

**错误码**
- 403: 无权限访问
- 500: 服务器内部错误

**示例请求**
```bash
curl -X GET "http://localhost:3000/api/user/admin/login-logs?page=1&pageSize=10&status=0" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**来源**
- [admin_handler.go](file://backed/internal/api/handler/admin_handler.go#L146-L184)
- [routes.go](file://backed/internal/api/routes.go#L92)

## 错误码说明

### 通用错误码
| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 200 | 操作成功 | 正常处理 |
| 400 | 参数错误 | 检查请求参数是否符合要求 |
| 401 | 未授权 | 检查Authorization头或重新登录 |
| 403 | 禁止访问 | 检查用户权限是否足够 |
| 404 | 资源不存在 | 检查URL或资源ID是否正确 |
| 500 | 服务器内部错误 | 联系管理员 |

### 用户模块错误码
| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 401 | 用户名或密码错误 | 检查用户名密码是否正确 |
| 403 | 账号已被禁用 | 联系管理员启用账号 |
| 404 | 用户不存在 | 检查用户ID是否正确 |

### 管理员模块错误码
| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 403 | 无权限访问 | 确认当前用户为管理员角色 |

### 业务逻辑错误码
| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 400 | 两次密码不一致 | 确保新密码和确认密码一致 |
| 400 | 拒绝原因不能为空 | 填写拒绝原因后再提交 |
| 404 | 未找到申请记录 | 确认是否已提交申请 |

**错误码来源**
- [response.go](file://backed/pkg/utils/response.go#L45-L68)
- [开发文档.md](file://开发文档.md#L277-L284)