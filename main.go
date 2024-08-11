package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

func GetValueWithKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	spew.Dump(vars)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/value/{key}", GetValueWithKeyHandler).Methods(http.MethodGet)
	fmt.Println("server started at port :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
