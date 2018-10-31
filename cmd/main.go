package main

import (
	"context"
	"ctco-dev/go-api-template/internal/app"
	"ctco-dev/go-api-template/internal/data/user"
	"ctco-dev/go-api-template/internal/log"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func main() {

	rootCtx := context.Background()

	var env app.Specification
	err := envconfig.Process("", &env)
	if err != nil {
		log.WithCtx(rootCtx).Panicf("env vars error: '%v'", err)
	}

	addMongo(rootCtx)

	someApp := app.NewApp(env)

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {

			reqID := uuid.NewV4().String()[0:8]
			reqCtx := log.NewContext(rootCtx, logrus.Fields{"reqID": reqID})
			reqCtx, cancel := context.WithTimeout(reqCtx, time.Second*10)
			defer cancel()

			resp, err := someApp.DoSomething(reqCtx)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.WithCtx(reqCtx).Error(err)
				w.Write([]byte("Cant't get new joke"))
				return
			}

			bytes, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Can't encode response"))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(bytes)

		})

	log.WithCtx(rootCtx).Infof("Server is running at: http://localhost:%d", env.Port)
	addr := fmt.Sprintf(":%d", env.Port)
	log.WithCtx(rootCtx).Fatal(http.ListenAndServe(addr, nil))
}

func addMongo(ctx context.Context) {
	var env user.Specification
	err := envconfig.Process("Users", &env)
	if err != nil {
		log.WithCtx(ctx).Panicf("env vars error: '%v'", err)
	}

	if env.Enabled {
		launchMongoEndpoint(ctx, env)
	}
}

func launchMongoEndpoint(ctx context.Context, env user.Specification) {
	client, err := user.NewMongoClient(env.MongodbHost)
	if err != nil {
		log.WithCtx(ctx).Panicf("Can't initialize mongodb", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.WithCtx(ctx).Panicf("Can't connect to mongodb", err)
	}

	service := user.NewService(user.NewMongoRepository(client, env.MongodbDatabase, env.MongodbCollectionUsers))
	user.LaunchController(ctx, env.Route, service)
	log.WithCtx(ctx).Info("Users service launched at ", env.Route)
}
