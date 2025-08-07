package types

import "sync"

type Config struct {
	Port        string `env:"PORT"`
	DatabaseUrl string `env:"DATABASE_URL"`
}

type Flag struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Rules   []Rule `json:"rules"`
}

type Rule struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type FlagStore struct {
	Flags map[int]Flag
	Mutex sync.Mutex
}
