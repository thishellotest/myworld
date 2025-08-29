package images

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"testing"
	"vbc/lib"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func getImageDimensions(filePath string) (int, int, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}

func Test122(t *testing.T) {
	filePath := "../tmp/000.jpg"
	//lib.GetImageDimensions(filePath)
	width, height, err := lib.GetImageDimensions(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Image dimensions: %dx%d\n", width, height)
}

func TestCalDimensions(t *testing.T) {

	aa, ss := CalDimensions(32323, 5553, 300, 500)
	fmt.Println(aa, ss)
	aa, ss = CalDimensions(10, 50, 300, 500)
	fmt.Println(aa, ss)
	aa, ss = CalDimensions(50, 10, 300, 500)
	fmt.Println(aa, ss)
}

func CalDimensions(sourceW, sourceH, MaxDestW, MaxDestH int) (int, int) {
	a1 := float64(sourceW) / float64(sourceH)
	a2 := float64(MaxDestW) / float64(MaxDestH)
	if a1 > a2 {
		aaa := (float64(MaxDestW) / float64(sourceW)) * float64(sourceH)
		return MaxDestW, int(math.Floor(aaa))
	} else {
		aaa := (float64(MaxDestH) / float64(sourceH)) * float64(sourceW)
		return int(math.Floor(aaa)), MaxDestH
	}
}
