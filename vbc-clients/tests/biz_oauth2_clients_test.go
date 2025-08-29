package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_Oauth2ClientUsecase_zohocrm_AuthUrl(t *testing.T) {
	a, _ := UT.Oauth2ClientUsecase.AuthUrl(biz.Oauth2_AppId_zohocrm)
	fmt.Println(a)
}

// http://localhost:8050/oauth2/callback?app_id=zohocrm&code=1000.2d3d7a1dc79535aaa6976b705b68ca32.13a77b2fec403b97b60f9f8f403a5cef&location=us&accounts-server=https%3A%2F%2Faccounts.zoho.com&
// http://localhost:8050/oauth2/callback?app_id=zohocrm&code=1000.3418f317edc9fbb8769534332fceda24.be122af1603aac13f235f1052d1b5e79&location=us&accounts-server=https%3A%2F%2Faccounts.zoho.com&
func Test_zohocrm_Exchange(t *testing.T) {
	code := "1000.3418f317edc9fbb8769534332fceda24.be122af1603aac13f235f1052d1b5e79"
	a, err := UT.Oauth2ClientUsecase.Exchange(biz.Oauth2_AppId_zohocrm, code)
	fmt.Println(a)
	fmt.Println(err)
	if a != nil {
		er := UT.Oauth2TokenUsecase.UpdateByToken(biz.Oauth2_AppId_zohocrm, a)
		lib.DPrintln(er)
	}
}

// http://localhost:8050/oauth2/callback?app_id=xero
func Test_Oauth2ClientUsecase_xero_AuthUrl(t *testing.T) {
	a, _ := UT.Oauth2ClientUsecase.AuthUrl(biz.Oauth2_AppId_xero)
	fmt.Println(a)
}

func Test_xero_Exchange(t *testing.T) {
	code := "_pd9aSpmvAS_M5b_IAh8763hFK_uCJ89tP4lhEJ9x-A"
	a, err := UT.Oauth2ClientUsecase.Exchange(biz.Oauth2_AppId_xero, code)
	fmt.Println(a)
	fmt.Println(err)
	if a != nil {
		er := UT.Oauth2TokenUsecase.UpdateByToken(biz.Oauth2_AppId_xero, a)
		lib.DPrintln(er)
	}
}

func Test_Oauth2ClientUsecase_google_AuthUrl(t *testing.T) {
	a, _ := UT.Oauth2ClientUsecase.AuthUrl(biz.Oauth2_AppId_google)
	fmt.Println(a)
}

func Test_Oauth2ClientUsecase_vbcapp_AuthUrl(t *testing.T) {
	a, _ := UT.Oauth2ClientUsecase.AuthUrl(biz.Oauth2_AppId_vbcapp)
	fmt.Println(a)
}

func Test_Oauth2ClientUsecase_microsoft_vbcapp_AuthUrl(t *testing.T) {
	a, err := UT.Oauth2ClientUsecase.AuthUrl(biz.Oauth2_AppId_microsoft_vbcapp)
	fmt.Println(err)
	fmt.Println(a)
}

func Test_microsoft_vbcapp_Exchange(t *testing.T) {
	code := "1.AWEBzH5KrmjHfkOUCIn4uZeJqkeAp1p4w2VJmzAReQjWQ5tiAXZhAQ.AgABBAIAAABVrSpeuWamRam2jAF1XRQEAwDs_wUA9P8ce-mBiQj2MBJQktqZ-F9FNLaU1cDASVZSnC_YURi9epHYcNxq5qzZ7-OUIPC_VuCd6CfPQDDfdIxy9N40D9JldNW4ACxa4q4GFUOFSUlOKfB0YEnXlioG3CwtAtvAyeOX-LUrA55QNzip0lB3TzIrstMIhtfHU0mb9clwQXrksrYh_auJmfNcis9Sz6EUIQVoRPXBVTcRRsEOcw7OlYWxWNiTswPVHouq1wPDnK3OdDQ6S3ihe6J6NpqzPp8Xf-npt-vs8BwcvsF6_2yXOQSEqjYMlsPPFJdlZNHcDsccvKJg5O0hNSC2X0mIALzcWFPi5Jym_NSxq2Ch7Pc5aiKiA1U1ZGgUKsMxVZrhtAj8qRwcYoxeV9uucLH3CPxcvQKV37R0PjAtCE-Tq6E68FG-1UTaatZaAoKBmg9MjrOQ7gHWUk77EfgnX2GJA26VUIMeyKZEYpHYv7LFX6x3Fd4JKCHgYnN9qQP8qWwAtvtW3FJq-Bjsu8abE7AGYAzg1N7IZ25ITq-n8_rlzUIE9_75Baxw3_xjjFs84sjUw0RsaKkhfRRLJkxLsMKEA0wSPIRuNfL7MZmSY8emOIg8IAmlgn1iQo_3nfKJEtX6UuojQTDzPYWgfPZIViM2PeLVvgWeTaw2kD52FcqFO5DjMdn-41LRhpysXFXuOHtPN4dr-ViLluwrDemv5T6sxvhyth2zljNalRCI9oWDIqpgXBljv7_mrKOR_lYBdiWe0Ms9EWgqoV0ozLapauMq"
	a, err := UT.Oauth2ClientUsecase.Exchange(biz.Oauth2_AppId_microsoft_vbcapp, code)
	fmt.Println(a)
	fmt.Println(err)
	fmt.Println("a AccessToken: ", a.AccessToken)
	fmt.Println("a RefreshToken: ", a.RefreshToken)
	if a != nil {
		er := UT.Oauth2TokenUsecase.UpdateByToken(biz.Oauth2_AppId_microsoft_vbcapp, a)
		lib.DPrintln(er)
	}
}

func Test_microsoft_vbcapp_userInfo(t *testing.T) {

	getUserEmail("eyJ0eXAiOiJKV1QiLCJub25jZSI6IkxQazV0cHhYWTZ6X0lzZmY5Ni04a0dQOHQ0WnpPWGdOVUUwVVZEVFNrbHMiLCJhbGciOiJSUzI1NiIsIng1dCI6IllUY2VPNUlKeXlxUjZqekRTNWlBYnBlNDJKdyIsImtpZCI6IllUY2VPNUlKeXlxUjZqekRTNWlBYnBlNDJKdyJ9.eyJhdWQiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20iLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9hZTRhN2VjYy1jNzY4LTQzN2UtOTQwOC04OWY4Yjk5Nzg5YWEvIiwiaWF0IjoxNzM4NzU0NTg4LCJuYmYiOjE3Mzg3NTQ1ODgsImV4cCI6MTczODc2MDEyMCwiYWNjdCI6MCwiYWNyIjoiMSIsImFpbyI6IkFZUUFlLzhaQUFBQTNjNjRzY1ZjSTNxOWZjUUlSdjY0ZHRLMzVXc3daOGlSTktvR3pONExuWXdJbzRUeUF2b3FtcVBtSHhTLzYxaDZWN3NIeFNiZDdDWDUrNGlxcWI2Z01tRHFBbHo5ZGJzMitoRnJDQzlCT25ReWQwNGlEY2J3Mk1FZHZZNXVLQk82WjM1anpvc2lvY1JRR2hRNERMeUZiaWNMYUdGeVQzbGZyMFN2Q1JERGxQUT0iLCJhbXIiOlsicHdkIiwibWZhIl0sImFwcF9kaXNwbGF5bmFtZSI6IlZCQyBBcHAiLCJhcHBpZCI6IjVhYTc4MDQ3LWMzNzgtNDk2NS05YjMwLTExNzkwOGQ2NDM5YiIsImFwcGlkYWNyIjoiMSIsImZhbWlseV9uYW1lIjoiTGlhbyIsImdpdmVuX25hbWUiOiJHYXJ5IiwiaWR0eXAiOiJ1c2VyIiwiaXBhZGRyIjoiNDYuMy4yNDAuMTA1IiwibmFtZSI6IkdhcnkgTGlhbyIsIm9pZCI6IjY0YTU1MDdjLTlhNDItNDVlNC1iMjI5LTcyMDk4MzJhYjVmYyIsInBsYXRmIjoiNSIsInB1aWQiOiIxMDAzMjAwNDQzMkM5RDMzIiwicmgiOiIxLkFXRUJ6SDVLcm1qSGZrT1VDSW40dVplSnFnTUFBQUFBQUFBQXdBQUFBQUFBQUFCaUFYWmhBUS4iLCJzY3AiOiJlbWFpbCBwcm9maWxlIFVzZXIuUmVhZCBvcGVuaWQiLCJzaWQiOiIwMDFmMjhhOS1iNzBlLTA0MmUtMjYyNS1mOTVlMzVmZjRlODUiLCJzaWduaW5fc3RhdGUiOlsia21zaSJdLCJzdWIiOiJ5MUxhU3VjN2hteUxOLVY1UzN0ZWlfeFVTN0c5NF9zQUswaUplUzdaLVRFIiwidGVuYW50X3JlZ2lvbl9zY29wZSI6Ik5BIiwidGlkIjoiYWU0YTdlY2MtYzc2OC00MzdlLTk0MDgtODlmOGI5OTc4OWFhIiwidW5pcXVlX25hbWUiOiJnbGlhb0B2ZXRiZW5lZml0c2NlbnRlci5vbm1pY3Jvc29mdC5jb20iLCJ1cG4iOiJnbGlhb0B2ZXRiZW5lZml0c2NlbnRlci5vbm1pY3Jvc29mdC5jb20iLCJ1dGkiOiJZQmhYZE1UNUFVQ1RyazF2czY2MUFBIiwidmVyIjoiMS4wIiwid2lkcyI6WyI2OTA5MTI0Ni0yMGU4LTRhNTYtYWE0ZC0wNjYwNzViMmE3YTgiLCJiNzlmYmY0ZC0zZWY5LTQ2ODktODE0My03NmIxOTRlODU1MDkiXSwieG1zX2lkcmVsIjoiMSA4IiwieG1zX3N0Ijp7InN1YiI6InVibGRUZVc1anFFYWRySWlCNEtVM0tUbVIyaW90Yzg0WEljYXNEc1dYRFUifSwieG1zX3RjZHQiOjE3Mzg2Mjk4MDR9.duCQb8KdTF4lwj2uxKagnG4YhwfPjY7jAEqqPJsD8UyeZm-MMFyZKOy4xDUyM-otns6WSK605gngzCxEGcd3qR-JFczUYU9LGMww2EZ9BL50Y2HnGxfT-6oEnQcfcKE95DV7RJMEn6ftySXDky4HkPe_x2z9zi8mtvDX02H3B0CWS27ifYbrPk9T-VHu3Fg64yG8nex9lGB37vz9VBU05XAzr81_jE--G1wCr7btSgLMs2hxxVtzT2cwo89e4uygV3L0hXuf3hDnh3bM1YbtWjsQ-JqvLomIPVocqxZfmL_e_4oRQLVOjMtufevZNEg8KoG0DGe-5w8SSKI8IHr7-A")
}

type User struct {
	Mail string `json:"mail"`
}

func getUserEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	aaa, err := io.ReadAll(resp.Body)
	fmt.Println("err:", err)
	fmt.Println("aaa:", string(aaa))

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}

	return user.Mail, nil
}

// http://localhost:8050/oauth2/callback?app_id=google&code=4/0AfJohXkXTCYtxgCPqfJxx0bYPw50zkXmDjEAg6AAEVuD1Fude6EN0kUISKBtQRfvPOIf_A&scope=https://www.googleapis.com/auth/drive.readonly%20https://www.googleapis.com/auth/drive%20https://www.googleapis.com/auth/drive.file%20https://www.googleapis.com/auth/spreadsheets.readonly%20https://www.googleapis.com/auth/spreadsheets
func Test_google_Exchange(t *testing.T) {
	code := "4/0AfJohXkm6uLJtLRlaZ6LPYd7RBOmUh3Naz-AXyipNCx8w1cvhP3rImF0XfJfPGOGM09zog"
	code = "4/0AfJohXku4IUImg1eP2t4UdHTj90zxMzWOJWVtpbn3RHTJXlBr2oyMT2doa0CSQC3_LSakQ"
	code = "4/0AfJohXkXTCYtxgCPqfJxx0bYPw50zkXmDjEAg6AAEVuD1Fude6EN0kUISKBtQRfvPOIf_A"
	code = "4/0AfJohXnzzBY00uIJcUqfRMrFf2KAUnGjROCvMZTC-rycwqYXSqjHNHp9v50WmK8qZRDzLg"
	code = "4/0AfJohXlMlU_lfkFQPSImd8duxWg7-RjD8ViwsGpZVVWgQHblLEV74Zho2Gp7LNA2mzjI6A"
	a, err := UT.Oauth2ClientUsecase.Exchange(biz.Oauth2_AppId_google, code)
	fmt.Println(a)
	fmt.Println(err)
	if a != nil {
		er := UT.Oauth2TokenUsecase.UpdateByToken(biz.Oauth2_AppId_google, a)
		lib.DPrintln(er)
	}
}

func Test_Oauth2ClientUsecase_google_RefreshAccessToken(t *testing.T) {
	client, err := UT.Oauth2ClientUsecase.GetByAppId(biz.Oauth2_AppId_google)
	lib.DPrintln(err)
	token, err := UT.Oauth2TokenUsecase.GetByClientId(client.ClientId)
	lib.DPrintln(token)
	err = UT.Oauth2TokenUsecase.RefreshAccessToken(client, token)
	lib.DPrintln(err)
}

func Test_Oauth2ClientUsecase_adobesign_AuthUrl(t *testing.T) {
	a, _ := UT.Oauth2ClientUsecase.AuthUrl(biz.Oauth2_AppId_adobesign)
	fmt.Println(a)
}

func Test_adobesign_Exchange(t *testing.T) {
	code := "CBNCKBAAHBCAABAA0ejCzGwuOrBgJeUYgmY4wJnf9cuVoy8B"
	a, err := UT.Oauth2ClientUsecase.Exchange(biz.Oauth2_AppId_adobesign, code)
	fmt.Println(a)
	fmt.Println(err)
	if a != nil {
		er := UT.Oauth2TokenUsecase.UpdateByToken(biz.Oauth2_AppId_adobesign, a)
		lib.DPrintln(er)
	}
}

func Test_Oauth2ClientUsecase_adobesign_RefreshAccessToken(t *testing.T) {
	client, err := UT.Oauth2ClientUsecase.GetByAppId(biz.Oauth2_AppId_adobesign)
	lib.DPrintln(err)
	token, err := UT.Oauth2TokenUsecase.GetByClientId(client.ClientId)
	lib.DPrintln(token)
	err = UT.Oauth2TokenUsecase.RefreshAccessToken(client, token)
	lib.DPrintln(err)
}

func Test_Oauth2ClientUsecase_box_AuthUrl(t *testing.T) {
	a, _ := UT.Oauth2ClientUsecase.AuthUrl(biz.Oauth2_AppId_box)
	fmt.Println(a)
}

func Test_box_Exchange(t *testing.T) {
	code := "DVZduYu1WEZKen8ct8VEUDJe3WdgFGTk"
	a, err := UT.Oauth2ClientUsecase.Exchange(biz.Oauth2_AppId_box, code)
	fmt.Println(a)
	fmt.Println(err)
	if a != nil {
		er := UT.Oauth2TokenUsecase.UpdateByToken(biz.Oauth2_AppId_box, a)
		lib.DPrintln(er)
	}
}

// https://app.asana.com/-/oauth_authorize?response_type=code&client_id=1206234321575740&redirect_uri=http%3A%2F%2Flocalhost%3A8050%2Foauth2%2Fcallback%3Fapp_id%3Dasana&state=<STATE_PARAM>
func Test_Oauth2ClientUsecase_asana_AuthUrl(t *testing.T) {
	a, _ := UT.Oauth2ClientUsecase.AuthUrl(biz.Oauth2_AppId_asana)
	fmt.Println(a)
}

func Test_Oauth2ClientUsecase_AuthUrl(t *testing.T) {
	a, _ := UT.Oauth2ClientUsecase.AuthUrl(biz.Oauth2_AppId_docusign)
	fmt.Println(a)
}

func Test_Oauth2ClientUsecase_Exchange(t *testing.T) {
	code := "eyJ0eXAiOiJNVCIsImFsZyI6IlJTMjU2Iiwia2lkIjoiNjgxODVmZjEtNGU1MS00Y2U5LWFmMWMtNjg5ODEyMjAzMzE3In0.AQ0AAAABAAYABwCAXSerGAzcSAgAgOmt8hgM3EgCAPOiKEGmO3tCo6-5_Xg_y2kVAAEAAAAYAAQAAAAFAAAACgAAAB0AAAACAAAADQAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyIgAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyEQABMACAXSerGAzcSBIAAQAAAAsAAABpbnRlcmFjdGl2ZTcAMgn50XP8hUCOHLtxH44mWx8AAQAAAD0AAAA0MTI4YTJmMy0zYmE2LTQyN2ItYTNhZi1iOWZkNzgzZmNiNjk7MjAyNC0wMS0wMyAwNDo1ODo0N1o7Mzsy.gXMXKFnJl5Sk6S2FJdt7Jk5NmVVT69lJf849ErE85mqcyqW_bvacG6Z0V5eAxdUU5ESwnhZ47BRm1Z4RZOkPR5toWETR1RJNQ797qBnKrZ0J1wuKZH1aE2ceJF664b5C1jOZO67zlv_Ec_jcajwwjdjprtrqNs5wDBH4aVeZoNfNx6mUA8RmHpKjXRFr2bnYGCuD0qWQPN1DqTAE6JcHFLlkTxuDjiiCWive0G-66eosrISOhGV1FUEn2mn2atZ4ZInHxqXG0c8RfcvZeJwvhR8A1LwnuSRYE5Q64Pt3Ocn0E2xxXAS4c-nmgahIXjGONaUmsPCa_XMM_temZnVocA"
	code = "eyJ0eXAiOiJNVCIsImFsZyI6IlJTMjU2Iiwia2lkIjoiNjgxODVmZjEtNGU1MS00Y2U5LWFmMWMtNjg5ODEyMjAzMzE3In0.AQwAAAABAAYABwAAlZ9FSQzcSAgAACEmjUkM3EgCAPOiKEGmO3tCo6-5_Xg_y2kVAAEAAAAYAAEAAAAFAAAADQAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyIgAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyMAAAlZ9FSQzcSBIAAQAAAAsAAABpbnRlcmFjdGl2ZTcAMgn50XP8hUCOHLtxH44mWx8AAQAAAD0AAAA0MTI4YTJmMy0zYmE2LTQyN2ItYTNhZi1iOWZkNzgzZmNiNjk7MjAyNC0wMS0wMyAxMDo0Njo0Mlo7Mzsy.OWIzU5StvC9pEGfT0Qvbtz19zSL7_KfLErl2I7iQIhgvqBvyCwF-5ROc7eDPp0vO4vAwOSIMO-bTfEGEoPxpZ4r8w0JHd8HB1x9-Rq48mjaaArALTD5VBrI-9yYpsVEIRaJwy3qnJeYWc0wMOccbFGdNPXwZM3JIL68o-vSV-bQU7fGJyX0DQAS17ozsmhl6KnWoxxPFOb447DBTzj0A9HryIYyf4LOG0xMqe_C5wPf_hD7fuuyxmUi56uJucOPIFO5ovc3oZ87h5gczCSS-yXedfFA1lhll294kFrq5y6BMOcHS5BPinidVxQdCQAIHapoEgM6wArvTYdJqoFiEFw"
	code = "eyJ0eXAiOiJNVCIsImFsZyI6IlJTMjU2Iiwia2lkIjoiNjgxODVmZjEtNGU1MS00Y2U5LWFmMWMtNjg5ODEyMjAzMzE3In0.AQoAAAABAAYABwAABw_XSQzcSAgAAJOVHkoM3EgCAPOiKEGmO3tCo6-5_Xg_y2kVAAEAAAAYAAEAAAAFAAAADQAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyIgAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyNwAyCfnRc_yFQI4cu3EfjiZbMAAAlZ9FSQzcSA.fm8B-nL3K7F79mW3m_ai5nqFGYLLlfSpX1NlN6-r4oyI6_4yhFWwc7vNAq9MH49DUTrpYAlOEIJUAHjSrAnqLeH4hfCXtdi20fTsrJ59WuzO6n-KkDsuIQu4z3TcPKB7BUqhHxpfLFGd4Bywp27eoxiMp1XstAfd_R-9RQdo0J5x8k6Gz6RAecL6DeXcayPVufD3deWNDfrr6zGCvHugF5bn_pfKask4u1FHdRorCtzChxaJzDFbjPTwseMDKLiEOkV8EWsLNB3UoMazTaCMU-o0S_U3VJEi6HorhpQs4XAiu-5ddpgg8IHt1B1CCe97XIUOwFzwu8lG5jgez3FS-w"
	code = "eyJ0eXAiOiJNVCIsImFsZyI6IlJTMjU2Iiwia2lkIjoiNjgxODVmZjEtNGU1MS00Y2U5LWFmMWMtNjg5ODEyMjAzMzE3In0.AQoAAAABAAYABwCAZy50awzcSAgAgPO0u2sM3EgCAPOiKEGmO3tCo6-5_Xg_y2kVAAEAAAAYAAEAAAAFAAAADQAkAAAANDNmZjNiZmMtNGEzMS00YTViLTkxNjQtNDM2NjU3MmVlZGZkIgAkAAAANDNmZjNiZmMtNGEzMS00YTViLTkxNjQtNDM2NjU3MmVlZGZkNwAyCfnRc_yFQI4cu3EfjiZbMAAAPFTVZwzcSA.42gCVgn6uaovYZT6c6WkatRWbJrORX42U-ayLYs_nw6xeo0ObRvqKBL3vzqJ2mIwj36pa2WrD8S2m7coitBTm-5SLaUwkt4hJfJFCC5KGcEeEb8OfS9yhyiimg9rdnTwHBirk3R_6P5aflVkmQRFz-uQ4aoXt1b7OYAQPCuedwceHnmE_j3qLbZENtY2sABYwiJ8C0l5XH7Qmtqto67FVPwuRdG-gwvUiviyJsPSuiAuXsx-Ft7ymTr8Xw0UF_rq0KWw9h0Tmy72SrKiky7CLRdECc-F8iD1K2Znr1o5LEmU-wfz6QhpUM8DJuDaA2lW9z4nZ_r1YlpTKNkKHl6qUA"
	code = "eyJ0eXAiOiJNVCIsImFsZyI6IlJTMjU2Iiwia2lkIjoiNjgxODVmZjEtNGU1MS00Y2U5LWFmMWMtNjg5ODEyMjAzMzE3In0.AQoAAAABAAYABwAAWmDPawzcSAgAAObmFmwM3EgCAPOiKEGmO3tCo6-5_Xg_y2kVAAEAAAAYAAEAAAAFAAAADQAkAAAANDNmZjNiZmMtNGEzMS00YTViLTkxNjQtNDM2NjU3MmVlZGZkIgAkAAAANDNmZjNiZmMtNGEzMS00YTViLTkxNjQtNDM2NjU3MmVlZGZkNwAyCfnRc_yFQI4cu3EfjiZbMAAAPFTVZwzcSA.peJtKRZbGclLtQ4sXxtH5LjJdNd9HOHGrIj7XezVOmz4frGeyKibP_A8yWy7pKmsfv3Igfziahxt-7lwdbSPi1EQcJx8kKSSkJwFy6UZaP7F3F9b34qeVwaSaaR1MOWEdTddvfVYePcWbEvyHtkd8YxvhcvzR_AlnaBRNogyZmt9AAuEjTuL1xUKjlpTJ9eBNV7tWk3MQxyUKZe4vbd986WTCd3PWnqmmGemxkUa6YnKqwGX86csmnMAuJWTRUO_BuZnwd5ZLAIf93gYVmBdW8fyu8pJtmvL9VGfiLDbLTaVs6IBIlpsMfjmSN57w68Ylkce8_9il6TK5CTNfPAbow"
	code = "eyJ0eXAiOiJNVCIsImFsZyI6IlJTMjU2Iiwia2lkIjoiNjgxODVmZjEtNGU1MS00Y2U5LWFmMWMtNjg5ODEyMjAzMzE3In0.AQwAAAABAAYABwCAn7vk0AzcSAgAgCtCLNEM3EgCAPOiKEGmO3tCo6-5_Xg_y2kVAAEAAAAYAAEAAAAFAAAADQAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyIgAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyMACAn7vk0AzcSBIAAQAAAAsAAABpbnRlcmFjdGl2ZTcAMgn50XP8hUCOHLtxH44mWx8AAQAAAD0AAAA0MTI4YTJmMy0zYmE2LTQyN2ItYTNhZi1iOWZkNzgzZmNiNjk7MjAyNC0wMS0wNCAwMjo1NzozMVo7Mzsy.wlnQjoSk6RAs3HCfS-rur5f4PJQfYKln4IyCJfN1bpZZOnr7c-D6nbe4u5LZvZvOxA8L94E2DJjdchrmwfwcvQusDMQhzkqXjrZrMY28UIrGjAcmwJF-4LLY8tX7uyEdGn4c2PjlkIoumF3CT51fMQjnVcPfo0WO8M9RbIKrywjVyFgHNXfztGrfuM4C4aXtyJWpn_T2ijf1Vlq3S15MCNId6IhMBhAPxbKuSrMQu1A9Fk-T0f6eNjNZcGzOwbcfrW9BoVkN-ez7Oc3vDEEH5eR1ba6ZJFHN8T57rwLSNSdyllPhYAy3HDTO5kWwDEOcP8MjZEg7MZPYF-fvU5J3aA"
	code = "eyJ0eXAiOiJNVCIsImFsZyI6IlJTMjU2Iiwia2lkIjoiNjgxODVmZjEtNGU1MS00Y2U5LWFmMWMtNjg5ODEyMjAzMzE3In0.AQoAAAABAAYABwAA4YX-uw3cSAgAAG0MRrwN3EgCAPOiKEGmO3tCo6-5_Xg_y2kVAAEAAAAYAAEAAAAFAAAADQAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyIgAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyNwAyCfnRc_yFQI4cu3EfjiZbMACANBfHuw3cSA.nXwZ5mHBIvH3ocOy6AZ77MJNxBMgcd2Ris6NszfNB258cP39lV1AjiTDwRVMfuDisaP201EDJ_vfwpA3b5akzKz2iVqJeYqZTRYbBd5ugL7LyyBWDZP8-MqAdJ4LoQBbbIcJAN_CqaQP0MP8Vju24czfWLLYxhvtuVB7Enw1HfJEFqOHtlzKRYPghPXkfWfT_E1it7rMjIojSote28BOR0XYii4Qh_pySF1EtwiO4FBIT2ahsJ-nMMLSJUSjTj-4tT-kdI3soTDcJXNPqCxIDYearn-sNqAT22LkaO__2SgqoXnOridfrveoihIbr952L7VW3JGFw-DUvlaUHl0NoQ"
	code = "eyJ0eXAiOiJNVCIsImFsZyI6IlJTMjU2Iiwia2lkIjoiNjgxODVmZjEtNGU1MS00Y2U5LWFmMWMtNjg5ODEyMjAzMzE3In0.AQoAAAABAAYABwAADfwG-A3cSAgAAJmCTvgN3EgCAPOiKEGmO3tCo6-5_Xg_y2kVAAEAAAAYAAEAAAAFAAAADQAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyIgAkAAAAMDZmMmVhYjQtMjk1Yi00N2QxLTkxZjItZDZmNTVhNzhkZWIyNwAyCfnRc_yFQI4cu3EfjiZbMACAEOQj9w3cSA.H43kDtornCLdBd7Tw2GJ2LfqYEopElaWAs3Bn8s8my3Vokh8j8LfZnomE0kA6pgHx1_-2LEjuuMl3ziVBi6XEoMS1ScyMd2If7jS43TW7hBsTyEmV7_Pk3OiqPJhqKKWtiR_lwZB8Qxdk53vq6nAGuaPuZRThVTyuxvlMGGBHBjW6gmwvmwbJ8fUoRedSamTeFB_cJr1Zu5G-iF0p6GHiN9GHWThSOqsQCRZ0QA1qKjfYfat3DHXX9LUVsNriRu3FAb3kbtK32ekxkCSAW3vVW9biSfoF6VVfh4jUqRfreY0jjzuF0CrvDI94QodTasZJaQ0oqfbNj1DLC6Avg3WLQ"
	code = "eyJ0eXAiOiJNVCIsImFsZyI6IlJTMjU2Iiwia2lkIjoiNjgxODVmZjEtNGU1MS00Y2U5LWFmMWMtNjg5ODEyMjAzMzE3In0.AQwAAAABAAYABwCAjOYEkg7cSAgAgBhtTJIO3EgCAISVWpu26mtFsI9QA8_F86gVAAEAAAAYAAEAAAAFAAAADQAkAAAAYTZlZGMwOTktMjJkNi00OGY3LWJmMTktZGE4MGUxNzg5YjRiIgAkAAAAYTZlZGMwOTktMjJkNi00OGY3LWJmMTktZGE4MGUxNzg5YjRiMACAjOYEkg7cSBIAAQAAAAsAAABpbnRlcmFjdGl2ZTcARTbDZhFZuUKuXYUTuBbMPh8AAgAAAD0AAAA0MTI4YTJmMy0zYmE2LTQyN2ItYTNhZi1iOWZkNzgzZmNiNjk7MjAyNC0wMS0wNCAxNDozNDozNVo7MzsyPQAAADliNWE5NTg0LWVhYjYtNDU2Yi1iMDhmLTUwMDNjZmM1ZjNhODsyMDI0LTAxLTA2IDA4OjMyOjI5WjszOzI.qqs42oeEahPtKmgCfdGteQ0hNYKsl6dnnZqiE0w7mIOLRi61u3VovvTH6B7Q7in211C4Mp56aZhJw0AWtXsalf8fvb-W0YRspxqMFChYchg1PSQ3rpXK0nPO8KQ6p1KYNcVz3QrppeEj4KLJThnEkCMBI7kd68ImcxwEq_HxECKWgqt67-uNkNd0GjJt3ivwXT8cxrELiKcb3oUIzIWfYPuawHBA7_nrL_dujsiJDkegNAGFf40Hq2jilxFFU3oN9xhJ94H0RG3KiftoOiiJwBMwXJo2idqVZIJZ0ND-iZVZTJZXFCCHGNANtav4oBnmi2GDq9nNB5xNgncB1kOlDw"
	a, err := UT.Oauth2ClientUsecase.Exchange(biz.Oauth2_AppId_docusign, code)
	fmt.Println(a)
	fmt.Println(err)
	if a != nil {
		er := UT.Oauth2TokenUsecase.UpdateByToken(biz.Oauth2_AppId_docusign, a)
		lib.DPrintln(er)
	}
}

func Test_Oauth2ClientUsecase_RefreshAccessToken(t *testing.T) {
	client, err := UT.Oauth2ClientUsecase.GetByAppId(biz.Oauth2_AppId_zoom)
	lib.DPrintln(err)
	token, err := UT.Oauth2TokenUsecase.GetByClientId(client.ClientId)
	lib.DPrintln(token)
	err = UT.Oauth2TokenUsecase.RefreshAccessToken(client, token)
	lib.DPrintln(err)
}

func Test_Oauth2TokenUsecase_WaitingRefreshToken(t *testing.T) {
	a, err := UT.Oauth2TokenUsecase.WaitingRefreshToken()
	lib.DPrintln(a, err)
}
