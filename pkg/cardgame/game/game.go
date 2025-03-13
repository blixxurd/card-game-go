// Package game defines interfaces and implementations for card games.
package game

import (
	"github.com/blixxurd/card-game-go/pkg/cardgame/deck"
	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
)

// GameState represents the current state of a game
type GameState string

const (
	// GameStateInitialized indicates the game has been initialized but not started
	GameStateInitialized GameState = "initialized"

	// GameStateInProgress indicates the game is currently in progress
	GameStateInProgress GameState = "in_progress"

	// GameStateCompleted indicates the game has been completed
	GameStateCompleted GameState = "completed"

	// GameStatePaused indicates the game has been paused
	GameStatePaused GameState = "paused"
)

// Action represents a player action in a game
type Action interface {
	// Type returns the type of action
	Type() string

	// PlayerID returns the ID of the player performing the action
	PlayerID() string

	// Params returns the parameters for the action
	Params() map[string]interface{}
}

// Game is the base interface that all game implementations must satisfy.
// It provides methods for managing a card game.
type Game interface {
	// ID returns a unique identifier for the game
	ID() string

	// Name returns the name of the game
	Name() string

	// Initialize sets up the game with the provided configuration
	Initialize(config map[string]interface{}) error

	// Start begins the game
	Start() error

	// AddPlayer adds a player to the game
	AddPlayer(p player.Player) error

	// RemovePlayer removes a player from the game
	RemovePlayer(playerID string) error

	// Players returns all players in the game
	Players() []player.Player

	// CurrentPlayer returns the player whose turn it is
	CurrentPlayer() player.Player

	// NextTurn advances to the next player's turn
	NextTurn() error

	// State returns the current state of the game
	State() GameState

	// Deck returns the main deck used in the game
	Deck() deck.Deck

	// ProcessAction processes a player action
	ProcessAction(action Action) error

	// IsValidAction checks if an action is valid in the current game state
	IsValidAction(action Action) bool

	// AllowedActions returns the actions that are currently allowed
	AllowedActions(playerID string) []string

	// Winners returns the winners of the game (if any)
	Winners() []player.Player

	// Reset resets the game to its initial state
	Reset() error

	// Data returns custom game data as a map
	Data() map[string]interface{}

	// SetData sets custom game data
	SetData(key string, value interface{})
}

// GameType represents the type of game (e.g., poker)
type GameType string

const (
	// PokerGame represents a poker game
	PokerGame GameType = "poker"
)

// GameFactory is an interface for creating games of different types
type GameFactory interface {
	// CreateGame creates a new game based on the provided parameters
	CreateGame(params map[string]interface{}) (Game, error)

	// Type returns the type of games this factory creates
	Type() GameType
}
