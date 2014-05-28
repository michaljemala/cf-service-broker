// Copyright 2014, The cf-service-broker Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that
// can be found in the LICENSE file.

package broker

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type PasswordGenerator interface {
	GeneratePassword() (string, error)
}

var RandomPasswordGenerator = &randomGenerator{16, base64.URLEncoding}

type randomGenerator struct {
	strength int
	encoding *base64.Encoding
}

func (g *randomGenerator) randomKey() []byte {
	k := make([]byte, g.strength)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}

func (g *randomGenerator) GeneratePassword() (string, error) {
	if key := g.randomKey(); key == nil {
		return "", errors.New("Failed to generate random key.")
	} else {
		return g.encoding.EncodeToString(key), nil
	}
}
