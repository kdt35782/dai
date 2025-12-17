# 基于国密加密的网上看诊系统 - 后端

## 项目简介

本项目是基于国密算法（SM2/SM3/SM4）的网上问诊系统后端，使用Golang开发，实现了用户注册登录、医生认证、在线问诊、电子病历等核心功能，所有敏感数据均采用国密算法加密存储和传输。

## 项目结构

```
backed/
├── cmd/
│   └── main.go              # 应用入口
├── internal/
│   ├── api/
│   │   ├── handler/         # API处理器
│   │   │   ├── user_handler.go
│   │   │   ├── admin_handler.go
│   │   │   ├── consultation_handler.go
│   │   │   ├── record_handler.go
│   │   │   ├── notification_handler.go
│   │   │   ├── file_handler.go
│   │   │   └── key_handler.go
│   │   └── routes.go        # 路由注册
│   ├── service/             # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── admin_service.go
│   │   ├── consultation_service.go
│   │   ├── record_service.go
│   │   └── notification_service.go
│   ├── repository/          # 数据访问层
│   │   ├── user_repository.go
│   │   ├── doctor_application_repository.go
│   │   ├── login_log_repository.go
│   │   ├── consultation_repository.go
│   │   ├── record_repository.go
│   │   └── notification_repository.go
│   ├── model/               # 数据模型
│   │   └── models.go
│   ├── middleware/          # 中间件
│   │   └── auth.go
│   └── crypto/              # 国密加密工具
│       └── crypto.go
├── pkg/
│   ├── config/              # 配置加载
│   │   └── config.go
│   ├── database/            # 数据库连接
│   │   └── database.go
│   └── utils/               # 工具函数
│       ├── jwt.go
│       └── response.go
├── config/
│   └── config.yaml          # 配置文件
├── init.sql                 # 数据库初始化SQL
├── go.mod                   # 依赖管理
└── README.md                # 项目说明
```

## 技术栈

- **Web框架**: Gin v1.9.1
- **ORM**: GORM v1.25.5
- **数据库**: MySQL 8.0
- **缓存**: Redis (go-redis/v9)
- **国密算法**: tjfoc/gmsm v1.4.1 (SM2/SM3/SM4)
- **JWT**: golang-jwt/jwt/v5 v5.2.0
- **配置管理**: Viper v1.18.2

## 核心功能

### 1. 用户模块
- ✅ 用户注册/登录（SM3密码加密）
- ✅ 用户信息管理（SM4敏感信息加密）
- ✅ 密码修改
- ✅ 申请成为医生
- ✅ 医生列表查询

### 2. 管理员模块
- ✅ 审核医生申请
- ✅ 获取医生申请列表
- ✅ 用户管理（查询、禁用/启用）
- ✅ 登录日志查询

### 3. 问诊模块
- ✅ 创建问诊（支持AI辅助诊断）
- ✅ 问诊列表查询
- ✅ 问诊详情查看
- ✅ 医生接诊
- ✅ 完成问诊

### 4. 病历模块
- ✅ 病历列表查询
- ✅ 病历详情查看
- ✅ 数据完整性验证（SM3哈希）

### 5. 通知模块
- ✅ 通知列表
- ✅ 未读数量统计
- ✅ 标记已读

### 6. 国密加密
- ✅ SM3密码哈希加密
- ✅ SM4对称加密（邮箱、手机号、身份证等）
- ✅ SM2非对称加密（密钥管理）

## 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Redis (可选)

### 安装步骤

#### 1. 克隆项目

```bash
cd backed
```

#### 2. 安装依赖

```bash
go mod tidy
```

#### 3. 配置数据库

编辑 `config/config.yaml` 文件，修改数据库连接信息：

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password  # 修改为你的数据库密码
  database: SM
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100
```

#### 4. 初始化数据库

执行 `init.sql` 文件创建数据库和表：

```bash
mysql -u root -p < init.sql
```

或在MySQL客户端中执行：

```sql
source /path/to/backed/init.sql
```

#### 5. 修改加密密钥（重要！）

编辑 `config/config.yaml`，修改SM4加密密钥：

```yaml
crypto:
  sm4_key: 0123456789abcdef0123456789abcdef  # 32位16进制字符串，请务必修改
```

⚠️ **安全提示**：生产环境中必须使用随机生成的密钥，不要使用示例密钥！

#### 6. 运行项目

```bash
go run cmd/main.go
```

服务器默认运行在 `http://localhost:3000`

## 配置说明

### config.yaml 完整配置

```yaml
server:
  port: 3000           # 服务端口
  mode: debug          # debug | release

database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: SM
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-jwt-secret-key-change-in-production  # JWT密钥，请修改
  expires_in: 7200  # Token过期时间（秒）

crypto:
  sm4_key: 0123456789abcdef0123456789abcdef  # SM4密钥，请修改
  sm2_private_key: ""  # 自动生成
  sm2_public_key: ""   # 自动生成

upload:
  max_size: 10485760  # 文件上传最大10MB
  allowed_types:
    - image/jpeg
    - image/png
    - image/jpg
  upload_path: ./uploads
```

## API文档

详细的API文档请参考项目根目录下的 `开发文档.md`

### 主要API端点

**用户模块**
- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录
- `GET /api/user/info` - 获取用户信息
- `PUT /api/user/profile` - 更新用户信息
- `PUT /api/user/password` - 修改密码
- `POST /api/user/apply-doctor` - 申请成为医生
- `GET /api/user/doctors` - 获取医生列表

**管理员模块**
- `PUT /api/user/admin/review-doctor` - 审核医生申请
- `GET /api/user/admin/doctor-applications` - 获取医生申请列表
- `GET /api/user/admin/users` - 获取用户列表
- `PUT /api/user/admin/status` - 禁用/启用用户
- `GET /api/user/admin/login-logs` - 获取登录日志

**问诊模块**
- `POST /api/consultation/create` - 创建问诊
- `GET /api/consultation/list` - 获取问诊列表
- `GET /api/consultation/detail` - 获取问诊详情
- `POST /api/consultation/accept` - 医生接诊
- `POST /api/consultation/finish` - 完成问诊

**病历模块**
- `GET /api/record/list` - 获取病历列表
- `GET /api/record/detail` - 获取病历详情

**通知模块**
- `GET /api/notification/list` - 获取通知列表
- `GET /api/notification/unread-count` - 获取未读数量
- `PUT /api/notification/mark-read` - 标记已读

**文件上传**
- `POST /api/file/upload` - 文件上传

**密钥管理**
- `POST /api/key/generate` - 生成密钥对

## 国密加密说明

### SM3哈希算法
- 用于密码加密（双重哈希：前端+后端加盐）
- 用于数据完整性验证

### SM4对称加密
- 邮箱、手机号、身份证号
- 真实姓名
- 登录IP
- 主诉、诊断、处方等医疗数据

### SM2非对称加密
- 用于密钥交换
- 用于敏感数据传输

## 默认账号

系统初始化后会创建一个默认管理员账号：

- 用户名：`admin`
- 密码：`Admin123!@#`

⚠️ **请在首次登录后立即修改密码！**

## 开发指南

### 添加新的API接口

1. 在 `internal/model/models.go` 中定义数据模型
2. 在 `internal/repository/` 中创建数据访问层
3. 在 `internal/service/` 中实现业务逻辑
4. 在 `internal/api/handler/` 中创建处理器
5. 在 `internal/api/routes.go` 中注册路由

### 代码规范

- 遵循Go官方代码规范
- 使用有意义的变量和函数名
- 添加必要的注释
- 敏感数据必须加密存储

## 常见问题

### 1. 数据库连接失败

检查 `config/config.yaml` 中的数据库配置是否正确，确保MySQL服务已启动。

### 2. 编译错误

执行 `go mod tidy` 更新依赖。

### 3. Token过期

检查JWT配置中的 `expires_in` 设置，默认为2小时（7200秒）。

### 4. 文件上传失败

确保 `uploads` 目录存在且有写入权限：

```bash
mkdir uploads
chmod 755 uploads
```

## 生产环境部署

### 1. 修改配置

```yaml
server:
  mode: release  # 切换到生产模式
```

### 2. 编译

```bash
go build -o sm-medical cmd/main.go
```

### 3. 运行

```bash
./sm-medical
```

### 4. 使用systemd管理（Linux）

创建服务文件 `/etc/systemd/system/sm-medical.service`：

```ini
[Unit]
Description=SM Medical Backend Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/backed
ExecStart=/path/to/backed/sm-medical
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl start sm-medical
sudo systemctl enable sm-medical
```

## 性能优化建议

1. 启用Redis缓存
2. 使用连接池
3. 添加数据库索引
4. 启用GZIP压缩
5. 使用CDN加速静态资源

## 安全建议

1. ✅ 使用HTTPS（建议配置国密SSL证书）
2. ✅ 定期更换加密密钥
3. ✅ 启用防火墙
4. ✅ 限制API访问频率
5. ✅ 定期备份数据库
6. ✅ 使用强密码策略
7. ✅ 及时更新依赖包

## 许可证

MIT License

## 联系方式

如有问题，请提交Issue或联系项目维护者。
