package driverhttp

import (
	"net/http"

	cf_http_handlers "github.com/cloudfoundry-incubator/cf_http/handlers"
	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"
)

// At present, Docker ignores HTTP status codes, and requires errors to be returned in the response body.  To
// comply with this API, we will return 200 in all cases
const (
	statusInternalServerError = http.StatusOK
	statusOK                  = http.StatusOK
)

func NewHandler(logger lager.Logger, client voldriver.Driver) (http.Handler, error) {
	return rata.NewRouter(voldriver.Routes, rata.Handlers{
		voldriver.ActivateRoute: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger := logger.Session("activate")
			logger.Info("start")
			defer logger.Info("end")

			activateResponse := client.Activate(logger)
			// ok to eat error as we should be removing error from the Info func signature
			if activateResponse.Err != "" {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, activateResponse)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, statusOK, activateResponse)
		}),

		voldriver.GetRoute: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger := logger.Session("get-mount")
			logger.Info("start")
			defer logger.Info("end")

			var getRequest voldriver.GetRequest
			err := unmarshallJSON(logger, req.Body, &getRequest)
			if err != nil {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, voldriver.GetResponse{Err: err.Error()})
				return
			}

			getResponse := client.Get(logger, getRequest)
			if getResponse.Err != "" {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, getResponse)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, statusOK, getResponse)
		}),

		voldriver.CreateRoute: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger := logger.Session("handle-create")
			logger.Info("start")
			defer logger.Info("end")

			var createRequest voldriver.CreateRequest
			err := unmarshallJSON(logger, req.Body, &createRequest)
			if err != nil {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, voldriver.ErrorResponse{Err: err.Error()})
				return
			}

			createResponse := client.Create(logger, createRequest)
			if createResponse.Err != "" {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, createResponse)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, statusOK, createResponse)
		}),

		voldriver.MountRoute: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger := logger.Session("handle-mount")
			logger.Info("start")
			defer logger.Info("end")

			var mountRequest voldriver.MountRequest
			err := unmarshallJSON(logger, req.Body, &mountRequest)
			if err != nil {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, voldriver.MountResponse{Err: err.Error()})
				return
			}

			mountResponse := client.Mount(logger, mountRequest)
			if mountResponse.Err != "" {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, mountResponse)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, statusOK, mountResponse)
		}),

		voldriver.UnmountRoute: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger := logger.Session("handle-unmount")
			logger.Info("start")
			defer logger.Info("end")

			var unmountRequest voldriver.UnmountRequest
			err := unmarshallJSON(logger, req.Body, &unmountRequest)
			if err != nil {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, voldriver.ErrorResponse{Err: err.Error()})
				return
			}

			unmountResponse := client.Unmount(logger, unmountRequest)
			if unmountResponse.Err != "" {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, unmountResponse)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, statusOK, unmountResponse)
		}),

		voldriver.RemoveRoute: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger := logger.Session("handle-remove")
			logger.Info("start")
			defer logger.Info("end")

			var removeRequest voldriver.RemoveRequest
			err := unmarshallJSON(logger, req.Body, &removeRequest)
			if err != nil {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, voldriver.ErrorResponse{Err: err.Error()})
				return
			}

			removeResponse := client.Remove(logger, removeRequest)
			if removeResponse.Err != "" {
				cf_http_handlers.WriteJSONResponse(w, statusInternalServerError, removeResponse)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, statusOK, removeResponse)
		}),
	})
}
