// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package idempotent

import (
	"net/http"
	"unicode/utf8"
)

const (
	// maxIdempotencyKeyLength is the longest X-Idempotency-Key string legnth allowed.
	maxIdempotencyKeyLength = 36
)

// Recorder offers a method to determine if a given key has been
// seen before or not. Each invocation of SeenBefore needs to
// record each key found, but there's no minimum duration required.
type Recorder interface {
	SeenBefore(key string) bool
}

// FromRequest extracts the idempotency key from HTTP headers and records its presence in
// the provided Recorder.
func FromRequest(req *http.Request, rec Recorder) (key string, seen bool) {
	key = truncate(req.Header.Get("X-Idempotency-Key"))
	if key == "" {
		return "", false
	}
	return key, rec.SeenBefore(key)
}

// SeenBefore sets a HTTP response code as an error for previously seen idempotency keys.
func SeenBefore(w http.ResponseWriter) {
	w.WriteHeader(http.StatusPreconditionFailed)
}

func truncate(s string) string {
	if utf8.RuneCountInString(s) > maxIdempotencyKeyLength {
		return s[:maxIdempotencyKeyLength]
	}
	return s
}
