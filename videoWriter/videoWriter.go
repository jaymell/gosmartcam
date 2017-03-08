package videoWriter

type VideoWriter interface {
	WriteVideo([]frameReader.Frame)	
}