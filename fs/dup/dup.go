// Program demonstrates dup file descriptor behavior and Pwrite, Pread system calls to
// read write in file with offset.
// Usage: go run dup.go test
package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"encoding/binary"

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

	flags := unix.O_RDWR | unix.O_CREAT
	fd, err := unix.Open(fn, flags, 0644)
	if err != nil {
		log.Fatalf("26: %s", err)
	}
	defer func() {
		if err := unix.Close(fd); err != nil {
			log.Fatalf("30: %s", err)
		}
	}()

	var wg sync.WaitGroup

	worker := func(fd int, workAtEnd bool) {
		defer wg.Done()

		fdOwn, err := unix.Dup(fd)
		if err != nil {
			log.Fatalf("41: %s", err)
		}

		var o int64
		var m uint32 = 1
		if workAtEnd {
			o = 33
			m = 100
		}

		b := make([]byte, 32)

		binary.LittleEndian.PutUint32(b, 2*m)
		if _, err := unix.Pwrite(fdOwn, b, o); err != nil {
			log.Fatalf("55: %s", err)
		}

		r1 := make([]byte, 32)
		if _, err := unix.Pread(fdOwn, r1, o); err != nil {
			log.Fatalf("60: %s", err)
		}

		binary.LittleEndian.PutUint32(b, 3*m)
		if _, err := unix.Pwrite(fdOwn, b, o); err != nil {
			log.Fatalf("65: %s", err)
		}

		r2 := make([]byte, 32)
		if _, err := unix.Pread(fdOwn, r2, o); err != nil {
			log.Fatalf("70: %s", err)
		}

		binary.LittleEndian.PutUint32(b, 5*m)
		if _, err := unix.Pwrite(fdOwn, b, o); err != nil {
			log.Fatalf("75: %s", err)
		}

		r3 := make([]byte, 32)
		if _, err := unix.Pread(fdOwn, r3, o); err != nil {
			log.Fatalf("80: %s", err)
		}

		sum := binary.LittleEndian.Uint32(r1) + binary.LittleEndian.Uint32(r2) + binary.LittleEndian.Uint32(r3)
		fmt.Println(sum)
	}

	wg.Add(2)

	go worker(fd, false)
	go worker(fd, true)

	wg.Wait()
}
