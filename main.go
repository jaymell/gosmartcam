package main

import "encoding/json"
import "fmt"
import "os"
import "github.com/jaymell/gosmartcam/frameReader"
import "github.com/jaymell/gosmartcam/util"
import "github.com/lazywei/go-opencv/opencv"

type config struct {
	CaptureFormat string
	VideoSource   string
	FPS float32
}

const FRAME_BUF_SIZE = 256

func loadConfig(f *os.File) (*config, error) {

	decoder := json.NewDecoder(f)
	config := config{}

	err := decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode json: ", err)
	}

	return &config, nil
}

func writeTestJpeg(cfg config) {
	fReader, err := frameReader.NewBJFrameReader(cfg.VideoSource, 
		                                         cfg.CaptureFormat, 
		                                         "", 
		                                         cfg.FPS, 
		                                         frameQueue)
	frame, _ := fReader.GetFrame()
	jpg, err := util.ByteSlicetoJpeg(frame)
	opencv.SaveImage("/tmp/out.jpeg", jpg, 0)
}


func run() error {

	f, err := os.Open("config.js")
	if err != nil {
		return fmt.Errorf("Unable to open config file: ", err)
	}

	cfg, err := loadConfig(f)
	if err != nil {
		return fmt.Errorf("Unable to load config: ", err)
	}

    frameQueue := make(chan *frameReader.Frame, FRAME_BUF_SIZE)
	fReader, err := frameReader.NewBJFrameReader(cfg.VideoSource, 
		                                         cfg.CaptureFormat, 
		                                         "", 
		                                         cfg.FPS, 
		                                         frameQueue)
	if err != nil {
		return fmt.Errorf("Unable to instantiate frame reader")
	}

	writeTestJpeg(cfg)
	
    //fReader.Test()
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
