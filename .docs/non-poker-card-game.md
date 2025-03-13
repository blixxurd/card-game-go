# Implementing a Non-Poker Card Game

This guide demonstrates how to add a non-poker card game (e.g., Blackjack) to the card game framework.

## Step-by-Step Implementation

### 1. Define a new game interface

Create a new interface that extends the base `Game` interface:

```go
// pkg/cardgame/blackjack/blackjack_game.go
package blackjack

import (
    "github.com/blixxurd/card-game-go/pkg/cardgame/card"
    "github.com/blixxurd/card-game-go/pkg/cardgame/game"
    "github.com/blixxurd/card-game-go/pkg/cardgame/player"
)

// BlackjackActionType represents the type of action in blackjack
type BlackjackActionType string

const (
    ActionHit    BlackjackActionType = "hit"
    ActionStand  BlackjackActionType = "stand"
    ActionDouble BlackjackActionType = "double"
    ActionSplit  BlackjackActionType = "split"
)

// BlackjackAction represents an action in a blackjack game
type BlackjackAction interface {
    game.Action
    ActionType() BlackjackActionType
}

// BlackjackGame is an interface for blackjack games
type BlackjackGame interface {
    game.Game
    // Add blackjack-specific methods
    DealerHand() []card.Card
    DealerUpCard() card.Card
    IsBlackjack(playerID string) bool
    IsBusted(playerID string) bool
    ProcessBlackjackAction(action BlackjackAction) error
    AllowedBlackjackActions(playerID string) []BlackjackActionType
    // etc.
}
```

### 2. Create a factory for your game type

```go
// pkg/cardgame/blackjack/blackjack_factory.go
package blackjack

import (
    "errors"
    "github.com/blixxurd/card-game-go/pkg/cardgame/game"
)

// BlackjackGameFactory creates blackjack games
type BlackjackGameFactory struct {
    // Configuration options
}

// NewBlackjackGameFactory creates a new BlackjackGameFactory
func NewBlackjackGameFactory() *BlackjackGameFactory {
    return &BlackjackGameFactory{}
}

// CreateGame creates a new blackjack game
func (f *BlackjackGameFactory) CreateGame(params map[string]interface{}) (game.Game, error) {
    id, ok := params["id"].(string)
    if !ok {
        return nil, errors.New("invalid game parameters: missing id")
    }

    name, ok := params["name"].(string)
    if !ok {
        return nil, errors.New("invalid game parameters: missing name")
    }

    minBet, ok := params["min_bet"].(int)
    if !ok {
        minBet = 5 // Default minimum bet
    }

    return NewStandardBlackjackGame(id, name, minBet), nil
}

// Type returns the type of games this factory creates
func (f *BlackjackGameFactory) Type() game.GameType {
    return "blackjack"
}
```

### 3. Implement the game logic

```go
// pkg/cardgame/blackjack/standard_blackjack.go
package blackjack

import (
    "errors"
    "github.com/blixxurd/card-game-go/pkg/cardgame/card"
    "github.com/blixxurd/card-game-go/pkg/cardgame/deck"
    "github.com/blixxurd/card-game-go/pkg/cardgame/game"
    "github.com/blixxurd/card-game-go/pkg/cardgame/player"
)

// StandardBlackjackGame implements the BlackjackGame interface
type StandardBlackjackGame struct {
    id          string
    name        string
    players     []player.Player
    dealerHand  []card.Card
    playerHands map[string][]card.Card
    deck        *deck.StandardDeck
    minBet      int
    bets        map[string]int
    state       game.GameState
    winners     []player.Player
    customData  map[string]interface{}
}

// NewStandardBlackjackGame creates a new standard blackjack game
func NewStandardBlackjackGame(id, name string, minBet int) *StandardBlackjackGame {
    return &StandardBlackjackGame{
        id:          id,
        name:        name,
        players:     make([]player.Player, 0),
        dealerHand:  make([]card.Card, 0),
        playerHands: make(map[string][]card.Card),
        deck:        deck.NewStandardDeck(),
        minBet:      minBet,
        bets:        make(map[string]int),
        state:       game.GameStateInitialized,
        customData:  make(map[string]interface{}),
    }
}

// Implement the game.Game interface
func (g *StandardBlackjackGame) ID() string {
    return g.id
}

func (g *StandardBlackjackGame) Name() string {
    return g.name
}

func (g *StandardBlackjackGame) Initialize(config map[string]interface{}) error {
    // Initialize the game with the provided configuration
    return nil
}

func (g *StandardBlackjackGame) Start() error {
    if len(g.players) < 1 {
        return errors.New("not enough players")
    }

    g.deck.Shuffle()

    // Deal two cards to each player
    for _, p := range g.players {
        cards, err := g.deck.DrawMany(2)
        if err != nil {
            return err
        }
        g.playerHands[p.ID()] = cards
    }

    // Deal two cards to the dealer
    dealerCards, err := g.deck.DrawMany(2)
    if err != nil {
        return err
    }
    g.dealerHand = dealerCards

    g.state = game.GameStateInProgress
    return nil
}

// Implement BlackjackGame-specific methods
func (g *StandardBlackjackGame) DealerHand() []card.Card {
    return g.dealerHand
}

func (g *StandardBlackjackGame) DealerUpCard() card.Card {
    if len(g.dealerHand) > 0 {
        return g.dealerHand[0]
    }
    return nil
}

func (g *StandardBlackjackGame) IsBlackjack(playerID string) bool {
    hand, ok := g.playerHands[playerID]
    if !ok || len(hand) != 2 {
        return false
    }

    // Check if the hand is an Ace and a 10-value card
    hasAce := false
    hasTenValue := false

    for _, c := range hand {
        if c.Rank() == 1 { // Ace
            hasAce = true
        } else if c.Rank() >= 10 { // 10, J, Q, K
            hasTenValue = true
        }
    }

    return hasAce && hasTenValue
}

func (g *StandardBlackjackGame) IsBusted(playerID string) bool {
    hand, ok := g.playerHands[playerID]
    if !ok {
        return false
    }

    // Calculate hand value
    value := g.calculateHandValue(hand)
    return value > 21
}

func (g *StandardBlackjackGame) calculateHandValue(hand []card.Card) int {
    value := 0
    aces := 0

    for _, c := range hand {
        rank := c.Rank()
        if rank == 1 { // Ace
            aces++
            value += 11
        } else if rank >= 10 { // 10, J, Q, K
            value += 10
        } else {
            value += rank
        }
    }

    // Adjust for aces if needed
    for aces > 0 && value > 21 {
        value -= 10
        aces--
    }

    return value
}

func (g *StandardBlackjackGame) ProcessBlackjackAction(action BlackjackAction) error {
    if g.state != game.GameStateInProgress {
        return errors.New("game is not in progress")
    }

    playerID := action.PlayerID()
    hand, ok := g.playerHands[playerID]
    if !ok {
        return errors.New("player not found")
    }

    switch action.ActionType() {
    case ActionHit:
        card, err := g.deck.Draw()
        if err != nil {
            return err
        }
        g.playerHands[playerID] = append(hand, card)

        // Check if player busted
        if g.IsBusted(playerID) {
            // Handle player bust
        }

    case ActionStand:
        // Player stands, move to next player or dealer's turn
        
    case ActionDouble:
        // Double the bet and draw one more card
        bet, ok := g.bets[playerID]
        if !ok {
            return errors.New("no bet found for player")
        }
        g.bets[playerID] = bet * 2

        card, err := g.deck.Draw()
        if err != nil {
            return err
        }
        g.playerHands[playerID] = append(hand, card)

    case ActionSplit:
        // Handle splitting pairs
        // ...

    default:
        return errors.New("invalid action type")
    }

    return nil
}

func (g *StandardBlackjackGame) AllowedBlackjackActions(playerID string) []BlackjackActionType {
    hand, ok := g.playerHands[playerID]
    if !ok {
        return nil
    }

    actions := []BlackjackActionType{ActionHit, ActionStand}

    // Check if doubling is allowed (only on first two cards)
    if len(hand) == 2 {
        actions = append(actions, ActionDouble)
    }

    // Check if splitting is allowed (only with a pair)
    if len(hand) == 2 && hand[0].Rank() == hand[1].Rank() {
        actions = append(actions, ActionSplit)
    }

    return actions
}

// Implement remaining methods from the game.Game interface
// ...
```

### 4. Create a concrete action implementation

```go
// pkg/cardgame/blackjack/blackjack_action.go
package blackjack

import (
    "fmt"
)

// StandardBlackjackAction implements the BlackjackAction interface
type StandardBlackjackAction struct {
    actionType BlackjackActionType
    playerID   string
    params     map[string]interface{}
}

// NewBlackjackAction creates a new blackjack action
func NewBlackjackAction(actionType BlackjackActionType, playerID string) *StandardBlackjackAction {
    return &StandardBlackjackAction{
        actionType: actionType,
        playerID:   playerID,
        params:     make(map[string]interface{}),
    }
}

// Type returns the type of the action
func (a *StandardBlackjackAction) Type() string {
    return string(a.actionType)
}

// PlayerID returns the ID of the player performing the action
func (a *StandardBlackjackAction) PlayerID() string {
    return a.playerID
}

// Params returns the parameters for the action
func (a *StandardBlackjackAction) Params() map[string]interface{} {
    return a.params
}

// ActionType returns the type of blackjack action
func (a *StandardBlackjackAction) ActionType() BlackjackActionType {
    return a.actionType
}

// SetParam sets a parameter for the action
func (a *StandardBlackjackAction) SetParam(key string, value interface{}) {
    a.params[key] = value
}

// String returns a string representation of the action
func (a *StandardBlackjackAction) String() string {
    return fmt.Sprintf("%s by player %s", a.actionType, a.playerID)
}
```

### 5. Register your game factory

If you have a central registry for game factories:

```go
// In your application initialization
gameRegistry := NewGameRegistry()
gameRegistry.RegisterFactory("blackjack", blackjack.NewBlackjackGameFactory())
```

### 6. Create tests for your implementation

```go
// pkg/cardgame/blackjack/blackjack_test.go
package blackjack_test

import (
    "testing"

    "github.com/blixxurd/card-game-go/pkg/cardgame/blackjack"
    "github.com/blixxurd/card-game-go/pkg/cardgame/player"
)

// TestBlackjackGameCreation tests the creation of a blackjack game
func TestBlackjackGameCreation(t *testing.T) {
    // Create a new factory
    factory := blackjack.NewBlackjackGameFactory()

    // Create a game
    game, err := factory.CreateGame(map[string]interface{}{
        "id":      "test",
        "name":    "Test Game",
        "min_bet": 10,
    })
    if err != nil {
        t.Errorf("Failed to create game: %v", err)
    }

    // Check game properties
    if game.ID() != "test" {
        t.Errorf("Expected game ID 'test', got '%s'", game.ID())
    }

    if game.Name() != "Test Game" {
        t.Errorf("Expected game name 'Test Game', got '%s'", game.Name())
    }

    // Add a player
    player1 := player.NewPlayer("p1", "Player 1")
    err = game.AddPlayer(player1)
    if err != nil {
        t.Errorf("Failed to add player: %v", err)
    }

    // Start the game
    err = game.Start()
    if err != nil {
        t.Errorf("Failed to start game: %v", err)
    }

    // Check that the game is in progress
    if game.State() != "in_progress" {
        t.Errorf("Expected game state 'in_progress', got '%s'", game.State())
    }

    // Cast to BlackjackGame to access specific methods
    blackjackGame, ok := game.(blackjack.BlackjackGame)
    if !ok {
        t.Errorf("Failed to cast game to BlackjackGame")
        return
    }

    // Check dealer's up card
    dealerUpCard := blackjackGame.DealerUpCard()
    if dealerUpCard == nil {
        t.Errorf("Expected dealer up card, got nil")
    }

    // Test allowed actions
    actions := blackjackGame.AllowedBlackjackActions("p1")
    if len(actions) == 0 {
        t.Errorf("Expected at least one allowed action")
    }

    // Add more tests...
}
```

## Example Usage

Here's how you would use your blackjack implementation:

```go
// Create a blackjack game
blackjackFactory := blackjack.NewBlackjackGameFactory()
blackjackGame, err := blackjackFactory.CreateGame(map[string]interface{}{
    "id":      "bj1",
    "name":    "My Blackjack Game",
    "min_bet": 10,
})
if err != nil {
    panic(err)
}

// Add players
player1 := player.NewPlayer("p1", "Player 1")
blackjackGame.AddPlayer(player1)

// Start the game
blackjackGame.Start()

// Cast to BlackjackGame to access specific methods
bjGame := blackjackGame.(blackjack.BlackjackGame)

// Check if player has blackjack
if bjGame.IsBlackjack("p1") {
    fmt.Println("Player 1 has blackjack!")
}

// Process a hit action
hitAction := blackjack.NewBlackjackAction(blackjack.ActionHit, "p1")
bjGame.ProcessBlackjackAction(hitAction)

// Check if player busted
if bjGame.IsBusted("p1") {
    fmt.Println("Player 1 busted!")
}

// Process a stand action
standAction := blackjack.NewBlackjackAction(blackjack.ActionStand, "p1")
bjGame.ProcessBlackjackAction(standAction)
```

## Key Considerations for Blackjack

When implementing Blackjack, pay special attention to these aspects:

1. **Hand Evaluation**: Implement proper logic for calculating hand values, considering that aces can be worth 1 or 11.
2. **Dealer Logic**: The dealer must follow specific rules (e.g., hit on 16 or less, stand on 17 or more).
3. **Special Actions**: Implement special actions like doubling down, splitting pairs, and insurance.
4. **Payout Logic**: Implement correct payout logic (e.g., blackjack typically pays 3:2).

## Best Practices

1. **Follow the interfaces**: Ensure your implementation satisfies all methods required by the interfaces.
2. **Maintain separation of concerns**: Keep game logic separate from UI or network code.
3. **Use composition**: Leverage existing components like deck and card implementations.
4. **Write tests**: Create comprehensive tests for your game implementation.
5. **Document your code**: Add clear comments and documentation for your implementation.
6. **Handle errors gracefully**: Implement proper error handling for all operations.
7. **Consider concurrency**: If your game will be used in a concurrent environment, ensure thread safety. 