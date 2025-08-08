package types

type Config struct {
	Port        string `env:"PORT" json:"port,omitempty"`
	DatabaseUrl string `env:"DATABASE_URL" json:"database_url,omitempty"`
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
