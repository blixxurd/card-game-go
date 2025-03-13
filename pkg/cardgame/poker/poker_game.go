// Package poker provides interfaces and implementations for poker games.
package poker

import (
	"errors"

	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
	"github.com/blixxurd/card-game-go/pkg/cardgame/game"
	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
	"github.com/blixxurd/card-game-go/pkg/cardgame/table"
)

// Error constants
var (
	// ErrInvalidVariant is returned when an invalid poker variant is specified
	ErrInvalidVariant = errors.New("invalid poker variant")
	// ErrInvalidGameParams is returned when invalid game parameters are provided
	ErrInvalidGameParams = errors.New("invalid game parameters")
)

// PokerGameState represents the state of a poker game
type PokerGameState string

// PokerActionType represents the type of action a player can take in a poker game
type PokerActionType string

const (
	// ActionCheck represents a check action
	ActionCheck PokerActionType = "check"
	// ActionBet represents a bet action
	ActionBet PokerActionType = "bet"
	// ActionCall represents a call action
	ActionCall PokerActionType = "call"
	// ActionRaise represents a raise action
	ActionRaise PokerActionType = "raise"
	// ActionFold represents a fold action
	ActionFold PokerActionType = "fold"
	// ActionAllIn represents an all-in action
	ActionAllIn PokerActionType = "all_in"
)

// PokerAction represents an action in a poker game
type PokerAction interface {
	game.Action
	// Amount returns the amount of the bet or raise
	Amount() int
	// ActionType returns the type of poker action
	ActionType() PokerActionType
}

// PokerGame is an interface that extends the base Game interface with poker-specific methods
type PokerGame interface {
	game.Game
	// Table returns the table used in the game
	Table() *table.Table
	// CommunityCards returns the community cards
	CommunityCards() []card.Card
	// ActivePlayers returns the players who are still in the hand
	ActivePlayers() []player.Player
	// Pot returns the current pot size
	Pot() int
	// CurrentBet returns the current bet amount
	CurrentBet() int
	// MinRaise returns the minimum raise amount
	MinRaise() int
	// MaxBet returns the maximum bet amount for a player
	MaxBet(playerID string) (int, error)
	// ProcessPokerAction processes a poker-specific action
	ProcessPokerAction(action PokerAction) error
	// IsValidPokerAction checks if a poker action is valid
	IsValidPokerAction(action PokerAction) bool
	// AllowedPokerActions returns the poker actions that are currently allowed for a player
	AllowedPokerActions(playerID string) []PokerActionType
	// BettingRound returns the current betting round
	BettingRound() int
	// AdvanceBettingRound advances to the next betting round
	AdvanceBettingRound() error
	// DealCommunityCards deals community cards according to the current game state
	DealCommunityCards() error
	// EvaluateHands evaluates the hands of all active players
	EvaluateHands() error
	// AwardPots awards the pots to the winners
	AwardPots() error
}

// PokerVariant defines the interface for different poker variants
type PokerVariant interface {
	// Name returns the name of the poker variant
	Name() string
	// Description returns a description of the poker variant
	Description() string
	// CreateGame creates a new instance of the poker variant
	CreateGame(id, name string, smallBlind, bigBlind int, bettingStructure table.BettingStructure) (PokerGame, error)
	// MaxPlayers returns the maximum number of players allowed
	MaxPlayers() int
	// MinPlayers returns the minimum number of players required
	MinPlayers() int
	// HandSize returns the number of cards in a player's hand
	HandSize() int
	// CommunityCardCount returns the number of community cards used
	CommunityCardCount() int
	// BettingRounds returns the number of betting rounds
	BettingRounds() int
	// StateSequence returns the sequence of game states
	StateSequence() []PokerGameState
}

// PokerGameFactory is a factory for creating poker games
type PokerGameFactory struct {
	variants map[string]PokerVariant
}

// NewPokerGameFactory creates a new PokerGameFactory
func NewPokerGameFactory() *PokerGameFactory {
	return &PokerGameFactory{
		variants: make(map[string]PokerVariant),
	}
}

// RegisterVariant registers a poker variant with the factory
func (f *PokerGameFactory) RegisterVariant(variant PokerVariant) {
	f.variants[variant.Name()] = variant
}

// CreatePokerGame creates a new poker game of the specified variant
func (f *PokerGameFactory) CreatePokerGame(variant string, id, name string, smallBlind, bigBlind int, bettingStructure table.BettingStructure) (PokerGame, error) {
	if v, ok := f.variants[variant]; ok {
		return v.CreateGame(id, name, smallBlind, bigBlind, bettingStructure)
	}
	return nil, ErrInvalidVariant
}

// AvailableVariants returns the names of all registered poker variants
func (f *PokerGameFactory) AvailableVariants() []string {
	variants := make([]string, 0, len(f.variants))
	for name := range f.variants {
		variants = append(variants, name)
	}
	return variants
}

// Type implements the GameFactory interface
func (f *PokerGameFactory) Type() game.GameType {
	return game.PokerGame
}

// CreateGame implements the GameFactory interface
func (f *PokerGameFactory) CreateGame(params map[string]interface{}) (game.Game, error) {
	variant, ok := params["variant"].(string)
	if !ok {
		return nil, ErrInvalidGameParams
	}

	id, ok := params["id"].(string)
	if !ok {
		return nil, ErrInvalidGameParams
	}

	name, ok := params["name"].(string)
	if !ok {
		return nil, ErrInvalidGameParams
	}

	smallBlind, ok := params["small_blind"].(int)
	if !ok {
		return nil, ErrInvalidGameParams
	}

	bigBlind, ok := params["big_blind"].(int)
	if !ok {
		return nil, ErrInvalidGameParams
	}

	bettingStructureInt, ok := params["betting_structure"].(int)
	if !ok {
		bettingStructureInt = int(table.NoLimit)
	}

	bettingStructure := table.BettingStructure(bettingStructureInt)

	pokerGame, err := f.CreatePokerGame(variant, id, name, smallBlind, bigBlind, bettingStructure)
	if err != nil {
		return nil, err
	}

	return pokerGame, nil
}
