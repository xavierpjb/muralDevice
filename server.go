package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"muraldevice/artifact"
	"muraldevice/imageDist"
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
	client := artifact.Dbdriver()
	arh := artifact.NewARH(client)
	artifactHandler := artifact.New(fs, arh)

	http.HandleFunc("/artifact", artifactHandler.HandleArtifacts)
	http.HandleFunc("/muralInfo", muralHandler.GetSoftwareSummary)

	http.Handle("/image", imageDist.ImageDistributor())
	http.HandleFunc("/", getterPoster)

	//clean up connection on application shutdown
	defer closeCons(context.TODO(), *client)
	cleanupOnExit(context.TODO(), *client)

	fmt.Println("Done with setup")
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
