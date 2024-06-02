package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/defgadget/gopokedex/internal/pokecache"
)

type LocationAreas []LocationArea

type LocationArea struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	GameIndex            int                    `json:"game_index"`
	EncounterMethodRates []EncounterMethodRates `json:"encounter_method_rates"`
	Location             Location               `json:"location"`
	Names                []Names                `json:"names"`
	PokemonEncounters    []PokemonEncounters    `json:"pokemon_encounters"`
}
type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type VersionDetails struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}
type EncounterMethodRates struct {
	EncounterMethod EncounterMethod  `json:"encounter_method"`
	VersionDetails  []VersionDetails `json:"version_details"`
}
type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Names struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}

type Method struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type EncounterDetails struct {
	MinLevel        int    `json:"min_level"`
	MaxLevel        int    `json:"max_level"`
	ConditionValues []any  `json:"condition_values"`
	Chance          int    `json:"chance"`
	Method          Method `json:"method"`
}
type EncounterVersionDetails struct {
	Version          Version            `json:"version"`
	MaxChance        int                `json:"max_chance"`
	EncounterDetails []EncounterDetails `json:"encounter_details"`
}
type PokemonEncounters struct {
	Pokemon        Pokemon                   `json:"pokemon"`
	VersionDetails []EncounterVersionDetails `json:"version_details"`
}

var cache pokecache.Cache = *pokecache.NewCache(time.Second * 10)

func GetLocationAreas(url string, c *Config) (LocationAreas, error) {
	data, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		cache.Add(url, data)
	}

	err := json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c.Results, nil
}

func GetAreaPokemon(location string) ([]Pokemon, error) {
	url := fmt.Sprintf("%s/%s/location-area/%s", api, version, location)
	data, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Get failed on ", url, err)
			return nil, err
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Reading body failed: ", err)
			return nil, err
		}
		cache.Add(url, data)
	}

	area := &LocationArea{}
	err := json.Unmarshal(data, area)
	if err != nil {
		fmt.Println("Json failed: ", err)
		return nil, err
	}

	pokemon := make([]Pokemon, 1)
	for _, encounter := range area.PokemonEncounters {
		pokemon = append(pokemon, encounter.Pokemon)
	}

	return pokemon, nil

}
