# Pokedex CLI

A simple command-line Pokédex application that lets you explore locations, catch Pokémon, and manage your collection using the PokéAPI.

## Features

- **Command History**: Use UP/DOWN arrow keys to navigate through previous commands
- **Tab Completion**: Press TAB to autocomplete commands and common Pokémon/location names
- **Location Exploration**: Browse Pokémon locations and discover what Pokémon can be found there
- **Pokémon Catching**: Attempt to catch Pokémon with realistic success rates based on their strength
- **Pokédex Management**: Keep track of your caught Pokémon and inspect their details
- **Caching System**: Fast responses with automatic data caching

## Installation

```bash
git clone <repository-url>
cd pokedex
go build
./pokedex
```

## Commands

- `help` - Show all available commands
- `map` - Display 20 location areas to explore
- `mapb` - Go back to previous 20 locations
- `explore <location-name>` - See what Pokémon are in a specific location
- `catch <pokemon-name>` - Try to catch a Pokémon
- `inspect <pokemon-name>` - View details of a caught Pokémon
- `pokedx` - List all your caught Pokémon
- `exit` - Quit the application

## Usage Examples

```
Pokedex > map
1 => canalave-city-area
2 => eterna-city-area
...

Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - tentacool
 - tentacruel
 - magikarp
 - gyarados

Pokedx > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!

Pokedx > inspect pikachu
Name: pikachu
Height: 4
Weight: 60
Stats:
  - hp: 35
  - attack: 55
  - defense: 40
  - special-attack: 50
  - special-defense: 50
  - speed: 90
Types:
  - electric

Pokedx > pokedx
Your Pokedx:
  - pikachu
```

## Tips

- Use arrow keys to cycle through command history
- Press TAB for command and name suggestions
- Stronger Pokémon (higher base experience) are harder to catch
- All data is cached for faster subsequent requests
- Commands are case-insensitive

Enjoy building your Pokédex collection!
