package handlers

import (
	"net/http"

	"github.com/shivanshkc/ledgerkeep/src/configs"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"
)

// BasicHandler serves the basic information about the application.
func BasicHandler(writer http.ResponseWriter, request *http.Request) {
	conf := configs.Get()

	responseBody := &httputils.ResponseBodyDTO{
		StatusCode: http.StatusOK,
		CustomCode: "OK",
		Data: map[string]interface{}{
			"name":    conf.Application.Name,
			"version": conf.Application.Version,
		},
	}

	response := &httputils.ResponseDTO{Status: http.StatusOK, Body: responseBody}
	httputils.WriteAndLog(request.Context(), writer, response, logger.Get())
}
