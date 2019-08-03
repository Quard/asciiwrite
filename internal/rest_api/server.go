package rest_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-pkgz/rest"
	"github.com/quard/asciiwrite/internal/storage"
	"github.com/thedevsaddam/govalidator"
)

type Opts struct {
	Host      string `short:"b" long:"bind" env:"HOST" default:"0.0.0.0"`
	Port      int    `short:"p" long:"port" env:"PORT" default:"5000"`
	AuthToken string `short:"t" long:"auth-token"`
}

type RestAPIServer struct {
	opts    Opts
	storage storage.FontStorage
}

func NewRestAPIServer(opts Opts, storage storage.FontStorage) (RestAPIServer, error) {
	srv := RestAPIServer{opts: opts, storage: storage}

	return srv, nil
}

func (srv RestAPIServer) Run() {
	srv.initCustomValidators()

	router := srv.getRouter()

	listenParams := fmt.Sprintf("%s:%d", srv.opts.Host, srv.opts.Port)
	log.Printf("listen: %s", listenParams)
	log.Fatal(http.ListenAndServe(listenParams, router))
}

func (srv RestAPIServer) getRouter() chi.Router {
	router := chi.NewRouter()
	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/print/", srv.Print)
		r.Get("/fonts/", srv.FontNames)

		r.With(srv.authMiddleware).Post("/font/upload/", srv.FontUpload)
		// r.With(srv.authMiddleware).Get("/font/{name}/", srv.GetFont)
		// r.With(srv.authMiddleware).Delete("/font/{name}/", srv.GetFont)
	})

	return router
}

func (srv RestAPIServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		authToken := request.Header.Get("Authorization")
		if authToken != srv.opts.AuthToken {
			response.WriteHeader(http.StatusForbidden)
		} else {
			next.ServeHTTP(response, request)
		}
	})
}

func (srv RestAPIServer) initCustomValidators() {
	govalidator.AddCustomRule("font_not_exists", func(field, rule, message string, value interface{}) error {
		var err error
		fontName := value.(string)
		fontExists, err := srv.storage.IsExist(fontName)
		if err == nil && fontExists {
			return fmt.Errorf("font with name '%s' already exists", fontName)
		}

		return err
	})
}

func responseValidationErrors(response http.ResponseWriter, validationError url.Values) {
	response.WriteHeader(http.StatusBadRequest)
	err := map[string]interface{}{"validationError": validationError}
	json.NewEncoder(response).Encode(err)
}

func responseBadRequest(response http.ResponseWriter, request *http.Request, err error) {
	response.WriteHeader(http.StatusBadRequest)
	rest.RenderJSON(response, request, rest.JSON{"error": err.Error()})
}
