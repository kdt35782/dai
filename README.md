# 基于国密加密的网上看诊系统

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-2.x-4FC08D?logo=vue.js)](https://vuejs.org/)
[![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white)](https://www.mysql.com/)

> 一个采用国密算法(SM2/SM3/SM4)保障数据安全的在线医疗问诊平台，集成智能分诊、AI辅助诊断、实时聊天等功能。

## 📋 项目简介

本系统是一个完整的在线医疗问诊解决方案，采用前后端分离架构开发。核心特色是全面应用国产密码算法(国密算法)，确保患者隐私数据和医疗信息的安全性，符合国家网络安全法及数据安全法要求。

### 核心亮点

- 🔐 **国密加密体系**: 全链路采用SM2/SM3/SM4算法保护数据安全
- 🤖 **AI智能诊断**: 基于机器学习的辅助诊断系统
- 🏥 **智能分诊**: 自动匹配最优医生，支持负载均衡
- 💬 **实时通讯**: WebSocket实现医患即时沟通
- 📋 **电子处方**: 标准化处方管理与药品推荐
- 📊 **数据可视化**: 医生工作负载、患者健康统计

## 🏗️ 系统架构

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│  前端应用   │─────▶│  后端API    │─────▶│  数据库     │
│  (uniapp)   │      │  (Gin)      │      │  (MySQL)    │
└─────────────┘      └─────────────┘      └─────────────┘
       │                    │                     │
       │                    ▼                     │
       │            ┌─────────────┐               │
       └───────────▶│  国密加密   │◀──────────────┘
                    │  SM2/3/4    │
                    └─────────────┘
```

### 技术栈

**后端 (backed/)**
- **语言**: Go 1.21+
- **框架**: Gin v1.9.1
- **ORM**: GORM v1.25.5
- **数据库**: MySQL 8.0
- **缓存**: Redis
- **国密**: tjfoc/gmsm v1.4.1
- **认证**: JWT (golang-jwt/jwt v5)

**前端 (fonted/)**
- **框架**: uniapp + Vue 2.x
- **UI**: uni-ui 组件库
- **国密**: sm-crypto
- **通信**: WebSocket + HTTP

**AI训练 (ml_training/)**
- **语言**: Python 3.8+
- **框架**: scikit-learn, pandas

## 🚀 快速开始

### 环境要求

- Go 1.21 或更高版本
- MySQL 8.0+
- Node.js 14+ (前端开发)
- HBuilderX (uniapp开发工具，可选)

### 1. 克隆项目

```bash
git clone https://github.com/yourusername/sm-medical-system.git
cd sm-medical-system
```

### 2. 后端部署

#### 安装依赖
```bash
cd backed
go mod tidy
```

#### 配置数据库
编辑 `backed/config/config.yaml`:
```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password  # 修改为你的密码
  database: SM
```

#### 初始化数据库
```bash
# 方式1: 使用init.sql
mysql -u root -p < backed/init.sql

# 方式2: 执行所有SQL文件
mysql -u root -p SM < database/SM_user_module.sql
mysql -u root -p SM < database/chat_message.sql
mysql -u root -p SM < database/medicine_database.sql
mysql -u root -p SM < database/triage_enhancement.sql
```

#### 启动后端服务
```bash
cd backed
go run cmd/main.go
```

服务运行在 `http://localhost:3000`

### 3. 前端部署

#### 安装依赖
```bash
cd fonted
npm install
```

#### 配置API地址
编辑 `fonted/utils/config.js`，修改 `BASE_URL`:
```javascript
export const BASE_URL = 'http://localhost:3000'
```

#### 运行项目

**使用HBuilderX**:
1. 导入 `fonted` 目录
2. 点击"运行" → "运行到浏览器" → "Chrome"

**使用命令行**:
```bash
npm run dev:h5
```

访问 `http://localhost:8080`

## 📚 功能模块

### 用户模块
- ✅ 用户注册/登录 (支持患者、医生双角色)
- ✅ 个人信息管理
- ✅ 密码修改 (SM3双重哈希)
- ✅ 医生认证申请

### 问诊模块
- ✅ 创建问诊 (支持选择医生或智能分诊)
- ✅ AI智能诊断 (基于症状分析)
- ✅ 医生接诊
- ✅ 实时聊天 (WebSocket)
- ✅ 完成问诊 (诊断+处方一体化)

### 智能分诊
- ✅ 4级规则优先级匹配
- ✅ 科室智能匹配
- ✅ 医生负载均衡
- ✅ 在线状态优先

### 处方管理
- ✅ 电子处方开具
- ✅ 药品搜索与推荐
- ✅ 用药信息管理
- ✅ 处方详情查看

### 病历管理
- ✅ 电子病历自动生成
- ✅ 数据完整性验证 (SM3哈希)
- ✅ 病历列表查询
- ✅ 病历详情查看

### 通知系统
- ✅ 系统通知推送
- ✅ 未读消息统计
- ✅ 消息已读标记

### 管理后台
- ✅ 用户管理 (查询、禁用/启用)
- ✅ 登录日志查询
- ✅ 数据统计分析

## 🔐 国密加密应用

### SM3 哈希算法
- **密码加密**: 前端+后端双重SM3哈希加盐
- **数据完整性**: 病历数据防篡改验证
- **敏感数据摘要**: 关键信息校验

### SM4 对称加密
- **用户信息**: 邮箱、手机号、身份证号、真实姓名
- **医疗数据**: 主诉、症状、诊断、处方
- **登录信息**: IP地址
- **通信数据**: 聊天消息内容

### SM2 非对称加密
- **密钥交换**: 前后端密钥协商
- **数字签名**: 敏感操作验证

## 📖 API文档

详细API文档请查看:
- [后端API总览](API接口总览.md)
- [开发文档](开发文档.md)
- [处方API文档](处方API文档.md)
- [聊天功能实现](CHAT_IMPLEMENTATION.md)

### 核心接口示例

```bash
# 用户注册
POST /api/user/register

# 用户登录
POST /api/user/login

# 创建问诊
POST /api/consultation/create

# 智能分诊
POST /api/triage/auto-assign

# 开具处方
POST /api/consultation/finish
```

## 📁 项目结构

```
.
├── backed/                 # 后端代码
│   ├── cmd/               # 应用入口
│   ├── internal/          # 核心业务逻辑
│   │   ├── api/          # API处理器
│   │   ├── service/      # 业务服务层
│   │   ├── repository/   # 数据访问层
│   │   ├── model/        # 数据模型
│   │   ├── crypto/       # 国密加密
│   │   └── middleware/   # 中间件
│   ├── pkg/              # 公共包
│   ├── config/           # 配置文件
│   └── init.sql          # 数据库初始化
│
├── fonted/               # 前端代码
│   ├── pages/           # 页面组件
│   │   ├── login/      # 登录注册
│   │   ├── consultation/ # 问诊模块
│   │   ├── chat/       # 聊天模块
│   │   ├── prescription/ # 处方模块
│   │   └── user/       # 用户中心
│   ├── utils/          # 工具函数
│   │   ├── request.js  # HTTP请求
│   │   ├── crypto.js   # 国密加密
│   │   └── storage.js  # 本地存储
│   └── pages.json      # 页面配置
│
├── database/            # 数据库脚本
│   ├── SM_user_module.sql
│   ├── chat_message.sql
│   ├── medicine_database.sql
│   └── triage_enhancement.sql
│
├── ml_training/         # AI模型训练
│   ├── train_model.py
│   └── export_training_data.py
│
└── README.md           # 项目说明
```

## 🔧 配置说明

### 后端配置 (backed/config/config.yaml)

```yaml
server:
  port: 3000
  mode: debug  # release 生产环境

database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: SM

jwt:
  secret: your-secret-key  # ⚠️ 生产环境请修改
  expires_in: 7200

crypto:
  sm4_key: 0123456789abcdef0123456789abcdef  # ⚠️ 请修改为随机密钥
```

### 前端配置 (fonted/utils/config.js)

```javascript
export const BASE_URL = 'http://localhost:3000'
export const WS_URL = 'ws://localhost:3000'
```

## 🧪 测试

### 默认账户

**管理员**:
- 用户名: `admin`
- 密码: `Admin123!@#`

**测试医生**:
- 用户名: `dr_zhangwei`
- 密码: `123456`

**测试患者**:
- 用户名: `patient_test`
- 密码: `123456`

⚠️ **首次登录后请立即修改密码!**

### 测试工具

- [API测试工具](API测试工具.html) - 在线接口测试
- [注册测试工具](测试注册工具.html) - 注册功能测试
- [SM3加密测试](测试SM3加密.html) - 国密加密验证

## 📊 数据库设计

### 核心表结构

- **SM_user**: 用户表 (支持患者/医生/管理员)
- **SM_consultation**: 问诊记录表
- **SM_medical_record**: 电子病历表
- **SM_chat_message**: 聊天消息表
- **SM_prescription**: 处方表
- **SM_medicine**: 药品库
- **SM_notification**: 通知表
- **SM_login_log**: 登录日志表

详见 [开发文档](开发文档.md)

## 🚢 部署指南

### 生产环境部署

1. **修改配置**
   ```yaml
   server:
     mode: release
   ```

2. **编译后端**
   ```bash
   cd backed
   go build -o sm-medical cmd/main.go
   ```

3. **前端打包**
   ```bash
   cd fonted
   npm run build:h5
   ```

4. **使用Nginx部署前端**
   ```nginx
   server {
       listen 80;
       server_name your-domain.com;
       
       location / {
           root /path/to/fonted/dist/build/h5;
           index index.html;
       }
       
       location /api {
           proxy_pass http://localhost:3000;
       }
   }
   ```

### Docker部署 (可选)

```bash
# 构建后端镜像
cd backed
docker build -t sm-medical-backend .

# 运行容器
docker run -d -p 3000:3000 sm-medical-backend
```

## 🔒 安全建议

1. ✅ 修改默认密钥 (SM4_KEY, JWT_SECRET)
2. ✅ 使用HTTPS协议 (建议国密SSL证书)
3. ✅ 启用防火墙限制数据库访问
4. ✅ 定期备份数据库
5. ✅ 实施API访问频率限制
6. ✅ 定期更新依赖包
7. ✅ 使用强密码策略

## 🤝 贡献指南

欢迎提交Issue和Pull Request!

1. Fork本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 📄 开源协议

本项目采用 [MIT](LICENSE) 协议开源

## 📧 联系方式

- 项目主页: https://github.com/yourusername/sm-medical-system
- 问题反馈: https://github.com/yourusername/sm-medical-system/issues
- 邮箱: your.email@example.com

## 🙏 致谢

- [tjfoc/gmsm](https://github.com/tjfoc/gmsm) - 国密算法Go语言实现
- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [uniapp](https://uniapp.dcloud.io/) - 跨平台应用框架
- [GORM](https://gorm.io/) - Go ORM框架

## 📝 更新日志

### v1.0.0 (2025-12-17)
- ✅ 完成基础用户系统
- ✅ 实现国密加密体系
- ✅ 开发问诊核心功能
- ✅ 集成AI辅助诊断
- ✅ 实现智能分诊系统
- ✅ 完善处方管理模块
- ✅ 实时聊天功能上线

---

**⭐ 如果这个项目对你有帮助，请给个Star支持一下！**
