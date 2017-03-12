package gosmartcam

import "fmt"
import "log"
import "time"
import "github.com/lazywei/go-opencv/opencv"

type CV2FrameDiffMotionDetector struct {
	background *OpenCVFrame
	current    *OpenCVFrame
}

func (md *CV2FrameDiffMotionDetector) DetectMotion() (contours *opencv.Seq) {
	return
}

func (md *CV2FrameDiffMotionDetector) SetCurrent(frame *OpenCVFrame) {
	if md.current != nil {
		md.background = md.current		
	}
	md.current = frame
}

func (md *CV2FrameDiffMotionDetector) Delta() (delta *opencv.IplImage) {
	if md.background == nil || md.current == nil {
		return 
	}
	delta = md.current.image.Clone()
	opencv.AbsDiff(md.background.image, md.current.image, delta)
	return
}

// This type handles the
type OpenCVMotionRunner struct {
	md             MotionDetector
	imageChan      FrameChan
	lastMotionTime time.Time
	motionTimeout  uint
	videoWriter    VideoWriter
	videoBuffer    []*OpenCVFrame
	frame          *OpenCVFrame
}

func NewOpenCVMotionRunner(md MotionDetector,
	imageChan FrameChan,
	motionTimeout uint,
	videoWriter VideoWriter) *OpenCVMotionRunner {

	var lastMotionTime time.Time
	var videoBuffer []*OpenCVFrame
	var frame *OpenCVFrame

	return &OpenCVMotionRunner{
		md:             md,
		imageChan:      imageChan,
		lastMotionTime: lastMotionTime,
		motionTimeout:  motionTimeout,
		videoWriter:    videoWriter,
		videoBuffer:    videoBuffer,
		frame:          frame,
	}
}

// Not currently used
func (mr *OpenCVMotionRunner) getBSFrame() *OpenCVFrame {
	f := mr.imageChan.PopFrame()
	frame := f.(*BSFrame)
	return frame.ToOpenCVFrame()
}

// Not currently used
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

	barf := mr.md.(*CV2FrameDiffMotionDetector)
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
		barf.SetCurrent(mr.frame)
		delta := barf.Delta()
		// win.ShowImage(mr.frame.image)
		if delta != nil {
			win.ShowImage(delta)
			opencv.WaitKey(1)
		} else {
			fmt.Println("wtf")
		}

	}

	return nil
}

func (mr *OpenCVMotionRunner) HandleMotion() {

}

func (mr *OpenCVMotionRunner) HandleMotionTimeout() {

}
