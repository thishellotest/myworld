package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_UserHttpUsecase_aaa222(t *testing.T) {
	tUser, _ := UT.TUsecase.DataById(biz.Kind_users, 4)
	//a, _ := UT.FieldUsecase.CacheStructByKind(biz.Kind_users)
	//for _, v := range a.Records {
	//	lib.DPrintln(v.FieldName, " ", v.FieldType)
	//}
	lib.DPrintln("====:", tUser.CustomFields.NumberValueByNameBasic(biz.UserFieldName_status))
}
func Test_UserHttpUsecase_BizVerifyEmailOutbox(t *testing.T) {

	/*
		INSERT INTO `users` (`id`, `email`, `microsoft_email`, `password`, `full_name`, `name`, `first_name`, `last_name`, `title`, `mobile`, `mail_username`, `mail_password`, `pic_url`, `profile_gid`, `role_gid`, `created_at`, `updated_at`, `deleted_at`, `gid`, `dialpad_userid`, `dialpad_phonenumber`, `created_time`, `modified_time`, `created_by`, `modified_by`, `status`, `timezone_id`, `biz_deleted_at`)
		VALUES
			(17, 'dharris@vetbenefitscenter.com', '', '', 'Debra Harris', 'Debra Harris', 'Debra', 'Harris', 'Veteran Services Manager', '+1 720-600-2216', 'dharris@vetbenefitscenter.com', 'ixkkqczucgqfhnqd', '', '441f19d51858417cb948cc286ef1b585', '', 1736204400, 1743475747, 0, 'c9ce3ecee21640e7978a373c08d21292', '6374258379636736', '+17206002216', '2025-01-07T06:53:45+08:00', '2025-02-05T13:40:26+08:00', '', '6159272000000453669', 1, '', 0);

	*/
	tUser, _ := UT.TUsecase.DataById(biz.Kind_users, 4)
	aa := biz.UserFacade{
		TData: *tUser,
	}
	aaa, err := UT.UserHttpUsecase.BizVerifyEmailOutbox(aa, "6159272000000453669")
	lib.DPrintln(err)
	lib.DPrintln(aaa)
}

func Test_UserHttpUsecase_BizSyncDailpad(t *testing.T) {

	/*
		INSERT INTO `users` (`id`, `email`, `microsoft_email`, `password`, `full_name`, `name`, `first_name`, `last_name`, `title`, `mobile`, `mail_username`, `mail_password`, `pic_url`, `profile_gid`, `role_gid`, `created_at`, `updated_at`, `deleted_at`, `gid`, `dialpad_userid`, `dialpad_phonenumber`, `created_time`, `modified_time`, `created_by`, `modified_by`, `status`, `timezone_id`, `biz_deleted_at`)
		VALUES
			(17, 'dharris@vetbenefitscenter.com', '', '', 'Debra Harris', 'Debra Harris', 'Debra', 'Harris', 'Veteran Services Manager', '+1 720-600-2216', 'dharris@vetbenefitscenter.com', 'ixkkqczucgqfhnqd', '', '441f19d51858417cb948cc286ef1b585', '', 1736204400, 1743475747, 0, 'c9ce3ecee21640e7978a373c08d21292', '6374258379636736', '+17206002216', '2025-01-07T06:53:45+08:00', '2025-02-05T13:40:26+08:00', '', '6159272000000453669', 1, '', 0);

	*/
	tUser, _ := UT.TUsecase.DataById(biz.Kind_users, 4)
	aa := biz.UserFacade{
		TData: *tUser,
	}
	aaa, err := UT.UserHttpUsecase.BizSyncDailpad(aa, "c9ce3ecee21640e7978a373c08d21292")
	lib.DPrintln(err)
	lib.DPrintln(aaa)
}

func Test_FormatPhoneNumberV2(t *testing.T) {
	a, err := biz.FormatPhoneNumberV2("+13109719600")
	lib.DPrintln(a, err)
}
