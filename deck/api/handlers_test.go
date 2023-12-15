package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cardGame/deck/dao"
	"cardGame/deck/model"
	"cardGame/deck/service"
)

func TestDeckHandler_HealthCheck(t *testing.T) {
	handler := NewDeckHandler(nil, nil)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HealthCheck(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Health check handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("Health check handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeckHandler_CreateDeck(t *testing.T) {
	storage := dao.NewDeckStorage()
	service := service.NewDeckService(storage)
	handler := NewDeckHandler(service, storage)

	t.Run("Create New Deck", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/deck?shuffled=true&cards=AS,KD,AC,2C,KH", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.CreateDeck(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("CreateDeck handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response CreateDeckResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if response.Shuffled != true {
			t.Errorf("CreateDeck handler returned wrong shuffle status: got %v want true", response.Shuffled)
		}

		if response.Remaining != 5 {
			t.Errorf("CreateDeck handler returned wrong remaining cards: got %v want 5", response.Remaining)
		}
	})

	t.Run("Create Existing Deck", func(t *testing.T) {
		existingDeck := service.CreateDeck(false, "")
		req, err := http.NewRequest("GET", "/deck?deckId="+existingDeck.ID.String(), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.CreateDeck(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("CreateDeck handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response model.Deck
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if response.ID != existingDeck.ID {
			t.Errorf("CreateDeck handler returned wrong deck ID: got %v want %v", response.ID, existingDeck.ID)
		}

		if response.Remaining != existingDeck.Remaining {
			t.Errorf("CreateDeck handler returned wrong remaining cards: got %v want %v", response.Remaining, existingDeck.Remaining)
		}
	})
}
