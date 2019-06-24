package gateway

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/market/config"
	"github.com/tryffel/market/modules"
	"github.com/tryffel/market/modules/Error"
	"github.com/tryffel/market/modules/middleware"
	"github.com/tryffel/market/modules/request"
	"github.com/tryffel/market/modules/response"
	"github.com/tryffel/market/modules/util"
	"github.com/tryffel/market/storage"
	"net/http"
	"time"
)

type Api struct {
	modules.Task
	store         *storage.Store
	server        *http.Server
	privateRouter *mux.Router
	publicRouter  *mux.Router
	auth          middleware.Auth
}

func NewApi(conf *config.Config, store *storage.Store) (modules.Tasker, error) {
	a := &Api{
		store: store,
	}
	a.Name = "Api"
	a.Loop = a.loop
	a.Init()

	a.publicRouter = mux.NewRouter().PathPrefix(config.ApiV1BasePath).Subrouter()
	a.privateRouter = a.publicRouter.PathPrefix("").Subrouter()

	a.server = &http.Server{
		Handler:      a.publicRouter,
		Addr:         conf.Api.ListenTo,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	var err error
	a.auth, err = middleware.NewAuth(conf, store)
	if err != nil {
		err = Error.Wrap(&err, "failed to initialize auth middleware")
		return a, err
	}

	a.addRoutes()
	a.Initialized = true
	return a, nil
}

func (a *Api) addRoutes() {

	a.privateRouter.Use(a.auth.Authorize)
	a.publicRouter.HandleFunc("/version", a.version)

	a.privateRouter.HandleFunc("/secret/version", a.version)

}

func (a *Api) loop() {
	go a.stopServer()

	err := a.server.ListenAndServe()
	logrus.Info("Api server stopped")
	if err != nil {
		if err.Error() != "http: Server closed" {
			err = Error.Wrap(&err, "failed to run api server")
			logrus.Error(err)
		}
	}
}

func (a *Api) stopServer() {
	<-a.ChanStop
	ctx := context.Background()
	err := a.server.Shutdown(ctx)
	if err != nil {
		err = Error.Wrap(&err, "failed to stop api server")
		logrus.Error(err)
	}
}

func (a *Api) version(w http.ResponseWriter, r *http.Request) {
	resp := response.NewHttp(w)
	req := request.NewHttp(r)
	util.GetServerInfo(req, resp)
}
