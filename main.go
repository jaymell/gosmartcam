package main

import "encoding/json"
import "fmt"
import "log"
import "os"
import "github.com/jaymell/gosmartcam/gosmartcam"
import "github.com/lazywei/go-opencv/opencv"


const FRAME_BUF_SIZE = 512

type config struct {
	CaptureFormat string
	VideoSource   string
	FPS float32
	MotionTimeout uint

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

func writeTestJpeg1(fReader gosmartcam.BJFrameReader) (error) {
	frame, err := fReader.ReadFrame()
	if err != nil {
		return fmt.Errorf("Failed to read frame: %v", err)
	}
	img := frame.Image().([]byte)
	jpg, err := gosmartcam.ByteSlicetoJpeg(img)
	if err != nil {
		return fmt.Errorf("gosmartcam.ByteSlicetoJpeg failed: %v", err)
	}
	newJpg := opencv.FromImage(*jpg)
	opencv.SaveImage("/tmp/out.jpeg", newJpg, 0)

	return nil
}

func writeTestJpeg2(fReader gosmartcam.BJFrameReader) (error) {
	frame, err := fReader.ReadFrame()
	if err != nil {
		return fmt.Errorf("Failed to read frame: %v", err)
	}	
	img := frame.Image().([]byte)
	jpg := opencv.DecodeImageMem(img)
	opencv.SaveImage("/tmp/out.jpeg", jpg, 0)

	return nil
}

func dumpFrametoFile(fReader gosmartcam.BJFrameReader) (error) {
	frame, err := fReader.ReadFrame()
	if err != nil {
		return fmt.Errorf("Failed to read frame: %v", err)
	}
	f, err := os.Create("/tmp/out.jpeg")
	if err != nil {
		return fmt.Errorf("Failed to create output file: %v", err)
	}
	img := frame.Image().([]byte)
	f.Write(img)
	return nil
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

    frameChan := make(gosmartcam.BSFrameChan, FRAME_BUF_SIZE)
    videoChan := make(gosmartcam.BSFrameChan, FRAME_BUF_SIZE)
    motionChan := make(gosmartcam.BSFrameChan, FRAME_BUF_SIZE)

	fReader, err := gosmartcam.NewBJFrameReader(cfg.VideoSource, 
		                                         cfg.CaptureFormat, 
		                                         "", 
		                                         cfg.FPS, 
		                                         frameChan)
	if err != nil {
		return fmt.Errorf("Unable to instantiate frame reader")
	}

    vw := gosmartcam.OpenCVVideoWriter{FPS: cfg.FPS}
	motionRunner := gosmartcam.NewOpenCVMotionRunner(motionChan,
		cfg.MotionTimeout,
		vw)

	go motionRunner.Run()

	go func(fReader gosmartcam.FrameReader, frameChan gosmartcam.FrameChan) {
		for {
			frame, err := fReader.ReadFrame()
			if err != nil {
				log.Println("Failed to read frame: %v", err)
			}	
			frameChan.PushFrame(frame)
		}
	}(fReader, frameChan)

	for {
		log.Println("getting frame")
		frame := frameChan.PopFrame()
		log.Println("got frame")
		videoChan.PushFrame(frame)
		motionChan.PushFrame(frame)
	}
	
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
