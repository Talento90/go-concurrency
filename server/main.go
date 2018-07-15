package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func processJob(ctx context.Context) chan bool {
	ch := make(chan bool)
	sleep := time.Duration(rand.Intn(8)) * time.Second

	go func() {
		defer close(ch)
		log.Printf("Job Start")

		select {
		case <-time.After(sleep):
			ch <- true
			log.Printf("Job Done")
		case <-ctx.Done():
			log.Printf("Job Cancelled")
		}
	}()

	return ch
}

// error: context canceled
// error: context deadline exceeded

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start request")
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)

	defer func() {
		log.Printf("start cancel")
		cancel()
		log.Printf("end cancel")
	}()

	select {
	case <-ctx.Done():
		log.Printf("Timeout request")

		if ctx.Err() != nil {
			log.Printf("error: %s", ctx.Err())
		}

		w.Write([]byte("Cancelled"))
		return

	case <-processJob(ctx):
		w.Write([]byte("HELLO WORLD"))
		log.Printf("End request")
	}

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
