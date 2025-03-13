// Package hand defines interfaces and implementations for card hands in various card games.
package hand

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
)

// StandardHand represents a basic hand of cards
type StandardHand struct {
	cards []card.Card
}

// NewStandardHand creates a new empty hand
func NewStandardHand() *StandardHand {
	return &StandardHand{
		cards: make([]card.Card, 0),
	}
}

// AddCard adds a card to the hand
func (h *StandardHand) AddCard(c card.Card) error {
	if c == nil {
		return errors.New("cannot add nil card to hand")
	}
	h.cards = append(h.cards, c)
	return nil
}

// AddCards adds multiple cards to the hand
func (h *StandardHand) AddCards(cards []card.Card) error {
	if cards == nil {
		return errors.New("cannot add nil cards to hand")
	}

	for _, c := range cards {
		if err := h.AddCard(c); err != nil {
			return err
		}
	}
	return nil
}

// RemoveCard removes a specific card from the hand
func (h *StandardHand) RemoveCard(c card.Card) (bool, error) {
	if c == nil {
		return false, errors.New("cannot remove nil card from hand")
	}

	for i, handCard := range h.cards {
		if handCard.Equal(c) {
			h.cards = append(h.cards[:i], h.cards[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

// HasCard checks if the hand contains a specific card
func (h *StandardHand) HasCard(c card.Card) bool {
	if c == nil {
		return false
	}

	for _, handCard := range h.cards {
		if handCard.Equal(c) {
			return true
		}
	}
	return false
}

// Size returns the number of cards in the hand
func (h *StandardHand) Size() int {
	return len(h.cards)
}

// IsEmpty returns true if the hand has no cards
func (h *StandardHand) IsEmpty() bool {
	return len(h.cards) == 0
}

// Clear removes all cards from the hand
func (h *StandardHand) Clear() {
	h.cards = h.cards[:0]
}

// Cards returns all cards in the hand (without removing them)
func (h *StandardHand) Cards() []card.Card {
	result := make([]card.Card, len(h.cards))
	copy(result, h.cards)
	return result
}

// String returns a string representation of the hand
func (h *StandardHand) String() string {
	if h.IsEmpty() {
		return "Empty Hand"
	}

	var sb strings.Builder
	for i, c := range h.cards {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(c.String())
	}
	return sb.String()
}

// Clone creates a deep copy of the hand
func (h *StandardHand) Clone() Hand {
	clone := &StandardHand{
		cards: make([]card.Card, len(h.cards)),
	}

	for i, c := range h.cards {
		clone.cards[i] = c.Clone()
	}

	return clone
}

// StandardHandFactory creates standard hands
type StandardHandFactory struct{}

// CreateHand creates a new standard hand
func (f *StandardHandFactory) CreateHand(params map[string]interface{}) (Hand, error) {
	hand := NewStandardHand()

	// Check if we should add initial cards
	if initialCards, ok := params["cards"].([]card.Card); ok {
		if err := hand.AddCards(initialCards); err != nil {
			return nil, fmt.Errorf("failed to add initial cards: %w", err)
		}
	}

	return hand, nil
}

// Type returns the type of hands this factory creates
func (f *StandardHandFactory) Type() HandType {
	return PokerHand // Using PokerHand as the default for standard hands
}
