package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"urllookupservice/common"
	"urllookupservice/datastore"

	"github.com/gorilla/mux"
	cache "github.com/hashicorp/golang-lru"
)

type Server struct {
	Addr   string
	// In memory data store
	urlMap           sync.Map
	lruInMemoryCache *cache.ARCCache
	Router           *mux.Router
	cb               *datastore.Couchbase
}

func NewServer(addr string) *Server {
	s :=  &Server{
		Addr:addr,
		Router: mux.NewRouter(),
	}
	c, err := cache.NewARC(common.MaxCacheSize)
	handleError(err)
	s.lruInMemoryCache = c
	s.cb, _ = datastore.New()
	return s
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func (s *Server) StartServer() {
	fmt.Printf("Starting listening at port %s\n", s.Addr)
	log.Fatal(http.ListenAndServe(s.Addr, s.Router))
}

func (s *Server) checkInCache(key string) bool {
	fmt.Printf("Checking cache for key :%s from cache\n", key)
	v, ok := s.lruInMemoryCache.Get(key)
	if ok {
		return v.(bool)
	}
	return ok
}

func (s *Server) checkInDataStore(key string) bool {
	if s.cb.Offline {
		fmt.Println("Couchbase server is not available")
		return false
	}
	fmt.Printf("Checking datastore for key :%s from cache\n", key)
	return s.cb.Get(key)
}

func (s *Server) IsMalwareURLv2(w http.ResponseWriter, r *http.Request)  {
	urlParam := mux.Vars(r)["url"]
	domain, err := common.ParseDomainName(urlParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Url parsing error: %s", err.Error())))
		return
	}
	ok := s.checkInCache(domain) || s.checkInDataStore(domain)
	s.lruInMemoryCache.Add(domain, ok)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatBool(ok)))
}

func (s *Server) AddMalwareURL(w http.ResponseWriter, r *http.Request)  {
	urlParam := mux.Vars(r)["url"]
	domain, err := common.ParseDomainName(urlParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Url parsing error: %s", err.Error())))
		return
	}
	s.lruInMemoryCache.Add(domain, true)
	if !s.cb.Offline {
		s.cb.PutOrPost(domain)
	} else {
		fmt.Println("Couchbase server is not available")
	}
}

func (s *Server) RemoveMalwareURL(w http.ResponseWriter, r *http.Request) {
	urlParam := mux.Vars(r)["url"]
	domain, err := common.ParseDomainName(urlParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Url parsing error: %s", err.Error())))
		return
	}
	s.lruInMemoryCache.Remove(domain)
	if !s.cb.Offline {
		s.cb.Delete(domain)
	} else {
		fmt.Println("Couchbase server is not available")
	}
}

func (s *Server) IsMalwareURL(w http.ResponseWriter, r *http.Request) {
	urlParam := mux.Vars(r)["url"]
	fmt.Printf("Checking if url:%s belongs to malware URLs\n", urlParam)
	domain, err := common.ParseDomainName(urlParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Url parsing error: %s", err.Error())))
		return
	}
	_, ok := s.urlMap.Load(domain)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatBool(ok)))
}

func (s *Server) AddUrlToMalwareList(w http.ResponseWriter, r *http.Request) {
	urlParam := mux.Vars(r)["url"]
	fmt.Printf("Received request to add URL:%s to malware url list\n", urlParam)
	domain, err := common.ParseDomainName(urlParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Url parsing error: %s", err.Error())))
		return
	}
	s.urlMap.Store(domain, true)
	w.Write([]byte(fmt.Sprintf("Successfully added domain:%s to malware list", domain)))
	w.WriteHeader(http.StatusOK)
}

func (s *Server) RemoveUrlFromMalwareList(w http.ResponseWriter, r *http.Request) {
	urlParam := mux.Vars(r)["url"]
	fmt.Printf("Received request to delete URL:%s to malware url list\n", urlParam)

	domain, err := common.ParseDomainName(urlParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Url parsing error: %s", err.Error())))
		return
	}
	s.urlMap.Delete(domain)
	w.Write([]byte(fmt.Sprintf("Successfully deleted domain:%s to malware list", domain)))
	w.WriteHeader(http.StatusOK)
}
