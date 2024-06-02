package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/defgadget/gopokedex/internal/pokeapi"
)

type Command struct {
	name        string
	description string
	callback    func(*pokeapi.Config, string) error
}

func commands() map[string]Command {
	return map[string]Command{
		"map": {
			name:        "map",
			description: "Travel forward through the world",
			callback:    mapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Travel back to where you previously came from",
			callback:    mapb,
		},
		"explore": {
			name:        "explore",
			description: "Discover pokemon in location -- explore <location>",
			callback:    explore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch pokemon -- catch <pokemon>",
			callback:    catch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon you have caught -- inspect <pokemon>",
			callback:    inspect,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func commandHelp(c *pokeapi.Config, arg string) error {
	_ = c
	_ = arg
	commands := commands()
	fmt.Println()
	fmt.Println("Welcome to the Pokedex")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()

	return nil
}

func commandExit(c *pokeapi.Config, arg string) error {
	_ = c
	_ = arg
	os.Exit(0)
	return nil
}

func explore(c *pokeapi.Config, arg string) error {
	location := arg
	pokemon, err := pokeapi.GetAreaPokemon(location)
	if err != nil {
		fmt.Println("failed to explore: ", err)
	}
	fmt.Println("Exploring ", location, "...")
	printPokemon(pokemon)
	return err
}

var pokedex = make(map[string]pokeapi.Pokemon)

func catch(c *pokeapi.Config, pokemon string) error {
	_ = c
	pokemonDetails, err := pokeapi.GetPokemonDetails(pokemon)
	if err != nil {
		fmt.Println("catch errored: ", err)
		return err
	}

	difficulty := 0
	// if pokemonDetails.BaseExperience <= 20 {
	// 	difficulty = 25
	// }
	// if pokemonDetails.BaseExperience >= 50 {
	// 	difficulty = 75
	// }

	fmt.Printf("Throwing a Pokeball at %s\n", pokemon)
	if rand.Intn(100) > difficulty {
		fmt.Printf("You just caught %s!!\n", pokemon)
		pokedex[pokemon] = pokemonDetails
		return nil
	}

	fmt.Printf("Sorry... %s got away. Better luck next time!\n", pokemon)
	return nil
}

func inspect(c *pokeapi.Config, pokemon string) error {
	_ = c
	details, ok := pokedex[pokemon]
	if !ok {
		fmt.Printf("You haven't caught a %s\n", pokemon)
		return nil
	}
	printDetails(details)
	return nil
}

func printDetails(details pokeapi.Pokemon) {
	name := details.Name
	height := details.Height
	weight := details.Weight
	stats := details.Stats
	types := details.Types

	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Weight: %d\n", weight)
	fmt.Println("Stats:")
	for _, stat := range stats {
		fmt.Printf("\t-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range types {
		fmt.Printf("\t-%s\n", t.Type.Name)
	}
}

func mapf(c *pokeapi.Config, arg string) error {
	locations, err := pokeapi.GetLocationAreas(c.Next, c)
	printLocations(locations)
	return err
}

func mapb(c *pokeapi.Config, arg string) error {
	if c.Previous == "" || strings.Contains(c.Previous, "offset=0") {
		return errors.New("Already on first page")
	}

	locations, err := pokeapi.GetLocationAreas(c.Previous, c)
	printLocations(locations)
	return err
}

func printLocations(locations pokeapi.LocationAreas) {
	for _, location := range locations {
		fmt.Println(location.Name)
	}
}

func printPokemon(pokemon []pokeapi.Pokemon) {
	fmt.Println("Found Pokemon:")
	for _, p := range pokemon {
		fmt.Println(" - ", p.Name)
	}
}

func Run() {
	commands := commands()
	config := pokeapi.NewDefaultConfig()
	arg := ""
	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanWords)
		scanner.Scan()
		input := scanner.Text()
		if command, ok := commands[input]; ok {
			switch input {
			case "explore":
				scanner.Scan()
				arg = scanner.Text()
			case "catch":
				scanner.Scan()
				arg = scanner.Text()
			case "inspect":
				scanner.Scan()
				arg = scanner.Text()
			}
			if err := command.callback(config, arg); err != nil {
				fmt.Println("command failed: ", err)
			}
		} else {
			fmt.Println("No commands match: ", input)
			commandHelp(config, arg)
		}
	}
}
