// package uploader

// import (
// 	"os"
// 	"net/http"
// )

// type FSWriter struct {
// 	fileSystem os.
// }

//Create a struct which has a fs and operate on that filesystem

package artifact

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/afero"
)

type Artifact struct {
	fileSystem afero.Fs
}

func New(fileSystem afero.Fs) Artifact {
	a := Artifact{fileSystem}
	return a
}

func (a Artifact) HandleArtifacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Get method for artifact not yet setup")
	case "POST":
		src, header, err := r.FormFile("file")
		if err != nil {
			//handle error
			log.Fatal(err)
		}
		defer src.Close()
		//you probably want to make sure header.Filename is unique and
		// use filepath.Join to put it somewhere else.
		dst, err := a.fileSystem.OpenFile("../artifacts/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			//handle error
			fmt.Println("err in readfile")
			log.Fatal(err)
		}
		defer dst.Close()
		io.Copy(dst, src)
	//do other stuff
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}
