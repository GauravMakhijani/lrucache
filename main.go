package main

import (
	"fmt"
	"log"
	app "lrucache/internal"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type errorResponse struct {
	Error string `json:"error"`
}

func main() {

	// create a new router
	r := mux.NewRouter()

	// create a new cache service
	cacheService := app.ProvideNewCache()

	// create a new websocket server
	server := newWebsocketServer(cacheService)

	// handle websocket connection
	r.HandleFunc("/ws", server.websocketHandler)

	// handle http requests
	r.HandleFunc("/cache/{key}", GetValueWithKeyHandler(cacheService)).Methods(http.MethodGet)
	r.HandleFunc("/", InsertValueHandler(cacheService)).Methods(http.MethodPost)
	r.HandleFunc("/initialize", InitializeCacheHandler(cacheService)).Methods(http.MethodPost)
	r.HandleFunc("/cache/{key}", DeleteKeyHandler(cacheService)).Methods(http.MethodDelete)
	r.HandleFunc("/capacity", GetCacheCapacityHandler(cacheService)).Methods(http.MethodGet)

	// cors handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
	})

	handler := c.Handler(r)

	fmt.Println("server started at port :3000")
	log.Fatal(http.ListenAndServe(":3001", handler))

}
