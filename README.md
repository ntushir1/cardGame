# Card Game
This project implements a simple set of card game API in Golang using Gorilla Mux for routing.

## Prerequisites

Before you begin, ensure you have the following installed:

- Go (version 1.13 or later)
- Git

## Setup

```
   git clone <repository-url>
```

## Install dependencies

  ```
  go get -u ./...
  ```

## Build the project using the following command

```
  go build
```

This will generate an executable file (e.g., deck) in the project directory.

## Run the project:

```
  ./deck
```

The server will start, and you can access the API at http://localhost:8080.


# Running Tests

## Run tests using the following command:
```
go test ./...
```

# API Endpoints

### Health Check: http://localhost:8080/health (GET)
### Create a new deck: http://localhost:8080/deck (GET)
### Draw cards from a deck: http://localhost:8080/deck/{deckID}/draw?count=2 (GET)
# Deck API

This API allows you to manage decks of playing cards.

## Health Check

Check the health of the API.

- **URL:** `/health`
- **Method:** `GET`
- **Response:**
  - Status: 200 OK
  - Body: "OK"

## Create a New Deck

Create a new deck of playing cards.

- **URL:** `/deck`
- **Method:** `GET`
- **Query Parameters:**
  - `shuffled` (optional): Shuffle the deck. Default is `false`.
  - `cards` (optional): Comma-separated list of cards to include in the deck.
- **Response:**
  - Status: 200 OK
  - Body Example:
    ```json
    {
      "deck_id": "b2bc11b8-9ab4-11ee-8065-acde48001122",
      "shuffled": true,
      "remaining": 52
    }
    ```

## Draw Cards from a Deck

Draw a specified number of cards from a deck.

- **URL:** `/deck/{deckID}/draw`
- **Method:** `GET`
- **Path Parameters:**
  - `deckID` (required): The ID of the deck.
- **Query Parameters:**
  - `count` (required): The number of cards to draw.
- **Response:**
  - Status: 200 OK
  - Body Example:
    ```json
    {
      "cards": [
        {"value": "ACE", "suit": "SPADES", "code": "AS"},
        {"value": "2", "suit": "HEARTS", "code": "2H"}
      ]
    }
    ```






