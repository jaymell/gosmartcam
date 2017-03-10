package motion

import "log"
import "time"
import "github.com/lazywei/go-opencv/opencv"
import "github.com/jaymell/gosmartcam/frameReader"
import "github.com/jaymell/gosmartcam/videoWriter"

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

type OpenCVFrame struct {
	*frameReader.Frame
	IplImage *opencv.IplImage
}

type OpenCVMotionRunner struct {
	motionDetector *MotionDetector
	imageQueue     chan *frameReader.Frame
	lastMotionTime time.Time
	motionTimeout  uint
	videoWriter    *videoWriter.VideoWriter
	videoBuffer    []frameReader.Frame
	frame          *frameReader.Frame
}

func NewOpenCVFrame(f *frameReader.Frame) *OpenCVFrame {
	return &OpenCVFrame{
		frameReader.Frame: f,
		IplImage: opencv.DecodeImageMem(frame.Image),
	}
}

func NewOpenCVMotionRunner(motionDetector *MotionDetector,
	imageQueue chan *frameReader.Frame,
	motionTimeout uint,
	videoWriter *videoWriter.VideoWriter) *OpenCVMotionRunner {

	var lastMotionTime time.Time
	var videoBuffer []frameReader.Frame
	var frame *frameReader.Frame

	return &{
		motionDetector: motionDetector,
		imageQueue: imageQueue,
		lastMotionTime: lastMotionTime,
		motionTimeout: motionTimeout,
		videoWriter: videoWriter,
		videoBuffer: videoBuffer,
		frame: frame,
	}
}

func (mr *OpenCVMotionRunner) Run() {
	log.Println("Starting motion detection... ")
	inMotion := false
	var frame OpenCVFrame
	for {
		f := <- imageQueue
        frame = NewOpenCVFrame(f)
	}
}

func (mr *OpenCVMotionRunner) HandleMotion() {

}

func (mr *OpenCVMotionRunner) HandleMotionTimeout() {

}


