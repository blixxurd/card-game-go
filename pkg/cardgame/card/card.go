// Package card defines the core interfaces and implementations for cards in various card games.
package card

// Card is the base interface that all card implementations must satisfy.
// It provides methods for identifying, comparing, and displaying cards.
type Card interface {
	// ID returns a unique identifier for the card
	ID() string

	// Name returns a human-readable name for the card
	Name() string

	// String returns a string representation of the card
	String() string

	// Equal checks if two cards are equivalent
	Equal(Card) bool

	// Clone creates a deep copy of the card
	Clone() Card
}

// CardType represents the type of card (e.g., standard playing card)
type CardType string

const (
	// StandardCard represents a standard playing card
	StandardCard CardType = "standard"
)

// CardFactory is an interface for creating cards of different types
type CardFactory interface {
	// CreateCard creates a new card based on the provided parameters
	CreateCard(params map[string]interface{}) (Card, error)

	// Type returns the type of cards this factory creates
	Type() CardType
}
