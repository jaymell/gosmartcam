package motion

import "github.com/lazywei/go-opencv/opencv"

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
}

