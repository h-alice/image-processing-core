package config

import (
	"errors"
	op "imagetools/operation" // Grab `EncoderOption` from operation package.
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var ErrNotImplemented = errors.New("operation not implemented")

type OutputOptionConfig op.EncoderOption

type OutputConfig struct {
	Format     string              `yaml:"format"`     // Output file format
	NameSuffix string              `yaml:"nameSuffix"` // Output file name suffix
	NamePrefix string              `yaml:"namePrefix"` // Output file name prefix
	Options    *OutputOptionConfig `yaml:"options"`    // Encoder option
}

// Generate output file name.
func (ocf OutputConfig) GenerateFileName(input_name string) string {

	original_ext := filepath.Ext(input_name)                     // Get file extension.
	original_name := filepath.Base(input_name)                   // Get file name.
	stem := original_name[:len(original_name)-len(original_ext)] // Get file name w/o extension.

	fileSuffix := ""

	switch strings.ToLower(ocf.Format) {
	case "jpeg":
		fileSuffix = ".jpg" // Use JPG instead of JPEG.
	case "":
		fileSuffix = original_ext // Output format not specified: keep original extension.
	default:
		fileSuffix = "." + ocf.Format // Use specified output format.

	}
	full_file := ocf.NamePrefix + stem + ocf.NameSuffix + fileSuffix

	return full_file
}

type ResizeConfig struct {
	Width     int     `yaml:"width"`     // Output image width
	Height    int     `yaml:"height"`    // Output image height
	Factor    float32 `yaml:"factor"`    // Resize factor
	Algorithm string  `yaml:"algorithm"` // Resize algorithm
}

type ProcessProfileConfig struct {
	ProfileName string        `yaml:"profileName"` // Profile identifier
	ICC         string        `yaml:"icc"`         // ICC profile to embed
	Resize      *ResizeConfig `yaml:"resize"`      // Resize option
	Output      *OutputConfig `yaml:"output"`      // Output file configuration
}

type OutputDirConfig struct {
	DirName *string `yaml:"dirName"` // Output directory name
}

type ConfigFile struct {
	Profiles []ProcessProfileConfig `yaml:"profiles"` // List of profile configurations
}

// Load config file from path.
func LoadConfig(config_path string) (*ConfigFile, error) {

	raw_config, err := os.ReadFile(config_path) // Read raw config file.
	if err != nil {
		return nil, err
	}

	// Converting JSON to config structure.
	var conf ConfigFile                     // Parsed config placeholder.
	err = yaml.Unmarshal(raw_config, &conf) // Convert JSON to structure.
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
