package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/junjunjd/GoQueryBlogs/pkg/handler"
	"github.com/junjunjd/GoQueryBlogs/pkg/model"
)

var (
	post1 = model.PostJSON{
		Id:         1,
		Author:     "Rylee Paul",
		AuthorId:   9,
		Likes:      960,
		Popularity: 0.13,
		Reads:      50361,
		Tags:       []string{"tech", "health"},
	}
	post2 = model.PostJSON{
		Id:         2,
		Author:     "Zackery Turner",
		AuthorId:   12,
		Likes:      469,
		Popularity: 0.68,
		Reads:      90406,
		Tags:       []string{"startups", "tech", "history"},
	}
	post4 = model.PostJSON{
		Id:         4,
		Author:     "Elisha Friedman",
		AuthorId:   8,
		Likes:      728,
		Popularity: 0.88,
		Reads:      19645,
		Tags:       []string{"science", "design", "tech"},
	}
	post12 = model.PostJSON{
		Id:         12,
		Author:     "Adalyn Blevins",
		AuthorId:   11,
		Likes:      590,
		Popularity: 0.32,
		Reads:      80351, Tags: []string{"tech"},
	}
)

// test get one tag
func TestGetOneTag(t *testing.T) {
	tag := "tech"
	postsResp := handler.GetTagsConcurrently(tag)
	postsRespExp := handler.GetOneTag(tag)
	if !reflect.DeepEqual(postsResp, postsRespExp) {
		t.Errorf("``````````````````````````````TestGetOneTag result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, postsRespExp)
	}
}

func TestDeduplicate(t *testing.T) {
	postsResp := model.PostsResp{
		Posts: []model.PostJSON{post1, post2, post1, post2, post4, post12, post4, post12},
	}
	postsRespExp := model.PostsResp{
		Posts: []model.PostJSON{post1, post2, post4, post12},
	}
	postsResp.Deduplicate()
	if !reflect.DeepEqual(postsResp, postsRespExp) {
		t.Errorf("``````````````````````````````TestCombinePosts result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, postsRespExp)
	}
}

func TestSortBy(t *testing.T) {
	postsResp := model.PostsResp{
		Posts: []model.PostJSON{post4, post12, post1, post2},
	}
	sortByDefaultExp := model.PostsResp{
		Posts: []model.PostJSON{post1, post2, post4, post12},
	}

	postsResp.SortBy("", "")
	if !reflect.DeepEqual(postsResp, sortByDefaultExp) {
		t.Errorf("``````````````````````````````TestsortByDefault result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, sortByDefaultExp)
	}

	postsResp.SortBy("id", "asc")
	if !reflect.DeepEqual(postsResp, sortByDefaultExp) {
		t.Errorf("``````````````````````````````TestsortById result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, sortByDefaultExp)
	}

	postsResp.SortBy("reads", "")
	SortByReadsExp := model.PostsResp{
		Posts: []model.PostJSON{post4, post1, post12, post2},
	}
	if !reflect.DeepEqual(postsResp, SortByReadsExp) {
		t.Errorf("``````````````````````````````TestSortByReads result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, SortByReadsExp)
	}

	postsResp.SortBy("likes", "")
	SortByLikesExp := model.PostsResp{
		Posts: []model.PostJSON{post2, post12, post4, post1},
	}
	if !reflect.DeepEqual(postsResp, SortByLikesExp) {
		t.Errorf("``````````````````````````````TestSortByLikes result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, SortByLikesExp)
	}

	postsResp.SortBy("popularity", "")
	SortByPopExp := model.PostsResp{
		Posts: []model.PostJSON{post1, post12, post2, post4},
	}
	if !reflect.DeepEqual(postsResp, SortByPopExp) {
		t.Errorf("``````````````````````````````TestSortByPop result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, SortByPopExp)
	}

	postsResp.SortBy("reads", "desc")
	SortByReadsExp = model.PostsResp{
		Posts: []model.PostJSON{post2, post12, post1, post4},
	}
	if !reflect.DeepEqual(postsResp, SortByReadsExp) {
		t.Errorf("``````````````````````````````TestSortByReads result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, SortByReadsExp)
	}

	postsResp.SortBy("likes", "desc")
	SortByLikesExp = model.PostsResp{
		Posts: []model.PostJSON{post1, post4, post12, post2},
	}
	if !reflect.DeepEqual(postsResp, SortByLikesExp) {
		t.Errorf("``````````````````````````````TestSortByLikes result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, SortByLikesExp)
	}

	postsResp.SortBy("popularity", "desc")
	SortByPopExp = model.PostsResp{
		Posts: []model.PostJSON{post4, post2, post12, post1},
	}
	if !reflect.DeepEqual(postsResp, SortByPopExp) {
		t.Errorf("``````````````````````````````TestSortByPop result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, SortByPopExp)
	}
}

// test get multiple tags concurrently
func TestGetTagsConcurrently(t *testing.T) {
	tags := "tech,science,design"
	// sequentially GetOneTag at a time
	var postsRespExp model.PostsResp
	for _, tag := range strings.Split(tags, ",") {
		postsRespExp.Posts = append(postsRespExp.Posts, handler.GetOneTag(tag).Posts...)
	}
	postsRespExp.SortBy("", "")

	// GetTagsConcurrently
	postsResp := handler.GetTagsConcurrently(tags)
	postsResp.SortBy("", "")

	if !reflect.DeepEqual(postsResp, postsRespExp) {
		t.Errorf("``````````````````````````````TestCombinePosts result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, postsRespExp)
	}
}

// integration test
func TestQuery(t *testing.T) {
	tags := "tech,science"
	// GetTagsConcurrently
	postsResp := handler.Query(tags, "likes", "desc")

	// sequentially GetOneTag at a time
	var postsRespExp model.PostsResp
	for _, tag := range strings.Split(tags, ",") {
		postsRespExp.Posts = append(postsRespExp.Posts, handler.GetOneTag(tag).Posts...)
	}
	postsRespExp.Deduplicate()
	postsRespExp.SortBy("likes", "desc")

	if !reflect.DeepEqual(postsResp, postsRespExp) {
		t.Errorf("``````````````````````````````TestQuery result incorrect, ``````````````````````````````got: %v, ``````````````````````````````want: %v.", postsResp, postsRespExp)
	}
}
