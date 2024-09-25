# Primitive Card Game Framework in Go

This project implements a flexible card game framework with a focus on:
- Deck management
- Hand Management
- Flexible ruleset implementation

We test this primitive class by implementing a simulation around a Texas Hold Em' poker ruleset.

This project is primarily an experiment to
1. Test the Go programming Language's viability for turn based gaming
2. Expand my knowledge of architecture around turn based card games
3. Experiment with details around casino games / games of chance in the context of software engineering

## Structure

- \`main.go\`: Main application demonstrating the use of the framework
- \`cardgame/\`: Package containing core card game logic
- \`pokerhand/\`: Package for poker hand evaluation

## Usage

To run this project, follow these steps:

1. If you haven't initialized the Go module yet, do so with:
   \`\`\`
   go mod init
   \`\`\`

4. Run the project:
   \`\`\`
   go run main.go
   \`\`\`

5. If you encounter any "undefined" errors, try running:
   \`\`\`
   go mod tidy
   \`\`\`


## Features

- Card and deck management
- Hand dealing and verification
- Poker hand evaluation (including games with community cards)