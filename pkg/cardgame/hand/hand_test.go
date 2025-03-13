package hand

import (
	"strings"
	"testing"

	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
)

func TestStandardHand_NewStandardHand(t *testing.T) {
	hand := NewStandardHand()

	// A new hand should be empty
	if !hand.IsEmpty() {
		t.Errorf("Expected new hand to be empty")
	}

	if hand.Size() != 0 {
		t.Errorf("Expected new hand size to be 0, got %d", hand.Size())
	}
}

func TestStandardHand_AddCard(t *testing.T) {
	hand := NewStandardHand()
	testCard := card.PlayingCard{Suit: card.Spades, Value: 1}

	// Add a card
	err := hand.AddCard(testCard)
	if err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}

	// Check that the hand size increased
	if hand.Size() != 1 {
		t.Errorf("Expected hand size to be 1, got %d", hand.Size())
	}

	// Check that the hand is not empty
	if hand.IsEmpty() {
		t.Errorf("Expected hand not to be empty after adding a card")
	}

	// Check that the card is in the hand
	if !hand.HasCard(testCard) {
		t.Errorf("Expected hand to contain the added card")
	}

	// Try to add a nil card
	err = hand.AddCard(nil)
	if err == nil {
		t.Errorf("Expected error when adding nil card, got nil")
	}
}

func TestStandardHand_AddCards(t *testing.T) {
	hand := NewStandardHand()
	testCards := []card.Card{
		card.PlayingCard{Suit: card.Spades, Value: 1},
		card.PlayingCard{Suit: card.Hearts, Value: 10},
		card.PlayingCard{Suit: card.Diamonds, Value: 5},
	}

	// Add multiple cards
	err := hand.AddCards(testCards)
	if err != nil {
		t.Fatalf("Failed to add cards: %v", err)
	}

	// Check that the hand size increased
	if hand.Size() != len(testCards) {
		t.Errorf("Expected hand size to be %d, got %d", len(testCards), hand.Size())
	}

	// Check that all cards are in the hand
	for _, c := range testCards {
		if !hand.HasCard(c) {
			t.Errorf("Expected hand to contain card %v", c)
		}
	}

	// Try to add nil cards
	err = hand.AddCards(nil)
	if err == nil {
		t.Errorf("Expected error when adding nil cards, got nil")
	}

	// Try to add a slice with a nil card
	nilCards := []card.Card{nil}
	err = hand.AddCards(nilCards)
	if err == nil {
		t.Errorf("Expected error when adding slice with nil card, got nil")
	}
}

func TestStandardHand_RemoveCard(t *testing.T) {
	hand := NewStandardHand()
	testCard := card.PlayingCard{Suit: card.Spades, Value: 1}

	// Add a card to remove later
	err := hand.AddCard(testCard)
	if err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}

	// Remove the card
	removed, err := hand.RemoveCard(testCard)
	if err != nil {
		t.Fatalf("Failed to remove card: %v", err)
	}

	if !removed {
		t.Errorf("Expected card to be removed, but RemoveCard returned false")
	}

	// Check that the hand is now empty
	if !hand.IsEmpty() {
		t.Errorf("Expected hand to be empty after removing the only card")
	}

	// Try to remove a card that's not in the hand
	removed, err = hand.RemoveCard(testCard)
	if err != nil {
		t.Fatalf("Failed to attempt removing non-existent card: %v", err)
	}

	if removed {
		t.Errorf("Expected RemoveCard to return false for non-existent card, got true")
	}

	// Try to remove a nil card
	_, err = hand.RemoveCard(nil)
	if err == nil {
		t.Errorf("Expected error when removing nil card, got nil")
	}
}

func TestStandardHand_HasCard(t *testing.T) {
	hand := NewStandardHand()
	testCard := card.PlayingCard{Suit: card.Spades, Value: 1}
	otherCard := card.PlayingCard{Suit: card.Hearts, Value: 10}

	// Check that the hand doesn't have any cards initially
	if hand.HasCard(testCard) {
		t.Errorf("Expected empty hand not to have any cards")
	}

	// Add a card
	err := hand.AddCard(testCard)
	if err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}

	// Check that the hand has the added card
	if !hand.HasCard(testCard) {
		t.Errorf("Expected hand to have the added card")
	}

	// Check that the hand doesn't have a different card
	if hand.HasCard(otherCard) {
		t.Errorf("Expected hand not to have a card that wasn't added")
	}

	// Check behavior with nil card
	if hand.HasCard(nil) {
		t.Errorf("Expected HasCard to return false for nil card")
	}
}

func TestStandardHand_Clear(t *testing.T) {
	hand := NewStandardHand()
	testCards := []card.Card{
		card.PlayingCard{Suit: card.Spades, Value: 1},
		card.PlayingCard{Suit: card.Hearts, Value: 10},
		card.PlayingCard{Suit: card.Diamonds, Value: 5},
	}

	// Add some cards
	err := hand.AddCards(testCards)
	if err != nil {
		t.Fatalf("Failed to add cards: %v", err)
	}

	// Clear the hand
	hand.Clear()

	// Check that the hand is empty
	if !hand.IsEmpty() {
		t.Errorf("Expected hand to be empty after clearing")
	}

	if hand.Size() != 0 {
		t.Errorf("Expected hand size to be 0 after clearing, got %d", hand.Size())
	}

	// Check that none of the original cards are in the hand
	for _, c := range testCards {
		if hand.HasCard(c) {
			t.Errorf("Expected hand not to contain card %v after clearing", c)
		}
	}
}

func TestStandardHand_Cards(t *testing.T) {
	hand := NewStandardHand()
	testCards := []card.Card{
		card.PlayingCard{Suit: card.Spades, Value: 1},
		card.PlayingCard{Suit: card.Hearts, Value: 10},
		card.PlayingCard{Suit: card.Diamonds, Value: 5},
	}

	// Check that a new hand has no cards
	if len(hand.Cards()) != 0 {
		t.Errorf("Expected new hand to have no cards, got %d", len(hand.Cards()))
	}

	// Add some cards
	err := hand.AddCards(testCards)
	if err != nil {
		t.Fatalf("Failed to add cards: %v", err)
	}

	// Check that Cards() returns all the added cards
	cards := hand.Cards()
	if len(cards) != len(testCards) {
		t.Errorf("Expected Cards() to return %d cards, got %d", len(testCards), len(cards))
	}

	// Check that the cards are in the same order they were added
	for i, c := range testCards {
		if !c.Equal(cards[i]) {
			t.Errorf("Card at index %d differs from expected", i)
		}
	}

	// Verify that Cards() returns a copy, not the original slice
	originalCards := hand.Cards()
	hand.AddCard(card.PlayingCard{Suit: card.Clubs, Value: 7})
	newCards := hand.Cards()

	if len(originalCards) == len(newCards) {
		t.Errorf("Expected Cards() to return a copy of the cards slice")
	}
}

func TestStandardHand_String(t *testing.T) {
	hand := NewStandardHand()

	// Check that an empty hand has a reasonable string representation
	emptyString := hand.String()
	if emptyString == "" {
		t.Errorf("Expected non-empty string representation for empty hand")
	}

	// Add some cards
	testCards := []card.Card{
		card.PlayingCard{Suit: card.Spades, Value: 1},
		card.PlayingCard{Suit: card.Hearts, Value: 10},
	}
	err := hand.AddCards(testCards)
	if err != nil {
		t.Fatalf("Failed to add cards: %v", err)
	}

	// Check that the string representation includes the cards
	handString := hand.String()
	if handString == "" {
		t.Errorf("Expected non-empty string representation for hand with cards")
	}

	// The string should contain some representation of each card
	for _, c := range testCards {
		if !strings.Contains(handString, c.String()) {
			t.Errorf("Expected string representation to include card %v", c)
		}
	}
}

func TestStandardHand_Clone(t *testing.T) {
	hand := NewStandardHand()
	testCards := []card.Card{
		card.PlayingCard{Suit: card.Spades, Value: 1},
		card.PlayingCard{Suit: card.Hearts, Value: 10},
		card.PlayingCard{Suit: card.Diamonds, Value: 5},
	}

	// Add some cards
	err := hand.AddCards(testCards)
	if err != nil {
		t.Fatalf("Failed to add cards: %v", err)
	}

	// Clone the hand
	clonedHand := hand.Clone()

	// Check that the clone has the same size
	if clonedHand.Size() != hand.Size() {
		t.Errorf("Expected cloned hand size to be %d, got %d", hand.Size(), clonedHand.Size())
	}

	// Check that the clone has the same cards
	originalCards := hand.Cards()
	clonedCards := clonedHand.Cards()

	for i, originalCard := range originalCards {
		if !originalCard.Equal(clonedCards[i]) {
			t.Errorf("Card at index %d differs between original and clone", i)
		}
	}

	// Modify the original hand
	hand.AddCard(card.PlayingCard{Suit: card.Clubs, Value: 7})

	// Check that the clone was not affected
	if clonedHand.Size() != len(testCards) {
		t.Errorf("Expected cloned hand size to remain %d, got %d", len(testCards), clonedHand.Size())
	}
}

func TestStandardHandFactory_CreateHand(t *testing.T) {
	factory := &StandardHandFactory{}

	// Create a hand with default parameters
	hand, err := factory.CreateHand(nil)
	if err != nil {
		t.Fatalf("Failed to create hand: %v", err)
	}

	// Check that it's a standard hand
	if !hand.IsEmpty() {
		t.Errorf("Expected new hand to be empty")
	}

	// Check that the Type() method returns the correct value
	if factory.Type() != PokerHand {
		t.Errorf("Expected factory type to be PokerHand, got %v", factory.Type())
	}
}
