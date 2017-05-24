package instrument

import "net/http"

// Trace has hooks that are called around the handler
// The main difference is that a ResponseWriterDelegate is provided
// that will track & provide the status code & amount written to the ResponseWriter.
type Trace struct {
	BeforeHandler func(http.ResponseWriter, *http.Request)
	AfterHandler  func(http.ResponseWriter, *http.Request, ResponseWriterDelegate)
}

// Handler should really be the outter most handler of your http.Handler chain.
// This way the ResponseWriterDelegate will accurately reflect the status & written amounts.
func Handler(trace *Trace, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if trace.BeforeHandler != nil {
			trace.BeforeHandler(w, r)
		}

		delegate := newResponseDelegator(w)
		next.ServeHTTP(delegate, r)

		if trace.AfterHandler != nil {
			trace.AfterHandler(w, r, delegate)
		}
	})
}
