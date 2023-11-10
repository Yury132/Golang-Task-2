package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	formatJSON = "json"
)

type Config struct {
	Server struct {
		Host        string `envconfig:"SERVER_HOST" default:":9000"`
		MetricsBind string `envconfig:"BIND_METRICS" default:":9090"`
		HealthHost  string `envconfig:"BIND_HEALTH" default:":9091"`
	}

	Service struct {
		LogLevel  string `envconfig:"LOGGER_LEVEL" default:"debug"`
		LogFormat string `envconfig:"LOGGER_FORMAT" default:"console"`
	}

	DB struct {
		Address  string `envconfig:"DB_ADDRESS" default:"localhost"`
		Name     string `envconfig:"DB_NAME" default:"mydb"`
		User     string `envconfig:"DB_USER" default:"root"`
		Password string `envconfig:"DB_PASSWORD" default:"mydbpass"`
		Port     int    `envconfig:"DB_PORT" default:"5432"`
		MaxConn  int    `envconfig:"DB_MAX_CONN" default:"15"`
	}

	NATS struct {
		URL string `envconfig:"NATS_URL" default:"nats://localhost:4222"`
	}
}

func Parse() (*Config, error) {
	var cfg = new(Config)
	// Устанавливаем значения переменных окружения
	err := envconfig.Process("", cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Логгер
func (cfg Config) Logger() (logger zerolog.Logger) {
	level := zerolog.InfoLevel
	if newLevel, err := zerolog.ParseLevel(cfg.Service.LogLevel); err == nil {
		level = newLevel
	}

	var out io.Writer = os.Stdout
	if cfg.Service.LogFormat != formatJSON {
		out = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return zerolog.New(out).Level(level).With().Caller().Timestamp().Logger()
}

// Получаем адрес в БД
func (cfg Config) GetDBConnString() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s sslmode=disable user=%s password=%s",
		cfg.DB.Address, cfg.DB.Port, cfg.DB.Name, cfg.DB.User, cfg.DB.Password,
	)
}

// func (cfg Config) PgPoolConfig() (*pgxpool.Config, error) {
// 	poolCfg, err := pgxpool.ParseConfig(fmt.Sprintf(
// 		"host=%s port=%d dbname=%s sslmode=disable user=%s password=%s pool_max_conns=%d",
// 		cfg.DB.Address, cfg.DB.Port, cfg.DB.Name, cfg.DB.User, cfg.DB.Password, cfg.DB.MaxConn,
// 	))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return poolCfg, nil
// }

// Конфигурация подключения к БД
func (cfg Config) PgPoolConfig() (*pgxpool.Config, error) {
	poolCfg, err := pgxpool.ParseConfig(fmt.Sprintf("%s pool_max_conns=%d", cfg.GetDBConnString(), cfg.DB.MaxConn))
	if err != nil {
		return nil, err
	}

	return poolCfg, nil
}

// Стрим для nats
func (cfg Config) NewStream(ctx context.Context, js jetstream.JetStream) (jetstream.Stream, error) {
	streamCfg := jetstream.StreamConfig{
		Name: "EVENTS",
		// Очередь
		Retention: jetstream.WorkQueuePolicy,
		// Топики
		Subjects: []string{"media.>"},
	}

	stream, err := js.CreateStream(ctx, streamCfg)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
