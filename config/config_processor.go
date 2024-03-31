package config

import (
	"encoding/json"
	"errors"
	op "imagetools/operation"
	"os"
	"path/filepath"
	"strings"
)

var ErrNotImplemented = errors.New("operation not implemented")

type OutputOptionConfig op.EncoderOption

type OutputConfig struct {
	Format     string              // Output file format
	NameSuffix string              // Output file name suffix
	NamePrefix string              // Output file name prefix
	Options    *OutputOptionConfig // Encoder option
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
	Width     int     // Output image width
	Height    int     // Output image height
	Factor    float32 // Resize factor
	Algorithm string  // Interpolation algorithm
}

type ProfileConfig struct {
	ProfileName string        // Profile identifier
	ICC         string        // ICC profile to embed
	Resize      *ResizeConfig // Resize option
	Output      *OutputConfig // Output file configuraion
}

type OutputDirConfig struct {
	DirName *string
}

type ConfigFile struct {
	// Output   *OutputDirConfig
	Profiles []ProfileConfig
}

// Load config file from path.
func ConfigLoader(config_path string) (*ConfigFile, error) {

	raw_config, err := os.ReadFile(config_path) // Read raw config file.
	if err != nil {
		return nil, err
	}

	// Converting JSON to config structure.
	var conf ConfigFile                     // Parsed config placeholder.
	err = json.Unmarshal(raw_config, &conf) // Convert JSON to structure.
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
