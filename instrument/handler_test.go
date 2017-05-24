package instrument

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerEvents(t *testing.T) {
	beforeCalled := false
	afterCalled := false

	tracer := &Trace{
		BeforeHandler: func(resp http.ResponseWriter, req *http.Request, info *TraceInfo) {
			require.NotNil(t, info.Timestamp)
			require.Equal(t, "get", info.Method)
			require.True(t, info.RequestSizeBytes > 0)
			require.Equal(t, int64(0), info.ResponseSizeBytes)
			beforeCalled = true
		},
		AfterHandler: func(resp http.ResponseWriter, req *http.Request, info *TraceInfo) {
			require.NotNil(t, info.Timestamp)
			require.Equal(t, "get", info.Method)
			require.Equal(t, "200", info.StatusCode)
			require.Equal(t, 200, info.Status)
			require.True(t, info.RequestSizeBytes > 0)
			require.True(t, info.ResponseSizeBytes > 0)

			afterCalled = true
		},
	}

	h := Handler(tracer, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	srv := httptest.NewServer(h)
	defer srv.Close()

	_, err := http.Get(srv.URL)
	require.Nil(t, err)
	require.True(t, beforeCalled)
	require.True(t, afterCalled)
}

func ExampleHandler() {
	tracer := &Trace{
		BeforeHandler: func(http.ResponseWriter, *http.Request, *TraceInfo) {
			fmt.Printf("Before\n")
		},
		AfterHandler: func(http.ResponseWriter, *http.Request, *TraceInfo) {
			fmt.Printf("After\n")
		},
	}

	h := Handler(tracer, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
