package service

import (
	"errors"
	"log"
	"sm-medical/internal/model"
	"sm-medical/internal/repository"
)

// TriageService 智能分诊服务
type TriageService struct {
	userRepo         *repository.UserRepository
	consultationRepo *repository.ConsultationRepository
}

func NewTriageService() *TriageService {
	return &TriageService{
		userRepo:         repository.NewUserRepository(),
		consultationRepo: repository.NewConsultationRepository(),
	}
}

// TriageRule 分诊规则
type TriageRule struct {
	Priority int    // 优先级(数字越小优先级越高)
	Name     string // 规则名称
	Match    func(doctor *model.User, recommendedDept string) bool
	Score    func(doctor *model.User) int
}

// AutoAssignDoctor 自动分配医生
// recommendedDept: AI推荐的科室
func (s *TriageService) AutoAssignDoctor(consultationID int64, recommendedDept string) (*model.User, string, error) {
	log.Printf("[智能分诊] 开始为问诊 %d 分配医生,推荐科室: %s", consultationID, recommendedDept)

	// 获取所有可用医生
	doctors, err := s.getAvailableDoctors()
	if err != nil {
		return nil, "", errors.New("无可用医生")
	}

	if len(doctors) == 0 {
		return nil, "", errors.New("暂无医生在线,请稍后重试")
	}

	// 分诊规则优先级
	rules := []TriageRule{
		// 规则1: 在线 + 科室匹配 + 负载最低
		{
			Priority: 1,
			Name:     "在线科室匹配",
			Match: func(d *model.User, dept string) bool {
				return d.IsOnline == 1 && d.DoctorDept == dept && d.CurrentConsultationCount < d.MaxConsultationCount
			},
			Score: func(d *model.User) int {
				return d.CurrentConsultationCount // 负载越低分数越低,越优先
			},
		},
		// 规则2: 科室匹配(不在线但负载允许)
		{
			Priority: 2,
			Name:     "科室匹配",
			Match: func(d *model.User, dept string) bool {
				return d.DoctorDept == dept && d.CurrentConsultationCount < d.MaxConsultationCount
			},
			Score: func(d *model.User) int {
				if d.IsOnline == 1 {
					return d.CurrentConsultationCount
				}
				return d.CurrentConsultationCount + 100 // 不在线的权重降低
			},
		},
		// 规则3: 在线医生(不考虑科室,负载最低)
		{
			Priority: 3,
			Name:     "在线医生兜底",
			Match: func(d *model.User, dept string) bool {
				return d.IsOnline == 1 && d.CurrentConsultationCount < d.MaxConsultationCount
			},
			Score: func(d *model.User) int {
				return d.CurrentConsultationCount
			},
		},
		// 规则4: 任意可用医生(最后兜底)
		{
			Priority: 4,
			Name:     "任意可用医生",
			Match: func(d *model.User, dept string) bool {
				return d.CurrentConsultationCount < d.MaxConsultationCount
			},
			Score: func(d *model.User) int {
				return d.CurrentConsultationCount + 200
			},
		},
	}

	// 按优先级匹配
	for _, rule := range rules {
		matchedDoctors := make([]*model.User, 0)
		for i := range doctors {
			if rule.Match(&doctors[i], recommendedDept) {
				matchedDoctors = append(matchedDoctors, &doctors[i])
			}
		}

		if len(matchedDoctors) > 0 {
			// 在匹配的医生中选择负载最低的
			bestDoctor := matchedDoctors[0]
			bestScore := rule.Score(bestDoctor)

			for _, d := range matchedDoctors[1:] {
				score := rule.Score(d)
				if score < bestScore {
					bestDoctor = d
					bestScore = score
				}
			}

			log.Printf("[智能分诊] 匹配成功 - 规则: %s, 医生: %s(%s), 科室: %s, 负载: %d/%d",
				rule.Name, bestDoctor.RealName, bestDoctor.Username, bestDoctor.DoctorDept,
				bestDoctor.CurrentConsultationCount, bestDoctor.MaxConsultationCount)

			return bestDoctor, rule.Name, nil
		}
	}

	return nil, "", errors.New("所有医生已达最大负载,请稍后重试")
}

// getAvailableDoctors 获取可用医生列表
func (s *TriageService) getAvailableDoctors() ([]model.User, error) {
	// 获取所有正常状态的医生
	doctors, err := s.userRepo.FindByRole("doctor")
	if err != nil {
		return nil, err
	}

	// 过滤禁用的医生
	available := make([]model.User, 0)
	for _, d := range doctors {
		if d.Status == 0 { // 0表示正常
			available = append(available, d)
		}
	}

	return available, nil
}

// UpdateDoctorOnlineStatus 更新医生在线状态
func (s *TriageService) UpdateDoctorOnlineStatus(doctorID int64, isOnline bool) error {
	doctor, err := s.userRepo.FindByID(doctorID)
	if err != nil {
		return err
	}

	if isOnline {
		doctor.IsOnline = 1
	} else {
		doctor.IsOnline = 0
	}

	return s.userRepo.Update(doctor)
}

// UpdateDoctorWorkload 更新医生工作负载
func (s *TriageService) UpdateDoctorWorkload(doctorID int64, delta int) error {
	doctor, err := s.userRepo.FindByID(doctorID)
	if err != nil {
		return err
	}

	newCount := doctor.CurrentConsultationCount + delta
	if newCount < 0 {
		newCount = 0
	}

	doctor.CurrentConsultationCount = newCount
	return s.userRepo.Update(doctor)
}
