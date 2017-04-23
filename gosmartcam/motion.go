package gosmartcam

import "fmt"
import "log"
import "time"
import "unsafe"
import "github.com/lazywei/go-opencv/opencv"

// pre-processing image (e.g., downsampling, grayscale) before motion detection
func opencvPreProcessImage(img *opencv.IplImage) (*opencv.IplImage) {
    gray := opencv.CreateImage(img.Width(),
                img.Height(),
                img.Depth(), 
                1)
	opencv.CvtColor(img, gray, opencv.CV_BGR2GRAY)
	opencv.Smooth(gray, gray, opencv.CV_BLUR, 3, 3, 0, 0)
	return gray
}

// given *opencv.Seq and image, draw all the contours
func opencvDrawRectangles(img *opencv.IplImage, contours *opencv.Seq) {
	for c := contours; c != nil; c = c.HNext() {
		rect := opencv.BoundingRect(unsafe.Pointer(c))
		opencv.Rectangle(img, 
			opencv.Point{ rect.X(), rect.Y() }, 
			opencv.Point{ rect.X() + rect.Width(), rect.Y() + rect.Height() },
			opencv.ScalarAll(255.0), 
			1, 1, 0)
	}
}

// return contours that meet the threshold
func opencvFindContours(img *opencv.IplImage, threshold float64) *opencv.Seq {
	defaultThresh := 100.0
	if threshold == 0.0 {
		threshold = defaultThresh
	}

	contours := img.FindContours(opencv.CV_RETR_LIST, opencv.CV_CHAIN_APPROX_SIMPLE, opencv.Point{0, 0})

	// if contours == nil {
	// 	return nil
	// }

	// defer contours.Release()
	/* latest failure: 
	c := contours
	is_first := true
	for ; c != nil;  {
		next := c.HNext()
		if is_first != true {
			area := opencv.ContourArea(c, opencv.WholeSeq(), 0)
			if area < threshold {
				fmt.Println("jaymell releasing")
				c.Release()
			}			
		} else {
			is_first = false
		}
		c = next
	}
	*/
	return contours
}

// helper to determine text position
func positionText(width, height int) (int, int) {
	xPos := int(float64(width) * 0.05)
	yPos := int(float64(height) * 0.9)
	return xPos, yPos
}

func opencvPutText(image *opencv.IplImage, text string) {
	xPos, yPos := positionText(image.Width(), image.Height())
	color := opencv.NewScalar(255, 255, 255, 255)
	font := opencv.InitFont(opencv.CV_FONT_HERSHEY_SIMPLEX, 1, 1, 0, 1, 8)
	font.PutText(image, text, opencv.Point{xPos,yPos}, color)
}

type OpenCVFrameDiffMotionDetector struct {
	background *OpenCVFrame
	current    *OpenCVFrame
}

func (md *OpenCVFrameDiffMotionDetector) SetCurrent(frame *OpenCVFrame) {
	if md.current != nil {
		if md.background != nil {
			md.background.image.Release()	
		}
		md.background = md.current	
	}
	frame.image = opencvPreProcessImage(frame.image)
	md.current = frame

}

// return diff of current frame and background frame, else nil
func (md *OpenCVFrameDiffMotionDetector) Delta() (*opencv.IplImage) {
	if md.background == nil || md.current == nil {
		return nil
	}
	delta := opencv.CreateImage(md.current.image.Width(),
		md.current.image.Height(),
		md.current.image.Depth(), 
		md.current.image.Channels())
	opencv.AbsDiff(md.background.image, md.current.image, delta)
	return delta
}

func (md *OpenCVFrameDiffMotionDetector) DetectMotion() (contours *opencv.Seq) {
	delta := md.Delta()
	defer delta.Release()
	if delta == nil {
		return		
	}
    opencv.Threshold(delta, delta, float64(25), 255, opencv.CV_THRESH_BINARY)
    opencv.Dilate(delta, delta, nil, 2)
	// contours = delta.FindContours(opencv.CV_RETR_EXTERNAL, opencv.CV_CHAIN_APPROX_SIMPLE, opencv.Point{0, 0})
    contours = opencvFindContours(delta, 0.0)
	return
}

// OpenCV-focused implementation
// of the main motion detection loop 
type OpenCVMotionRunner struct {
	md             MotionDetector
	imageChan      FrameChan
	lastMotionTime time.Time
	motionTimeout  time.Duration
	videoWriter    VideoWriter
	videoBuffer    []*OpenCVFrame
	frame          *OpenCVFrame
}

func NewOpenCVMotionRunner(md MotionDetector,
	imageChan FrameChan,
	motionTimeout string,
	videoWriter VideoWriter) (*OpenCVMotionRunner, error) {

	var lastMotionTime time.Time
	var frame *OpenCVFrame

    dur, err := time.ParseDuration(motionTimeout)
    if err != nil {
		return nil, fmt.Errorf("Unable to parse timeout entry")
    }

	// FIXME: better way to initialize frame buffer:
	videoBuffer := make([]*OpenCVFrame, 0, 400)

	return &OpenCVMotionRunner{
		md:             md,
		imageChan:      imageChan,
		lastMotionTime: lastMotionTime,
		motionTimeout:  dur,
		videoWriter:    videoWriter,
		videoBuffer:    videoBuffer,
		frame:          frame,
	}, nil
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

func (mr *OpenCVMotionRunner) handleMotion(contours *opencv.Seq, win *opencv.Window) {
	mr.lastMotionTime = mr.frame.Time()
	opencvDrawRectangles(mr.frame.image, contours)
	opencvPutText(mr.frame.image, mr.frame.Time().Format(time.RFC3339))
	win.ShowImage(mr.frame.image)
	opencv.WaitKey(1)
	mr.videoBuffer = append(mr.videoBuffer, mr.frame)
}

func (mr *OpenCVMotionRunner) handleMotionTimeout() {
	log.Println("writing video")
	var nilTime time.Time
	mr.lastMotionTime = nilTime
	fmt.Println("this is the buffer: ", mr.videoBuffer)
	for _, frame := range mr.videoBuffer {
		frame.image.Release()
	}
	mr.videoBuffer = make([]*OpenCVFrame, 0, 400)

}

func (mr *OpenCVMotionRunner) motionIsTimedOut() bool {
	if mr.lastMotionTime.IsZero() {
		return false
	}
	return mr.frame.Time().Sub(mr.lastMotionTime) >= mr.motionTimeout
}

func (mr *OpenCVMotionRunner) Run() error {
	log.Println("Starting motion detection... ")

	// inMotion := false
	win := opencv.NewWindow("Live Feed")
	var motionWin *opencv.Window

	defer win.Destroy()

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
		
		win.ShowImage(mr.frame.image)
		opencv.WaitKey(1)

		md.SetCurrent(mr.frame.Copy().(*OpenCVFrame))
		contours := md.DetectMotion()
		if contours != nil {
			log.Println("motion detected")
			motionWin = opencv.NewWindow("Motion Detected")
			mr.handleMotion(contours, motionWin)
			contours.Release()
		} else if mr.motionIsTimedOut() == true {
			mr.handleMotionTimeout()
			motionWin.Destroy()
		} else if mr.lastMotionTime.IsZero() == false {
			mr.videoBuffer = append(mr.videoBuffer, mr.frame)
		} else {
			mr.frame.image.Release()			
		}
	}
	return nil
}
