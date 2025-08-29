package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"os"
	"vbc/internal/conf"
)

type WordbuzUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	BoxUsecase    *BoxUsecase
	BoxbuzUsecase *BoxbuzUsecase
}

func NewWordbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	BoxUsecase *BoxUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
) *WordbuzUsecase {
	uc := &WordbuzUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		BoxUsecase:    BoxUsecase,
		BoxbuzUsecase: BoxbuzUsecase,
	}

	return uc
}

func (c *WordbuzUsecase) GetPersonalStatementsDocxByCase(tClient TData, tCase TData) (personalStatementsVo PersonalStatementsVo, err error) {

	//_, boxFileId, err := c.BoxbuzUsecase.PersonalStatementDocFileBoxFileId(&tClient, &tCase)
	//if err != nil {
	//	return PersonalStatementsVo{}, err
	//}
	_, _, boxFileId, err := c.BoxbuzUsecase.RealtimeCPersonalStatementsDocxFileId(tClient, tCase)

	if boxFileId == "" {
		return PersonalStatementsVo{}, errors.New("PersonalStatementDocFileId is empty")
	}
	return c.GetPersonalStatementsDocx(boxFileId)
}

func (c *WordbuzUsecase) GetPersonalStatementsDocx(boxFileId string) (personalStatementsVo PersonalStatementsVo, err error) {

	file, path, err := c.BoxUsecase.DownloadToLocal(boxFileId, "docx")
	if err != nil {
		return PersonalStatementsVo{}, err
	}
	filenamePath := path + "/" + file
	defer os.Remove(filenamePath)
	text, err := ReadDocxText(filenamePath)
	if err != nil {
		return PersonalStatementsVo{}, err
	}
	personalStatementsVo, err = SplitPersonalStatementsString(text)
	return
}
