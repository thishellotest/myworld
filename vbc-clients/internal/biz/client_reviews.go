package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

/*
CREATE TABLE `client_reviews` (

	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`full_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`branch_of_service` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`rating_goal_met` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`recommendation_consent` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`communication_rating` int(11) NOT NULL DEFAULT '0',
	`overall_rating` int(11) NOT NULL DEFAULT '0',
	`testimonial_text` text COLLATE utf8mb4_unicode_ci,
	`improvement_feedback` text COLLATE utf8mb4_unicode_ci,
	`allow_testimonial_use` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`status` tinyint(4) NOT NULL DEFAULT '1',
	`created_at` int(11) NOT NULL DEFAULT '0',
	`updated_at` int(11) NOT NULL DEFAULT '0',
	`deleted_at` int(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='client_reviews';
*/
const (
	ClientReview_RatingGoalMet_Yes         = "Yes"
	ClientReview_RatingGoalMet_No          = "No"
	ClientReview_RecommendationConsent_Yes = "Yes"
	ClientReview_RecommendationConsent_No  = "No"
	ClientReview_AllowTestimonialUse_Yes   = "Yes"
	ClientReview_AllowTestimonialUse_No    = "No"
	ClientReview_Status_Enable             = 1
	ClientReview_Status_Disable            = 0
)

type ClientReviewEntity struct {
	ID                    int32 `gorm:"primaryKey"`
	FullName              string
	SubmissionDate        string
	BranchOfService       string
	RatingGoalMet         string
	RecommendationConsent string
	CommunicationRating   int32
	OverallRating         int32
	TestimonialText       string
	ImprovementFeedback   string
	AllowTestimonialUse   string
	Status                int
	Sort                  int
	CreatedAt             int64
	UpdatedAt             int64
	DeletedAt             int64
}

type ClientReviewVo struct {
	ID                  int32  `json:"id"`
	FullName            string `json:"name"`
	BranchOfService     string `json:"branch"`
	CommunicationRating int32  `json:"communication_rating"`
	TestimonialText     string `json:"quote"`
}

func FormatFullNameToShort(fullname string) string {
	if fullname == "" {
		return ""
	}
	res := strings.Split(fullname, " ")
	if len(res) == 1 {
		return lib.Capitalize(strings.ToLower(res[0]))
	}
	c := res[len(res)-1]
	runes := []rune(c)
	lastFirstLetter := strings.ToUpper(string(runes[0]))
	return fmt.Sprintf("%s %s.", lib.Capitalize(strings.ToLower(res[0])), lastFirstLetter)
}
func (c *ClientReviewEntity) ToClientReviewVo() ClientReviewVo {

	clientReviewVo := ClientReviewVo{
		ID:                  c.ID,
		FullName:            FormatFullNameToShort(c.FullName),
		BranchOfService:     c.BranchOfService,
		CommunicationRating: c.CommunicationRating,
		TestimonialText:     c.TestimonialText,
	}
	return clientReviewVo
}

func (ClientReviewEntity) TableName() string {
	return "client_reviews"
}

type ClientReviewUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[ClientReviewEntity]
}

func NewClientReviewUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *ClientReviewUsecase {
	uc := &ClientReviewUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *ClientReviewUsecase) ImportExcel() error {
	file := "/Users/garyliao/code/vbc-clients/tests/docx/Ed Edits - Share_Your_Experience2025-06-16_21_17_58 (2).xlsx"
	f, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get value from cell by given worksheet name and cell reference.
	//cell, err := f.GetCellValue("Sheet1", "B2")
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return err
	}
	for k, row := range rows {
		if len(row) > 0 {
			if strings.Index(row[0], "Submission Date") >= 0 {
				continue
			}
		}
		clientReviewEntity := ClientReviewEntity{}
		clientReviewEntity.Sort = 1000 - k
		for k, colCell := range row {
			if k == 0 {
				submissionDate, err := time.ParseInLocation(configs.TimeFormatDate, colCell, configs.VBCDefaultLocation)
				if err != nil {
					panic(err)
				}
				clientReviewEntity.SubmissionDate = submissionDate.Format(time.DateOnly)

			} else if k == 1 {
				clientReviewEntity.FullName = strings.TrimSpace(colCell)
			} else if k == 2 {

			} else if k == 3 {
				clientReviewEntity.BranchOfService = strings.TrimSpace(colCell)
			} else if k == 4 {
				if colCell == "Yes" {
					clientReviewEntity.RatingGoalMet = ClientReview_RatingGoalMet_Yes
				} else {
					clientReviewEntity.RatingGoalMet = ClientReview_RatingGoalMet_No
				}
			} else if k == 5 {
				if colCell == "Yes" {
					clientReviewEntity.RecommendationConsent = ClientReview_RecommendationConsent_Yes
				} else {
					clientReviewEntity.RecommendationConsent = ClientReview_RecommendationConsent_No
				}
			} else if k == 6 {
				a, _ := strconv.ParseInt(colCell, 10, 32)
				if a >= 0 && a <= 5 {
					clientReviewEntity.CommunicationRating = int32(a)
				}
			} else if k == 7 {
				a, _ := strconv.ParseInt(colCell, 10, 32)
				if a >= 0 && a <= 5 {
					clientReviewEntity.OverallRating = int32(a)
				}
			} else if k == 8 {
				clientReviewEntity.TestimonialText = strings.TrimSpace(colCell)
			} else if k == 9 {
				clientReviewEntity.ImprovementFeedback = strings.TrimSpace(colCell)
			} else if k == 10 {
				if colCell == "Yes" {
					clientReviewEntity.AllowTestimonialUse = ClientReview_AllowTestimonialUse_Yes
				} else {
					clientReviewEntity.AllowTestimonialUse = ClientReview_AllowTestimonialUse_No
				}

			} else if k == 11 {
				if colCell == "yes" {
					clientReviewEntity.Status = ClientReview_Status_Enable
				} else {
					clientReviewEntity.Status = ClientReview_Status_Disable
				}
			}
		}
		a, err := c.GetByCond(Eq{"full_name": clientReviewEntity.FullName})
		if err != nil {
			panic(err)
		}
		if a == nil {
			err = c.CommonUsecase.DB().Save(&clientReviewEntity).Error
			if err != nil {
				c.log.Error(err)
			}
		}
		fmt.Println()
	}
	return nil
}
