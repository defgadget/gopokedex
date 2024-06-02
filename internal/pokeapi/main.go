package pokeapi

import "fmt"

const (
	api     string = "https://pokeapi.co/api"
	version string = "v2"
)

type Config struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  LocationAreas `json:"results"`
}

func NewDefaultConfig() *Config {
	return &Config{Next: fmt.Sprintf("%s/%s/location-area", api, version)}
}
