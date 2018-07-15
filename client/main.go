package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const requestNumber = 1

func request(wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)

	if err != nil {
		log.Printf("%v", err)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient

	res, err := client.Do(req)

	if err != nil {
		log.Printf("%v", err)
		return
	}

	fmt.Printf("%v\n", res.StatusCode)
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(requestNumber)

	for i := 0; i < requestNumber; i++ {
		go request(&wg)
	}

	wg.Wait()
}
