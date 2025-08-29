package biz

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

type BoxcontractUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	MapUsecase       *MapUsecase
	TUsecase         *TUsecase
	BoxUsecase       *BoxUsecase
	DataComboUsecase *DataComboUsecase
}

func NewBoxcontractUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	TUsecase *TUsecase,
	BoxUsecase *BoxUsecase,
	DataComboUsecase *DataComboUsecase) *BoxcontractUsecase {
	uc := &BoxcontractUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		MapUsecase:       MapUsecase,
		TUsecase:         TUsecase,
		BoxUsecase:       BoxUsecase,
		DataComboUsecase: DataComboUsecase,
	}
	return uc
}

// ContractFolderId 获取客户合同存储的文件夹id
func (c *BoxcontractUsecase) ContractFolderId(clientId int32) (contractFolderId string, err error) {
	key := fmt.Sprintf("%s%d", Map_ClientContractBoxFolderId, clientId)
	contractFolderId, err = c.MapUsecase.GetForString(key)
	if err != nil {
		return "", err
	}
	if contractFolderId != "" {
		return
	}
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientId)
	if err != nil {
		return "", err
	}
	if tClientCase == nil {
		return "", errors.New("ContractFolderId tClientCase is nil.")
	}

	_, tContactFields, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		return "", err
	}
	if tContactFields == nil {
		return "", errors.New("ContractFolderId tContactFields is nil.")
	}

	folderName := ClientContractFolderNameForBox(tContactFields.TextValueByNameBasic("first_name"),
		tContactFields.TextValueByNameBasic("last_name"), clientId)

	contractFolderId, err = c.BoxUsecase.CreateFolder(folderName, c.conf.Box.ClientContractsId)
	if err != nil {
		return "", err
	}
	if len(contractFolderId) == 0 {
		return "", errors.New("contractFolderId is empty.")
	}
	err = c.MapUsecase.Set(key, contractFolderId)
	// 合同保存
	if err != nil {
		return "", err
	}
	return contractFolderId, nil
}

// AmSignedAgreementFolderId 文件夹id
func (c *BoxcontractUsecase) AmSignedAgreementFolderId(caseId int32) (boxFolderId string, err error) {
	key := MapKeyClientCaseAmSignedAgreementBoxFolderId(caseId)
	boxFolderId, err = c.MapUsecase.GetForString(key)
	if err != nil {
		return "", err
	}
	if boxFolderId != "" {
		return
	}
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return "", err
	}
	if tCase == nil {
		return "", errors.New("tCase is nil.")
	}
	folderName := "AM - Signed Agreement"

	contractFolderId, err := c.ContractFolderId(caseId)
	if err != nil {
		return "", err
	}
	boxFolderId, err = c.BoxUsecase.CreateFolder(folderName, contractFolderId)
	if err != nil {
		return "", err
	}
	if len(boxFolderId) == 0 {
		return "", errors.New("boxFolderId is empty.")
	}
	err = c.MapUsecase.Set(key, boxFolderId)
	// 合同保存
	if err != nil {
		return "", err
	}
	return boxFolderId, nil
}

// AmSignedVA2122aBoxFolderId 文件夹id
func (c *BoxcontractUsecase) AmSignedVA2122aBoxFolderId(caseId int32) (boxFolderId string, err error) {
	key := MapKeyClientCaseAmSignedVA2122aBoxFolderId(caseId)
	boxFolderId, err = c.MapUsecase.GetForString(key)
	if err != nil {
		return "", err
	}
	if boxFolderId != "" {
		return
	}
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return "", err
	}
	if tCase == nil {
		return "", errors.New("tCase is nil.")
	}
	folderName := "AM - Signed VA 21-22a"

	contractFolderId, err := c.ContractFolderId(caseId)
	if err != nil {
		return "", err
	}

	boxFolderId, err = c.BoxUsecase.CreateFolder(folderName, contractFolderId)
	if err != nil {
		return "", err
	}
	if len(boxFolderId) == 0 {
		return "", errors.New("boxFolderId is empty.")
	}
	err = c.MapUsecase.Set(key, boxFolderId)
	// 合同保存
	if err != nil {
		return "", err
	}
	return boxFolderId, nil
}
