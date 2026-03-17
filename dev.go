package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"portfolio/internal/builder"
)

const reloadScript = `<script>
(function() {
	var es = new EventSource('/live-reload');
	es.onmessage = function(e) {
		if (e.data === 'reload') { location.reload(); }
	};
	es.onerror = function() {
		setTimeout(function() { location.reload(); }, 1000);
	};
})();
</script>`

var (
	clientsMu sync.Mutex
	clients   = map[chan struct{}]struct{}{}
)

func runDevServer(cfg builder.Config) {
	// Initial build
	if err := builder.Build(cfg); err != nil {
		log.Printf("initial build error: %v", err)
	} else {
		log.Println("initial build done")
	}

	// File watcher
	go watchFiles(cfg)

	// Routes
	http.Handle("/live-reload", http.HandlerFunc(sseHandler))
	http.Handle("/", injectMiddleware(http.FileServer(http.Dir(cfg.OutputDir))))

	log.Println("dev server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan struct{}, 1)
	clientsMu.Lock()
	clients[ch] = struct{}{}
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, ch)
		clientsMu.Unlock()
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ch:
			fmt.Fprint(w, "data: reload\n\n")
			flusher.Flush()
		}
	}
}

func broadcast() {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for ch := range clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

type bufferedWriter struct {
	http.ResponseWriter
	buf    bytes.Buffer
	status int
}

func (bw *bufferedWriter) WriteHeader(status int) { bw.status = status }
func (bw *bufferedWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}

func injectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bw := &bufferedWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(bw, r)

		ct := w.Header().Get("Content-Type")
		if strings.Contains(ct, "text/html") {
			body := bw.buf.String()
			if idx := strings.LastIndex(body, "</body>"); idx != -1 {
				body = body[:idx] + reloadScript + body[idx:]
			}
			w.Header().Del("Content-Length")
			w.WriteHeader(bw.status)
			fmt.Fprint(w, body)
		} else {
			w.WriteHeader(bw.status)
			w.Write(bw.buf.Bytes())
		}
	})
}

func watchFiles(cfg builder.Config) {
	snapshots := map[string]time.Time{}

	// Seed initial snapshot
	_ = takeSnapshot([]string{cfg.ContentDir, "internal"}, snapshots)

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		newSnap := map[string]time.Time{}
		_ = takeSnapshot([]string{cfg.ContentDir, "internal"}, newSnap)

		if changed(snapshots, newSnap) {
			snapshots = newSnap
			if err := builder.Build(cfg); err != nil {
				log.Printf("rebuild error: %v", err)
			} else {
				log.Println("rebuilt")
				broadcast()
			}
		}
	}
}

func takeSnapshot(dirs []string, snap map[string]time.Time) error {
	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			snap[path] = info.ModTime()
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func changed(old, new map[string]time.Time) bool {
	if len(old) != len(new) {
		return true
	}
	for k, v := range new {
		if ov, ok := old[k]; !ok || !ov.Equal(v) {
			return true
		}
	}
	return false
}
