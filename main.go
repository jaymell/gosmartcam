package main

import "encoding/json"
import "fmt"
import "os"
import "github.com/jaymell/gosmartcam/frameReader"

type config struct {
	CaptureFormat string
	VideoSource   string
}

func loadConfig(f *os.File) (*config, error) {

	decoder := json.NewDecoder(f)
	config := config{}

	err := decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode json: ", err)
	}

	return &config, nil
}

func run() error {

	f, err := os.Open("config.js")
	if err != nil {
		return fmt.Errorf("Unable to open config file: ", err)
	}

	config, err := loadConfig(f)
	if err != nil {
		return fmt.Errorf("Unable to load config: ", err)
	}

	fReader, err := frameReader.NewBJFrameReader(config.VideoSource, config.CaptureFormat, "")
	if err != nil {
		return fmt.Errorf("Unable to instantiate frame reader")
	}

    fReader.Test()
	return nil
}

func main() {

	fmt.Println("Starting...")
	err := run()
	if err != nil {
		fmt.Println("Failed: ", err)
		os.Exit(1)
	}
}
