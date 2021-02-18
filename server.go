package main

import (
	"log"
	"net/http"

	"muraldevice/artifact"
	"muraldevice/mural"

	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	softJSON, err := fs.Open("containerFiles/software.json")
	if err != nil {
		log.Panic(err)
	}
	muralHandler := mural.New(softJSON)
	artifactHandler := artifact.New(fs)

	http.HandleFunc("/artifact", artifactHandler.HandleArtifacts)
	http.HandleFunc("/muralInfo", muralHandler.GetSoftwareSummary)
	http.HandleFunc("/", getterPoster)

	http.ListenAndServe(":42069", nil)
}

func getterPoster(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "artifact.html")
}
