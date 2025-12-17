package service

import (
	"fmt"
	"math"
	"strings"
)

// performAIDiagnosis 执行AI智能诊断 - 增强版V2
func (s *ConsultationService) performAIDiagnosis(chiefComplaint string, symptoms map[string]interface{}) AIResult {
	riskScore := 0.0 // 改用浮点数以支持加权计算
	possibleDiseases := []string{}
	diseaseScores := make(map[string]float64) // 疾病评分映射
	detailedAnalysis := make(map[string]string)
	lifestyleAdvice := []string{}
	recommendedDept := "全科"
	urgencyLevel := "normal"
	followUpAdvice := ""
	symptomCount := 0 // 统计有效症状数量
	
	// 提取症状数据
	age, _ := symptoms["age"].(float64)
	bloodPressure, _ := symptoms["bloodPressure"].(string)
	heartRate, _ := symptoms["heartRate"].(float64)
	temperature, _ := symptoms["temperature"].(float64)
	bloodSugar, _ := symptoms["bloodSugar"].(float64)
	otherSymptoms, _ := symptoms["otherSymptoms"].(string)
	
	// 合并主诉和其他症状，确保完整分析患者的所有描述
	combinedSymptoms := chiefComplaint
	if otherSymptoms != "" {
		combinedSymptoms = chiefComplaint + " " + otherSymptoms
	}
	
	// 预先声明血压变量（用于后续症状组合判断）
	var systolic, diastolic int
	
	// 1. 血压分析 - 改进版:考虑年龄因素和脉压差
	if bloodPressure != "" {
		fmt.Sscanf(bloodPressure, "%d/%d", &systolic, &diastolic)
		symptomCount++
		
		// 计算脉压差(正常30-40mmHg)
		pulsePressure := systolic - diastolic
		
		// 年龄校正:老年人血压标准可适当放宽
		ageAdjustment := 0
		if age > 65 {
			ageAdjustment = 10 // 65岁以上,标准放宽10mmHg
		}
		
		if systolic >= 180 || diastolic >= 110 {
			baseScore := 40.0
			if age > 60 {
				baseScore *= 1.3 // 老年人高血压风险更高
			}
			riskScore += baseScore
			diseaseScores["高血压3级"] = baseScore
			possibleDiseases = append(possibleDiseases, "高血压3级(重度)")
			detailedAnalysis["心血管系统"] = fmt.Sprintf("血压%s(脉压差%dmmHg),达到3级高血压标准,存在严重心血管风险。建议立即就医,可能需要药物干预。", bloodPressure, pulsePressure)
			recommendedDept = "心血管内科/心内科"
			urgencyLevel = "urgent"
			lifestyleAdvice = append(lifestyleAdvice, "低盐饮食,每日食盐不超过6克", "避免剧烈运动和情绪激动", "每日监测血压2-3次")
		} else if systolic >= (160-ageAdjustment) || diastolic >= (100-ageAdjustment/2) {
			baseScore := 30.0
			if age > 60 {
				baseScore *= 1.2
			}
			riskScore += baseScore
			diseaseScores["高血压2级"] = baseScore
			possibleDiseases = append(possibleDiseases, "高血压2级(中度)")
			detailedAnalysis["心血管系统"] = fmt.Sprintf("血压%s,达到2级高血压标准,需要积极控制。建议3天内就诊心内科。", bloodPressure)
			recommendedDept = "心血管内科"
			urgencyLevel = "attention"
			lifestyleAdvice = append(lifestyleAdvice, "减少钠摄入,多食新鲜蔬果", "适度有氧运动,如散步、游泳", "保持良好睡眠,每晚7-8小时")
		} else if systolic >= (140-ageAdjustment) || diastolic >= (90-ageAdjustment/2) {
			riskScore += 20.0
			diseaseScores["高血压1级"] = 20.0
			possibleDiseases = append(possibleDiseases, "高血压1级(轻度)")
			detailedAnalysis["心血管系统"] = fmt.Sprintf("血压%s,达到1级高血压标准。建议1-2周内就诊,评估是否需要药物治疗。", bloodPressure)
			recommendedDept = "心内科"
			lifestyleAdvice = append(lifestyleAdvice, "控制体重,BMI保持18.5-24", "戒烟限酒,尤其避免醉酒", "减少压力,保持心情愉悦")
		} else if systolic < 90 || diastolic < 60 {
			riskScore += 15.0
			diseaseScores["低血压"] = 15.0
			possibleDiseases = append(possibleDiseases, "低血压")
			detailedAnalysis["心血管系统"] = fmt.Sprintf("血压%s,低于正常范围。可能引起头晕、乏力等症状。建议就诊查明原因。", bloodPressure)
			recommendedDept = "心内科"
			lifestyleAdvice = append(lifestyleAdvice, "增加食盐和水分摄入", "避免突然起立或长时间站立", "适度运动,增强体质")
		} else {
			detailedAnalysis["心血管系统"] = fmt.Sprintf("血压%s,处于正常范围内,请继续保持。", bloodPressure)
		}
		
		// 脉压差异常警告
		if pulsePressure > 60 {
			riskScore += 10.0
			detailedAnalysis["心血管系统"] += fmt.Sprintf(" ⚠️脉压差%dmmHg偏大,可能提示动脉硬化,建议检查。", pulsePressure)
		} else if pulsePressure < 20 {
			riskScore += 8.0
			detailedAnalysis["心血管系统"] += fmt.Sprintf(" ⚠️脉压差%dmmHg偏小,可能提示心功能不全。", pulsePressure)
		}
	}
	
	// 2. 心率分析 - 改进版:考虑年龄和体温影响
	if heartRate > 0 {
		symptomCount++
		
		// 年龄因素:老年人正常心率偏低
		normalMax := 100.0
		normalMin := 60.0
		if age > 65 {
			normalMax = 95.0
			normalMin = 55.0
		} else if age < 18 {
			normalMax = 110.0 // 儿童心率偏快
		}
		
		// 发热时心率会升高(体温每升高1℃,心率增加10-15次)
		feverAdjustment := 0.0
		if temperature > 37.0 {
			feverAdjustment = (temperature - 37.0) * 12
		}
		
		adjustedHeartRate := heartRate - feverAdjustment
		
		if heartRate > normalMax {
			baseScore := 15.0
			if heartRate > 120 {
				baseScore = 25.0
				detailedAnalysis["心率分析"] = fmt.Sprintf("心率%.0f次/分,明显过快。可能与心脏疾病、甲亢、贫血或情绪紧张有关。建议尽快就诊。", heartRate)
				urgencyLevel = "attention"
				diseaseScores["心动过速"] = baseScore
			} else {
				if feverAdjustment > 0 {
					baseScore = 8.0 // 发热导致的心率加快风险较低
					detailedAnalysis["心率分析"] = fmt.Sprintf("心率%.0f次/分,考虑体温%.1f℃因素,校正后约%.0f次/分。发热可导致心率加快,属正常生理反应。", heartRate, temperature, adjustedHeartRate)
				} else {
					detailedAnalysis["心率分析"] = fmt.Sprintf("心率%.0f次/分,轻度过快。可能与运动、焦虑、咖啡因摄入有关。建议休息后重新测量。", heartRate)
				}
			}
			riskScore += baseScore
			possibleDiseases = append(possibleDiseases, "心动过速")
			lifestyleAdvice = append(lifestyleAdvice, "避免咖啡因和烟酒", "减少压力,保证睡眠", "如伴有心慌、胸闷,请立即就医")
		} else if heartRate < normalMin {
			baseScore := 10.0
			if heartRate < 45 {
				baseScore = 20.0
				urgencyLevel = "attention"
			}
			riskScore += baseScore
			diseaseScores["心动过缓"] = baseScore
			possibleDiseases = append(possibleDiseases, "心动过缓")
			detailedAnalysis["心率分析"] = fmt.Sprintf("心率%.0f次/分,低于正常范围。如果您是运动员可能属正常,否则建议就诊检查。", heartRate)
			lifestyleAdvice = append(lifestyleAdvice, "注意是否有头晕、乏力等症状", "如服用降压药,请咨询医生")
		} else {
			detailedAnalysis["心率分析"] = fmt.Sprintf("心率%.0f次/分,处于正常范围(%.0f-%.0f次/分)。", heartRate, normalMin, normalMax)
		}
	}
	
	// 3. 体温分析
	if temperature > 0 {
		if temperature >= 39.0 {
			riskScore += 25
			possibleDiseases = append(possibleDiseases, "高热")
			detailedAnalysis["体温分析"] = fmt.Sprintf("体温%.1f℃,已达高热标准。可能是严重感染,建议立即就医。请多饮水,可物理降温。", temperature)
			recommendedDept = "感染科/呼吸内科"
			urgencyLevel = "urgent"
			lifestyleAdvice = append(lifestyleAdvice, "大量饮水,每小时2000ml以上", "物理降温:温水擦浴、冰敷", "避免盖厚被,保持通风", "如体温持续不降或出现抽搐,立即拨打120")
		} else if temperature >= 38.0 {
			riskScore += 20
			possibleDiseases = append(possibleDiseases, "中度发热")
			detailedAnalysis["体温分析"] = fmt.Sprintf("体温%.1f℃,中度发热。建议24-48小时内就诊,明确发热原因。", temperature)
			recommendedDept = "呼吸内科/全科"
			urgencyLevel = "attention"
			lifestyleAdvice = append(lifestyleAdvice, "多饮水,促进代谢", "清淡饮食,易消化为主", "充分休息,避免劳累")
		} else if temperature >= 37.3 {
			riskScore += 10
			possibleDiseases = append(possibleDiseases, "低热")
			detailedAnalysis["体温分析"] = fmt.Sprintf("体温%.1f℃,轻度发热。可能是感冒早期或疲劳引起,建议观察24小时。", temperature)
			lifestyleAdvice = append(lifestyleAdvice, "注意保暖,避免着凉", "多休息,保证睡眠", "如体温持续上升,请就诊")
		} else if temperature < 36.0 {
			riskScore += 15
			possibleDiseases = append(possibleDiseases, "体温过低")
			detailedAnalysis["体温分析"] = fmt.Sprintf("体温%.1f℃,低于正常范围。可能与甲减、休克等有关,建议就诊检查。", temperature)
			urgencyLevel = "attention"
			lifestyleAdvice = append(lifestyleAdvice, "注意保暖,增加衣物", "适当运动,促进血液循环", "如伴有乏力、嘴唇发紫,请立即就医")
		} else {
			detailedAnalysis["体温分析"] = fmt.Sprintf("体温%.1f℃,处于正常范围(36.1-37.2℃)。", temperature)
		}
	}
	
	// 4. 血糖分析
	if bloodSugar > 0 {
		if bloodSugar >= 11.1 {
			riskScore += 30
			possibleDiseases = append(possibleDiseases, "糖尿病风险(血糖过高)")
			detailedAnalysis["血糖分析"] = fmt.Sprintf("血糖%.1fmmol/L,达到糖尿病诊断标准。建议尽快就诊内分泌科,进行糖化血红蛋白等检查。", bloodSugar)
			recommendedDept = "内分泌科"
			urgencyLevel = "urgent"
			lifestyleAdvice = append(lifestyleAdvice, "控制饮食,减少糖分和精制碳水摄入", "增加运动,每天至少30分钟有氧运动", "监测血糖,建议购买血糖仪", "多饮水,预防脱水")
		} else if bloodSugar >= 7.0 {
			riskScore += 20
			possibleDiseases = append(possibleDiseases, "空腹血糖异常")
			detailedAnalysis["血糖分析"] = fmt.Sprintf("血糖%.1fmmol/L,高于正常范围。建议就诊做糖耐量试验,排除糖尿病前期。", bloodSugar)
			recommendedDept = "内分泌科"
			urgencyLevel = "attention"
			lifestyleAdvice = append(lifestyleAdvice, "减少糖分摄入,控制饮料和甜食", "增加膳食纤维,多吃全谷物和蔬菜", "控制体重,超重者应减肥")
		} else if bloodSugar < 3.9 {
			riskScore += 25
			possibleDiseases = append(possibleDiseases, "低血糖")
			detailedAnalysis["血糖分析"] = fmt.Sprintf("血糖%.1fmmol/L,低于正常范围。如有头晕、出汗、心慌症状,请立即进食糖分。", bloodSugar)
			urgencyLevel = "urgent"
			lifestyleAdvice = append(lifestyleAdvice, "立即补充糖分:果汁、糖果或饼干", "如服用降糖药,请咨询医生调整剂量", "定时进餐,避免长时间空腹")
		} else {
			detailedAnalysis["血糖分析"] = fmt.Sprintf("血糖%.1fmmol/L,处于正常范围(3.9-6.1mmol/L)。", bloodSugar)
		}
	}
	
	// 5. 年龄因素
	if age > 0 {
		if age > 60 {
			riskScore += 5
			detailedAnalysis["年龄评估"] = fmt.Sprintf("您已%.0f岁,属于老年人群,建议每年进行全面体检,重点关注心血管、糖尿病、骨质疏松等问题。", age)
			followUpAdvice = "建议每年进行一次全面体检,包括血常规、生化全套、心电图、胸片、腹部B超等。"
		} else if age > 40 {
			detailedAnalysis["年龄评估"] = fmt.Sprintf("您%.0f岁,建议每年进行常规体检,关注三高(高血压、高血糖、高血脂)和肿瘤筛查。", age)
			followUpAdvice = "建议每1-2年进行一次健康体检,重点筛查慢性病风险。"
		}
	}
	
	// 6. 其他症状关键词分析 - 增强版：症状组合判断（含主诉）
	isEmergency := false
	symptomFlags := make(map[string]bool) // 症状标记
	
	if combinedSymptoms != "" {
		symptomsLower := strings.ToLower(combinedSymptoms)
		
		// 提取症状标记（用于组合判断）
		symptomFlags["头痛"] = strings.Contains(combinedSymptoms, "头痛") || strings.Contains(combinedSymptoms, "头晕")
		symptomFlags["发热"] = strings.Contains(combinedSymptoms, "发热") || strings.Contains(combinedSymptoms, "发烧")
		symptomFlags["胸痛"] = strings.Contains(symptomsLower, "胸痛") || strings.Contains(symptomsLower, "胸闷")
		symptomFlags["呼吸困难"] = strings.Contains(symptomsLower, "呼吸困难") || strings.Contains(symptomsLower, "气促")
		symptomFlags["咳嗽"] = strings.Contains(symptomsLower, "咳嗽")
		symptomFlags["腹痛"] = strings.Contains(symptomsLower, "腹痛")
		symptomFlags["恶心"] = strings.Contains(symptomsLower, "恶心") || strings.Contains(symptomsLower, "呕吐")
		
		// 紧急情况关键词
		emergencyKeywords := []string{"胸痛", "呼吸困难", "意识障碍", "意识不清", "晕厥", "抽搐", "喉头水肿", "大量出血", "剧烈腹痛"}
		for _, keyword := range emergencyKeywords {
			if strings.Contains(combinedSymptoms, keyword) {
				isEmergency = true
				riskScore += 40
				urgencyLevel = "emergency"
				break
			}
		}
		
		// ========== 症状组合分析（提升准确性）==========
		
		// 组合1：高血压+头痛+头晕 → 高血压脑病风险
		if (systolic >= 140 || diastolic >= 90) && symptomFlags["头痛"] {
			riskScore += 10.0
			if _, exists := diseaseScores["高血压"]; exists {
				diseaseScores["高血压"] += 10.0
			}
			detailedAnalysis["症状组合分析"] = "高血压伴头痛/头晕，需警惕高血压脑病，建议尽快就医检查。"
			if urgencyLevel == "normal" {
				urgencyLevel = "attention"
			}
		}
		
		// 组合2：胸痛+呼吸困难 → 心肺疾病高危
		if symptomFlags["胸痛"] && symptomFlags["呼吸困难"] {
			riskScore += 25.0
			if !isEmergency {
				isEmergency = true
				urgencyLevel = "emergency"
			}
			detailedAnalysis["高危组合"] = "胸痛合并呼吸困难是心梗、肺栓塞等危重症的典型表现！请立即拨打120或前往急诊！"
			recommendedDept = "急诊科/心内科"
			possibleDiseases = append(possibleDiseases, "急性心肌梗死风险")
			diseaseScores["急性心肌梗死风险"] = 50.0
		}
		
		// 组合3：发热+咳嗽+呼吸困难 → 肺炎
		if symptomFlags["发热"] && symptomFlags["咳嗽"] && (symptomFlags["呼吸困难"] || temperature >= 38.5) {
			riskScore += 15.0
			possibleDiseases = append(possibleDiseases, "肺炎")
			diseaseScores["肺炎"] = 25.0
			detailedAnalysis["呼吸系统症状"] = "发热+咳嗽+呼吸困难三联征，高度怀疑肺炎。建议拍胸片或CT，查血常规、CRP。"
			recommendedDept = "呼吸内科"
			lifestyleAdvice = append(lifestyleAdvice, "卧床休息，多饮水", "保持室内通风", "如呼吸困难加重，立即就医")
		}
		
		// 组合4：腹痛+恶心呕吐+发热 → 急性阑尾炎等急腹症
		if symptomFlags["腹痛"] && symptomFlags["恶心"] && (symptomFlags["发热"] || temperature >= 37.5) {
			riskScore += 20.0
			possibleDiseases = append(possibleDiseases, "急性阑尾炎")
			diseaseScores["急性阑尾炎"] = 30.0
			detailedAnalysis["消化系统症状"] += " 腹痛伴发热、恶心呕吐，需警惕急性阑尾炎、胆囊炎等急腹症，建议12小时内就诊。"
			urgencyLevel = "urgent"
			recommendedDept = "普外科/急诊科"
		}
		
		// 组合5：高血糖+多饮多尿 → 糖尿病
		if bloodSugar >= 11.1 && (strings.Contains(combinedSymptoms, "口渴") || strings.Contains(combinedSymptoms, "多尿")) {
			riskScore += 15.0
			if _, exists := diseaseScores["糖尿病"]; exists {
				diseaseScores["糖尿病"] += 15.0
			} else {
				possibleDiseases = append(possibleDiseases, "糖尿病典型症状")
				diseaseScores["糖尿病典型症状"] = 35.0
			}
			detailedAnalysis["血糖分析"] += " 伴有多饮多尿症状，糖尿病诊断明确，需立即就诊内分泌科。"
		}
		
		// 组合6：低血糖+心慌出汗 → 低血糖反应
		if bloodSugar > 0 && bloodSugar < 3.9 && (strings.Contains(combinedSymptoms, "心慌") || strings.Contains(combinedSymptoms, "出汗")) {
			riskScore += 10.0
			if _, exists := diseaseScores["低血糖"]; exists {
				diseaseScores["低血糖"] += 10.0
			}
			detailedAnalysis["血糖分析"] += " 典型的低血糖症状（心慌、出汗），请立即补充糖分！"
			lifestyleAdvice = append(lifestyleAdvice, "立即进食含糖食物：糖果、果汁、蜂蜜", "15分钟后复查血糖", "如无改善，立即就医")
		}
		
		// ========== 单症状分析 ==========
		
		// 心血管相关
		if strings.Contains(symptomsLower, "胸痛") || strings.Contains(symptomsLower, "胸闷") {
			riskScore += 25
			possibleDiseases = append(possibleDiseases, "心血管疾病风险")
			detailedAnalysis["心血管症状"] = "胸痛/胸闷是心脏疾病的常见症状,可能提示冠心病、心绞痛等。建议立即就医,必要时做心电图检查。"
			recommendedDept = "心血管内科/急诊科"
			if urgencyLevel != "emergency" {
				urgencyLevel = "urgent"
			}
		}
		if strings.Contains(symptomsLower, "呼吸困难") || strings.Contains(symptomsLower, "气促") {
			riskScore += 20
			possibleDiseases = append(possibleDiseases, "呼吸系统异常")
			detailedAnalysis["呼吸系统症状"] = "呼吸困难可能提示哮喘、肺炎、心功能不全等。建议及时就医,进行胸部X线或CT检查。"
			recommendedDept = "呼吸内科/急诊科"
		}
		
		// 神经系统
		if strings.Contains(symptomsLower, "头晕") || strings.Contains(symptomsLower, "头痛") {
			riskScore += 10
			if strings.Contains(symptomsLower, "剧烈") || strings.Contains(symptomsLower, "持续") {
				riskScore += 10
				possibleDiseases = append(possibleDiseases, "偏头痛/高血压头痛")
				detailedAnalysis["神经系统症状"] = "剧烈或持续头痛/头晕需要引起重视。可能与高血压、偏头痛、脑血管疾病有关。建议1-3天内就诊神经内科。"
				if recommendedDept == "全科" {
					recommendedDept = "神经内科"
				}
			} else {
				detailedAnalysis["神经系统症状"] = "轻度头晕/头痛,可能与睡眠不足、压力、颈椎病有关。建议注意休息,如症状加重请就诊。"
			}
		}
		if strings.Contains(symptomsLower, "意识") || strings.Contains(symptomsLower, "晕厥") || strings.Contains(symptomsLower, "抽搐") {
			riskScore += 35
			possibleDiseases = append(possibleDiseases, "严重神经系统症状")
			detailedAnalysis["紧急情况"] = "意识障碍、晕厥或抽搐属于紧急情况!建议立即拨打120急救电话,或前往最近医院急诊科。"
			recommendedDept = "急诊科"
			urgencyLevel = "emergency"
		}
		
		// 消化系统
		if strings.Contains(symptomsLower, "腹痛") || strings.Contains(symptomsLower, "呕吐") {
			riskScore += 15
			if strings.Contains(symptomsLower, "剧烈") {
				riskScore += 10
				possibleDiseases = append(possibleDiseases, "急腹症")
				detailedAnalysis["消化系统症状"] = "剧烈腹痛需要立即就医,排除阑尾炎、胆囊炎、肠梗阻等急症。建议禁食,尽快前往急诊科。"
				recommendedDept = "消化内科/急诊科"
				urgencyLevel = "urgent"
			} else {
				detailedAnalysis["消化系统症状"] = "腹痛/呕吐可能与胃炎、胃溃疡、饮食不当有关。建议清淡饮食,如症状持续请就诊消化内科。"
				if recommendedDept == "全科" {
					recommendedDept = "消化内科"
				}
			}
		}
		if strings.Contains(symptomsLower, "腹泻") {
			riskScore += 10
			possibleDiseases = append(possibleDiseases, "消化系统疾病")
			detailedAnalysis["消化系统症状"] = "腹泻可能与肠炎、食物中毒、菌群失调有关。建议多饮水防脱水,清淡饮食,如超过3天请就诊。"
			lifestyleAdvice = append(lifestyleAdvice, "补充水分和电解质,可饮用口服补液盐", "清淡饮食,避免油腻和辣椒", "注意食品卫生")
		}
		
		// 感染症状
		if strings.Contains(symptomsLower, "咳嗽") || strings.Contains(symptomsLower, "咽痛") {
			riskScore += 10
			possibleDiseases = append(possibleDiseases, "呼吸道感染")
			detailedAnalysis["感染症状"] = "咳嗽/咽痛是上呼吸道感染的常见症状。建议多休息、多喝水,如伴有高热或症状加重请就诊。"
			if recommendedDept == "全科" {
				recommendedDept = "呼吸内科/耳鼻喉科"
			}
			lifestyleAdvice = append(lifestyleAdvice, "多喝温水,每天2000ml以上", "注意保暖和通风", "佩戴口罩,避免传染他人")
		}
		
		// 疼痛相关
		if strings.Contains(symptomsLower, "疼痛") {
			riskScore += 8
			detailedAnalysis["疼痛症状"] = "疼痛需要根据具体部位判断。建议详细记录疼痛特点(部位、性质、时间),就诊时告知医生。"
		}
		
		// 皮肤科症状
		if strings.Contains(symptomsLower, "瘙痒") || strings.Contains(symptomsLower, "痒") || strings.Contains(combinedSymptoms, "瘙痒") {
			riskScore += 12
			possibleDiseases = append(possibleDiseases, "皮肤瘙痒症")
			
			if strings.Contains(combinedSymptoms, "全身") {
				riskScore += 8
				possibleDiseases = append(possibleDiseases, "全身性皮肤病/过敏反应")
				diseaseScores["全身性皮肤病/过敏反应"] = 20.0
				detailedAnalysis["皮肤症状"] = "全身皮肤瘙痒可能与过敏、湿疹、荨麻疹、肝胆疾病、肾脏疾病、内分泌疾病等有关。建议尽快就诊皮肤科,进行过敏原检测和相关检查,明确病因。"
				recommendedDept = "皮肤科"
				urgencyLevel = "attention"
				lifestyleAdvice = append(lifestyleAdvice, "避免搔抓,防止皮肤破损感染", "穿着宽松透气的纯棉衣物", "避免接触可能的过敏原(海鲜、芒果、花粉等)", "保持皮肤清洁干燥,可用温水擦浴", "避免使用刺激性化妆品和洗涤用品")
			} else {
				diseaseScores["局部皮肤瘙痒"] = 12.0
				detailedAnalysis["皮肤症状"] = "局部皮肤瘙痒可能与接触性皮炎、虫咬、真菌感染等有关。建议观察瘙痒部位是否有红肿、皮疹等,如症状持续3天以上请就诊皮肤科。"
				if recommendedDept == "全科" {
					recommendedDept = "皮肤科"
				}
				lifestyleAdvice = append(lifestyleAdvice, "避免搔抓患处", "保持局部清洁干燥", "可用冷敷缓解瘙痒")
			}
		}
		
		if strings.Contains(symptomsLower, "皮疹") || strings.Contains(symptomsLower, "红疹") || strings.Contains(combinedSymptoms, "皮疹") {
			riskScore += 15
			possibleDiseases = append(possibleDiseases, "皮疹/皮肤病")
			diseaseScores["皮疹/皮肤病"] = 15.0
			detailedAnalysis["皮肤症状"] = "皮疹可能提示过敏、病毒感染、细菌感染、自身免疫性疾病等。建议就诊皮肤科,必要时进行皮肤活检。"
			if recommendedDept == "全科" {
				recommendedDept = "皮肤科"
			}
			lifestyleAdvice = append(lifestyleAdvice, "避免抓挠皮疹部位", "注意观察皮疹变化", "拍照记录皮疹形态,就诊时提供给医生")
		}
		
		if strings.Contains(symptomsLower, "湿疹") || strings.Contains(combinedSymptoms, "湿疹") {
			riskScore += 12
			possibleDiseases = append(possibleDiseases, "湿疹")
			diseaseScores["湿疹"] = 12.0
			detailedAnalysis["皮肤症状"] = "湿疹是常见的慢性皮肤病,可能与过敏体质、环境因素有关。建议就诊皮肤科,进行规范治疗。"
			if recommendedDept == "全科" {
				recommendedDept = "皮肤科"
			}
			lifestyleAdvice = append(lifestyleAdvice, "避免接触刺激物", "保持皮肤湿润,可使用保湿霜", "穿着纯棉衣物")
		}
		
		if strings.Contains(symptomsLower, "过敏") || strings.Contains(combinedSymptoms, "过敏") {
			riskScore += 10
			possibleDiseases = append(possibleDiseases, "过敏反应")
			diseaseScores["过敏反应"] = 10.0
			detailedAnalysis["过敏症状"] = "过敏反应可能由食物、药物、环境因素引起。建议就诊皮肤科或变态反应科,进行过敏原检测。"
			if recommendedDept == "全科" {
				recommendedDept = "皮肤科/变态反应科"
			}
			lifestyleAdvice = append(lifestyleAdvice, "记录可能接触的过敏原", "停止使用可疑物品", "如出现呼吸困难、面部肿胀,立即就医")
		}
	}
	
	// 确保风险分数在0-100范围内,并根据症状数量进行加权调整
	if symptomCount > 0 {
		// 多症状协同效应:症状越多,风险指数增长
		if symptomCount >= 4 {
			riskScore *= 1.2 // 4个及以上症状,风险提升20%
		} else if symptomCount >= 3 {
			riskScore *= 1.1 // 3个症状,风险提升10%
		}
	}
	
	// 年龄额外风险评估
	if age > 70 && riskScore > 30 {
		riskScore *= 1.15 // 70岁以上高危患者风险再提升15%
	}
	
	// 限制在0-100范围
	if riskScore > 100 {
		riskScore = 100
	}
	finalRiskScore := int(math.Round(riskScore))
	
	// 生成诊断建议 - 改进版:按评分排序疾病
	var diagnosis string
	if len(possibleDiseases) > 0 {
		// 如果有疾病评分,按分数排序
		if len(diseaseScores) > 0 {
			sortedDiseases := make([]string, 0)
			for _, disease := range possibleDiseases {
				if score, exists := diseaseScores[disease]; exists && score > 15 {
					sortedDiseases = append(sortedDiseases, disease)
				}
			}
			if len(sortedDiseases) > 0 {
				diagnosis = "AI初步分析:根据您的症状,可能存在以下情况:" + strings.Join(sortedDiseases, "、") + "。"
			} else {
				diagnosis = "AI初步分析:根据您的症状,可能存在:" + strings.Join(possibleDiseases, "、") + "。"
			}
		} else {
			diagnosis = "AI初步分析:根据您的症状,可能存在以下情况:" + strings.Join(possibleDiseases, "、") + "。"
		}
		if len(possibleDiseases) > 1 {
			diagnosis += "注意:多种症状共存可能增加疾病复杂性,建议尽快就医。"
		}
	} else {
		diagnosis = "AI初步分析:基于您提供的信息,症状相对较轻,暂未发现明显异常。但如果您感觉不适,仍建议就医进行专业评估。"
	}
	
	// 生成就医建议 - 改进版:更精确的风险分级
	var suggestions string
	if isEmergency || urgencyLevel == "emergency" {
		suggestions = "⚠️ 紧急提醒:您的症状可能存在严重风险!建议立即拨打120急救电话,或立即前往最近医院急诊科。在等待救护车期间,请保持呼吸道畅通,不要移动患者。"
	} else if finalRiskScore >= 70 || urgencyLevel == "urgent" {
		suggestions = fmt.Sprintf("🚑 高度关注:您的症状风险评分为%d分,建议尽快就医(推荐科室:%s)。请在24-48小时内前往医院进行专业检查和治疗。就诊前请避免剧烈运动,注意休息。如症状突然加重,请立即就医。", finalRiskScore, recommendedDept)
	} else if finalRiskScore >= 40 || urgencyLevel == "attention" {
		suggestions = fmt.Sprintf("🏥 建议就诊:您的症状风险评分为%d分,需要引起重视。建议在3-7天内前往%s就诊,进行进一步检查。同时注意观察症状变化,如症状加重或出现新症状,请提前就医。", finalRiskScore, recommendedDept)
	} else if finalRiskScore >= 20 {
		suggestions = fmt.Sprintf("📝 注意观察:您的症状风险评分为%d分,建议注意观察症状变化。如症状持续超过3天或逐渐加重,请及时就医。同时请保持良好的生活作息,注意饮食健康,增强免疫力。", finalRiskScore)
	} else {
		suggestions = "✅ 持续关注:您的症状相对较轻,建议多休息、多喝水,保持良好的生活习惯。如症状持续不缓解或您感觉不适,建议就医咨询。预防胜于治疗,保持健康的生活方式很重要。"
	}
	
	// 添加复诊建议
	if followUpAdvice == "" {
		if urgencyLevel == "urgent" || urgencyLevel == "emergency" {
			followUpAdvice = "建议就诊后严格遵医嘱,按时复诊和用药。"
		} else if urgencyLevel == "attention" {
			followUpAdvice = "建议就诊后1-2周进行复诊,评估治疗效果。"
		} else {
			followUpAdvice = "如症状改善不明显,建议1个月内复诊。保持健康的生活方式。"
		}
	}
	
	// 添加通用健康建议
	if len(lifestyleAdvice) == 0 {
		lifestyleAdvice = append(lifestyleAdvice, "保持规律作息，每天睡眠7-8小时", "均衡饮食，多吃蔬菜水果", "适量运动，每周至少3次，每次30分钟", "保持心情愉悦，学会减压")
	}
		
	// ========== 计算诊断置信度和数据完整度 ==========
		
	// 1. 数据完整度评估
	dataPoints := 0
	filledPoints := 0
		
	if age > 0 {
		dataPoints++
		filledPoints++
	}
	if bloodPressure != "" {
		dataPoints++
		filledPoints++
	}
	if heartRate > 0 {
		dataPoints++
		filledPoints++
	}
	if temperature > 0 {
		dataPoints++
		filledPoints++
	}
	if bloodSugar > 0 {
		dataPoints++
		filledPoints++
	}
	if otherSymptoms != "" {
		dataPoints++
		filledPoints++
	}
	dataPoints = 6 // 总共6个关键数据点
		
	dataCompleteness := float64(filledPoints) / float64(dataPoints)
		
	// 2. 诊断置信度评估
	confidence := 0.0
		
	// 基础置信度（基于数据完整度）
	confidence = dataCompleteness * 0.4 // 最多40%
		
	// 症状数量加成
	if symptomCount >= 4 {
		confidence += 0.25
	} else if symptomCount >= 3 {
		confidence += 0.15
	} else if symptomCount >= 2 {
		confidence += 0.10
	}
		
	// 有明确疾病诊断加成
	if len(possibleDiseases) > 0 {
		confidence += 0.15
	}
		
	// 症状组合匹配加成
	if _, exists := detailedAnalysis["症状组合分析"]; exists {
		confidence += 0.10
	}
	if _, exists := detailedAnalysis["高危组合"]; exists {
		confidence += 0.15
	}
		
	// 限制在0-1范围
	if confidence > 1.0 {
		confidence = 1.0
	}
		
	// 如果数据不完整，降低置信度
	if dataCompleteness < 0.5 {
		confidence *= 0.7 // 数据不足时，置信度打7折
	}
	
	return AIResult{
		RiskScore:        finalRiskScore,
		Diagnosis:        diagnosis,
		Suggestions:      suggestions,
		PossibleDiseases: possibleDiseases,
		RecommendedDept:  recommendedDept,
		UrgencyLevel:     urgencyLevel,
		DetailedAnalysis: detailedAnalysis,
		LifestyleAdvice:  lifestyleAdvice,
		FollowUpAdvice:   followUpAdvice,
		Confidence:       math.Round(confidence*100) / 100, // 保留2位小数
		DataCompleteness: math.Round(dataCompleteness*100) / 100,
	}
}
