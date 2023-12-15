package main_test

import (
	"cardGame/deck/model"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"cardGame/deck/api"
	"cardGame/deck/dao"
	"cardGame/deck/service"
)

func configureRoutes(deckHandler *api.DeckHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", deckHandler.HealthCheck).Methods("GET")
	router.HandleFunc("/deck", deckHandler.CreateDeck).Methods("GET")
	router.HandleFunc("/deck/{deckID}/draw", deckHandler.DrawCards).Methods("GET")

	return router
}

func TestMain(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Main Suite")
}

var _ = ginkgo.Describe("Deck API", func() {
	var (
		deckStorage *dao.DeckStorage
		deckService *service.DeckService
		deckHandler *api.DeckHandler
		router      *mux.Router
	)

	ginkgo.BeforeEach(func() {
		deckStorage = dao.NewDeckStorage()
		deckService = service.NewDeckService(deckStorage)
		deckHandler = api.NewDeckHandler(deckService, deckStorage)
		router = configureRoutes(deckHandler)
	})

	ginkgo.It("should respond with 'OK' for health check", func() {
		req, err := http.NewRequest("GET", "/health", nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body.String()).To(gomega.Equal("OK"))
	})

	ginkgo.It("should create a new deck", func() {
		req, err := http.NewRequest("GET", "/deck", nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))

		var response api.CreateDeckResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		gomega.Expect(response.DeckID).NotTo(gomega.BeEmpty())
		gomega.Expect(response.Shuffled).To(gomega.BeFalse())
		gomega.Expect(response.Remaining).To(gomega.Equal(52))
	})

	ginkgo.It("should draw cards from a deck", func() {
		deckIDStr := "b2bc11b8-9ab4-11ee-8065-acde48001122"
		deckID, err := uuid.Parse(deckIDStr)
		deckStorage.SaveDeck(model.Deck{
			ID:        deckID,
			Shuffled:  false,
			Remaining: 13,
			Cards: []model.Card{
				{Value: "ACE", Suit: "SPADES", Code: "AS"},
				{Value: "2", Suit: "SPADES", Code: "2S"},
				{Value: "3", Suit: "SPADES", Code: "3S"},
				{Value: "4", Suit: "SPADES", Code: "4S"},
				{Value: "5", Suit: "SPADES", Code: "5S"},
				{Value: "6", Suit: "SPADES", Code: "6S"},
				{Value: "7", Suit: "SPADES", Code: "7S"},
				{Value: "8", Suit: "SPADES", Code: "8S"},
				{Value: "9", Suit: "SPADES", Code: "9S"},
				{Value: "10", Suit: "SPADES", Code: "10S"},
				{Value: "JACK", Suit: "SPADES", Code: "JS"},
				{Value: "QUEEN", Suit: "SPADES", Code: "QS"},
				{Value: "KING", Suit: "SPADES", Code: "KS"},
			},
		})

		req, err := http.NewRequest("GET", "/deck/"+deckID.String()+"/draw?count=2", nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))

		var response map[string][]model.Card
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		drawnCards, found := response["cards"]
		gomega.Expect(found).To(gomega.BeTrue())
		gomega.Expect(len(drawnCards)).To(gomega.Equal(2))
	})

	ginkgo.It("should return an error for invalid deck ID during card draw", func() {
		req, err := http.NewRequest("GET", "/deck/invalid-deck-id/draw?count=2", nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
	})

	ginkgo.It("should return an error for invalid count parameter during card draw", func() {
		deckIDStr := "b2bc11b8-9ab4-11ee-8065-acde48001122"
		deckID, err := uuid.Parse(deckIDStr)
		deckStorage.SaveDeck(model.Deck{
			ID:        deckID,
			Shuffled:  false,
			Remaining: 10,
			Cards: []model.Card{
				{Value: "ACE", Suit: "SPADES", Code: "AS"},
				{Value: "2", Suit: "SPADES", Code: "2S"},
				{Value: "3", Suit: "SPADES", Code: "3S"},
				{Value: "4", Suit: "SPADES", Code: "4S"},
				{Value: "5", Suit: "SPADES", Code: "5S"},
				{Value: "6", Suit: "SPADES", Code: "6S"},
				{Value: "7", Suit: "SPADES", Code: "7S"},
				{Value: "8", Suit: "SPADES", Code: "8S"},
				{Value: "9", Suit: "SPADES", Code: "9S"},
				{Value: "10", Suit: "SPADES", Code: "10S"},
			},
		})

		req, err := http.NewRequest("GET", "/deck/"+deckID.String()+"/draw?count=invalid-count", nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
	})

	ginkgo.It("should return an error for drawing cards from a non-existing deck", func() {
		req, err := http.NewRequest("GET", "/deck/"+uuid.New().String()+"/draw?count=2", nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusNotFound))
	})

	ginkgo.It("should return an error for drawing more cards than remaining in the deck", func() {
		deckIDStr := "b2bc11b8-9ab4-11ee-8065-acde48001122"
		deckID, err := uuid.Parse(deckIDStr)
		deckStorage.SaveDeck(model.Deck{
			ID:        deckID,
			Shuffled:  false,
			Remaining: 2,
			Cards: []model.Card{
				{Value: "ACE", Suit: "SPADES", Code: "AS"},
				{Value: "2", Suit: "SPADES", Code: "2S"},
			},
		})

		req, err := http.NewRequest("GET", "/deck/"+deckID.String()+"/draw?count=3", nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
	})

	ginkgo.It("should shuffle the deck when creating a new shuffled deck", func() {
		req, err := http.NewRequest("GET", "/deck?shuffled=true", nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))

		var response api.CreateDeckResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		gomega.Expect(response.Shuffled).To(gomega.BeTrue())
	})
})
