package main

import (
	"bytes"
	"flag"
	"fmt"
	"imagetools/config"
	conf "imagetools/config"
	"io"
	"log"
	"os"
	"path/filepath"
)

type OutputOptionConfig operation.EncoderOption

type OutputConfig struct {
	Format     string
	NameSuffix string
	NamePrefix string
	Options    *OutputOptionConfig
}

type ResizeConfig struct {
	Width     int
	Height    int
	Factor    float32
	Algorithm string
}

type ProfileConfig struct {
	Profile string
	ICC     string
	Resize  *ResizeConfig
	Output  *OutputConfig
}

type OutputDirConfig struct {
	DirName *string
}

type ConfigFile struct {
	// Output   *OutputDirConfig
	Profiles []ProfileConfig
}

func generateOutputFileName(input_file string, output_profile *OutputConfig) (string, error) {

	var err error = nil
	if output_profile == nil {
		err = fmt.Errorf("no output section")
		return "", err
	}

	profile := *output_profile

	original_ext := filepath.Ext(input_file)
	original_name := filepath.Base(input_file)
	stem := original_name[:len(original_name)-len(original_ext)]

	fileSuffix := ""

	switch strings.ToLower(profile.Format) {
	case "jpeg":
		fileSuffix = ".jpg"
	case "":
		fileSuffix = original_ext
	default:
		fileSuffix = "." + profile.Format

	}
	full_file := profile.NamePrefix + stem + profile.NameSuffix + fileSuffix
	return full_file, nil
}

func processResize(in image.Image, profile *ResizeConfig) (out image.Image, err error) {
	err = nil
	if profile == nil {
		out = in
		return
	}

	algo := (*profile).Algorithm

	switch {
	case (*profile).Factor != 0.0:
		out = operation.ResizeImageByFactor(in, algo, (*profile).Factor)
		return

	case (*profile).Width != 0:
		out = operation.ResizeImageByWidth(in, algo, (*profile).Width)
		return

	case (*profile).Height != 0:
		out = operation.ResizeImageByHeight(in, algo, (*profile).Height)
		return

	default:
		err = fmt.Errorf("operation not implemented")
		return
	}
}

func embeddingIcc(w io.Writer, in io.Reader, profile ProfileConfig) (err error) {

	var input_image bytes.Buffer
	_, err = io.Copy(&input_image, in)
	if err != nil {
		log.Printf("[x] Can't copy image data to buffer: %s\n", err)
		return
	}
	input_reader := bytes.NewReader(input_image.Bytes())
	parsed_image, err := jpeg.ParseJpeg(input_reader)

	if err != nil {
		return
	}

	icc_data, err := icc.GetEmbeddedProfile(profile.ICC)
	if err != nil {
		return
	}

	seg_data, err := icc.JpegIccSegment(icc_data)
	if err != nil {
		return
	}

	segment := jpeg.NewJpegSegment(jpeg.APP2, seg_data)
	err = parsed_image.InsertIccSegment(segment)
	if err != nil {
		return
	}
	_, err = parsed_image.WriteTo(w)

	if err != nil {
		return
	}

	return
}

func processFile(out io.Writer, in io.Reader, profile ProfileConfig) (err error) {

	img, _, err := operation.Decode(in)
	if err != nil {
		log.Printf("[x] Error while decoding image: %v", err)
		return
	}

	// Do resize
	resized, err := processResize(img, profile.Resize)
	if err != nil {
		log.Printf("[x] Error while resizing image: %v", err)
		return
	}

	// Do encode
	var buf bytes.Buffer
	err = operation.Encode(&buf, &resized, profile.Output.Format, (*operation.EncoderOption)(profile.Output.Options))
	if err != nil {
		log.Printf("[x] Error while encoding image: %v", err)
		return
	}

	// Do embed ICC
	err = embeddingIcc(out, &buf, profile)
	if err != nil {
		log.Printf("[x] Error while embedding ICC profile: %v", err)
		return
	}

	return

}

// Get profile from home directory.
func defaultProfileFilePath(profile_name string) (path string, err error) {

	if profile_name == "" {
		profile_name = "default"
	}
	profile_name = fmt.Sprintf("%s.json", profile_name)

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	profile_dir := filepath.Join(home, ".imgtools")
	profile_file := filepath.Join(home, ".imgtools", profile_name)

	_, err = os.Stat(profile_dir)
	if err != nil {
		if os.IsNotExist(err) {

			// Making profile directory if not exists.
			err = os.Mkdir(profile_dir, 0777)
			if err != nil {
				return "", err
			}

			// Writing default profile.
			err = os.WriteFile(profile_file, []byte(default_profile), 0777)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return profile_file, err
}

func main() {

	// Parse command line.
	ptr_config_path := flag.String("config", "", "config file path")
	flag.Parse()

	// If path not specified, load defaultprofile from home directory.
	config_path := *ptr_config_path
	if *ptr_config_path == "" {

		log.Printf("[!] Using default config file.")

		var err error
		config_path, err = defaultProfileFilePath("default")
		if err != nil {
			log.Fatalf("[x] Error while getting default config: %s\n", err)
		}
	}

	conf, err := conf.ConfigLoader(config_path)
	if err != nil {
		log.Fatalf("[x] Cannot load config file: %s\n", err)
	}

	for _, f := range flag.Args() { // Iterate through input images.

		raw_bytes, err := os.ReadFile(f)
		if err != nil {
			log.Printf("[x] Error while reading file: %s\n", err)
			continue
		}

		tasks := make(chan struct{}, len(conf.Profiles))
		for _, pf := range conf.Profiles { // Apply all profile to input image.

			go func(pf config.ProfileConfig) {

				// TODO: Output dir.
				output_dir := filepath.Dir(f)
				if pf.Output == nil {
					log.Println("[x] No output section.")
					return
				}

				outfile_name := pf.Output.GenerateFileName(f)

				outputbuf := bytes.NewBuffer([]byte{})
				err = pf.ProcessFile(outputbuf, bytes.NewBuffer(raw_bytes))
				if err != nil {
					log.Printf("[x] An error occurred while processing image: %s\n", err)
					return
				}

				output_full_path := filepath.Join(output_dir, outfile_name)
				ofp, err := os.OpenFile(output_full_path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
				defer func() {
					ofp.Close()
				}()

				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("Writing output %s -> %s [%s]\n", f, output_full_path, pf.ProfileName)

				output_length := outputbuf.Len()

				written, err := io.Copy(ofp, outputbuf)

				if err != nil {
					log.Fatalln(err)
				} else if written != int64(output_length) {
					err = fmt.Errorf("written length mismatch")
					log.Fatal(err)
				}

				tasks <- struct{}{}
			}(pf)

		}

		for i := 0; i < len(conf.Profiles); i++ {
			log.Println(i)
			<-tasks
		}
	}
}

var default_profile string = `{"profiles": []}`
