package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/afero"
)

type Updater struct {
	fileSystem afero.Fs
}

func New(fileSystem afero.Fs) Updater {
	a := Updater{fileSystem}
	return a
}

func (u Updater) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			// http.Error(r, err.Error(), status)
			fmt.Println(status)

		}
	}()
	switch r.Method {
	case "POST":
		fmt.Println("post called")
		r.ParseMultipartForm(32 << 20)
		fmt.Println(r.MultipartForm.File)
		cont := r.MultipartForm.Value["filer"]
		_ = ioutil.WriteFile("containerFiles/mural_dev.tar.gz", []byte(cont[0]), 0666)
		stdout, _ := u.fileSystem.OpenFile("containerFiles/update", os.O_WRONLY, 0600)
		stdout.Write([]byte("cd /home/ubuntu/ && touch thisIsFromGo && docker-compose down && cd containerFiles && (docker load < mural_dev.tar.gz) && cd .. && docker-compose up -d"))
		stdout.Close()
		// stdout.Write([]byte("cd /home/ubuntu/ && docker-compose down && cd containerFiles && (docker load < mural_dev.tar.gz) && cd .. &&  docker-compoes up -d"))

	}

}
