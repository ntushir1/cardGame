package service

import (
	"cardGame/deck/dao"
	"cardGame/deck/model"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type DeckService struct {
	mu      sync.Mutex
	storage *dao.DeckStorage
}

func NewDeckService(storage *dao.DeckStorage) *DeckService {
	return &DeckService{
		storage: storage,
	}
}

func (s *DeckService) CreateDeck(shuffled bool, cards string) model.Deck {
	newDeck := model.NewDeck(shuffled, cards)
	s.storage.SaveDeck(newDeck)
	return newDeck
}

func (s *DeckService) GetDeck(deckID uuid.UUID) (model.Deck, bool) {
	deck, err := s.storage.GetDeck(deckID)
	if err != true {
		return model.Deck{}, false
	}
	return deck, true
}

func (s *DeckService) DrawCards(deck model.Deck, count int) ([]model.Card, error) {
	deck, err := s.storage.GetDeck(deck.ID)
	if err != true {
		return nil, fmt.Errorf("Invalid Deck ID")
	}

	drawnCards, err1 := deck.DrawCards(count)
	if err1 != true {
		return nil, fmt.Errorf("Not enough cards remaining in the deck")
	}

	s.storage.SaveDeck(deck)
	return drawnCards, nil
}
