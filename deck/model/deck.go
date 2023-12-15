package model

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

type Deck struct {
	ID        uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	Cards     []Card    `json:"cards"`
}

func NewDeck(shuffled bool, cards string) Deck {
	var deckID uuid.UUID
	deckID, _ = uuid.NewUUID()

	var allCards []Card
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING", "ACE"}
	suits := []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}

	for _, suit := range suits {
		for _, value := range values {
			card := Card{Value: value, Suit: suit, Code: value[:1] + strings.ToUpper(suit[:1])}
			allCards = append(allCards, card)
		}
	}

	if cards != "" {
		allCards = filterDeck(allCards, cards)
	}

	if shuffled {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(allCards), func(i, j int) {
			allCards[i], allCards[j] = allCards[j], allCards[i]
		})
	}

	return Deck{
		ID:        deckID,
		Shuffled:  shuffled,
		Remaining: len(allCards),
		Cards:     allCards,
	}
}

func filterDeck(allCards []Card, cards string) []Card {
	cardCodes := strings.Split(cards, ",")
	var filteredDeck []Card

	for _, code := range cardCodes {
		for _, card := range allCards {
			if card.Code == strings.ToUpper(code) {
				filteredDeck = append(filteredDeck, card)
				break
			}
		}
	}

	return filteredDeck
}

func (d *Deck) DrawCards(count int) ([]Card, bool) {
	if count > d.Remaining {
		return nil, false
	}

	drawnCards := d.Cards[:count]
	d.Cards = d.Cards[count:]
	d.Remaining -= count

	return drawnCards, true
}
