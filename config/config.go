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
			Net             string `yaml:"net" toml:"net" env:"instance_net" flag:"instance-net"`
			DataPath        string `yaml:"data_path" toml:"data_path" env:"instance_data_path" flag:"instance-data-path"`
			DefaultHostname string `yaml:"default_hostname" toml:"default_hostname" env:"default_hostname" flag:"default-hostname"`
		}
	}
	PortAllocator struct {
		MinPort int `yaml:"min_port" toml:"min_port" env:"min_port" flag:"min-port"`
		MaxPort int `yaml:"max_port" toml:"max_port" env:"max_port" flag:"max-port"`
	}
}
