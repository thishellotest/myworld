package tests

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"testing"
	"time"
	"vbc/configs"
	"vbc/internal/biz"
	"vbc/internal/config_vbc"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/esign"
	"vbc/lib/esign/v2.1/envelopes"
	"vbc/lib/uuid"
)

func TestTmp1(t *testing.T) {

	tasks, err := UT.TUsecase.ListByCond(biz.Kind_client_tasks, And(Eq{"biz_deleted_at": 0, "what_id_gid": ""},
		In("status",
			config_zoho.ClientTaskStatus_Waitingforinput,
			config_zoho.ClinetTaskStatus_NotStarted,
			config_zoho.ClientTaskStatus_Deferred,
			config_zoho.ClientTaskStatus_InProgress),
		Expr("status not like '"+biz.ClientTaskSubject_ITFExpirationWithPrefix+"%'"),
		Expr("status != ''"),
		Expr("status != ''")))
	lib.DPrintln(tasks, err)
}

func Test_uuid(t *testing.T) {
	a := uuid.UuidWithoutStrike()
	lib.DPrintln(a)
	a = uuid.UuidWithoutStrike()
	lib.DPrintln(a)
}

//
//func getLocalCredential() (*esign.OAuth2Credential, error) {
//	os.Setenv("DOCUSIGN_Token", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IjY4MTg1ZmYxLTRlNTEtNGNlOS1hZjFjLTY4OTgxMjIwMzMxNyJ9.eyJUb2tlblR5cGUiOjUsIklzc3VlSW5zdGFudCI6MTcwNDE2NDA4NCwiZXhwIjoxNzA0MTkyODg0LCJVc2VySWQiOiI0MTI4YTJmMy0zYmE2LTQyN2ItYTNhZi1iOWZkNzgzZmNiNjkiLCJzaXRlaWQiOjEsInNjcCI6WyJpbXBlcnNvbmF0aW9uIiwiZXh0ZW5kZWQiLCJzaWduYXR1cmUiLCJjb3JzIiwiY2xpY2subWFuYWdlIiwiY2xpY2suc2VuZCIsIm9yZ2FuaXphdGlvbl9yZWFkIiwiZ3JvdXBfcmVhZCIsInBlcm1pc3Npb25fcmVhZCIsInVzZXJfcmVhZCIsInVzZXJfd3JpdGUiLCJhY2NvdW50X3JlYWQiLCJkb21haW5fcmVhZCIsImlkZW50aXR5X3Byb3ZpZGVyX3JlYWQiLCJ1c2VyX2RhdGFfcmVkYWN0IiwiZHRyLnJvb21zLnJlYWQiLCJkdHIucm9vbXMud3JpdGUiLCJkdHIuZG9jdW1lbnRzLnJlYWQiLCJkdHIuZG9jdW1lbnRzLndyaXRlIiwiZHRyLnByb2ZpbGUucmVhZCIsImR0ci5wcm9maWxlLndyaXRlIiwiZHRyLmNvbXBhbnkucmVhZCIsImR0ci5jb21wYW55LndyaXRlIiwicm9vbV9mb3JtcyIsIm5vdGFyeV93cml0ZSIsIm5vdGFyeV9yZWFkIiwic3ByaW5nX3JlYWQiLCJzcHJpbmdfd3JpdGUiXSwiYXVkIjoiNDNmZjNiZmMtNGEzMS00YTViLTkxNjQtNDM2NjU3MmVlZGZkIiwiYXpwIjoiNDNmZjNiZmMtNGEzMS00YTViLTkxNjQtNDM2NjU3MmVlZGZkIiwiaXNzIjoiaHR0cHM6Ly9hY2NvdW50LWQuZG9jdXNpZ24uY29tLyIsInN1YiI6IjQxMjhhMmYzLTNiYTYtNDI3Yi1hM2FmLWI5ZmQ3ODNmY2I2OSIsImFtciI6WyJpbnRlcmFjdGl2ZSJdLCJhdXRoX3RpbWUiOjE3MDQxNjQwNzcsInB3aWQiOiJkMWY5MDkzMi1mYzczLTQwODUtOGUxYy1iYjcxMWY4ZTI2NWIifQ.vrLIGlUYkpa6F79mclt21X2QieCsro-l58KT6S2HF2KSf0Y15Pk0mPtKAkf-6JA8-flOSAji4CHtpXKMVgAwEdvKZpwFFK_ZpAwRe18r4PruwI727xUAvI0z8NOJKOY4CWbAmBXwt_qpJnCDYXLmul5qO2hXcG5XUV2crLnQyPW_qiHwNCioQ76Na8lYkuFUgwyIgPJC1rCxhTYNBRY3h0CdEOn0IpBf25AMB8C9fPzlO8s0WdmLM3hvG1xYCsclZ5z0V78fhKTU_JvmFuXJQwnt42Wp1r8UbAB56Rpx3HxmdwphoEJV8c45pLvhpIf5a3IC4Zcr2R3K3NOM68DIYw")
//	if tk, ok := os.LookupEnv("DOCUSIGN_Token"); ok {
//		acctID, _ := os.LookupEnv("DOCUSIGN_AccountID")
//		return esign.TokenCredential(tk, true).WithAccountID(acctID), nil
//	}
//
//	if jwtConfigJSON, ok := os.LookupEnv("DOCUSIGN_JWTConfig"); ok {
//		jwtAPIUserName, ok := os.LookupEnv("DOCUSIGN_JWTAPIUser")
//		if !ok {
//			return nil, fmt.Errorf("expected DOCUSIGN_JWTAPIUser environment variable with DOCUSIGN_JWTConfig=%s", jwtConfigJSON)
//		}
//
//		buffer, err := ioutil.ReadFile(jwtConfigJSON)
//		if err != nil {
//			return nil, fmt.Errorf("%s open: %v", jwtConfigJSON, err)
//		}
//		var cfg *esign.JWTConfig
//		if err = json.Unmarshal(buffer, &cfg); err != nil {
//			return nil, fmt.Errorf("%v", err)
//		}
//		return cfg.Credential(jwtAPIUserName, nil, nil)
//	}
//	return nil, nil
//}

type Credential struct {
	oauth2Token *biz.Oauth2TokenEntity
}

func (c *Credential) AuthDo(ctx context.Context, op *esign.Op) (*http.Response, error) {

	if op.Version == nil {
		return nil, errors.New("no api version set for op")
	}
	req, err := op.CreateRequest()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.oauth2Token.TokenType+" "+c.oauth2Token.AccessToken)
	var rawUrl string

	/*
		demoHost = "demo.docusign.net"
		baseHost = "www.docusign.net"
	*/
	//fmt.Println(op.Path, op.QueryOpts)
	if configs.IsDev() {
		rawUrl = fmt.Sprintf("https://demo.docusign.net/restapi/v2.1/accounts/%s/%s", "60e36fd1-0481-40c3-b7ca-c1ca4776bd87", op.Path)
	} else {
		rawUrl = fmt.Sprintf("https://www.docusign.net/%s", op.Path)
	}
	//fmt.Println("rawUrl:", rawUrl)

	//a, err := lib.Request("GET", rawUrl, nil, map[string]string{
	//	"Authorization": c.oauth2Token.TokenType + " " + c.oauth2Token.AccessToken,
	//})
	//fmt.Println(c.oauth2Token.TokenType, c.oauth2Token.AccessToken)
	//lib.DPrintln("+++++", a, err)
	//fmt.Println(*a, err)

	// /restapi/v2.1/accounts/{accountId}/envelopes/{envelopeId}/documents

	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	//fmt.Println(u.Host, u.Path, u.RawPath)
	req.URL.Scheme = u.Scheme
	req.URL.Host = u.Host
	req.URL.Path = u.Path
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	return client.Do(req)

}

func NewCredential() *Credential {
	return &Credential{}
}

func Test_ac(t *testing.T) {

	// 06f2eab4-295b-47d1-91f2-d6f55a78deb2 company
	// 43ff3bfc-4a31-4a5b-9164-4366572eedfd home
	oauth2Token, _ := UT.Oauth2TokenUsecase.GetByClientId("06f2eab4-295b-47d1-91f2-d6f55a78deb2")
	//cred.oauth2Token = oauth2Token
	//envelopes.ListStatusChangesOp
	cred := biz.NewDocuSignCredential(oauth2Token, "60e36fd1-0481-40c3-b7ca-c1ca4776bd87")
	srv := envelopes.New(cred)

	now := time.Now()
	now = now.AddDate(0, 0, -110)
	a, err := srv.ListStatusChanges().FromDate(now).Do(context.Background())
	fmt.Println(err)
	lib.DPrintln(a)
}

func Test_DocumentsList1(t *testing.T) {

	cred := NewCredential()
	// 06f2eab4-295b-47d1-91f2-d6f55a78deb2
	// 43ff3bfc-4a31-4a5b-9164-4366572eedfd
	oauth2Token, _ := UT.Oauth2TokenUsecase.GetByClientId("06f2eab4-295b-47d1-91f2-d6f55a78deb2")
	cred.oauth2Token = oauth2Token
	//envelopes.ListStatusChangesOp
	srv := envelopes.New(cred)
	a, e := srv.DocumentsList("2adc1cbb-5616-4810-bfa4-b7e658cb28ce").Do(context.Background())

	// /envelopes/2adc1cbb-5616-4810-bfa4-b7e658cb28ce/attachments
	// /envelopes/2adc1cbb-5616-4810-bfa4-b7e658cb28ce/documents
	//a, e := srv.DocumentsGet("/envelopes/2adc1cbb-5616-4810-bfa4-b7e658cb28ce/attachments",
	//	"2adc1cbb-5616-4810-bfa4-b7e658cb28ce").Do(context.Background())

	lib.DPrintln(a, e)
	//l, err := srv.Get("2adc1cbb-5616-4810-bfa4-b7e658cb28ce").
	//	Do(context.Background())
	//lib.DPrintln(l)
	//lib.DPrintln(err)
}

func Test_accc(t *testing.T) {
	//a := UT.TaskUsecase.RunTaskJob
	//d := a(context.Background())
}

func Test_aaa(t *testing.T) {
	str := `{
	"First Name": "NN",
	"Last Name": "LL2",
	"Street Address": "address val 1",
	"City": "SH",
	"State": "Arizona",
	"Zip Code": "20021",
	"Phone Number": "123-233-4421",
	"SSN": "123-45-2234",
	"Date of Birth": "2023-09-07",
	"Overall Rating": "10",
	"Branch of Service ": ["Army", "Navy", "Marine Corps", "Air Force", "Space Force", "Coast Guard", "National Oceanic and Atmospheric Administration", "Public Health Service", "Army National Guard", "Air National Guard"],
	"Have you retired from US Military services?": "Yes",
	"Military Toxic Exposures": ["Agent Orange Exposure", "Gulf War Illness", "Burn Pits and Other Airborne Hazards", "Illness Due to Toxic Drinking Water at Camp Lejeune", "\"Atomic Veterans\" and Radiation Exposure", "Amyotrophic Lateral Sclerosis (ALS)"]
}`
	a := lib.ToTypeMapByString(str)

	asanaField := config_vbc.GetAsanaCustomFields()

	BranchofService := a.Get("Branch of Service ")
	BranchofServiceList := lib.InterfaceToTDef[[]string](BranchofService, nil)
	for _, v := range BranchofServiceList {
		gid := asanaField.GetByName("Branch").GetEnumGidByName(v)
		if gid != "" {
			fmt.Println(gid)
			break
		}
	}
	return
	c := a.Get("Military Toxic Exposures")
	d := lib.InterfaceToTDef[[]string](c, nil)

	for _, v := range d {
		gid := asanaField.GetByName(v).GetEnumGidByName("Yes")
		fmt.Println(gid, "====")
	}

	lib.DPrintln(d)
}

func Test_time1(T *testing.T) {
	ti, err := time.ParseInLocation(time.DateOnly, "2024-09-10", configs.GetVBCDefaultLocation())
	lib.DPrintln(err)
	currentTime := time.Now().In(configs.GetVBCDefaultLocation())
	currentTimeStr := currentTime.Format(time.DateOnly)
	currentTimeStr = "2024-09-11"
	currentTime, _ = time.ParseInLocation(time.DateOnly, currentTimeStr, configs.GetVBCDefaultLocation())
	if currentTime.After(ti) {
		lib.DPrintln("超时")
	} else {
		lib.DPrintln("11")
	}
}

func TestTemp2(t *testing.T) {
	tasks, err := UT.TUsecase.ListByCond(biz.Kind_client_tasks, And(Eq{"biz_deleted_at": 0, "what_id_gid": ""},
		In("status",
			config_zoho.ClientTaskStatus_Waitingforinput,
			config_zoho.ClinetTaskStatus_NotStarted,
			config_zoho.ClientTaskStatus_Deferred,
			config_zoho.ClientTaskStatus_InProgress),
		Expr("subject like '"+biz.ClientTaskSubject_ITFExpirationWithPrefix+"%'")))
	lib.DPrintln(tasks, err)
}

func Test_GetFilingDateByItfExpiration(t *testing.T) {
	a := biz.GetFilingDateByItfExpiration("2026-07-07")
	lib.DPrintln(a)
	a = biz.GetFilingDateByItfExpiration("2028-02-29")
	lib.DPrintln(a)
}

func Test_GetDiffDaysFilingDateByItfExpiration(t *testing.T) {
	a := biz.GetDiffDaysFilingDateByItfExpiration("2026-07-07", "2026-07-07")
	lib.DPrintln(a)
	a = biz.GetDiffDaysFilingDateByItfExpiration("2028-02-29", "2028-02-29")
	lib.DPrintln(a)
}
