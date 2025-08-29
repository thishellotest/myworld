package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

/*
CREATE TABLE `statement_conditions` (

	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`case_id` int(11) NOT NULL DEFAULT '0',
	`uuid` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`front_value` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`condition_value` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
	`behind_value` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
	`sort` int(11) NOT NULL DEFAULT '1000',
	`created_at` int(11) NOT NULL DEFAULT '0',
	`updated_at` int(11) NOT NULL DEFAULT '0',
	`deleted_at` int(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`),
	KEY `uniq` (`case_id`,`uuid`(191))

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='statement_conditions';
*/
type StatementConditionEntity struct {
	ID             int32 `gorm:"primaryKey"`
	CaseId         int32
	Uuid           string
	FrontValue     string
	ConditionValue string
	BehindValue    string
	Category       string
	Sort           int
	CreatedAt      int64
	UpdatedAt      int64
	DeletedAt      int64
}

func (c *StatementConditionEntity) ToStatementCondition() (statementCondition StatementCondition) {
	statementCondition.StatementConditionId = c.ID
	//statementCondition.StatementConditionUuid = c.Uuid
	statementCondition.Sort = c.Sort
	statementCondition.ConditionValue = c.ConditionValue
	statementCondition.FrontValue = c.FrontValue
	statementCondition.BehindValue = c.BehindValue
	statementCondition.Category = c.Category
	return statementCondition
}

func (c *StatementConditionEntity) ToStatementConditionJobUuid() string {
	return GenStatementConditionJobUuid(c.CaseId, c.ID)
	//return fmt.Sprintf("%s:%d:%d", AiAssistantJobBizType_statementCondition, c.CaseId, c.ID)
}

func GenStatementConditionJobUuid(caseId int32, StatementConditionId int32) string {
	return fmt.Sprintf("%s:%d:%d", AiAssistantJobBizType_statementCondition, caseId, StatementConditionId)
}

func GenAllStatementsJobUuid(caseId int32) string {
	return fmt.Sprintf("%s:%d", AiAssistantJobBizType_allStatements, caseId)
}

func GenDocEmailJobUuid(caseId int32) string {
	return fmt.Sprintf("%s:%d", AiAssistantJobBizType_genDocEmail, caseId)
}

func (c *StatementConditionEntity) ToCondition() (condition string) {
	if c.FrontValue != "" {
		condition = c.FrontValue + " - "
	}
	condition += c.ConditionValue
	if c.BehindValue != "" {
		condition += " (" + c.BehindValue + ")"
	}
	return condition
}

func (StatementConditionEntity) TableName() string {
	return "statement_conditions"
}

type StatementConditionUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[StatementConditionEntity]
}

func NewStatementConditionUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *StatementConditionUsecase {
	uc := &StatementConditionUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func StatementConditionsToTextNoDivide(res []*StatementConditionEntity) (text string, count int) {
	//text := ""
	for _, v := range res {
		aa := v.ToCondition()
		aa = "- **" + aa + "**"
		if text == "" {
			text = aa
		} else {
			text += "\n" + aa
		}
	}
	return text, len(res)
}

func (c *StatementConditionUsecase) AllConditions(caseId int32) ([]*StatementConditionEntity, error) {
	return c.AllByCondWithOrderBy(Eq{"case_id": caseId, "deleted_at": 0}, "sort asc", 1000)
}

func (c *StatementConditionUsecase) GetCondition(caseId int32, statementConditionId int32) (*StatementConditionEntity, error) {
	return c.GetByCond(Eq{"case_id": caseId, "id": statementConditionId, "deleted_at": 0})
}

func (c *StatementConditionUsecase) Upsert(caseId int32, Uuid string, FrontValue string, ConditionValue string, BehindValue string, Sort int, category string) (entity *StatementConditionEntity, err error) {

	entity, err = c.GetByCond(Eq{"case_id": caseId, "uuid": Uuid})
	if err != nil {
		return nil, err
	}
	if entity == nil {
		entity = &StatementConditionEntity{
			CaseId:    caseId,
			Uuid:      Uuid,
			CreatedAt: time.Now().Unix(),
		}
	}
	entity.Sort = Sort
	entity.FrontValue = FrontValue
	entity.ConditionValue = ConditionValue
	entity.BehindValue = BehindValue
	entity.Sort = Sort
	entity.DeletedAt = 0
	entity.UpdatedAt = time.Now().Unix()
	entity.Category = category
	err = c.CommonUsecase.DB().Save(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
