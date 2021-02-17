package mural

import (
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
	fmt.Println("setup")

}

func teardownAll() {
	fmt.Println("teardown")

}

func TestGetMuralInfo(t *testing.T) {

	fs := afero.NewOsFs()
	softJSON, err := fs.Open("../containerFiles/software.json")
	if err != nil {

	}
	muralHandler := New(softJSON)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(muralHandler.GetSoftwareSummary)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestNotGetMuralInfo(t *testing.T) {

	fs := afero.NewOsFs()
	softJSON, err := fs.Open("../containerFiles/software.json")
	if err != nil {

	}
	muralHandler := New(softJSON)
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(muralHandler.GetSoftwareSummary)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `This endpoint only supports get`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
