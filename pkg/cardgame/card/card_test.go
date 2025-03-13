package card

import (
	"testing"
)

func TestPlayingCard_ID(t *testing.T) {
	tests := []struct {
		name     string
		card     PlayingCard
		expected string
	}{
		{
			name:     "Ace of Spades",
			card:     PlayingCard{Suit: Spades, Value: 1},
			expected: "0-1",
		},
		{
			name:     "King of Hearts",
			card:     PlayingCard{Suit: Hearts, Value: 13},
			expected: "1-13",
		},
		{
			name:     "Two of Diamonds",
			card:     PlayingCard{Suit: Diamonds, Value: 2},
			expected: "2-2",
		},
		{
			name:     "Jack of Clubs",
			card:     PlayingCard{Suit: Clubs, Value: 11},
			expected: "3-11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.ID(); got != tt.expected {
				t.Errorf("PlayingCard.ID() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPlayingCard_Name(t *testing.T) {
	tests := []struct {
		name     string
		card     PlayingCard
		expected string
	}{
		{
			name:     "Ace of Spades",
			card:     PlayingCard{Suit: Spades, Value: 1},
			expected: "Ace of Spades",
		},
		{
			name:     "King of Hearts",
			card:     PlayingCard{Suit: Hearts, Value: 13},
			expected: "King of Hearts",
		},
		{
			name:     "Two of Diamonds",
			card:     PlayingCard{Suit: Diamonds, Value: 2},
			expected: "2 of Diamonds",
		},
		{
			name:     "Jack of Clubs",
			card:     PlayingCard{Suit: Clubs, Value: 11},
			expected: "Jack of Clubs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.Name(); got != tt.expected {
				t.Errorf("PlayingCard.Name() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPlayingCard_String(t *testing.T) {
	tests := []struct {
		name     string
		card     PlayingCard
		expected string
	}{
		{
			name:     "Ace of Spades",
			card:     PlayingCard{Suit: Spades, Value: 1},
			expected: "A♠",
		},
		{
			name:     "King of Hearts",
			card:     PlayingCard{Suit: Hearts, Value: 13},
			expected: "K♥",
		},
		{
			name:     "Two of Diamonds",
			card:     PlayingCard{Suit: Diamonds, Value: 2},
			expected: "2♦",
		},
		{
			name:     "Jack of Clubs",
			card:     PlayingCard{Suit: Clubs, Value: 11},
			expected: "J♣",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.String(); got != tt.expected {
				t.Errorf("PlayingCard.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPlayingCard_Equal(t *testing.T) {
	card1 := PlayingCard{Suit: Spades, Value: 1}
	card2 := PlayingCard{Suit: Spades, Value: 1}
	card3 := PlayingCard{Suit: Hearts, Value: 1}
	card4 := PlayingCard{Suit: Spades, Value: 2}

	tests := []struct {
		name     string
		card     PlayingCard
		other    Card
		expected bool
	}{
		{
			name:     "Same card",
			card:     card1,
			other:    card2,
			expected: true,
		},
		{
			name:     "Different suit",
			card:     card1,
			other:    card3,
			expected: false,
		},
		{
			name:     "Different value",
			card:     card1,
			other:    card4,
			expected: false,
		},
		{
			name:     "Same instance",
			card:     card1,
			other:    card1,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.Equal(tt.other); got != tt.expected {
				t.Errorf("PlayingCard.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPlayingCard_Clone(t *testing.T) {
	original := PlayingCard{Suit: Spades, Value: 1}
	clone := original.Clone()

	if !original.Equal(clone) {
		t.Errorf("Clone not equal to original: original=%v, clone=%v", original, clone)
	}

	// Verify it's a deep copy by modifying the clone
	cloneCard, ok := clone.(PlayingCard)
	if !ok {
		t.Fatalf("Clone is not a PlayingCard: %T", clone)
	}

	cloneCard.Value = 2
	if original.Equal(cloneCard) {
		t.Errorf("Clone should be a deep copy, but modifying clone affected original")
	}
}

func TestPlayingCardFactory_CreateCard(t *testing.T) {
	factory := &PlayingCardFactory{}

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		expected    PlayingCard
	}{
		{
			name: "Valid card",
			params: map[string]interface{}{
				"suit":  int(Hearts),
				"value": 10,
			},
			expectError: false,
			expected:    PlayingCard{Suit: Hearts, Value: 10},
		},
		{
			name: "Invalid suit",
			params: map[string]interface{}{
				"suit":  5,
				"value": 10,
			},
			expectError: true,
		},
		{
			name: "Invalid value",
			params: map[string]interface{}{
				"suit":  int(Diamonds),
				"value": 15,
			},
			expectError: true,
		},
		{
			name: "Missing suit",
			params: map[string]interface{}{
				"value": 10,
			},
			expectError: true,
		},
		{
			name: "Missing value",
			params: map[string]interface{}{
				"suit": int(Clubs),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card, err := factory.CreateCard(tt.params)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			playingCard, ok := card.(PlayingCard)
			if !ok {
				t.Fatalf("Expected PlayingCard but got %T", card)
			}

			if playingCard.Suit != tt.expected.Suit || playingCard.Value != tt.expected.Value {
				t.Errorf("Got card %v, expected %v", playingCard, tt.expected)
			}
		})
	}
}

func TestSuit_String(t *testing.T) {
	tests := []struct {
		name     string
		suit     Suit
		expected string
	}{
		{"Spades", Spades, "Spades"},
		{"Hearts", Hearts, "Hearts"},
		{"Diamonds", Diamonds, "Diamonds"},
		{"Clubs", Clubs, "Clubs"},
		{"Invalid", Suit(10), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.suit.String(); got != tt.expected {
				t.Errorf("Suit.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSuit_Symbol(t *testing.T) {
	tests := []struct {
		name     string
		suit     Suit
		expected string
	}{
		{"Spades", Spades, "♠"},
		{"Hearts", Hearts, "♥"},
		{"Diamonds", Diamonds, "♦"},
		{"Clubs", Clubs, "♣"},
		{"Invalid", Suit(10), "?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.suit.Symbol(); got != tt.expected {
				t.Errorf("Suit.Symbol() = %v, want %v", got, tt.expected)
			}
		})
	}
}
