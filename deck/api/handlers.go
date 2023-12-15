package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"cardGame/deck/dao"
	"cardGame/deck/service"
)

type DeckHandler struct {
	DeckService *service.DeckService
	DeckStorage *dao.DeckStorage
}

type CreateDeckResponse struct {
	DeckID    uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
}

func (h *DeckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func NewDeckHandler(deckService *service.DeckService, deckStorage *dao.DeckStorage) *DeckHandler {
	return &DeckHandler{
		DeckService: deckService,
		DeckStorage: deckStorage,
	}
}

func (h *DeckHandler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	deckIDParam := r.URL.Query().Get("deckId")

	if deckIDParam != "" {
		h.handleExistingDeck(w, deckIDParam)
		return
	}

	h.handleNewDeck(w, r)
}

func (h *DeckHandler) handleExistingDeck(w http.ResponseWriter, deckIDParam string) {
	deckID, err := uuid.Parse(deckIDParam)
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	existingDeck, found := h.DeckService.GetDeck(deckID)
	if !found {
		http.Error(w, "Deck not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingDeck)
}

func (h *DeckHandler) handleNewDeck(w http.ResponseWriter, r *http.Request) {
	cards := r.URL.Query().Get("cards")
	shuffled, _ := strconv.ParseBool(r.URL.Query().Get("shuffled"))
	newDeck := h.DeckService.CreateDeck(shuffled, cards)
	response := CreateDeckResponse{
		DeckID:    newDeck.ID,
		Shuffled:  newDeck.Shuffled,
		Remaining: newDeck.Remaining,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DeckHandler) DrawCards(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deckID, err := uuid.Parse(vars["deckID"])
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		http.Error(w, "Invalid count parameter", http.StatusBadRequest)
		return
	}

	deck, found := h.DeckStorage.GetDeck(deckID)
	if !found {
		http.Error(w, "Deck not found", http.StatusNotFound)
		return
	}

	drawnCards, err := h.DeckService.DrawCards(deck, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"cards": drawnCards})
}
