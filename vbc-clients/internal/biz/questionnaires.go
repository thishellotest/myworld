package biz

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
)

type QuestionnairesEntity struct {
	ID            int32 `gorm:"primaryKey"`
	JotformFormId string
	Title         string
	BaseTitle     string
	JsonData      string
	IsIntake      int
	CreatedAt     int64
	UpdatedAt     int64
	DeletedAt     int64
}

func (QuestionnairesEntity) TableName() string {
	return "questionnaires"
}

type QuestionnairesUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	TUsecase         *TUsecase
	DataComboUsecase *DataComboUsecase
	DBUsecase[QuestionnairesEntity]
}

func NewQuestionnairesUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	DataComboUsecase *DataComboUsecase) *QuestionnairesUsecase {
	uc := &QuestionnairesUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		DataComboUsecase: DataComboUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *QuestionnairesUsecase) HttpList(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizHttpList(body.GetInt("case_id"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func GetQuestionnairesItemByFormId(FormId string) *QuestionnairesItem {
	for k, v := range QuestionnairesListConfigs {
		if v.FormId == FormId {
			r := QuestionnairesListConfigs[k]
			return &r
		}
	}
	return nil
}

type QuestionnairesItem struct {
	Title     string   `json:"title"`
	BaseTitle string   `json:"base_title"`
	FormId    string   `json:"form_id"`
	FileNames []string `json:"file_names"`
}

const (
	QuestionnairesInitialIntake_FormId       = "240926345150149"
	QuestionnairesUpdateQuestionnaire_FormId = "250947787651169"
)

var QuestionnairesListConfigs = []QuestionnairesItem{
	{
		Title:     "Initial Intake Questionnaire",
		BaseTitle: "Initial Intake",
		FormId:    QuestionnairesInitialIntake_FormId,
		FileNames: []string{},
	},
	{
		Title:     "Ankle Questionnaire",
		BaseTitle: "Ankle",
		FormId:    "240887024317154",
		FileNames: []string{
			"content.answers.85.answer",
			"content.answers.84.answer",
		},
	},
	{
		Title:     "Artery and Vein Questionnaire",
		BaseTitle: "Artery and Vein",
		FormId:    "250707310735148",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.205.answer",
		},
	},
	{
		Title:     "Arthritis Questionnaire",
		BaseTitle: "Arthritis",
		FormId:    "241148305913149",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.205.answer",
		},
	},
	{
		Title:     "Back Questionnaire",
		BaseTitle: "Back",
		FormId:    "240905843097159",
		FileNames: []string{
			"content.answers.97.answer",
			"content.answers.98.answer",
		},
	},
	{
		Title:     "Blood Conditions Questionnaire",
		BaseTitle: "Blood Conditions",
		FormId:    "252248306909158",
		FileNames: []string{
			"content.answers.101.answer",
			"content.answers.100.answer",
		},
	},
	{
		Title:     "Bone Questionnaire",
		BaseTitle: "Bone",
		FormId:    "241128913538155",
		FileNames: []string{
			"content.answers.97.answer",
			"content.answers.98.answer",
		},
	},
	{
		Title:     "Chronic Fatigue Syndrome Questionnaire",
		BaseTitle: "Chronic Fatigue Syndrome",
		FormId:    "251474489842166",
		FileNames: []string{
			"content.answers.15.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Diabetes Questionnaire",
		BaseTitle: "Diabetes",
		FormId:    "240906996093165",
		FileNames: []string{
			"content.answers.101.answer",
			"content.answers.100.answer",
		},
	},
	{
		Title:     "Ear Questionnaire",
		BaseTitle: "Ear",
		FormId:    "241150533456147",
		FileNames: []string{
			"content.answers.89.answer",
			"content.answers.91.answer",
		},
	},
	{
		Title:     "Elbow Questionnaire",
		BaseTitle: "Elbow",
		FormId:    "240917935177163",
		FileNames: []string{
			"content.answers.172.answer",
			"content.answers.173.answer",
		},
	},
	{
		Title:     "Epilepsy Questionnaire",
		BaseTitle: "Epilepsy",
		FormId:    "241499010039150",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Esophageal Conditions Questionnaire",
		BaseTitle: "Esophageal Conditions",
		FormId:    "240908764581162",
		FileNames: []string{
			"content.answers.15.answer",
			"content.answers.71.answer",
		},
	},
	{
		Title:     "Eye Questionnaire",
		BaseTitle: "Eye",
		FormId:    "241499180510152",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.205.answer",
		},
	},
	{
		Title:     "Foot Questionnaire",
		BaseTitle: "Foot",
		FormId:    "240908012657152",
		FileNames: []string{
			"content.answers.100.answer",
			"content.answers.101.answer",
		},
	},
	{
		Title:     "Gallbladder Questionnaire",
		BaseTitle: "Gallbladder",
		FormId:    "250707993189168",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.108.answer",
		},
	},
	{
		Title:     "Gynecological Conditions Questionnaire",
		BaseTitle: "Gynecological Conditions",
		FormId:    "241498350441153",
		FileNames: []string{
			"content.answers.101.answer",
			"content.answers.100.answer",
		},
	},
	{
		Title:     "Hand and Finger Questionnaire",
		BaseTitle: "Hand and Finger",
		FormId:    "240919081022146",
		FileNames: []string{
			"content.answers.192.answer",
			"content.answers.193.answer",
		},
	},
	{
		Title:     "Headaches and Migraines Questionnaire",
		BaseTitle: "Headaches and Migraines",
		FormId:    "240808698395069",
		FileNames: []string{
			"content.answers.53.answer",
			"content.answers.52.answer",
		},
	},
	{
		Title:     "Hearing Loss and Tinnitus Questionnaire",
		BaseTitle: "Hearing Loss and Tinnitus",
		FormId:    "240917546196162",
		FileNames: []string{
			"content.answers.89.answer",
			"content.answers.91.answer",
		},
	},
	{
		Title:     "Heart Questionnaire",
		BaseTitle: "Heart",
		FormId:    "240924069988169",
		FileNames: []string{
			"content.answers.211.answer",
			"content.answers.212.answer",
		},
	},
	{
		Title:     "Hernia Questionnaire",
		BaseTitle: "Hernia",
		FormId:    "242387777008163",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.205.answer",
		},
	},
	{
		Title:     "Hip and Thigh Questionnaire",
		BaseTitle: "Hip and Thigh",
		FormId:    "240911239839159",
		FileNames: []string{
			"content.answers.115.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Hypertension Questionnaire",
		BaseTitle: "Hypertension",
		FormId:    "240911062260141",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Intestinal Conditions Questionnaire",
		BaseTitle: "Intestinal Conditions",
		FormId:    "240914628058157",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.108.answer",
		},
	},
	{
		Title:     "Kidney Questionnaire",
		BaseTitle: "Kidney",
		FormId:    "242936332167156",
		FileNames: []string{
			"content.answers.101.answer",
			"content.answers.100.answer",
		},
	},
	{
		Title:     "Knee Questionnaire",
		BaseTitle: "Knee",
		FormId:    "240915419171152",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Male Reproductive Organ Conditions Questionnaire",
		BaseTitle: "Male Reproductive Organ Conditions",
		FormId:    "240916181206148",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.136.answer",
		},
	},
	{
		Title:     "Mental Disorders Questionnaire",
		BaseTitle: "Mental Disorders",
		FormId:    "240916451227151",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Mental Disorders Secondaries Questionnaire",
		BaseTitle: "Mental Disorders Secondaries",
		FormId:    "240996235351157",
		FileNames: []string{
			"content.answers.150.answer",
			"content.answers.100.answer",
		},
	},
	{
		Title:     "Miscellaneous Questionnaire",
		BaseTitle: "Miscellaneous",
		FormId:    "240925950270153",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.205.answer",
		},
	},
	{
		Title:     "Muscle Injuries Questionnaire", // 生成的pdf没有ID, 所以在文件名上加了ID，需要特殊处理
		BaseTitle: "Muscle Injuries",
		FormId:    "240920932879163",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.203.answer",
		},
	}, {
		Title:     "Neck Questionnaire",
		BaseTitle: "Neck",
		FormId:    "240916250200139",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.97.answer",
		},
	},
	{
		Title:     "Nerve Questionnaire",
		BaseTitle: "Nerve",
		FormId:    "240921025199152",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.171.answer",
		},
	},
	{
		Title:     "Post-Traumatic Stress Disorder (PTSD) Questionnaire", // 没有 claiming
		BaseTitle: "Post-Traumatic Stress Disorder (PTSD)",
		FormId:    "240916177462157",
		FileNames: []string{
			"content.answers.178.answer",
		},
	},
	{
		Title:     "Prostate Cancer Questionnaire",
		BaseTitle: "Prostate Cancer",
		FormId:    "240926220533146",
		FileNames: []string{
			"content.answers.140.answer",
			"content.answers.141.answer",
		},
	},
	{
		Title:     "Rectum and Anus Questionnaire",
		BaseTitle: "Rectum and Anus",
		FormId:    "240924814665159",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.239.answer",
		},
	},
	{
		Title:     "Respiratory Questionnaire",
		BaseTitle: "Respiratory",
		FormId:    "240916210529149",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Scar Questionnaire",
		BaseTitle: "Scar",
		FormId:    "240996877765178",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Shoulder Questionnaire",
		BaseTitle: "Shoulder",
		FormId:    "240916790171155",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Sinus Questionnaire",
		BaseTitle: "Sinus",
		FormId:    "240917105856156",
		FileNames: []string{
			"content.answers.15.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Skin Questionnaire",
		BaseTitle: "Skin",
		FormId:    "240917671665162",
		FileNames: []string{
			"content.answers.89.answer",
			"content.answers.91.answer",
		},
	},
	{
		Title:     "Sleep Apnea Questionnaire",
		BaseTitle: "Sleep Apnea",
		FormId:    "240917398302156",
		FileNames: []string{
			"content.answers.89.answer",
			"content.answers.128.answer",
		},
	},
	{
		Title:     "Stomach and Duodenal Conditions",
		BaseTitle: "Stomach and Duodenal Conditions",
		FormId:    "252089144305151",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.108.answer",
		},
	},
	{
		Title:     "Temporomandibular Joint (TMJ) Disorder Questionnaire",
		BaseTitle: "Temporomandibular Joint (TMJ) Disorder",
		FormId:    "240923742537156",
		FileNames: []string{
			"content.answers.208.answer",
			"content.answers.11.answer",
		},
	},
	{
		Title:     "Thyroid and Parathyroid Questionnaire",
		BaseTitle: "Thyroid and Parathyroid",
		FormId:    "242415464777161",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.205.answer",
		},
	},
	{
		Title:     "Total Disability Individual Unemployability (TDIU) Questionnaire", // 没有New Or Increase
		BaseTitle: "Total Disability Individual Unemployability (TDIU)",
		FormId:    "241718641709158",
		FileNames: []string{
			"content.answers.205.answer",
		},
	},
	{
		Title:     "Traumatic Brain Injury (TBI) Questionnaire",
		BaseTitle: "Traumatic Brain Injury (TBI)",
		FormId:    "241497331060148",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.205.answer",
		},
	},
	{
		Title:     "Update Questionnaire",
		BaseTitle: "Update Questionnaire",
		FormId:    "250947787651169",
		FileNames: []string{},
	},
	{
		Title:     "Urinary Tract Conditions Questionnaire",
		BaseTitle: "Urinary Tract Conditions",
		FormId:    "241147117607149",
		FileNames: []string{
			"content.answers.140.answer",
			"content.answers.141.answer",
		},
	},
	{
		Title:     "Wrist Questionnaire",
		BaseTitle: "Wrist",
		FormId:    "240918258523157",
		FileNames: []string{
			"content.answers.52.answer",
			"content.answers.190.answer",
		},
	},
}

func (c *QuestionnairesUsecase) BizHttpList(caseId int32) (lib.TypeMap, error) {

	if caseId <= 0 {
		return nil, errors.New("Parameter error")
	}
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}

	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return nil, err
	}

	var questionnaires lib.TypeList
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Initial Intake Questionnaire",
		"link":  c.LinkInitialIntake(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Ankle Questionnaire",
		"link":  c.LinkAnkleQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Artery and Vein Questionnaire",
		"link":  c.LinkArteryAndVeinQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Arthritis Questionnaire",
		"link":  c.LinkArthritisQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Back Questionnaire",
		"link":  c.LinkBackQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Blood Conditions Questionnaire",
		"link":  c.LinkBloodConditionsQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Bone Questionnaire",
		"link":  c.LinkBoneQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Chronic Fatigue Syndrome Questionnaire",
		"link":  c.LinkChronicFatigueSyndromeQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Diabetes Questionnaire",
		"link":  c.LinkDiabetesQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Ear Questionnaire",
		"link":  c.LinkEarQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Elbow Questionnaire",
		"link":  c.LinkElbowQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Epilepsy Questionnaire",
		"link":  c.LinkEpilepsyQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Esophageal Conditions Questionnaire",
		"link":  c.LinkEsophagealConditionsQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Eye Questionnaire",
		"link":  c.LinkEyeQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Foot Questionnaire",
		"link":  c.LinkFootQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Gallbladder Questionnaire",
		"link":  c.LinkGallbladderQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Gynecological Conditions Questionnaire",
		"link":  c.LinkGynecologicalConditionsQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Hand and Finger Questionnaire",
		"link":  c.LinkHandAndFingerQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Headaches and Migraines Questionnaire",
		"link":  c.LinkHeadachesAndMigrainesQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Hearing Loss and Tinnitus Questionnaire",
		"link":  c.LinkHearingLossAndTinnitusQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Heart Questionnaire",
		"link":  c.LinkHeartQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Hernia Questionnaire",
		"link":  c.LinkHerniaQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Hip and Thigh Questionnaire",
		"link":  c.LinkHipAndThighQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Hypertension Questionnaire",
		"link":  c.LinkHypertensionQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Intestinal Conditions Questionnaire",
		"link":  c.LinkIBSQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Kidney Questionnaire",
		"link":  c.LinkKidneyQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Knee Questionnaire",
		"link":  c.LinkKneeQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Male Reproductive Organ Conditions Questionnaire",
		"link":  c.LinkMaleReproductiveOrganConditionsQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Mental Disorders Questionnaire",
		"link":  c.LinkMentalDisordersQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Mental Disorders Secondaries Questionnaire",
		"link":  c.LinkMentalDisordersSecondariesQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Miscellaneous Questionnaire",
		"link":  c.LinkMiscellaneousQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Muscle Injuries Questionnaire",
		"link":  c.LinkMuscleInjuriesQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Neck Questionnaire",
		"link":  c.LinkNeckQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Nerve Questionnaire",
		"link":  c.LinkNerveQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Post-Traumatic Stress Disorder (PTSD) Questionnaire",
		"link":  c.LinkPTSDQuestionnaire(*tClient, *tCase),
	})

	//questionnaires = append(questionnaires, map[string]interface{}{
	//	"title": "Post-Traumatic Stress Disorder (PTSD) Secondary to Personal Assault",
	//	"link":  c.LinkPTSDSecondaryQuestionnaire(*tClient, *tCase),
	//})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Prostate Cancer Questionnaire",
		"link":  c.LinkProstateCancerQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Rectum and Anus Questionnaire",
		"link":  c.LinkRectumAndAnusQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Respiratory Questionnaire",
		"link":  c.LinkRespiratoryQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Scar Questionnaire",
		"link":  c.LinkScarQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Shoulder Questionnaire",
		"link":  c.LinkShoulderQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Sinus Questionnaire",
		"link":  c.LinkSinusQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Skin Questionnaire",
		"link":  c.LinkSkinQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Sleep Apnea Questionnaire",
		"link":  c.LinkSleepApneaQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Stomach and Duodenal Conditions",
		"link":  c.LinkStomachAndDuodenalConditions(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Temporomandibular Joint (TMJ) Disorder Questionnaire",
		"link":  c.LinkTMJQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Thyroid and Parathyroid Questionnaire",
		"link":  c.LinkThyroidAndParathyroidQuestionnaire(*tClient, *tCase),
	})

	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Total Disability Individual Unemployability (TDIU) Questionnaire",
		"link":  c.LinkTDIUQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Traumatic Brain Injury (TBI) Questionnaire",
		"link":  c.LinkTBIQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Update Questionnaire",
		"link":  c.LinkUpdateQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Urinary Tract Conditions Questionnaire",
		"link":  c.LinkUrinaryTractConditionsQuestionnaire(*tClient, *tCase),
	})
	questionnaires = append(questionnaires, map[string]interface{}{
		"title": "Wrist Questionnaire",
		"link":  c.LinkWristQuestionnaire(*tClient, *tCase),
	})

	data := make(lib.TypeMap)
	data.Set("case.deal_name", tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name))
	data.Set("questionnaires", questionnaires)
	return data, nil
}

const (
	JotformBaseUrl = "https://form.jotform.com"
)

// LinkInitialIntake ok
func (c *QuestionnairesUsecase) LinkInitialIntake(tClient TData, tCase TData) string {
	//Initial Intake
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	branch := tCase.CustomFields.TextValueByNameBasic("branch")
	retired := tCase.CustomFields.TextValueByNameBasic("retired")
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/%s?fullName[first]=%s&fullName[last]=%s&uniqueId=%d&branchOf=%s&didYou=%s",
		QuestionnairesInitialIntake_FormId, firstName, lastName, uniqueId, branch, retired)
	return url
}

// LinkAnkleQuestionnaire ok
func (c *QuestionnairesUsecase) LinkAnkleQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240887024317154?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

func (c *QuestionnairesUsecase) LinkArteryAndVeinQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/250707310735148?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

// LinkArthritisQuestionnaire ok
func (c *QuestionnairesUsecase) LinkArthritisQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/241148305913149?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

// LinkBackQuestionnaire ok
func (c *QuestionnairesUsecase) LinkBackQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240905843097159?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

func (c *QuestionnairesUsecase) LinkBloodConditionsQuestionnaire(tClient TData, tCase TData) string {

	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/252248306909158?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

// LinkBoneQuestionnaire ok
func (c *QuestionnairesUsecase) LinkBoneQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/241128913538155?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

// LinkChronicFatigueSyndromeQuestionnaire ok
func (c *QuestionnairesUsecase) LinkChronicFatigueSyndromeQuestionnaire(tClient TData, tCase TData) string {
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/251474489842166?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

// LinkDiabetesQuestionnaire ok
func (c *QuestionnairesUsecase) LinkDiabetesQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240906996093165?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

// LinkEarQuestionnaire ok
func (c *QuestionnairesUsecase) LinkEarQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/241150533456147?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

// LinkElbowQuestionnaire ok
func (c *QuestionnairesUsecase) LinkElbowQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240917935177163?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

// LinkEpilepsyQuestionnaire ok
func (c *QuestionnairesUsecase) LinkEpilepsyQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/241499010039150?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

// LinkEsophagealConditionsQuestionnaire ok
func (c *QuestionnairesUsecase) LinkEsophagealConditionsQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240908764581162?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

// LinkEyeQuestionnaire ok
func (c *QuestionnairesUsecase) LinkEyeQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/241499180510152?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkFootQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240908012657152?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkGallbladderQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/250707993189168?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkGynecologicalConditionsQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/241498350441153?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkHandAndFingerQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240919081022146?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkHeadachesAndMigrainesQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240808698395069?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkHearingLossAndTinnitusQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240917546196162?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkHeartQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240924069988169?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkHerniaQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/242387777008163?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkHipAndThighQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240911239839159?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkHypertensionQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240911062260141?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkIBSQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240914628058157?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkKidneyQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/242936332167156?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkKneeQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240915419171152?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkMaleReproductiveOrganConditionsQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240916181206148?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkMentalDisordersQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240916451227151?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkMentalDisordersSecondariesQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240996235351157?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkMiscellaneousQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240925950270153?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkMuscleInjuriesQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240920932879163?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkNeckQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240916250200139?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkNerveQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240921025199152?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkPTSDQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240916177462157?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkPTSDSecondaryQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/241095451969163?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkProstateCancerQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240926220533146?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkRectumAndAnusQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240924814665159?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkRespiratoryQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240916210529149?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkScarQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240996877765178?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkShoulderQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240916790171155?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkSinusQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240917105856156?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkSkinQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240917671665162?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkSleepApneaQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240917398302156?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkStomachAndDuodenalConditions(tClient TData, tCase TData) string {
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/252089144305151?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkTMJQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240923742537156?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkThyroidAndParathyroidQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/242415464777161?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkTDIUQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/241718641709158?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkTBIQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/241497331060148?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

func (c *QuestionnairesUsecase) LinkUpdateQuestionnaire(tClient TData, tCase TData) string {
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/250947787651169?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}

func (c *QuestionnairesUsecase) LinkUrinaryTractConditionsQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/241147117607149?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
func (c *QuestionnairesUsecase) LinkWristQuestionnaire(tClient TData, tCase TData) string {
	//Initial Intake
	// https://form.jotform.com/240926345150149
	// https://hipaa.jotform.com
	firstName := tClient.CustomFields.TextValueByNameBasic(FieldName_first_name)
	lastName := tClient.CustomFields.TextValueByNameBasic(FieldName_last_name)
	uniqueId := tCase.Id()
	url := fmt.Sprintf("https://form.jotform.com/240918258523157?fullName[first]=%s&fullName[last]=%s&uniqueId=%d", firstName, lastName, uniqueId)
	return url
}
