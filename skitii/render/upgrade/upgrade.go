package upgrade

import (
	"github.com/xinliangnote/go-gin-api/skitii/repository/mysql"
	"github.com/xinliangnote/go-gin-api/skitii/repository/redis"

	"go.uber.org/zap"
)

type handler struct {
	db     mysql.Repo
	logger *zap.Logger
	cache  redis.Repo
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) *handler {
	return &handler{
		logger: logger,
		cache:  cache,
		db:     db,
	}
}
