#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
导出训练数据脚本
从数据库导出训练数据到CSV文件，方便数据分析和模型训练
"""

import pymysql
import pandas as pd
import json
from datetime import datetime

DB_CONFIG = {
    'host': 'localhost',
    'port': 3306,
    'user': 'root',
    'password': 'your_password',  # 修改为实际密码
    'database': 'sm_medical',
    'charset': 'utf8mb4'
}

def export_training_data():
    """导出训练数据"""
    try:
        # 连接数据库
        conn = pymysql.connect(**DB_CONFIG)
        print("✅ 数据库连接成功")
        
        # 查询训练数据
        query = """
        SELECT 
            t.*,
            u.gender as patient_gender,
            TIMESTAMPDIFF(YEAR, u.birthday, t.created_at) as calculated_age,
            m.disease_category,
            m.severity_level,
            m.recommended_dept
        FROM SM_ai_training_data t
        LEFT JOIN SM_user u ON t.patient_id = u.id
        LEFT JOIN SM_disease_mapping m ON t.diagnosis_icd10 = m.icd10_code
        WHERE t.is_verified = 1
        ORDER BY t.created_at DESC
        """
        
        df = pd.read_sql(query, conn)
        
        if len(df) == 0:
            print("⚠️  没有找到训练数据")
            return
        
        print(f"📊 导出数据: {len(df)} 条记录")
        
        # 保存到CSV
        timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
        filename = f'training_data_{timestamp}.csv'
        df.to_csv(filename, index=False, encoding='utf-8-sig')
        
        print(f"✅ 数据已导出到: {filename}")
        
        # 统计信息
        print(f"\n📈 数据统计:")
        print(f"  总样本数: {len(df)}")
        print(f"  疾病类别数: {df['doctor_diagnosis'].nunique()}")
        print(f"  时间范围: {df['created_at'].min()} 至 {df['created_at'].max()}")
        
        print(f"\n📋 疾病分布:")
        print(df['doctor_diagnosis'].value_counts().head(10))
        
        conn.close()
        
    except Exception as e:
        print(f"❌ 导出失败: {e}")
        import traceback
        traceback.print_exc()

def check_data_quality():
    """检查数据质量"""
    try:
        conn = pymysql.connect(**DB_CONFIG)
        cursor = conn.cursor()
        
        print("\n" + "="*60)
        print("📊 数据质量检查")
        print("="*60)
        
        # 1. 总样本数
        cursor.execute("SELECT COUNT(*) FROM SM_ai_training_data")
        total = cursor.fetchone()[0]
        print(f"\n1. 总样本数: {total}")
        
        # 2. 已验证样本
        cursor.execute("SELECT COUNT(*) FROM SM_ai_training_data WHERE is_verified = 1")
        verified = cursor.fetchone()[0]
        print(f"2. 已验证样本: {verified} ({verified/total*100:.1f}%)")
        
        # 3. 数据质量分布
        cursor.execute("""
            SELECT data_quality, COUNT(*) as cnt 
            FROM SM_ai_training_data 
            GROUP BY data_quality
        """)
        print(f"\n3. 数据质量分布:")
        for row in cursor.fetchall():
            quality_map = {1: '高质量', 2: '中等', 3: '低质量'}
            print(f"   {quality_map.get(row[0], '未知')}: {row[1]}")
        
        # 4. 疾病类别数
        cursor.execute("""
            SELECT COUNT(DISTINCT doctor_diagnosis) 
            FROM SM_ai_training_data 
            WHERE doctor_diagnosis IS NOT NULL AND doctor_diagnosis != ''
        """)
        disease_count = cursor.fetchone()[0]
        print(f"\n4. 疾病类别数: {disease_count}")
        
        # 5. TOP疾病
        cursor.execute("""
            SELECT doctor_diagnosis, COUNT(*) as cnt 
            FROM SM_ai_training_data 
            WHERE doctor_diagnosis IS NOT NULL 
            GROUP BY doctor_diagnosis 
            ORDER BY cnt DESC 
            LIMIT 10
        """)
        print(f"\n5. TOP 10 疾病:")
        for idx, row in enumerate(cursor.fetchall(), 1):
            print(f"   {idx}. {row[0]}: {row[1]} 条")
        
        # 6. 数据完整性
        cursor.execute("""
            SELECT 
                SUM(CASE WHEN age IS NOT NULL THEN 1 ELSE 0 END) as has_age,
                SUM(CASE WHEN systolic_bp IS NOT NULL THEN 1 ELSE 0 END) as has_bp,
                SUM(CASE WHEN heart_rate IS NOT NULL THEN 1 ELSE 0 END) as has_hr,
                SUM(CASE WHEN temperature IS NOT NULL THEN 1 ELSE 0 END) as has_temp,
                SUM(CASE WHEN blood_sugar IS NOT NULL THEN 1 ELSE 0 END) as has_sugar,
                COUNT(*) as total
            FROM SM_ai_training_data
        """)
        row = cursor.fetchone()
        total = row[5]
        print(f"\n6. 数据完整性:")
        print(f"   年龄: {row[0]/total*100:.1f}%")
        print(f"   血压: {row[1]/total*100:.1f}%")
        print(f"   心率: {row[2]/total*100:.1f}%")
        print(f"   体温: {row[3]/total*100:.1f}%")
        print(f"   血糖: {row[4]/total*100:.1f}%")
        
        # 7. 建议
        print(f"\n" + "="*60)
        if verified < 50:
            print("⚠️  建议: 训练数据不足50条，建议继续收集")
        elif verified < 100:
            print("⚠️  建议: 训练数据较少，模型准确率可能不理想")
        elif verified < 200:
            print("✅ 建议: 数据量基本满足，可以开始训练")
        else:
            print("✅ 建议: 数据量充足，可以训练高质量模型")
        
        conn.close()
        
    except Exception as e:
        print(f"❌ 检查失败: {e}")

if __name__ == '__main__':
    print("🚀 AI训练数据导出工具\n")
    
    # 检查数据质量
    check_data_quality()
    
    # 导出数据
    print("\n")
    export_training_data()
