package main
import (
	"bytes"
	"fmt"
)

type file_and_bb struct {
	name string
	buffer *bytes.Buffer
}

func consumer(message <-chan *file_and_bb, done chan<- bool) {
	// for {
	// 	bb, ok := <- message
		// if !ok {
		// 	return
	// }
	for bb := range message {
		fmt.Println("got bb", bb.name, ":", bb.buffer.String())
	}
	done <- true
}

func main() {
	message := make(chan *file_and_bb, 3)
	done := make(chan bool)
	go consumer(message, done)

	for i := 0; i < 10; i++ {
		fb := &file_and_bb{"hello", new(bytes.Buffer)}
		fb.buffer.Write([]byte{0x61, 0x62, 0x63})
		fmt.Println("putting bb")
		message <- fb
	}
	close(message)
	fmt.Println("closed")
	<- done
	fmt.Println("finished")
}

