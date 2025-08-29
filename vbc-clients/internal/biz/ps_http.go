package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
)

type PsHttpUsecase struct {
	log                 *log.Helper
	conf                *conf.Data
	JWTUsecase          *JWTUsecase
	Awsclaude3Usecase   *Awsclaude3Usecase
	AiResultUsecase     *AiResultUsecase
	AiUsecase           *AiUsecase
	ConditionbuzUsecase *ConditionbuzUsecase
	TUsecase            *TUsecase
	AiTaskUsecase       *AiTaskUsecase
	DataComboUsecase    *DataComboUsecase
	StatementUsecase    *StatementUsecase
	RelasLogUsecase     *RelasLogUsecase
	ConditionUsecase    *ConditionUsecase
	CommonUsecase       *CommonUsecase
	PsbuzUsecase        *PsbuzUsecase
	AiTaskJobUsecase    *AiTaskJobUsecase
	AiTaskbuzUsecase    *AiTaskbuzUsecase
}

func NewPsHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	Awsclaude3Usecase *Awsclaude3Usecase,
	AiResultUsecase *AiResultUsecase,
	AiUsecase *AiUsecase,
	ConditionbuzUsecase *ConditionbuzUsecase,
	TUsecase *TUsecase,
	AiTaskUsecase *AiTaskUsecase,
	DataComboUsecase *DataComboUsecase,
	StatementUsecase *StatementUsecase,
	RelasLogUsecase *RelasLogUsecase,
	ConditionUsecase *ConditionUsecase,
	CommonUsecase *CommonUsecase,
	PsbuzUsecase *PsbuzUsecase,
	AiTaskJobUsecase *AiTaskJobUsecase,
	AiTaskbuzUsecase *AiTaskbuzUsecase) *PsHttpUsecase {
	return &PsHttpUsecase{
		log:                 log.NewHelper(logger),
		conf:                conf,
		JWTUsecase:          JWTUsecase,
		Awsclaude3Usecase:   Awsclaude3Usecase,
		AiResultUsecase:     AiResultUsecase,
		AiUsecase:           AiUsecase,
		ConditionbuzUsecase: ConditionbuzUsecase,
		TUsecase:            TUsecase,
		AiTaskUsecase:       AiTaskUsecase,
		DataComboUsecase:    DataComboUsecase,
		StatementUsecase:    StatementUsecase,
		RelasLogUsecase:     RelasLogUsecase,
		ConditionUsecase:    ConditionUsecase,
		CommonUsecase:       CommonUsecase,
		PsbuzUsecase:        PsbuzUsecase,
		AiTaskJobUsecase:    AiTaskJobUsecase,
		AiTaskbuzUsecase:    AiTaskbuzUsecase,
	}
}

func (c *PsHttpUsecase) HandleUpdateStatement(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	caseId := body.GetInt("case_id")
	data, err := c.BizHandleUpdateStatement(caseId)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *PsHttpUsecase) BizHandleUpdateStatement(caseId int32) (lib.TypeMap, error) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("Parameter error")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return nil, err
	}
	if tClient == nil {
		return nil, errors.New("tClient is nil")
	}

	err = c.PsbuzUsecase.HandleUpdateStatement(*tClient, *tCase)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *PsHttpUsecase) GenerateDocument(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	caseId := body.GetInt("case_id")
	data, err := c.BizGenerateDocument(caseId)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *PsHttpUsecase) BizGenerateDocument(caseId int32) (lib.TypeMap, error) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("Parameter error")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return nil, err
	}
	if tClient == nil {
		return nil, errors.New("tClient is nil")
	}
	//summary, err := c.AiTaskbuzUsecase.HandleVeteranSummary(context.TODO(), tCase)
	//if err != nil {
	//	return nil, err
	//}
	err = c.StatementUsecase.GenerateDocument(*tCase, *tClient)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
