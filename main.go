package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

const port = "8080"
const rootfilepath = "."

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}
func (cfg *apiConfig) showMetrics(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Swap(0)
	w.Write([]byte(fmt.Sprintf("Hits reset to %v", cfg.fileserverHits.Load())))
}

func handlerHealthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	cfg := apiConfig{}
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(rootfilepath)))))
	serveMux.HandleFunc("/healthz", handlerHealthz)
	serveMux.HandleFunc("/metrics", cfg.showMetrics)
	serveMux.HandleFunc("/reset", cfg.resetMetrics)
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
