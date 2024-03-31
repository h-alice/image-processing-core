package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {

	config, err := LoadConfig("test_resources/test1.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Logf("Config: %#v", config)
}
