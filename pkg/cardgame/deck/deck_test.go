package deck

import (
	"testing"

	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
)

func TestStandardDeck_NewStandardDeck(t *testing.T) {
	deck := NewStandardDeck()

	// A standard deck should have 52 cards
	if deck.Size() != 52 {
		t.Errorf("Expected deck size to be 52, got %d", deck.Size())
	}

	// Check that all cards are valid
	for _, c := range deck.Cards() {
		playingCard, ok := c.(card.PlayingCard)
		if !ok {
			t.Errorf("Expected card to be PlayingCard, got %T", c)
			continue
		}

		if playingCard.Value < 1 || playingCard.Value > 13 {
			t.Errorf("Card has invalid value: %d", playingCard.Value)
		}

		if playingCard.Suit < card.Spades || playingCard.Suit > card.Clubs {
			t.Errorf("Card has invalid suit: %d", playingCard.Suit)
		}
	}
}

func TestStandardDeck_Shuffle(t *testing.T) {
	deck := NewStandardDeck()

	// Store the original order
	originalOrder := make([]card.Card, len(deck.Cards()))
	copy(originalOrder, deck.Cards())

	// Shuffle the deck
	err := deck.Shuffle()
	if err != nil {
		t.Fatalf("Failed to shuffle deck: %v", err)
	}

	// Check that the order has changed
	// Note: There's a small chance this could fail randomly if the shuffle
	// happens to result in the same order, but it's extremely unlikely
	sameOrder := true
	for i, c := range deck.Cards() {
		if !c.Equal(originalOrder[i]) {
			sameOrder = false
			break
		}
	}

	if sameOrder {
		t.Errorf("Deck was not shuffled, order remains the same")
	}
}

func TestStandardDeck_Draw(t *testing.T) {
	deck := NewStandardDeck()

	originalSize := deck.Size()

	// Draw a card
	drawnCard, err := deck.Draw()
	if err != nil {
		t.Fatalf("Failed to draw card: %v", err)
	}

	// Check that the deck size decreased
	if deck.Size() != originalSize-1 {
		t.Errorf("Expected deck size to be %d, got %d", originalSize-1, deck.Size())
	}

	// Check that the drawn card is valid
	_, ok := drawnCard.(card.PlayingCard)
	if !ok {
		t.Errorf("Expected drawn card to be PlayingCard, got %T", drawnCard)
	}

	// Draw all remaining cards
	for i := 0; i < originalSize-1; i++ {
		_, err := deck.Draw()
		if err != nil {
			t.Fatalf("Failed to draw card at index %d: %v", i, err)
		}
	}

	// Check that the deck is now empty
	if !deck.IsEmpty() {
		t.Errorf("Expected deck to be empty after drawing all cards")
	}

	// Try to draw from an empty deck
	_, err = deck.Draw()
	if err == nil {
		t.Errorf("Expected error when drawing from empty deck, got nil")
	}
}

func TestStandardDeck_DrawN(t *testing.T) {
	deck := NewStandardDeck()

	originalSize := deck.Size()

	// Draw 5 cards
	drawnCards, err := deck.DrawN(5)
	if err != nil {
		t.Fatalf("Failed to draw cards: %v", err)
	}

	// Check that we got the right number of cards
	if len(drawnCards) != 5 {
		t.Errorf("Expected to draw 5 cards, got %d", len(drawnCards))
	}

	// Check that the deck size decreased
	if deck.Size() != originalSize-5 {
		t.Errorf("Expected deck size to be %d, got %d", originalSize-5, deck.Size())
	}

	// Try to draw more cards than are in the deck
	_, err = deck.DrawN(originalSize)
	if err == nil {
		t.Errorf("Expected error when drawing more cards than are in the deck, got nil")
	}

	// Try to draw a negative number of cards
	_, err = deck.DrawN(-1)
	if err == nil {
		t.Errorf("Expected error when drawing a negative number of cards, got nil")
	}
}

func TestStandardDeck_AddCard(t *testing.T) {
	deck := NewStandardDeck()

	// Draw a card to add back later
	drawnCard, err := deck.Draw()
	if err != nil {
		t.Fatalf("Failed to draw card: %v", err)
	}

	originalSize := deck.Size()

	// Add the card back
	err = deck.AddCard(drawnCard)
	if err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}

	// Check that the deck size increased
	if deck.Size() != originalSize+1 {
		t.Errorf("Expected deck size to be %d, got %d", originalSize+1, deck.Size())
	}

	// Check that the added card is in the deck
	found := false
	for _, c := range deck.Cards() {
		if c.Equal(drawnCard) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Added card not found in deck")
	}
}

func TestStandardDeck_AddCards(t *testing.T) {
	deck := NewStandardDeck()

	// Draw some cards to add back later
	drawnCards, err := deck.DrawN(5)
	if err != nil {
		t.Fatalf("Failed to draw cards: %v", err)
	}

	originalSize := deck.Size()

	// Add the cards back
	err = deck.AddCards(drawnCards)
	if err != nil {
		t.Fatalf("Failed to add cards: %v", err)
	}

	// Check that the deck size increased
	if deck.Size() != originalSize+5 {
		t.Errorf("Expected deck size to be %d, got %d", originalSize+5, deck.Size())
	}

	// Check that all added cards are in the deck
	for _, drawnCard := range drawnCards {
		found := false
		for _, c := range deck.Cards() {
			if c.Equal(drawnCard) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Added card not found in deck: %v", drawnCard)
		}
	}
}

func TestStandardDeck_RemoveCard(t *testing.T) {
	deck := NewStandardDeck()

	// Get a card to remove
	cards := deck.Cards()
	cardToRemove := cards[0]
	originalSize := deck.Size()

	// Remove the card
	removed, err := deck.RemoveCard(cardToRemove)
	if err != nil {
		t.Fatalf("Failed to remove card: %v", err)
	}

	if !removed {
		t.Errorf("Expected card to be removed, but RemoveCard returned false")
	}

	// Check that the deck size decreased
	if deck.Size() != originalSize-1 {
		t.Errorf("Expected deck size to be %d, got %d", originalSize-1, deck.Size())
	}

	// Check that the card is no longer in the deck
	for _, c := range deck.Cards() {
		if c.Equal(cardToRemove) {
			t.Errorf("Removed card still found in deck")
		}
	}

	// Try to remove a card that's not in the deck
	removed, err = deck.RemoveCard(cardToRemove)
	if err != nil {
		t.Fatalf("Failed to attempt removing non-existent card: %v", err)
	}

	if removed {
		t.Errorf("Expected RemoveCard to return false for non-existent card, got true")
	}
}

func TestStandardDeck_Reset(t *testing.T) {
	deck := NewStandardDeck()

	// Draw some cards
	_, err := deck.DrawN(10)
	if err != nil {
		t.Fatalf("Failed to draw cards: %v", err)
	}

	// Reset the deck
	err = deck.Reset()
	if err != nil {
		t.Fatalf("Failed to reset deck: %v", err)
	}

	// Check that the deck has 52 cards again
	if deck.Size() != 52 {
		t.Errorf("Expected deck size to be 52 after reset, got %d", deck.Size())
	}

	// Check that all cards are valid
	for _, c := range deck.Cards() {
		playingCard, ok := c.(card.PlayingCard)
		if !ok {
			t.Errorf("Expected card to be PlayingCard, got %T", c)
			continue
		}

		if playingCard.Value < 1 || playingCard.Value > 13 {
			t.Errorf("Card has invalid value after reset: %d", playingCard.Value)
		}

		if playingCard.Suit < card.Spades || playingCard.Suit > card.Clubs {
			t.Errorf("Card has invalid suit after reset: %d", playingCard.Suit)
		}
	}
}

func TestStandardDeck_Clone(t *testing.T) {
	deck := NewStandardDeck()

	// Shuffle the deck to get a unique order
	err := deck.Shuffle()
	if err != nil {
		t.Fatalf("Failed to shuffle deck: %v", err)
	}

	// Clone the deck
	clonedDeck := deck.Clone()

	// Check that the clone has the same size
	if clonedDeck.Size() != deck.Size() {
		t.Errorf("Expected cloned deck size to be %d, got %d", deck.Size(), clonedDeck.Size())
	}

	// Check that the clone has the same cards in the same order
	originalCards := deck.Cards()
	clonedCards := clonedDeck.Cards()

	for i, originalCard := range originalCards {
		if !originalCard.Equal(clonedCards[i]) {
			t.Errorf("Card at index %d differs between original and clone", i)
		}
	}

	// Modify the original deck
	_, err = deck.Draw()
	if err != nil {
		t.Fatalf("Failed to draw card: %v", err)
	}

	// Check that the clone was not affected
	if clonedDeck.Size() != 52 {
		t.Errorf("Expected cloned deck size to remain 52, got %d", clonedDeck.Size())
	}
}

func TestStandardDeckFactory_CreateDeck(t *testing.T) {
	factory := &StandardDeckFactory{}

	// Create a deck with default parameters
	deck, err := factory.CreateDeck(nil)
	if err != nil {
		t.Fatalf("Failed to create deck: %v", err)
	}

	// Check that it's a standard deck
	if deck.Size() != 52 {
		t.Errorf("Expected deck size to be 52, got %d", deck.Size())
	}

	// Create a deck with custom parameters
	deck, err = factory.CreateDeck(map[string]interface{}{
		"shuffled": true,
	})
	if err != nil {
		t.Fatalf("Failed to create shuffled deck: %v", err)
	}

	// Check that it's a standard deck
	if deck.Size() != 52 {
		t.Errorf("Expected deck size to be 52, got %d", deck.Size())
	}

	// We can't directly test that it's shuffled since we don't know the original order,
	// but we can check that the Type() method returns the correct value
	if factory.Type() != StandardDeck {
		t.Errorf("Expected factory type to be StandardDeck, got %v", factory.Type())
	}
}
