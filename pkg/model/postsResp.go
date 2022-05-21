package model

import "sort"

type PostJSON struct {
	Id         uint64   `json:"id"`
	Author     string   `json:"author"`
	AuthorId   uint64   `json:"authorId"`
	Likes      uint64   `json:"likes"`
	Popularity float64  `json:"popularity"`
	Reads      uint64   `json:"reads"`
	Tags       []string `json:"tags"`
}

type PostsResp struct {
	Posts []PostJSON `json:"posts"`
}

// Remove all the repeated posts
func (p *PostsResp) Deduplicate() {
	var posts []PostJSON
	postsMap := make(map[uint64]struct{}) // create a map with 'id' as key
	empty := struct{}{}                   // a shortcut variable used as value

	for _, post := range p.Posts {
		if _, found := postsMap[post.Id]; !found { // check if post already exists
			postsMap[post.Id] = empty
			posts = append(posts, post)
		}
	}

	p.Posts = posts // update field Posts after removing repeated posts
}

// By default, posts are sorted by 'id' in ascending order
func (p *PostsResp) SortByDefault() {
	sort.Slice(p.Posts, func(i, j int) bool { return p.Posts[i].Id < p.Posts[j].Id })
}

// Use SliceStable to sort. When the sortBy elements are equal, posts are ordered by 'id' in ascending order
// Posts are first sorted using SortByDefault()
func (p *PostsResp) SortBy(sortByVal string, directionVal string) {
	if directionVal == "desc" {
		if sortByVal == "" || sortByVal == "id" {
			sort.Slice(p.Posts, func(i, j int) bool { return p.Posts[i].Id > p.Posts[j].Id })
		} else {
			p.SortByDefault() // Posts are first sorted by 'id' in ascending order. Thus when the sortBy elements are equal, posts are ordered by 'id'
			if sortByVal == "reads" {
				sort.SliceStable(p.Posts, func(i, j int) bool { return p.Posts[i].Reads > p.Posts[j].Reads })
			} else if sortByVal == "likes" {
				sort.SliceStable(p.Posts, func(i, j int) bool { return p.Posts[i].Likes > p.Posts[j].Likes })
			} else if sortByVal == "popularity" {
				sort.SliceStable(p.Posts, func(i, j int) bool { return p.Posts[i].Popularity > p.Posts[j].Popularity })
			}
		}
	} else if directionVal == "" || directionVal == "asc" {
		p.SortByDefault() // Posts are first sorted by 'id' in ascending order. Thus when the sortBy elements are equal, posts are ordered by 'id'
		if sortByVal == "reads" {
			sort.SliceStable(p.Posts, func(i, j int) bool { return p.Posts[i].Reads < p.Posts[j].Reads })
		} else if sortByVal == "likes" {
			sort.SliceStable(p.Posts, func(i, j int) bool { return p.Posts[i].Likes < p.Posts[j].Likes })
		} else if sortByVal == "popularity" {
			sort.SliceStable(p.Posts, func(i, j int) bool { return p.Posts[i].Popularity < p.Posts[j].Popularity })
		}
	}
}
