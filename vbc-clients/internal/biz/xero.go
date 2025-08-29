package biz

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
)

const (
	Xero_Reference = ""
)

type XeroUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	Oauth2TokenUsecase *Oauth2TokenUsecase
	MapUsecase         *MapUsecase
	LogUsecase         *LogUsecase
	FeeUsecase         *FeeUsecase
	DataComboUsecase   *DataComboUsecase
	ClientCaseUsecase  *ClientCaseUsecase
	DataEntryUsecase   *DataEntryUsecase
}

func NewXeroUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	Oauth2TokenUsecase *Oauth2TokenUsecase,
	MapUsecase *MapUsecase,
	LogUsecase *LogUsecase,
	FeeUsecase *FeeUsecase,
	DataComboUsecase *DataComboUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	DataEntryUsecase *DataEntryUsecase) *XeroUsecase {
	uc := &XeroUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		Oauth2TokenUsecase: Oauth2TokenUsecase,
		MapUsecase:         MapUsecase,
		LogUsecase:         LogUsecase,
		FeeUsecase:         FeeUsecase,
		DataComboUsecase:   DataComboUsecase,
		ClientCaseUsecase:  ClientCaseUsecase,
		DataEntryUsecase:   DataEntryUsecase,
	}

	return uc
}

func (c *XeroUsecase) Headers() (map[string]string, error) {

	token, err := c.Oauth2TokenUsecase.GetAccessToken(Oauth2_AppId_xero)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Authorization": "Bearer " + token, "Xero-Tenant-Id": c.conf.Xero.TenantId}, nil
}

// Accounts 获取帐号类型对应 AccountCode
func (c *XeroUsecase) Accounts() error {
	headers, err := c.Headers()
	if err != nil {
		return err
	}

	api := fmt.Sprintf("%s/api.xro/2.0/Accounts", c.conf.Xero.ApiUrl)
	a, _, err := lib.Request("GET", api, nil, headers)
	lib.DPrintln(err)
	lib.DPrintln(a)
	return nil
}

func (c *XeroUsecase) BizCreateOrUpdateContact(tContact *TData) (contactId string, err error) {
	if tContact == nil {
		return "", errors.New("tContact is nil")
	}

	key := fmt.Sprintf("%s%d", Map_XeroContactId, tContact.CustomFields.NumberValueByNameBasic("id"))
	contactId, err = c.MapUsecase.GetForString(key)
	if err != nil {
		return "", err
	}
	res, err := c.CreateOrUpdateContact(tContact, contactId)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("Xero response is nil")
	}
	data := lib.ToTypeMapByString(*res)
	responseContactId := data.GetString("ContactID")
	if responseContactId == "" {
		return "", errors.New("responseContactId is empty.")
	}
	er := c.MapUsecase.Set(key, responseContactId)
	if er != nil {
		return "", er
	}
	return responseContactId, nil
}

// CreateOrUpdateContact 目前验证的结果有可能通过姓名去重， 有重名风险
func (c *XeroUsecase) CreateOrUpdateContact(tContact *TData, ContactID string) (*string, error) {
	/*
		说明：
		1. 通过修改 ContactNumber， 或FirstName LastName Name 并没有产生新的id（）可能是手机或邮件一致的原因
		2. 取消手机号，修改 ContactNumber和email, 但是同名 结果：通过姓名一样修改了对应信息
		3. 全新名子和email，无手机号： 产生了新的contactID
		4. 已经发送的invoice信息，不能随着Contact Data修改而更新。
		5. 状态为Draft的invoice信息，也不会随着Contact Data修改而更新
		6. 通过传递ContactID：可以更新客户信息
	*/
	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}
	//if tClientCase == nil {
	//	return nil, errors.New("tClientCase is nil")
	//}
	if tContact == nil {
		return nil, errors.New("tContact is nil")
	}

	tContactFields := tContact.CustomFields
	/*
		 "Phones": [
		        {
		          "PhoneType": "MOBILE",
		          "PhoneNumber": "555-1212",
		          "PhoneAreaCode": "415"
		        }
		      ],


		"Addresses": [
		          {
		            "AddressType": "STREET",
		            "City": "",
		            "Region": "",
		            "PostalCode": "",
		            "Country": ""
		          },
		          {
		            "AddressType": "POBOX",
		            "AddressLine1": "919 Rigley St",
		            "City": "Chula Vista",
		            "Region": "California",
		            "PostalCode": "91911",
		            "Country": "United States"
		          }
		        ],
	*/
	var Addresses lib.TypeList
	address := make(lib.TypeMap)
	dbAddr := tContactFields.TextValueByNameBasic(FieldName_address)
	dbAddr = strings.TrimRight(dbAddr, ",")
	dbAddr = strings.TrimRight(dbAddr, ", ")
	address.Set("AddressType", "POBOX")
	address.Set("AddressLine1", dbAddr)
	address.Set("PostalCode", tContactFields.TextValueByNameBasic(FieldName_zip_code))
	address.Set("City", tContactFields.TextValueByNameBasic(FieldName_city))
	address.Set("Region", tContactFields.TextValueByNameBasic(FieldName_state))
	address.Set("Country", "USA")
	Addresses = append(Addresses, address)

	/*
		48d47f54-f562-4618-a386-cf9b8a5c3963:  gengling.liao@hotmail.com Ling Li
	*/
	ContactNumber := InterfaceToString(tContactFields.NumberValueByNameBasic("id"))

	api := fmt.Sprintf("%s/api.xro/2.0/Contacts", c.conf.Xero.ApiUrl)

	contact := make(lib.TypeMap)
	if ContactID != "" {
		contact.Set("ContactID", ContactID)
	}
	contact.Set("ContactNumber", ContactNumber)
	contact.Set("Name", GenFullName(tContactFields.TextValueByNameBasic("first_name"), tContactFields.TextValueByNameBasic("last_name")))
	contact.Set("FirstName", tContactFields.TextValueByNameBasic("first_name"))
	contact.Set("LastName", tContactFields.TextValueByNameBasic("last_name"))
	contact.Set("EmailAddress", tContactFields.TextValueByNameBasic("email"))
	contact.Set("Addresses", Addresses)
	paymentTerms := lib.ToTypeMapByString(`{
        "Bills": {
          "Day": 7,
          "Type": "OFCURRENTMONTH"
        },
        "Sales": {
          "Day": 7,
          "Type": "DAYSAFTERBILLMONTH"
        }
      }`)
	contact.Set("PaymentTerms", paymentTerms)
	payload := `{"Contacts": [` + contact.ToString() + `]}`

	/*
	   	payload := `{
	     "Contacts": [
	       {
	   	  "ContactID":"48d47f54-f562-4618-a386-cf9b8a5c3963",
	         "ContactNumber":"17",
	         "Name": "Lin In",
	         "FirstName": "Lin",
	         "LastName": "In",
	         "EmailAddress": "ling@foxmail.com",

	         "PaymentTerms": {
	           "Bills": {
	             "Day": 7,
	             "Type": "OFCURRENTMONTH"
	           },
	           "Sales": {
	             "Day": 7,
	             "Type": "DAYSAFTERBILLMONTH"
	           }
	         }
	       }
	     ]
	   }`
	*/
	res, _, err := lib.Request("POST", api, []byte(payload), headers)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("res is nil")
	}
	a := lib.ToTypeMapByString(*res)
	Contacts := lib.ToTypeList(a.Get("Contacts"))
	for k, _ := range Contacts {
		val := InterfaceToString(Contacts[k])

		return &val, nil
	}
	return nil, errors.New("Xero response is error")
}

func (c *XeroUsecase) GetContacts() error {

	headers, err := c.Headers()
	if err != nil {
		return err
	}

	api := fmt.Sprintf("%s/api.xro/2.0/Contacts", c.conf.Xero.ApiUrl)
	a, _, err := lib.Request("GET", api, nil, headers)
	lib.DPrintln(err)
	lib.DPrintln(a)
	return nil
}

func (c *XeroUsecase) GetInvoice(invoiceId string) error {
	api := fmt.Sprintf("%s/api.xro/2.0/Invoices/%s", c.conf.Xero.ApiUrl, invoiceId)

	headers, err := c.Headers()
	if err != nil {
		return err
	}
	_, _, err = lib.Request("GET", api, nil, headers)
	return err
}

//
//func (c *XeroUsecase) DoBizCreateInvoice(tClient *TData) error {
//	InvoiceID, InvoiceNumber, err := c.BizCreateInvoice(tClient)
//	if err != nil {
//		return err
//	}
//	clientId := tClient.CustomFields.NumberValueByNameBasic("id")
//	return c.LogUsecase.SaveLog(clientId, Log_FromType_Xero_CreateInvoice, map[string]interface{}{
//		"InvoiceID":     InvoiceID,
//		"InvoiceNumber": InvoiceNumber,
//	})
//}

func (c *XeroUsecase) BizAmCreateInvoice(tClientCase *TData) (InvoiceID string, InvoiceNumber string, err error) {
	if tClientCase == nil {
		return "", "", errors.New("tClientCase is nil.")
	}
	clientCaseFields := tClientCase.CustomFields

	newRating := clientCaseFields.NumberValueByNameBasic("new_rating")
	if newRating <= 0 {
		return "", "", errors.New("newRating is wrong. " + InterfaceToString(tClientCase.Id()))
	}

	usePrimaryCaseCalc := true
	if HasEnabledPrimaryCase(clientCaseFields.TextValueByNameBasic("client_gid")) {
		usePrimaryCaseCalc, _, err = c.FeeUsecase.UsePrimaryCaseCalc(tClientCase)
		if err != nil {
			return "", "", err
		}
	}
	amountStr := clientCaseFields.TextValueByNameBasic(FieldName_am_invoice_amount)
	amountFloat, err := strconv.ParseFloat(amountStr, 10)
	if err != nil {
		return "", "", err
	}
	amount := int(amountFloat)
	//math.Ceil(amountFloat)
	if amount <= 0 {
		return "", "", errors.New(fmt.Sprintf("The \"AM Invoice Amount value\" for CaseId %d is 0 ", tClientCase.Id()))
	}
	contractSource := clientCaseFields.TextValueByNameBasic(FieldName_ContractSource)
	if contractSource != ContractSource_AM {
		return "", "", errors.New(fmt.Sprintf("The client's contract is not AM.  %d", tClientCase.Id()))
	}

	tContact, _, err := c.DataComboUsecase.Client(clientCaseFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		return "", "", err
	}
	if tContact == nil {
		return "", "", errors.New("tContact is nil")
	}

	ContactID, err := c.BizCreateOrUpdateContact(tContact)
	if err != nil {
		return "", "", err
	}

	CreateInvoiceRes, err := c.CreateInvoice(float64(amount), ContactID, Xero_AM_BrandingThemeID)
	if err != nil {
		return "", "", err
	}
	if CreateInvoiceRes == nil {
		return "", "", errors.New("CreateInvoiceRes is nil.")
	}
	CreateInvoiceResMap := lib.ToTypeMapByString(*CreateInvoiceRes)
	if CreateInvoiceResMap.GetString("Status") != "OK" {
		return "", "", errors.New("Status:" + CreateInvoiceResMap.GetString("Status") + " is not ok.")
	}
	Invoices := lib.ToTypeList(CreateInvoiceResMap.Get("Invoices"))
	for k, _ := range Invoices {
		InvoiceID = Invoices[k].GetString("InvoiceID")
		InvoiceNumber = Invoices[k].GetString("InvoiceNumber")

		if HasEnabledPrimaryCase(clientCaseFields.TextValueByNameBasic("client_gid")) {
			if usePrimaryCaseCalc { // 说明是primary case需要入库
				entity := make(lib.TypeMap)
				entity.Set(FieldName_is_primary_case, Is_primary_case_YES)
				entity.Set(FieldName_gid, clientCaseFields.TextValueByNameBasic("gid"))
				c.log.Info("EnabledPrimaryCase: ", entity)
				c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry(entity), FieldName_gid, nil)
			}
		}

		return InvoiceID, InvoiceNumber, nil
	}
	return "", "", errors.New("InvoiceID and InvoiceNumber is error.")
}

func (c *XeroUsecase) BizCreateInvoice(tClientCase *TData) (InvoiceID string, InvoiceNumber string, err error) {
	if tClientCase == nil {
		return "", "", errors.New("tClientCase is nil.")
	}
	clientCaseFields := tClientCase.CustomFields
	clientCaseId := clientCaseFields.NumberValueByNameBasic("id")

	newRating := clientCaseFields.NumberValueByNameBasic("new_rating")
	if newRating <= 0 {
		return "", "", errors.New("newRating is wrong.")
	}

	var amount int
	usePrimaryCaseCalc := true
	if HasEnabledPrimaryCase(clientCaseFields.TextValueByNameBasic("client_gid")) {
		usePrimaryCaseCalc, _, err = c.FeeUsecase.UsePrimaryCaseCalc(tClientCase)
		if err != nil {
			return "", "", err
		}
		if !usePrimaryCaseCalc {
			amount, err = c.FeeUsecase.NotPrimaryCaseAmount(tClientCase, newRating)
			if err != nil {
				return "", "", err
			}
		}
	}

	if usePrimaryCaseCalc {
		clientCaseContractBasicDataVo, err := c.ClientCaseUsecase.ClientCaseContractBasicDataVoById(clientCaseId)
		if err != nil {
			return "", "", err
		}
		effectiveCurrentRating := clientCaseContractBasicDataVo.EffectiveCurrentRating
		if effectiveCurrentRating < 0 {
			return "", "", errors.New("effectiveCurrentRating is wrong.")
		}
		amount, err = c.FeeUsecase.InvoiceAmount(clientCaseId, clientCaseContractBasicDataVo.ActiveDuty, effectiveCurrentRating, newRating)
		if err != nil {
			return "", "", err
		}
	}
	tContact, _, err := c.DataComboUsecase.Client(clientCaseFields.TextValueByNameBasic("client_gid"))

	if err != nil {
		return "", "", err
	}
	if tContact == nil {
		return "", "", errors.New("tContact is nil")
	}

	ContactID, err := c.BizCreateOrUpdateContact(tContact)
	if err != nil {
		return "", "", err
	}

	if amount <= 0 {
		return "", "", errors.New(fmt.Sprintf("The amount calculated by CaseId: %d is 0 ", tClientCase.Id()))
	}

	CreateInvoiceRes, err := c.CreateInvoice(float64(amount), ContactID, Xero_VBC_BrandingThemeID)
	if err != nil {
		return "", "", err
	}
	if CreateInvoiceRes == nil {
		return "", "", errors.New("CreateInvoiceRes is nil.")
	}
	CreateInvoiceResMap := lib.ToTypeMapByString(*CreateInvoiceRes)
	if CreateInvoiceResMap.GetString("Status") != "OK" {
		return "", "", errors.New("Status:" + CreateInvoiceResMap.GetString("Status") + " is not ok.")
	}
	Invoices := lib.ToTypeList(CreateInvoiceResMap.Get("Invoices"))
	for k, _ := range Invoices {
		InvoiceID = Invoices[k].GetString("InvoiceID")
		InvoiceNumber = Invoices[k].GetString("InvoiceNumber")

		if HasEnabledPrimaryCase(clientCaseFields.TextValueByNameBasic("client_gid")) {
			if usePrimaryCaseCalc { // 说明是primary case需要入库
				entity := make(lib.TypeMap)
				entity.Set(FieldName_is_primary_case, Is_primary_case_YES)
				entity.Set(FieldName_gid, clientCaseFields.TextValueByNameBasic("gid"))
				c.log.Info("EnabledPrimaryCase: ", entity)
				c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry(entity), FieldName_gid, nil)
			}
		}

		return InvoiceID, InvoiceNumber, nil
	}
	return "", "", errors.New("InvoiceID and InvoiceNumber is error.")
}

/*
	{
	  "Id": "e413260c-1dab-4d07-87c8-a9297fcc87b0",
	  "Status": "OK",
	  "ProviderName": "VBC",
	  "DateTimeUTC": "\/Date(1752500079517)\/",
	  "BrandingThemes": [
	    {
	      "BrandingThemeID": "24badf25-21ae-49c6-b96f-817402428a79",
	      "Name": "Standard VBC",
	      "LogoUrl": "https://in.xero.com/logo?id=ZXcwS0lDQWlieUk2SUNJME16bGlOamcxWWkwd05EbG1MVFJqTnpNdE9UZGxaaTFrWkRkbU5EZzFabVZrTTJRaUxBMEtJQ0FpWmlJNklDSTRZbVExTmpNek5pMWhZalF6TFRSbU1tUXRZamt3TXkxa05HUmpOMlEzTkRJMk5tUWlEUXA5LXZRdVh4WUduWVpnPQ",
	      "Type": "INVOICE",
	      "SortOrder": 0,
	      "CreatedDateUTC": "\/Date(1708483081437+0000)\/"
	    },
	    {
	      "BrandingThemeID": "79c4c18e-a119-476a-b714-5894732ca231",
	      "Name": "VBC for AM",
	      "LogoUrl": "https://in.xero.com/logo?id=ZXcwS0lDQWlieUk2SUNJME16bGlOamcxWWkwd05EbG1MVFJqTnpNdE9UZGxaaTFrWkRkbU5EZzFabVZrTTJRaUxBMEtJQ0FpWmlJNklDSmpObUl5TkdOa01TMDRNelV5TFRRNVpXVXRPREprWVMxaFkySXpZakpoTURVMVltWWlEUXA5LWJMMkFhZ0V4VlFFPQ",
	      "Type": "INVOICE",
	      "SortOrder": 1,
	      "CreatedDateUTC": "\/Date(1751826249670+0000)\/"
	    }
	  ]
	}
*/
func (c *XeroUsecase) BrandingThemes() (*string, error) {
	api := fmt.Sprintf("%s/api.xro/2.0/BrandingThemes", c.conf.Xero.ApiUrl)
	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}
	res, _, err := lib.Request("GET", api, nil, headers)
	return res, err
}

const (
	Xero_VBC_BrandingThemeID = "24badf25-21ae-49c6-b96f-817402428a79"
	Xero_AM_BrandingThemeID  = "79c4c18e-a119-476a-b714-5894732ca231"
)

/*
https://developer.xero.com/documentation/api/accounting/invoices
*/
func (c *XeroUsecase) CreateInvoice(UnitAmount float64, ContactID string, BrandingThemeID string) (*string, error) {
	/*
		说明：
		1. 一样的数据，在没有InvoiceNumber时，会产生新数据
		2. 通过设置InvoiceID 和InvoiceNumber，会更新结果
	*/
	if UnitAmount <= 0 {
		return nil, errors.New("UnitAmount is wrong")
	}
	api := fmt.Sprintf("%s/api.xro/2.0/Invoices", c.conf.Xero.ApiUrl)

	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}

	Date := time.Now()
	DueDate := Date.Add(6 * 24 * time.Hour)
	//DueDate := Date
	// AccountCode 4300
	// TaxType
	// ContactID 48d47f54-f562-4618-a386-cf9b8a5c3963 gengling.liao@hotmail.com
	// 金额不要设置太大，防止误操作
	// InvoiceNumber: 创建时，不要添加， 更新时添写查看情况

	/*
	   "InvoiceID": "a5a650cb-daf8-4dbd-9a70-c958301cbc0d",
	      "InvoiceNumber": "INV-2403043",
	*/

	str := `{
  "Invoices": [
    {
	
      "Type": "ACCREC",
      "Contact": {
        "ContactID": "` + ContactID + `"
      },
      "LineItems": [
        {
          "Description": "Client Service Fee",
          "Quantity": 1,
          "UnitAmount": ` + InterfaceToString(UnitAmount) + `,
          "AccountCode": "4300",
          "TaxType": "NONE",
          "LineAmount": ` + InterfaceToString(UnitAmount) + `
        }
      ],
      "Date": "` + Date.Format("2006-01-02") + `",
      "DueDate": "` + DueDate.Format("2006-01-02") + `",
      "Reference": "",
	  "CurrencyCode": "USD",
      "Status": "DRAFT",
	  "BrandingThemeID":"` + BrandingThemeID + `"	
    }
  ]
}`
	// https://developer.xero.com/documentation/api/accounting/invoices

	res, _, err := lib.Request("POST", api, []byte(str), headers)

	c.LogUsecase.SaveLog(0, Log_FromType_Xero_InvoicesApi, map[string]interface{}{
		"str": str,
		"res": res,
		"err": err,
	})

	return res, err
}
