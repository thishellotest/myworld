package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type AdobesignSyncTaskUsecase struct {
	CommonUsecase          *CommonUsecase
	log                    *log.Helper
	DataEntryUsecase       *DataEntryUsecase
	AdobeSignUsecase       *AdobeSignUsecase
	ClientAgreementUsecase *ClientAgreementUsecase
	SyncTask
}

func NewAdobesignSyncTaskUsecase(CommonUsecase *CommonUsecase,
	logger log.Logger,
	AdobeSignUsecase *AdobeSignUsecase,
	DataEntryUsecase *DataEntryUsecase, ClientAgreementUsecase *ClientAgreementUsecase,
) *AdobesignSyncTaskUsecase {

	adobesignSyncTaskUsecase := &AdobesignSyncTaskUsecase{
		CommonUsecase:          CommonUsecase,
		log:                    log.NewHelper(logger),
		AdobeSignUsecase:       AdobeSignUsecase,
		DataEntryUsecase:       DataEntryUsecase,
		ClientAgreementUsecase: ClientAgreementUsecase,
	}
	adobesignSyncTaskUsecase.SyncTask.RedisQueue = Redis_sync_adobesign_tasks_queue
	adobesignSyncTaskUsecase.SyncTask.RedisProcessing = Redis_sync_adobesign_tasks_processing
	adobesignSyncTaskUsecase.SyncTask.RedisClient = CommonUsecase.RedisClient()
	adobesignSyncTaskUsecase.SyncTask.Log = log.NewHelper(logger)
	adobesignSyncTaskUsecase.SyncTask.Handle = adobesignSyncTaskUsecase.HandleTask

	return adobesignSyncTaskUsecase
}

// HandleTask 同步单个task
func (c *AdobesignSyncTaskUsecase) HandleTask(ctx context.Context, agreementId string) error {

	agreement, err := c.AdobeSignUsecase.GetAgreement(ctx, agreementId)
	if err != nil {
		return err
	}
	return c.ClientAgreementUsecase.Update(agreement)
}
