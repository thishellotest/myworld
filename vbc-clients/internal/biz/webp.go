package biz

import (
	"github.com/chai2010/webp"
	"github.com/go-kratos/kratos/v2/log"
	"image/jpeg"
	"os"
	"vbc/internal/conf"
)

type WebpUsecase struct {
	log             *log.Helper
	conf            *conf.Data
	CommonUsecase   *CommonUsecase
	ResourceUsecase *ResourceUsecase
}

func NewWebpUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	ResourceUsecase *ResourceUsecase,
) *WebpUsecase {
	uc := &WebpUsecase{
		log:             log.NewHelper(logger),
		CommonUsecase:   CommonUsecase,
		conf:            conf,
		ResourceUsecase: ResourceUsecase,
	}

	return uc
}

func (c *WebpUsecase) JpgToWebp(jpgFilePath string, webpFilePath string) (width int, height int, err error) {

	file, err := os.Open(jpgFilePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	// 解码 JPG 图像
	img, err := jpeg.Decode(file)
	if err != nil {
		return 0, 0, err
	}
	// 获取图像尺寸
	bounds := img.Bounds()
	width = bounds.Dx()  // 宽度
	height = bounds.Dy() // 高度
	// 创建输出 WebP 文件
	outFile, err := os.Create(webpFilePath)
	if err != nil {
		return 0, 0, err
	}
	defer outFile.Close()
	// 编码为 WebP 格式，质量可设置为 0~100
	err = webp.Encode(outFile, img, &webp.Options{Quality: 80})
	if err != nil {
		return 0, 0, err
	}
	return
}

func (c *WebpUsecase) TestWebp() error {
	jpgFile := c.ResourceUsecase.ResPath() + "/123.jpg"
	// 打开 JPG 文件
	file, err := os.Open(jpgFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 解码 JPG 图像
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	// 创建输出 WebP 文件
	outFile, err := os.Create("output.webp")
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 编码为 WebP 格式，质量可设置为 0~100
	err = webp.Encode(outFile, img, &webp.Options{Quality: 80})
	if err != nil {
		return err
	}
	return nil
}
