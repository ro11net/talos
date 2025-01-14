// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package configloader provides methods to load Talos config.
package configloader

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/siderolabs/talos/pkg/machinery/config"
	"github.com/siderolabs/talos/pkg/machinery/config/configloader/internal/decoder"
	"github.com/siderolabs/talos/pkg/machinery/config/container"
	_ "github.com/siderolabs/talos/pkg/machinery/config/types" //nolint:revive
)

// ErrNoConfig is returned when no configuration was found in the input.
var ErrNoConfig = errors.New("config not found")

// newConfig initializes and returns a Configurator.
func newConfig(r io.Reader) (config config.Provider, err error) {
	dec := decoder.NewDecoder()

	var buf bytes.Buffer

	// preserve the original contents
	r = io.TeeReader(r, &buf)

	manifests, err := dec.Decode(r)
	if err != nil {
		return nil, err
	}

	if len(manifests) == 0 {
		return nil, ErrNoConfig
	}

	return container.NewReadonly(buf.Bytes(), manifests...)
}

// NewFromFile will take a filepath and attempt to parse a config file from it.
func NewFromFile(filepath string) (config.Provider, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer f.Close() //nolint:errcheck

	return newConfig(f)
}

// NewFromStdin initializes a config provider by reading from stdin.
func NewFromStdin() (config.Provider, error) {
	return newConfig(os.Stdin)
}

// NewFromBytes will take a byteslice and attempt to parse a config file from it.
func NewFromBytes(source []byte) (config.Provider, error) {
	return newConfig(bytes.NewReader(source))
}
