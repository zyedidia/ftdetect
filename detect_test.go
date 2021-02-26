package ftdetect_test

import (
	"testing"

	"github.com/zyedidia/ftdetect"
)

func TestDetectors(t *testing.T) {
	god := &ftdetect.Detector{
		Exts:     []string{".go"},
		File:     nil,
		Header:   nil,
		Priority: 0,
		Name:     "go",
	}
	shd := &ftdetect.Detector{
		Exts:     []string{".sh"},
		Files:    []string{".shellcfg"},
		File:     ftdetect.MustRegex("(\\.bash|\\.ash|bashrc|bash_aliases|bash_functions|profile|bash-fc\\.|Pkgfile|pkgmk.conf|rc.conf|PKGBUILD|.ebuild\\$|APKBUILD)"),
		Header:   ftdetect.MustRegex("^#!.*/(env +)?(ba)?(a)?(mk)?sh( |$)"),
		Priority: 0,
		Name:     "shell",
	}

	ds := make(ftdetect.Detectors)
	ds.RegisterDetector(god)
	ds.RegisterDetector(shd)

	data, err := ds.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	loaded, err := ftdetect.LoadDetectors(data)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		filename string
		header   string
		filetype string
	}{
		{"hello.bash", "", "shell"},
		{"hello.go", "", "go"},
		{"test", "#!/bin/bash", "shell"},
		{".shellcfg", "", "shell"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			dsd := ds.Detect(tt.filename, []byte(tt.header))
			if dsd == nil || dsd.Name != tt.filetype {
				t.Error("ds detected wrong filetype")
			}
			loadedd := loaded.Detect(tt.filename, []byte(tt.header))
			if loadedd == nil || loadedd.Name != tt.filetype {
				t.Error("loaded detected wrong filetype")
			}
		})
	}
}
