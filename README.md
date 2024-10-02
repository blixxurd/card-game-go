# Card Game Framework in Go

This project implements a flexible and composable card game framework with a primary focus on poker to start. The framework is designed with flexibility in mind, allowing for potential expansion to other card games and online multiplayer functionality.

## Key Features

- Robust deck and hand management
- Comprehensive poker hand evaluation
- Texas Hold'em game simulation
- Basics of a WebSocket module intentded for future online play
- Flexible architecture for easy extension to other card games

## Project Goals

1. Demonstrate Go's capabilities for implementing turn-based games
2. Explore architecture patterns for card game development
3. Experiment with casino-style games and games of chance in a software engineering context
4. Provide a foundation for potential online multiplayer card games

## Structure

- `cmd/main.go`: Main application demonstrating the use of the framework (Currently runs a holdem simulation)
- `internal/cardgame/`: Package containing core card game logic (cards, decks, game management)
- `internal/pokerhand/`: Package for poker hand evaluation
- `internal/games/`: Package for specific game implementations (currently Texas Hold'em)
- `internal/net/`: Package for networking capabilities (WebSocket implementation)

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
- WebSocket infrastructure for potential online play

## Future Directions

- Implement additional poker variants (e.g., Omaha, Seven-Card Stud)
- Expand to other card games (e.g., Blackjack, Bridge)
- Develop a full-fledged online multiplayer system
- Create a CLI or GUI for interactive gameplay

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).