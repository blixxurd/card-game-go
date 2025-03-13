# Implementing a New Poker Variant

This guide demonstrates how to add a new poker variant (e.g., Omaha) to the card game framework.

## Step-by-Step Implementation

### 1. Create a new package

Create a new package for your variant:

```
pkg/cardgame/poker/omaha/
```

### 2. Implement the core game logic

```go
// omaha.go
package omaha

import (
    "github.com/blixxurd/card-game-go/pkg/cardgame/deck"
    "github.com/blixxurd/card-game-go/pkg/cardgame/player"
    "github.com/blixxurd/card-game-go/pkg/cardgame/table"
    "github.com/blixxurd/card-game-go/pkg/cardgame/card"
)

// OmahaGame represents a game of Omaha Hold'em
type OmahaGame struct {
    id            string
    name          string
    players       []player.Player
    table         *table.Table
    mainDeck      *deck.StandardDeck
    communityCards []card.Card
    // Other necessary fields
}

// NewOmahaGame creates a new Omaha game
func NewOmahaGame(id, name string, smallBlind, bigBlind int) *OmahaGame {
    t := table.NewTable(smallBlind, bigBlind, table.NoLimit)
    return &OmahaGame{
        id:            id,
        name:          name,
        players:       make([]player.Player, 0),
        table:         t,
        mainDeck:      deck.NewStandardDeck(),
        communityCards: make([]card.Card, 0, 5),
    }
}

// Implement necessary methods for game logic
func (g *OmahaGame) ID() string {
    return g.id
}

func (g *OmahaGame) Name() string {
    return g.name
}

// Add more methods as needed...
```

### 3. Create a variant implementation

Implement the `PokerVariant` interface:

```go
// omaha_variant.go
package omaha

import (
    "github.com/blixxurd/card-game-go/pkg/cardgame/poker"
    "github.com/blixxurd/card-game-go/pkg/cardgame/table"
)

// OmahaVariant implements the PokerVariant interface for Omaha
type OmahaVariant struct{}

// NewOmahaVariant creates a new OmahaVariant
func NewOmahaVariant() *OmahaVariant {
    return &OmahaVariant{}
}

// Name returns the name of the poker variant
func (v *OmahaVariant) Name() string {
    return "Omaha Hold'em"
}

// Description returns a description of the poker variant
func (v *OmahaVariant) Description() string {
    return "Omaha Hold'em is a community card poker variant where each player is dealt four private cards and must use exactly two of them with three of the five community cards."
}

// CreateGame creates a new instance of the poker variant
func (v *OmahaVariant) CreateGame(id, name string, smallBlind, bigBlind int, bettingStructure table.BettingStructure) (poker.PokerGame, error) {
    omahaGame := NewOmahaGame(id, name, smallBlind, bigBlind)
    return NewOmahaGameAdapter(omahaGame), nil
}

// MaxPlayers returns the maximum number of players allowed
func (v *OmahaVariant) MaxPlayers() int {
    return 9
}

// MinPlayers returns the minimum number of players required
func (v *OmahaVariant) MinPlayers() int {
    return 2
}

// HandSize returns the number of cards in a player's hand
func (v *OmahaVariant) HandSize() int {
    return 4  // Omaha uses 4 hole cards
}

// CommunityCardCount returns the number of community cards used
func (v *OmahaVariant) CommunityCardCount() int {
    return 5
}

// BettingRounds returns the number of betting rounds
func (v *OmahaVariant) BettingRounds() int {
    return 4
}

// StateSequence returns the sequence of game states
func (v *OmahaVariant) StateSequence() []poker.PokerGameState {
    return []poker.PokerGameState{
        "initial",
        "pre_flop",
        "flop",
        "turn",
        "river",
        "showdown",
        "complete",
    }
}
```

### 4. Create an adapter

Implement an adapter that satisfies the `PokerGame` interface:

```go
// omaha_adapter.go
package omaha

import (
    "github.com/blixxurd/card-game-go/pkg/cardgame/card"
    "github.com/blixxurd/card-game-go/pkg/cardgame/deck"
    "github.com/blixxurd/card-game-go/pkg/cardgame/game"
    "github.com/blixxurd/card-game-go/pkg/cardgame/player"
    "github.com/blixxurd/card-game-go/pkg/cardgame/poker"
    "github.com/blixxurd/card-game-go/pkg/cardgame/table"
)

// OmahaGameAdapter adapts an OmahaGame to implement the PokerGame interface
type OmahaGameAdapter struct {
    omahaGame *OmahaGame
}

// NewOmahaGameAdapter creates a new OmahaGameAdapter
func NewOmahaGameAdapter(omahaGame *OmahaGame) *OmahaGameAdapter {
    return &OmahaGameAdapter{
        omahaGame: omahaGame,
    }
}

// Implement all methods required by the PokerGame interface
func (g *OmahaGameAdapter) ID() string {
    return g.omahaGame.ID()
}

func (g *OmahaGameAdapter) Name() string {
    return g.omahaGame.Name()
}

// Implement remaining methods...
```

### 5. Register your variant

Register your variant with the `PokerGameFactory`:

```go
// In your application initialization (e.g., main.go)
package main

import (
    "github.com/blixxurd/card-game-go/pkg/cardgame/poker"
    "github.com/blixxurd/card-game-go/pkg/cardgame/poker/holdem"
    "github.com/blixxurd/card-game-go/pkg/cardgame/poker/omaha"
    "github.com/blixxurd/card-game-go/pkg/cardgame/table"
)

func main() {
    // Create a factory and register variants
    factory := poker.NewPokerGameFactory()
    factory.RegisterVariant(holdem.NewHoldemVariant())
    factory.RegisterVariant(omaha.NewOmahaVariant())

    // Create a game using the factory
    game, err := factory.CreatePokerGame("Omaha Hold'em", "game1", "My Omaha Game", 5, 10, table.NoLimit)
    if err != nil {
        panic(err)
    }

    // Use the game...
}
```

### 6. Create tests

Write tests for your new variant:

```go
// omaha_test.go
package omaha_test

import (
    "testing"

    "github.com/blixxurd/card-game-go/pkg/cardgame/player"
    "github.com/blixxurd/card-game-go/pkg/cardgame/poker/omaha"
    "github.com/blixxurd/card-game-go/pkg/cardgame/table"
)

// TestOmahaGameCreation tests the creation of an Omaha game
func TestOmahaGameCreation(t *testing.T) {
    // Create a new variant
    variant := omaha.NewOmahaVariant()

    // Check variant properties
    if variant.Name() != "Omaha Hold'em" {
        t.Errorf("Expected name to be 'Omaha Hold'em', got '%s'", variant.Name())
    }

    if variant.HandSize() != 4 {
        t.Errorf("Expected hand size to be 4, got %d", variant.HandSize())
    }

    // Create a game
    game, err := variant.CreateGame("test", "Test Game", 5, 10, table.NoLimit)
    if err != nil {
        t.Errorf("Failed to create game: %v", err)
    }

    // Check game properties
    if game.ID() != "test" {
        t.Errorf("Expected game ID 'test', got '%s'", game.ID())
    }

    // Add more tests...
}
```

## Key Differences for Omaha

When implementing Omaha, pay special attention to these key differences from Texas Hold'em:

1. **Hand Size**: Players receive 4 hole cards instead of 2.
2. **Hand Evaluation**: Players must use exactly 2 hole cards and 3 community cards.
3. **Betting Structure**: While the betting rounds are the same as Texas Hold'em, the strategy is different due to the hand composition rules.

## Best Practices

1. **Follow the interfaces**: Ensure your implementation satisfies all methods required by the interfaces.
2. **Maintain separation of concerns**: Keep game logic separate from UI or network code.
3. **Use composition**: Leverage existing components like deck and card implementations.
4. **Write tests**: Create comprehensive tests for your game implementation.
5. **Document your code**: Add clear comments and documentation for your implementation. 