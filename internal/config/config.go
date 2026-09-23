package config

import (
	"fmt"
	"os"

	loadserver "github.com/ManoloEsS/load_balancer/internal/load_server"
	"go.yaml.in/yaml/v3"
)

const Config_file string = "load_balancer.yaml"

type Config struct {
	Servers   []loadserver.Server `yaml:"servers"`
	Algorithm string              `yaml:"algorithm"`
}

func NewConfig(filepath string) (*Config, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file contents: %w", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(contents, config)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config contents: %w", err)
	}

	return config, nil
}
