package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/polisko/gitcommits"
)

var wait time.Duration = time.Second * 30
var authToken string

func main() {
	var ok bool
	// check, if auth token is there
	if authToken, ok = os.LookupEnv("AUTH_TOKEN"); !ok {
		log.Fatalf("Server needs environment variable %s to be set", "AUTH_TOKEN")
	}

	r := mux.NewRouter()
	// Add your routes as needed

	r.HandleFunc("/{owner}/{repo}/{branch}/{commit}", getCommits)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	// wait for signal
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)

}

func getCommits(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*3000)
	defer cancel()
	cli, err := gitcommits.NewGitCommits(authToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating git client %s", err), http.StatusInternalServerError)
		return
	}
	cli.DefaultBranch = params["branch"]
	cli.DefaultOwner = params["owner"]
	cli.DefaultRepo = params["repo"]
	c, err := cli.FindCommitWithCtx(ctx, params["commit"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ctx.Err() != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	list, err := cli.ListCommitsWithCtx(ctx, *c)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s getting list of commits %s", err, params["commit"]), http.StatusBadRequest)
		return
	}
	if ctx.Err() != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(list)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling to json %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(res))

}
