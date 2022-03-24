package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

var parallelLimit = flag.Int("parallel", 10, "number of parallel requests, should be > 0")

func getContentHash(requestURL string) (string, error) {
	url, err := url.Parse(requestURL)
	if err != nil {
		return "", fmt.Errorf("can't parse url: %s", err)
	}
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	resp, err := http.Get(url.String())
	if err != nil {
		return "", fmt.Errorf("can't get: %s", err)
	}
	defer resp.Body.Close()
	hash := md5.New()
	io.Copy(hash, resp.Body)

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func processURLs(urls []string, parallelLimit int, cb func(string)) {
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, parallelLimit)
	for _, url := range urls {
		url := url
		wg.Add(1)
		go func() {
			sem <- struct{}{}
			defer func() {
				wg.Done()
				<-sem
			}()

			hash, err := getContentHash(url)
			if err != nil {
				log.Fatalf("can't get content of %s: %s", url, err)
			}

			cb(fmt.Sprintf("%s %s", url, hash))
		}()
	}
	wg.Wait()
}

func main() {
	flag.Parse()
	if *parallelLimit <= 0 {
		flag.Usage()
		return
	}
	processURLs(flag.Args(), *parallelLimit, func(s string) {
		fmt.Println(s)
	})
}
