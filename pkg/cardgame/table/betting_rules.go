package table

import "errors"

// BettingRules defines the interface for different betting structures.
type BettingRules interface {
	// ValidateBet checks if a bet is valid according to the betting rules.
	// It returns the validated bet amount and an error if the bet is invalid.
	ValidateBet(currentBet, lastRaise, playerChips, betAmount, potSize int) (int, error)

	// GetMinRaise returns the minimum raise amount.
	GetMinRaise(currentBet, lastRaise int) int

	// GetMaxBet returns the maximum bet amount a player can make.
	GetMaxBet(playerChips, currentBet, potSize int) int

	// GetStructureType returns the betting structure type.
	GetStructureType() BettingStructure
}

// NoLimitRules implements the BettingRules interface for No Limit betting.
type NoLimitRules struct{}

// ValidateBet validates a bet according to No Limit rules.
func (r *NoLimitRules) ValidateBet(currentBet, lastRaise, playerChips, betAmount, potSize int) (int, error) {
	// In No Limit, a player can bet any amount up to their chip stack
	if betAmount > playerChips {
		// If the player is trying to bet more than they have, they're going all-in
		return playerChips, nil
	}

	// If this is a raise, it must be at least the size of the last raise
	if betAmount > currentBet {
		minRaise := r.GetMinRaise(currentBet, lastRaise)
		if betAmount < minRaise && betAmount < playerChips {
			// Only enforce minimum raise if the player has enough chips
			return 0, errors.New("raise must be at least the size of the last raise")
		}
	}

	return betAmount, nil
}

// GetMinRaise returns the minimum raise amount for No Limit.
func (r *NoLimitRules) GetMinRaise(currentBet, lastRaise int) int {
	// In No Limit, the minimum raise is the size of the last raise or the big blind
	if lastRaise > 0 {
		return currentBet + lastRaise
	}
	// If there was no previous raise, the minimum raise is the big blind
	// This is a simplification - in a real implementation, we would need to
	// pass the big blind amount as a parameter
	return currentBet * 2
}

// GetMaxBet returns the maximum bet amount for No Limit.
func (r *NoLimitRules) GetMaxBet(playerChips, currentBet, potSize int) int {
	// In No Limit, a player can bet any amount up to their chip stack
	return playerChips
}

// GetStructureType returns the betting structure type.
func (r *NoLimitRules) GetStructureType() BettingStructure {
	return NoLimit
}

// PotLimitRules implements the BettingRules interface for Pot Limit betting.
type PotLimitRules struct{}

// ValidateBet validates a bet according to Pot Limit rules.
func (r *PotLimitRules) ValidateBet(currentBet, lastRaise, playerChips, betAmount, potSize int) (int, error) {
	// In Pot Limit, a player can bet any amount up to the pot size
	maxBet := r.GetMaxBet(playerChips, currentBet, potSize)
	if betAmount > maxBet {
		// If the player is trying to bet more than allowed, cap it at the maximum
		return maxBet, nil
	}

	// If this is a raise, it must be at least the size of the last raise
	if betAmount > currentBet {
		minRaise := r.GetMinRaise(currentBet, lastRaise)
		if betAmount < minRaise && betAmount < playerChips {
			// Only enforce minimum raise if the player has enough chips
			return 0, errors.New("raise must be at least the size of the last raise")
		}
	}

	return betAmount, nil
}

// GetMinRaise returns the minimum raise amount for Pot Limit.
func (r *PotLimitRules) GetMinRaise(currentBet, lastRaise int) int {
	// In Pot Limit, the minimum raise is the size of the last raise or the big blind
	if lastRaise > 0 {
		return currentBet + lastRaise
	}
	// If there was no previous raise, the minimum raise is the big blind
	return currentBet * 2
}

// GetMaxBet returns the maximum bet amount for Pot Limit.
func (r *PotLimitRules) GetMaxBet(playerChips, currentBet, potSize int) int {
	// In Pot Limit, a player can bet any amount up to the pot size plus the current bet
	// The formula is: pot size + current bet
	maxBet := potSize + currentBet

	// If the player doesn't have enough chips, they can only bet what they have
	if maxBet > playerChips {
		maxBet = playerChips
	}

	return maxBet
}

// GetStructureType returns the betting structure type.
func (r *PotLimitRules) GetStructureType() BettingStructure {
	return PotLimit
}

// FixedLimitRules implements the BettingRules interface for Fixed Limit betting.
type FixedLimitRules struct {
	// SmallBet is the bet size for early betting rounds
	SmallBet int

	// BigBet is the bet size for later betting rounds
	BigBet int

	// MaxRaises is the maximum number of raises allowed per betting round
	MaxRaises int

	// CurrentRound is the current betting round (0-indexed)
	CurrentRound int
}

// NewFixedLimitRules creates a new FixedLimitRules instance.
func NewFixedLimitRules(smallBet, bigBet, maxRaises int) *FixedLimitRules {
	return &FixedLimitRules{
		SmallBet:     smallBet,
		BigBet:       bigBet,
		MaxRaises:    maxRaises,
		CurrentRound: 0,
	}
}

// ValidateBet validates a bet according to Fixed Limit rules.
func (r *FixedLimitRules) ValidateBet(currentBet, lastRaise, playerChips, betAmount, potSize int) (int, error) {
	// In Fixed Limit, a player can only bet the fixed amount
	fixedBet := r.GetFixedBetSize()

	// If the player doesn't have enough chips, they go all-in
	if betAmount > playerChips {
		return playerChips, nil
	}

	// If the player is trying to bet something other than the fixed amount or 0 (check/call)
	if betAmount != fixedBet && betAmount != 0 && betAmount != playerChips {
		return 0, errors.New("bet amount must be the fixed limit")
	}

	return betAmount, nil
}

// GetMinRaise returns the minimum raise amount for Fixed Limit.
func (r *FixedLimitRules) GetMinRaise(currentBet, lastRaise int) int {
	// In Fixed Limit, the raise amount is fixed
	return r.GetFixedBetSize()
}

// GetMaxBet returns the maximum bet amount for Fixed Limit.
func (r *FixedLimitRules) GetMaxBet(playerChips, currentBet, potSize int) int {
	// In Fixed Limit, the bet amount is fixed
	fixedBet := r.GetFixedBetSize()
	if fixedBet > playerChips {
		fixedBet = playerChips
	}
	return fixedBet
}

// GetStructureType returns the betting structure type.
func (r *FixedLimitRules) GetStructureType() BettingStructure {
	return FixedLimit
}

// GetFixedBetSize returns the fixed bet size for the current round.
func (r *FixedLimitRules) GetFixedBetSize() int {
	// In Fixed Limit, the bet size depends on the betting round
	// Typically, the first two rounds use the small bet, and the last two rounds use the big bet
	if r.CurrentRound < 2 {
		return r.SmallBet
	}
	return r.BigBet
}

// SetCurrentRound sets the current betting round.
func (r *FixedLimitRules) SetCurrentRound(round int) {
	r.CurrentRound = round
}

// NewBettingRules creates a new BettingRules instance based on the specified structure.
func NewBettingRules(structure BettingStructure, smallBet, bigBet, maxRaises int) BettingRules {
	switch structure {
	case NoLimit:
		return &NoLimitRules{}
	case PotLimit:
		return &PotLimitRules{}
	case FixedLimit:
		return NewFixedLimitRules(smallBet, bigBet, maxRaises)
	default:
		return &NoLimitRules{}
	}
}
