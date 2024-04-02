package main

import (
	"net/http"
)

func main() {
	jh := NewJobsHandler()
	http.ListenAndServe(":8080", jh.mux)

    // TODO: write a goroutine to loop all values and delete the ones that are concluded after some time.
}

