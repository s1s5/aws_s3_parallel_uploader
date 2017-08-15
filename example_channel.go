package main
import (
	"bytes"
	"fmt"
)

func consumer(message <-chan *bytes.Buffer, done chan<- bool) {
	// for {
	// 	bb, ok := <- message
		// if !ok {
		// 	return
	// }
	for bb := range message {
		fmt.Println("got bb", bb.String())
	}
	done <- true
}

func main() {
	message := make(chan *bytes.Buffer, 1)
	done := make(chan bool)
	go consumer(message, done)

	for i := 0; i < 10; i++ {
		bb := new(bytes.Buffer)
		bb.Write([]byte{0x61, 0x62, 0x63})
		fmt.Println("putting bb")
		message <- bb
	}
	close(message)
	fmt.Println("closed")
	<- done
	fmt.Println("finished")
}

