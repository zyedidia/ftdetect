//go:build !ftdetect_custom

package ftdetect

import _ "embed"

//go:embed detectors/detectors.dat
var defaultDetectors []byte

// LoadDefaultDetectors returns a set of detectors for many programming languages.
func LoadDefaultDetectors() Detectors {
	d, _ := LoadDetectors(defaultDetectors)
	return d
}
