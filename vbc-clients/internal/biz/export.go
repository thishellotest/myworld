package biz

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"vbc/internal/conf"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

type ExportUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	DocEmailUsecase  *DocEmailUsecase
	TUsecase         *TUsecase
	DataComboUsecase *DataComboUsecase
	WordUsecase      *WordUsecase
	GopdfUsecase     *GopdfUsecase
	StatementUsecase *StatementUsecase
}

func NewExportUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	DocEmailUsecase *DocEmailUsecase,
	TUsecase *TUsecase,
	DataComboUsecase *DataComboUsecase,
	WordUsecase *WordUsecase,
	GopdfUsecase *GopdfUsecase,
	StatementUsecase *StatementUsecase,
) *ExportUsecase {
	uc := &ExportUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		DocEmailUsecase:  DocEmailUsecase,
		TUsecase:         TUsecase,
		DataComboUsecase: DataComboUsecase,
		WordUsecase:      WordUsecase,
		GopdfUsecase:     GopdfUsecase,
		StatementUsecase: StatementUsecase,
	}

	return uc
}

const (
	Export_type_all_statements = "all_statements"
	Export_type_statement      = "statement"
	Export_type_docemail       = "docemail"
)

func (c *ExportUsecase) BizHttpAllStatements(caseGid string) ([]byte, string, error) {

	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		return nil, "", errors.New("Incorrect parameters")
	}
	tClient, _, _ := c.DataComboUsecase.ClientWithCase(*tCase)
	if tClient == nil {
		return nil, "", errors.New("Incorrect parameters")
	}

	statementDetail, err := c.StatementUsecase.GetListStatementDetail(false, *tClient, *tCase, 0)
	if err != nil {
		return nil, "", err
	}
	ioReader, err := c.GopdfUsecase.CreatePersonalStatementsPDFForAiV1(tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name), statementDetail, 0)
	if err != nil {
		return nil, "", err
	}
	if ioReader == nil {
		return nil, "", errors.New("The content might still be under processing and has not been found yet")
	}
	r, err := io.ReadAll(ioReader)
	if err != nil {
		return nil, "", err
	}
	if len(r) == 0 {
		return nil, "", errors.New("The content might still be under processing and has not been found yet")
	}
	name := GenPersonalStatementsFileNamePdf(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
		tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())
	return r, name, nil
}

func (c *ExportUsecase) BizHttpStatement(caseGid string, statementConditionId int32) ([]byte, string, error) {

	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		return nil, "", errors.New("Incorrect parameters")
	}
	tClient, _, _ := c.DataComboUsecase.ClientWithCase(*tCase)
	if tClient == nil {
		return nil, "", errors.New("Incorrect parameters")
	}

	statementDetail, err := c.StatementUsecase.GetListStatementDetail(false, *tClient, *tCase, 0)
	if err != nil {
		return nil, "", err
	}
	ioReader, err := c.GopdfUsecase.CreatePersonalStatementsPDFForAiV1(tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name), statementDetail, statementConditionId)
	if err != nil {
		return nil, "", err
	}
	if ioReader == nil {
		return nil, "", errors.New("The content might still be under processing and has not been found yet")
	}
	r, err := io.ReadAll(ioReader)
	if err != nil {
		return nil, "", err
	}
	if len(r) == 0 {
		return nil, "", errors.New("The content might still be under processing and has not been found yet")
	}
	name := GenPersonalStatementsFileNamePdf(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
		tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())
	return r, fmt.Sprintf("%d-%s", statementConditionId, name), nil
}

func (c *ExportUsecase) BizHttpDocemail(caseGid string) ([]byte, string, error) {

	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		return nil, "", errors.New("Incorrect parameters")
	}
	tClient, _, _ := c.DataComboUsecase.ClientWithCase(*tCase)
	if tClient == nil {
		return nil, "", errors.New("Incorrect parameters")
	}
	ioReader, err := c.DocEmailUsecase.DocEmailResultWordByCase(*tCase, *tClient)
	if err != nil {
		return nil, "", err
	}
	if ioReader == nil {
		return nil, "", errors.New("The content might still be under processing and has not been found yet")
	}
	r, err := io.ReadAll(ioReader)
	if err != nil {
		return nil, "", err
	}
	if len(r) == 0 {
		return nil, "", errors.New("The content might still be under processing and has not been found yet")
	}
	name := GenDocPDFEmailFileName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name),
		tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())
	return r, name, nil
}

func (c *ExportUsecase) Http(ctx *gin.Context) {
	caseGid := ctx.Query("gid")
	bizType := ctx.Query("type")
	statementConditionIdStr := ctx.Query("statement_condition_id")
	statementConditionId, _ := strconv.ParseInt(statementConditionIdStr, 10, 32)
	var err error
	var docBytes []byte
	var docName string

	if Export_type_docemail == bizType {
		docBytes, docName, err = c.BizHttpDocemail(caseGid)
	} else if Export_type_all_statements == bizType {
		docBytes, docName, err = c.BizHttpAllStatements(caseGid)
	} else if Export_type_statement == bizType {
		docBytes, docName, err = c.BizHttpStatement(caseGid, int32(statementConditionId))
	} else {
		err = errors.New("The download URL is incorrect")
	}

	//docxBytes, err := generateDocxBytes()
	if err != nil {
		docName = "export.docx"
		wordLineList := WordLineList{
			{
				Type:  WordLine_Type_Normal,
				Value: err.Error(),
			},
		}
		errIOReader, _ := c.WordUsecase.CreateDocEmailWord(wordLineList)
		if errIOReader != nil {
			docBytes, _ = io.ReadAll(errIOReader)
		}
	}

	// Set appropriate headers based on file type
	ctx.Header("Content-Disposition", `attachment; filename="`+docName+`"`)
	if strings.HasSuffix(docName, ".pdf") {
		// For PDF files
		ctx.Header("Content-Type", "application/pdf")
		ctx.Data(http.StatusOK, "application/pdf", docBytes)
	} else {
		// For Word files (other types)
		ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", docBytes)
	}

	//lib.DPrintln(caseGid, bizType)
}
