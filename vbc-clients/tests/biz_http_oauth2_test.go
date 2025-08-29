package tests

import (
	"context"
	"testing"
)

func Test_HttpOauth2Usecase_aa(t *testing.T) {

	code := "4/0AcvDMrDN5tRFkUUGC8DwKyT5cY2V_kpxFl9eY0OA1ORh28s53VDFt-lJOa57V7oBDdMmUQ"
	UT.HttpOauth2Usecase.BizVbcapp(context.Background(), code)
}
