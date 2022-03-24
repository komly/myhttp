package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestGetHash(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}))
	defer s.Close()

	expected := "098f6bcd4621d373cade4e832627b4f6"
	hash, err := getContentHash(s.URL)
	if err != nil {
		t.Fatalf("getContentHash error: %s", err)
		return
	}
	if hash != expected {
		t.Fatalf("got: %s, want: %s", hash, expected)
	}
}

func TestProcessURLs(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}))
	defer s.Close()

	s2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test2"))
	}))
	defer s2.Close()

	expected := []string{
		s.URL + " 098f6bcd4621d373cade4e832627b4f6",
		s2.URL + " ad0234829205b9033196ba818f7a872b",
	}
	res := make([]string, 0)
	mu := sync.Mutex{}
	processURLs([]string{s.URL, s2.URL}, 10, func(s string) {
		mu.Lock()
		res = append(res, s)
		mu.Unlock()
	})
	sort.Sort(sort.StringSlice(res))
	sort.Sort(sort.StringSlice(expected))

	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("got: %+v, want: %+v", res, expected)
	}

}
