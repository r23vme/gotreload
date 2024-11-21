// gotreload is a hot reload for web apps.
// Inject Script() into <body> on every page you wish to be hot-reloaded.
package gotreload

import (
	"fmt"
	"net/http"
)

const script string = `<script>
		let status;
		events = new EventSource("%s");
		events.onerror = () => {
			status = 0;
		}
		events.onopen = () => {
			if (status == 0) { location.reload(); }
		};
	</script>`

type gr struct {
	URL   string // URL of SSE endpoint.
	Retry string // Retry interval, ms
}

// New returns new instance of gr with default values:
// URL:   "/__gotreload",
// Retry: "500",
func New() gr {
	return gr{
		URL:   "/__gotreload",
		Retry: "500",
	}
}

// WithURL returns a copy of gr with new URL value.
func (s gr) WithURL(u string) gr {
	s.URL = u
	return s
}

// WithRetry returns a copy of gr with new Retry value.
func (s gr) WithRetry(r string) gr {
	s.Retry = r
	return s
}

// Script returns <script> which is supposed
// to be injected into <body> on every page to be gotreloaded.
func (s gr) Script() string {
	return fmt.Sprintf(script, s.URL)
}

// ServeHTTP provides SSE handler.
func (s gr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, _ := w.(http.Flusher)
	blocker := make(chan struct{})
	for {
		fmt.Fprintf(w, "retry: %s\n\n", s.Retry)
		flusher.Flush()
		<-blocker
	}
}
