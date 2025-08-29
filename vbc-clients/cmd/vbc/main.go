package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/server"
	"vbc/lib"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs/config_dev.yaml", "config path, eg: -conf config_1.yaml")
	flag.Parse()
	configs.InitApp(configs.App_vbc)

}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, job *server.Job) *kratos.App {

	server.VbcJobManager = lib.NewJobManager(func(ctx context.Context) {
		// 所有任务停止zoho版本
		err := job.Run(ctx)
		if err != nil {
			panic(err)
		}
		//}
	})
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			//gs,
			hs,
		),
		kratos.BeforeStop(func(ctx context.Context) error {
			return nil
		}),
	)
}

func main() {

	fmt.Println("main_flagconf_env:", configs.AppEnv())
	confYaml := ""
	if configs.AppEnv() == configs.ENV_PROD {
		confYaml = "/app/configs/config_prod.yaml"
	} else if configs.AppEnv() == configs.ENV_TEST {
		confYaml = "/app/configs/config_test.yaml"
	} else {
		confYaml = flagconf
	}

	if configs.IsProd() && configs.IsJobTypeLargeMemory() {
		confYaml = flagconf
	}
	fmt.Println("confYaml", confYaml)
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(confYaml),
		),
	)
	defer c.Close()

	logHelper := log.NewHelper(logger)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	bc.Server.Http.Addr = "0.0.0.0:" + os.Getenv("PORT")

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	go func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	jobManagerContext := context.Background()
	err = server.VbcJobManager.Start(jobManagerContext)
	if err != nil {
		panic(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.VbcJobManager.Cancel()
	time.Sleep(6 * time.Second)
	logHelper.Info("The container stopped.")
}
