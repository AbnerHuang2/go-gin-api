package handler

import (
	"github.com/xinliangnote/go-gin-api/cmd"
	"github.com/xinliangnote/go-gin-api/internal/router"
	"net/http"
)

var (
	s *router.Server
)

func init() {
	s = cmd.Start()
}

/**

vercel 入口

所有的路由都从里进入，相当于平时的 main 函数
*/

func Handler(w http.ResponseWriter, r *http.Request) {

	s.Mux.ServeHTTP(w, r)
}
