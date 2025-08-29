package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"vbc/internal/conf"
	"vbc/lib"
)

type AdobepdfUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
}

func NewAdobepdfUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *AdobepdfUsecase {
	uc := &AdobepdfUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	return uc
}

func (c *AdobepdfUsecase) Host() string {
	return "https://pdf-services.adobe.io"
}

func (c *AdobepdfUsecase) ClientId() string {
	return "893715e9cefe4951842e009f89be091b"
}

func (c *AdobepdfUsecase) SecretKey() string {
	return "p8e-DMz8wH0bRnsPcCyMnDpyjQUzwOud39z-"
}

func (c *AdobepdfUsecase) Token() string {
	return "eyJhbGciOiJSUzI1NiIsIng1dSI6Imltc19uYTEta2V5LWF0LTEuY2VyIiwia2lkIjoiaW1zX25hMS1rZXktYXQtMSIsIml0dCI6ImF0In0.eyJpZCI6IjE3NDQwNzc4MzU3NzJfODA5MTNhNmItYzZlNi00Y2E0LWEzMmEtZjYyYjE2NDJlYTA1X3VlMSIsIm9yZyI6IjE5Q0IyMDhCNjdGNDZGODAwQTQ5NUM1NEBBZG9iZU9yZyIsInR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJjbGllbnRfaWQiOiI4OTM3MTVlOWNlZmU0OTUxODQyZTAwOWY4OWJlMDkxYiIsInVzZXJfaWQiOiIxODdEMjE3RjY3RjQ3NkMyMEE0OTVFQUZAdGVjaGFjY3QuYWRvYmUuY29tIiwiYXMiOiJpbXMtbmExIiwiYWFfaWQiOiIxODdEMjE3RjY3RjQ3NkMyMEE0OTVFQUZAdGVjaGFjY3QuYWRvYmUuY29tIiwiY3RwIjozLCJtb2kiOiJmNDhjMDZlIiwiZXhwaXJlc19pbiI6Ijg2NDAwMDAwIiwiY3JlYXRlZF9hdCI6IjE3NDQwNzc4MzU3NzIiLCJzY29wZSI6IkRDQVBJLG9wZW5pZCxBZG9iZUlEIn0.I5QZJseWeDEbm7x42haEXhn8i2zugkgUoN_7EG5-g7NBQZCtbcW9HSRA8SV8AgpzgsrrKb5NcTnfyZW6Sw-wwiNjT1GUrXwTZ-bQ-OrrPxO2N29HbnddlaMgurdhQm7Irk8y0BeialvASi53XX-k0OjmLykVQk6ZWjvnU3mPlaXsWQavrUPvz2tJyjLekpWx3F_Ce6gu2HmUIgN5WfqKJdXWEwWpzeqNFKFluiUaCqhd5vu8UyTUaQVU4HtMFHjjAHCYN6kg_ATWlbFT3gRFZGiBtOPciWmZ7ayj7sTG6Qe7vttg0ZIQmjdPTSliUEIqt7TmLhsgEFa2BIh1SNIVpQ"
}

// GetUploadPreSignedURIAssets Get upload pre-signed URI mediaType application/pdf
// 返回：{"assetID":"urn:aaid:AS:UE1:8429a1cd-0318-4679-af1f-453c340254d8","uploadUri":"https://dcplatformstorageservice-prod-us-east-1.s3-accelerate.amazonaws.com/893715e9cefe4951842e009f89be091b_187D217F67F476C20A495EAF%40techacct.adobe.com/50506501-3aa0-41c2-b42b-a566f5b04140?X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPD%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLWVhc3QtMSJHMEUCIEZmlm2KKLo%2FjOp4TYWMSpV54%2FSCsXW2gxEqg1OiKrXRAiEA96SA2APFkZJocldmcdnC44dAbfM9XOLT1OXBTNfsxK4qmAUIaRAAGgw0MjA1MzM0NDU5ODIiDBWrlk9gEBKBZK2mFSr1BDRkSyY9Pj%2Fa8DYY%2FYxUnj%2BWEzL%2BwISScXn5D5pat6zD5Q2zNZj7fC8zwbrWfW8WIO13S9iOlzoN1c0ZvZ2QiRaa44LJXEu8enAD3%2B2nsf4CgnW%2BWMYXkMestVXVUcX2Lg%2FlHiKQSOtisC0B4M5EFMKjkmCE%2BiwtubrbvwviMkY8DaJfgCIEgAOXYLoDJRwR%2B%2B82dR%2FfgXgynhZYu11yejv2cz9xxlRUkhT23NKxBWJM6KJGzJCSIwdHdVlISZ1E5ITV2IBHWUKgYzMOcFjeIfgfVImojiaR5XVjedC0WvB14lrstNwuqcdpYvwM63pMvYxzjMxH2nMhwvo%2FegU7XMGaQaapfnfEYgoOF8vHmtcsb7emkP9oM3bwGv7J%2FpSdjxPFCjp64CfJTyTKjzTiqe5h%2FxNY%2BUDYMpcT5QtBbrxXVBiUjaApyOYDssQlqmvLwvb0xhBgBCeXhoh8KYP3BDMtwNMoAG8t4bmuRHN7wJ%2Fz35XeaD9LbDK%2FRlbNpDtv3UyE%2Fou6XIlDI2vDxfeuvpRG5lKSdhg2yjiLsrbJ7yCgpM5N3I3YS2YfJf3b5mHXPr0vnpo1aZZywz1yKhSBCJ1gaCAB1H%2FNcRkGipIg98FCbuCJ8TTOO45rnO9wGj7MQP7TpWxnqiwcYx20NJ2P9UCsPp%2BM4IYSHNuludtKRXVMIBHUgLXr1NIxpkh9p9n4k9U%2BV%2Burxtyy5UikbBg3vouAo6mSoJVPE23qmih7hxTSxy%2FbaS%2FJnL%2BbQZJf61ZwbNFGVMdJSyxY%2F63zoTG2mdYu8IcOlOI%2Feg7q2u2EmkaPyM3IVcChGlKlkxCzNnTHUEfzqxYiMKLK0b8GOpsB6VPbsgUnHAhe62Kmrq6vmtKNU4Is8iQzlGS2BVlZZ9VZw5ZkwZKqI9JpLEnPUf9se%2FK11%2FDaJgh7dYmpZkJOAv0yHxNNNhGwThCN%2B3WUX62SSL0kNT8HuBJ4hpyD69PaDLjZIZk7fDrEz4Ni0HwM00zID0uB5HNBPpddWLqXsININ4wTUdZRkgbuZrKsj6NGOzLjUG6UcIPS4V0%3D\u0026X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Date=20250408T020514Z\u0026X-Amz-SignedHeaders=content-type%3Bhost\u0026X-Amz-Expires=3600\u0026X-Amz-Credential=ASIAWD2N7EVPAR2VKQMJ%2F20250408%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Signature=c0f6cfe824df9c3c8fb7624da8ef2c442ae771f93eec00cf2a977d27f42bc66e"}
func (c *AdobepdfUsecase) GetUploadPreSignedURIAssets(mediaType string) (lib.TypeMap, error) {
	url := fmt.Sprintf("%s/assets", c.Host())
	params := make(lib.TypeMap)
	params.Set("mediaType", mediaType)
	res, _, err := lib.Request("POST", url, params.ToBytes(), map[string]string{
		"X-API-Key":     c.ClientId(),
		"Authorization": "Bearer " + c.Token(),
		"Content-Type":  "application/json",
	})
	if err != nil {
		return nil, err
	}
	if res != nil {
		return lib.ToTypeMapByString(*res), nil
	}
	return nil, nil
}

// Put
/*
curl --location -g --request PUT 'https://dcplatformstorageservice-prod-us-east-1.s3-accelerate.amazonaws.com/b37fd583-1ab6-4f49-99ef-d716180b5de4?X-Amz-Security-Token={{Placeholder for X-Amz-Security-Token}}&X-Amz-Algorithm={{Placeholder for X-Amz-Algorithm}}&X-Amz-Date={{Placeholder for X-Amz-Date}}&X-Amz-SignedHeaders={{Placeholder for X-Amz-SignedHeaders}}&X-Amz-Expires={{Placeholder for X-Amz-Expires}}&X-Amz-Credential={{Placeholder for X-Amz-Credential}}&X-Amz-Signature={{Placeholder for X-Amz-Signature}}' \
--header 'Content-Type: application/pdf' \
--data-binary '@{{Placeholder for file path}}'
contentType: application/pdf
*/
func (c *AdobepdfUsecase) Put(uri string, reader io.Reader, contentType string, contentLength int64) error {
	a, b, err := lib.RequestStream("PUT", uri, reader, map[string]string{
		"Content-Type": contentType,
	}, contentLength)
	lib.DPrintln(a, b)
	return err
}

func (c *AdobepdfUsecase) ExportPDFFormData(assetID string) (lib.TypeMap, error) {

	url := fmt.Sprintf("%s/operation/getformdata", c.Host())
	params := `{
  "assetID": "` + assetID + `",
  "notifiers": [
    {
      "type": "CALLBACK",
      "data": {
        "url": "https:///webhook/post_source?from=adobeExportPdf",
        "headers": {
          "x-api-key": "dummykey",
          "access-token": "dummytoken"
        }
      }
    }
  ]
}`

	res, _, err := lib.Request("POST", url, []byte(params), map[string]string{
		"X-API-Key":     c.ClientId(),
		"Authorization": "Bearer " + c.Token(),
		"Content-Type":  "application/json",
	})
	if err != nil {
		return nil, err
	}
	if res != nil {
		return lib.ToTypeMapByString(*res), nil
	}
	return nil, nil
}

func (c *AdobepdfUsecase) ImportPDFFormData(assetID string) (lib.TypeMap, error) {
	url := fmt.Sprintf("%s/operation/setformdata", c.Host())

	params := `{
  "assetID": "` + assetID + `",
  "jsonFormFieldsData": {"F[0]": {
		    "Page_1[0]": {
		        "MailingAddress_NumberAndStreet[0]": "123 Port Trinity",
		        "MailingAddress_City[0]": "123 Vista",
		        "DOB_Day[0]": "12",
		        "EMAIL_ADDRESS[1]": "com",
		        "REMARKS[0]": "12344",
		        "SocialSecurityNumber_LastFourNumbers[0]": "",
		        "Veterans_Beneficiary_First_Name[0]": "Mar123 ",
		        "TelephoneNumber_FirstThreeNumbers[0]": "123",
		        "SocialSecurityNumber_FirstThreeNumbers[0]": "123",
		        "TelephoneNumber_SecondThreeNumbers[0]": "227",
		        "International_Phone_Number[0]": "",
		        "Middle_Initial1[0]": "V",
		        "Veterans_Service_Number_If_Applicable[0]": "",
		        "TelephoneNumber_LastFourNumbers[0]": "7362",
		        "Veterans_DOB_Month[0]": "12",
		        "VA_File_Number_If_Applicable[0]": "561530503",
		        "Last_Name[0]": "Sese",
		        "EMAIL_ADDRESS[0]": "red2robinchel@yahoo.",
		        "MailingAddress_ZIPOrPostalCode_FirstFiveNumbers[0]": "91913",
		        "MailingAddress_ZIPOrPostalCode_LastFourNumbers[0]": "",
		        "DOB_Year[0]": "1959",
		        "MailingAddress_Country[0]": "US",
		        "MailingAddress_ApartmentOrUnitNumber[0]": "",
		        "MailingAddress_StateOrProvince[0]": "CA",
		        "SocialSecurityNumber_SecondTwoNumbers[0]": "53"
		    },
		    "#subform[1]": {
		        "Date_Signed_Month[0]": "10",
		        "SocialSecurityNumber_FirstThreeNumbers[0]": "561",
		        "Date_Signed_Day[0]": "11",
		        "REMARKS[0]": "blood sugar ",
		        "SocialSecurityNumber_LastFourNumbers[0]": "0503",
		        "Date_Signed_Year[0]": "2024",
		        "SocialSecurityNumber_SecondTwoNumbers[0]": "53"
		    }
		}},
  "notifiers": [
    {
      "type": "CALLBACK",
      "data": {
        "url": "https:///webhook/post_source?from=adobeImportPdfFormData",
        "headers": {
          "x-api-key": "dummykey",
          "access-token": "dummytoken"
        }
      }
    }
  ]
}`
	res, httpCode, err := lib.Request("POST", url, []byte(params), map[string]string{
		"X-API-Key":     c.ClientId(),
		"Authorization": "Bearer " + c.Token(),
		"Content-Type":  "application/json",
	})

	if err != nil {
		return nil, err
	}
	if res != nil {
		c.log.Info("res: ", *res, " httpCode: ", httpCode)
		return lib.ToTypeMapByString(*res), nil
	}
	return nil, nil
}
