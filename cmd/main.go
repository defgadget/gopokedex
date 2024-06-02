package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/defgadget/gopokedex/internal/pokeapi"
)

func mapf(c *pokeapi.Config) error {
	locations, err := pokeapi.GetLocationAreas(c.Next, c)
	printLocations(locations)
	return err
}

func mapb(c *pokeapi.Config) error {
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

type Command struct {
	name        string
	description string
	callback    func(*pokeapi.Config) error
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

func commandHelp(c *pokeapi.Config) error {
	_ = c
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

func commandExit(c *pokeapi.Config) error {
	_ = c
	os.Exit(0)
	return nil
}

func Run() {
	commands := commands()
	config := pokeapi.NewDefaultConfig()
	for {
		fmt.Print("pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		if command, ok := commands[input]; ok {
			if err := command.callback(config); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("No commands match: ", input)
			commandHelp(config)
		}
	}
}
