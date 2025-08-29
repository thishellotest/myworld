package tests

import (
	"context"
	"os"
	"testing"
	"vbc/lib"
)

func Test_GoogleDriveUsecase_ListItemsInFolder(t *testing.T) {

	// 1xwolHdQbmw-_uVz3Xxi-AaBlUoJpxgGR : prod folder
	// 1KZo9l1sOENV4s8DJbL2vIx8FBiBIznww : test folder
	r, err := UT.GoogleDriveUsecase.ListItemsInFolder(context.TODO(), "1KZo9l1sOENV4s8DJbL2vIx8FBiBIznww")
	lib.DPrintln("err:", err)
	lib.DPrintln(r, err)
}

// 示例文件：https://veteranbenefitscenter.app.box.com/folder/263409020836

func Test_GoogleDriveUsecase_CreateFolder(t *testing.T) {
	aaa, err := UT.GoogleDriveUsecase.CreateFolder(context.TODO(), "a", "abc")
	lib.DPrintln(aaa, err)
	lib.DPrintln("err:", err)
}

func Test_GoogleDriveUsecase_UploadFile(t *testing.T) {

	file, err := os.Open("./res/a.pdf")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	aaa, err := UT.GoogleDriveUsecase.UploadFile(context.TODO(), UT.Conf.GoogleDrive.PaymentsFolderId, "cc.pdf", file)

	lib.DPrintln(aaa, err)
	lib.DPrintln("err:", err)
}
