// Program demonstrates atomic property of append flag of open kernel call.
// Usage: go run atomic-append.go [-a] test-file
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

func main() {
	fn := os.Args[len(os.Args)-1]
	if _, err := os.Stat(fn); err == nil {
		err := os.Remove(fn)
		if err != nil {
			log.Fatal(err)
		}
	}

	a := flag.Bool("a", false, "Append")
	flag.Parse()

	flags := unix.O_WRONLY | unix.O_CREAT
	if *a {
		flags |= unix.O_APPEND
	}

	passCh := make(chan struct{})
	var wg sync.WaitGroup

	write := func(line string) {
		defer wg.Done()

		<-passCh

		fd, err := unix.Open(fn, flags, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := unix.Close(fd); err != nil {
				log.Fatal(err)
			}
		}()

		// move to end
		_, err = unix.Seek(fd, 0, 2)
		if err != nil {
			log.Fatal(err)
		}

		_, err = unix.Write(fd, []byte(line+"\n"))
		if err != nil {
			log.Fatal(err)
		}
	}

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go write(fmt.Sprintf("Worker %d", i))
	}

	close(passCh)
	wg.Wait()

	content, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%s", content)
}
