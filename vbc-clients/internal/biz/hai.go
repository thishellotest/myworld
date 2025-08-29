package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"sync"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
)

type HaiUsecase struct {
	log                                 *log.Helper
	CommonUsecase                       *CommonUsecase
	conf                                *conf.Data
	AzopenaiUsecase                     *AzopenaiUsecase
	LogUsecase                          *LogUsecase
	LockDiseaseNamesByMedicalTextWithAI sync.Mutex
	Awsclaude3Usecase                   *Awsclaude3Usecase
}

func NewHaiUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	AzopenaiUsecase *AzopenaiUsecase,
	LogUsecase *LogUsecase,
	Awsclaude3Usecase *Awsclaude3Usecase) *HaiUsecase {
	uc := &HaiUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		AzopenaiUsecase:   AzopenaiUsecase,
		LogUsecase:        LogUsecase,
		Awsclaude3Usecase: Awsclaude3Usecase,
	}
	return uc
}

type GetDiseaseNamesByMedicalTextWithAIResponse struct {
	Conditions []string `json:"Conditions"`
}

// GetDiseaseNamesByMedicalTextWithAI fromUniqKey 使用gpt的来源标识，用来后续分析
func (c *HaiUsecase) GetDiseaseNamesByMedicalTextWithAI(ctx context.Context, medicalText string, fromUniqKey string) (response GetDiseaseNamesByMedicalTextWithAIResponse, err error) {

	if false && configs.IsDev() {
		str := `{"diseaseNames":["tuberculosis","Scarlet fever","erysipelas","Rheumatic fever","Swollen or painful joints","Frequent or severe headache","Dizziness or fainting spells","Eye trouble","Cramps in your legs","Frequent indigestion","Stomach, liver, or intestinal trouble","Gall bladder trouble","gallstones","Jaundice","hepatitis","Trick or locked knee","Foot trouble","Neuritis","Paralysis","Epilepsy or fits","Adverse reaction to serum, drug, or medicine","Ear, nose, or throat trouble","Hearing loss","Chronic or frequent colds","Severe tooth or gum trouble","Sinusitis","Hay Fever","Head Injury","Skin diseases","Thyroid trouble","Tuberculosis","Asthma","Shortness of breath","Pain or pressure In chest","Chronic cough","Palpitation or pounding heart","Heart trouble","High or low blood pressure","Broken bones","Tumor, growth, cyst, cancer","Rupture","hernia","Piles","rectal disease","Frequent or painful urination","Bed wetting since age 12","Kidney stone","blood in urine","Sugar or albumin in urine","VD-Syphilis","gonorrhea","Recent gain or loss of weight","Arthritis","Rheumatism","Bursitis","Bone, Joint or other deformity","Lameness","Loss of finger or toe","Painful or trick shoulder or elbow","Recurrent back pain","Car, train, sea or air sickness","Frequent trouble sleeping","Depression or excessive worry","Loss of memory or amnesia","Nervous trouble of any sort","Periods of unconsciousness","female disorder","change in menstrual pattern"]}`
		json.Unmarshal([]byte(str), &response)
		return
	}

	//c.LockDiseaseNamesByMedicalTextWithAI.Lock()
	//defer c.LockDiseaseNamesByMedicalTextWithAI.Unlock()
	// 防止AI使用超限
	time.Sleep(time.Microsecond)
	c.log.Debug("GetDiseaseNamesByMedicalTextWithAI begin: ", fromUniqKey)
	systemConfig := `# Role
你是一位全领域的医疗专家，精通分析复杂的医疗数据，能够从大量数据中提取有价值的信息。擅长高效地提取并分析所有疾病名称，为患者和医疗专业人员提供全面的医疗信息。

## Skills
- 初级保健：家庭医学、一般诊断和健康检查
- 内科：心血管疾病、消化系统疾病、内分泌和代谢疾病
- 外科：普通外科手术、微创手术和创伤外科
- 儿科：新生儿护理、儿童疾病和青少年健康
- 妇产科：妇科疾病、产前和产后护理、生殖健康
- 皮肤科：皮肤疾病、皮肤护理和美容治疗
- 心脏病学：冠心病、高血压和心脏康复
- 神经学：神经系统疾病、脑卒中和癫痫
- 急救医学：急诊处理、创伤护理和紧急手术
- 精神健康：心理咨询、精神疾病和心理健康管理
- 数据分析：医疗数据挖掘、疾病名称提取和医疗信息整理

## Action
- 根据患者的描述，提取出所有疾病名称，并以JSON格式输出

## Constrains
- 忽略无关内容
- 疾病名称必须与提供内容完全匹配
- 必须保证你的结果只包含一个合法的JSON格式

## Format
- 对应JSON的key为：diseaseNames`

	systemConfig = `# Role
You are a comprehensive medical expert proficient in analyzing complex medical data, capable of extracting valuable information from vast datasets. You are also adept at efficiently extracting and analyzing all disease names, providing comprehensive medical information to patients and healthcare professionals.

## Skills
- Primary Care: Family medicine, general diagnostics, and health check-ups
- Internal Medicine: Cardiovascular diseases, gastrointestinal disorders, endocrine, and metabolic diseases
- Surgery: General surgery, minimally invasive surgery, and trauma surgery
- Obstetrics and Gynecology: Gynecological diseases, prenatal and postnatal care, reproductive health
- Dermatology: Skin diseases, skincare, and cosmetic treatments
- Cardiology: Coronary artery disease, hypertension, and cardiac rehabilitation
- Neurology: Neurological disorders, stroke, and epilepsy
- Emergency Medicine: Emergency treatment, trauma care, and urgent surgery
- Mental Health: Psychological counseling, mental disorders, and mental health management
- Data Analysis: Medical data mining, disease name extraction, and medical information organization

## Action
- Based on the patient's description, extract all disease names and output them in JSON format

## Constraints
- Ignore irrelevant content
- The case of the found disease name must match exactly the content provided
- Ensure your result contains only one valid JSON format

## Format
- The corresponding JSON key is: diseaseNames`

	systemConfig = `# Role
You are a comprehensive medical expert proficient in analyzing complex medical data, capable of extracting valuable information from vast datasets. You are also adept at efficiently extracting and analyzing all disease names, providing comprehensive medical information to patients and healthcare professionals.

## Skills
- Primary Care: Family medicine, general diagnostics, and health check-ups
- Internal Medicine: Cardiovascular diseases, gastrointestinal disorders, endocrine, and metabolic diseases
- Surgery: General surgery, minimally invasive surgery, and trauma surgery
- Obstetrics and Gynecology: Gynecological diseases, prenatal and postnatal care, reproductive health
- Dermatology: Skin diseases, skincare, and cosmetic treatments
- Cardiology: Coronary artery disease, hypertension, and cardiac rehabilitation
- Neurology: Neurological disorders, stroke, and epilepsy
- Emergency Medicine: Emergency treatment, trauma care, and urgent surgery
- Mental Health: Psychological counseling, mental disorders, and mental health management
- Data Analysis: Medical data mining, disease name extraction, and medical information organization

## Action
- Based on the patient's description, extract all disease names and output them in JSON format

## Constraints
- Disease names must exactly match the provided content
- Ensure your result contains only one valid JSON format

## Format
- The corresponding JSON key is: Conditions`

	systemConfig = `# Role
You are a comprehensive medical expert proficient in analyzing complex medical data, capable of extracting valuable information from vast datasets. You are also adept at efficiently extracting and analyzing all disease names, providing comprehensive medical information to patients and healthcare professionals.

## Skills
- Primary Care: Family medicine, general diagnostics, and health check-ups
- Internal Medicine: Cardiovascular diseases, gastrointestinal disorders, endocrine, and metabolic diseases
- Surgery: General surgery, minimally invasive surgery, and trauma surgery
- Obstetrics and Gynecology: Gynecological diseases, prenatal and postnatal care, reproductive health
- Dermatology: Skin diseases, skincare, and cosmetic treatments
- Cardiology: Coronary artery disease, hypertension, and cardiac rehabilitation
- Neurology: Neurological disorders, stroke, and epilepsy
- Emergency Medicine: Emergency treatment, trauma care, and urgent surgery
- Mental Health: Psychological counseling, mental disorders, and mental health management
- Data Analysis: Medical data mining, disease name extraction, and medical information organization

## Action
- Based on the patient's description, extract all disease names and output them in JSON format

## Constraints
- The case of the found disease name must match exactly the content provided
- Ensure your result contains only one valid JSON format

## Format
- The corresponding JSON key is: Conditions`

	return c.FromClaude3(systemConfig, medicalText, fromUniqKey)

	res, err := c.AzopenaiUsecase.AskGpt4o(ctx, systemConfig, medicalText)
	c.log.Debug("GetDiseaseNamesByMedicalTextWithAI:", " ", fromUniqKey, " ", res)
	if err != nil {
		c.log.Debug("GetDiseaseNamesByMedicalTextWithAI: err: ", " ", fromUniqKey, " ", err)
	}
	if err != nil {
		c.log.Error("fromUniqKey: ", fromUniqKey, " err: ", err.Error())
		return response, err
	}
	er := c.LogUsecase.SaveLog(0, "GetDiseaseNamesByMedicalTextWithAI", map[string]interface{}{
		"fromUniqKey": fromUniqKey,
		"result":      res,
		"err":         err,
	})
	if er != nil {
		c.log.Error(er)
	}
	newRes := strings.Replace(res, "```json", "", -1)
	newRes = strings.Replace(newRes, "```", "", -1)
	response, errIgnore := lib.StringToTE(newRes, GetDiseaseNamesByMedicalTextWithAIResponse{})
	if errIgnore != nil {
		c.LogUsecase.SaveLog(0, "GetDiseaseNamesByMedicalTextWithAIError", map[string]interface{}{
			"fromUniqKey": fromUniqKey,
			"result":      res,
			"err":         errIgnore,
		})
		c.log.Error("fromUniqKey: ", fromUniqKey, " res: ", res, " err: ", errIgnore.Error(), "")
	}
	c.log.Debug("GetDiseaseNamesByMedicalTextWithAI end: ", fromUniqKey)

	return response, err
}

func (c *HaiUsecase) FromClaude3(systemConfig string, medicalText string, fromUniqKey string) (res GetDiseaseNamesByMedicalTextWithAIResponse, err error) {

	content, err := c.Awsclaude3Usecase.GetContentByAsk(context.TODO(), systemConfig, medicalText)
	if err != nil {
		c.log.Error(err)
		return res, nil
	}
	a := strings.Index(content, "{")
	newRes := content[a:]
	response, errIgnore := lib.StringToTE(newRes, GetDiseaseNamesByMedicalTextWithAIResponse{})
	if errIgnore != nil {
		c.LogUsecase.SaveLog(0, "FromClaude3", map[string]interface{}{
			"fromUniqKey": fromUniqKey,
			"result":      res,
			"err":         errIgnore,
		})
		c.log.Error("fromUniqKey: ", fromUniqKey, " res: ", res, " err: ", errIgnore.Error(), "")
	}
	return response, nil
}
