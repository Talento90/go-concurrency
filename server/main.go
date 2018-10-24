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

		select {
		case <-time.After(sleep):
			ch <- true
			log.Printf("Job Finished")
			return
		case <-ctx.Done():
			log.Printf("Job was aborted")
		}
	}()

	return ch
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()
	defer func(start time.Time) {
		log.Print(time.Since(start))
	}(time.Now())

	select {
	case <-ctx.Done():
		log.Printf("Error: %v", ctx.Err())

		if ctx.Err() == context.Canceled {
			w.Write([]byte("Cancelled"))
		}

		if ctx.Err() == context.DeadlineExceeded {
			w.Write([]byte("Timeout"))
		}
	case <-processJob(ctx):
		w.Write([]byte("Job Done"))
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
