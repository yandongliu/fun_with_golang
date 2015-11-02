/*
	http://matt.aimonetti.net/posts/2012/11/27/real-life-concurrency-in-go/
*/

package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var urls = []string{
	"http://www.google.com",
	"http://www.yahoo.com",
	"http://www.microsoft.com",
}

type HttpResponse struct {
	url      string
	response *http.Response
	err      error
	elapse   time.Duration
}

func asyncHttpGet(urls []string) []*HttpResponse {
	var wg sync.WaitGroup
	responses := []*HttpResponse{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fmt.Printf("Fetching %s \n", url)
			start := time.Now()
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			elapsed := time.Since(start)
			fmt.Printf("elapsed %s \n", elapsed)
			responses = append(responses, &HttpResponse{url, resp, err, elapsed})
		}(url)
	}

	wg.Wait()
	return responses
}

func main() {
	results := asyncHttpGet(urls)
	for _, result := range results {
		fmt.Printf("%s status: %s %s\n", result.url, result.response.Status, result.elapse)
	}
}
