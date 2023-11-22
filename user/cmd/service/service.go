package service

import (
	"context"
	"flag"
	"fmt"
	"net"
	http1 "net/http"

	_ "net/http/pprof"

	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"

	// "time"

	endpoint "user/pkg/endpoint"
	grpc "user/pkg/grpc"
	pb "user/pkg/grpc/pb"
	service "user/pkg/service"
	postgresStore "user/pkg/store/postgres"

	// oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	endpoint1 "github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	group "github.com/oklog/oklog/pkg/group"
	opentracinggo "github.com/opentracing/opentracing-go"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	grpc1 "google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	// "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var tracer opentracinggo.Tracer
var logger log.Logger

// Define our flags. Your service probably won't need to bind listeners for
// all* supported transports, but we do it here for demonstration purposes.
var fs = flag.NewFlagSet("user", flag.ExitOnError)
var debugAddr = fs.String("debug-addr", ":8080", "Debug and metrics listen address")
var httpAddr = fs.String("http-addr", ":8081", "HTTP listen address")
var grpcAddr = fs.String("grpc-addr", ":8082", "gRPC listen address")
var thriftAddr = fs.String("thrift-addr", ":8083", "Thrift listen address")
var thriftProtocol = fs.String("thrift-protocol", "binary", "binary, compact, json, simplejson")
var thriftBuffer = fs.Int("thrift-buffer", 0, "0 for unbuffered")
var thriftFramed = fs.Bool("thrift-framed", false, "true to enable framing")
var otelGrpcUrl = fs.String("otel-grpc-url", "", "Enable opentelemetry tracing via a grpc exporter URL. Default: localhost:4317")

var logLevel = fs.String("log-level", "INFO", "Sets log levels. Valid values are DEBUG, INFO, WARN, ERROR")
var logFormat = fs.String("log-format", "fmt", "Log format type. Valid values: fmt, json")
var postgresConnString = fs.String("postgres-conn-url", "", "Enable postgres store and connects to the provided connection string url.")

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			fmt.Println(setting.Key)
			if setting.Key == "vcs.revision" {
				fmt.Println("Setting", setting.Value)
				return setting.Value
			}
		}
	}
	return ""
}()

func Run() {
	fs.Parse(os.Args[1:])

	switch strings.ToLower(*logFormat) {
	case "json":
		logger = log.NewJSONLogger(os.Stderr)
	default:
		logger = log.NewLogfmtLogger(os.Stderr)

	}

	switch *logLevel {
	case "DEBUG":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "WARN":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "ERROR":
		logger = level.NewFilter(logger, level.AllowError())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	logger.Log("Version", Commit)

	//  Determine which tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency
	if *otelGrpcUrl != "" {
		logger.Log("tracer", "OTEL", "GRPC_URL", *otelGrpcUrl)

		tp, err := initTracerProvider(*otelGrpcUrl)

		if err != nil {
			logger.Log(err)
		}
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				logger.Log("Tracer Provider Shutdown: %v", err)
			}
		}()

		// mp := initMeterProvider()
		// defer func() {
		// 	if err := mp.Shutdown(context.Background()); err != nil {
		// 		logger.Log("Error shutting down meter provider: %v", err)
		// 	}
		// }()

		// err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
		// if err != nil {
		// 	logger.Log(err)
		// }
	}

	var ps *postgresStore.PostgresStore

	if *postgresConnString != "" {
		logger.Log("tracer", "PostgresStore", "URL", postgresConnString)
		var err error
		ps, err = postgresStore.NewPostgresStore(*postgresConnString, log.With(logger, "store", "postgres"))

		if err != nil {
			logger.Log("tracer", "PostgresStore", "URL", err)
			os.Exit(1)
		}
	} else {
		logger.Log("tracer", "Only postgres store supported now. Please enable using --postgres-conn-url")
		os.Exit(1)
	}

	svc, err := service.New(ps, ps, log.With(logger, "service", "userService"), getServiceMiddleware(logger))

	if err != nil {
		logger.Log("tracer", "Failed to start service")
		os.Exit(1)
	}

	eps := endpoint.New(svc, getEndpointMiddleware(logger))
	g := createService(eps)
	initMetricsEndpoint(g)
	initCancelInterrupt(g)
	logger.Log("exit", g.Run())

}
func initGRPCHandler(endpoints endpoint.Endpoints, g *group.Group) {
	options := defaultGRPCOptions(logger, tracer)
	// Add your GRPC options here

	grpcServer := grpc.NewGRPCServer(endpoints, options)
	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "gRPC", "addr", *grpcAddr)
		baseServer := grpc1.NewServer(
			grpc1.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		)
		pb.RegisterUserServer(baseServer, grpcServer)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
	})

}
func getServiceMiddleware(logger log.Logger) (mw []service.Middleware) {
	mw = []service.Middleware{}
	// Append your middleware here

	return
}
func getEndpointMiddleware(logger log.Logger) (mw map[string][]endpoint1.Middleware) {
	mw = map[string][]endpoint1.Middleware{}
	// Add you endpoint middleware here

	return
}

func newResource(serviceName, serviceVersion string) (*sdkresource.Resource, error) {
	return sdkresource.Merge(sdkresource.Default(),
		sdkresource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

func initMetricsEndpoint(g *group.Group) {

	http1.DefaultServeMux.Handle("/metrics", promhttp.Handler())
	debugListener, err := net.Listen("tcp", *debugAddr)
	if err != nil {
		logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "debug/HTTP", "addr", *debugAddr)
		return http1.Serve(debugListener, http1.DefaultServeMux)
	}, func(error) {
		debugListener.Close()
	})
}

func initTracerProvider(traceEndpoint string) (*sdktrace.TracerProvider, error) {
	resource, err := newResource("user", "0.1.0")

	if err != nil {
		logger.Log("msg", "Failed to create resource")
	}

	ctx := context.Background()

	if traceEndpoint == "" {
		traceEndpoint = "localhost:4317"
	}

	logger.Log("traceEndpoint: ", traceEndpoint)

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(traceEndpoint))

	if err != nil {
		logger.Log("OTLP Trace gRPC Creation: %v", err)
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func initMeterProvider() *sdkmetric.MeterProvider {
	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
	if err != nil {
		logger.Log("new otlp metric grpc exporter failed: %v", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(initResource()),
	)
	otel.SetMeterProvider(mp)
	return mp
}

var (
	resource          *sdkresource.Resource
	initResourcesOnce sync.Once
)

func initResource() *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
		)
		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)
	})
	return resource
}

func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}
