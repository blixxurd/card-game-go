// Package deck defines interfaces and implementations for card decks in various card games.
package deck

import (
	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
)

// Deck is the base interface that all deck implementations must satisfy.
// It provides methods for managing a collection of cards.
type Deck interface {
	// Shuffle randomizes the order of cards in the deck
	Shuffle() error

	// Draw removes and returns the top card from the deck
	Draw() (card.Card, error)

	// DrawN removes and returns the top n cards from the deck
	DrawN(n int) ([]card.Card, error)

	// AddCard adds a card to the deck
	AddCard(c card.Card) error

	// AddCards adds multiple cards to the deck
	AddCards(cards []card.Card) error

	// RemoveCard removes a specific card from the deck
	RemoveCard(c card.Card) (bool, error)

	// Size returns the number of cards in the deck
	Size() int

	// IsEmpty returns true if the deck has no cards
	IsEmpty() bool

	// Reset returns the deck to its initial state
	Reset() error

	// Cards returns all cards in the deck (without removing them)
	Cards() []card.Card

	// Clone creates a deep copy of the deck
	Clone() Deck
}

// DeckType represents the type of deck (e.g., standard 52-card deck)
type DeckType string

const (
	// StandardDeck represents a standard 52-card deck
	StandardDeck DeckType = "standard"
)

// DeckFactory is an interface for creating decks of different types
type DeckFactory interface {
	// CreateDeck creates a new deck based on the provided parameters
	CreateDeck(params map[string]interface{}) (Deck, error)

	// Type returns the type of decks this factory creates
	Type() DeckType
}
