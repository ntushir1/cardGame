package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"cardGame/deck/api"
	"cardGame/deck/dao"
	"cardGame/deck/service"
)

func TestHealthCheckHandler(t *testing.T) {
	deckStorage := dao.NewDeckStorage()
	deckService := service.NewDeckService(deckStorage)
	deckHandler := api.NewDeckHandler(deckService, deckStorage)

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err, "Error creating health check request")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deckHandler.HealthCheck)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	assert.Equal(t, "OK", rr.Body.String(), "Expected body to be 'OK'")
}

func TestCreateDeckHandler(t *testing.T) {
	deckStorage := dao.NewDeckStorage()
	deckService := service.NewDeckService(deckStorage)
	deckHandler := api.NewDeckHandler(deckService, deckStorage)

	req, err := http.NewRequest("GET", "/deck", nil)
	assert.NoError(t, err, "Error creating create deck request")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deckHandler.CreateDeck)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Expected Content-Type header to be 'application/json'")

}

func TestDrawCardsHandler(t *testing.T) {
	deckStorage := dao.NewDeckStorage()
	deckService := service.NewDeckService(deckStorage)
	deckHandler := api.NewDeckHandler(deckService, deckStorage)

	req, err := http.NewRequest("GET", "/deck/{deckID}/draw", nil)
	assert.NoError(t, err, "Error creating draw cards request")

	deckID := "your-deck-id"
	req = mux.SetURLVars(req, map[string]string{"deckID": deckID})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deckHandler.DrawCards)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected BAD Request")
}
