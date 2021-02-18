package artifact

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/afero"
)

// ArtifactModel represent the model to imitate an artifact (file/mp4 etc)
type ArtifactModel struct {
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
		fmt.Fprintf(w, "Get method for artifact not yet setup")
	case "POST":
		// Check for valid JSON Body
		if r.Body == nil {
			log.Println("Received an empty body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Read body
		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			log.Println(readErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Parse JSON
		artif := ArtifactModel{}
		jsonErr := json.Unmarshal(body, &artif)
		if jsonErr != nil {

			log.Println("Received Invalid JSON")
			log.Println(jsonErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Save artifact to fs
		fsErr := a.saveToFs(artif)
		if fsErr != nil {
			log.Println("Error saving to filesystem")
			log.Println(fsErr)
			w.WriteHeader(http.StatusBadRequest)
			return

		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

// Method for getting an artifact model and storing it to fs
func (a Artifact) saveToFs(entry ArtifactModel) error {
	unbased, err := base64.StdEncoding.DecodeString(entry.File)
	if err != nil {
		return err
	}

	r := bytes.NewReader(unbased)
	// Will need to add a factory for handling different file types. Leaving as png for pr
	im, err := jpeg.Decode(r)
	if err != nil {
		return err
	}

	f, err := a.fileSystem.OpenFile("containerFiles/artifacts/"+genFileName(".jpeg"), os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	err = jpeg.Encode(f, im, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}
	log.Println("Saved filed to fs")
	return nil
}

func genFileName(ext string) string {
	currentTime := time.Now().UTC()

	return currentTime.Format("2006-01-02 15:04:05.000000000") + ext
}
