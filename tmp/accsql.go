package main

import (
  "database/sql"
  "fmt"
  _ "github.com/lib/pq"
//  "log"

)

type Post struct{
  Id      int
  Content string
  Author  string
}

var Db *sql.DB

// connect
func init(){
  var err error
  Db, err = sql.Open("postgres", "user=test dbname=test_db password=test sslmode=disable")
  if err != nil{
    panic(err)
  }
}

func (post *Post) Create()(err error){
  statement := "insert into posts (content, author) values ($1, $2) returning id"
  stmt, err := Db.Prepare(statement)
  if err != nil {
    return
  }
  defer stmt.Close()
  err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
  return
}

func GetPost(id int) (post Post, err error){
  post = Post{}
  err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
  return
}

// get all posts
func Posts(limit int) (posts []Post, err error) {
	rows, err := Db.Query("select id, content, author from posts limit $1", limit)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Content, &post.Author)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

func main(){

  post := Post{Content: "Hello jojo!", Author: "jojo"}
  fmt.Println(post)
  post.Create()
  fmt.Println(post)

  readPost, _ := GetPost(post.Id)
//  if err != nil{
//    log.Fatal("the post is not exist.")
//  }
  fmt.Println(readPost)

  post2 := Post{Content: "welcome", Author: "loa"}
  readPost, _ = GetPost(post2.Id)
//  if err != nil{
//    log.Fatal("the post2 is not exist.")
//  }
  fmt.Println(readPost)

}
