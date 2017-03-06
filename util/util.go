package util

import "bytes"
import "fmt"
import "image/jpeg"


func ByteSlicetoJpeg(jSlice []byte) (*image.Image, error) {
	jReader := bytes.NewReader(jSlice)
	jpg, err := jpeg.Decode(jReader)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode jpeg") 
	}
    return &jpg, nil
}
