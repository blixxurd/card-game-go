// Package table provides structures and functionality for managing poker tables,
// including seating arrangements, chip management, and pot tracking.
package table

import (
	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
)

// Seat represents a position at a poker table occupied by a player.
type Seat struct {
	// Position is the seat's position at the table (0-indexed)
	Position int

	// Player is the player occupying this seat
	Player player.Player

	// ChipStack is the number of chips the player has at this seat
	ChipStack int

	// IsOccupied indicates whether the seat is currently occupied
	IsOccupied bool

	// IsActive indicates whether the player is active in the current hand
	IsActive bool

	// HasFolded indicates whether the player has folded in the current hand
	HasFolded bool

	// IsAllIn indicates whether the player is all-in in the current hand
	IsAllIn bool

	// CurrentBet is the amount the player has bet in the current betting round
	CurrentBet int

	// TotalBet is the total amount the player has bet in the current hand
	TotalBet int
}

// NewSeat creates a new seat at the specified position.
func NewSeat(position int) *Seat {
	return &Seat{
		Position:   position,
		IsOccupied: false,
		IsActive:   false,
		HasFolded:  false,
		IsAllIn:    false,
		CurrentBet: 0,
		TotalBet:   0,
	}
}

// Occupy sets a player in this seat with the specified chip stack.
func (s *Seat) Occupy(player player.Player, chipStack int) {
	s.Player = player
	s.ChipStack = chipStack
	s.IsOccupied = true
}

// Vacate removes the player from this seat.
func (s *Seat) Vacate() {
	s.Player = nil
	s.IsOccupied = false
	s.Reset()
}

// Reset resets the seat's state for a new hand.
func (s *Seat) Reset() {
	s.IsActive = s.IsOccupied && s.ChipStack > 0
	s.HasFolded = false
	s.IsAllIn = false
	s.CurrentBet = 0
	s.TotalBet = 0
}

// PlaceBet places a bet of the specified amount.
// Returns the actual amount bet (which may be less if the player doesn't have enough chips)
// and whether the player went all-in.
func (s *Seat) PlaceBet(amount int) (int, bool) {
	if !s.IsOccupied || !s.IsActive || s.HasFolded || s.IsAllIn {
		return 0, false
	}

	// If player doesn't have enough chips, they go all-in
	if s.ChipStack <= amount {
		amount = s.ChipStack
		s.IsAllIn = true
	}

	s.ChipStack -= amount
	s.CurrentBet += amount
	s.TotalBet += amount

	// Double-check if the player is all-in
	if s.ChipStack == 0 {
		s.IsAllIn = true
	}

	return amount, s.IsAllIn
}

// Fold marks the player as folded.
func (s *Seat) Fold() {
	if s.IsOccupied && s.IsActive && !s.HasFolded {
		s.HasFolded = true
	}
}

// AddChips adds chips to the player's stack.
func (s *Seat) AddChips(amount int) {
	if s.IsOccupied {
		s.ChipStack += amount
	}
}

// CanBet returns whether the player can place a bet.
func (s *Seat) CanBet() bool {
	return s.IsOccupied && s.IsActive && !s.HasFolded && !s.IsAllIn && s.ChipStack > 0
}

// ResetBetForRound resets the current bet for a new betting round.
func (s *Seat) ResetBetForRound() {
	s.CurrentBet = 0
}
