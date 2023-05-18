package main

import (
	"bytes"
	"flag"
	"fmt"
	"imagetools/config"
	"io"
	"log"
	"os"
	"path/filepath"
)

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
