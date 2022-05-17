package ftdetect

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"path/filepath"
	"regexp"
)

type regex struct {
	compiled *regexp.Regexp
	Regex    string
}

// Regex compiles a serializable regular expression.
func Regex(s string) (*regex, error) {
	r, err := regexp.Compile(s)
	if err != nil {
		return nil, err
	}
	return &regex{
		Regex:    s,
		compiled: r,
	}, nil
}

// MustRegex is equivalent to Regex but panics if the regular expression is invalid.
func MustRegex(s string) *regex {
	r, err := Regex(s)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *regex) Match(b []byte) bool {
	if r == nil {
		return false
	}
	if r.compiled == nil {
		var err error
		r.compiled, err = regexp.Compile(r.Regex)
		if err != nil {
			return false
		}
	}
	return r.compiled.Match(b)
}

func (r *regex) MatchString(s string) bool {
	if r == nil {
		return false
	}
	if r.compiled == nil {
		var err error
		r.compiled, err = regexp.Compile(r.Regex)
		if err != nil {
			return false
		}
	}
	return r.compiled.MatchString(s)
}

// A Detector defines a language and how it should be detected via extensions,
// special files, a file regex, and a header regex. A language does not need to
// provide every detection mechanism. In fact, most languages only need to
// provide an extension, which makes detection very efficient.
type Detector struct {
	Exts     []string
	Files    []string
	File     *regex
	Header   *regex
	Priority int // 0 is lowest priority
	Name     string
}

// Detectors is a set of languages that are supported. It is a map from
// extension/filename to the list of detectors that are registered for those
// values.
type Detectors map[string][]*Detector

// LoadDetectors loads a set of languages from a serialized Detectors byte slice.
func LoadDetectors(b []byte) (Detectors, error) {
	ds := make(Detectors)
	fz, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return ds, err
	}
	dec := gob.NewDecoder(fz)
	err = dec.Decode(&ds)
	fz.Close()

	return ds, err
}

// LoadDetectorJson loads a detector from a json spec.
func LoadDetectorJson(data []byte) (*Detector, error) {
	var d Detector
	err := json.Unmarshal(data, &d)
	return &d, err
}

// Serialize writes the detector set to a byte slice so it can be saved.
func (ds Detectors) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	fz := gzip.NewWriter(&buf)
	enc := gob.NewEncoder(fz)
	err := enc.Encode(ds)
	fz.Close()
	return buf.Bytes(), err
}

// RegisterDetector adds a new detector to the set.
func (ds Detectors) RegisterDetector(d *Detector) {
	for _, ext := range d.Exts {
		ds[ext] = append(ds[ext], d)
	}
	for _, f := range d.Files {
		ds[f] = append(ds[f], d)
	}
	if len(d.Files) == 0 && len(d.Exts) == 0 {
		ds[""] = append(ds[""], d)
	}
}

// SetPriority sets the priority for a certain language. A higher priority will
// make this language match above others when there are match conflicts.
func (ds Detectors) SetPriority(lang string, priority int) {
	for _, arrs := range ds {
		for _, d := range arrs {
			if d.Name == lang {
				d.Priority = priority
				return
			}
		}
	}
}

// Detect returns the language that was detected from the filename and file
// header (first line of file), or nil if no matching language was found.
func (ds Detectors) Detect(filename string, header []byte) *Detector {
	ext := filepath.Ext(filename)
	arr, ok := ds[filename]
	if !ok {
		arr, ok = ds[ext]
	}
	if !ok {
		for _, arrs := range ds {
			arr = append(arr, arrs...)
		}
	}
	if len(arr) == 1 {
		// no conflicts
		return arr[0]
	}

	var best *Detector
	for _, d := range arr {
		if ok || d.File.MatchString(filename) || d.Header.Match(header) {
			if best == nil || d.Priority > best.Priority {
				best = d
			}
		}
	}
	return best
}
