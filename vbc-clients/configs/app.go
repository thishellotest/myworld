package configs

import (
	"fmt"
	"math"
	"os"
	"time"
	"vbc/lib"
)

/*
characters

 • abc
 • abc1

*/

var AppName = ""
var AppLaunchAt time.Time // 应用启动的时间

const (
	App_UnitTest = "unit-test"
	App_vbc      = "vbc"
)

const (
	TimeFormatDate         = "Jan 02, 2006"
	TimeFormatDate2        = "January 02, 2006"
	TimeFormatDateThisYear = "Jan 02"
	TimeFormatDateTime     = "January 2, 2006 03:04 PM"
)

const SubmissionToGoogleDriveFailedNotifyEmail = "venriquez@vetbenefitscenter.com"
const EdEmail = "ebunting@vetbenefitscenter.com"
const YnEmail = "ywang@vetbenefitscenter.com"
const GaryEmail = "glliao@vetbenefitscenter.com"

const DCFolderId = "263406803830"

var Domain = "https://base.vetbenefitscenter.com"

// EnableNewUpdateQuestionnaireVersion 开启新版本，使用多个任务
const EnableNewUpdateQuestionnaireVersion = true

const EnableAiAutoAssociationJotform = true // 先关闭，等收集数据
const NewPSGen = true

const EnableNewJotformName = true

// 停用zoho
const StoppedZoho = true

// CRM上线
const VBC_CRM_RELEASE = true

// 启用订阅管控制
const Enable_Unsubscribes_SMS = true

// 开启 debug方式， 短信不会实际发送，邮件发送到测试帐号； 上线后要改为 false(2025-02-09上线了)
const Enable_SMS_New_Version_Debug = false

// 开启新的版本优化-解决性能问题和KIND依赖自己的问题
const Enable_NewVersionForT = true

// 不使用client tasks(太复杂了)，直接使用client cases的due date
const Enable_Client_Task_ForCRM = false

// WorkflowDebugTest 开启此处时，邮箱为：liaogling@gmail.com lialing@foxmail.com 不发送邮件和短信，方便调试
// 此开关需保证任务状态都不影响业务
var WorkflowDebugTest = true

func IsWorkflowDebug(email string) bool {
	if WorkflowDebugTest && (email == "liaogling@gmail.com" || email == "lialing@foxmail.com") {
		return true
	}
	return false
}

// 非测试时，请设置为1倍（原速）， 测试时，可设置大于1的整数
var ZohoContactAndDealSyncSlowTimes = 1

// DebugMedTeamFormBoxSign 开启box测试，只允许部分通过
var DebugMedTeamFormBoxSign = false

var EnabledAmIntakeFormReminder = true

var EnabledContractReminder = true
var EnabledContractReminderBySMS = true

// EnabledDataEntryDependField 是否打开字段值依懒， 核心，如果有bug请设置为false
var EnabledDataEntryDependField = true

// EnabledTwoBySMS 开启第二部分运营发送, 注意：开启会影响线上, 先暂停，因为要改stages
var EnabledTwoBySMS = true

// EnabledPrimaryCase 是否开启primary case计算方法， 防止出bug后，不好回滚
var EnabledPrimaryCase = true

// EnabledDBPricingVersion 使用DB控制版本价格
var EnabledDBPricingVersion = true

// StopNotifyCaseInDebug 默认应该设置为 false;  当开启调试时，设置为true
var StopNotifyCaseInDebug = false

// UseOwnerSendingSMS 使用case的owner发送sms
var UseOwnerSendingSMS = false

var LoadLocation *time.Location
var VBCDefaultLocation *time.Location
var AppRuntimePath string

const ENV_PROD = "PROD"
const ENV_TEST = "TEST"
const ENV_DEV = "DEV"

func AppEnv() string {
	return os.Getenv("ENV")
}

func AppEnvType() string {
	return os.Getenv("ENV_TYPE")
}

const ENV_TYPE_DEV_TEST = "DEV_TEST"

const JOB_TYPE_DEFAULT = "DEFAULT"
const JOB_TYPE_LARGE_MEMORY = "LARGE_MEMORY"
const JOB_TYPE_QA = "QA"

func AppJobType() string {
	return os.Getenv("JOB_TYPE")
}

func IsJobTypeQA() bool {
	if AppJobType() == JOB_TYPE_QA {
		return true
	}
	return false
}

func IsJobTypeDefault() bool {
	if AppJobType() == JOB_TYPE_DEFAULT {
		return true
	}
	return false
}

func IsJobTypeLargeMemory() bool {
	if AppJobType() == JOB_TYPE_LARGE_MEMORY {
		return true
	}
	return false
}

func IsProd() bool {
	if AppEnv() == ENV_PROD {
		return true
	}
	return false
}

// IsCrmProhibitForProd crm禁止给prod使用
func IsCrmProhibitForProd() bool {
	return true
}

func IsTest() bool {
	if AppEnv() == ENV_TEST {
		return true
	}
	return false
}

func InitApp(appName string) {
	var err error
	LoadLocation = time.FixedZone("CST", 0)
	time.Local = LoadLocation
	AppName = appName
	AppLaunchAt = time.Now()

	AppRuntimePath, err = lib.RuntimePath()
	lib.DPrintln("AppRuntimePath:", AppRuntimePath)
	if err != nil {
		panic(err)
	}

	VBCDefaultLocation = California()
}

func GetVBCDefaultLocation() *time.Location {
	if VBCDefaultLocation == nil {
		VBCDefaultLocation = California()
	}
	return VBCDefaultLocation
}

// GetAppRuntimePath Example: /Users/.../bin
func GetAppRuntimePath() string {
	return AppRuntimePath
}

// CacheExpiredDurationDefault 缓存默认过期时间
var CacheExpiredDurationDefault = time.Second * 20
var CacheExpiredDuration5Seconds = time.Second * 15

func IsDev() bool {
	if AppEnv() == ENV_DEV {
		return true
	}
	return false
}

func California() *time.Location {
	la, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic("should not error on generating current time in la")
	}
	return la
}

func GetTempDir() string {
	return GetAppRuntimePath()
}

func MkdirTemp() (string, error) {
	tempDir, err := os.MkdirTemp(GetAppRuntimePath(), "")
	return tempDir, err
}

func PollingTransIdByTime(currentTime time.Time) (str string, beginAt time.Time) {
	oriTime := currentTime
	gapMinute := PollingGapMinute()
	if currentTime.Minute() == 0 {
		currentTime = currentTime.Add(-1 * time.Hour)
		oriTime = oriTime.Add(-time.Duration(gapMinute) * time.Minute)
	}

	minuteToGap := MinuteToGap(currentTime.Minute(), gapMinute)
	// 2006-01-02 15:04:05
	beginStr := oriTime.Format("2006-01-02 15:")
	a := gapMinute * (minuteToGap - 1)
	if a == 60 {
		a = 0
	}
	beginStr = fmt.Sprintf("%s%02d:00", beginStr, a)
	beginAt, _ = time.ParseInLocation("2006-01-02 15:04:05", beginStr, LoadLocation)
	str = fmt.Sprintf("%s_%d", currentTime.Format("2006010215"), minuteToGap)
	return str, beginAt
}

// MinuteToGap 分钟转为巡检的间隔
func MinuteToGap(minute int, gapMinute int) int {
	if minute == 0 {
		minute = 60
	}
	return int(math.Ceil(float64(minute) / float64(gapMinute)))
}

// PollingGapMinute 时间间隔，注意需要被60整除
func PollingGapMinute() int {
	return 5
}
