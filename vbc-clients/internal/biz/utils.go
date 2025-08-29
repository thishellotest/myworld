package biz

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"time"
	"vbc/configs"
	"vbc/lib"
)

// TimeToVBCDisplay 2024-03-29T23:06:04+08:00
func TimeToVBCDisplay(timeStr string) (string, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return "", err
	}
	cc := t.In(configs.GetVBCDefaultLocation()).Format("January 2, 2006 03:04 PM")
	return cc, nil
}

func TesReadExcel() {
	f, err := excelize.OpenFile("./VBC Conditions.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get value from cell by given worksheet name and cell reference.
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "====")
		}
		fmt.Println()
	}
}

func ToImages(fileNamePath string, destDir string) (records []string, err error) {
	return nil, nil
	//doc, err := fitz.New(fileNamePath)
	//if err != nil {
	//	return nil, err
	//}
	//defer doc.Close()
	//
	//imgUid := uuid.UuidWithoutStrike()
	//// Extract pages as images
	//for n := 0; n < doc.NumPage(); n++ {
	//	img, err := doc.Image(n)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	str := filepath.Join(destDir, fmt.Sprintf(imgUid+"%03d.jpg", n))
	//	f, err := os.Create(str)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	err = jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
	//	if err != nil {
	//		return nil, err
	//	}
	//	f.Close()
	//	records = append(records, str)
	//}
	//return records, nil
}

func EncryptSensitive(value string) (string, error) {
	key := configs.EnvSensitiveDataKey()
	if key == "" {
		return "", errors.New("SensitiveDataKey is empty")
	}
	return lib.EncryptToBase64([]byte(value), []byte(key))
}

func DecryptSensitive(value string) (string, error) {
	key := configs.EnvSensitiveDataKey()
	if key == "" {
		return "", errors.New("SensitiveDataKey is empty")
	}
	bytes, err := lib.DecryptFromBase64(value, []byte(key))
	return string(bytes), err
}
