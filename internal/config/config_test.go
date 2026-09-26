package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ManoloEsS/load_balancer/internal/backend"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	testDir := t.TempDir()

	var configtests = []struct {
		name           string
		filename       string
		fileContents   []byte
		expectedConfig *Config
		expectError    bool
	}{
		{
			name:         "Load config with no error",
			filename:     "load_balancer.yaml",
			fileContents: []byte("---\nservers:\n - address: \"localhost:8081\"\nalgorithm: \"round_robin\"\n..."),
			expectedConfig: &Config{
				Servers: []backend.Server{
					{
						Address: "localhost:8081",
					},
				},
				Algorithm: "round_robin",
			},
			expectError: false,
		},
		{
			name:           "Error for inexistent file",
			filename:       "wrong_filename.yaml",
			fileContents:   []byte("---\nservers:\n - address: \"localhost:8081\"\nalgorithm: \"round_robin\"\n..."),
			expectedConfig: nil,
			expectError:    true,
		},
	}

	for _, tt := range configtests {
		t.Run(tt.name, func(t *testing.T) {
			defaultFile := "load_balancer.yaml"
			defaultConfigPath := filepath.Join(testDir, defaultFile)

			pathToTestConfig := filepath.Join(testDir, tt.filename)
			_ = os.WriteFile(pathToTestConfig, tt.fileContents, 0644)

			testCfg, err := LoadConfig(defaultConfigPath)
			if tt.expectError {
				assert.Error(t, err, "expected error")
			}
			assert.Equal(t, tt.expectedConfig, testCfg, "test config does not match expected")
			os.Remove(pathToTestConfig)
		})

	}

}
