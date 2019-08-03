package rest_api

import (
	"net/http"

	"github.com/go-pkgz/rest"
)

func (srv RestAPIServer) FontNames(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	names, err := srv.storage.Names()
	if err != nil {
		responseBadRequest(response, request, err)
	} else {
		rest.RenderJSON(response, request, rest.JSON{"fonts": names})
	}
}
