//
// based on ideas & work from Prometheus
//
// Copyright 2017 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package instrument

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

// Creates a writer delegator
func newResponseDelegator(w http.ResponseWriter) ResponseWriterDelegate {
	delegate := &responseWriterDelegator{ResponseWriter: w}

	// which interfaces it implements?
	_, cnOk := w.(http.CloseNotifier)
	_, flOk := w.(http.Flusher)
	_, hjOk := w.(http.Hijacker)
	_, psOk := w.(http.Pusher)
	_, rfOk := w.(io.ReaderFrom)

	switch {
	case cnOk && flOk && hjOk && rfOk && psOk:
		return &fancyPushDelegator{
			fancyDelegator: &fancyDelegator{delegate},
			push:           &pushDelegator{delegate},
		}
	case cnOk && flOk && hjOk && rfOk: // no http.Pusher
		return &fancyDelegator{delegate}
	case psOk: // only http.Pusher
		return &pushDelegator{delegate}
	}

	return delegate
}

type fancyPushDelegator struct {
	push *pushDelegator

	*fancyDelegator
}

func (f *fancyPushDelegator) Push(target string, opts *http.PushOptions) error {
	return f.push.Push(target, opts)
}

type pushDelegator struct {
	*responseWriterDelegator
}

func (f *pushDelegator) Push(target string, opts *http.PushOptions) error {
	return f.ResponseWriter.(http.Pusher).Push(target, opts)
}

type fancyDelegator struct {
	*responseWriterDelegator
}

func (r *fancyDelegator) CloseNotify() <-chan bool {
	return r.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (r *fancyDelegator) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

func (r *fancyDelegator) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

func (r *fancyDelegator) ReadFrom(re io.Reader) (int64, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.(io.ReaderFrom).ReadFrom(re)
	r.written += n
	return n, err
}

type ResponseWriterDelegate interface {
	http.ResponseWriter

	Status() int
	Written() int64
}

type responseWriterDelegator struct {
	http.ResponseWriter

	status      int
	written     int64
	wroteHeader bool
}

func (r *responseWriterDelegator) Status() int {
	return r.status
}

func (r *responseWriterDelegator) Written() int64 {
	return r.written
}

func (r *responseWriterDelegator) WriteHeader(code int) {
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}
