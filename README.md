# HTTP Server Instrumentation

Small utility to retrieve callbacks before and after a handler is called and includes 
the delegate wrapper from Prometheus that nicely provides the Status code & amount of bytes written.

This code is not tied to any frameworks, only using the go standard library.

## How to Use


```
    import "github.com/jdamick/go-http-instrument/instrument"

    tracer := &instrument.Trace{
		BeforeHandler: func(http.ResponseWriter, *http.Request) {
			... Collect metrics here ...
		},
		AfterHandler: func(http.ResponseWriter, *http.Request, ResponseWriterDelegate) {
			... Collect more metrics here ...
		},
	}

	h := instrument.Handler(tracer, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", h))
```




Basically this abstracts the Instrumented Handler available in prometheus:

<https://github.com/prometheus/client_golang/blob/master/prometheus/promhttp/instrument_server.go>



