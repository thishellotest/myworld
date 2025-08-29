package biz

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	"vbc/lib"
)

type RollpoingJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[RollpoingEntity]
	BaseHandleCustom[RollpoingEntity]
	BoxUsecase            *BoxUsecase
	TaskCreateUsecase     *TaskCreateUsecase
	ClientEnvelopeUsecase *ClientEnvelopeUsecase
	BehaviorUsecase       *BehaviorUsecase
}

func NewRollpoingJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	BoxUsecase *BoxUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
	BehaviorUsecase *BehaviorUsecase) *RollpoingJobUsecase {
	uc := &RollpoingJobUsecase{
		log:                   log.NewHelper(logger),
		CommonUsecase:         CommonUsecase,
		conf:                  conf,
		BoxUsecase:            BoxUsecase,
		TaskCreateUsecase:     TaskCreateUsecase,
		ClientEnvelopeUsecase: ClientEnvelopeUsecase,
		BehaviorUsecase:       BehaviorUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandleCustom.DB = CommonUsecase.DB()
	uc.BaseHandleCustom.Log = log.NewHelper(logger)

	return uc
}

func (c *RollpoingJobUsecase) WaitingTasks(ctx context.Context) (*sql.Rows, error) {

	return c.CommonUsecase.DB().
		Table(RollpoingEntity{}.TableName()).
		Where("handle_status=? and next_at<=? and deleted_at=0",
			HandleStatus_waiting, time.Now().Unix()).Rows()

}

func (c *RollpoingJobUsecase) Handle(ctx context.Context, task *RollpoingEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	isDone, err := c.HandleExec(ctx, task)
	fmt.Println("RollpoingJobUsecase:Handle:", isDone, err)

	task.UpdatedAt = time.Now().Unix() // 解决修改无更新有一次sql的问题
	if err != nil {
		task.HandleStatus = HandleStatus_done
		task.HandleResult = HandleResult_failure
		task.HandleResultDetail = err.Error()
		c.log.Error(err.Error(), "Rollpoing task id: ", task.ID)
		return c.CommonUsecase.DB().Save(task).Error
	} else if isDone {
		task.HandleStatus = HandleStatus_done
		task.HandleResult = HandleResult_ok
		return c.CommonUsecase.DB().Save(task).Error
	} else if time.Now().Unix() > (task.CreatedAt + 365*24*3600) {
		task.HandleStatus = HandleStatus_done
		task.HandleResult = HandleResult_failure
		task.HandleResultDetail = "超时不在处理"
		return c.CommonUsecase.DB().Save(task).Error
	} else {
		if task.Vendor == Rollpoing_Vendor_boxsign {
			// 不要太快，因为api可能收费
			task.NextAt = time.Now().Unix() + 6*3600 // 每隔6小时查询一次：api limit 10万 per/month
		} else {
			task.NextAt = time.Now().Unix() + 3600 // 20秒后查询
		}

		return c.CommonUsecase.DB().Save(task).Error
	}
}

func (c *RollpoingJobUsecase) HandleExec(ctx context.Context, task *RollpoingEntity) (isDone bool, err error) {

	if task == nil {
		return false, errors.New("task is nil")
	}
	if task.Vendor == Rollpoing_Vendor_boxsign {

		clientEnvelope, err := c.ClientEnvelopeUsecase.GetByEnvelopeId(EsignVendor_box, task.VendorUniqId)
		if err != nil {
			return false, err
		}
		if clientEnvelope == nil {
			return false, errors.New("clientEnvelope is nil")
		}
		res, err := c.BoxUsecase.GetSignRequest(task.VendorUniqId)
		if res != nil {
			defer res.Body.Close()
		}
		if err != nil {
			return false, err
		}
		if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
			bs, err := io.ReadAll(res.Body)
			if err != nil {
				return false, err
			}
			return c.HandleExecBox(string(bs), clientEnvelope.ClientId, clientEnvelope)
		} else {
			if res.StatusCode == config_box.HttpCode_404 {
				clientEnvelope.IsSigned = ClientEnvelope_IsSigned_Cancelled
				clientEnvelope.UpdatedAt = time.Now().Unix()
				clientEnvelope.SignStatus = "StatusCode:" + InterfaceToString(res.StatusCode)
				er := c.CommonUsecase.DB().Save(&clientEnvelope).Error
				if er != nil {
					c.log.Error(er)
				}
			}
			return false, errors.New("StatusCode:" + InterfaceToString(res.StatusCode))
		}

	} else {
		return false, errors.New("Vendor类型不支持")
	}
}

func ClientEnvelopeTypeToBehaviorType(typ string) (behaviorType string, err error) {
	if typ == Type_FeeContract {
		return BehaviorType_complete_fee_schedule_contract, nil
	} else if typ == Type_MedicalTeamForms {
		return BehaviorType_complete_medical_team_forms_contract, nil
	} else if typ == Type_AmContract {
		return BehaviorType_complete_am_contract, nil
	}
	/*
		else if typ == Type_ReleaseOfInformation {
			return BehaviorType_complete_release_of_information_contract, nil
		} else if typ == Type_PatientPaymentForm {
			return BehaviorType_complete_patient_payment_form_contract, nil
		}*/
	return "", errors.New("ClientEnvelope Type is wrong.")
}

func (c *RollpoingJobUsecase) HandleExecBox(body string, clientId int32, clientEnvelopeEntity *ClientEnvelopeEntity) (isDone bool, err error) {

	resMap := lib.ToTypeMapByString(body)

	if clientEnvelopeEntity == nil {
		return false, errors.New("clientEnvelopeEntity is nil")
	}

	status := resMap.GetString("status")
	// 合同完成签属
	if config_box.IsBoxSignStatusFinishSigned(status) {
		BehaviorType, err := ClientEnvelopeTypeToBehaviorType(clientEnvelopeEntity.Type)
		if err != nil {
			return false, err
		}

		clientEnvelopeEntity.IsSigned = ClientEnvelope_IsSigned_Yes
		clientEnvelopeEntity.UpdatedAt = time.Now().Unix()
		clientEnvelopeEntity.SignStatus = status
		er := c.CommonUsecase.DB().Save(&clientEnvelopeEntity).Error
		if er != nil {
			c.log.Error(er)
		}

		// 处理行为
		err = c.BehaviorUsecase.Upsert(clientId, BehaviorType, time.Now(), "")
		if err != nil {
			return false, err
		}
		return true, nil
	} else if config_box.IsBoxSignStatusFinalizing(status) {

		clientEnvelopeEntity.IsSigned = ClientEnvelope_IsSigned_Cancelled
		clientEnvelopeEntity.UpdatedAt = time.Now().Unix()
		clientEnvelopeEntity.SignStatus = status
		er := c.CommonUsecase.DB().Save(&clientEnvelopeEntity).Error
		if er != nil {
			c.log.Error(er)
		}
		return true, errors.New("StatusFinalizing: " + status)
	} else {
		return false, nil
	}
}
