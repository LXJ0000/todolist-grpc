package bootstrap

import (
	"github.com/LXJ0000/todolist-grpc-gateway/pkg/log"
	"github.com/LXJ0000/todolist-grpc-gateway/pkg/snowflake"
	"github.com/LXJ0000/todolist-grpc-gateway/pkg/cache"
	"github.com/LXJ0000/todolist-grpc-gateway/pkg/orm"
)

type Application struct {
	Env *Env
	//Mongo mongo.Client
	Orm   orm.Database
	Cache cache.RedisCache
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	//app.Mongo = NewMongoDatabase(app.Env)
	app.Orm = NewOrmDatabase(app.Env)
	app.Cache = NewRedisCache(app.Env)
	logutil.Init(app.Env.AppEnv)
	snowflakeutil.Init(app.Env.SnowflakeStartTime, app.Env.SnowflakeMachineID)

	return *app
}
