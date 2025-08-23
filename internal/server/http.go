package server

import (
	v1 "exam_api/api/exam_api/v1"
	"exam_api/internal/conf"
	"exam_api/internal/middleware"
	"exam_api/internal/pkg/ilog"
	"exam_api/internal/service"
	"github.com/airunny/wiki-go-tools/env"
	"github.com/go-kratos/grpc-gateway/v2/protoc-gen-openapiv2/generator"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, exam *service.ExamService, logger log.Logger) *http.Server {
	serviceName := env.GetServiceName()
	var opts = []http.ServerOption{
		http.Filter(middleware.CORS(), ilog.LoggingHandler(serviceName, ilog.WithAccessLog())),
		http.Middleware(
			recovery.Recovery(),
			middleware.RequestIdWithHeader,
			middleware.TryParseHeader(),
			middleware.CheckExamTokenMiddleware(),
			validate.Validator(),
		),
		http.ErrorEncoder(middleware.ErrorEncoder),
		http.ResponseEncoder(middleware.ResponseEncoder),
		http.StrictSlash(false),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterExamServiceHTTPServer(srv, exam)
	openAPIHandler := openapiv2.NewHandler(openapiv2.WithGeneratorOptions(
		generator.UseJSONNamesForFields(false),
		generator.EnumsAsInts(true),
	))
	srv.HandlePrefix("/q/", openAPIHandler)

	srv.Handle("/metrics", promhttp.Handler())
	return srv
}
