package main

import (
	"fmt"
	"sync"
)

type SafeMapper struct {
	mu sync.Mutex
	v  map[string]bool
}

var safeMapper = SafeMapper{v: make(map[string]bool)}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func (m *SafeMapper) safeCheck(url string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.v[url]
	return ok
}

func (m *SafeMapper) safeAddKey(url string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.v[url] = true
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, status chan bool) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:

	if depth <= 0 {
		status <- false
		return
	}

	// check url appearance in safeMapper. If exists, url has already been crawled. If not, crawl and update safeMapper.
	if ok := safeMapper.safeCheck(url); ok {
		status <- true
		return
	}
	
	//crawl url
	body, urls, err := fetcher.Fetch(url)
	//add url to safeMapper no matter found or not found
	safeMapper.safeAddKey(url)
	if err != nil {
		fmt.Println(err)
		status <- false
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	
	//operate urls, further crawl
	statuses := make([]chan bool, len(urls))
	for i, u := range urls {
		statuses[i] = make(chan bool)
		go Crawl(u, depth-1, fetcher, statuses[i])
	}
	for _, childStatus := range statuses {
		<-childStatus
	}
	status <- true
	return
}

func main() {
	status := make(chan bool)
	go Crawl("https://golang.org/", 4, fetcher, status)
	<-status
	// important to wait for the channel status
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

