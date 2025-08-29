package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ContractbuzUsecase struct {
	log                    *log.Helper
	conf                   *conf.Data
	CommonUsecase          *CommonUsecase
	TUsecase               *TUsecase
	MapUsecase             *MapUsecase
	ClientCaseUsecase      *ClientCaseUsecase
	ZohobuzUsecase         *ZohobuzUsecase
	DataEntryUsecase       *DataEntryUsecase
	RevisionHistoryUsecase *RevisionHistoryUsecase
	MgmtPermissionUsecase  *MgmtPermissionUsecase
}

func NewContractbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	MapUsecase *MapUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	ZohobuzUsecase *ZohobuzUsecase,
	DataEntryUsecase *DataEntryUsecase,
	RevisionHistoryUsecase *RevisionHistoryUsecase,
	MgmtPermissionUsecase *MgmtPermissionUsecase,
) *ContractbuzUsecase {
	uc := &ContractbuzUsecase{
		log:                    log.NewHelper(logger),
		CommonUsecase:          CommonUsecase,
		conf:                   conf,
		TUsecase:               TUsecase,
		MapUsecase:             MapUsecase,
		ClientCaseUsecase:      ClientCaseUsecase,
		ZohobuzUsecase:         ZohobuzUsecase,
		DataEntryUsecase:       DataEntryUsecase,
		RevisionHistoryUsecase: RevisionHistoryUsecase,
		MgmtPermissionUsecase:  MgmtPermissionUsecase,
	}

	return uc
}

func (c *ContractbuzUsecase) VerifyWhetherCanModifyContract(tCase TData) (isOK bool, err error) {

	clientGid := tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid)
	otherCase, err := c.TUsecase.Data(Kind_client_cases, And(
		Eq{FieldName_biz_deleted_at: 0, DataEntry_deleted_at: 0, FieldName_client_gid: clientGid},
		Neq{DataEntry_gid: tCase.Gid()}))
	if err != nil {
		return false, err
	}
	if otherCase != nil {
		return false, errors.New("There are multiple cases and the contract cannot be modified")
	}
	// 37. Completed
	// 38. Terminated
	// 39. Dormant
	stage := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stage == config_vbc.Stages_Completed ||
		stage == config_vbc.Stages_Terminated ||
		stage == config_vbc.Stages_Dormant {
		return false, errors.New("During stage Completed, Terminated and Dormant, the contract is not allowed to be modified")
	}

	key := MapKeyClientCaseContractBasicData(tCase.Id())
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return false, err
	}
	if val == "" {
		return false, errors.New("There is no contract yet.")
	}
	return true, nil
}

type ContractbuzBizGetVo struct {
	CanModifyContract bool                          `json:"can_modify_contract"`
	ContractBasicData ClientCaseContractBasicDataVo `json:"contract_basic_data"`
	Gid               string                        `json:"gid"`
	DealName          string                        `json:"deal_name"`
}

func (c *ContractbuzUsecase) BizGet(caseId int32, userFacade UserFacade) (lib.TypeMap, error) {

	hasPermission, err := c.MgmtPermissionUsecase.Verify(userFacade, MgmtPermission_ReviseContract)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New(Error_UnauthorizedOperation)
	}

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("The case does not exist")
	}
	isOk, err := c.VerifyWhetherCanModifyContract(*tCase)
	if err != nil {
		return nil, err
	}
	var contractbuzBizGetVo ContractbuzBizGetVo
	if isOk {
		contractbuzBizGetVo.CanModifyContract = isOk
		clientCaseContractBasicDataVo, _ := c.ClientCaseUsecase.ClientCaseContractBasicDataVoById(tCase.Id())
		if clientCaseContractBasicDataVo == nil {
			return nil, errors.New("The contract is incorrect and cannot be operated")
		}
		contractbuzBizGetVo.ContractBasicData = *clientCaseContractBasicDataVo
		contractbuzBizGetVo.Gid = tCase.Gid()
		contractbuzBizGetVo.DealName = tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	}
	data := make(lib.TypeMap)
	data.Set("data", contractbuzBizGetVo)
	return data, nil
}

func (c *ContractbuzUsecase) BizSave(contractHttpSaveVo ContractHttpSaveVo, userFacade UserFacade) (lib.TypeMap, error) {

	hasPermission, err := c.MgmtPermissionUsecase.Verify(userFacade, MgmtPermission_ReviseContract)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New(Error_UnauthorizedOperation)
	}

	tCase, err := c.TUsecase.DataById(Kind_client_cases, contractHttpSaveVo.CaseId)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("The case does not exist")
	}
	isOk, err := c.VerifyWhetherCanModifyContract(*tCase)
	if err != nil {
		return nil, err
	}
	if isOk {
		clientCaseContractBasicDataVo, _ := c.ClientCaseUsecase.ClientCaseContractBasicDataVoById(tCase.Id())
		if clientCaseContractBasicDataVo == nil {
			return nil, errors.New("The contract is incorrect and cannot be operated")
		}

		oldBasicVo := *clientCaseContractBasicDataVo

		dataEntry := make(TypeDataEntry)
		if contractHttpSaveVo.ActiveDuty {
			if contractHttpSaveVo.ActiveDuty == clientCaseContractBasicDataVo.ActiveDuty {
				return nil, errors.New("No modifications at all.")
			}
			clientCaseContractBasicDataVo.ActiveDuty = true

			key := MapKeyClientCaseContractBasicData(tCase.Id())
			err = c.MapUsecase.Set(key, InterfaceToString(clientCaseContractBasicDataVo))
			if err != nil {
				return nil, err
			}
			dataEntry[DataEntry_gid] = tCase.Gid()
			dataEntry[FieldName_active_duty] = ActiveDuty_Yes
		} else {
			if !VerifyCurrentRating(contractHttpSaveVo.Rating) {
				return nil, errors.New("The rating setting is incorrect")
			}
			if contractHttpSaveVo.ActiveDuty == clientCaseContractBasicDataVo.ActiveDuty &&
				contractHttpSaveVo.Rating == clientCaseContractBasicDataVo.EffectiveCurrentRating {
				return nil, errors.New("No modifications at all.")
			}
			clientCaseContractBasicDataVo.ActiveDuty = false
			clientCaseContractBasicDataVo.EffectiveCurrentRating = contractHttpSaveVo.Rating
			dataEntry[DataEntry_gid] = tCase.Gid()
			dataEntry[FieldName_active_duty] = ActiveDuty_No
			dataEntry[FieldName_effective_current_rating] = clientCaseContractBasicDataVo.EffectiveCurrentRating
		}

		key := MapKeyClientCaseContractBasicData(tCase.Id())
		err = c.MapUsecase.Set(key, InterfaceToString(clientCaseContractBasicDataVo))
		if err != nil {
			return nil, err
		}

		_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)
		if err != nil {
			return nil, err
		}
		_, err = c.RevisionHistoryUsecase.Add(RevisionHistory_BizType_contract,
			tCase.Gid(),
			InterfaceToString(oldBasicVo),
			InterfaceToString(clientCaseContractBasicDataVo), userFacade.Gid())
		if err != nil {
			return nil, err
		}

		err = c.ZohobuzUsecase.HandleAmount(tCase.Id())
		if err != nil {
			return nil, err
		}
		err = c.ZohobuzUsecase.HandleClientCaseName(tCase.Id())
		if err != nil {
			return nil, err
		}

	}
	data := make(lib.TypeMap)
	return data, nil
}
