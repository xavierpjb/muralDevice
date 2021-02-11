package main

import (
	artifact "artifact"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	fmt.Print("setup")

}

func teardownAll() {
	fmt.Print("teardown:")

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

	// Send the file as its expected from server

	// Compare the saved file from server to the saved file from setup

	// Cleanup file (both from server and setup)

	fmt.Print("post artifactfuuuuuuck")
}

func setupFS() {

}
