package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

const (
	StatementComment_IsComplete_Yes = 1
	StatementComment_IsComplete_No  = 0
)

type StatementCommentEntity struct {
	ID                   int32 `gorm:"primaryKey"`
	CaseId               int32
	StatementConditionId int32
	StatementSection     string
	Text                 string
	IsComplete           int
	ModifiedBy           string
	CreatedAt            int64
	UpdatedAt            int64
	DeletedAt            int64
}
type StatementCommentVo struct {
	ID                   int32  `json:"id"`
	SenderName           string `json:"sender_name"`
	CreatedAt            int32  `json:"created_at"`
	StatementConditionId int32  `json:"statement_condition_id"`
	StatementSection     string `json:"statement_section"`
	Text                 string `json:"text"`
	ConditionLabel       string `json:"condition_label"`
	SectionLabel         string `json:"section_label"`
	ModifiedBy           string `json:"modified_by"`
	IsComplete           bool   `json:"is_complete"`
}
type ListStatementCommentVo []StatementCommentVo

func (c *StatementCommentEntity) ToStatementCommentVo(statementConditions map[int32]*StatementConditionEntity) (vo StatementCommentVo) {
	if c.ModifiedBy == "" {
		vo.SenderName = "Collaborator"
	} else {
		vo.SenderName = "Staff"
	}
	vo.ID = c.ID
	vo.CreatedAt = int32(c.CreatedAt)
	vo.StatementConditionId = c.StatementConditionId
	vo.StatementSection = c.StatementSection
	vo.Text = c.Text
	if v, ok := statementConditions[c.StatementConditionId]; ok {
		vo.ConditionLabel = v.ConditionValue
	}
	if c.StatementSection != "" {
		vo.SectionLabel = GetSectionTitleFromSectionType(c.StatementSection)
	}
	vo.ModifiedBy = c.ModifiedBy
	if c.IsComplete == StatementComment_IsComplete_Yes {
		vo.IsComplete = true
	} else {
		vo.IsComplete = false
	}
	return vo
}

func (StatementCommentEntity) TableName() string {
	return "statement_comments"
}

type StatementCommentUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[StatementCommentEntity]
}

func NewStatementCommentUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *StatementCommentUsecase {
	uc := &StatementCommentUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}
