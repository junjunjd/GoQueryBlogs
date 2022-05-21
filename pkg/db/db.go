package db

import (
	"sync"

	"github.com/junjunjd/GoQueryBlogs/pkg/model"
)

type Database struct {
	Resp  map[string]model.PostsResp
	Mutex *sync.RWMutex
}

// The assessment does not mention anything about duration. We assume that cached response does not expire.
func (db Database) Get(tag string) (model.PostsResp, bool) {
	db.Mutex.Lock()
	postsResp, found := db.Resp[tag]
	db.Mutex.Unlock()
	return postsResp, found
}

// The assessment does not mention anything about duration. We assume that cached response does not expire.
func (db *Database) Set(tag string, postsResp model.PostsResp) {
	db.Mutex.Lock()
	db.Resp[tag] = postsResp
	db.Mutex.Unlock()
}
