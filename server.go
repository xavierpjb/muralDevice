package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"muraldevice/artifact"
	"muraldevice/mural"

	"github.com/spf13/afero"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	fs := afero.NewOsFs()
	softJSON, err := fs.Open("containerFiles/software.json")
	if err != nil {
		log.Panic(err)
	}
	muralHandler := mural.New(softJSON)

	fmt.Println("will try to connect to db")
	//clean up connection on application shutdown
	// https://stackoverflow.com/questions/36432123/how-to-correctly-work-with-mongodb-session-in-go
	client, context := artifact.Dbdriver()
	artifactHandler := artifact.New(fs, *client)

	http.HandleFunc("/artifact", artifactHandler.HandleArtifacts)
	http.HandleFunc("/muralInfo", muralHandler.GetSoftwareSummary)
	http.HandleFunc("/", getterPoster)

	defer closeCons(context, *client)
	cleanupOnExit(context, *client)

	http.ListenAndServe(":42069", nil)
}

func getterPoster(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "artifact.html")
}

func cleanupOnExit(context context.Context, client mongo.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		closeCons(context, client)
		os.Exit(1)
	}()

}

func closeCons(context context.Context, client mongo.Client) {
	client.Disconnect(context)
}
