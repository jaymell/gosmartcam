package frameReader

import "fmt"

import "os"
import "time"
import "github.com/blackjack/webcam"


type FrameReader struct {
	cam *webcam.Webcam
}


func NewFrameReader(videoSource string, captureFormat string) (*FrameReader, error) {
	cam, err := webcam.Open(videoSource)

	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()

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

	_ = cam.StartStreaming()
	sizes := []webcam.FrameSize(cam.GetSupportedFrameSizes(*cFormat))
	size := sizes[len(sizes)-1]
	// size := sizes[0]
	_, _, _, _ = cam.SetImageFormat(*cFormat, uint32(size.MaxWidth), uint32(size.MaxHeight))
	timeout := uint32(5)

	t1 := time.Now()

	var numFrames int = 100
	for i := 0; i < numFrames; i++ {

		err = cam.WaitForFrame(timeout)
		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Fprint(os.Stderr, err.Error())
			continue
		default:
			panic(err.Error())
		}

		_, err := cam.ReadFrame()
		if err != nil {
			panic("Error getting frame: ")
		}
	}
	diff := time.Since(t1)
	fmt.Printf("Time elapsed: %f\n", diff.Seconds())
	fmt.Printf("Frames per second: %f\n", float64(float64(numFrames)/diff.Seconds()))

	return &FrameReader{
		cam: cam,
	}, nil
}
