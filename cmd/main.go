package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"poi-service/cmd/auth"
	"poi-service/cmd/data"
	"poi-service/cmd/download"
	"poi-service/cmd/handler"
	"time"
)

// used components
var (
	poiHandler         handler.PoiHandler
	dbHandler          handler.DbHandler
	credentialUsername string
	credentialPw       string
	authorizer         auth.Authorizer
	jwkStore           auth.JwkStore
	httpClient         download.HttpRequester
)

// init is the reserved golang function that will initialize all components once.
func init() {
	mongodb := os.Getenv("DATABASE_URL")
	dbHandler, err := handler.NewDbHandler(mongodb)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create dbHandler")
		return
	}
	poiHandler = handler.NewPoiHandler(dbHandler)
	httpClient = download.NewHttpRequester(http.DefaultClient)
	jwkCache := auth.JwkCache{}
	jwkCache.Init()
	jwkStore = auth.NewJwkStore("http://127.0.0.1:4444/", httpClient, jwkCache)
	authorizer = auth.NewAuthorizer(jwkStore)

	if poiHandler == nil {
		log.Fatal().Msg("poiHandler is nil")
		return
	}

	if authorizer == nil {
		log.Fatal().Msg("authorizer is nil")
		return
	}
}

func main() {
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, os.Interrupt)
	res := make(chan error, 1)
	defer close(res)

	port := os.Getenv("SERVICE_PORT")
	log.Printf("Listening in port %s", port)

	s := http.Server{
		Addr:    ":" + port,
		Handler: createRootHandler(),
	}
	go func() {
		res <- s.ListenAndServe()
	}()

	select {
	case <-quit:
		log.Info().Msg("user initiated termination of server started")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := s.Shutdown(ctx)
		if err != nil {
			log.Error().Err(err).Msg("graceful shutdown failed")
		}
	case err := <-res:
		log.Error().Err(err).Msg("server stopped with error")
	}
}

func createRootHandler() http.Handler {

	r := mux.NewRouter()
	api := r.PathPrefix("/v1").Subrouter()
	api.Use(authorizer.Authorize)
	api.HandleFunc("/pois/{id}", getPoi).Methods(http.MethodGet)
	api.HandleFunc("/pois", createPoi).Methods(http.MethodPost)
	api.HandleFunc("/pois/{id}", deletePoi).Methods(http.MethodDelete)
	api.HandleFunc("/pois/{id}", updatePoi).Methods(http.MethodPut)
	api.HandleFunc("/pois/list", listPoi).Methods(http.MethodPost)
	return r
}

func createPoi(rw http.ResponseWriter, r *http.Request) {
	var poi data.Poi
	if err := decode(r, &poi); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := poiHandler.Create(&poi)
	if err != nil {
		log.Warn().Err(err).Msg("createPoi failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := encode(rw, &id); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	return
}

func updatePoi(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := params["id"]; !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Warn().Msg("id not available in path")
		return
	}

	var poi data.Poi
	if err := decode(r, &poi); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := poiHandler.Update(data.Id(params["id"]), &poi); err != nil {
		log.Warn().Err(err).Msg("updatePoi failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	return
}

func getPoi(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := params["id"]; !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Warn().Msg("getPoi id not available in path")
		return
	}

	resp, err := poiHandler.Get(data.Id(params["id"]))
	if err != nil {
		log.Warn().Err(err).Msg("getPoi failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := encode(rw, &resp); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	return
}

func listPoi(rw http.ResponseWriter, r *http.Request) {
	var area data.SearchArea

	// if provided set a search area
	if r.ContentLength > 0 {
		if err := decode(r, &area); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	resp, err := poiHandler.Search(area)
	if err != nil {
		log.Warn().Err(err).Msg("listPoi failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := encode(rw, &resp); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	return
}

func deletePoi(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := params["id"]; !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Warn().Msg("id not available in path")
		return
	}

	if err := poiHandler.Delete(data.Id(params["id"])); err != nil {
		log.Warn().Err(err).Msg("deletePoi failed")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	return
}

func decode(r *http.Request, poi interface{}) (err error) {
	if poi == nil {
		return errors.New("is nil")
	}
	err = json.NewDecoder(r.Body).Decode(poi)
	if err != nil {
		log.Warn().Err(err).Msg("json decoding failed:")
	}
	return
}

func encode(rw http.ResponseWriter, poi interface{}) (err error) {
	if poi == nil {
		return errors.New("poi is nil")
	}
	err = json.NewEncoder(rw).Encode(poi)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	return
}
