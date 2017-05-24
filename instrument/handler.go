package instrument

import (
	"net/http"
	"time"
)

// Trace has hooks that are called around the handler
// The main difference is that a ResponseWriterDelegate is provided
// that will track & provide the status code & amount written to the ResponseWriter.
type Trace struct {
	BeforeHandler func(http.ResponseWriter, *http.Request, *TraceInfo)
	AfterHandler  func(http.ResponseWriter, *http.Request, *TraceInfo)
}

type TraceInfo struct {
	Timestamp         time.Time
	Method            string
	Status            int
	StatusCode        string
	RequestSizeBytes  int64
	ResponseSizeBytes int64
}

// Handler should really be the outter most handler of your http.Handler chain.
// This way the ResponseWriterDelegate will accurately reflect the status & written amounts.
func Handler(trace *Trace, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := &TraceInfo{
			Timestamp:        time.Now(),
			Method:           sanitizeMethod(r.Method),
			RequestSizeBytes: computeApproximateRequestSize(r),
		}

		if trace.BeforeHandler != nil {
			trace.BeforeHandler(w, r, info)
		}

		delegate := newResponseDelegator(w)
		next.ServeHTTP(delegate, r)

		if trace.AfterHandler != nil {
			info.Status = delegate.Status()
			info.StatusCode = sanitizeCode(delegate.Status())
			info.ResponseSizeBytes = delegate.Written()

			trace.AfterHandler(w, r, info)
		}
	})
}
