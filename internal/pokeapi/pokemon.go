package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Pokemon struct {
	ID                     int         `json:"id"`
	Name                   string      `json:"name"`
	BaseExperience         int         `json:"base_experience"`
	Height                 int         `json:"height"`
	IsDefault              bool        `json:"is_default"`
	Order                  int         `json:"order"`
	Weight                 int         `json:"weight"`
	Abilities              []Abilities `json:"abilities"`
	Forms                  []Forms     `json:"forms"`
	LocationAreaEncounters string      `json:"location_area_encounters"`
	URL                    string      `json:"url"`
	Stats                  []Stats     `json:"stats"`
	Types                  []Types     `json:"types"`
}
type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Abilities struct {
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
	Ability  Ability `json:"ability"`
}
type Forms struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

func GetPokemonDetails(pokemon string) (Pokemon, error) {
	url := fmt.Sprintf("%s/%s/pokemon/%s", api, version, pokemon)
	data, ok := cache.Get(url)
	empty := Pokemon{}
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Get failed on ", url, err)
			return empty, err
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Reading body failed: ", err)
			return empty, err
		}
		cache.Add(url, data)
	}

	pokemonDetails := Pokemon{}
	err := json.Unmarshal(data, &pokemonDetails)
	if err != nil {
		fmt.Println("Json failed: ", err)
		return empty, err
	}
	return pokemonDetails, nil
}
