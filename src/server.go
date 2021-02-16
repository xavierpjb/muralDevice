package main

import (
	"net/http"

	artifact "artifact"

	"github.com/spf13/afero"
)

func main() {
	artifactHandler := artifact.New(afero.NewOsFs())

	http.HandleFunc("/artifacts", artifactHandler.HandleArtifacts)
	http.HandleFunc("/", getterPoster)

	http.ListenAndServe(":8090", nil)
}

func getterPoster(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "artifact.html")
}
