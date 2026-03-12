package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var port = flag.Int("port", 8080, "port to listen on")

func main() {
	flag.Parse()

	http.HandleFunc("/", proxyHandler)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("CORS proxy listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("url")
	if target == "" {
		http.Error(w, "missing ?url= parameter", http.StatusBadRequest)
		return
	}

	parsed, err := url.ParseRequestURI(target)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	// Forward request headers, skipping hop-by-hop headers
	for key, vals := range r.Header {
		switch key {
		case "Host", "Connection", "Te", "Trailers", "Transfer-Encoding", "Upgrade":
			continue
		}
		for _, v := range vals {
			req.Header.Add(key, v)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "upstream request failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// Handle preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
