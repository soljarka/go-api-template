package user

import (
	"context"
	"ctco-dev/go-api-template/internal/log"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// Specification is a struct with config parameters for User service
type Specification struct {
	UseMongoDB             bool   `split_words:"true" required:"true" default:"true"`
	MongoDBUrl             string `split_words:"true" required:"true" default:"mongodb://localhost:27017"`
	MongoDBDatabase        string `split_words:"true" required:"true" default:"db"`
	MongoDBCollectionUsers string `split_words:"true" required:"true" default:"users"`
	UsersRoute             string `split_words:"true" required:"true" default:"/users"`
}

// LaunchController registers a new REST endpoint for User operations
func LaunchController(ctx context.Context, pattern string, service Service) {

	http.HandleFunc(
		pattern,
		func(w http.ResponseWriter, r *http.Request) {
			reqID := uuid.NewV4().String()[0:8]
			reqCtx := log.NewContext(ctx, logrus.Fields{"reqID": reqID})
			reqCtx, cancel := context.WithTimeout(reqCtx, time.Second*10)
			defer cancel()

			if r.Method == "GET" {
				id := r.URL.Query().Get("id")
				var users []*User
				var err error
				if id == "" {
					users, err = service.All(reqCtx)
					if err != nil {
						handleError(reqCtx, w, err, "Failed to load users")
						return
					}
				} else {
					usr, err := service.Find(reqCtx, id)
					if err != nil {
						handleError(reqCtx, w, err, "Failed to load user by id")
						return
					}
					users = append(users, usr)
				}

				bytes, err := json.Marshal(users)
				if err != nil {
					handleError(reqCtx, w, err, "Failed to parse users")
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Write(bytes)

			} else if r.Method == "POST" {
				usr, err := readUser(r)
				if err != nil {
					handleError(reqCtx, w, err, "Unable to parse body")
					return
				}

				id, err := service.Save(reqCtx, usr)
				if err != nil {
					handleError(reqCtx, w, err, "Failed to save user")
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(id))

			} else if r.Method == "DELETE" {
				id := r.URL.Query().Get("id")

				if id == "" {
					err := service.DeleteAll(reqCtx)
					if err != nil {
						handleError(reqCtx, w, err, "Failed to delete users")
						return
					}
				} else {
					err := service.Delete(reqCtx, id)
					if err != nil {
						handleError(reqCtx, w, err, "Failed to delete user")
						return
					}
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}
		})
}

func handleError(ctx context.Context, w http.ResponseWriter, err error, message string) {
	log.WithCtx(ctx).Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

func readUser(r *http.Request) (*User, error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var usr User
	err = json.Unmarshal(b, &usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}
