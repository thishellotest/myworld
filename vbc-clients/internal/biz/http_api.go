package biz

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"vbc/lib"
	"vbc/lib/builder"
)

type HttpApiUsecase struct {
	MailUsecase      *MailUsecase
	DataEntryUsecase *DataEntryUsecase
	TUsecase         *TUsecase
}

func NewHttpApiUsecase(MailUsecase *MailUsecase, DataEntryUsecase *DataEntryUsecase, TUsecase *TUsecase) *HttpApiUsecase {
	return &HttpApiUsecase{
		MailUsecase:      MailUsecase,
		DataEntryUsecase: DataEntryUsecase,
		TUsecase:         TUsecase,
	}
}

func (c *HttpApiUsecase) TestMailTpl(ctx *gin.Context) {

	reply := CreateReply()
	var request lib.TypeMap
	rawData, _ := ctx.GetRawData()
	json.Unmarshal(rawData, &request)

	tpl := request.GetString("tpl")
	id := request.GetInt("id")
	if len(tpl) == 0 || id <= 0 {
		reply.CommonStrError("Invalid parameter.")
		goto end
	} else {
		tplData, _ := c.TUsecase.Data(Kind_email_tpls, builder.Eq{"tpl": tpl})
		customerData, _ := c.TUsecase.Data(Kind_client_cases, builder.Eq{"id": id})
		err, _, _, _, _, _ := c.MailUsecase.SendEmailWithData(customerData, tplData, nil)
		if err != nil {
			reply.CommonError(err)
			goto end
		}
	}
	reply.Success()
end:
	ctx.JSON(200, reply)
}
