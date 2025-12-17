# 🤖 AI诊断机器学习集成方案

## 📋 方案概述

本方案通过机器学习算法训练AI诊断模型，目标准确率**85%+**，提升诊断精准度。

### 核心特性

- ✅ **自动数据收集**：问诊完成后自动收集医生确诊数据
- ✅ **多模型对比**：逻辑回归、随机森林、梯度提升树
- ✅ **自动调参**：网格搜索最优超参数
- ✅ **版本管理**：模型版本追踪和性能监控
- ✅ **在线预测**：训练好的模型可直接集成到Go后端

---

## 🗂️ 文件结构

```
ml_training/
├── train_model.py           # 模型训练脚本（主文件）
├── export_training_data.py  # 数据导出和质量检查
├── requirements.txt         # Python依赖
├── README.md                # 本文档
├── models/                  # 模型文件存储目录
│   └── random_forest_v20250101_120000.pkl
└── logs/                    # 训练日志

database/
└── ml_training_data.sql     # 数据库表结构
```

---

## 🚀 快速开始

### 1. 环境准备

#### 安装Python环境
```bash
# Windows
python --version  # 确认Python 3.8+

# 安装依赖
cd ml_training
pip install -r requirements.txt
```

#### 导入数据库表
```bash
# 在MySQL中执行
mysql -u root -p sm_medical < ../database/ml_training_data.sql
```

### 2. 配置数据库

编辑 `train_model.py` 和 `export_training_data.py`：

```python
DB_CONFIG = {
    'host': 'localhost',
    'port': 3306,
    'user': 'root',
    'password': 'your_actual_password',  # 修改这里
    'database': 'sm_medical',
    'charset': 'utf8mb4'
}
```

### 3. 收集训练数据

**方式一：自动收集（推荐）**

系统已配置触发器，问诊完成后自动收集数据到 `SM_ai_training_data` 表。

**方式二：手动插入测试数据**

```sql
-- 插入示例训练数据
INSERT INTO SM_ai_training_data (
    consultation_id, patient_id,
    age, gender, systolic_bp, diastolic_bp, heart_rate,
    temperature, blood_sugar,
    doctor_diagnosis, diagnosis_icd10,
    data_quality, is_verified
) VALUES 
(1001, 1001, 65, 1, 165, 105, 88, 36.8, 6.5, '高血压2级', 'I10', 1, 1),
(1002, 1002, 45, 0, 155, 95, 92, 37.2, 7.5, '高血压1级', 'I10', 1, 1),
(1003, 1003, 55, 1, 180, 110, 105, 36.5, 8.2, '高血压3级', 'I10', 1, 1),
-- ... 更多数据
```

### 4. 检查数据质量

```bash
python export_training_data.py
```

输出示例：
```
📊 数据质量检查
============================================================

1. 总样本数: 150
2. 已验证样本: 120 (80.0%)

3. 数据质量分布:
   高质量: 80
   中等: 40
   
4. 疾病类别数: 8

5. TOP 10 疾病:
   1. 高血压2级: 35 条
   2. 高血压1级: 28 条
   3. 糖尿病: 20 条
   ...

6. 数据完整性:
   年龄: 95.0%
   血压: 92.0%
   心率: 88.0%
   体温: 75.0%
   血糖: 70.0%

============================================================
✅ 建议: 数据量基本满足，可以开始训练
```

### 5. 训练模型

```bash
python train_model.py
```

训练过程：
```
============================================================
🚀 开始训练AI诊断模型
============================================================

📊 加载数据: 120 条记录
✅ 有效数据: 110 条
📋 疾病类别: 8 个
📈 类别分布:
高血压      35
糖尿病      20
感染        18
...

📊 数据划分:
  训练集: 88 条
  测试集: 22 条

🔧 训练逻辑回归模型...
✅ 最佳参数: {'C': 1, 'penalty': 'l2', 'solver': 'lbfgs'}
✅ 交叉验证准确率: 0.8250

📊 模型评估...
✅ 准确率 (Accuracy): 0.8636
✅ 精确率 (Precision): 0.8542
✅ 召回率 (Recall): 0.8636
✅ F1分数: 0.8588
✅ AUC: 0.9245

🌲 训练随机森林模型...
✅ 最佳参数: {'max_depth': 10, 'n_estimators': 100}
✅ 交叉验证准确率: 0.8750

📊 模型评估...
✅ 准确率 (Accuracy): 0.9091
✅ 精确率 (Precision): 0.9125
✅ 召回率 (Recall): 0.9091
✅ F1分数: 0.9078
✅ AUC: 0.9580

============================================================
🏆 最佳模型: random_forest
🎯 准确率: 0.9091
============================================================

✅ 模型已保存: ./models/random_forest_v20250116_153045.pkl
✅ 模型信息已保存到数据库

📊 特征重要性 TOP 10:
            feature  importance
         systolic_bp      0.2450
        diastolic_bp      0.1820
                 age      0.1350
          heart_rate      0.1120
         blood_sugar      0.0980
...

✅ 训练完成！

🎉 恭喜！模型准确率达到 90.91%，超过85%目标！
```

---

## 📊 数据库表说明

### 1. SM_ai_training_data（训练数据表）

存储用于模型训练的真实诊断数据。

**关键字段**：
- 输入特征：年龄、性别、血压、心率、体温、血糖等
- 医生标签：doctor_diagnosis（真实诊断结果）
- 数据质量：data_quality（1高/2中/3低）
- 验证状态：is_verified（是否医生复核）

### 2. SM_ai_model_version（模型版本表）

记录每个训练模型的版本和性能指标。

**关键字段**：
- model_name: 模型名称
- version: 版本号（如v20250116_153045）
- accuracy/precision/recall/f1_score: 性能指标
- model_file_path: 模型文件路径
- is_deployed: 是否部署

### 3. SM_ai_prediction_log（预测日志表）

记录每次AI预测的结果，用于后续模型评估。

### 4. SM_disease_mapping（疾病映射表）

标准化疾病名称和ICD-10编码。

---

## 🎯 训练策略

### 数据要求

| 数据量 | 训练建议 | 预期准确率 |
|-------|---------|-----------|
| < 50条 | 不建议训练，继续收集 | < 70% |
| 50-100条 | 可以训练，但效果一般 | 70-80% |
| 100-200条 | 推荐训练 | 80-85% |
| 200+条 | 优质训练 | **85-90%+** |

### 模型选择

1. **逻辑回归**：快速、可解释性强，适合小数据集
2. **随机森林**：鲁棒性好，特征重要性分析
3. **梯度提升树**：准确率最高，但可能过拟合

### 特征工程

系统自动提取以下特征：

**生理指标**（8个）：
- age, gender, systolic_bp, diastolic_bp, heart_rate, temperature, blood_sugar, bmi

**派生特征**（3个）：
- pulse_pressure（脉压差）
- mean_arterial_pressure（平均动脉压）
- age_group（年龄分组）

**症状特征**（5个）：
- has_headache（头痛）
- has_fever（发热）
- has_chest_pain（胸痛）
- has_cough（咳嗽）
- has_abdominal_pain（腹痛）

**既往史**（5个）：
- has_hypertension, has_diabetes, has_heart_disease, smoking_status, drinking_status

**总计**：21个特征

---

## 🔧 模型部署

### 查看训练好的模型

```sql
-- 查询所有模型版本
SELECT 
    id, model_name, version, model_type,
    accuracy, f1_score,
    is_deployed, created_at
FROM SM_ai_model_version
ORDER BY accuracy DESC;

-- 查询最佳模型
SELECT * FROM SM_ai_model_version
WHERE is_deployed = 0
ORDER BY accuracy DESC
LIMIT 1;
```

### 部署模型（待实现）

模型训练完成后，可通过以下方式集成到Go后端：

**方式一：HTTP API（推荐）**
```python
# 启动预测服务
python prediction_service.py --port 5000
```

**方式二：导出PMML格式**
```python
# 导出为PMML（跨语言模型格式）
from sklearn2pmml import sklearn2pmml
sklearn2pmml(model, "model.pmml")
```

**方式三：Go调用Python**
```go
// Go中通过exec调用Python脚本
cmd := exec.Command("python", "predict.py", "--input", jsonData)
output, _ := cmd.Output()
```

---

## 📈 性能监控

### 查看模型性能

```sql
-- 模型预测准确率统计
SELECT 
    mv.model_name,
    mv.version,
    COUNT(*) as total_predictions,
    SUM(CASE WHEN pl.is_correct = 1 THEN 1 ELSE 0 END) as correct_predictions,
    ROUND(SUM(CASE WHEN pl.is_correct = 1 THEN 1 ELSE 0 END) / COUNT(*) * 100, 2) as accuracy_percent
FROM SM_ai_prediction_log pl
JOIN SM_ai_model_version mv ON pl.model_version_id = mv.id
WHERE pl.is_correct IS NOT NULL
GROUP BY mv.id
ORDER BY accuracy_percent DESC;
```

### 模型对比

```sql
-- 对比不同模型版本
SELECT 
    model_type,
    COUNT(*) as versions,
    AVG(accuracy) as avg_accuracy,
    MAX(accuracy) as best_accuracy
FROM SM_ai_model_version
GROUP BY model_type
ORDER BY avg_accuracy DESC;
```

---

## 🔍 故障排查

### Q1: 数据量不足

**问题**：训练数据少于50条

**解决**：
```sql
-- 检查已完成问诊数
SELECT COUNT(*) FROM SM_consultation WHERE status = 2;

-- 检查是否有医生诊断
SELECT COUNT(*) FROM SM_consultation 
WHERE status = 2 AND doctor_diagnosis IS NOT NULL;

-- 手动触发数据收集
INSERT INTO SM_ai_training_data (...)
SELECT ... FROM SM_consultation WHERE ...;
```

### Q2: 模型准确率低

**问题**：准确率<85%

**优化方向**：
1. **增加数据量**：收集更多样本
2. **特征工程**：添加更多症状特征
3. **数据清洗**：移除低质量数据
4. **参数调优**：调整模型超参数

### Q3: 依赖安装失败

```bash
# 使用国内镜像
pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 或单独安装
pip install scikit-learn==1.3.0 -i https://pypi.tuna.tsinghua.edu.cn/simple
```

---

## 📝 最佳实践

### 1. 数据收集阶段

- ✅ 确保医生填写完整的诊断信息
- ✅ 定期审查数据质量（data_quality字段）
- ✅ 标准化疾病名称（使用ICD-10编码）

### 2. 模型训练阶段

- ✅ 数据量达到100+再训练
- ✅ 保留10-20%数据作为测试集
- ✅ 训练多个模型进行对比
- ✅ 记录训练日志和性能指标

### 3. 模型部署阶段

- ✅ 先在测试环境验证
- ✅ 逐步灰度发布（10% → 50% → 100%）
- ✅ 监控线上预测准确率
- ✅ 定期重新训练（如每月一次）

---

## 🎯 目标达成路径

### 阶段一：数据收集（1-2个月）

**目标**：收集200+条高质量诊断数据

**行动**：
- 每天完成5-10个问诊
- 医生必须填写完整诊断
- 定期检查数据质量

### 阶段二：模型训练（1周）

**目标**：训练出准确率85%+的模型

**行动**：
- 运行 `python export_training_data.py` 检查数据
- 运行 `python train_model.py` 训练模型
- 分析特征重要性，优化特征

### 阶段三：模型集成（2周）

**目标**：将模型集成到Go后端

**行动**：
- 搭建Python预测API服务
- Go后端调用Python API
- 前端展示ML预测结果

### 阶段四：持续优化（长期）

**目标**：准确率提升到90%+

**行动**：
- 每月重新训练模型
- 收集错误预测案例
- 增加更多特征维度

---

## ✅ 总结

通过本方案，您可以：

1. ✅ **自动收集**真实诊断数据
2. ✅ **训练高精度**机器学习模型（目标85%+）
3. ✅ **版本管理**模型性能和迭代
4. ✅ **监控评估**线上预测效果
5. ✅ **持续优化**提升诊断准确率

**预期效果**：

| 指标 | 当前（规则引擎） | 目标（机器学习） |
|------|---------------|----------------|
| 准确率 | 65-75% | **85-90%+** |
| 误报率 | 20-25% | **<10%** |
| 疾病覆盖 | 15种 | **30+种** |
| 响应时间 | <500ms | <200ms |

---

**开始训练您的AI医生吧！** 🚀
