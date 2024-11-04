package main

import "net/http"

func healthzfunc(w http.ResponseWriter, req *http.Request) {
	header := w.Header()
	body := []byte("OK")
	header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
