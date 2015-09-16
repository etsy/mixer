package main

import (
	"log"
	"net/http"
	"strconv"

	. "github.com/etsy/mixer/config"

	"github.com/etsy/mixer/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/etsy/mixer/handlers"
)

var router = mux.NewRouter()

func main() {
	err := Config.Load()
	if err != nil {
		log.Fatal("error reading or parsing config:", err)
	}

	h, _ := handlers.NewHandlers()

	router.NotFoundHandler = http.HandlerFunc(h.NotFoundHandler)

	router.HandleFunc("/mixer", h.IndexHandler)
	router.HandleFunc("/mixers", h.MixerHandler).Methods("GET")
	router.HandleFunc("/group/{group}", h.AllPeopleHandler).Methods("GET")
	router.HandleFunc("/people/{id:[0-9]+}", h.PersonHandler).Methods("GET")
	router.HandleFunc("/people/{id:[0-9]+}", h.PersonPutHandler).Methods("PUT")
	router.HandleFunc("/people", h.PersonPostHandler).Methods("POST")
	router.HandleFunc("/authuser", h.AuthHandler)
	router.HandleFunc("/redirect/staff/{username}", h.StaffRedirectHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", router)

	serverAddr := "0.0.0.0:" + strconv.Itoa(Config.Server.Port)
	log.Println("Listening on " + serverAddr)
	err = http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("Error starting web server:", err)
	}
}
