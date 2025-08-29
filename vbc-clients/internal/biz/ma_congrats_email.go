package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
)

type MaCongratsEmailUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	MapUsecase        *MapUsecase
	TUsecase          *TUsecase
	TaskCreateUsecase *TaskCreateUsecase
	DataComboUsecase  *DataComboUsecase
}

func NewMaCongratsEmailUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	TUsecase *TUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	DataComboUsecase *DataComboUsecase) *MaCongratsEmailUsecase {
	uc := &MaCongratsEmailUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		MapUsecase:        MapUsecase,
		TUsecase:          TUsecase,
		TaskCreateUsecase: TaskCreateUsecase,
		DataComboUsecase:  DataComboUsecase,
	}

	return uc
}

func (c *MaCongratsEmailUsecase) HandleInputTask(clientCaseId int32) error {

	key := fmt.Sprintf("%s%d", Map_MaCongratsEmail, clientCaseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val != "" { // 已经发送了
		return nil
	}
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return err
	}

	if tClientCase == nil {
		return errors.New("tClient is nil.")
	}
	stages := tClientCase.CustomFields.TextValueByNameBasic("stages")
	if stages != config_vbc.Stages_AwaitingPayment && stages != config_vbc.Stages_AmAwaitingPayment {
		return nil
	}

	_, tContactFields, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		return err
	}

	if tContactFields == nil {
		return errors.New("tContactFields is nil.")
	}

	email := tContactFields.TextValueByNameBasic(FieldName_email)
	if email == "" {
		return errors.New("email is empty.")
	}
	firstName := tContactFields.TextValueByNameBasic(FieldName_first_name)
	if firstName == "" {
		return errors.New("firstName is empty.")
	}
	effectiveCurrentRating := tClientCase.CustomFields.NumberValueByNameBasic(FieldName_effective_current_rating)
	newRating := tClientCase.CustomFields.NumberValueByNameBasic(FieldName_new_rating)
	if effectiveCurrentRating < 0 {
		return errors.New("effectiveCurrentRating is wrong.")
	}
	if newRating <= 0 {
		return nil
	}

	if tClientCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource) == ContractSource_AM {
		err = c.TaskCreateUsecase.CreateTaskMail(clientCaseId, MailGenre_AmCongratulationsNewRating, 0, nil, 0, "", "")
		if err != nil {
			return err
		}
	} else {
		err = c.TaskCreateUsecase.CreateTaskMail(clientCaseId, MailGenre_CongratulationsNewRating, 0, nil, 0, "", "")
		if err != nil {
			return err
		}
	}
	val = "1"
	return c.MapUsecase.Set(key, val)
}
