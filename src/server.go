package main

// Create a new package which can proved artifact methods
// Package no longer imported properly after removing go.mod file
// checkout golang.org/doc/code.html for further info on packages
import (
	"fmt"
	"net/http"

	arti "artifact"

	"github.com/spf13/afero"
)

func main() {
	artifactHandler := arti.New(afero.NewOsFs())

	http.HandleFunc("/artifacts", artifactHandler.HandleArtifacts)
	http.HandleFunc("/", getterPoster)

	http.ListenAndServe(":8090", nil)
}

func getterPoster(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get poster called")
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Println("get poster called in get method")

		http.ServeFile(w, r, "artifact.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		address := r.FormValue("address")
		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Address = %s\n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
