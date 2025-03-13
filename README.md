# Card Game Framework in Go

This project implements a flexible and composable card game framework with a primary focus on poker to start. The framework is designed with flexibility in mind, allowing for potential expansion to other card games and online multiplayer functionality.

## Key Features

- Robust deck and hand management
- Comprehensive poker hand evaluation
- Texas Hold'em game simulation with side pot handling
- Extensible architecture for different poker variants
- Flexible betting structures (No Limit, Pot Limit, Fixed Limit)
- Clean separation of concerns with well-defined interfaces
- Adapter pattern for integrating existing implementations with new interfaces

## Project Goals

1. Demonstrate Go's capabilities for implementing turn-based games
2. Explore architecture patterns for card game development
3. Experiment with casino-style games and games of chance in a software engineering context
4. Provide a foundation for potential online multiplayer card games
5. Showcase clean architecture principles and design patterns in Go

## Structure

- `cmd/main.go`: Main application demonstrating the use of the framework (Currently runs a holdem simulation)
- `pkg/cardgame/`: Package containing core card game interfaces and implementations
  - `card/`: Card interfaces and implementations
  - `deck/`: Deck interfaces and implementations
  - `hand/`: Hand interfaces and implementations
  - `player/`: Player interfaces and implementations
  - `game/`: Base game interfaces and types
  - `table/`: Table management for poker games
  - `poker/`: Poker-specific interfaces and factory
    - `pokerhand/`: Poker hand evaluation
    - `holdem/`: Texas Hold'em game implementation and adapter

## Usage

To run this project, follow these steps:

1. Ensure you have Go installed on your system.

2. Clone the repository:
   ```
   git clone https://github.com/yourusername/card-game-go.git
   cd card-game-go
   ```

3. Run the project:
   ```
   go run cmd/main.go
   ```

4. If you encounter any "undefined" errors, try running:
   ```
   go mod tidy
   ```

## Current Functionality

- Card and deck management with shuffling and drawing capabilities
- Hand dealing and verification
- Poker hand evaluation (including games with community cards)
- Texas Hold'em game simulation with proper side pot handling
- Extensibility framework for different poker variants
- Factory pattern for creating different types of poker games
- Well-defined interfaces for extending to other card games

## Architecture

The project follows a clean architecture approach with:

- Clear separation between interfaces and implementations
- Domain-driven design principles
- Composition over inheritance
- Dependency injection for flexible component wiring
- Adapter pattern for integrating existing implementations with new interfaces
- Factory pattern for creating different game variants

### Extensibility Framework

The project includes an extensibility framework that allows for:

1. **Different Poker Variants**: The `PokerVariant` interface defines the contract for implementing different poker variants like Texas Hold'em, Omaha, Seven-Card Stud, etc.

2. **Poker Game Factory**: The `PokerGameFactory` provides a way to register and create different poker variants.

3. **Adapter Pattern**: The `HoldemGameAdapter` adapts the existing `HoldemGame` implementation to the new `PokerGame` interface, allowing for backward compatibility.

4. **Betting Structures**: Support for different betting structures (No Limit, Pot Limit, Fixed Limit) through the `table` package.

5. **Game State Management**: Well-defined game states and transitions for different poker variants.

## Future Directions

### Additional Poker Variants
- **Omaha**: Implement Omaha Hold'em with four hole cards per player and rules for using exactly two hole cards
- **Seven-Card Stud**: Add support for stud poker games with no community cards
- **Five-Card Draw**: Implement classic draw poker with card replacement

### Expanded Game Types
- **Blackjack**: Implement dealer logic, splitting, doubling down, and insurance

### Online Multiplayer System
- **WebSocket Server**: Real-time communication for live gameplay
- **Integration with Player Accounts**: Leveraging the ID from a third party account service
- **Game Rooms**: Lobby system for creating and joining games
- **Spectator Mode**: Allow users to watch ongoing games

### User Interface Options
- **Command-Line Interface**: Text-based gameplay for quick testing
- **Websocket Interface**: Responsive websocket interface to play with before adding a UI
- 

### AI Players
- **Basic Strategy AI**: Rule-based AI for simple decision making
- **Statistical AI**: Decision making based on pot odds and expected value
- **Personality-Based AI**: Different AI styles (tight-aggressive, loose-passive, etc.)
- **Machine Learning Integration**: Train models on gameplay data for advanced AI
- **Difficulty Levels**: Configurable AI strength for appropriate challenge

### Performance Optimizations
- **Benchmarking**: Identify and address performance bottlenecks
- **Concurrency Improvements**: Better utilization of Go's concurrency features
- **Memory Optimization**: Reduce memory footprint for large-scale deployments

### Additional Features
- **Hand History & Logging**: Detailed logging and replay of previous hands