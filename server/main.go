package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func processJob(ctx context.Context) chan bool {
	ch := make(chan bool)

	go func() {
		log.Printf("Job Start")

		select {
		case <-time.After(time.Second * 5):
			ch <- true
			log.Printf("Job Done")
		case <-ctx.Done():
			log.Printf("Job Cancelled")
		}
	}()

	return ch
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start request")
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)

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
