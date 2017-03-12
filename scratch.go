package main

import "encoding/json"
import "fmt"
import "log"
import "os"
import "github.com/jaymell/gosmartcam/gosmartcam"
import "github.com/lazywei/go-opencv/opencv"

func writeTestJpeg1(fReader gosmartcam.BJFrameReader) error {
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

func writeTestJpeg2(fReader gosmartcam.BJFrameReader) error {
	frame, err := fReader.ReadFrame()
	if err != nil {
		return fmt.Errorf("Failed to read frame: %v", err)
	}
	img := frame.Image().([]byte)
	jpg := opencv.DecodeImageMem(img)
	opencv.SaveImage("/tmp/out.jpeg", jpg, 0)

	return nil
}

func dumpFrametoFile(fReader gosmartcam.BJFrameReader) error {
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
