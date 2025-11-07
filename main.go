package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	server := &http.Server{}
	server.Handler = mux
	server.Addr = ":8080"

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("./"))))
	mux.HandleFunc("/healthz", healthzHandler)

	server.ListenAndServe()
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
