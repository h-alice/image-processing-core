package config

import (
	"encoding/json"
	"errors"
	"image"
	icc "imagetools/icc"
	image_parser "imagetools/image_parser"
	op "imagetools/operation"
	"io"
	"log"
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

func (profile ProfileConfig) DoCrop(in image.Image) (image.Image, error) {
	return in, nil
}

// Embed ICC profile to image.
func (profile ProfileConfig) DoEmbedIcc(out io.Writer, in io.Reader) error {

	parsed_image, err := image_parser.Parse(in)
	if err != nil {
		return err
	}

	err = icc.EmbedIccProfile(profile.ICC, parsed_image)
	if err != nil {
		return err
	}

	_, err = parsed_image.WriteTo(out)

	return err
}

func (profile ProfileConfig) ProcessFile(out io.Writer, in io.Reader) error {

	// Procedure: Decode -> image ops -> encode -> segment ops -> write out

	working_image, err := op.CreateImageFromReader(in)
	if err != nil {
		log.Printf("[x] Error while creating image: %v", err)
		return err
	}

	// Do crop
	//img, err = profile.DoCrop(img)
	//if err != nil {
	//	log.Printf("[x] Error while cropping image: %v", err)
	//	return err
	//}

	output_image := working_image.
		Then(op.Decode()).
		ThenIf(profile.Resize.Factor != 0.0, op.ResizeImageByFactor(profile.Resize.Algorithm, profile.Resize.Factor)).
		ThenIf((profile.Resize.Factor == 0.0) && (profile.Resize.Width != 0), op.ResizeImageByWidth(profile.Resize.Algorithm, profile.Resize.Width)).
		ThenIf((profile.Resize.Factor == 0.0) && (profile.Resize.Width == 0) && (profile.Resize.Height != 0), op.ResizeImageByHeight(profile.Resize.Algorithm, profile.Resize.Height)).
		Then(op.Encode(profile.Output.Format, (*op.EncoderOption)(profile.Output.Options))).
		ThenIf((profile.Output.Format == "jpeg" || profile.Output.Format == "jpg") && profile.ICC != "", op.EmbedProfile(profile.ICC)) // Supports jpeg only for now, will be extended to other formats.

	if output_image.LastError() != nil {
		log.Printf("[x] Error while processing image: %v", output_image.LastError())
		return output_image.LastError()
	}

	// Write output to writer.
	output_image.Then(op.WriteImageToWriter(out))
	if output_image.LastError() != nil {
		log.Printf("[x] Error while writing image: %v", output_image.LastError())
		return output_image.LastError()
	}

	return nil
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
