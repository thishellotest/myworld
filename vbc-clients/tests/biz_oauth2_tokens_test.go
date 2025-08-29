package tests

import "testing"

func Test_handleRefreshToken_HandleRefreshToken(t *testing.T) {
	UT.Oauth2TokenUsecase.HandleRefreshToken()
}
