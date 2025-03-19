package index

import (
	"github.com/xinliangnote/go-gin-api/skitii/pkg/core"
	"github.com/xinliangnote/go-gin-api/skitii/repository/mysql"
	"github.com/xinliangnote/go-gin-api/skitii/repository/redis"

	"go.uber.org/zap"
)

type handler struct {
	logger *zap.Logger
	cache  redis.Repo
	db     mysql.Repo
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) *handler {
	return &handler{
		logger: logger,
		cache:  cache,
		db:     db,
	}
}

func (h *handler) Index() core.HandlerFunc {
	return func(ctx core.Context) {
		ctx.HTML("index", nil)
	}
}
