package artifact

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/afero"
)

// ArtifcatModel represent the model to imitate an artifact (file/mp4 etc)
type ArtifactModel struct {
	Name string
	File string
	Type string
}

// Artifact base constructor for construtor to fs
type Artifact struct {
	fileSystem afero.Fs
}

// New instantiates an artifact with passed in fs
func New(fileSystem afero.Fs) Artifact {
	a := Artifact{fileSystem}
	return a
}

// HandleArtifacts is used for processing the artifact upload endpoint
func (a Artifact) HandleArtifacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Get method for artifact not yet setupss")
	case "POST":
		if r.Body != nil {
			defer r.Body.Close()
		}

		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		artif := ArtifactModel{}
		jsonErr := json.Unmarshal(body, &artif)
		if jsonErr != nil {
			log.Fatal(jsonErr)

		}
		a.saveToFs(artif)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

// Method for getting an artifact model and storing it to fs
func (a Artifact) saveToFs(entry ArtifactModel) {
	unbased, err := base64.StdEncoding.DecodeString(entry.File)
	if err != nil {
		panic("Cannot decode b64")
	}

	r := bytes.NewReader(unbased)
	// Will need to add a factory for handling different file types. Leaving as png for pr
	im, err := png.Decode(r)
	if err != nil {
		panic("Bad png")
	}

	f, err := a.fileSystem.OpenFile("../artifacts/"+entry.Name, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		panic("Cannot open file")
	}

	png.Encode(f, im)
}
