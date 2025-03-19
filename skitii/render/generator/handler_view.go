package generator_handler

import (
	"github.com/xinliangnote/go-gin-api/skitii/pkg/core"
)

func (h *handler) HandlerView() core.HandlerFunc {
	return func(c core.Context) {
		c.HTML("generator_handler", nil)
	}
}
