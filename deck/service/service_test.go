package service

import (
	"cardGame/deck/dao"
	"cardGame/deck/model"
	"github.com/google/uuid"
	"testing"
)

func TestDeckService_CreateDeck(t *testing.T) {
	storage := dao.NewDeckStorage()
	service := NewDeckService(storage)

	createdDeck := service.CreateDeck(false, "")
	assertDeckProperties(t, createdDeck, 52, false)
}

func TestDeckService_GetDeck(t *testing.T) {
	storage := dao.NewDeckStorage()
	service := NewDeckService(storage)

	_, found := service.GetDeck(uuid.New())
	if found {
		t.Errorf("GetDeck failed: found a non-existing deck")
	}

	createdDeck := service.CreateDeck(false, "")
	retrievedDeck, found := service.GetDeck(createdDeck.ID)
	if !found {
		t.Errorf("GetDeck failed: deck not found")
	}
	assertDeckProperties(t, retrievedDeck, 52, false)
}

func TestDeckService_DrawCards(t *testing.T) {
	storage := dao.NewDeckStorage()
	service := NewDeckService(storage)

	_, err := service.DrawCards(model.Deck{ID: uuid.New()}, 3)
	if err == nil || err.Error() != "Invalid Deck ID" {
		t.Errorf("DrawCards failed: expected 'Invalid Deck ID' error, got %v", err)
	}

	createdDeck := service.CreateDeck(false, "")
	_, err = service.DrawCards(createdDeck, 60)
	if err == nil || err.Error() != "Not enough cards remaining in the deck" {
		t.Errorf("DrawCards failed: expected 'Not enough cards remaining in the deck' error, got %v", err)
	}

	drawnCards, err := service.DrawCards(createdDeck, 0)
	if err != nil {
		t.Errorf("DrawCards failed: unexpected error %v", err)
	}

	assertDeckProperties(t, createdDeck, 52, false)
	if len(drawnCards) != 0 {
		t.Errorf("DrawCards failed: expected 3 drawn cards, got %v", len(drawnCards))
	}
}

func assertDeckProperties(t *testing.T, deck model.Deck, remaining int, shuffled bool) {
	t.Helper()

	if deck.Remaining != remaining {
		t.Errorf("Unexpected remaining cards: got %v want %v", deck.Remaining, remaining)
	}

	if deck.Shuffled != shuffled {
		t.Errorf("Unexpected shuffled status: got %v want %v", deck.Shuffled, shuffled)
	}
}
