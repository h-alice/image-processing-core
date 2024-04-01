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

func TestLoadConfigAndParse(t *testing.T) {

	config, err := LoadConfig("test_resources/test1.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Profiles) == 0 {
		t.Fatalf("Expected at least one profile, got none.")
	}

	// Test first profile.

	if config.Profiles[0].ProfileName != "Profile1" {
		t.Fatalf("Expected profile name to be 'Profile1', got '%s'", config.Profiles[0].ProfileName)
	}

	if config.Profiles[0].ICC != "ADOBE RGB" {
		t.Fatalf("Expected ICC to be 'ADOBE RGB', got '%s'", config.Profiles[0].ICC)
	}

	// Config 1: Resize config.
	if config.Profiles[0].Resize.Width != 100 {
		t.Fatalf("Expected width to be 100, got '%d'", config.Profiles[0].Resize.Width)
	}

	if config.Profiles[0].Resize.Height != 200 {
		t.Fatalf("Expected height to be 200, got '%d'", config.Profiles[0].Resize.Height)
	}

	if config.Profiles[0].Resize.Factor != 0.9 {
		t.Fatalf("Expected factor to be 0.9, got '%f'", config.Profiles[0].Resize.Factor)
	}

	if config.Profiles[0].Resize.Algorithm != "catmullrom" {
		t.Fatalf("Expected algorithm to be 'catmullrom', got '%s'", config.Profiles[0].Resize.Algorithm)
	}

	// Config 1: Output config.
	if config.Profiles[0].Output.Format != "jpeg" {
		t.Fatalf("Expected format to be 'jpeg', got '%s'", config.Profiles[0].Output.Format)
	}

	if config.Profiles[0].Output.NamePrefix != "prefix1_" {
		t.Fatalf("Expected prefix to be 'prefix1_', got '%s'", config.Profiles[0].Output.NamePrefix)
	}

	if config.Profiles[0].Output.NameSuffix != "_suffix1" {
		t.Fatalf("Expected suffix to be '_suffix1', got '%s'", config.Profiles[0].Output.NameSuffix)
	}

	if config.Profiles[0].Output.Options.Quality != 80 {
		t.Fatalf("Expected quality to be 80, got '%d'", config.Profiles[0].Output.Options.Quality)
	}

}
