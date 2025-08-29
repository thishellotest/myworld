package lib

import (
	"strings"
	"testing"
)

func Test_MD5Hash(t *testing.T) {
	a := MD5Hash("aaa")
	DPrintln(a)
	a = MD5Hash("a\naa")
	DPrintln(a)
	a = MD5Hash("a\naa")
	DPrintln(a)
	a = MD5Hash("a\naa")
	DPrintln(a)
}

func Test_VerifyEmail(t *testing.T) {
	a := VerifyEmail("lialing@foxmail.com")
	DPrintln(a)
}

func Test_FormatNameWithVersion(t *testing.T) {
	aa := FormatNameWithVersion("name", "v2024-06-16_11-33")
	DPrintln(aa)

	aa = FormatNameWithVersion("ccc我是中文.pdf.pdf.doc", "v2024-06-16_11-33")
	DPrintln(aa)
}

func Test_StringToBytesFormat(t *testing.T) {
	a := StringToBytesFormat("我是")
	DPrintln(a)
}

func Test_InterfaceToString(t *testing.T) {
	var aaa *int32
	cc := InterfaceToString(aaa)
	DPrintln(cc)

}

func Test_IsValidEmail(t *testing.T) {
	a := IsValidEmail("a@qq.com.com")
	DPrintln(a)
	a = IsValidEmail("sss.@qq.com.com")
	DPrintln(a)
	a = IsValidEmail("a.s%+ss.a-c@csdfdsfsafd.qqsss")
	DPrintln(a)
}

func Test_IsValidURL(t *testing.T) {
	a := IsValidURL("http://www.baidu.com")
	DPrintln(a)
	a = IsValidURL("https:///工是")
	DPrintln(a)
}

func Test_GenerateSafePassword(t *testing.T) {
	a := GenerateSafePassword(8)
	DPrintln(a)
}

func Test_Capitalize(t *testing.T) {
	abc := Capitalize("aaabc")
	DPrintln(abc)

	firstName := "aInfo"
	firstName = Capitalize(strings.ToLower(firstName))
	DPrintln(firstName)
}

func Test_TruncateStringWithRune(t *testing.T) {
	a := TruncateStringWithRune("123456", 5)
	DPrintln(a)
}

func Test_AESEncrypt(t *testing.T) {
	a, b, err := AESEncrypt([]byte("fdsafdsafsdafcjsaojfds"), []byte("Eiy8yahcsohd3eeV"))
	DPrintln(err)
	DPrintln(string(a), string(b))
}

func Test_EncryptToBase64(t *testing.T) {
	a, err := EncryptToBase64([]byte("[OK]12😉😉😏我是。.&^%%2)((2!~/a/<<>"), []byte("Eiy8yahcsohd3eeV"))
	DPrintln(err)
	DPrintln(a)
	a, err = EncryptToBase64([]byte(""), []byte("Eiy8yahcsohd3eeV"))
	DPrintln(err)
	DPrintln(a)
}

func Test_DecryptFromBase64(t *testing.T) {
	// ad4c1dgZDf0836CVenpUrMEERoq9Xg==
	a, err := DecryptFromBase64("VPUk6Y9PjgPCJTRrhJDOQDjkymRIy77HjhTBsxRmcAXIy5K0BEMBecD52io/", []byte("Eiy8yahcsohd3eeV"))
	DPrintln(err)
	DPrintln(a)
}
