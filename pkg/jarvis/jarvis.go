package jarvis

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/fizx/baseplate.go/thriftbp"
	"github.com/getsentry/raven-go"
	"github.com/reddit/baseplate.go/edgecontext"
	"github.com/reddit/baseplate.go/log"
	"github.com/reddit/baseplate.go/metricsbp"
	"github.com/reddit/baseplate.go/runtimebp"
	"github.com/reddit/baseplate.go/secrets"
	"github.com/reddit/baseplate.go/tracing"
	"gopkg.in/yaml.v2"
)

var configPath = flag.String("config", "", "Path to the config file to use.")
var debug = flag.Bool("debug", false, "Should the service be run in 'debug' mode.")
var serverEndpoint = flag.String("bind", ":8080", "The endpoint `{ipAddress}:{port}` to bind the server to.")

type config struct {
	Metrics struct {
		Namespace string
		Endpoint  string
	}

	Redis struct {
		Endpoints []string
	}

	Secrets struct {
		Path string
	}

	Sentry struct {
		DSN         string
		Environment string
		SampleRate  float64
	}

	Tracing struct {
		Namespace     string
		Endpoint      string
		RecordTimeout time.Duration `yaml:"recordTimeout"`
		SampleRate    float64
	}
}

func StartBaseplateThrift(processor thrift.TProcessor) {
	prev, current := runtimebp.GOMAXPROCS(1, 50)
	fmt.Println("GOMAXPROCS:", prev, current)
	flag.Parse()

	logLevel, logger := initLogger(*debug)
	cfg, err := parseConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricsbp.M = metricsbp.NewStatsd(ctx, metricsbp.StatsdConfig{
		Prefix:   cfg.Metrics.Namespace,
		Address:  cfg.Metrics.Endpoint,
		LogLevel: logLevel,
	})

	secretsStore, err := initSecrets(ctx, cfg, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer secretsStore.Close()

	if _ = edgecontext.Init(edgecontext.Config{Store: secretsStore}); err != nil {
		log.Fatal(err)
	}
	if err = initTracing(cfg, logger, metricsbp.M); err != nil {
		log.Fatal(err)
	}
	if err = initSentry(cfg); err != nil {
		log.Fatal(err)
	}

	transport, err := thrift.NewTServerSocketTimeout(*serverEndpoint, 1*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	server := thrift.NewTSimpleServer4(
		processor,
		transport,
		thrift.NewTHeaderTransportFactory(nil),
		thrift.NewTHeaderProtocolFactory(),
	)
	server.SetForwardHeaders(thriftbp.HeadersToForward)
	server.SetLogger(thrift.Logger(logger))
	server.Listen()
	// handler, err := transport.NewThriftHandler(*serverEndpoint, time.Second*10, logger,
	// 	endpoints.MakeEndpoints(service.New(service.Config{
	// 		Secrets: secretsStore,
	// 		Redis: redisbp.NewMonitoredClusterFactory(
	// 			"redis",
	// 			redis.NewClusterClient(&redis.ClusterOptions{
	// 				Addrs: cfg.Redis.Endpoints,
	// 			}),
	// 		),
	// 	})),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Fatal(handler.Serve())
}

func parseConfig(path string) (config, error) {
	var cfg config
	if path == "" {
		return cfg, errors.New("No config path given.")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	log.Debugf("%#v", cfg)
	return cfg, nil
}

func initLogger(debug bool) (log.Level, log.Wrapper) {
	var logLevel log.Level
	if debug {
		logLevel = log.DebugLevel
	} else {
		logLevel = log.WarnLevel
	}
	log.InitLogger(logLevel)
	return logLevel, log.ZapWrapper(logLevel)
}

func initSecrets(ctx context.Context, cfg config, logger log.Wrapper) (*secrets.Store, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	secretsStore, err := secrets.NewStore(ctx, cfg.Secrets.Path, logger)
	if err != nil {
		return nil, err
	}
	return secretsStore, nil
}

func initTracing(cfg config, logger log.Wrapper, metrics *metricsbp.Statsd) error {
	if err := tracing.InitGlobalTracer(tracing.TracerConfig{
		ServiceName:      cfg.Tracing.Namespace,
		SampleRate:       cfg.Tracing.SampleRate,
		Logger:           logger,
		MaxRecordTimeout: cfg.Tracing.RecordTimeout,
		QueueName:        "main",
	}); err != nil {
		return err
	}

	tracing.RegisterCreateServerSpanHooks(
		metricsbp.CreateServerSpanHook{Metrics: metrics},
		tracing.ErrorReporterCreateServerSpanHook{},
	)
	return nil
}

func initSentry(cfg config) error {
	if err := raven.SetDSN(cfg.Sentry.DSN); err != nil {
		return err
	}
	if err := raven.SetSampleRate(float32(cfg.Sentry.SampleRate)); err != nil {
		return err
	}
	raven.SetEnvironment(cfg.Sentry.Environment)
	return nil
}
