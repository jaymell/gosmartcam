package gosmartcam

import "fmt"
import "log"
import "time"
import "github.com/lazywei/go-opencv/opencv"


type OpenCVMotionRunner struct {
	// motionDetector *MotionDetector
	imageChan      FrameChan
	lastMotionTime time.Time
	motionTimeout  uint
	videoWriter    VideoWriter
	videoBuffer    []*OpenCVFrame
	frame          *OpenCVFrame
}


//func NewOpenCVMotionRunner(motionDetector *MotionDetector,
func NewOpenCVMotionRunner(imageChan FrameChan,
	motionTimeout uint,
	videoWriter VideoWriter) *OpenCVMotionRunner {

	var lastMotionTime time.Time
	var videoBuffer []*OpenCVFrame
	var frame *OpenCVFrame

	return &OpenCVMotionRunner{
		// motionDetector: motionDetector,
		imageChan: imageChan,
		lastMotionTime: lastMotionTime,
		motionTimeout: motionTimeout,
		videoWriter: videoWriter,
		videoBuffer: videoBuffer,
		frame: frame,
	}
}

func (mr *OpenCVMotionRunner) getBSFrame() *OpenCVFrame {
	f := mr.imageChan.PopFrame()
	frame := f.(*BSFrame)
	return frame.ToOpenCVFrame()
}

func (mr *OpenCVMotionRunner) getOpenCVFrame() *OpenCVFrame {
	f := mr.imageChan.PopFrame()
	frame := f.(*OpenCVFrame)
	return frame
}

func (mr *OpenCVMotionRunner) Run() error {
	log.Println("Starting motion detection... ")

	// inMotion := false
	win := opencv.NewWindow("GoOpenCV: VideoPlayer")
	defer win.Destroy()

	for {
		f := mr.imageChan.PopFrame()
		switch f := f.(type) {
		default:
			return fmt.Errorf("Unknown frame type")
		case *BSFrame:
			mr.frame = f.ToOpenCVFrame()
		case *OpenCVFrame:
			mr.frame = f
		}		
        win.ShowImage(mr.frame.image)
        opencv.WaitKey(1)
	}

	return nil
}

func (mr *OpenCVMotionRunner) HandleMotion() {

}

func (mr *OpenCVMotionRunner) HandleMotionTimeout() {

}
