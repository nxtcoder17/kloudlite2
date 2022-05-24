package framework

import (
	"fmt"
	"go.uber.org/fx"
	"kloudlite.io/apps/console/internal/app"
	"kloudlite.io/pkg/cache"
	"kloudlite.io/pkg/config"
	rpc "kloudlite.io/pkg/grpc"
	httpServer "kloudlite.io/pkg/http-server"
	"kloudlite.io/pkg/logger"
	loki_server "kloudlite.io/pkg/loki-server"
	"kloudlite.io/pkg/messaging"
	mongo_db "kloudlite.io/pkg/repos"
	rcn "kloudlite.io/pkg/res-change-notifier"
)

type GrpcInfraConfig struct {
	InfraGrpcHost string `env:"INFRA_HOST" required:"true"`
	InfraGrpcPort string `env:"INFRA_PORT" required:"true"`
}

func (e *GrpcInfraConfig) GetGCPServerURL() string {
	return e.InfraGrpcHost + ":" + e.InfraGrpcPort
}

type GrpcFinanceConfig struct {
	FinanceGrpcHost string `env:"FINANCE_HOST" required:"true"`
	FinanceGrpcPort string `env:"FINANCE_PORT" required:"true"`
}

func (e *GrpcFinanceConfig) GetGCPServerURL() string {
	return e.FinanceGrpcHost + ":" + e.FinanceGrpcPort
}

type GrpcAuthConfig struct {
	AuthGrpcHost string `env:"AUTH_HOST" required:"true"`
	AuthGrpcPort string `env:"AUTH_PORT" required:"true"`
}

func (e *GrpcAuthConfig) GetGCPServerURL() string {
	return e.AuthGrpcHost + ":" + e.AuthGrpcPort
}

type GrpcCIConfig struct {
	CIGrpcHost string `env:"CI_HOST" required:"true"`
	CIGrpcPort string `env:"CI_PORT" required:"true"`
}

func (e *GrpcCIConfig) GetGCPServerURL() string {
	return e.CIGrpcHost + ":" + e.CIGrpcPort
}

type IAMGRPCEnv struct {
	IAMServerHost string `env:"IAM_SERVER_HOST"`
	IAMServerPort uint16 `env:"IAM_SERVER_PORT"`
}

func (e *IAMGRPCEnv) GetGCPServerURL() string {
	return fmt.Sprintf("%v:%v", e.IAMServerHost, e.IAMServerPort)
}

type LogServerEnv struct {
	LokiServerUrl string `env:"LOKI_URL" required:"true"`
	LogServerPort uint64 `env:"LOG_SERVER_PORT" required:"true"`
}

func (l *LogServerEnv) GetLokiServerUrl() string {
	return l.LokiServerUrl
}

func (l *LogServerEnv) GetLogServerPort() uint64 {
	return l.LogServerPort
}

type Env struct {
	MongoUri      string `env:"MONGO_URI" required:"true"`
	RedisHosts    string `env:"REDIS_HOSTS" required:"true"`
	RedisUserName string `env:"REDIS_USERNAME"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	MongoDbName   string `env:"MONGO_DB_NAME" required:"true"`
	KafkaBrokers  string `env:"KAFKA_BOOTSTRAP_SERVERS" required:"true"`
	Port          uint16 `env:"PORT" required:"true"`
	IsDev         bool   `env:"DEV" default:"false" required:"true"`
	CorsOrigins   string `env:"ORIGINS" required:"true"`
	GrpcPort      uint16 `env:"GRPC_PORT" required:"true"`
	NotifierUrl   string `env:"NOTIFIER_URL" required:"true"`
}

func (e *Env) GetBrokers() string {
	return e.KafkaBrokers
}

func (e *Env) GetHttpPort() uint16 {
	return e.Port
}

func (e *Env) GetHttpCors() string {
	return e.CorsOrigins
}

func (e *Env) RedisOptions() (hosts, username, password string) {
	return e.RedisHosts, e.RedisUserName, e.RedisPassword
}

func (e *Env) GetMongoConfig() (url string, dbName string) {
	return e.MongoUri, e.MongoDbName
}

func (e *Env) GetGRPCPort() uint16 {
	return e.GrpcPort
}

func (e *Env) GetNotifierUrl() string {
	return e.NotifierUrl
}

var Module = fx.Module("framework",
	config.EnvFx[Env](),
	config.EnvFx[LogServerEnv](),
	config.EnvFx[IAMGRPCEnv](),
	config.EnvFx[GrpcInfraConfig](),
	config.EnvFx[GrpcAuthConfig](),
	config.EnvFx[GrpcFinanceConfig](),
	config.EnvFx[GrpcCIConfig](),

	fx.Provide(logger.NewLogger),
	rcn.NewFxResourceChangeNotifier[*Env](),
	rpc.NewGrpcServerFx[*Env](),
	rpc.NewGrpcClientFx[*IAMGRPCEnv, app.IAMClientConnection](),
	rpc.NewGrpcClientFx[*GrpcInfraConfig, app.InfraClientConnection](),
	rpc.NewGrpcClientFx[*GrpcAuthConfig, app.AuthClientConnection](),
	rpc.NewGrpcClientFx[*GrpcCIConfig, app.CIClientConnection](),
	rpc.NewGrpcClientFx[*GrpcFinanceConfig, app.FinanceClientConnection](),
	mongo_db.NewMongoClientFx[*Env](),
	messaging.NewKafkaClientFx[*Env](),
	cache.NewRedisFx[*Env](),
	httpServer.NewHttpServerFx[*Env](),
	loki_server.NewLogServerFx[*LogServerEnv](), // will provide log server and loki client
	app.Module,
)
