package update

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/afero"
)

// Updater struct to handle files to be used for updating containers
type Updater struct {
	fileSystem afero.Fs
}

// New instantiates updater
func New(fileSystem afero.Fs) Updater {
	a := Updater{fileSystem}
	return a
}

// HandleUpdate is used for the endpoint
func (u Updater) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			fmt.Println(status)
		}
	}()
	switch r.Method {
	case "POST":
		fmt.Println("Update endpoint called, Parsing file")
		r.ParseMultipartForm(32 << 20)
		fmt.Println(r.MultipartForm.File)
		cont := r.MultipartForm.Value["filer"]

		fmt.Println("Writing tar file to fs")
		err = afero.WriteFile(u.fileSystem, "containerFiles/mural_dev.tar.gz", []byte(cont[0]), 0666)
		if err != nil {
			fmt.Println("There was an issue saving to fs")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("Writing update instruction so system pipe")
		stdout, _ := u.fileSystem.OpenFile("containerFiles/update", os.O_WRONLY, 0600)
		stdout.Write([]byte("cd /home/ubuntu/ && touch thisIsFromGo && docker-compose down && cd containerFiles && (docker load < mural_dev.tar.gz) && cd .. && docker-compose up -d"))
		fmt.Println("Done writing, containers will start shutting down")
		stdout.Close()
	}

}
