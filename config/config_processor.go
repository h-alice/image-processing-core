package config

import (
	"errors"
	"fmt"
	op "imagetools/operation" // Grab `EncoderOption` from operation package.
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var ErrNotImplemented = errors.New("operation not implemented")

// OutputOptionConfig structure for output file.
type OutputOptionConfig op.EncoderOption

// Config structure for output file.
//
// Format: Output file format. Either `jpeg` or `png`.
//
// NameSuffix: Suffix of the output file name.
//
// NamePrefix: Prefix of the output file name.
//
// Options: Encoder option. For jpeg use and supports only `Quality` option.
type OutputConfig struct {
	Format     string              `yaml:"format"`     // Output file format
	NameSuffix string              `yaml:"nameSuffix"` // Output file name suffix
	NamePrefix string              `yaml:"namePrefix"` // Output file name prefix
	Options    *OutputOptionConfig `yaml:"options"`    // Encoder option
}

// Config structure for resizing image.
//
// Width: Output image width.
//
// Height: Output image height.
//
// Factor: Resize factor.
//
// Algorithm: Resize algorithm. Either `nearestneighbor`, `catmullrom`, or `approxbilinear`.
//
// NOTE: The `Factor` is prioritized over `Width` and `Height`.
type ResizeConfig struct {
	Width     int     `yaml:"width"`     // Output image width
	Height    int     `yaml:"height"`    // Output image height
	Factor    float32 `yaml:"factor"`    // Resize factor
	Algorithm string  `yaml:"algorithm"` // Resize algorithm
}

// Config structure for processing profile.
//
// ProfileName: Profile identifier.
//
// ICC: ICC profile to embed.
//
// Resize: Resize option.
//
// Output: Output file configuration.
type ProcessProfileConfig struct {
	ProfileName string        `yaml:"profileName"` // Profile identifier
	ICC         string        `yaml:"icc"`         // ICC profile to embed
	Resize      *ResizeConfig `yaml:"resize"`      // Resize option
	Output      *OutputConfig `yaml:"output"`      // Output file configuration
}

// Currently not used.
type OutputDirConfig struct {
	DirName *string `yaml:"dirName"` // Output directory name
}

// Config structure for config file.
//
// Profiles: List of profile configurations.
type ConfigFileRoot struct {
	Profiles []ProcessProfileConfig `yaml:"profiles"` // List of profile configurations
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

// Load config file from path.
func LoadConfigFromFile(config_path string) (*ConfigFileRoot, error) {

	raw_config, err := os.ReadFile(config_path) // Read raw config file.
	if err != nil {
		return nil, err
	}

	// Converting JSON to config structure.
	var conf ConfigFileRoot                 // Parsed config placeholder.
	err = yaml.Unmarshal(raw_config, &conf) // Convert JSON to structure.
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func (pf ConfigFileRoot) ConfigPrettyPrint(config *ConfigFileRoot) string {

	var output string = "" // Placeholder for output.

	for _, pf := range config.Profiles { // Iterate through profiles.

		output += "Profile Name: " + pf.ProfileName + "\n"
		output += "\tICC Profile: " + pf.ICC + "\n"

		if pf.Resize != nil {
			output += "\tResizing Configuration:\n"
			output += "\t\tResize Width: " + fmt.Sprintf("%d", pf.Resize.Width) + "\n"
			output += "\t\tResize Height: " + fmt.Sprintf("%d", pf.Resize.Height) + "\n"
			output += "\t\tResize Factor: " + fmt.Sprintf("%f.2", pf.Resize.Factor) + "\n"
			output += "\t\tResize Algorithm: " + pf.Resize.Algorithm + "\n"
		}

		if pf.Output != nil {
			output += "\tOutput Configuration:\n"
			output += "\t\tOutput Format: " + pf.Output.Format + "\n"
			output += "\t\tOutput Name Prefix: " + pf.Output.NamePrefix + "\n"
			output += "\t\tOutput Name Suffix: " + pf.Output.NameSuffix + "\n"
			if pf.Output.Options != nil {
				output += "\t\tEncoder Quality: " + fmt.Sprintf("%d", pf.Output.Options.Quality) + "\n"
			}
		}
	}
	return output
}
