package config

type Config struct {
	Tmdb tmdbConfig `mapstructure:"tmdb"`
	Db   dbConfig   `mapstructure:"db"`
}

type tmdbConfig struct {
	ApiKey string `mapstructure:"api_key"`
}

type dbConfig struct {
	Driver string `mapstructure:"driver"`
	Uri    string `mapstructure:"uri"`
}
