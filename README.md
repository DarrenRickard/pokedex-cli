# pokedex-cli

A simple command-line Pokédex built in Go.  
Search, explore, and catch Pokémon directly from your terminal using the [PokeAPI](https://pokeapi.co/).

## Features
- Fetch Pokémon details by name.
- List available Pokémon in the region.
- Simulate catching Pokémon based on base experience.
- View caught Pokémon in your Pokédex.

## Installation
```bash
git clone https://github.com/darrenrickard/pokedex-cli.git
```
## Go Installation 
### Webi install
- macOS/Linux 
```bash
curl -sS https://webi.sh/golang | sh; \
source ~/.config/envman/PATH.env
```
- Windows
```bash
curl.exe https://webi.ms/golang | powershell
```
### Official install
https://go.dev/doc/install

## Dependencies
- Go 1.20+

- Internet connection (for API calls)

## Usage
```bash
cd path/to/pokedex-cli
go run .
```

## Available commands
```bash
help                    # Show all commands
map                     # Show next set of world locations
mapb                    # Show previous set of world locations
explore <location>      # View existing Pokemon within a world location
catch <pokemon>         # Attempt to catch Pokemon
inspect <pokemon>       # View details of a Pokemon you've caught
pokedex                 # List all caught Pokemon
ccache                  # View cached locations
```

