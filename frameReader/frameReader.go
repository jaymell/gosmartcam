package frameReader

import "fmt"
import "os"
import "time"
import "github.com/blackjack/webcam"

type Frame struct {
	image []byte
	time  time.Time
	width uint32
	height uint32
}

type FrameReader interface {
	GetFrame() (frame Frame, err error)
}

type BJFrameReader struct {
	cam *webcam.Webcam
	width uint32
	height uint32
	pixelFormat webcam.PixelFormat
}

func NewBJFrameReader(videoSource string, captureFormat string, size string) (*BJFrameReader, error) {

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

	code, width, height, err := cam.SetImageFormat(*cFormat, uint32(s.MaxWidth), uint32(s.MaxHeight))
	if err != nil {
		return nil, fmt.Errorf("Failed to set format/size")
	}

	return &BJFrameReader{
		cam: cam,
		width: width,
		height: height,
		pixelFormat: code,
	}, nil

}

func (fr *BJFrameReader) Test() {

	var numFrames int = 100
	timeout := uint32(5)
	
	_ = fr.cam.StartStreaming()
	t1 := time.Now()
	for i := 0; i < numFrames; i++ {

		err := fr.cam.WaitForFrame(timeout)
		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Fprint(os.Stderr, err.Error())
			continue
		default:
			panic(err.Error())
		}

		_, err = fr.cam.ReadFrame()
		if err != nil {
			panic("Error getting frame: ")
		}
	}
	diff := time.Since(t1)
	fmt.Printf("Time elapsed: %f\n", diff.Seconds())
	fmt.Printf("Frames per second: %f\n", float64(float64(numFrames)/diff.Seconds()))
}