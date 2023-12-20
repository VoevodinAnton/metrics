package config

type Logger struct {
	Encoding    string `mapstructure:"encoding"`
	Level       string `mapstructire:"level"`
	Development bool   `mapstructure:"development"`
}

type Postgres struct {
	DatabaseDSN string
}

type Server struct {
	Address string
}
