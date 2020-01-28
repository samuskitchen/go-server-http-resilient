package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func slowAPICall(ctx context.Context) string {
	d := rand.Intn(5)
	select {
	case <-ctx.Done():
		log.Printf("slowAPICall was supposed to take %v seconds, but was canceled.", d)
		return ""
	case <-time.After(time.Duration(d) * time.Second):
		log.Printf("Slow API call done after %v seconds.\n", d)
		return "foobar"
	}
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	result := slowAPICall(r.Context())
	io.WriteString(w, result+"\n")
}

func main() {

	srv := http.Server{
		Addr:         ":8085",
		WriteTimeout: 5 * time.Second,
		Handler:      http.TimeoutHandler(http.HandlerFunc(slowHandler), 1*time.Second, "Timeout!"),
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
