// Package json stores dad's corpus in the file system.
package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Store satisfies dad's Storage interface.
type Store struct{}

// Data defines JSON encoding.
type Data struct {
	Grams map[string][]string
}

// Save to disk as JSON.
func (s *Store) Save(path string, grams map[string][]string) error {
	d := &Data{
		Grams: grams,
	}
	var buf = new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "    ")
	if err := enc.Encode(d); err != nil {
		return errors.Wrap(err, "json encoding failed")
	}
	if err := ioutil.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return errors.Wrap(err, "write fail failed")
	}
	return nil
}

// Load from disk.
func (s *Store) Load(path string) (map[string][]string, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var d *Data
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, err
	}
	return d.Grams, nil
}
