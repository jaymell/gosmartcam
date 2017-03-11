package videoWriter

import "fmt"

type VideoWriter interface {
	WriteVideo([]interface{}) error
}

type OpenCVVideoWriter struct {
	// Format
	// CloudWriter
	// Path
	FPS float32
}

func (w OpenCVVideoWriter) WriteVideo(buf []interface{}) error {
	fmt.Println(buf)
	return nil
} 
