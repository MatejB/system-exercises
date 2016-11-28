package main

import (
	"flag"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	a := flag.Bool("a", false, "Append")
	flag.Parse()

	flags := unix.O_WRONLY | unix.O_CREAT
	if *a {
		flags |= unix.O_APPEND
	} else {
		flags |= unix.O_TRUNC
	}

	fn := os.Args[len(os.Args)-1]
	fd, err := unix.Open(fn, flags, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := unix.Close(fd); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		b := make([]byte, 256)

		// read stdin
		n, err := unix.Read(int(os.Stdin.Fd()), b)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			return
		}

		//write stdout
		_, err = unix.Write(int(os.Stdout.Fd()), b)
		if err != nil {
			log.Fatal(err)
		}

		// write to file
		_, err = unix.Write(fd, b)
		if err != nil {
			log.Fatal(err)
		}
	}
}
