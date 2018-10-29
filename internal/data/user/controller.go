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

// Specification holds parameters for User service
type Specification struct {
	UseMongoDB             bool   `split_words:"true" required:"true" default:"true"`
	MongoDBUrl             string `split_words:"true" required:"true" default:"mongodb://localhost:27017"`
	MongoDBDatabase        string `split_words:"true" required:"true" default:"db"`
	MongoDBCollectionUsers string `split_words:"true" required:"true" default:"users"`
	UsersRoute             string `split_words:"true" required:"true" default:"/users"`
}

//LaunchController registers a new REST endpoint for User operations
func LaunchController(ctx context.Context, pattern string, s Service) {
	http.HandleFunc(
		pattern,
		func(w http.ResponseWriter, r *http.Request) {
			reqID := uuid.NewV4().String()[0:8]
			reqCtx := log.NewContext(ctx, logrus.Fields{"reqID": reqID})
			reqCtx, cancel := context.WithTimeout(reqCtx, time.Second*10)
			defer cancel()

			if r.Method == "GET" {
				load(reqCtx, s, w, r)
			} else if r.Method == "POST" {
				save(reqCtx, s, w, r)
			} else if r.Method == "DELETE" {
				delete(reqCtx, s, w, r)
			}
		})
}

func load(ctx context.Context, s Service, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var users []*User
	var err error
	if id == "" {
		users, err = s.All(ctx)
		if err != nil {
			handleError(ctx, w, err, "Failed to load")
			return
		}
	} else {
		usr, err := s.Find(ctx, id)
		if err != nil {
			handleError(ctx, w, err, "Failed to load")
			return
		}
		users = append(users, usr)
	}

	bytes, err := json.Marshal(users)
	if err != nil {
		handleError(ctx, w, err, "Failed to parse users")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func save(ctx context.Context, s Service, w http.ResponseWriter, r *http.Request) {
	usr, err := readUser(r)
	if err != nil {
		handleError(ctx, w, err, "Unable to parse body")
		return
	}

	id, err := s.Save(ctx, usr)
	if err != nil {
		handleError(ctx, w, err, "Failed to save user")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

func delete(ctx context.Context, s Service, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var err error
	if id == "" {
		err = s.DeleteAll(ctx)
	} else {
		err = s.Delete(ctx, id)
	}
	if err != nil {
		handleError(ctx, w, err, "Failed to delete")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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

	var usr User
	err = json.Unmarshal(b, &usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}
