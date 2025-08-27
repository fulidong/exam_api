package main

import (
	"exam_api/internal/middleware"
	"exam_api/internal/pkg/ijwt"
	"exam_api/internal/pkg/isnowflake"
	"flag"
	"github.com/airunny/wiki-go-tools/ilog"
	"os"
	"time"

	"exam_api/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
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
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	logger, closer := ilog.NewLogger(id, Name)
	defer closer.Close()

	// 初始化jwt
	middleware.JWT = ijwt.NewSecureJWT(bc.Data.Jwt.AccessSecret, bc.Data.Jwt.ExamSecret, []ijwt.SecurityOption{
		ijwt.WithAccessExpiry(time.Minute * time.Duration(bc.Data.Jwt.AccessTokenExpireMinutes)),
	}...)
	// 初始化雪花算法
	snowFlake, err := isnowflake.NewSnowflake(1)
	if err != nil {
		panic(err)
	}
	isnowflake.SnowFlake = snowFlake
	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
