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

func TestDetectUnnamed(t *testing.T) {
	// Detectors with no exts or files share the "" key, which an unnamed
	// buffer also looks up. That lookup alone must not count as a match.
	gitd := &ftdetect.Detector{
		File: ftdetect.MustRegex("COMMIT_EDITMSG$"),
		Name: "git",
	}
	shd := &ftdetect.Detector{
		Header: ftdetect.MustRegex("^#!.*/(env +)?(ba)?(a)?(mk)?sh( |$)"),
		Name:   "shell",
	}

	ds := make(ftdetect.Detectors)
	ds.RegisterDetector(gitd)
	ds.RegisterDetector(shd)

	if d := ds.Detect("", nil); d != nil {
		t.Errorf("empty unnamed buffer detected as %s", d.Name)
	}
	if d := ds.Detect("", []byte("hello")); d != nil {
		t.Errorf("unnamed buffer detected as %s", d.Name)
	}
	if d := ds.Detect("", []byte("#!/bin/sh")); d == nil || d.Name != "shell" {
		t.Error("unnamed buffer with shebang not detected as shell")
	}
}
