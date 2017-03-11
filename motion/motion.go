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


type OpenCVMotionRunner struct {
	// motionDetector *MotionDetector
	imageQueue     chan *frameReader.Frame
	lastMotionTime time.Time
	motionTimeout  uint
	videoWriter    videoWriter.VideoWriter
	videoBuffer    []frameReader.Frame
	frame          *frameReader.Frame
}


type OvenCVFrameDiffMotionDetector struct {

}


//func NewOpenCVMotionRunner(motionDetector *MotionDetector,
func NewOpenCVMotionRunner(imageQueue chan *frameReader.Frame,
	motionTimeout uint,
	videoWriter videoWriter.VideoWriter) *OpenCVMotionRunner {

	var lastMotionTime time.Time
	var videoBuffer []frameReader.Frame
	var frame *frameReader.Frame

	return &OpenCVMotionRunner{
		// motionDetector: motionDetector,
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

	// inMotion := false
	win := opencv.NewWindow("GoOpenCV: VideoPlayer")
	defer win.Destroy()

	var frame *frameReader.OpenCVFrame
	for {
		log.Println("motion: getting frame")
		f := <- mr.imageQueue
        frame = f.ToOpenCVFrame()
        win.ShowImage(frame.Image)
        opencv.WaitKey(1)
	}
}

func (mr *OpenCVMotionRunner) HandleMotion() {

}

func (mr *OpenCVMotionRunner) HandleMotionTimeout() {

}


