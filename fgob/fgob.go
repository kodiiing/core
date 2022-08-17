// fgob implement a simple abstraction on top of gob encoding
// provided by https://pkg.go.dev/encoding/gob
package fgob

import (
	"bytes"
	"encoding/gob"
)

// Marshal returns the gob encoding of v.
func Marshal(v any) ([]byte, error) {
	var out bytes.Buffer
	err := gob.NewEncoder(&out).Encode(v)
	if err != nil {
		return []byte{}, err
	}

	return out.Bytes(), nil
}

// Unmarshal parses the gob encoding of data and stores the result in v.
func Unmarshal(data []byte, v any) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(v)
}
