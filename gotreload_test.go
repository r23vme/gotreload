package gotreload_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/r23vme/gotreload"
)

func Example() {
	mux := http.NewServeMux()

	gr := gotreload.New()

	// Change default SSE-endpoint url.
	// gr = gotreload.WithURL("/__gotreload")

	// Change default SSE retry interval, ms.
	// gr = gotreload.WithURL("500")

	// Attach SSE handler.
	mux.Handle("GET "+gr.URL, gr)

	// Page to be reloaded.
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Cache-Control", "no-cache")

		fmt.Fprint(w, `<html><head><meta name="color-scheme" content="light dark"></head><body>`)
		fmt.Fprintf(w, "<span style='color: gray'>page got reloaded at:</span> %s\n", time.Now().Format(time.StampMilli))

		// Write Script to the end of the body.
		w.Write([]byte(gr.Script()))

		fmt.Fprint(w, "</body></html>")
	})

	http.ListenAndServe("localhost:8080", mux)

	// Output:
	// Open localhost:8080 in browser, make some changes,
	// restart go app, browser will reload page automatically.
}
