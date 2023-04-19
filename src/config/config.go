package config

type Config struct {
	Tmdb        tmdbConfig        `mapstructure:"tmdb"`
	Db          dbConfig          `mapstructure:"db"`
	Meilisearch meilisearchConfig `mapstructure:"meilisearch"`
}

type tmdbConfig struct {
	ApiKey string `mapstructure:"api_key"`
}

type dbConfig struct {
	Driver string `mapstructure:"driver"`
	Uri    string `mapstructure:"uri"`
}

type meilisearchConfig struct {
	Host   string `mapstructure:"host"`
	ApiKey string `mapstructure:"api_key"`
}
