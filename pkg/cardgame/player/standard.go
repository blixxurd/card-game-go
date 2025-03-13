// Package player defines interfaces and implementations for players in various card games.
package player

import (
	"fmt"
	"sync"

	"github.com/blixxurd/card-game-go/pkg/cardgame/hand"
)

// StandardPlayer represents a basic player in a card game
type StandardPlayer struct {
	id        string
	name      string
	hand      hand.Hand
	score     int
	active    bool
	data      map[string]interface{}
	dataMutex sync.RWMutex
}

// NewStandardPlayer creates a new player with the given ID and name
func NewStandardPlayer(id, name string) *StandardPlayer {
	return &StandardPlayer{
		id:     id,
		name:   name,
		hand:   hand.NewStandardHand(),
		score:  0,
		active: true,
		data:   make(map[string]interface{}),
	}
}

// ID returns a unique identifier for the player
func (p *StandardPlayer) ID() string {
	return p.id
}

// Name returns the player's name
func (p *StandardPlayer) Name() string {
	return p.name
}

// SetName sets the player's name
func (p *StandardPlayer) SetName(name string) {
	p.name = name
}

// Hand returns the player's current hand
func (p *StandardPlayer) Hand() hand.Hand {
	return p.hand
}

// SetHand sets the player's hand
func (p *StandardPlayer) SetHand(h hand.Hand) {
	p.hand = h
}

// Score returns the player's current score
func (p *StandardPlayer) Score() int {
	return p.score
}

// SetScore sets the player's score
func (p *StandardPlayer) SetScore(score int) {
	p.score = score
}

// AddToScore adds to the player's score
func (p *StandardPlayer) AddToScore(points int) {
	p.score += points
}

// IsActive returns whether the player is active in the current game
func (p *StandardPlayer) IsActive() bool {
	return p.active
}

// SetActive sets whether the player is active
func (p *StandardPlayer) SetActive(active bool) {
	p.active = active
}

// Data returns custom player data as a map
func (p *StandardPlayer) Data() map[string]interface{} {
	p.dataMutex.RLock()
	defer p.dataMutex.RUnlock()

	// Create a copy to avoid concurrent modification issues
	result := make(map[string]interface{}, len(p.data))
	for k, v := range p.data {
		result[k] = v
	}

	return result
}

// SetData sets custom player data
func (p *StandardPlayer) SetData(key string, value interface{}) {
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()

	p.data[key] = value
}

// Clone creates a deep copy of the player
func (p *StandardPlayer) Clone() Player {
	p.dataMutex.RLock()
	defer p.dataMutex.RUnlock()

	clone := &StandardPlayer{
		id:     p.id,
		name:   p.name,
		hand:   p.hand.Clone(),
		score:  p.score,
		active: p.active,
		data:   make(map[string]interface{}, len(p.data)),
	}

	// Copy the data map
	for k, v := range p.data {
		clone.data[k] = v
	}

	return clone
}

// StandardPlayerFactory creates standard players
type StandardPlayerFactory struct{}

// CreatePlayer creates a new standard player
func (f *StandardPlayerFactory) CreatePlayer(params map[string]interface{}) (Player, error) {
	id, ok := params["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid id parameter")
	}

	name, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid name parameter")
	}

	player := NewStandardPlayer(id, name)

	// Set initial score if provided
	if score, ok := params["score"].(int); ok {
		player.SetScore(score)
	}

	// Set active status if provided
	if active, ok := params["active"].(bool); ok {
		player.SetActive(active)
	}

	return player, nil
}

// Type returns the type of players this factory creates
func (f *StandardPlayerFactory) Type() PlayerType {
	return HumanPlayer
}
