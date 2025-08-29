package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type PsbuzUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	BoxbuzUsecase    *BoxbuzUsecase
	AiTaskUsecase    *AiTaskUsecase
	BoxUsecase       *BoxUsecase
	WordbuzUsecase   *WordbuzUsecase
	StatementUsecase *StatementUsecase
}

func NewPsbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	AiTaskUsecase *AiTaskUsecase,
	BoxUsecase *BoxUsecase,
	WordbuzUsecase *WordbuzUsecase,
	StatementUsecase *StatementUsecase,
) *PsbuzUsecase {
	uc := &PsbuzUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		BoxbuzUsecase:    BoxbuzUsecase,
		AiTaskUsecase:    AiTaskUsecase,
		BoxUsecase:       BoxUsecase,
		WordbuzUsecase:   WordbuzUsecase,
		StatementUsecase: StatementUsecase,
	}

	return uc
}

func (c *PsbuzUsecase) DoGetPersonalStatementsDocxByCaseForUpdateStatement(tClient TData, tCase TData) (personalStatementsVo PersonalStatementsVo, err error) {

	// 先判断文件名称是否存在
	dCPersonalStatementsFolderId, ClientPSSourceFileName, boxFileId, err := c.StatementUsecase.DocClientPSSourceBoxFileId(tClient, tCase)
	if err != nil {
		return PersonalStatementsVo{}, err
	}
	if boxFileId == "" {

		_, _, psBoxFileId, err := c.BoxbuzUsecase.RealtimeCPersonalStatementsDocxFileId(tClient, tCase)
		if err != nil {
			return PersonalStatementsVo{}, err
		}
		if psBoxFileId == "" {
			return PersonalStatementsVo{}, errors.New("PersonalStatementDocFileId is empty")
		}
		boxFileId, err = c.BoxUsecase.CopyFileNewFileNameReturnFileId(psBoxFileId, ClientPSSourceFileName, dCPersonalStatementsFolderId)
		if err != nil {
			return PersonalStatementsVo{}, err
		}
	}
	return c.WordbuzUsecase.GetPersonalStatementsDocx(boxFileId)
}

func (c *PsbuzUsecase) HandleUpdateStatement(tClient TData, tCase TData) error {

	personalStatementsVo, err := c.DoGetPersonalStatementsDocxByCaseForUpdateStatement(tClient, tCase)
	if err != nil {
		return err
	}

	for k, _ := range personalStatementsVo.Statements {

		a, err := c.AiTaskUsecase.GetByCond(Eq{"from_type": AiTaskFromType_update_statement,
			"case_id":       tCase.Id(),
			"deleted_at":    0,
			"serial_number": k,
		})
		if err != nil {
			return err
		}
		if a == nil {
			var personalStatementOneVo PersonalStatementOneVo
			personalStatementOneVo.BaseInfo = personalStatementsVo.BaseInfo
			personalStatementOneVo.Statement = personalStatementsVo.Statements[k]
			personalStatementOneVo.DealName = personalStatementsVo.DealName
			_, err = c.AiTaskUsecase.CreateUpdateStatementTask(tCase, personalStatementOneVo, k)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//
//func (c *PsbuzUsecase) HandleUpdateStatementOld(tClient TData, tCase TData) error {
//
//	a, err := c.AiTaskUsecase.GetByCond(Eq{"from_type": AiTaskFromType_update_statement, "case_id": tCase.Id(), "deleted_at": 0})
//	if err != nil {
//		return err
//	}
//	if a == nil {
//		personalStatementsVo, err := c.WordbuzUsecase.GetPersonalStatementsDocxByCase(tClient, tCase)
//		if err != nil {
//			return err
//		}
//		_, err = c.AiTaskUsecase.CreateUpdateStatementTask(tCase, personalStatementsVo)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
