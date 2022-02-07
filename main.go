package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/asoorm/jsonplaceholder/model"
	"github.com/go-chi/chi"
)

var (
	userData    []model.User
	commentData []model.Comment
	postData    []model.Post
)

func loadFromDB(fileName string, m string) {
	myFile, err := os.Open(fileName)
	if err != nil {
		panic(err.Error())
	}
	defer myFile.Close()

	jsonParser := json.NewDecoder(myFile)

	switch m {
	case "users":
		if err = jsonParser.Decode(&userData); err != nil {
			panic(err.Error())
		}
	case "posts":
		if err = jsonParser.Decode(&postData); err != nil {
			panic(err.Error())
		}
	case "comments":
		if err = jsonParser.Decode(&commentData); err != nil {
			panic(err.Error())
		}
	default:
		panic("unknown m")
	}
}

func main() {
	loadFromDB("db/users.json", "users")
	loadFromDB("db/posts.json", "posts")
	loadFromDB("db/comments.json", "comments")

	muxer := chi.NewMux()

	muxer.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		jsBytes, _ := json.Marshal(userData)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsBytes)
	})

	muxer.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userIDString := chi.URLParam(r, "id")
		userID, _ := strconv.Atoi(userIDString)

		jsBytes, _ := json.Marshal(userData[userID-1])
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsBytes)
	})

	muxer.Get("/users/{id}/posts", func(w http.ResponseWriter, r *http.Request) {
		userIDString := chi.URLParam(r, "id")
		userID, _ := strconv.Atoi(userIDString)

		var userPosts []model.Post

		for _, post := range postData {
			if post.UserId == userID {
				userPosts = append(userPosts, post)
			}
		}

		jsBytes, _ := json.Marshal(userPosts)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsBytes)
	})

	muxer.Get("/posts/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		postIDString := chi.URLParam(r, "id")
		postID, _ := strconv.Atoi(postIDString)

		var postComments []model.Comment

		for _, comment := range commentData {
			if comment.PostId == postID {
				postComments = append(postComments, comment)
			}
		}

		jsBytes, _ := json.Marshal(postComments)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsBytes)
	})

	http.ListenAndServe(":8000", muxer)
}
