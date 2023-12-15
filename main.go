package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"cardGame/deck/api"
	"cardGame/deck/dao"
	"cardGame/deck/service"
)

func main() {
	deckStorage := dao.NewDeckStorage()
	deckService := service.NewDeckService(deckStorage)
	deckHandler := api.NewDeckHandler(deckService, deckStorage)

	router := configureRoutes(deckHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func configureRoutes(deckHandler *api.DeckHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", deckHandler.HealthCheck).Methods("GET")
	router.HandleFunc("/deck", deckHandler.CreateDeck).Methods("GET")
	router.HandleFunc("/deck/{deckID}/draw", deckHandler.DrawCards).Methods("GET")

	return router
}
