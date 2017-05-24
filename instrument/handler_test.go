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
		BeforeHandler: func(http.ResponseWriter, *http.Request) {
			beforeCalled = true
		},
		AfterHandler: func(http.ResponseWriter, *http.Request, ResponseWriterDelegate) {
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
		BeforeHandler: func(http.ResponseWriter, *http.Request) {
			fmt.Printf("Before\n")
		},
		AfterHandler: func(http.ResponseWriter, *http.Request, ResponseWriterDelegate) {
			fmt.Printf("After\n")
		},
	}

	h := Handler(tracer, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
