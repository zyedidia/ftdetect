// +build ignore

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/zyedidia/ftdetect"
)

func main() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	ds := make(ftdetect.Detectors)
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}

		data, err := ioutil.ReadFile(f.Name())
		if err != nil {
			log.Printf("%s: read: %v\n", f.Name(), err)
			continue
		}

		var d ftdetect.Detector
		err = json.Unmarshal(data, &d)
		if err != nil {
			log.Printf("%s: unmarshal: %v\n", f.Name(), err)
			continue
		}

		ds.RegisterDetector(&d)
	}

	data, err := ds.Serialize()
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile("detectors.dat", data, 0666)
}
