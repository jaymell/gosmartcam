package gosmartcam

import "fmt"
import "log"
import "time"
import "github.com/lazywei/go-opencv/opencv"

type CV2FrameDiffMotionDetector struct {
	background *OpenCVFrame
	current    *OpenCVFrame
	delta *opencv.IplImage // to prevent repeatedly invoking CreateImage
}

func (md *CV2FrameDiffMotionDetector) DetectMotion() (contours *opencv.Seq) {
	return
}

func (md *CV2FrameDiffMotionDetector) SetCurrent(frame *OpenCVFrame) {
	if md.current != nil {
		if md.background != nil {
			md.background.image.Release()	
		}
		md.background = md.current	
	}
	md.current = frame
}

func cv2preProcessFrame(src *OpenCVFrame) (processed *OpenCVFrame) {
	log.Println("processing frame")
	processed = src.Copy().(*OpenCVFrame)
    gray := opencv.CreateImage(src.image.Width(),
				src.image.Height(),
				src.image.Depth(), 
				1)
	defer gray.Release()

    blurred := opencv.CreateImage(src.image.Width(),
				src.image.Height(),
				src.image.Depth(), 
				1)

	// grayscale:
	log.Println("applying grayscale")	
	opencv.CvtColor(src.image, gray, opencv.CV_BGR2GRAY)

	// blur:
	log.Println("applying smoothing")	
	opencv.Smooth(gray, blurred, opencv.CV_BLUR, 3, 3, 0, 0)

	processed.image = blurred
	return
}

// return diff of current frame and background frame, else nil
func (md *CV2FrameDiffMotionDetector) Delta() (*opencv.IplImage) {
	if md.background == nil || md.current == nil {
		return nil
	}
	if md.delta == nil {
		md.delta = opencv.CreateImage(md.current.image.Width(),
								   md.current.image.Height(),
								   md.current.image.Depth(), 
								   md.current.image.Channels())
	}
	opencv.AbsDiff(md.background.image, md.current.image, md.delta)
	return md.delta
}

// OpenCV-focused implementation
// of the main motion detection loop 
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

	md := mr.md.(*CV2FrameDiffMotionDetector)
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

	    mdFrame := cv2preProcessFrame(mr.frame)
		md.SetCurrent(mdFrame)
		delta := md.Delta()

		if delta != nil {
			win.ShowImage(delta)
			opencv.WaitKey(1)
		} else {
			fmt.Println("wtf")
		}
		mr.frame.image.Release()
	}

	return nil
}

func (mr *OpenCVMotionRunner) HandleMotion() {

}

func (mr *OpenCVMotionRunner) HandleMotionTimeout() {

}
