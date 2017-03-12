package gosmartcam

import "time"
import "fmt"
import "github.com/lazywei/go-opencv/opencv"
import "github.com/blackjack/webcam"

// abstract type for frames
type Frame interface {
	Image() interface{}
	Copy() Frame
}

// Frame with byte slice image
type BSFrame struct {
	image  []byte
	Time   time.Time
	Width  uint32
	Height uint32
}

func (f *BSFrame) Image() interface{} {
	return f.image
}

func (f *BSFrame) ToOpenCVFrame() *OpenCVFrame {
	img := opencv.DecodeImageMem(f.image)
	return &OpenCVFrame{
		image:  img,
		Time:   f.Time,
		Width:  f.Width,
		Height: f.Height,
	}
}

func (f *BSFrame) Copy() Frame {
	newImage := make([]byte, len(f.image))
	fmt.Println("length: ", len(f.image))
	copy(newImage, f.image)
	fmt.Println("length: ", len(newImage))

	return &BSFrame{
		image:  newImage,
		Time:   f.Time,
		Width:  f.Width,
		Height: f.Height,
	}
}

// Frame using OpenCV's default image type
type OpenCVFrame struct {
	image  *opencv.IplImage
	Time   time.Time
	Width  uint32
	Height uint32
}

func (f *OpenCVFrame) Image() interface{} {
	return f.image
}

func (f *OpenCVFrame) Copy() Frame {
	newImage := *f.image
	return &OpenCVFrame{
		image:  &newImage,
		Time:   f.Time,
		Width:  f.Width,
		Height: f.Height,
	}
}


// FrameReader interface defines the object
// that reads frames from camera in a loop
type FrameReader interface {
	ReadFrame() (Frame, error)
	Run()
}

// FrameReader based on Blackjack's webcam library
type BJFrameReader struct {
	cam         *webcam.Webcam
	width       uint32
	height      uint32
	pixelFormat webcam.PixelFormat
	frameChan   BSFrameChan
	fps         float32
}

// MotionDetector is the interface implemented by
// various motion detection algorithms
type MotionDetector interface {
	DetectMotion() *opencv.Seq
}

// MotionRunner is the interface for
// the object that runs the motion detection
// loop and has associated methods for drawing contours
// and handling video
type MotionRunner interface {
	Run()
	HandleMotion()
	HandleMotionTimeout()
}

type OvenCVFrameDiffMotionDetector struct {
}

type FrameChan interface {
	PopFrame() Frame
	PushFrame(Frame)
}

type BSFrameChan chan *BSFrame

func (fc BSFrameChan) PopFrame() (frame Frame) {
	frame = <-fc
	return
}

func (fc BSFrameChan) PushFrame(f Frame) {
	frame := f.(*BSFrame)
	fc <- frame
}

type OpenCVFrameChan chan *OpenCVFrame

func (fc OpenCVFrameChan) PopFrame() (frame Frame) {
	frame = <-fc
	return
}

func (fc OpenCVFrameChan) PushFrame(f Frame) {
	frame := f.(*OpenCVFrame)
	fc <- frame
}
