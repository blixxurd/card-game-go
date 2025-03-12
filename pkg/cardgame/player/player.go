// Package player defines interfaces and implementations for players in various card games.
package player

import (
	"github.com/blixxurd/card-game-go/pkg/cardgame/hand"
)

// Player is the base interface that all player implementations must satisfy.
// It provides methods for managing a player in a card game.
type Player interface {
	// ID returns a unique identifier for the player
	ID() string

	// Name returns the player's name
	Name() string

	// SetName sets the player's name
	SetName(name string)

	// Hand returns the player's current hand
	Hand() hand.Hand

	// SetHand sets the player's hand
	SetHand(h hand.Hand)

	// Score returns the player's current score
	Score() int

	// SetScore sets the player's score
	SetScore(score int)

	// AddToScore adds to the player's score
	AddToScore(points int)

	// IsActive returns whether the player is active in the current game
	IsActive() bool

	// SetActive sets whether the player is active
	SetActive(active bool)

	// Data returns custom player data as a map
	Data() map[string]interface{}

	// SetData sets custom player data
	SetData(key string, value interface{})

	// Clone creates a deep copy of the player
	Clone() Player
}

// PlayerType represents the type of player (e.g., human, AI, etc.)
type PlayerType string

const (
	// HumanPlayer represents a human player
	HumanPlayer PlayerType = "human"

	// AIPlayer represents an AI player
	AIPlayer PlayerType = "ai"

	// NetworkPlayer represents a player connected over a network
	NetworkPlayer PlayerType = "network"
)

// PlayerFactory is an interface for creating players of different types
type PlayerFactory interface {
	// CreatePlayer creates a new player based on the provided parameters
	CreatePlayer(params map[string]interface{}) (Player, error)

	// Type returns the type of players this factory creates
	Type() PlayerType
}
