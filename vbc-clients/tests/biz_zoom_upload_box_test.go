package tests

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"vbc/lib"
)

func Test_UploadFileToBoxUseChunks(t *testing.T) {
	filePath := "/tmp/VSCode-darwin-arm64.zip"
	fileName := "VSCode-darwin-arm64.zip"
	boxFolderId := "267683828516"
	err := UT.BoxUsecase.UploadFileToBoxUseChunks(filePath, fileName, boxFolderId)
	lib.DPrintln(err)
}

func Test_CommitUpload(t *testing.T) {
	//UT.BoxUsecase.CommitUpload("https://upload.app.box.com/api/2.0/files/upload_sessions/5B4A8BD7E088EEEBB3231E2A3292C2B2/commit")
}

func Test_ZoomUploadBoxUsecase_UploadToBox(t *testing.T) {
	url := "https://us06web.zoom.us/rec/download/Au386IkDQnYnDbYvq5dXsw1cvyNHH2CdgWBtB2eHatpEfvc8mkYwOri4BxP2KCMqNzVZoDmeUP7ZxIM.mumk6MD1Lup9BT4F"
	err := UT.ZoomUploadBoxUsecase.UploadToBox(url,
		"267683828516",
		"GMT20240813-203606_Recording_1920x1080.mp4",
		376674194,
	)
	lib.DPrintln(err)
}

func Test_aaaccc(t *testing.T) {
	// Defining reader using NewReader method
	reader := strings.NewReader("Geeks")

	// Defining buffer of specified length
	// using make keyword
	buffer := make([]byte, 4)

	for {
		// Calling ReadFull method with its parameters
		n, err := io.ReadFull(reader, buffer)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			panic(err)
		}

		// Prints output
		fmt.Printf("Number of bytes in the buffer: %d\n", n)
		fmt.Printf("Content in buffer: %s\n", buffer[:n])
		if err == io.EOF {
			break
		}
	}
	//此时只剩s没有读取，即大小为1，当设置大于1时报错：panic: unexpected EOF
	//buffer2 := make([]byte, 1)
	//n, err = io.ReadFull(reader, buffer2)
	//// If error is not nil then panics
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Prints output
	//fmt.Printf("Number of bytes in the buffer: %d\n", n)
	//fmt.Printf("Content in buffer: %s\n", buffer2)

}
