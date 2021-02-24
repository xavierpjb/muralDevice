package artifact

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
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
func getMockDB(t *testing.T) *MockIArtifactRepositoryHandler {
	mockCtrl := gomock.NewController(t)
	mockObj := NewMockIArtifactRepositoryHandler(mockCtrl)
	return mockObj

}
func TestGetArtifact(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	mockArtiDB := getMockDB(t)
	mockArtiDB.EXPECT().RetrieveList()
	artifactHandler := New(afero.NewMemMapFs(), mockArtiDB)
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

func TestValidPostArtifact(t *testing.T) {
	appFS := afero.NewMemMapFs()
	mockArtiDB := getMockDB(t)
	mockArtiDB.EXPECT().Create(gomock.Any())
	artifactHandler := New(appFS, mockArtiDB)
	handler := http.HandlerFunc(artifactHandler.HandleArtifacts)

	req := generatePostRequest()
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Send the file as its expected from server
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func generatePostRequest() *http.Request {
	smallJPG := "/9j/4AAQSkZJRgABAQAAAQABAAD//gAfQ29tcHJlc3NlZCBieSBqcGVnLXJlY29tcHJlc3P/2wCEAAQEBAQEBAQEBAQGBgUGBggHBwcHCAwJCQkJCQwTDA4MDA4MExEUEA8QFBEeFxUVFx4iHRsdIiolJSo0MjRERFwBBAQEBAQEBAQEBAYGBQYGCAcHBwcIDAkJCQkJDBMMDgwMDgwTERQQDxAUER4XFRUXHiIdGx0iKiUlKjQyNEREXP/CABEIAAIAAgMBIgACEQEDEQH/xAAUAAEAAAAAAAAAAAAAAAAAAAAH/9oACAEBAAAAAD7/xAAUAQEAAAAAAAAAAAAAAAAAAAAH/9oACAECEAAAAEL/xAAUAQEAAAAAAAAAAAAAAAAAAAAF/9oACAEDEAAAACf/xAAWEAEBAQAAAAAAAAAAAAAAAAABACH/2gAIAQEAAT8ADC//xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oACAECAQE/AH//xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oACAEDAQE/AH//2Q=="

	var jsonStr = []byte(`{"file":"` + smallJPG + `", "type": "type"}`)
	req, err := http.NewRequest("POST", "rand", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	return req
}

func TestInvalidPostArtifact(t *testing.T) {
	appFS := afero.NewMemMapFs()
	artifactHandler := New(appFS, getMockDB(t))
	handler := http.HandlerFunc(artifactHandler.HandleArtifacts)

	req := generateBadArtiPostRequest()
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Send the file as its expected from server
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func generateBadArtiPostRequest() *http.Request {
	smallJPG := "/9j/InvalidJPG"

	var jsonStr = []byte(`{"file":"` + smallJPG + `", "type": "type"}`)
	req, err := http.NewRequest("POST", "rand", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	return req
}

// test for invalid json
func TestInvalidJson(t *testing.T) {
	appFS := afero.NewMemMapFs()
	artifactHandler := New(appFS, getMockDB(t))
	handler := http.HandlerFunc(artifactHandler.HandleArtifacts)

	req := generateInvalidJSONPostRequest()
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Send the file as its expected from server
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	req = generateBodylessJSON()
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func generateInvalidJSONPostRequest() *http.Request {

	var jsonStr = []byte("invalid json")
	req, err := http.NewRequest("POST", "rand", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	return req

}

func generateBodylessJSON() *http.Request {
	req, err := http.NewRequest("POST", "rand", nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	return req

}

// test for invalid b64

// test for invalid jpeg
