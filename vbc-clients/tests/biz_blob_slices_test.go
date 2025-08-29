package tests

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_BlobSliceUsecase_HasFinish(t *testing.T) {
	aaa, er := UT.BlobSliceUsecase.HasFinish("aa")
	lib.DPrintln(aaa, er)
}

func Test_BlobSliceUsecase_HasFinish1(t *testing.T) {
	aaa, er := UT.BlobSliceUsecase.HasFinish("5e1140c7fd0f48948088529b7539065b")
	lib.DPrintln(aaa, er)
}

func Test_BlobSliceUsecase_GetOcrResultByGid(t *testing.T) {
	a, er := UT.BlobSliceUsecase.GetOcrResultByGid(context.TODO(), "fb73178ae71f46b486494a227c06a8f2")
	lib.DPrintln(er)
	content := a.GetContent()
	lib.DPrintln(content)
}

func Test_BlobSliceUsecase_Prompt(t *testing.T) {

	res, _ := UT.BlobSliceUsecase.AllByCond(Eq{"blob_gid": "e1723f0a921341fc8e2acf274d007732"})
	content := ""
	for _, v := range res {
		content += fmt.Sprintf("Page %s:\n%s\n______\n\n", v.SliceId, v.OcrResultContent)
	}
	aa, er := os.OpenFile("/tmp/prompt_3.log", os.O_CREATE|os.O_WRONLY, 0644)
	if er != nil {
		panic(er)
	}
	io.WriteString(aa, content)
	defer aa.Close()
}

func Test_BlobSliceUsecase_GetBlobNameUrlVo(t *testing.T) {
	a, _ := UT.BlobSliceUsecase.GetByCond(Eq{"id": 1238})
	b1, b2, _ := a.GetBlobNameUrlVo(UT.CommonUsecase, UT.CommonUsecase.Log, UT.AzstorageUsecase)
	lib.DPrintln(b1.Url)
	lib.DPrintln(b2.Url)
}
