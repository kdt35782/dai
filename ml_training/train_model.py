#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
AI诊断模型训练脚本
支持逻辑回归、随机森林、XGBoost等多种模型
目标：准确率达到85%+
"""

import os
import json
import pickle
import numpy as np
import pandas as pd
from datetime import datetime
from sklearn.model_selection import train_test_split, GridSearchCV, cross_val_score
from sklearn.preprocessing import StandardScaler, LabelEncoder
from sklearn.linear_model import LogisticRegression
from sklearn.ensemble import RandomForestClassifier, GradientBoostingClassifier
from sklearn.metrics import (
    accuracy_score, precision_score, recall_score, f1_score,
    roc_auc_score, confusion_matrix, classification_report
)
import pymysql
import warnings
warnings.filterwarnings('ignore')

# ================== 配置 ==================
DB_CONFIG = {
    'host': 'localhost',
    'port': 3306,
    'user': 'root',
    'password': 'your_password',  # 修改为实际密码
    'database': 'sm_medical',
    'charset': 'utf8mb4'
}

MODEL_DIR = './models'
LOG_DIR = './logs'
os.makedirs(MODEL_DIR, exist_ok=True)
os.makedirs(LOG_DIR, exist_ok=True)

# 疾病分类映射
DISEASE_MAPPING = {
    '高血压': ['高血压1级', '高血压2级', '高血压3级', '轻度高血压', '中度高血压', '重度高血压'],
    '低血压': ['低血压', '血压偏低'],
    '心律失常': ['心动过速', '心动过缓', '心率过快', '心率过慢'],
    '糖尿病': ['糖尿病', '高血糖', '空腹血糖异常'],
    '低血糖': ['低血糖', '血糖偏低'],
    '感染': ['发热', '高热', '感冒', '上呼吸道感染'],
    '消化系统': ['急性胃炎', '胃炎', '腹痛', '腹泻'],
    '神经系统': ['偏头痛', '头痛', '头晕']
}


class MedicalAITrainer:
    """医疗AI模型训练器"""
    
    def __init__(self, db_config):
        self.db_config = db_config
        self.conn = None
        self.scaler = StandardScaler()
        self.label_encoder = LabelEncoder()
        self.feature_names = []
        self.model = None
        self.model_type = None
        
    def connect_db(self):
        """连接数据库"""
        try:
            self.conn = pymysql.connect(**self.db_config)
            print("✅ 数据库连接成功")
            return True
        except Exception as e:
            print(f"❌ 数据库连接失败: {e}")
            return False
    
    def load_training_data(self, min_samples_per_class=10):
        """从数据库加载训练数据"""
        if not self.conn:
            if not self.connect_db():
                return None, None
        
        # 查询高质量训练数据
        query = """
        SELECT 
            age, gender, systolic_bp, diastolic_bp, heart_rate,
            temperature, blood_sugar, bmi,
            symptom_keywords, symptom_severity,
            has_hypertension, has_diabetes, has_heart_disease,
            smoking_status, drinking_status,
            doctor_diagnosis, diagnosis_icd10
        FROM SM_ai_training_data
        WHERE is_verified = 1 
          AND data_quality IN (1, 2)
          AND doctor_diagnosis IS NOT NULL
          AND doctor_diagnosis != ''
        ORDER BY created_at DESC
        """
        
        try:
            df = pd.read_sql(query, self.conn)
            print(f"📊 加载数据: {len(df)} 条记录")
            
            if len(df) < 50:
                print(f"⚠️  数据量不足({len(df)}条)，建议至少50条，当前使用规则引擎")
                return None, None
            
            # 数据预处理
            df = self._preprocess_data(df)
            
            # 过滤小样本类别
            disease_counts = df['disease_category'].value_counts()
            valid_diseases = disease_counts[disease_counts >= min_samples_per_class].index
            df = df[df['disease_category'].isin(valid_diseases)]
            
            print(f"✅ 有效数据: {len(df)} 条")
            print(f"📋 疾病类别: {len(valid_diseases)} 个")
            print(f"📈 类别分布:\n{df['disease_category'].value_counts()}")
            
            # 分离特征和标签
            X = df.drop(['doctor_diagnosis', 'diagnosis_icd10', 'disease_category'], axis=1)
            y = df['disease_category']
            
            self.feature_names = X.columns.tolist()
            
            return X, y
            
        except Exception as e:
            print(f"❌ 数据加载失败: {e}")
            return None, None
    
    def _preprocess_data(self, df):
        """数据预处理"""
        # 1. 处理缺失值
        numeric_cols = ['age', 'systolic_bp', 'diastolic_bp', 'heart_rate', 
                       'temperature', 'blood_sugar', 'bmi', 'symptom_severity']
        for col in numeric_cols:
            if col in df.columns:
                df[col] = df[col].fillna(df[col].median())
        
        # 2. 处理症状关键词（文本特征）
        if 'symptom_keywords' in df.columns:
            # 提取关键症状特征
            df['has_headache'] = df['symptom_keywords'].apply(
                lambda x: 1 if isinstance(x, str) and ('头痛' in x or '头晕' in x) else 0
            )
            df['has_fever'] = df['symptom_keywords'].apply(
                lambda x: 1 if isinstance(x, str) and ('发热' in x or '发烧' in x) else 0
            )
            df['has_chest_pain'] = df['symptom_keywords'].apply(
                lambda x: 1 if isinstance(x, str) and ('胸痛' in x or '胸闷' in x) else 0
            )
            df['has_cough'] = df['symptom_keywords'].apply(
                lambda x: 1 if isinstance(x, str) and '咳嗽' in x else 0
            )
            df['has_abdominal_pain'] = df['symptom_keywords'].apply(
                lambda x: 1 if isinstance(x, str) and '腹痛' in x else 0
            )
            df = df.drop('symptom_keywords', axis=1)
        
        # 3. 标准化疾病诊断（映射到大类）
        df['disease_category'] = df['doctor_diagnosis'].apply(self._map_disease_category)
        
        # 4. 处理布尔值
        bool_cols = ['has_hypertension', 'has_diabetes', 'has_heart_disease']
        for col in bool_cols:
            if col in df.columns:
                df[col] = df[col].fillna(0).astype(int)
        
        # 5. 处理分类变量
        if 'gender' in df.columns:
            df['gender'] = df['gender'].fillna(0).astype(int)
        if 'smoking_status' in df.columns:
            df['smoking_status'] = df['smoking_status'].fillna(0).astype(int)
        if 'drinking_status' in df.columns:
            df['drinking_status'] = df['drinking_status'].fillna(0).astype(int)
        if 'symptom_severity' in df.columns:
            df['symptom_severity'] = df['symptom_severity'].fillna(5).astype(int)
        
        # 6. 计算派生特征
        if 'systolic_bp' in df.columns and 'diastolic_bp' in df.columns:
            df['pulse_pressure'] = df['systolic_bp'] - df['diastolic_bp']  # 脉压差
            df['mean_arterial_pressure'] = (df['systolic_bp'] + 2 * df['diastolic_bp']) / 3  # 平均动脉压
        
        if 'age' in df.columns:
            df['age_group'] = pd.cut(df['age'], bins=[0, 18, 40, 60, 100], 
                                    labels=[0, 1, 2, 3]).astype(int)
        
        return df
    
    def _map_disease_category(self, diagnosis):
        """将具体疾病映射到大类"""
        if not isinstance(diagnosis, str):
            return '其他'
        
        for category, keywords in DISEASE_MAPPING.items():
            for keyword in keywords:
                if keyword in diagnosis:
                    return category
        return '其他'
    
    def train_logistic_regression(self, X_train, y_train):
        """训练逻辑回归模型"""
        print("\n🔧 训练逻辑回归模型...")
        
        param_grid = {
            'C': [0.001, 0.01, 0.1, 1, 10, 100],
            'penalty': ['l2'],
            'solver': ['lbfgs', 'saga'],
            'max_iter': [1000]
        }
        
        lr = LogisticRegression(random_state=42, multi_class='multinomial')
        grid_search = GridSearchCV(lr, param_grid, cv=5, scoring='accuracy', n_jobs=-1)
        grid_search.fit(X_train, y_train)
        
        print(f"✅ 最佳参数: {grid_search.best_params_}")
        print(f"✅ 交叉验证准确率: {grid_search.best_score_:.4f}")
        
        return grid_search.best_estimator_
    
    def train_random_forest(self, X_train, y_train):
        """训练随机森林模型"""
        print("\n🌲 训练随机森林模型...")
        
        param_grid = {
            'n_estimators': [50, 100, 200],
            'max_depth': [5, 10, 15, None],
            'min_samples_split': [2, 5, 10],
            'min_samples_leaf': [1, 2, 4]
        }
        
        rf = RandomForestClassifier(random_state=42, n_jobs=-1)
        grid_search = GridSearchCV(rf, param_grid, cv=5, scoring='accuracy', n_jobs=-1)
        grid_search.fit(X_train, y_train)
        
        print(f"✅ 最佳参数: {grid_search.best_params_}")
        print(f"✅ 交叉验证准确率: {grid_search.best_score_:.4f}")
        
        return grid_search.best_estimator_
    
    def train_gradient_boosting(self, X_train, y_train):
        """训练梯度提升树模型"""
        print("\n🚀 训练梯度提升树模型...")
        
        param_grid = {
            'n_estimators': [50, 100, 200],
            'learning_rate': [0.01, 0.05, 0.1],
            'max_depth': [3, 5, 7],
            'min_samples_split': [2, 5],
            'subsample': [0.8, 1.0]
        }
        
        gb = GradientBoostingClassifier(random_state=42)
        grid_search = GridSearchCV(gb, param_grid, cv=5, scoring='accuracy', n_jobs=-1)
        grid_search.fit(X_train, y_train)
        
        print(f"✅ 最佳参数: {grid_search.best_params_}")
        print(f"✅ 交叉验证准确率: {grid_search.best_score_:.4f}")
        
        return grid_search.best_estimator_
    
    def evaluate_model(self, model, X_test, y_test):
        """评估模型性能"""
        print("\n📊 模型评估...")
        
        y_pred = model.predict(X_test)
        y_pred_proba = model.predict_proba(X_test) if hasattr(model, 'predict_proba') else None
        
        # 基本指标
        accuracy = accuracy_score(y_test, y_pred)
        precision = precision_score(y_test, y_pred, average='weighted', zero_division=0)
        recall = recall_score(y_test, y_pred, average='weighted', zero_division=0)
        f1 = f1_score(y_test, y_pred, average='weighted', zero_division=0)
        
        # 混淆矩阵
        cm = confusion_matrix(y_test, y_pred)
        
        # AUC (多分类)
        auc = 0.0
        if y_pred_proba is not None and len(np.unique(y_test)) > 1:
            try:
                auc = roc_auc_score(y_test, y_pred_proba, multi_class='ovr', average='weighted')
            except:
                pass
        
        print(f"\n✅ 准确率 (Accuracy): {accuracy:.4f}")
        print(f"✅ 精确率 (Precision): {precision:.4f}")
        print(f"✅ 召回率 (Recall): {recall:.4f}")
        print(f"✅ F1分数: {f1:.4f}")
        if auc > 0:
            print(f"✅ AUC: {auc:.4f}")
        
        print(f"\n混淆矩阵:\n{cm}")
        
        # 分类报告
        print(f"\n详细分类报告:")
        print(classification_report(y_test, y_pred, zero_division=0))
        
        return {
            'accuracy': accuracy,
            'precision': precision,
            'recall': recall,
            'f1_score': f1,
            'auc_score': auc,
            'confusion_matrix': cm.tolist()
        }
    
    def save_model(self, model, metrics, model_type='random_forest'):
        """保存模型"""
        timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
        version = f"v{timestamp}"
        
        # 保存模型文件
        model_filename = f"{model_type}_{version}.pkl"
        model_path = os.path.join(MODEL_DIR, model_filename)
        
        model_data = {
            'model': model,
            'scaler': self.scaler,
            'label_encoder': self.label_encoder,
            'feature_names': self.feature_names,
            'model_type': model_type,
            'version': version,
            'metrics': metrics,
            'trained_at': datetime.now().isoformat()
        }
        
        with open(model_path, 'wb') as f:
            pickle.dump(model_data, f)
        
        print(f"\n✅ 模型已保存: {model_path}")
        
        # 保存到数据库
        self._save_model_to_db(model_filename, metrics, model_type, version)
        
        return model_path
    
    def _save_model_to_db(self, model_filename, metrics, model_type, version):
        """保存模型信息到数据库"""
        try:
            cursor = self.conn.cursor()
            
            sql = """
            INSERT INTO SM_ai_model_version (
                model_name, version, model_type,
                accuracy, precision_score, recall_score, f1_score, auc_score,
                confusion_matrix, model_file_path,
                description
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            """
            
            cursor.execute(sql, (
                'medical_diagnosis_ai',
                version,
                model_type,
                metrics['accuracy'],
                metrics['precision'],
                metrics['recall'],
                metrics['f1_score'],
                metrics['auc_score'],
                json.dumps(metrics['confusion_matrix']),
                model_filename,
                f"自动训练模型 - {model_type}"
            ))
            
            self.conn.commit()
            print(f"✅ 模型信息已保存到数据库")
            
        except Exception as e:
            print(f"⚠️  保存模型信息失败: {e}")
    
    def train_all_models(self):
        """训练并比较所有模型"""
        print("="*60)
        print("🚀 开始训练AI诊断模型")
        print("="*60)
        
        # 1. 加载数据
        X, y = self.load_training_data()
        if X is None or y is None:
            print("❌ 数据加载失败，训练终止")
            return
        
        # 2. 数据标准化
        X_scaled = self.scaler.fit_transform(X)
        y_encoded = self.label_encoder.fit_transform(y)
        
        # 3. 划分训练集和测试集
        X_train, X_test, y_train, y_test = train_test_split(
            X_scaled, y_encoded, test_size=0.2, random_state=42, stratify=y_encoded
        )
        
        print(f"\n📊 数据划分:")
        print(f"  训练集: {len(X_train)} 条")
        print(f"  测试集: {len(X_test)} 条")
        
        # 4. 训练多个模型
        models = {}
        model_metrics = {}
        
        # 逻辑回归
        lr_model = self.train_logistic_regression(X_train, y_train)
        models['logistic_regression'] = lr_model
        model_metrics['logistic_regression'] = self.evaluate_model(lr_model, X_test, y_test)
        
        # 随机森林
        rf_model = self.train_random_forest(X_train, y_train)
        models['random_forest'] = rf_model
        model_metrics['random_forest'] = self.evaluate_model(rf_model, X_test, y_test)
        
        # 梯度提升树
        gb_model = self.train_gradient_boosting(X_train, y_train)
        models['gradient_boosting'] = gb_model
        model_metrics['gradient_boosting'] = self.evaluate_model(gb_model, X_test, y_test)
        
        # 5. 选择最佳模型
        best_model_type = max(model_metrics, key=lambda k: model_metrics[k]['accuracy'])
        best_model = models[best_model_type]
        best_metrics = model_metrics[best_model_type]
        
        print("\n" + "="*60)
        print(f"🏆 最佳模型: {best_model_type}")
        print(f"🎯 准确率: {best_metrics['accuracy']:.4f}")
        print("="*60)
        
        # 6. 保存最佳模型
        model_path = self.save_model(best_model, best_metrics, best_model_type)
        
        # 7. 特征重要性分析（如果支持）
        if hasattr(best_model, 'feature_importances_'):
            self._plot_feature_importance(best_model)
        
        print("\n✅ 训练完成！")
        
        return best_model, best_metrics
    
    def _plot_feature_importance(self, model):
        """分析特征重要性"""
        if not hasattr(model, 'feature_importances_'):
            return
        
        importance = model.feature_importances_
        feature_importance = pd.DataFrame({
            'feature': self.feature_names,
            'importance': importance
        }).sort_values('importance', ascending=False)
        
        print("\n📊 特征重要性 TOP 10:")
        print(feature_importance.head(10).to_string(index=False))
    
    def close(self):
        """关闭数据库连接"""
        if self.conn:
            self.conn.close()
            print("\n✅ 数据库连接已关闭")


def main():
    """主函数"""
    trainer = MedicalAITrainer(DB_CONFIG)
    
    try:
        best_model, metrics = trainer.train_all_models()
        
        if metrics['accuracy'] >= 0.85:
            print(f"\n🎉 恭喜！模型准确率达到 {metrics['accuracy']:.2%}，超过85%目标！")
        else:
            print(f"\n⚠️  模型准确率 {metrics['accuracy']:.2%}，未达到85%目标")
            print("建议：")
            print("  1. 收集更多训练数据（当前可能不足）")
            print("  2. 增加更多特征工程")
            print("  3. 调整模型超参数")
        
    except Exception as e:
        print(f"\n❌ 训练失败: {e}")
        import traceback
        traceback.print_exc()
    
    finally:
        trainer.close()


if __name__ == '__main__':
    main()
