package main

import (
	"fmt"
	"sort"
	"time"
)

/*
Build a Facebook-like news feed system ranking posts by relevance and time.

user should be able to get news posts by relevance (category).
user should be able to get news post by created_at.

user hits the newsFeed(with filter relevance or time or both) -> FBApp NewsFeed() -> PostService -> GetPosts(...filters) -> PostRepo apply filters -> return []posts

Post
id
category
message
created_at

FBApp
postService

IFilter

SortFilter
asc
field

ValueFilter
field
value

*/
type FBApp struct {
	postService PostService 
}


func (fb *FBApp) AddPost(post Post) {
	fb.postService.AddPost(post)
}

func (fb *FBApp) GetPost(postId int64) Post  {
	return fb.postService.GetPost(postId)
}

func (fb *FBApp) NewsFeed(filters ...Ifilter) []Post {
	return fb.postService.NewsFeed(filters...)
}


type PostService struct {
	repo PostRepo
}

func (ps *PostService) AddPost(post Post) {
	ps.repo.AddPost(post)
}

func (ps *PostService) GetPost(postId int64) Post  {
	return ps.repo.GetPost(postId)
}

func (ps *PostService) NewsFeed(filters ...Ifilter) []Post {
	return ps.repo.GetNewsFeed(filters...)
}


type PostRepo struct {
	db IDB
}

func (p *PostRepo) AddPost(post Post) {
	p.db.AddPost(post)
}

func (p *PostRepo) GetPost(postId int64) Post {
	return p.db.GetPost(postId)
}

func(p *PostRepo) GetNewsFeed(filters ...Ifilter) []Post {
	return p.db.Filter(filters...)
}

type IDB interface{
	Filter(...Ifilter) []Post 
	AddPost(post Post)
	GetPost(postId int64) Post 
}

type InMemoryDB struct {
	posts map[int64]Post
}

func (md *InMemoryDB) AddPost(post Post) {
	md.posts[post.Id] = post 
}

func (md *InMemoryDB) GetPost(postId int64) Post {
	post, ok := md.posts[postId]
	if ok {
		return post 
	}
	return Post{}
} 

func (md *InMemoryDB) Filter(filters ...Ifilter) []Post {
	posts := md.posts
	arr := make([]Post, 0)
	for _, v := range posts {
		arr = append(arr, v)
	}

	for _,v := range filters {
		arr = v.Apply(arr)
	}

	return arr
}

type Post struct {
	Id int64 
	message string
	category Category 
	createdAt int64
}

type Category int64

const (
	SPORTS Category = iota
	POLITICS 
	MOVIES
)

type FieldGetter func(Post) any 

var postFieldGetters = map[string]FieldGetter{
	"id": func(p Post) any {return  p.Id},
	"message" :func(p Post) any {return p.message},
	"category": func(p Post) any {return p.category},
	"createdAt": func(p Post) any {return p.createdAt},
}

type Ifilter interface{
	Apply(posts []Post) []Post
}

type SortFilter struct {
	asc bool
	field string 
}

func(sf *SortFilter) Apply(arr []Post) []Post {
	getter, ok := postFieldGetters[sf.field]
	if !ok {
		fmt.Println("Incorrect Field ", sf.field)
		return make([]Post,0)
	}

	sort.Slice(arr, func(i, j int) bool {
		ai := getter(arr[i])
		aj := getter(arr[j])
		if(sf.asc) {
			return fmt.Sprint(ai) < fmt.Sprint(aj) 
		}

		return fmt.Sprint(ai) > fmt.Sprint(aj)
	})

	return arr 
}

type ValueFilter struct {
	field string 
	value string 
}

func(vf *ValueFilter) Apply(arr []Post) []Post {
	getter, ok := postFieldGetters[vf.field]
	if !ok {
		fmt.Println("incorrect field ", vf.field)
	}
 
	out := make([]Post , 0)
	for _,v := range arr {
		if fmt.Sprint(getter(v)) == vf.value {
			out = append(out, v)
		}
	}

	return out
 }

 func NewFbApp() FBApp {
	prepo := PostRepo{
		db: &InMemoryDB{
			posts: make(map[int64]Post),
		},
	}
	ps := PostService{
	repo: prepo,
	}
	return FBApp{
		postService: ps ,
	}

 }

func main() {
	fb := NewFbApp()
	
	post1 := Post{
		Id: 1,
		message: "1st post",
		category: POLITICS,
		createdAt: time.Now().Unix(),
	}
		post2 := Post{
		Id: 2,
		message: "2st post",
		category: MOVIES,
		createdAt: time.Now().Unix()+1,
	}
		post3 := Post{
		Id: 3,
		message: "3st post",
		category: POLITICS,
		createdAt: time.Now().Unix()+2,
	}
		post4 := Post{
		Id: 4,
		message: "4st post",
		category: SPORTS,
		createdAt: time.Now().Unix()+3,
	}

	fb.AddPost(post1)
		fb.AddPost(post2)
			fb.AddPost(post3)
				fb.AddPost(post4)
	
	filters := []Ifilter{
		// &ValueFilter{
		// 	field: "category",
		// 	value: fmt.Sprint(POLITICS),
		// },
		&SortFilter{
			field: "createdAt",
			asc: true,
		},
	}
	fmt.Println(fb.NewsFeed(filters...))
}