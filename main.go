package main

import "fmt"
import "github.com/jaymell/gosmartcam/frameReader"

func main() {
	videoSource := "/dev/video3"
	fr, _ := frameReader.NewFrameReader(videoSource)
	fmt.Println(fr)
}
