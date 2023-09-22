package issuer

import (
	"github.com/pygzfei/issuer-gateway/issuer/internal/config"
	"github.com/pygzfei/issuer-gateway/issuer/internal/database"
	"github.com/pygzfei/issuer-gateway/issuer/internal/errs"
	"github.com/pygzfei/issuer-gateway/issuer/internal/handler"
	"github.com/pygzfei/issuer-gateway/issuer/internal/svc"
	"github.com/pygzfei/issuer-gateway/issuer/middleware"
	"github.com/pygzfei/issuer-gateway/pkg/acme"
	"github.com/pygzfei/issuer-gateway/pkg/logger"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	xhttp "github.com/zeromicro/x/http"
	"net/http"
)

func Run(conPath string) {

	var c config.Config
	conf.MustLoad(conPath, &c)

	logger.Init(c.Env)
	database.Init(&c)

	server := rest.MustNewServer(c.Issuer.RestConf,
		middleware.Cors(),
		rest.WithNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			xhttp.JsonBaseResponseCtx(r.Context(), w, errs.NotFoundException)
		})),
	)

	defer server.Stop()

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/.well-known/acme-challenge/:token",
		Handler: acme.AcceptChallenge(),
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logx.Infof("Starting server at %s:%d...", c.Issuer.Host, c.Issuer.Port)

	server.Start()
}