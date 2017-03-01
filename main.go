package main

import "github.com/jaymell/gosmartcam/frameReader"

func main() {
	videoSource := "/dev/video1"
	_, _ = frameReader.NewFrameReader(videoSource)
}
