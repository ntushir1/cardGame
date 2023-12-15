package dao

import (
	"cardGame/deck/model"
	"github.com/google/uuid"
	"sync"
)

type DeckStorage struct {
	mu    sync.Mutex
	decks map[uuid.UUID]model.Deck
}

func NewDeckStorage() *DeckStorage {
	return &DeckStorage{
		decks: make(map[uuid.UUID]model.Deck),
	}
}

func (s *DeckStorage) SaveDeck(deck model.Deck) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.decks[deck.ID] = deck
}

func (s *DeckStorage) GetDeck(deckID uuid.UUID) (model.Deck, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	deck, ok := s.decks[deckID]
	return deck, ok
}
