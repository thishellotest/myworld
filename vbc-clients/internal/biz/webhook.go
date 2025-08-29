package biz

type WebhookUsecase struct {
	CommonUsecase *CommonUsecase
}

func NewWebhookUsecase(CommonUsecase *CommonUsecase) *WebhookUsecase {

	return &WebhookUsecase{
		CommonUsecase: CommonUsecase,
	}
}
