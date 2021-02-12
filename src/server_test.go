package main

import (
	artifact "artifact"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func TestMain(m *testing.M) {
	setupAll()
	code := m.Run()
	teardownAll()
	os.Exit(code)
}

func setupAll() {
	fmt.Println("setup")

}

func teardownAll() {
	fmt.Println("teardown")

}
func TestGetArtifact(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	artifactHandler := artifact.New(afero.NewMemMapFs())
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(artifactHandler.HandleArtifacts)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `Get method for artifact not yet setup`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	fmt.Print("get artifact")

}

func TestPostArtifact(t *testing.T) {
	// Setup filesystem to contain image that we want to send
	// appFS := afero.NewMemMapFs()
	expectedFileName, expectedFileContent := "fileToUpload", "file content"
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, expectedFileName, []byte(expectedFileContent), 0644)
	// appFS.MkdirAll("artifacts", 0755)
	fmt.Println("at the artifact post req creation")
	req := generatePostRequest(appFS, "artifact", "fileToUpload", "file")
	fmt.Println("after the artifact post req creation")
	artifactHandler := artifact.New(appFS)
	root, err := appFS.Open("/")
	if err != nil {
		fmt.Println("fuck")

	}
	fmt.Println("Below is fs")
	fmt.Println(root.Readdirnames(10))
	handler := http.HandlerFunc(artifactHandler.HandleArtifacts)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Send the file as its expected from server
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	fileContent, err := afero.ReadFile(appFS, "../artifacts/"+expectedFileName)

	if string(fileContent) != expectedFileContent {
		fmt.Print("They're equal")
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(fileContent), expectedFileContent)
	}

}

func setupFS() {

}

func generatePostRequest(fs afero.Fs, url string, filename string, filetype string) *http.Request {
	file, err := fs.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(filetype, filepath.Base(file.Name()))

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	fmt.Println("file found and added to request body")

	return request
}
