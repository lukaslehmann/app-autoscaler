package server

import (
	"autoscaler/cf"
	"autoscaler/db"
	"autoscaler/scalingengine/config"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"github.com/gorilla/mux"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"

	"fmt"
	"net/http"
)

const (
	PathScale            = "/v1/apps/{appid}/scale"
	PathScalingHistories = "/v1/apps/{appid}/scaling_histories"
	RouteNameScale       = "scale"
	RouteNameHistoreis   = "histories"
)

type VarsFunc func(w http.ResponseWriter, r *http.Request, vars map[string]string)

func (vh VarsFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vh(w, r, vars)
}

func NewServer(logger lager.Logger, conf config.ServerConfig, cfc cf.CfClient, policyDB db.PolicyDB, historyDB db.HistoryDB) ifrit.Runner {
	handler := NewScalingHandler(logger, cfc, policyDB, historyDB, clock.NewClock())

	r := mux.NewRouter()
	r.Methods("POST").Path(PathScale).Handler(VarsFunc(handler.HandleScale)).Name(RouteNameScale)
	r.Methods("GET").Path(PathScalingHistories).Handler(VarsFunc(handler.GetScalingHistories)).Name(RouteNameHistoreis)

	addr := fmt.Sprintf("0.0.0.0:%d", conf.Port)
	return http_server.New(addr, r)
}