package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"io"
	"regexp"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	. "vbc/lib/builder"
)

func DocEmailExtractHealthIssues(text string) (result []string) {
	rgx := regexp.MustCompile(`(?m)(\d+)\.\s([^:]+):\s(.*?)\n{0,1}$`)
	matches := rgx.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			//fmt.Printf("%s. %s: %s\n", match[1], match[2], match[3])
			r := fmt.Sprintf("%s: %s", match[2], match[3])
			result = append(result, r)
		}
	}
	return
}

type DocEmailUsecase struct {
	log                    *log.Helper
	conf                   *conf.Data
	CommonUsecase          *CommonUsecase
	AiTaskUsecase          *AiTaskUsecase
	AiResultUsecase        *AiResultUsecase
	DataComboUsecase       *DataComboUsecase
	WordUsecase            *WordUsecase
	BoxbuzUsecase          *BoxbuzUsecase
	BoxUsecase             *BoxUsecase
	MapUsecase             *MapUsecase
	TUsecase               *TUsecase
	PersonalWebformUsecase *PersonalWebformUsecase
}

func NewDocEmailUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	AiTaskUsecase *AiTaskUsecase,
	AiResultUsecase *AiResultUsecase,
	DataComboUsecase *DataComboUsecase,
	WordUsecase *WordUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	BoxUsecase *BoxUsecase,
	MapUsecase *MapUsecase,
	TUsecase *TUsecase,
	PersonalWebformUsecase *PersonalWebformUsecase,
) *DocEmailUsecase {
	uc := &DocEmailUsecase{
		log:                    log.NewHelper(logger),
		CommonUsecase:          CommonUsecase,
		conf:                   conf,
		AiTaskUsecase:          AiTaskUsecase,
		AiResultUsecase:        AiResultUsecase,
		DataComboUsecase:       DataComboUsecase,
		WordUsecase:            WordUsecase,
		BoxbuzUsecase:          BoxbuzUsecase,
		BoxUsecase:             BoxUsecase,
		MapUsecase:             MapUsecase,
		TUsecase:               TUsecase,
		PersonalWebformUsecase: PersonalWebformUsecase,
	}

	return uc
}

func (c *DocEmailUsecase) DocEmailResult(ParseResult string, tClient TData) (wordLineList WordLineList, err error) {

	extractHealthIssuesResult := DocEmailExtractHealthIssues(ParseResult)
	if len(extractHealthIssuesResult) == 0 {
		return nil, errors.New("extractHealthIssuesResult is wrong")
	}

	wordLineList = append(wordLineList, WordDocEmailTop...)
	for _, v := range extractHealthIssuesResult {
		wordLineList = append(wordLineList, WordLine{
			Type:  WordLine_Type_List,
			Value: v,
		})
	}
	wordLineList = append(wordLineList, WordDocEmailBottom...)
	wordLineList = append(wordLineList, WordLine{
		Type:  WordLine_Type_Normal,
		Value: tClient.CustomFields.TextValueByNameBasic(FieldName_full_name),
	})
	return wordLineList, nil
}

func (c *DocEmailUsecase) DocEmailResultTextByCase(tCase TData, tClient TData) (string, error) {

	aiResultEntity, err := c.DocEmailResultAiResult(tCase.Id())
	if err != nil {
		return "", err
	}
	if aiResultEntity == nil {
		return "", nil
	}
	return c.DocEmailResultText(aiResultEntity.ParseResult, tClient)
}

func (c *DocEmailUsecase) DocEmailResultText(ParseResult string, tClient TData) (string, error) {
	wordLineList, err := c.DocEmailResult(ParseResult, tClient)
	if err != nil {
		return "", err
	}
	return wordLineList.ToString(), nil
}

func (c *DocEmailUsecase) DocEmailResultWordByCase(tCase TData, tClient TData) (io.Reader, error) {

	aiResultEntity, err := c.DocEmailResultAiResult(tCase.Id())
	if err != nil {
		return nil, err
	}
	var ParseResult string
	if aiResultEntity == nil {
		return nil, nil
	} else {
		ParseResult = aiResultEntity.ParseResult
	}
	return c.DocEmailResultWord(ParseResult, tClient)
}
func (c *DocEmailUsecase) DocEmailResultWord(ParseResult string, tClient TData) (io.Reader, error) {
	wordLineList, _ := c.DocEmailResult(ParseResult, tClient)
	//if err != nil {
	//	return nil, err
	//}
	return c.WordUsecase.CreateDocEmailWord(wordLineList)
}

func (c *DocEmailUsecase) DocEmailResultAiResult(caseId int32) (aiResult *AiResultEntity, err error) {

	key := MapKeyPSCurrentDocEmailAiResultId(caseId)

	aiResultId, err := c.MapUsecase.GetForInt(key)
	if err != nil {
		return nil, err
	}
	if aiResultId == 0 {
		aiTask, err := c.AiTaskUsecase.GetByCondWithOrderBy(Eq{"case_id": caseId,
			"deleted_at":    0,
			"from_type":     AiTaskFromType_generate_doc_email,
			"handle_status": HandleStatus_done,
			"handle_result": HandleResult_ok,
		}, "id desc")
		if err != nil {
			return nil, err
		}
		if aiTask == nil {
			return nil, nil
		}
		if aiTask.CurrentResultId == 0 {
			return nil, nil
		}
		er := c.SetLatestDocEmailResult(caseId, aiTask.CurrentResultId)
		if er != nil {
			c.log.Error(er)
		}

		aiResultId = aiTask.CurrentResultId
	}
	aiResult, _ = c.AiResultUsecase.GetByCond(Eq{"id": aiResultId})
	return aiResult, nil
}

func (c *DocEmailUsecase) SetLatestDocEmailResult(caseId int32, aiResultId int32) error {
	key := MapKeyPSCurrentDocEmailAiResultId(caseId)
	c.MapUsecase.SetInt(key, int(aiResultId))

	go func() {

		//c.StatementUsecase.
		useWebForm, err := c.PersonalWebformUsecase.IsUseNewPersonalWebForm(caseId)
		if err != nil {
			c.log.Error(err)
			return
		}
		if !useWebForm {
			tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
			if err != nil {
				c.log.Error(err)
			} else {
				err = c.HandleDocEmailToBox(tCase)
				if err != nil {
					c.log.Error(err)
				}
			}
		}
	}()

	return nil
}

func (c *DocEmailUsecase) HandleDocEmailToBox(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}

	tClient, _, _ := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if tClient == nil {
		return errors.New("tClient is nil")
	}

	//aiTask, err := c.AiTaskUsecase.GetByCondWithOrderBy(Eq{"case_id": tCase.Id(),
	//	"deleted_at":    0,
	//	"from_type":     AiTaskFromType_generate_doc_email,
	//	"handle_status": HandleStatus_done,
	//	"handle_result": HandleResult_ok,
	//}, "id desc")
	//if err != nil {
	//	return err
	//}
	//if aiTask == nil {
	//	return errors.New("aiTask is nil")
	//}
	//if aiTask.CurrentResultId == 0 {
	//	return errors.New("aiTask.CurrentResultId is 0")
	//}

	//aiResult, _ := c.AiResultUsecase.GetByCond(Eq{"id": aiTask.CurrentResultId})

	aiResult, _ := c.DocEmailResultAiResult(tCase.Id())
	if aiResult == nil {
		return errors.New("aiResult is nil")
	}

	wordReader, err := c.DocEmailResultWord(aiResult.ParseResult, *tClient)
	//wordLineList, err := c.DocEmailResult(aiResult.ParseResult, *tClient)
	//if err != nil {
	//	return err
	//}
	//
	//wordReader, err := c.WordUsecase.CreateDocEmailWord(wordLineList)
	if err != nil {
		return err
	}
	dCPersonalStatementsFolderId, aiDocEmailFileName, boxFileId, err := c.DocEmailBoxFileId(tClient, tCase)
	if err != nil {
		return err
	}
	if boxFileId == "" {
		boxFileId, err = c.BoxUsecase.UploadFile(dCPersonalStatementsFolderId, wordReader, aiDocEmailFileName)
		if err != nil {
			return err
		}
	} else {
		_, err = c.BoxUsecase.UploadFileVersion(boxFileId, wordReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DocEmailUsecase) DocEmailBoxFileId(tClient *TData, tCase *TData) (dCPersonalStatementsFolderId string, aiDocEmailFileName string, boxFileId string, err error) {

	if tClient == nil {
		return "", "", "", errors.New("tClient is nil")
	}
	if tCase == nil {
		return "", "", "", errors.New("tCase is nil")
	}

	aiDocEmailFileName = GenDocEmailFileNameAuto(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name), tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

	dCPersonalStatementsFolderId, err = c.BoxbuzUsecase.DCPersonalStatementsFolderId(tCase)
	if err != nil {
		return "", "", "", err
	}
	if dCPersonalStatementsFolderId == "" {
		return "", "", "", errors.New("dCPersonalStatementsFolderId is empty")
	}
	resItems, err := c.BoxUsecase.ListItemsInFolderFormat(dCPersonalStatementsFolderId)
	if err != nil {
		return "", "", "", err
	}
	for _, v := range resItems {
		resId := v.GetString("id")
		resType := v.GetString("type")
		resName := v.GetString("name")
		if resType == string(config_box.BoxResType_file) {
			if resName == aiDocEmailFileName {
				boxFileId = resId
				break
			}
		}
	}
	return
}
