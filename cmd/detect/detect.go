package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/zyedidia/ftdetect"
)

var detectorFile = flag.String("detector", "", "detector file")

func main() {
	flag.Parse()
	args := flag.Args()

	var ds ftdetect.Detectors
	if *detectorFile != "" {
		data, err := ioutil.ReadFile(*detectorFile)
		if err != nil {
			log.Fatal(err)
		}
		ds, err = ftdetect.LoadDetectors(data)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		ds = ftdetect.LoadDefaultDetectors()
	}

	if len(args) <= 0 {
		log.Fatal("no file provided")
	}

	f, err := os.Open(args[0])
	var first []byte
	if err == nil {
		// if the file exists, we can extract the header to improve the guess.
		r := bufio.NewReader(f)
		first, _ = r.ReadSlice('\n')
		first = bytes.TrimSpace(first)
	}
	d := ds.Detect(args[0], first)

	if d == nil {
		fmt.Println("unknown")
	} else {
		fmt.Println(d.Name)
	}
}
