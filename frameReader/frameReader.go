package frameReader

import "fmt"
import "time"
import "github.com/blackjack/webcam"

type Frame struct {
	Image []byte
	Time  time.Time
	Width uint32
	Height uint32
}

type FrameReader interface {
	GetFrame() (*Frame, error)
	Run()
}

type BJFrameReader struct {
	cam *webcam.Webcam
	width uint32
	height uint32
	pixelFormat webcam.PixelFormat
	frameQueue chan *Frame
	fps float32
}

func NewBJFrameReader(videoSource string, 
	                  captureFormat string, 
	                  size string, 
	                  fps float32, 
	                  frameQueue chan *Frame) (*BJFrameReader, error) {

	cam, err := webcam.Open(videoSource)

	if err != nil {
		panic(err.Error())
	}

	var cFormat *webcam.PixelFormat

	for i, v := range cam.GetSupportedFormats() {
		if v == captureFormat {
			cFormat = &i
			break
		}
	}

	if cFormat == nil {
		return nil, fmt.Errorf("CaptureFormat not supported")
	}

  var s webcam.FrameSize
	sizes := []webcam.FrameSize(cam.GetSupportedFrameSizes(*cFormat))
	if size == "" {
		s = sizes[len(sizes)-1]
	}
	// FIXME -- support passing non-empty size:
	// else { }

	code, width, height, err := cam.SetImageFormat(*cFormat, uint32(s.MaxWidth), uint32(s.MaxHeight))
	if err != nil {
		return nil, fmt.Errorf("Failed to set format/size")
	}

	// turn camera on
	err = cam.StartStreaming()
	if err != nil {
		return nil, fmt.Errorf("Failed to start streaming: %v", err)
	}

	return &BJFrameReader{
		cam: cam,
		width: width,
		height: height,
		pixelFormat: code,
		fps: fps,
		frameQueue: frameQueue,
	}, nil

}

func (fr *BJFrameReader) GetFrame() (*Frame, error) {

	timeout := uint32(5)

	err := fr.cam.WaitForFrame(timeout)
	switch err.(type) {
	case nil:
	case *webcam.Timeout:
		return nil, fmt.Errorf(err.Error())
	default:
		panic(err.Error())
	}

	frame, err := fr.cam.ReadFrame()
	if err != nil {
		return nil, fmt.Errorf("Error getting frame: %v", err)
	}

	return &Frame{
		Image: frame,
		Time: time.Now(),
		Width: fr.width,
		Height: fr.height,
	}, nil
    
}

func (fr *BJFrameReader) Run() {
	for {
		frame, err := fr.GetFrame()
		if err != nil {
			fmt.Println("Error getting frame: ", err.Error())
		} else {
			fr.frameQueue <- frame
		}
		d := time.Duration( 1 / fr.fps * float32(time.Second) )
		time.Sleep(d)
	}

}

func (fr *BJFrameReader) Test() {

	var numFrames int = 1000
	
	t1 := time.Now()
	for i := 0; i < numFrames; i++ {
		_, err := fr.GetFrame()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	diff := time.Since(t1)
	fmt.Printf("Time elapsed: %f\n", diff.Seconds())
	fmt.Printf("Frames per second: %f\n", float64(float64(numFrames)/diff.Seconds()))
}