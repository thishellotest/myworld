//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package tests

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"vbc/internal/biz"
	"vbc/internal/conf"
	"vbc/internal/data"
)

// initApp init kratos application.
func initApp(*conf.Data, log.Logger) (*UnittestApp, func(), error) {
	panic(wire.Build(
		biz.ProviderSet,
		data.ProviderSet,
		//packet.ProviderSet,
		newApp))
}
