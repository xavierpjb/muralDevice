package artifact

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/afero"
)

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
		body, err := getJSONBody(r)
		if err != nil {
			log.Println("we got an error")
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get Json from body
		artif, err := unmarshalBody(body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !artif.IsPersistable() {
			missingParams := "Params missing from request body. Should include username, file, filetype and datetime"
			log.Println(missingParams)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, missingParams)
			return
		}

		// Save artifact to fs
		fileURL, fileType, fsErr := a.saveToFs(*artif)
		if fsErr != nil {
			log.Println("Error saving to filesystem")
			log.Println(fsErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		artifPersisted := RepositoryModel{URL: fileURL, FileType: fileType, UploadDateTime: artif.UploadDateTime.UTC(), Username: artif.Username}
		a.artifactRepositoryHandler.Create(artifPersisted)
		log.Println("Artifact Post request fulfilled")

	case "DELETE":
		log.Println("Delete called")
		// Check for valid JSON Body
		body, err := getJSONBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		artif, err := unmarshalDeleteJSON(body)

		if err != nil {
			log.Println(err)
			return
		}
		if !artif.IsDeleteable() {
			missingParams := "Params missing from request body. Should include url and username"
			log.Println(missingParams)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, missingParams)
			return
		}

		a.artifactRepositoryHandler.Delete(*artif)
		filename := strings.Replace(artif.URL, "/image?source=", "", 1)
		fmt.Println("file to delete is")
		fmt.Println(filename)

		a.deleteFromFS(filename)
		log.Println("Artifact delete request fulfilled")

	default:
		fmt.Fprintf(w, "Sorry, only GET, DELET and POST methods are supported.")
	}

}

func getJSONBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		fmt.Println("body has issues")
		return nil, errors.New("Received an empty body")
	}
	defer r.Body.Close()

	// Read body
	body, readErr := ioutil.ReadAll(r.Body)
	fmt.Println("we should get a read err")
	if readErr != nil {
		fmt.Println("read err")
		return nil, readErr
	}
	return body, nil
}

func unmarshalBody(body []byte) (*ArtifactModel, error) {
	var artif ArtifactModel
	jsonErr := json.Unmarshal(body, &artif)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &artif, nil
}

func unmarshalDeleteJSON(body []byte) (*DeleteModel, error) {
	var delete DeleteModel
	jsonErr := json.Unmarshal(body, &delete)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &delete, nil
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
	fileURL := genFileName(entry, fileType)

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

func (a Artifact) deleteFromFS(filename string) {
	err := a.fileSystem.Remove("containerFiles/artifacts/" + filename)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func genFileName(entry ArtifactModel, ext string) string {
	currentTime := entry.UploadDateTime.UTC()
	return currentTime.Format("2006-01-02 15:04:05.000000000") + ext
}
