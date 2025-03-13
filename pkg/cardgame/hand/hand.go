// Package hand defines interfaces and implementations for card hands in various card games.
package hand

import (
	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
)

// Hand is the base interface that all hand implementations must satisfy.
// It provides methods for managing a player's collection of cards.
type Hand interface {
	// AddCard adds a card to the hand
	AddCard(c card.Card) error

	// AddCards adds multiple cards to the hand
	AddCards(cards []card.Card) error

	// RemoveCard removes a specific card from the hand
	RemoveCard(c card.Card) (bool, error)

	// HasCard checks if the hand contains a specific card
	HasCard(c card.Card) bool

	// Size returns the number of cards in the hand
	Size() int

	// IsEmpty returns true if the hand has no cards
	IsEmpty() bool

	// Clear removes all cards from the hand
	Clear()

	// Cards returns all cards in the hand (without removing them)
	Cards() []card.Card

	// String returns a string representation of the hand
	String() string

	// Clone creates a deep copy of the hand
	Clone() Hand
}

// HandType represents the type of hand (e.g., poker hand)
type HandType string

const (
	// PokerHand represents a poker hand
	PokerHand HandType = "poker"
)

// HandFactory is an interface for creating hands of different types
type HandFactory interface {
	// CreateHand creates a new hand based on the provided parameters
	CreateHand(params map[string]interface{}) (Hand, error)

	// Type returns the type of hands this factory creates
	Type() HandType
}
