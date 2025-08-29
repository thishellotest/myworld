package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib"
)

const (
	Psychiatric_DBQ_Fee     = 400
	Medical_DBQ_Base_Fee    = 450 // includes up to 3 DBQs(少于等于3个时)
	Medical_DBQ_Per_Fee     = 100 // 每增加一个+100的费用
	Opinion_Letters_Per_Fee = 75  // $75 per letter
)

type MedicalDbqCostFee struct {
	PsychiatricDBQFee    int `json:"psychiatric_dbq_fee"`
	MedicalDBQBaseFee    int `json:"medical_dbq_base_fee"`
	MedicalDBQPerFee     int `json:"medical_dbq_per_fee"`
	OpinionLettersPerFee int `json:"opinion_letters_per_fee"`
}

func MedicalDbqCostFeeConfig() (medicalDbqCostFee MedicalDbqCostFee) {
	medicalDbqCostFee.PsychiatricDBQFee = Psychiatric_DBQ_Fee
	medicalDbqCostFee.MedicalDBQBaseFee = Medical_DBQ_Base_Fee
	medicalDbqCostFee.MedicalDBQPerFee = Medical_DBQ_Per_Fee
	medicalDbqCostFee.OpinionLettersPerFee = Opinion_Letters_Per_Fee
	return medicalDbqCostFee
}

type MedicalDbqCost struct {
	FeeConfig MedicalDbqCostFee      `json:"fee_config"`
	UserInfo  MedicalDbqCostUserInfo `json:"user_info"`
}

type MedicalDbqCostUserInfo struct {
	HasPsychiatricDBQ   bool `json:"has_psychiatric_dbq"`
	MedicalDBQCount     int  `json:"medical_dbq_count"`
	OpinionLettersCount int  `json:"opinion_letters_count"`
}

type MedicalDbqCostUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	MapUsecase    *MapUsecase
}

func NewMedicalDbqCostUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	MapUsecase *MapUsecase,
) *MedicalDbqCostUsecase {
	uc := &MedicalDbqCostUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		MapUsecase:    MapUsecase,
	}
	return uc
}

func (c *MedicalDbqCostUsecase) GetMedicalDbqCost(caseId int32) (medicalDbqCost MedicalDbqCost) {
	medicalDbqCost.FeeConfig = MedicalDbqCostFeeConfig()
	medicalDbqCost.UserInfo = c.MedicalDbqCostUserInfoByCaseId(caseId)
	return medicalDbqCost
}

func (c *MedicalDbqCostUsecase) SaveMedicalDbqCostUserInfo(caseId int32, medicalDbqCostUserInfo MedicalDbqCostUserInfo) {
	key := MapKeyMedicalDbqCost(caseId)
	c.MapUsecase.Set(key, InterfaceToString(medicalDbqCostUserInfo))
}

func (c *MedicalDbqCostUsecase) MedicalDbqCostUserInfoByCaseId(caseId int32) MedicalDbqCostUserInfo {

	var medicalDbqCostUserInfo MedicalDbqCostUserInfo
	key := MapKeyMedicalDbqCost(caseId)
	val, _ := c.MapUsecase.GetForString(key)
	if val != "" {
		medicalDbqCostUserInfo = lib.StringToTDef(val, MedicalDbqCostUserInfo{})
	}
	return medicalDbqCostUserInfo
}

// MedicalDbqCostCalculator 这个公式放到前端计算
func MedicalDbqCostCalculator(hasPsychiatricDBQ bool, medicalDBQCount int, opinionLettersCount int) int {
	total := 0
	if hasPsychiatricDBQ {
		total += Psychiatric_DBQ_Fee
	}
	if medicalDBQCount > 0 {
		if medicalDBQCount <= 3 {
			total += Medical_DBQ_Base_Fee
		} else if medicalDBQCount > 3 {
			total += Medical_DBQ_Base_Fee
			total += (medicalDBQCount - 3) * Medical_DBQ_Per_Fee
		}
	}
	if opinionLettersCount > 0 {
		total += opinionLettersCount * Opinion_Letters_Per_Fee
	}
	return total
}
