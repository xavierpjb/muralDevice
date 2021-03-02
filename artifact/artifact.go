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
	"strconv"
	"time"

	"github.com/spf13/afero"
)

// ArtifactModel represent the model to imitate an artifact (file/mp4 etc)
type ArtifactModel struct {
	Username string
	File     string
	Type     string
}

// IsPersistable checks that the properties need to persist an artifact are all present
func (a ArtifactModel) IsPersistable() bool {
	return a.File != "" && a.Type != "" && a.Username != ""
}

// Artifact base constructor for construtor to fs
type Artifact struct {
	fileSystem                afero.Fs
	artifactRepositoryHandler IRepositoryHandler
}

// New instantiates an artifact with passed in fs
func New(fileSystem afero.Fs, arh IRepositoryHandler) Artifact {
	a := Artifact{fileSystem, arh}
	return a
}

// HandleArtifacts is used for processing the artifact upload endpoint
func (a Artifact) HandleArtifacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("Artifacts requested")
		var isSuccesful bool
		var pageInt int64
		pageInt, isSuccesful = getIntParam(r, "page")
		if !isSuccesful || pageInt < 1 {
			log.Println("Defaulting page to 1")
			pageInt = 1
		}

		var perPageInt int64
		perPageInt, isSuccesful = getIntParam(r, "perPage")
		if !isSuccesful || perPageInt < 1 {
			log.Println("Defaulting perPage to 5")
			perPageInt = 5
		}

		var entries []RepositoryModel
		log.Printf("Fetching artifacts with params page:%d, perPage:%d\n", pageInt, perPageInt)
		entries = a.artifactRepositoryHandler.RetrieveList(pageInt, perPageInt)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entries)
		log.Println("Fulfilled artifact request")

	case "POST":
		log.Println("Artifact post requested")
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

		if artif.IsPersistable() {
			missingParams := "Params missing from request body. Should include username, file, filetype"
			log.Println(missingParams)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, missingParams)
			return
		}

		// Save artifact to fs
		fileURL, fileType, fsErr := a.saveToFs(artif)
		if fsErr != nil {
			log.Println("Error saving to filesystem")
			log.Println(fsErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		artifPersisted := RepositoryModel{URL: fileURL, FileType: fileType, UploadDateTime: time.Now(), Username: artif.Username}
		a.artifactRepositoryHandler.Create(artifPersisted)
		log.Println("Artifact Post request fulfilled")

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func getIntParam(r *http.Request, param string) (int64, bool) {
	paramVal := r.URL.Query()[param]
	var err error
	var paramInt int64
	var isSuccesful bool
	//Turn check for query param into a func
	if len(paramVal) > 0 {
		paramInt, err = strconv.ParseInt(paramVal[0], 10, 64)
		if err != nil {
			log.Println("Could not process page param")
		} else {
			isSuccesful = true
		}
	} else {
		log.Println("page query param not found")
	}
	return paramInt, isSuccesful

}

// Method for getting an artifact model and storing it to fs
func (a Artifact) saveToFs(entry ArtifactModel) (string, string, error) {
	unbased, err := base64.StdEncoding.DecodeString(entry.File)
	if err != nil {
		return "", "", err
	}

	r := bytes.NewReader(unbased)
	// Will need to add a factory for handling different file types. Leaving as png for pr
	im, err := jpeg.Decode(r)
	if err != nil {
		return "", "", err
	}
	fileType := ".jpeg"
	fileURL := genFileName(fileType)

	f, err := a.fileSystem.OpenFile("containerFiles/artifacts/"+fileURL, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return "", "", err
	}

	err = jpeg.Encode(f, im, &jpeg.Options{Quality: 100})
	if err != nil {
		return "", "", err
	}
	log.Println("Saved filed to fs")
	return "/image?source=" + fileURL, fileType, nil
}

func genFileName(ext string) string {
	currentTime := time.Now().UTC()

	return currentTime.Format("2006-01-02 15:04:05.000000000") + ext
}
