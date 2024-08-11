package main

import (
	"encoding/json"
	"errors"
	"log"
	app "lrucache/internal"
	"net/http"

	"github.com/gorilla/mux"
)

func GetValueWithKeyHandler(cacheService app.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key := vars["key"]

		resp, err := cacheService.GetKeyValue(key)
		if err != nil {
			log.Println("error getting key from cache - ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusOK)
	}

}

func GetCacheCapacityHandler(cacheService app.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := cacheService.GetCacheCapacity()
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusOK)
	}
}

func InsertValueHandler(cacheService app.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the key and value from request body

		var input app.CacheItem
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println("error in decoding request body - ", err)
			err = errors.New("error bad request")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
			return
		}

		resp, err := cacheService.Insert(input)
		if err != nil {
			log.Println("error inserting key in cache - ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusCreated)
	}
}

func InitializeCacheHandler(cacheService app.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("initializing cache")
		var input app.InitializeCacheInput
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println("error in decoding request body - ", err)
			err = errors.New("error bad request")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
			return
		}

		log.Println("input capacity - ", input.Capacity)

		err = cacheService.InitializeCache(input)
		if err != nil {
			log.Println("error inserting key in cache - ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
			return
		}

		log.Println("cache initialized")
		w.WriteHeader(http.StatusAccepted)
	}
}

func DeleteKeyHandler(cacheService app.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		cacheService.RemoveFromCache(key)
		w.WriteHeader(http.StatusOK)
	}
}
