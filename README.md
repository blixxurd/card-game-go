# Card Game Framework in Go

This project implements a flexible and composable card game framework with a primary focus on poker to start. The framework is designed with flexibility in mind, allowing for potential expansion to other card games and online multiplayer functionality.

## Key Features

- Robust deck and hand management
- Comprehensive poker hand evaluation
- Texas Hold'em game simulation
- Flexible architecture for easy extension to other card games
- Clean separation of concerns with well-defined interfaces

## Project Goals

1. Demonstrate Go's capabilities for implementing turn-based games
2. Explore architecture patterns for card game development
3. Experiment with casino-style games and games of chance in a software engineering context
4. Provide a foundation for potential online multiplayer card games

## Structure

- `cmd/main.go`: Main application demonstrating the use of the framework (Currently runs a holdem simulation)
- `pkg/cardgame/`: Package containing core card game interfaces and implementations
  - `card/`: Card interfaces and implementations
  - `deck/`: Deck interfaces and implementations
  - `hand/`: Hand interfaces and implementations
  - `player/`: Player interfaces and implementations
  - `poker/`: Poker-specific implementations
    - `pokerhand/`: Poker hand evaluation
    - `holdem/`: Texas Hold'em game implementation

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
- Texas Hold'em game simulation
- Well-defined interfaces for extending to other card games

## Architecture

The project follows a clean architecture approach with:

- Clear separation between interfaces and implementations
- Domain-driven design principles
- Composition over inheritance
- Dependency injection for flexible component wiring

## Future Directions

- Implement additional poker variants (e.g., Omaha, Seven-Card Stud)
- Expand to other card games (e.g., Blackjack, Bridge)
- Develop a full-fledged online multiplayer system
- Create a CLI or GUI for interactive gameplay
- Add comprehensive unit tests for all components