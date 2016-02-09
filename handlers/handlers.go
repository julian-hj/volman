package handlers

import (
	"net/http"

	"github.com/cloudfoundry-incubator/cf_http"
	"github.com/cloudfoundry-incubator/volman/delegate"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"
)

var Routes = rata.Routes{
	{Path: "/v1/drivers", Method: "GET", Name: "drivers"},
}

func New(logger lager.Logger) (http.Handler, error) {

	var handlers = rata.Handlers{
		"drivers": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			client := delegate.NewLocalClient()
			drivers, _ := client.ListDrivers(logger)
			cf_http.WriteJSONResponse(w, http.StatusOK, drivers)
		}),
	}

	return rata.NewRouter(Routes, handlers)
}
