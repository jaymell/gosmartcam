package frameReader

import "fmt"
import "strconv"

import "os"
import "path"
import "github.com/blackjack/webcam"


type FrameReader struct {
  cam *webcam.Webcam
}


func NewFrameReader(videoSource string) (*FrameReader, error) {
	cam, err := webcam.Open(videoSource)

	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()

	format_desc := cam.GetSupportedFormats()
	var formats []webcam.PixelFormat
	for f := range format_desc {
		formats = append(formats, f)
	}

  // format := formats[1]
  // sizes := []webcam.FrameSize(cam.GetSupportedFrameSizes(format))
  // size := sizes[len(sizes)-1]
  // _, _, _, _ = cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
  // _ = cam.StartStreaming()
  //  timeout := uint32(5)
  //  err = cam.WaitForFrame(timeout)
  //  frame, _ := cam.ReadFrame()
  //  f, _ := os.Create(path.Join("/tmp", "out.jpg"))
  //  f.Write(frame)





    _ = cam.StartStreaming()
    for i:=0 ; i < 100; i++  {
	  format := formats[1]
	  sizes := []webcam.FrameSize(cam.GetSupportedFrameSizes(format))
	  size := sizes[len(sizes)-1]
	  _, _, _, _ = cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
	  _ = cam.StartStreaming()
	   timeout := uint32(5)
		err = cam.WaitForFrame(timeout)
		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Fprint(os.Stderr, err.Error())
			continue
		default:
			panic(err.Error())
		}

        frame, err := cam.ReadFrame()
        if err != nil {
        	panic("Error getting frame: ")
        }
        f, err := os.Create(path.Join("/tmp", strconv.Itoa(i) + ".jpg"))
		if err != nil {
			panic("Error creating file: ")
		}
		n, err := f.Write(frame)
		if err != nil {
			panic("Error writing file: ")
		}
		fmt.Printf("Wrote %d bytes\n", n)
		err = f.Close()
		if err != nil {
			panic("Error closing file: ")
		}

    
    }





	// for i, value := range formats {
	// 	fmt.Fprintf(os.Stderr, "[%d] %s\n", i+1, format_desc[value])
	// 	sizes := []webcam.FrameSize(cam.GetSupportedFrameSizes(value)) 
	// 	size := sizes[len(sizes)-1]
	// 	fmt.Fprintf(os.Stderr, "[%s]\n", size.GetString())
	//     _, _, _, _ = cam.SetImageFormat(value, uint32(size.MaxWidth), uint32(size.MaxHeight))
 //        timeout := uint32(5)

	// 	err = cam.WaitForFrame(timeout)
	// 	switch err.(type) {
	// 	case nil:
	// 	case *webcam.Timeout:
	// 		fmt.Fprint(os.Stderr, err.Error())
	// 		continue
	// 	default:
	// 		panic(err.Error())
	// 	}

 //        frame, err := cam.ReadFrame()
 //        if err != nil {
 //        	panic("Error getting frame: ")
 //        }
	// 	f, err := os.Create(path.Join("/tmp", format_desc[value]))
	// 	if err != nil {
	// 		panic("Error creating file: ")
	// 	}
	// 	n, err := f.Write(frame)
	// 	if err != nil {
	// 		panic("Error writing file: ")
	// 	}
	// 	fmt.Printf("Wrote %d bytes\n", n)
	// 	err = f.Close()
	// 	if err != nil {
	// 		panic("Error closing file: ")
	// 	}
	// }

	return &FrameReader{
		cam: cam,
	}, nil
}

