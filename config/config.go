package config

type ValdockConfig struct {
	Server struct {
		Port int    `yaml:"port" toml:"port" env:"PORT" flag:"port"`
		Host string `yaml:"host" toml:"host" env:"HOST" flag:"host"`
	}
	KV struct {
		Url string `yaml:"url" toml:"url" env:"VALKEY_URL" flag:"valkey-url"`
	}
	Docker struct {
		Instance struct {
			Net string `yaml:"net" toml:"net" env:"instance_net" flag:"instance-net"`
		}
	}
}
