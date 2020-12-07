package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/afero"
)

var AppFS = afero.NewOsFs()

func main() {
	http.HandleFunc("/artifacts", artifacts)
	http.HandleFunc("/", getterPoster)

	http.ListenAndServe(":8090", nil)
}

func artifacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Get method for artifact not yet setup")
	case "POST":
		in, header, err := r.FormFile("file")
		if err != nil {
			//handle error
			log.Fatal(err)
		}
		defer in.Close()
		//you probably want to make sure header.Filename is unique and
		// use filepath.Join to put it somewhere else.
		out, err := AppFS.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			//handle error
			fmt.Println("err in readfile")
			log.Fatal(err)
		}
		defer out.Close()
		io.Copy(out, in)
	//do other stuff
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func getterPoster(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
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
