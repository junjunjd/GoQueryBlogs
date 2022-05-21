package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/junjunjd/GoQueryBlogs/pkg/db"
	"github.com/junjunjd/GoQueryBlogs/pkg/model"
)

type ErrorType uint8

const (
	TagsNotPresent ErrorType = iota
	SortByInvalid
	DirectionInvalid
)

// example url. Change this to the URL you want to fetch data from
const url string = "https://example.com/blog/posts?tag="

// database for server side caching
var (
	database = db.Database{
		Resp:  make(map[string]model.PostsResp),
		Mutex: &sync.RWMutex{},
	}
)

// query one tag at a time
func GetOneTag(tag string) model.PostsResp {
	resp, err := http.Get(url + tag)
	if err != nil {
		log.Fatal("http Get error when querying one tag: ", err)
	}

	var postsResp model.PostsResp
	err = json.NewDecoder(resp.Body).Decode(&postsResp)
	if err != nil {
		log.Fatal("json Decode error when decoding response: ", err)
	}

	return postsResp
}

// Cache response after querying one tag
func CacheTag(tag string) model.PostsResp {
	postsResp, found := database.Get(tag)
	if !found {
		postsResp = GetOneTag(tag)
		database.Set(tag, postsResp)
	}
	return postsResp
}

// query all tags concurrently
func GetTagsConcurrently(tags string) model.PostsResp {
	var postsResp model.PostsResp
	var mutex sync.Mutex
	wg := sync.WaitGroup{}
	for _, tag := range strings.Split(tags, ",") { // make concurrent query requests to the API
		wg.Add(1)
		go func(tag string) {
			mutex.Lock()
			postsResp.Posts = append(postsResp.Posts, CacheTag(tag).Posts...)
			mutex.Unlock()
			wg.Done()
		}(tag)
	}
	wg.Wait()
	return postsResp
}

func Query(tags string, sortByVal string, directionVal string) model.PostsResp {
	postsResp := GetTagsConcurrently(tags)
	postsResp.Deduplicate()
	postsResp.SortBy(sortByVal, directionVal)
	return postsResp
}

func ErrorResp(errType ErrorType) []byte {
	var resp map[string]string
	switch errType {
	case 0:
		resp = map[string]string{"error": "Tags parameter is required"}
	case 1:
		resp = map[string]string{"error": "sortBy parameter is invalid"}
	case 2:
		resp = map[string]string{"error": "direction parameter is invalid"}
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened when marshaling ErrorResp to JSON. Err: %s", err)
	}
	return jsonResp
}

func Ping(w http.ResponseWriter, r *http.Request) {
	resp := map[string]bool{"success": true}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened when marshaling Ping resp to JSON. Err: %s", err)
	}
	w.Write(jsonResp)
}

func Posts(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query().Get("tags") // parse all tags into a string
	if tags == "" {
		w.WriteHeader(400)
		w.Write(ErrorResp(TagsNotPresent))
		return
	}

	sortByVal := r.URL.Query().Get("sortBy") // parse sortBy parameter into a string
	if !(sortByVal == "" || sortByVal == "id" || sortByVal == "reads" || sortByVal == "likes" || sortByVal == "popularity") {
		w.WriteHeader(400)
		w.Write(ErrorResp(SortByInvalid))
		return
	}

	directionVal := r.URL.Query().Get("direction") // parse directionBy parameter into a string
	if !(directionVal == "" || directionVal == "asc" || directionVal == "desc") {
		w.WriteHeader(400)
		w.Write(ErrorResp(DirectionInvalid))
		return
	}

	resp := Query(tags, sortByVal, directionVal)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened when marshaling Query resp to JSON. Err: %s", err)
	}
	w.Write(jsonResp)
}

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/ping", Ping).Methods("GET")
	router.HandleFunc("/api/posts", Posts).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
