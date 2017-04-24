package gosmartcam

// import "fmt"
import "strconv"
import "github.com/lazywei/go-opencv/opencv"

type VideoWriter interface {
	WriteVideo(interface{}) error
}

type OpenCVVideoWriter struct {
	// Format string
	// CloudWriter
	// Path
	FPS float32
}

func (w OpenCVVideoWriter) WriteVideo(buffer interface{}) error {
	buf := buffer.([]*OpenCVFrame)
	// FIXME: don't hard code this!
	fourcc , _ := strconv.Atoi("MP42")
	vw := opencv.NewVideoWriter("/tmp/out.avi",
			fourcc,
			w.FPS,
			int(buf[0].Width),
			int(buf[0].Height),
			1)
	for _, v := range buf {
		_ = vw.WriteFrame(v.image)
	}
	vw.Release()
	return nil
}
