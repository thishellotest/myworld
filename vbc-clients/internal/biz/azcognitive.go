package biz

import (
	"context"
	"encoding/base64"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"io"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
)

type AzcognitiveUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	LogUsecase    *LogUsecase
}

func NewAzcognitiveUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	LogUsecase *LogUsecase) *AzcognitiveUsecase {
	uc := &AzcognitiveUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		LogUsecase:    LogUsecase,
	}
	return uc
}

func AzcognitiveKey() string {
	return configs.EnvAzureCognitiveKey()
}

func AzcognitiveUrl() string {
	if configs.IsProd() {
		return "https://documentintelligenceeu2s0.cognitiveservices.azure.com"
	}
	return "https://documentintelligenceeu2s0.cognitiveservices.azure.com"
	return "https://documentintelligenceeus2.cognitiveservices.azure.com"
}

func (c *AzcognitiveUsecase) PrebuiltRead(reader io.Reader) (operationLocation string, err error) {

	// https://learn.microsoft.com/zh-cn/azure/ai-services/computer-vision/how-to/call-analyze-image-40?pivots=programming-language-rest-api
	// 映像分析模型自定义：推理: $2/1K 个事务 自定义图像分类(预览)
	//  定价：https://azure.microsoft.com/zh-cn/pricing/details/cognitive-services/computer-vision/
	// 5次：请求一样数据 2024-04-26，验证花费 0.08元 computervision/imageanalysis:analyze  1000/16元
	//{your-resource-endpoint}.cognitiveservices.azure.com/formrecognizer/documentModels/prebuilt-layout:analyze?api-version=2023-07-31&features=ocrHighResolution
	api := AzcognitiveUrl() + "/formrecognizer/documentModels/prebuilt-read:analyze?api-version=2023-07-31"
	params := make(lib.TypeMap)

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)
	params.Set("base64Source", encoded)
	resp, err := lib.RequestDoTimeout("POST", api, params.ToBytes(), map[string]string{
		"Ocp-Apim-Subscription-Key": AzcognitiveKey(),
	}, time.Minute*5)

	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("resp is nil")
	}
	if resp.StatusCode != 202 {
		return "", errors.New("resp.StatusCode: " + InterfaceToString(resp.StatusCode))
	}
	apimRequestId := resp.Header.Get("Apim-Request-Id")
	operationLocation = resp.Header.Get("Operation-Location")

	c.log.Debug("apimRequestId: ", apimRequestId)
	c.log.Debug(InterfaceToString(resp.Header))
	c.log.Debug(InterfaceToString(resp.StatusCode))

	er := c.LogUsecase.SaveLog(0, "PrebuiltRead", map[string]interface{}{
		"operationLocation": operationLocation,
	})
	if er != nil {
		c.log.Error(er)
	}
	/*
		resp.Header: {"Apim-Request-Id":["4f9fe37a-a1df-4703-a8c8-113e3bf33930"],"Content-Length":["0"],"Date":["Tue, 02 Jul 2024 07:43:40 GMT"],"Operation-Location":["https://documentintelligenceeu2s0.cognitiveservices.azure.com/formrecognizer/documentModels/prebuilt-read/analyzeResults/4f9fe37a-a1df-4703-a8c8-113e3bf33930?api-version=2023-07-31"],"Strict-Transport-Security":["max-age=31536000; includeSubDomains; preload"],"X-Content-Type-Options":["nosniff"],"X-Envoy-Upstream-Service-Time":["173"],"X-Ms-Region":["East US 2"]}
		resp Code: 202
		body:  err:
	*/
	return operationLocation, nil
}

func (c *AzcognitiveUsecase) GetPrebuiltReadResultWithBlock(ctx context.Context, operationLocation string) (res lib.TypeMap, err error) {

	ticker := time.NewTimer(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-ticker.C:
			return nil, errors.New(operationLocation + ": GetPrebuiltReadResultBlock: Timeout")
		default:
			res, _, err := lib.HTTPJsonWithHeaders("GET", operationLocation, nil, map[string]string{
				"Ocp-Apim-Subscription-Key": AzcognitiveKey(),
			})
			if err != nil {
				return nil, err
			}
			if res == nil {
				return nil, errors.New("res is nil")
			} else {
			}
			typeMap := lib.ToTypeMapByString(*res)
			if typeMap.GetString("status") == "succeeded" {
				return typeMap, nil
			} else if typeMap.GetString("status") == "running" {
				c.log.Debug("res: ", InterfaceToString(res))
			} else {
				return nil, errors.New(typeMap.GetString("status") + " : " + typeMap.GetString("createdDateTime"))
			}
			time.Sleep(time.Second * 3)
		}
	}
	return nil, errors.New("GetPrebuiltReadResultWithBlock unknown error")
}
