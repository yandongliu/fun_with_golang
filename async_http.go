/*
	http://matt.aimonetti.net/posts/2012/11/27/real-life-concurrency-in-go/
*/

package main

import (
	"fmt"
	"net/http"
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
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}
	for _, url := range urls {
		go func(url string) {
			fmt.Printf("Fetching %s \n", url)
			start := time.Now()
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			elapsed := time.Since(start)
			fmt.Printf("elapsed %s \n", elapsed)
			ch <- &HttpResponse{url, resp, err, elapsed}
		}(url)
	}

	for {
		select {
		case r := <-ch:
			fmt.Printf("%s was fetched\n", r.url)
			if r.err != nil {
				fmt.Printf("has error: %s\n", r.err)
			}
			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
	return responses
}

func main() {
	results := asyncHttpGet(urls)
	for _, result := range results {
		fmt.Printf("%s status: %s %s\n", result.url, result.response.Status, result.elapse)
	}
}
