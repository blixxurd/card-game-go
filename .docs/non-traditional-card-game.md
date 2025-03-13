
# Implementing a non-traditional card game 

To implement a game that doesn't use standard playing cards (e.g., Uno), you'll need to create custom card types and game logic:

1. **Define custom card types**:
   ```go
   // pkg/cardgame/uno/card.go
   package uno

   import (
       "fmt"
       "github.com/blixxurd/card-game-go/pkg/cardgame/card"
   )

   // UnoColor represents the color of an Uno card
   type UnoColor string

   const (
       Red    UnoColor = "red"
       Blue   UnoColor = "blue"
       Green  UnoColor = "green"
       Yellow UnoColor = "yellow"
       Wild   UnoColor = "wild"
   )

   // UnoValue represents the value of an Uno card
   type UnoValue string

   const (
       Zero       UnoValue = "0"
       One        UnoValue = "1"
       Two        UnoValue = "2"
       // ... other numbers ...
       Nine       UnoValue = "9"
       Skip       UnoValue = "skip"
       Reverse    UnoValue = "reverse"
       DrawTwo    UnoValue = "draw_two"
       Wild4      UnoValue = "wild_draw_four"
       WildColor  UnoValue = "wild_color"
   )

   // UnoCard represents a card in the Uno game
   type UnoCard struct {
       color UnoColor
       value UnoValue
   }

   // NewUnoCard creates a new Uno card
   func NewUnoCard(color UnoColor, value UnoValue) *UnoCard {
       return &UnoCard{
           color: color,
           value: value,
       }
   }

   // Color returns the color of the card
   func (c *UnoCard) Color() UnoColor {
       return c.color
   }

   // Value returns the value of the card
   func (c *UnoCard) Value() UnoValue {
       return c.value
   }

   // SetColor sets the color of the card (for wild cards)
   func (c *UnoCard) SetColor(color UnoColor) {
       c.color = color
   }

   // String returns a string representation of the card
   func (c *UnoCard) String() string {
       return fmt.Sprintf("%s %s", c.color, c.value)
   }

   // Implement the card.Card interface
   func (c *UnoCard) ID() string {
       return fmt.Sprintf("%s-%s", c.color, c.value)
   }

   func (c *UnoCard) Name() string {
       return c.String()
   }

   func (c *UnoCard) Rank() int {
       // Assign numeric values based on card type
       switch c.value {
       case Zero:
           return 0
       case One:
           return 1
       // ... other numbers ...
       case Nine:
           return 9
       case Skip, Reverse, DrawTwo:
           return 20
       case WildColor:
           return 50
       case Wild4:
           return 50
       default:
           return 0
       }
   }

   func (c *UnoCard) Suit() string {
       return string(c.color)
   }

   func (c *UnoCard) FaceUp() bool {
       return true
   }

   func (c *UnoCard) SetFaceUp(faceUp bool) {
       // Uno cards are always face up in hand
   }

   func (c *UnoCard) Clone() card.Card {
       return &UnoCard{
           color: c.color,
           value: c.value,
       }
   }
   ```

2. **Create a custom deck implementation**:
   ```go
   // pkg/cardgame/uno/deck.go
   package uno

   import (
       "math/rand"
       "time"

       "github.com/blixxurd/card-game-go/pkg/cardgame/card"
       "github.com/blixxurd/card-game-go/pkg/cardgame/deck"
   )

   // UnoDeck represents a deck of Uno cards
   type UnoDeck struct {
       cards []card.Card
       rng   *rand.Rand
   }

   // NewUnoDeck creates a new Uno deck
   func NewUnoDeck() *UnoDeck {
       d := &UnoDeck{
           cards: make([]card.Card, 0, 108), // Uno has 108 cards
           rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
       }
       d.initialize()
       return d
   }

   // initialize creates all the cards in the Uno deck
   func (d *UnoDeck) initialize() {
       // Add number cards (0-9) in each color
       // Each color has one 0 and two of each 1-9
       colors := []UnoColor{Red, Blue, Green, Yellow}
       
       for _, color := range colors {
           // Add one zero card
           d.cards = append(d.cards, NewUnoCard(color, Zero))
           
           // Add two of each 1-9 card
           for _, value := range []UnoValue{One, Two, /* ... */, Nine} {
               d.cards = append(d.cards, NewUnoCard(color, value))
               d.cards = append(d.cards, NewUnoCard(color, value))
           }
           
           // Add two of each action card (Skip, Reverse, Draw Two)
           for _, value := range []UnoValue{Skip, Reverse, DrawTwo} {
               d.cards = append(d.cards, NewUnoCard(color, value))
               d.cards = append(d.cards, NewUnoCard(color, value))
           }
       }
       
       // Add wild cards (4 of each type)
       for i := 0; i < 4; i++ {
           d.cards = append(d.cards, NewUnoCard(Wild, WildColor))
           d.cards = append(d.cards, NewUnoCard(Wild, Wild4))
       }
   }

   // Implement the deck.Deck interface
   func (d *UnoDeck) Cards() []card.Card {
       return d.cards
   }

   func (d *UnoDeck) Shuffle() {
       // Fisher-Yates shuffle
       for i := len(d.cards) - 1; i > 0; i-- {
           j := d.rng.Intn(i + 1)
           d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
       }
   }

   func (d *UnoDeck) Draw() (card.Card, error) {
       if len(d.cards) == 0 {
           return nil, deck.ErrEmptyDeck
       }
       
       card := d.cards[len(d.cards)-1]
       d.cards = d.cards[:len(d.cards)-1]
       return card, nil
   }

   func (d *UnoDeck) DrawMany(n int) ([]card.Card, error) {
       if len(d.cards) < n {
           return nil, deck.ErrNotEnoughCards
       }
       
       cards := d.cards[len(d.cards)-n:]
       d.cards = d.cards[:len(d.cards)-n]
       return cards, nil
   }

   func (d *UnoDeck) AddCard(c card.Card) {
       d.cards = append(d.cards, c)
   }

   func (d *UnoDeck) AddCards(cards []card.Card) {
       d.cards = append(d.cards, cards...)
   }

   func (d *UnoDeck) Count() int {
       return len(d.cards)
   }

   func (d *UnoDeck) Clone() deck.Deck {
       clone := &UnoDeck{
           cards: make([]card.Card, len(d.cards)),
           rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
       }
       
       for i, c := range d.cards {
           clone.cards[i] = c.Clone()
       }
       
       return clone
   }
   ```

3. **Define game-specific actions and states**:
   ```go
   // pkg/cardgame/uno/uno_game.go
   package uno

   import (
       "errors"
       "github.com/blixxurd/card-game-go/pkg/cardgame/game"
       "github.com/blixxurd/card-game-go/pkg/cardgame/player"
   )

   // UnoGameState represents the state of an Uno game
   type UnoGameState string

   const (
       StateInitial  UnoGameState = "initial"
       StatePlaying  UnoGameState = "playing"
       StateComplete UnoGameState = "complete"
   )

   // UnoActionType represents the type of action in Uno
   type UnoActionType string

   const (
       ActionPlayCard UnoActionType = "play_card"
       ActionDrawCard UnoActionType = "draw_card"
       ActionSayUno   UnoActionType = "say_uno"
       ActionChallenge UnoActionType = "challenge"
   )

   // UnoAction represents an action in an Uno game
   type UnoAction struct {
       actionType UnoActionType
       playerID   string
       cardID     string
       color      UnoColor // For wild cards
   }

   // NewUnoAction creates a new Uno action
   func NewUnoAction(actionType UnoActionType, playerID string, cardID string, color UnoColor) *UnoAction {
       return &UnoAction{
           actionType: actionType,
           playerID:   playerID,
           cardID:     cardID,
           color:      color,
       }
   }

   // Type returns the type of the action
   func (a *UnoAction) Type() string {
       return string(a.actionType)
   }

   // PlayerID returns the ID of the player performing the action
   func (a *UnoAction) PlayerID() string {
       return a.playerID
   }

   // Params returns the parameters for the action
   func (a *UnoAction) Params() map[string]interface{} {
       return map[string]interface{}{
           "card_id": a.cardID,
           "color":   string(a.color),
       }
   }
   ```

4. **Implement the game logic**:
   ```go
   // UnoGame represents a game of Uno
   type UnoGame struct {
       id            string
       name          string
       players       []player.Player
       playerHands   map[string][]card.Card
       drawPile      *UnoDeck
       discardPile   []card.Card
       currentPlayer int
       direction     int // 1 for clockwise, -1 for counter-clockwise
       state         UnoGameState
       topCard       *UnoCard
       winners       []player.Player
       customData    map[string]interface{}
   }

   // NewUnoGame creates a new Uno game
   func NewUnoGame(id, name string) *UnoGame {
       return &UnoGame{
           id:            id,
           name:          name,
           players:       make([]player.Player, 0),
           playerHands:   make(map[string][]card.Card),
           drawPile:      NewUnoDeck(),
           discardPile:   make([]card.Card, 0),
           currentPlayer: 0,
           direction:     1, // Start clockwise
           state:         StateInitial,
           customData:    make(map[string]interface{}),
       }
   }

   // Implement the game.Game interface and Uno-specific methods
   // ...

   // Start begins the game
   func (g *UnoGame) Start() error {
       if len(g.players) < 2 {
           return errors.New("not enough players")
       }

       g.drawPile.Shuffle()

       // Deal 7 cards to each player
       for _, p := range g.players {
           cards, err := g.drawPile.DrawMany(7)
           if err != nil {
               return err
           }
           g.playerHands[p.ID()] = cards
       }

       // Draw the first card for the discard pile
       // (Ensure it's not a wild card or action card for the first turn)
       var firstCard card.Card
       for {
           var err error
           firstCard, err = g.drawPile.Draw()
           if err != nil {
               return err
           }

           unoCard, ok := firstCard.(*UnoCard)
           if !ok {
               return errors.New("invalid card type")
           }

           // If it's a regular number card, use it
           if unoCard.color != Wild && unoCard.value != Skip && 
              unoCard.value != Reverse && unoCard.value != DrawTwo {
               break
           }

           // Otherwise, put it back in the deck and try again
           g.drawPile.AddCard(firstCard)
           g.drawPile.Shuffle()
       }

       g.discardPile = append(g.discardPile, firstCard)
       g.topCard = firstCard.(*UnoCard)
       g.state = StatePlaying

       return nil
   }

   // ProcessAction processes a player action
   func (g *UnoGame) ProcessAction(action game.Action) error {
       if g.state != StatePlaying {
           return errors.New("game is not in playing state")
       }

       if action.PlayerID() != g.players[g.currentPlayer].ID() {
           return errors.New("not your turn")
       }

       unoAction, ok := action.(*UnoAction)
       if !ok {
           return errors.New("invalid action type")
       }

       switch unoAction.actionType {
       case ActionPlayCard:
           return g.playCard(unoAction)
       case ActionDrawCard:
           return g.drawCard(unoAction)
       case ActionSayUno:
           return g.sayUno(unoAction)
       case ActionChallenge:
           return g.challenge(unoAction)
       default:
           return errors.New("invalid action type")
       }
   }

   // Additional methods for game logic
   // ...
   ```

5. **Create a factory for your game**:
   ```go
   // UnoGameFactory creates Uno games
   type UnoGameFactory struct{}

   // NewUnoGameFactory creates a new UnoGameFactory
   func NewUnoGameFactory() *UnoGameFactory {
       return &UnoGameFactory{}
   }

   // CreateGame creates a new Uno game
   func (f *UnoGameFactory) CreateGame(params map[string]interface{}) (game.Game, error) {
       id, ok := params["id"].(string)
       if !ok {
           return nil, errors.New("invalid game parameters: missing id")
       }

       name, ok := params["name"].(string)
       if !ok {
           return nil, errors.New("invalid game parameters: missing name")
       }

       return NewUnoGame(id, name), nil
   }

   // Type returns the type of games this factory creates
   func (f *UnoGameFactory) Type() game.GameType {
       return "uno"
   }
   ```

6. **Register your game factory**:
   ```go
   // In your application initialization
   gameRegistry := NewGameRegistry()
   gameRegistry.RegisterFactory("uno", NewUnoGameFactory())
   ```

7. **Example usage**:
   ```go
   // Create an Uno game
   unoFactory := uno.NewUnoGameFactory()
   unoGame, err := unoFactory.CreateGame(map[string]interface{}{
       "id":   "uno1",
       "name": "My Uno Game",
   })
   if err != nil {
       panic(err)
   }

   // Add players
   player1 := player.NewPlayer("p1", "Player 1")
   player2 := player.NewPlayer("p2", "Player 2")
   unoGame.AddPlayer(player1)
   unoGame.AddPlayer(player2)

   // Start the game
   unoGame.Start()

   // Play a card
   playAction := uno.NewUnoAction(uno.ActionPlayCard, "p1", "red-5", "")
   unoGame.ProcessAction(playAction)

   // Draw a card
   drawAction := uno.NewUnoAction(uno.ActionDrawCard, "p2", "", "")
   unoGame.ProcessAction(drawAction)
   ```

### Adapting the Framework

When implementing a non-traditional card game like Uno, you may need to adapt or extend certain parts of the framework:

1. **Custom Card Types**: Create your own card types that implement the `card.Card` interface.
2. **Custom Deck Implementation**: Implement a deck specific to your game's needs.
3. **Game-Specific Rules**: Implement the unique rules and mechanics of your game.
4. **Custom Actions**: Define actions specific to your game.
5. **Hand Evaluation**: Create custom logic for determining winning conditions.

The framework's modular design allows you to reuse components like player management and game state while customizing the card-specific aspects to fit your game's requirements.

### Example: Using Your Non-Traditional Card Game

Here's a complete example of how to use your Uno implementation in an application:

```go
package main

import (
    "fmt"
    "github.com/blixxurd/card-game-go/pkg/cardgame/player"
    "github.com/blixxurd/card-game-go/pkg/cardgame/uno"
)

func main() {
    // Create an Uno game
    unoFactory := uno.NewUnoGameFactory()
    unoGame, err := unoFactory.CreateGame(map[string]interface{}{
        "id":   "uno1",
        "name": "My Uno Game",
    })
    if err != nil {
        panic(err)
    }

    // Add players
    player1 := player.NewPlayer("p1", "Player 1")
    player2 := player.NewPlayer("p2", "Player 2")
    player3 := player.NewPlayer("p3", "Player 3")
    player4 := player.NewPlayer("p4", "Player 4")
    
    unoGame.AddPlayer(player1)
    unoGame.AddPlayer(player2)
    unoGame.AddPlayer(player3)
    unoGame.AddPlayer(player4)

    // Start the game
    err = unoGame.Start()
    if err != nil {
        panic(err)
    }

    fmt.Println("Game started!")
    
    // Game loop (simplified)
    for unoGame.State() == "playing" {
        currentPlayer := unoGame.CurrentPlayer()
        fmt.Printf("It's %s's turn\n", currentPlayer.Name())
        
        // Get the player's hand
        hand := unoGame.(*uno.UnoGame).PlayerHand(currentPlayer.ID())
        fmt.Printf("Your hand: %v\n", hand)
        
        // Get the top card
        topCard := unoGame.(*uno.UnoGame).TopCard()
        fmt.Printf("Top card: %s\n", topCard)
        
        // For this example, we'll just have the player draw a card
        drawAction := uno.NewUnoAction(uno.ActionDrawCard, currentPlayer.ID(), "", "")
        err = unoGame.ProcessAction(drawAction)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
        }
        
        // Check if the player can play a card from their hand
        hand = unoGame.(*uno.UnoGame).PlayerHand(currentPlayer.ID())
        for _, card := range hand {
            unoCard := card.(*uno.UnoCard)
            if unoCard.Color() == topCard.Color() || unoCard.Value() == topCard.Value() || unoCard.Color() == uno.Wild {
                // Play the card
                playAction := uno.NewUnoAction(uno.ActionPlayCard, currentPlayer.ID(), card.ID(), "")
                if unoCard.Color() == uno.Wild {
                    // Choose a color for wild cards
                    playAction = uno.NewUnoAction(uno.ActionPlayCard, currentPlayer.ID(), card.ID(), uno.Red)
                }
                
                err = unoGame.ProcessAction(playAction)
                if err != nil {
                    fmt.Printf("Error playing card: %v\n", err)
                } else {
                    fmt.Printf("%s played %s\n", currentPlayer.Name(), card)
                    break
                }
            }
        }
        
        // Check if the player has one card left
        if len(unoGame.(*uno.UnoGame).PlayerHand(currentPlayer.ID())) == 1 {
            // Say "UNO!"
            unoAction := uno.NewUnoAction(uno.ActionSayUno, currentPlayer.ID(), "", "")
            unoGame.ProcessAction(unoAction)
            fmt.Printf("%s says UNO!\n", currentPlayer.Name())
        }
        
        // Check if the player has won
        if len(unoGame.(*uno.UnoGame).PlayerHand(currentPlayer.ID())) == 0 {
            fmt.Printf("%s has won the game!\n", currentPlayer.Name())
            break
        }
        
        // For this example, we'll just break after one turn
        break
    }
}
```

## Key Considerations for Uno

When implementing Uno, pay special attention to these aspects:

1. **Card Matching Rules**: Implement proper logic for matching cards by color or value.
2. **Special Card Effects**: Handle the effects of special cards like Skip, Reverse, Draw Two, and Wild cards.
3. **Direction Management**: Implement the ability to change the direction of play with Reverse cards.
4. **UNO Declaration**: Implement the rule that players must declare "UNO" when they have one card left.
5. **Scoring System**: Implement the scoring system based on the cards left in players' hands when someone goes out.

## Best Practices

1. **Follow the interfaces**: Ensure your implementation satisfies all methods required by the interfaces.
2. **Maintain separation of concerns**: Keep game logic separate from UI or network code.
3. **Use composition**: Leverage existing components like player management and game state.
4. **Write tests**: Create comprehensive tests for your game implementation.
5. **Document your code**: Add clear comments and documentation for your implementation.
6. **Handle errors gracefully**: Implement proper error handling for all operations.
7. **Consider concurrency**: If your game will be used in a concurrent environment, ensure thread safety.

By following these guidelines, you can implement virtually any card game within the framework, even those that use non-traditional cards or have unique mechanics.