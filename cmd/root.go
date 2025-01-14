package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	_ "github.com/ihezebin/oneness"
	"github.com/ihezebin/oneness/logger"
	"github.com/ihezebin/promisys/component/cache"
	"github.com/ihezebin/promisys/component/email"
	"github.com/ihezebin/promisys/component/oss"
	"github.com/ihezebin/promisys/component/pubsub"
	"github.com/ihezebin/promisys/component/storage"
	"github.com/ihezebin/promisys/config"
	"github.com/ihezebin/promisys/cron"
	"github.com/ihezebin/promisys/domain/repository"
	"github.com/ihezebin/promisys/domain/service"
	"github.com/ihezebin/promisys/server"
	"github.com/ihezebin/promisys/worker"
	"github.com/ihezebin/promisys/worker/example"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	configPath string
)

func Run(ctx context.Context) error {

	app := &cli.App{
		Name:    "promisys",
		Version: "v1.0.1",
		Usage:   "Rapid construction template of Web service based on DDD architecture",
		Authors: []*cli.Author{
			{Name: "hezebin", Email: "ihezebin@qq.com"},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Destination: &configPath,
				Name:        "config", Aliases: []string{"c"},
				Value: "./config/config.toml",
				Usage: "config file path (default find file from pwd and exec dir",
			},
		},
		Before: func(c *cli.Context) error {
			if configPath == "" {
				return errors.New("config path is empty")
			}

			conf, err := config.Load(configPath)
			if err != nil {
				return errors.Wrapf(err, "load config error, path: %s", configPath)
			}
			logger.Debugf(ctx, "load config: %+v", conf.String())

			if err = initComponents(ctx, conf); err != nil {
				return errors.Wrap(err, "init components error")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			worker.Register(example.NewExampleWorker())
			worker.Run(ctx)
			defer worker.Wait(ctx)

			if err := cron.Run(ctx); err != nil {
				logger.WithError(err).Fatalf(ctx, "cron run error")
			}

			if err := server.Run(ctx, config.GetConfig().Port); err != nil {
				logger.WithError(err).Fatalf(ctx, "server run error, port: %d", config.GetConfig().Port)
			}

			return nil
		},
	}

	return app.Run(os.Args)
}

func initComponents(ctx context.Context, conf *config.Config) error {
	// init logger
	if conf.Logger != nil {
		logger.ResetLoggerWithOptions(
			logger.WithServiceName(conf.ServiceName),
			logger.WithPrettyCallerHook(),
			logger.WithTimestampHook(),
			logger.WithLevel(conf.Logger.Level),
			//logger.WithLocalFsHook(filepath.Join(conf.Pwd, conf.Logger.Filename)),
			// 每天切割，保留 3 天的日志
			logger.WithRotateLogsHook(filepath.Join(conf.Pwd, conf.Logger.Filename), time.Hour*24, time.Hour*24*3),
		)
	}

	// init storage
	if conf.MysqlDsn != "" {
		if err := storage.InitMySQLStorageClient(ctx, conf.MysqlDsn); err != nil {
			return errors.Wrap(err, "init mysql storage client error")
		}
	}
	if conf.MongoDsn != "" {
		if err := storage.InitMongoStorageClient(ctx, conf.MongoDsn); err != nil {
			return errors.Wrap(err, "init mongo storage client error")
		}
	}

	// init oss
	if conf.OSSDsn != "" {
		if err := oss.Init(conf.OSSDsn); err != nil {
			return errors.Wrap(err, "init oss client error")
		}
	}

	// init cache
	cache.InitMemoryCache(time.Minute*5, time.Minute)
	if conf.Redis != nil {
		if err := cache.InitRedisCache(ctx, conf.Redis.Addrs, conf.Redis.Password); err != nil {
			return errors.Wrap(err, "init redis cache client error")
		}
	}

	// init repository
	if conf.MysqlDsn != "" || conf.MongoDsn != "" {
		repository.Init()
	}

	// init pubsub
	if conf.Pulsar != nil {
		if err := pubsub.InitPulsarClient(conf.Pulsar.Url); err != nil {
			return errors.Wrap(err, "init pulsar client error")
		}
	}
	if conf.Kafka != nil {
		if err := pubsub.InitKafkaConn(ctx, conf.Kafka.Address, conf.Kafka.Topic, conf.Kafka.Partition); err != nil {
			return errors.Wrap(err, "init kafka client error")
		}
	}

	// init email
	if conf.Email != nil {
		if err := email.Init(*conf.Email); err != nil {
			return errors.Wrap(err, "init email client error")
		}
	}

	// init service
	service.Init()

	return nil
}
