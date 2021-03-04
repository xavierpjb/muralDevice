package update

import (
	"fmt"
	"io/ioutil"
	"net/http"

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
	}

}
