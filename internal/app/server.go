package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const addr string = "localhost:8080"

type Counter struct {
	store     map[int]url.URL
	storeLock *sync.Mutex
}

func Server() *http.Server {
	counter := Counter{
		store:     make(map[int]url.URL, 0),
		storeLock: &sync.Mutex{},
	}
	return &http.Server{
		Addr:    addr,
		Handler: &counter,
	}
}

func (s *Counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if isPost(r) {
		s.getUrlByPost(w, r)
		return
	}

	if isGet(r) {
		s.getUrlByGet(w, r)
		return
	}

	http.NotFound(w, r)
}

func isPost(r *http.Request) bool {
	return r.URL.Path == "/" && r.Method == http.MethodPost
}

func isGet(r *http.Request) bool {
	pathParts := strings.Split(r.URL.Path, "/")
	return len(pathParts) == 2
}

func (s *Counter) getUrlByPost(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, errRead := buf.ReadFrom(r.Body)
	if errRead != nil {
		http.Error(w, "Cannot read request", http.StatusBadRequest)
		return
	}
	url, _ := url.Parse(buf.String())
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, fmt.Sprintf("http://%v/%d", addr, s.persistUrl(*url)))
}

func (s *Counter) getUrlByGet(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	url := s.retrieveUrl(id)
	if url == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Add("Location", url.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Counter) persistUrl(url url.URL) int {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	id := len(s.store)
	s.store[id] = url

	return id
}

func (s *Counter) retrieveUrl(id int) *url.URL {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	url, found := s.store[id]
	if !found {
		return nil
	}

	return &url
}
