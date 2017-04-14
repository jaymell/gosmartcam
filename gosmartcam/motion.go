package gosmartcam

import "fmt"
import "log"
import "time"
import "unsafe"
import "github.com/lazywei/go-opencv/opencv"

func opencvPreProcessFrame(src *OpenCVFrame) (processed *OpenCVFrame) {
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

// given *opencv.Seq and image, draw all the contours
func opencvDrawRectangles(img *opencv.IplImage, contours *opencv.Seq) {
	for ; contours != nil; contours = contours.HNext() {
		rect := opencv.BoundingRect(unsafe.Pointer(contours))
		log.Println("this1: ", rect.X(), rect.Y())
		log.Println("this2: ", rect.X() + rect.Width(), rect.Y() + rect.Height())
		opencv.Rectangle(img, 

			// opencv.Point{ rect.X() + rect.Width(), rect.Y() }, 
			// opencv.Point{ rect.X() , rect.Y() + rect.Height() },

			opencv.Point{ rect.X(), rect.Y() }, 
			opencv.Point{ rect.X() + rect.Width(), rect.Y() + rect.Height() },

			// opencv.Point{100, 50},
			// opencv.Point{200, 200},

			opencv.ScalarAll(255.0), 
			1, 1, 0)
	}
}

type OpenCVFrameDiffMotionDetector struct {
	background *OpenCVFrame
	current    *OpenCVFrame
	delta *opencv.IplImage // to prevent repeatedly invoking CreateImage
}

func (md *OpenCVFrameDiffMotionDetector) DetectMotion() (contours *opencv.Seq) {
	delta := md.Delta()
	if delta == nil {
		return		
	} 
	contours = delta.FindContours(opencv.CV_RETR_TREE, opencv.CV_CHAIN_APPROX_SIMPLE, opencv.Point{0, 0})
    opencv.Threshold(delta, delta, float64(100), 255, opencv.CV_THRESH_BINARY)
    opencv.Dilate(delta, delta, nil, 2)
	win := opencv.NewWindow("wtf")
	opencvDrawRectangles(delta, contours)
	win.ShowImage(delta)
	opencv.WaitKey(1)
	return
}

func (md *OpenCVFrameDiffMotionDetector) SetCurrent(frame *OpenCVFrame) {
	if md.current != nil {
		if md.background != nil {
			md.background.image.Release()	
		}
		md.background = md.current	
	}
	md.current = frame
}

// return diff of current frame and background frame, else nil
func (md *OpenCVFrameDiffMotionDetector) Delta() (*opencv.IplImage) {
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

func (mr *OpenCVMotionRunner) handleMotion(contours *opencv.Seq) {
	win := opencv.NewWindow("Motion Feed")
	opencvDrawRectangles(mr.frame.image, contours)
	win.ShowImage(mr.frame.image)
	opencv.WaitKey(1)
 	// optional: draw contours
}

func (mr *OpenCVMotionRunner) Run() error {
	log.Println("Starting motion detection... ")

	// inMotion := false
	// win := opencv.NewWindow("Live Feed")
	// defer win.Destroy()

	md := mr.md.(*OpenCVFrameDiffMotionDetector)
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

	    mdFrame := opencvPreProcessFrame(mr.frame)
		md.SetCurrent(mdFrame)
		contours := md.DetectMotion()
		if contours != nil {
			log.Println("got some contours")
			mr.handleMotion(contours)
		} 

		// win.ShowImage(mr.frame.image)
		// opencv.WaitKey(1)

		mr.frame.image.Release()
	}

	return nil
}

func (mr *OpenCVMotionRunner) HandleMotionTimeout() {

}
