package config

type Config struct {
	Tmdb        TmdbConfig        `mapstructure:"tmdb"`
	Db          DbConfig          `mapstructure:"db"`
	Meilisearch MeilisearchConfig `mapstructure:"meilisearch"`
	Redis       RedisConfig
}

type TmdbConfig struct {
	ApiKey string `mapstructure:"api_key"`
}

type DbConfig struct {
	Driver   string `mapstructure:"driver"`
	Uri      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type MeilisearchConfig struct {
	Host   string `mapstructure:"host"`
	ApiKey string `mapstructure:"api_key"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}
