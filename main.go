package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
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

	reqResRecorder := func(writer io.Writer) func(next http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqBytes, _ := httputil.DumpRequest(r, true)
				writer.Write([]byte("===\n"))
				writer.Write(reqBytes)
				writer.Write([]byte("\n\n"))

				next.ServeHTTP(w, r)
			})
		}
	}

	muxer := chi.NewMux()
	muxer.Use(reqResRecorder(os.Stdout))

	muxer.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		jsBytes, _ := json.Marshal(userData)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsBytes)
	})

	muxer.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		u := model.User{}
		err := decoder.Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if u.Id == 0 {
			// we need a user id

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing user id"))
			return
		}

		for _, user := range userData {
			if u.Id == user.Id {

				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("user already exists"))
				return
			}
		}

		userData = append(userData, u)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(http.StatusText(http.StatusCreated)))
	})

	muxer.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userIDString := chi.URLParam(r, "id")
		userID, err := strconv.Atoi(userIDString)
		if err != nil {
			panic(err.Error())
		}

		for _, u := range userData {
			if u.Id == userID {

				jsBytes, _ := json.Marshal(u)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsBytes)

				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	})

	muxer.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userIDString := chi.URLParam(r, "id")
		userID, _ := strconv.Atoi(userIDString)

		for index, u := range userData {
			if u.Id == userID {
				newUsers := make([]model.User, 0)
				newUsers = append(newUsers, userData[:index]...)
				newUsers = append(newUsers, userData[index+1:]...)

				userData = newUsers

				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
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
