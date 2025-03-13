// Package holdem provides a Texas Hold'em poker game implementation.
package holdem

import (
	"github.com/blixxurd/card-game-go/pkg/cardgame/poker"
	"github.com/blixxurd/card-game-go/pkg/cardgame/table"
)

// Define Texas Hold'em specific game states as poker.PokerGameState
const (
	// HoldemStateInitial is the initial state of the game
	HoldemStateInitial poker.PokerGameState = "initial"
	// HoldemStatePreFlop is the pre-flop state
	HoldemStatePreFlop poker.PokerGameState = "pre_flop"
	// HoldemStateFlop is the flop state
	HoldemStateFlop poker.PokerGameState = "flop"
	// HoldemStateTurn is the turn state
	HoldemStateTurn poker.PokerGameState = "turn"
	// HoldemStateRiver is the river state
	HoldemStateRiver poker.PokerGameState = "river"
	// HoldemStateShowdown is the showdown state
	HoldemStateShowdown poker.PokerGameState = "showdown"
	// HoldemStateComplete is the complete state
	HoldemStateComplete poker.PokerGameState = "complete"
)

// HoldemVariant implements the PokerVariant interface for Texas Hold'em
type HoldemVariant struct{}

// NewHoldemVariant creates a new HoldemVariant
func NewHoldemVariant() *HoldemVariant {
	return &HoldemVariant{}
}

// Name returns the name of the poker variant
func (v *HoldemVariant) Name() string {
	return "Texas Hold'em"
}

// Description returns a description of the poker variant
func (v *HoldemVariant) Description() string {
	return "Texas Hold'em is a community card poker variant where each player is dealt two private cards and five community cards are dealt face up."
}

// CreateGame creates a new instance of the poker variant
func (v *HoldemVariant) CreateGame(id, name string, smallBlind, bigBlind int, bettingStructure table.BettingStructure) (poker.PokerGame, error) {
	holdemGame := NewHoldemGame(id, name, smallBlind, bigBlind)
	return NewHoldemGameAdapter(holdemGame), nil
}

// MaxPlayers returns the maximum number of players allowed
func (v *HoldemVariant) MaxPlayers() int {
	return 9 // Standard poker table size
}

// MinPlayers returns the minimum number of players required
func (v *HoldemVariant) MinPlayers() int {
	return 2 // Need at least 2 players for a game
}

// HandSize returns the number of cards in a player's hand
func (v *HoldemVariant) HandSize() int {
	return 2 // Each player gets 2 hole cards
}

// CommunityCardCount returns the number of community cards used
func (v *HoldemVariant) CommunityCardCount() int {
	return 5 // 5 community cards (flop, turn, river)
}

// BettingRounds returns the number of betting rounds
func (v *HoldemVariant) BettingRounds() int {
	return 4 // Pre-flop, flop, turn, river
}

// StateSequence returns the sequence of game states
func (v *HoldemVariant) StateSequence() []poker.PokerGameState {
	return []poker.PokerGameState{
		HoldemStateInitial,
		HoldemStatePreFlop,
		HoldemStateFlop,
		HoldemStateTurn,
		HoldemStateRiver,
		HoldemStateShowdown,
		HoldemStateComplete,
	}
}
