package model

import (
	"testing"
)

func TestNewDeck(t *testing.T) {
	t.Run("Default Deck", func(t *testing.T) {
		deck := NewDeck(false, "")
		assertDeckProperties(t, deck, 52, false)
	})

	t.Run("Shuffled Deck", func(t *testing.T) {
		deck := NewDeck(true, "")
		assertDeckProperties(t, deck, 52, true)
	})

	t.Run("Partial Deck", func(t *testing.T) {
		deck := NewDeck(false, "AS,KD,AC,2C,KH")
		assertDeckProperties(t, deck, 5, false)
	})
}

func TestDrawCards(t *testing.T) {
	t.Run("Draw Valid Cards", func(t *testing.T) {
		deck := NewDeck(false, "")
		drawnCards, success := deck.DrawCards(3)

		if !success {
			t.Errorf("DrawCards returned unexpected success status: got %v want true", success)
		}

		if len(drawnCards) != 3 {
			t.Errorf("DrawCards returned wrong number of cards: got %v want 3", len(drawnCards))
		}

		assertDeckProperties(t, deck, 49, false)
	})

	t.Run("Draw Too Many Cards", func(t *testing.T) {
		deck := NewDeck(false, "")
		drawnCards, success := deck.DrawCards(60)

		if success {
			t.Errorf("DrawCards returned unexpected success status: got %v want false", success)
		}

		if drawnCards != nil {
			t.Errorf("DrawCards returned unexpected cards: got %v want nil", drawnCards)
		}

		assertDeckProperties(t, deck, 52, false)
	})
}

func assertDeckProperties(t *testing.T, deck Deck, remaining int, shuffled bool) {
	t.Helper()

	if deck.Remaining != remaining {
		t.Errorf("Unexpected remaining cards: got %v want %v", deck.Remaining, remaining)
	}

	if deck.Shuffled != shuffled {
		t.Errorf("Unexpected shuffled status: got %v want %v", deck.Shuffled, shuffled)
	}

	if shuffled && !isShuffled(deck.Cards) {
		t.Errorf("Deck is marked as shuffled but cards are not shuffled")
	}
}

func isShuffled(cards []Card) bool {
	for i := 1; i < len(cards); i++ {
		if cards[i-1].Code == cards[i].Code {
			return false
		}
	}
	return true
}
