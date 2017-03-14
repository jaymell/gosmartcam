package main

import "encoding/json"
import "fmt"
import "log"
import "os"
import "github.com/jaymell/gosmartcam/gosmartcam"
import _ "net/http/pprof"
import "net/http"

const FRAME_BUF_SIZE = 1024

type config struct {
	CaptureFormat string
	VideoSource   string
	FPS           float32
	MotionTimeout uint
	MotionDetector string
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

func loadMotionDetector(cfg *config) (md gosmartcam.MotionDetector, err error) {
	switch m := cfg.MotionDetector; m {
	default:
		err = fmt.Errorf("Unknown MotionDetector type")
		return
	case "CV2FrameDiffMotionDetector":
		md = &gosmartcam.CV2FrameDiffMotionDetector{}
		return
	}
}

func run() (err error) {

	f, err := os.Open("config.js")
	if err != nil {
		return fmt.Errorf("Unable to open config file: ", err)
	}

	cfg, err := loadConfig(f)
	if err != nil {
		return fmt.Errorf("Unable to load config: ", err)
	}

	frameChan := make(gosmartcam.BSFrameChan, FRAME_BUF_SIZE)
	// videoChan := make(gosmartcam.BSFrameChan, FRAME_BUF_SIZE)
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
	md, err := loadMotionDetector(cfg)
	if err != nil {
		return
	}
	motionRunner := gosmartcam.NewOpenCVMotionRunner(md,
		motionChan,
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
		// videoChan.PushFrame(frame.Copy())
		motionChan.PushFrame(frame)

		// TODO: figure out exactly why
		// this breaks so weirdly:
		// videoChan.PushFrame(frame)
		// motionChan.PushFrame(frame)
	}

	return nil
}

func main() {
	
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	
	fmt.Println("Starting...")
	err := run()
	if err != nil {
		fmt.Println("Failed: ", err)
		os.Exit(1)
	}
}
