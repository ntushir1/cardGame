package dao

import (
	"cardGame/deck/model"
	"github.com/google/uuid"
	"testing"
)

func TestDeckStorage_SaveDeck(t *testing.T) {
	storage := NewDeckStorage()

	deck := model.NewDeck(false, "")
	storage.SaveDeck(deck)

	savedDeck, found := storage.GetDeck(deck.ID)
	if !found {
		t.Errorf("SaveDeck failed: deck not found in storage")
	}
	if savedDeck.ID != deck.ID {
		t.Errorf("SaveDeck failed: expected deck ID %v, got %v", deck.ID, savedDeck.ID)
	}
	if savedDeck.Remaining != deck.Remaining {
		t.Errorf("SaveDeck failed: expected remaining cards %v, got %v", deck.Remaining, savedDeck.Remaining)
	}
}

func TestDeckStorage_GetDeck(t *testing.T) {
	storage := NewDeckStorage()

	_, found := storage.GetDeck(uuid.New())
	if found {
		t.Errorf("GetDeck failed: found a non-existing deck")
	}

	deck := model.NewDeck(false, "")
	storage.SaveDeck(deck)

	savedDeck, found := storage.GetDeck(deck.ID)
	if !found {
		t.Errorf("GetDeck failed: deck not found in storage")
	}
	if savedDeck.ID != deck.ID {
		t.Errorf("GetDeck failed: expected deck ID %v, got %v", deck.ID, savedDeck.ID)
	}
	if savedDeck.Remaining != deck.Remaining {
		t.Errorf("GetDeck failed: expected remaining cards %v, got %v", deck.Remaining, savedDeck.Remaining)
	}
}
