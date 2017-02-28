package frameReader

import "fmt"
import "os"
import "path"
import "github.com/blackjack/webcam"

type FrameReader struct {
  cam *webcam.Webcam
}


func NewFrameReader(videoSource string) (*FrameReader, error) {
	cam, err := webcam.Open(videoSource)

	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()

	format_desc := cam.GetSupportedFormats()
	var formats []webcam.PixelFormat
	for f := range format_desc {
		formats = append(formats, f)
	}

    _ = cam.StartStreaming()
	for i, value := range formats {
		fmt.Fprintf(os.Stderr, "[%d] %s\n", i+1, format_desc[value])
		sizes := []webcam.FrameSize(cam.GetSupportedFrameSizes(value)) 
		size := sizes[len(sizes)-1]
	    _, _, _, _ = cam.SetImageFormat(value, uint32(size.MaxWidth), uint32(size.MaxHeight))
        timeout := uint32(5)
        err = cam.WaitForFrame(timeout)
        frame, _ := cam.ReadFrame()
		f, _ := os.Create(path.Join("/tmp", format_desc[value]))
		f.Write(frame)
	}

	return &FrameReader{
		cam: cam,
	}, nil
}

